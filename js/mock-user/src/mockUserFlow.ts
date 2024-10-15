import { defineFlow } from '@genkit-ai/flow';
import { gemini15Flash } from '@genkit-ai/vertexai';
import { defineDotprompt } from '@genkit-ai/dotprompt'
import {MockUserFlowInputSchema, MockUserFlowOutputSchema} from './mockUserFlowTypes'
import { MockUserFlowPrompt } from './prompts';

export const MockUserPrompt = defineDotprompt(
    {
      name: 'mockUserFlow',
      model: gemini15Flash,
      input: {
        schema: MockUserFlowInputSchema,
      },
      output: {
        format: 'json',
        schema: MockUserFlowOutputSchema,
      },  
    }, 
    MockUserFlowPrompt
)
  export const MockUserFlow = defineFlow(
    {
      name: 'mockUserFlow',
      inputSchema: MockUserFlowInputSchema,
      outputSchema: MockUserFlowOutputSchema
    },
    async (input) => {
      try {
        console.log("Generating response...", input);
        const response = await MockUserPrompt.generate({ input: input });
        console.log(response.output(0))
        return response.output(0);
      } catch (error) {
        console.error("Error generating response:", error, input);
        return { 
          relevantMovies: [],
          answer: "",
          justification: ""
         }; 
      }
    }
  );
  