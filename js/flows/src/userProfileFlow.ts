import { gemini15Flash } from '@genkit-ai/vertexai';
import {UserProfileFlowInputSchema, UserProfileFlowOutputSchema} from './userProfileTypes'
import { UserProfilePromptText } from './prompts';
import { ai } from './genkitConfig'

export const UserProfileFlowPrompt = ai.definePrompt(
    {
      name: 'userProfileFlowPrompt',
      model: gemini15Flash,
      input: {
        schema: UserProfileFlowInputSchema,
      },
      output: {
        format: 'json',
      },  
    }, 
    UserProfilePromptText)
  
  export const UserProfileFlow = ai.defineFlow(
    {
      name: 'userProfileFlow',
      inputSchema: UserProfileFlowInputSchema,
      outputSchema: UserProfileFlowOutputSchema
    },
    async (input) => {
      try {
        const response = await UserProfileFlowPrompt({ query: input.query, agentMessage: input.agentMessage });
        return JSON.parse(response.text);
      } catch (error) {
        console.error("Error generating response:", error);
        return { 
          profileChangeRecommendations: [],
          justification: ""
         }; 
      }
    }
  );
  