package flows

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/invopop/jsonschema"

	types "github.com/movie-guru/pkg/types"
)

func GetUserProfileFlow(ctx context.Context, model ai.Model) (*genkit.Flow[*types.UserProfileFlowInput, *types.UserProfileFlowOutput, struct{}], error) {

	prefPrompt, err := dotprompt.Define("userProfileFlow",
		`
	You are a user's movie profiling expert focused on uncovering users' enduring likes and dislikes. 
     Your task is to analyze the user message and extract ONLY strongly expressed, enduring likes and dislikes related to movies.
     Once you extract any new likes or dislikes from the current query respond with the items you extracted with:
		  1. the category (ACTOR, DIRECTOR, GENRE, OTHER)
		  2. the item value
		  3. your reason behind the choice
		  4. the sentiment of the user has about the item (POSITIVE, NEGATIVE).
		
      Guidelines:
      1. Strong likes and dislikes Only: Add or Remove ONLY items expressed with strong language indicating long-term enjoyment or aversion (e.g., "love," "hate," "can't stand,", "always enjoy"). Ignore mild or neutral items (e.g., "like,", "okay with," "fine", "in the mood for", "do not feel like").
      2. Distinguish current state of mind vs. Enduring likes and dislikes:  Be very cautious when interpreting statements. Focus only on long-term likes or dislikes while ignoring current state of mind. If the user expresses wanting to watch a specific type of movie or actor NOW, do NOT assume it's an enduring like unless they explicitly state it. For example, "I want to watch a horror movie movie with Christina Appelgate" is a current desire, NOT an enduring preference for horror movies or Christina Appelgate.
      3. Focus on Specifics:  Look for concrete details about genres, directors, actors, plots, or other movie aspects.
      4. Give an explanation as to why you made the choice.
        
		Inputs: 
		1. Optional Message 0 from agent: {{agentMessage}}
		2. Required Message 1 from user: {{query}}
		`,

		dotprompt.Config{
			Model:        model,
			InputSchema:  jsonschema.Reflect(types.UserProfileFlowInput{}),
			OutputSchema: jsonschema.Reflect(types.UserProfileFlowOutput{}),
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
	userProfileFlow := genkit.DefineFlow("userProfileFlow", func(ctx context.Context, input *types.UserProfileFlowInput) (*types.UserProfileFlowOutput, error) {
		userProfileFlowOutput := &types.UserProfileFlowOutput{
			ModelOutputMetadata: &types.ModelOutputMetadata{
				SafetyIssue: false,
			},
			ProfileChangeRecommendations: make([]*types.ProfileChangeRecommendation, 0),
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
				userProfileFlowOutput = &types.UserProfileFlowOutput{
					ModelOutputMetadata: &types.ModelOutputMetadata{
						SafetyIssue: true,
					},
				}
				return userProfileFlowOutput, nil

			} else {
				return nil, err

			}
		}
		t := resp.Text()
		err = json.Unmarshal([]byte(t), &userProfileFlowOutput)
		if err != nil {
			return nil, err
		}
		return userProfileFlowOutput, nil
	})
	return userProfileFlow, nil
}
