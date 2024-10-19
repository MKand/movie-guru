import { defineFlow } from '@genkit-ai/flow';
import { gemini15Flash } from '@genkit-ai/vertexai';
import { defineDotprompt } from '@genkit-ai/dotprompt'
import {UserProfileFlowInputSchema, UserProfileFlowOutputSchema} from './userProfileTypes'


export const UserProfileFlowPrompt = defineDotprompt(
    {
      name: 'UserProfileFlowPrompt',
      model: gemini15Flash,
      input: {
        schema: UserProfileFlowInputSchema,
      },
      output: {
        format: 'json',
        schema: UserProfileFlowOutputSchema,
      },  
    }, 
    `You are a user's movie profiling expert focused on uncovering users' enduring likes and dislikes. 
  Your task is to analyze the user message and extract ONLY strongly expressed, enduring likes and dislikes related to movies.
  Once you extract any new likes or dislikes from the current query respond with the items you extracted with:
   1. the category (ACTOR, DIRECTOR, GENRE, OTHER)
   2. the item value
   3. your reason behind the choice
   4. the sentiment of the user has about the item (POSITIVE, NEGATIVE).
   
  Guidelines:
  1. Strong likes and dislikes Only: Add or Remove ONLY items expressed with strong language indicating long-term enjoyment or aversion (e.g., "love," "hate," "can't stand,", "always enjoy"). Ignore mild or neutral items (e.g., "like,", "okay with," "fine", "in the mood for", "do not feel like").
  2. Distinguish current state of mind vs. Enduring likes and dislikes:  Focus only on long-term likes or dislikes while ignoring current state of mind. 
  
  Examples:
   ---
   userMessage: "I want to watch a horror movie with Christina Appelgate" 
   output: profileChangeRecommendations:[]
   ---
   userMessage: "I love horror movies and want to watch one with Christina Appelgate" 
   output: profileChangeRecommendations=[
   item: horror,
   category: genre,
   reason: The user specifically stated they love horror indicating a strong preference. They are looking for one with Christina Appelgate, which is a current desire and not an enduring preference.
   sentiment: POSITIVE]
   ---
   userMessage: "Show me some action films" 
   output: profileChangeRecommendations:[]
   ---
   userMessage: "I dont feel like watching an action film" 
   output: profileChangeRecommendations:[]
   ---
   userMessage: "I dont like action films" 
   output: profileChangeRecommendations=[
   item: action,
   category: genre,
   reason: The user specifically states they don't like action films which is a statement that expresses their long term disklike for action films.
   sentiment: NEGATIVE]
Here are the inputs:
	1. Optional Message 0 from agent: {{agentMessage}}
	2. Required Message 1 from user: {{userQuery}}
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
  