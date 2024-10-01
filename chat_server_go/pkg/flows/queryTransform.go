package flows

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/vertexai/genai"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/invopop/jsonschema"

	types "github.com/movie-guru/pkg/types"
)

func GetQueryTransformFlow(ctx context.Context, model ai.Model, prompt string) (*genkit.Flow[*types.QueryTransformFlowInput, *types.QueryTransformFlowOutput, struct{}], error) {

	queryTransformPrompt, err := dotprompt.Define("queryTransformFlow",
		prompt,

		dotprompt.Config{
			Model:        model,
			InputSchema:  jsonschema.Reflect(types.QueryTransformFlowInput{}),
			OutputSchema: jsonschema.Reflect(types.QueryTransformFlowOutput{}),
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
	queryTransformFlow := genkit.DefineFlow("queryTransformFlow", func(ctx context.Context, input *types.QueryTransformFlowInput) (*types.QueryTransformFlowOutput, error) {
		// Default output
		queryTransformFlowOutput := &types.QueryTransformFlowOutput{
			ModelOutputMetadata: &types.ModelOutputMetadata{
				SafetyIssue:   false,
				Justification: "",
			},
			TransformedQuery: "",
			Intent:           types.USERINTENT(types.UNCLEAR),
		}

		// INSTRUCTIONS:
		// 1. Call this prompt with the necessary input and get the output.
		// 2. The output should then be tranformed into the type  QueryTransformFlowOutput and stored in the variable queryTransformFlowOutput
		// 3. Handle any errors that may arise.

		resp, err := queryTransformPrompt.Generate(ctx,
			&dotprompt.PromptRequest{
				Variables: input,
			},
			nil,
		)
		if err != nil {
			if blockedErr, ok := err.(*genai.BlockedError); ok {
				log.Println("Request was blocked:", blockedErr)
				queryTransformFlowOutput = &types.QueryTransformFlowOutput{
					ModelOutputMetadata: &types.ModelOutputMetadata{
						SafetyIssue: true,
					},
					TransformedQuery: "",
				}
				return queryTransformFlowOutput, nil

			} else {
				return nil, err

			}
		}
		t := resp.Text()
		err = json.Unmarshal([]byte(t), &queryTransformFlowOutput)
		if err != nil {
			return nil, err
		}

		return queryTransformFlowOutput, nil
	})
	return queryTransformFlow, nil
}
