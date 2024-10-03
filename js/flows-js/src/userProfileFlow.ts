import { defineFlow } from '@genkit-ai/flow';
import { gemini15Flash } from '@genkit-ai/vertexai';
import { defineDotprompt } from '@genkit-ai/dotprompt'
import {UserProfileFlowInputSchema, UserProfileFlowOutputSchema} from './userProfileTypes'
import { UserProfilePromptText } from './prompts';


export const UserProfileFlowPrompt = defineDotprompt(
    {
      name: 'userProfileFlow',
      model: gemini15Flash,
      input: {
        schema: UserProfileFlowInputSchema,
      },
      output: {
        format: 'json',
        schema: UserProfileFlowOutputSchema,
      },  
    }, 
    UserProfilePromptText)
  
  export const UserProfileFlow = defineFlow(
    {
      name: 'userProfileFlow',
      inputSchema: UserProfileFlowInputSchema,
      outputSchema: UserProfileFlowOutputSchema
    },
    async (input) => {
      try {
        const response = await UserProfileFlowPrompt.generate({ input: input });
        console.log(response.output(0))
        return response.output(0);
      } catch (error) {
        console.error("Error generating response:", error);
        return { 
          profileChangeRecommendations: [],
          justification: ""
         }; 
      }
    }
  );
  