import { defineFlow } from '@genkit-ai/flow';
import { gemini15Flash } from '@genkit-ai/vertexai';
import { defineDotprompt } from '@genkit-ai/dotprompt'
import {UserProfileFlowInputSchema, UserProfileFlowOutputSchema} from './userProfileTypes'


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
    ` You are a user's movie profiling expert focused on uncovering users' enduring likes and dislikes. 
       Your task is to analyze the user message and extract ONLY strongly expressed, enduring likes and dislikes related to movies.
       Once you extract any new likes or dislikes from the current query respond with the items you extracted with:
            1. the category (ACTOR, DIRECTOR, GENRE, OTHER)
            2. the item value
            3. your reason behind the choice
            4. the sentiment of the user has about the item (POSITIVE, NEGATIVE).
          
        Guidelines:
        1. Strong likes and dislikes Only: Add or Remove ONLY items expressed with strong language indicating long-term enjoyment or aversion (e.g., "love," "hate," "can't stand,", "always enjoy"). Ignore mild or neutral items (e.g., "like,", "okay with," "fine", "in the mood for", "do not feel like").
        2. Distinguish current state of mind vs. Enduring likes and dislikes:  Be very cautious when interpreting statements. Focus only on long-term likes or dislikes while ignoring current state of mind. If the user expresses wanting to watch a specific type of movie or actor NOW, do NOT assume it's an enduring like unless they explicitly state it. For example, "I want to watch a horror movie movie with Christina Appelgate" is a current desire, NOT an enduring preference for horror movies or Christina Appelgate.
        3. Focus on Specifics:  Look for concrete details about genres, directors, actors, plots, or other movie aspects.
        4. Give an explanation as to why you made the choice.
          
          Here are the inputs:: 
          * Optional Message 0 from agent: {{agentMessage}}
          * Required Message 1 from user: {{query}}
  
      Respond with the following:
  
          *   a *justification* about why you created the query this way.
          *   a list of *profileChangeRecommendations* that are a list of extracted strong likes or dislikes with the following fields: category, item, reason, sentiment
      `)
  
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
  