import { defineFlow } from '@genkit-ai/flow';
import { gemini15Flash } from '@genkit-ai/vertexai';
import { defineDotprompt } from '@genkit-ai/dotprompt'
import {QueryTransformFlowInputSchema, QueryTransformFlowOutputSchema} from './queryTransformTypes'
import { QueryTransformPromptText } from './prompts';
import * as z from 'zod';

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
   QueryTransformPromptText
)

// Implement the QueryTransformFlow
export const QueryTransformFlow = defineFlow(
{
  name: 'queryTransformFlow',
  inputSchema: z.string(), // what should this be?
  outputSchema: z.string(), // what should this be?
},
async (input) => {
  return "Hello World"
}
    
);
  