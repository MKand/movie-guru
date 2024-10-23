import { defineFlow } from '@genkit-ai/flow';
import { gemini15Flash } from '@genkit-ai/vertexai';
import { defineDotprompt, promptRef } from '@genkit-ai/dotprompt'
import {QueryTransformFlowInputSchema, QueryTransformFlowOutputSchema} from './queryTransformTypes'
import { QueryTransformPromptText } from './prompts';

// Defining the dotPrompt
export const QueryTransformPrompt = defineDotprompt(
  {
    name: 'QueryTransformFlowPrompt',
    model: gemini15Flash,
    input: {
      schema: QueryTransformFlowInputSchema,
    },
    output: {
      format: 'json',
      schema: QueryTransformFlowOutputSchema,
    },  
  }, 
 QueryTransformPromptText
)

export const QueryTransformFlow = defineFlow(
  {
    name: 'QueryTransformFlow',
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