package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/vertexai/genai"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/invopop/jsonschema"

	types "github.com/movie-guru/pkg/types"
)

func GetUserProfileFlow(ctx context.Context, model ai.Model) (*genkit.Flow[*types.ProfileAgentInput, *types.UserProfileAgentOutput, struct{}], error) {

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
			InputSchema:  jsonschema.Reflect(types.ProfileAgentInput{}),
			OutputSchema: jsonschema.Reflect(types.UserProfileAgentOutput{}),
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
	userPrefFlow := genkit.DefineFlow("userPreferencesFlow", func(ctx context.Context, input *types.ProfileAgentInput) (*types.UserProfileAgentOutput, error) {
		userPrefOutput := &types.UserProfileAgentOutput{
			ModelOutputMetadata: &types.ModelOutputMetadata{
				SafetyIssue: false,
			},
			ProfileChangeRecommendations: make([]*types.ProfileChangeRecommendation, 0),
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
				userPrefOutput = &types.UserProfileAgentOutput{
					ModelOutputMetadata: &types.ModelOutputMetadata{
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
