package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/vertexai/genai"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/invopop/jsonschema"
)

type PreferencesAgent struct {
	Model ai.Model
	Flow  *genkit.Flow[*ProfileAgentInput, *UserProfileAgentOutput, struct{}]
	DB    *sql.DB
}

func CreatePreferencesAgent(ctx context.Context, model ai.Model, db *sql.DB) (*PreferencesAgent, error) {
	flow, err := GetUserPrefFlow(ctx, model)
	if err != nil {
		return nil, err
	}
	return &PreferencesAgent{
		Model: model,
		Flow:  flow,
		DB:    db,
	}, nil
}

func (p *PreferencesAgent) Run(ctx context.Context, history *ChatHistory, user string) (*UserProfileOutput, error) {
	userProfile, err := getCurrentProfile(ctx, user, p.DB)
	userProfileOutput := &UserProfileOutput{
		UserProfile: userProfile,
		ChangesMade: false,
		ModelOutputMetadata: &ModelOutputMetadata{
			SafetyIssue:   false,
			Justification: "",
		},
	}
	if err != nil {
		return nil, err
	}
	agentMessage := ""
	if len(history.History) > 1 {
		agentMessage = history.History[len(history.History)-2].Content[0].Text
	}
	if err != nil {
		return nil, err
	}
	lastUserMessage, err := history.GetLastMessage()
	if err != nil {
		return nil, err
	}

	prefInput := ProfileAgentInput{Query: lastUserMessage, AgentMessage: agentMessage}
	resp, err := p.Flow.Run(ctx, &prefInput)
	if err != nil {
		return userProfileOutput, err
	}
	userProfileOutput.ChangesMade = resp.ChangesMade
	userProfileOutput.ModelOutputMetadata.Justification = resp.ModelOutputMetadata.Justification
	userProfileOutput.ModelOutputMetadata.SafetyIssue = resp.ModelOutputMetadata.SafetyIssue

	if resp.ChangesMade {
		updatedProfile, err := processProfileChanges(userProfile, resp.ProfileChangeRecommendations)
		if err != nil {
			return userProfileOutput, err
		}
		err = updatePreferences(ctx, updatedProfile, user, p.DB)
		if err != nil {
			return userProfileOutput, err
		}
		userProfileOutput.UserProfile = updatedProfile
	}
	return userProfileOutput, nil
}

