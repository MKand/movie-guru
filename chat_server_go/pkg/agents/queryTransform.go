package agents

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

func GetQueryTransformFlow(ctx context.Context, model ai.Model) (*genkit.Flow[*types.QueryTransformInput, *types.QueryTransformOutput, struct{}], error) {

	queryTransformPrompt, err := dotprompt.Define("queryTransform",
		`You are a search query refinement expert. Your goal is NOT to answer the user's question directly, but to craft the most effective raw query for a vector search engine to retrieve information relevant to a user's current request, taking into account their conversation history and known preferences.
                    Instructions:

                    1. Analyze Conversation History: Carefully examine the provided conversation history to understand the context and main topics the user is interested in. Identify the user's most recent question or request as the primary focus for the search query.

                    2. Incorporate the user's profle when relevant:
                    * Strong Likes: If the likes in the user's profile align directly with the current query, integrate them into the query to enhance results.
                    * Strong Dislikes: Only incorporate dislikes into the query if they directly conflict with or narrow down the user's request.
                    * Irrelevant like or disklike: If a like or dislike doesn't relate to the current query, exclude it from the search.

                    3. Prioritize User Intent: The user's current request should be the core of the search query. Don't let the user's profile overshadow the main topic the user is seeking information about.

                    4. Concise and Specific: Keep the query concise and specific to maximize the relevance of search results. Avoid adding unnecessary details or overly broad terms.					

					Here is the user profile. This expresses their long-term likes and dislikes:
                    {{userProfile}} 
					If the user is looking for movies to watch and isn't specific then you can incorporate information from their profile into the query. 

					This is the history of the conversation with the user so far to learn about the context of the conversation:
					{{history}} 
			
					This is the last message the user sent. Use this to understand the user's intent.:
					{{userMessage}}
					Use it to understand if the user is asking a question, or clarifying their request or responding to a question posed by the agent. If the user is doing neither, (eg: greeting you, or saying bye or just acknowledging the response) then set the variable userIntent accordingly. 

					 Your response should include the following main parts:

					justification: Justification for your answer
					transformedQuery: Your interpretation of the user's message
					userIntent: The userIntent can be GREET, END_CONVERSATION, REQUEST, RESPONSE, ACKNOWLEDGE or UNCLEAR. If the user is acknowleding what the agent last said, or remarking on it by saying something like OK or cool, then set the userIntent to ACKNOWLEDGE.
					If the user is asking something (eg: about a movie, or a clarification), then the userIntent is REQUEST. If the user is responding to the agent's question, then the userIntent is RESPONSE
					`,

		dotprompt.Config{
			Model:        model,
			InputSchema:  jsonschema.Reflect(types.QueryTransformInput{}),
			OutputSchema: jsonschema.Reflect(types.QueryTransformOutput{}),
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
	userPrefFlow := genkit.DefineFlow("QueryTransformFlow", func(ctx context.Context, input *types.QueryTransformInput) (*types.QueryTransformOutput, error) {
		transformedQuery := &types.QueryTransformOutput{
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
				transformedQuery = &types.QueryTransformOutput{
					ModelOutputMetadata: &types.ModelOutputMetadata{
						SafetyIssue: true,
					},
					TransformedQuery: "",
				}
				return transformedQuery, nil

			} else {
				return nil, err

			}
		}
		t := resp.Text()
		err = json.Unmarshal([]byte(t), &transformedQuery)
		if err != nil {
			return nil, err
		}

		return transformedQuery, nil
	})
	return userPrefFlow, nil
}
