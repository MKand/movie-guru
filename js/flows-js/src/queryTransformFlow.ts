import { defineFlow } from '@genkit-ai/flow';
import { gemini15Flash } from '@genkit-ai/vertexai';
import { defineDotprompt } from '@genkit-ai/dotprompt'
import {QueryTransformFlowInputSchema, QueryTransformFlowOutputSchema} from './queryTransformTypes'
import { QueryTransformPromptText } from './prompts';

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

        // INSTRUCTIONS:
        // 1. Create a flow called QueryTransform flow
        // 2. Call this prompt with the necessary input and get the output.
        // 3. The output should returned as type  QueryTransformFlowOutput

  