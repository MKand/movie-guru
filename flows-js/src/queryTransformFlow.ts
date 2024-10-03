import { defineFlow } from '@genkit-ai/flow';
import { gemini15Flash } from '@genkit-ai/vertexai';
import { defineDotprompt } from '@genkit-ai/dotprompt'
import {QueryTransformFlowInputSchema, QueryTransformFlowOutputSchema} from './queryTransformTypes'


export const QueryTransformPrompt = defineDotprompt(
    {
      name: 'queryTransformFlow',
      model: gemini15Flash,
      input: {
        schema: QueryTransformFlowInputSchema,
      },
      output: {
        format: 'json',
        schema: QueryTransformFlowOutputSchema,
      },  
    }, 
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
		`
)
  export const QueryTransformFlow = defineFlow(
    {
      name: 'queryTransformFlow',
      inputSchema: QueryTransformFlowInputSchema,
      outputSchema: QueryTransformFlowOutputSchema
    },
    async (input) => {
      try {
        const response = await QueryTransformPrompt.generate({ input: input });
        console.log(response.output(0))
        return response.output(0);
      } catch (error) {
        console.error("Error generating response:", error);
        return { 
          transformedQuery: "",
          userIntent: 'UNCLEAR',
          justification: ""
         }; 
      }
    }
  );
  