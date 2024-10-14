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

type MovieFlowInput struct {
	History          []*types.SimpleMessage `json:"history"`
	UserPreferences  *types.UserProfile     `json:"userPreferences"`
	ContextDocuments []*types.MovieContext  `json:"contextDocuments"`
	UserMessage      string                 `json:"userMessage"`
}

type MovieFlowOutput struct {
	Answer               string                 `json:"answer"`
	RelevantMoviesTitles []*types.RelevantMovie `json:"relevantMovies"`
	WrongQuery           bool                   `json:"wrongQuery,omitempty"`
	Justification        string                 `json:"justification,omitempty"`
}

func GetMovieFlow(ctx context.Context, model ai.Model, prompt string) (*genkit.Flow[*MovieFlowInput, *MovieFlowOutput, struct{}], error) {
	movieAgentPrompt, err := dotprompt.Define("movieFlow",
		prompt,

		dotprompt.Config{
			Model:        model,
			InputSchema:  jsonschema.Reflect(MovieFlowInput{}),
			OutputSchema: jsonschema.Reflect(MovieFlowOutput{}),
			OutputFormat: ai.OutputFormatJSON,
			GenerationConfig: &ai.GenerationCommonConfig{
				Temperature: 0.5,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	movieFlow := genkit.DefineFlow(
		"movieQAFlow",
		func(ctx context.Context, input *MovieFlowInput) (*MovieFlowOutput, error) {
			var movieFlowOutput *MovieFlowOutput
			resp, err := movieAgentPrompt.Generate(ctx,
				&dotprompt.PromptRequest{
					Variables: input,
				},
				nil,
			)
			if err != nil {
				if blockedErr, ok := err.(*genai.BlockedError); ok {
					fmt.Println("Request was blocked:", blockedErr)
					movieFlowOutput = &MovieFlowOutput{
						RelevantMoviesTitles: make([]*types.RelevantMovie, 0),
						WrongQuery:           false,
					}
					return movieFlowOutput, nil

				} else {
					return nil, err

				}
			}
			t := resp.Text()
			err = json.Unmarshal([]byte(t), &movieFlowOutput)
			if err != nil {
				return nil, err
			}
			return movieFlowOutput, nil
		},
	)
	return movieFlow, nil
}
