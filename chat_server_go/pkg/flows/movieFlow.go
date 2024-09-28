package flows

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

func GetMovieFlow(ctx context.Context, model ai.Model, prompt string) (*genkit.Flow[*types.MovieFlowInput, *types.MovieFlowOutput, struct{}], error) {
	movieAgentPrompt, err := dotprompt.Define("movieFlow",
		prompt,

		dotprompt.Config{
			Model:        model,
			InputSchema:  jsonschema.Reflect(types.MovieFlowInput{}),
			OutputSchema: jsonschema.Reflect(types.MovieFlowOutput{}),
			OutputFormat: ai.OutputFormatText,
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
		func(ctx context.Context, input *types.MovieFlowInput) (*types.MovieFlowOutput, error) {
			var movieFlowOutput *types.MovieFlowOutput
			resp, err := movieAgentPrompt.Generate(ctx,
				&dotprompt.PromptRequest{
					Variables: input,
				},
				nil,
			)
			if err != nil {
				if blockedErr, ok := err.(*genai.BlockedError); ok {
					fmt.Println("Request was blocked:", blockedErr)
					movieFlowOutput = &types.MovieFlowOutput{
						ModelOutputMetadata: &types.ModelOutputMetadata{
							SafetyIssue: true,
						},
						RelevantMoviesTitles: make([]*types.RelevantMovie, 0),
						WrongQuery:           false,
					}
					return movieFlowOutput, nil

				} else {
					return nil, err

				}
			}
			t := resp.Text()
			parsedJson, err := makeJsonMarshallable(t)
			if err != nil {
				if len(parsedJson) > 0 {
					log.Printf("Didn't get json resp from movie agent. %s", t)
				}
			}
			err = json.Unmarshal([]byte(parsedJson), &movieFlowOutput)
			if err != nil {
				return nil, err
			}
			return movieFlowOutput, nil
		},
	)
	return movieFlow, nil
}

func extractText(jsonText string) string {
	return ""
}
