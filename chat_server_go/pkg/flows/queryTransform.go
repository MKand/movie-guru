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

func GetQueryTransformFlow(ctx context.Context, model ai.Model) (*genkit.Flow[*types.QueryTransformFlowInput, *types.QueryTransformFlowOutput, struct{}], error) {

	queryTransformPrompt, err := dotprompt.Define("queryTransformFlow",
		`
		You are a search query refinement expert regarding movies and movie related information.  Your goal is to analyse the user's intent and create a short query for a vector search engine specialised in movie related information.
		If the user's intent doesn't require a search in the database then return an empty transformedQuery. For example: if the user is greeting you, or ending the conversation.
		You should NOT attempt to answer's the user's query.
		Instructions:

		1. Analyze the conversation history to understand the context and main topics. Focus on the user's most recent request. The history may be empty.
		2.  Use the user profile when relevant:
			*   Include strong likes if they align with the query.
			*   Include strong dislikes only if they conflict with or narrow the request.
			*   Ignore irrelevant likes or dislikes.
			*  The user may have no strong likes or dislikes
		3. Prioritize the user's current request as the core of the search query.
		4. Keep the transformed query concise and specific.
		5. Only use the information in the conversation history, the user's preferences and the current request to respond. Do not use other sources of information.
		6. If the user is talking about topics unrelated to movies, return an empty transformed query and state the intent as UNCLEAR.
		7. You have absolutely no knowledge of movies.

		Here are the inputs:
		* Conversation History (this may be empty):
			{{history}}
		* UserProfile (this may be empty):
			{{userProfile}}
		* User Message:
			{{userMessage}}

		Respond with the following:

		*   a *justification* about why you created the query this way.
		*   the *transformedQuery* which is the resulting refined search query.
		*   a *userIntent*, which is one of GREET, END_CONVERSATION, REQUEST, RESPONSE, ACKNOWLEDGE, UNCLEAR
		`,

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
