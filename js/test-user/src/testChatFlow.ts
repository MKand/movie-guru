import { defineFlow } from '@genkit-ai/flow';
import { gemini15Flash } from '@genkit-ai/vertexai';
import { defineDotprompt } from '@genkit-ai/dotprompt'
import {DummyUserFlowInputSchema, DummyUserFlowOutputSchema} from './testChatFlowTypes'
import { DummyUserFlowPrompt } from './prompts';

export const DummyUserPrompt = defineDotprompt(
    {
      name: 'dummyUserFlow',
      model: gemini15Flash,
      input: {
        schema: DummyUserFlowInputSchema,
      },
      output: {
        format: 'json',
        schema: DummyUserFlowOutputSchema,
      },  
    }, 
    DummyUserFlowPrompt
)
  export const DummyUserFlow = defineFlow(
    {
      name: 'dummyUserFlow',
      inputSchema: DummyUserFlowInputSchema,
      outputSchema: DummyUserFlowOutputSchema
    },
    async (input) => {
      try {
        console.log("Generating response...", input);
        const response = await DummyUserPrompt.generate({ input: input });
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
  