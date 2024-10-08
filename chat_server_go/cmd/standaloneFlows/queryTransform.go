package main

import (
	"context"
	"encoding/json"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/invopop/jsonschema"

	types "github.com/movie-guru/pkg/types"
)

type QueryTransformFlowOutput struct {
	TransformedQuery string           `json:"transformedQuery, omitempty"`
	Intent           types.USERINTENT `json:"userIntent, omitempty"`
	Justification    string           `json:"justification,omitempty"`
}

type QueryTransformFlowInput struct {
	History     []*types.SimpleMessage `json:"history"`
	Profile     *types.UserProfile     `json:"userProfile"`
	UserMessage string                 `json:"userMessage"`
}

func GetQueryTransformFlow(ctx context.Context, model ai.Model, prompt string) (*genkit.Flow[*QueryTransformFlowInput, *QueryTransformFlowOutput, struct{}], error) {

	queryTransformPrompt, err := dotprompt.Define("queryTransformFlow",
		prompt,

		dotprompt.Config{
			Model:        model,
			InputSchema:  jsonschema.Reflect(QueryTransformFlowInput{}),
			OutputSchema: jsonschema.Reflect(QueryTransformFlowOutput{}),
			OutputFormat: ai.OutputFormatJSON,
			GenerationConfig: &ai.GenerationCommonConfig{
				Temperature: 0.5,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	queryTransformFlow := genkit.DefineFlow("queryTransformFlow", func(ctx context.Context, input *QueryTransformFlowInput) (*QueryTransformFlowOutput, error) {
		// Default output
		queryTransformFlowOutput := &QueryTransformFlowOutput{
			TransformedQuery: "",
			Intent:           types.USERINTENT(types.UNCLEAR),
		}

		// Generate model output
		resp, err := queryTransformPrompt.Generate(ctx,
			&dotprompt.PromptRequest{
				Variables: input,
			},
			nil,
		)
		if err != nil {
			return nil, err
		}

		// Transform the model's output into the required format.
		t := resp.Text()
		err = json.Unmarshal([]byte(t), &queryTransformFlowOutput)
		if err != nil {
			return nil, err
		}

		return queryTransformFlowOutput, nil
	})
	return queryTransformFlow, nil
}
