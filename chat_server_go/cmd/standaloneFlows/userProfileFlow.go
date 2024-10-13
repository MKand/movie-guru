package main

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

type UserProfileFlowInput struct {
	Query        string `json:"query"`
	AgentMessage string `json:"agentMessage"`
}

type UserProfileFlowOutput struct {
	ProfileChangeRecommendations []*types.ProfileChangeRecommendation `json:"profileChangeRecommendations"`
	Justification                string                               `json:"justification,omitempty"`
}

func NewUserProfileFlowOuput() *UserProfileFlowOutput {
	return &UserProfileFlowOutput{
		ProfileChangeRecommendations: make([]*types.ProfileChangeRecommendation, 5),
	}
}

func GetUserProfileFlow(ctx context.Context, model ai.Model, prompt string) (*genkit.Flow[*UserProfileFlowInput, *UserProfileFlowOutput, struct{}], error) {

	prefPrompt, err := dotprompt.Define("userProfileFlow",
		prompt,

		dotprompt.Config{
			Model:        model,
			InputSchema:  jsonschema.Reflect(UserProfileFlowInput{}),
			OutputSchema: jsonschema.Reflect(UserProfileFlowOutput{}),
			OutputFormat: ai.OutputFormatJSON,
			GenerationConfig: &ai.GenerationCommonConfig{
				Temperature: 0.5,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	userProfileFlow := genkit.DefineFlow("userProfileFlow", func(ctx context.Context, input *UserProfileFlowInput) (*UserProfileFlowOutput, error) {
		userProfileFlowOutput := &UserProfileFlowOutput{
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
				userProfileFlowOutput = &UserProfileFlowOutput{}
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