func GetUserPrefFlow(ctx context.Context, model ai.Model) (*genkit.Flow[*ProfileAgentInput, *UserProfileAgentOutput, struct{}], error) {

	prefPrompt, err := dotprompt.Define("userPrefileAgent",
		`You are a user's movie profiling expert focused on uncovering users' enduring likes and dislikes. Analyze the user message and extract ONLY strongly expressed, enduring likes and dislikes related to movies.
         Once you extract any new likes or dislikes from the current query respond with the new profile items you extracted with the category (ACTOR, DIRECTOR, GENRE, OTHER), the item value, the reason  and the sentiment of the user has about the item (POSITIVE, NEGATIVE).
		
            Guidelines:

            1. Strong likes and dislikes Only: Add or Remove ONLY items expressed with strong language indicating long-term enjoyment or aversion (e.g., "love," "hate," "can't stand,", "always enjoy"). Ignore mild or neutral items (e.g., "like,", "okay with," "fine", "in the mood for", "do not feel like").
            2. Distinguish current state of mind vs. Enduring likes and dislikes:  Be very cautious when interpreting statements. Focus only on long-term likes or dislikes while ignoring current state of mind. If the user expresses wanting to watch a specific type of movie or actor NOW, do NOT assume it's an enduring like unless they explicitly state it. For example, "I want to watch a horror movie movie with Christina Appelgate" is a current desire, NOT an enduring preference for horror movies or Christina Appelgate.
            3. Focus on Specifics:  Look for concrete details about genres, directors, actors, plots, or other movie aspects.
            4. Exclude Vague Statements: Don't include vague statements like "good movies" or "bad movies."
			5. Do not remove or change anything in the input profile unless the user makes a statement that expresses an enduring change in the user's preference or aversion that is present in the input Profile. For example if the user's profile states that they like horror, and their statement is "I dont feel like watching a horror movie", that is not an enduring change, but is only their current state of mind. So don't update their profile based on that.
			If you do make changes in their Profile (move likes to dislikes or vice versa or delete existing items, justify this change)
			6. Be VERY conservative in interpreting likes and dislikes: If the user asks for a specific genre or actor, or plot, or director does not indicate a strong preference and therefore should not be added to the profile.
            user message: {{query}}
			The user may be responding to the agent's response {{agentMessage}} that preceeds the user message. If this present use this to inform you of the context of the conversation.

           Give an explanation as to why you made the choice.
		   Set the value of changesMade if you suggest changes in the input profile. Set to false otherwise.
`,

		dotprompt.Config{
			Model:        model,
			InputSchema:  jsonschema.Reflect(ProfileAgentInput{}),
			OutputSchema: jsonschema.Reflect(UserProfileAgentOutput{}),
			OutputFormat: ai.OutputFormatJSON,
			GenerationConfig: &ai.GenerationCommonConfig{
				Temperature: 0.5,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	// Define a simple flow that prompts an LLM to generate menu suggestions.
	userPrefFlow := genkit.DefineFlow("UserPreferencesFlow", func(ctx context.Context, input *ProfileAgentInput) (*UserProfileAgentOutput, error) {
		userPrefOutput := &UserProfileAgentOutput{
			ModelOutputMetadata: &ModelOutputMetadata{
				SafetyIssue: false,
			},
			ProfileChangeRecommendations: make([]*ProfileChangeRecommendation, 0),
			ChangesMade:                  false,
		}

		resp, err := prefPrompt.Generate(ctx,
			&dotprompt.PromptRequest{
				Variables: input,
			},
			nil,
		)
		if err != nil {
			if blockedErr, ok := err.(*genai.BlockedError); ok {
				fmt.Println("Request was blocked:", blockedErr)
				userPrefOutput = &UserProfileAgentOutput{
					ModelOutputMetadata: &ModelOutputMetadata{
						SafetyIssue: true,
					},
				}
				return userPrefOutput, nil

			} else {
				return nil, err

			}
		}
		t := resp.Text()
		err = json.Unmarshal([]byte(t), &userPrefOutput)
		if err != nil {
			return nil, err
		}
		log.Println(userPrefOutput)
		return userPrefOutput, nil
	})
	return userPrefFlow, nil
}

func getCurrentProfile(ctx context.Context, user string, db *sql.DB) (*UserProfile, error) {
	preferences := NewUserProfile()
	rows := db.QueryRowContext(ctx, `
	SELECT preferences FROM user_preferences 
	WHERE "user" = $1;`,
		user)
	var jsonData string
	err := rows.Scan(&jsonData)
	if err != nil {
		return preferences, err
	}
	err = json.Unmarshal([]byte(jsonData), &preferences)
	if err != nil {
		return preferences, err
	}
	return preferences, nil
}

func updatePreferences(ctx context.Context, newPref *UserProfile, user string, db *sql.DB) error {
	newPreferencesStr, err := json.Marshal(newPref)
	if err != nil {
		return err
	}
	query := `
        INSERT INTO user_preferences ("user", preferences)
        VALUES ($1, $2)
        ON CONFLICT ("user") DO UPDATE
        SET preferences = EXCLUDED.preferences;
    `

	// Execute the query (replace with your actual execute_query function)
	_, err = db.ExecContext(ctx, query, user, newPreferencesStr)
	if err != nil {
		return err
	}
	return nil
}

func deletePreferences(ctx context.Context, user string, db *sql.DB) error {
	query := `
		DELETE FROM user_preferences
		WHERE "user" = %1;
	`
	_, err := db.ExecContext(ctx, query, user)
	if err != nil {
		return err
	}
	return nil
}
