import { gemini15Flash } from '@genkit-ai/vertexai';
import {QueryTransformFlowOutput, QueryTransformFlowInputSchema, QueryTransformFlowOutputSchema} from './queryTransformTypes'
import { QueryTransformPromptText } from './prompts';
import { ai } from './genkitConfig'

export const QueryTransformPrompt = ai.definePrompt(
    {
      name: 'queryTransformFlowPrompt',
      model: gemini15Flash,
      input: {
        schema: QueryTransformFlowInputSchema,
      },
      output: {
        format: 'json',
      },  
    }, 
   QueryTransformPromptText
)
  export const QueryTransformFlow = ai.defineFlow(
    {
      name: 'queryTransformFlow',
      inputSchema: QueryTransformFlowInputSchema,
      outputSchema: QueryTransformFlowOutputSchema
    },
    async (input) => {
      try {
        const response = await QueryTransformPrompt({ history: input.history, userMessage: input.userMessage, userProfile: input.userProfile});
        const jsonResponse =  JSON.parse(response.text);
        const output: QueryTransformFlowOutput = {
          "transformedQuery":  jsonResponse.transformedQuery,
          "userIntent": jsonResponse.userIntent,
          "modelOutputMetadata": {
            "justification": jsonResponse.justification,
            "safetyIssue": jsonResponse.safetyIssue,
          }
        }
        return output
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
  