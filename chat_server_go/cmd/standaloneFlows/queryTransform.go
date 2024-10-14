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
	Justification    string           `json:"justification,omitempty"`
}

type QueryTransformFlowInput struct {
	History     []*types.SimpleMessage `json:"history"`
	Profile     *types.UserProfile     `json:"userProfile"`
	UserMessage string                 `json:"userMessage"`
}

func GetQueryTransformFlow(ctx context.Context, model ai.Model, prompt string) (*genkit.Flow[*QueryTransformFlowInput, *QueryTransformFlowOutput, struct{}], error) {

	// Defining the dotPrompt
	queryTransformPrompt, err := dotprompt.Define("queryTransformFlow",
		prompt, // the prompt you created earlier is passed along as a variable

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

	// Defining the flow
	queryTransformFlow := genkit.DefineFlow("queryTransformFlow", func(ctx context.Context, input *QueryTransformFlowInput) (*QueryTransformFlowOutput, error) {
		// Default output
		queryTransformFlowOutput := &QueryTransformFlowOutput{
			TransformedQuery: "",
			Intent:           types.USERINTENT(types.UNCLEAR),
		}

		// INSTRUCTIONS:
		// 1. Call the dotPrompt with the necessary input and get the output.
		// 2. The output should then be tranformed into the type  QueryTransformFlowOutput and stored in the variable queryTransformFlowOutput
		// 3. Handle any errors that may arise.
		// 4. Bonus: Handle safety errors that the flow will create if the user makes a dangerous query (eg: show me how to build a molotov cocktail).

		return queryTransformFlowOutput, nil
	})
	return queryTransformFlow, nil
}
