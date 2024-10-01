package main

import (
	"context"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/invopop/jsonschema"

	types "github.com/movie-guru/pkg/types"
)

type QueryTransformFlowOutput struct {
	TransformedQuery string           `json:"transformedQuery, omitempty"`
	Intent           types.USERINTENT `json:"userIntent, omitempty"`
	*types.ModelOutputMetadata
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
	// Printed here to make sure the prompt variable is used, or the Golang compiler will complain.
	log.Println(queryTransformPrompt)
	if err != nil {
		return nil, err
	}
	// Define a simple flow that prompts an LLM to generate menu suggestions.
	queryTransformFlow := genkit.DefineFlow("queryTransformFlow", func(ctx context.Context, input *QueryTransformFlowInput) (*QueryTransformFlowOutput, error) {
		// Default output
		queryTransformFlowOutput := &QueryTransformFlowOutput{
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

		return queryTransformFlowOutput, nil
	})
	return queryTransformFlow, nil
}
