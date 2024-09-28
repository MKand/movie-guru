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

func GetQueryTransformFlow(ctx context.Context, model ai.Model, prompt string) (*genkit.Flow[*types.QueryTransformFlowInput, *types.QueryTransformFlowOutput, struct{}], error) {

	queryTransformPrompt, err := dotprompt.Define("queryTransform",
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
		queryTransformFlowOutput := &types.QueryTransformFlowOutput{
			ModelOutputMetadata: &types.ModelOutputMetadata{
				SafetyIssue:   false,
				Justification: "",
			},
			TransformedQuery: "",
			Intent:           types.USERINTENT(types.UNCLEAR),
		}

		resp, err := queryTransformPrompt.Generate(ctx,
			&dotprompt.PromptRequest{
				Variables: input,
			},
			nil,
		)
		if err != nil {
			if blockedErr, ok := err.(*genai.BlockedError); ok {
				fmt.Println("Request was blocked:", blockedErr)
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
