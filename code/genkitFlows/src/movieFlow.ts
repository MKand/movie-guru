/**
 * Copyright 2025 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { ai } from './genkitConfig'
import { MovieFlowInputSchema, MovieFlowOutputSchema } from './movieFlowTypes'
import { GenerationBlockedError } from 'genkit';


/**
 * Prompt file: js/flows/prompts/movie.prompt
 * 
 * This prompt takes the original user input, user preferences, relevant documents retrieved from the database,
 * and conversation history to make a recommendation to the user.
 * 
 * Input schema: MovieFlowInputSchema
 * Output schema: MovieFlowOutputSchema
 * 
 * The MovieGuru development team, uses a variant system to version our prompts. The "v2" variant corresponds to the movie.v2.prompt file. 
 * To use the default variant -- ai.prompt('movie')
 * To use a variant -- ai.prompt('movie', {variant: 'v2'})
 * 
 * ATTENTION: Variant v2 is currently being tested with Gemini 2.5 PRO, if it is not performing well, please revert to the default variant.
 */

const use_pro_prompt = process.env.MOVIEFLOW_PRO_PROMPT || "false"
export var makeMovieRecommendation = ai.prompt('movie');
if (use_pro_prompt == "true"){
  makeMovieRecommendation = ai.prompt('movie', {variant : 'promodel'});
}

export const MovieFlow = ai.defineFlow(
  {
    name: 'movieQAFlow',
    inputSchema: MovieFlowInputSchema,
    outputSchema: MovieFlowOutputSchema,
    

  },
  async (input) => {
    const defaultOutput = MovieFlowOutputSchema.parse({})
    try {
      const response = await makeMovieRecommendation({ history: input.history, userPreferences: input.userPreferences, userMessage: input.userMessage, contextDocuments: input.contextDocuments });
      const safeOutput = response.output ?? defaultOutput;
      const output = MovieFlowOutputSchema.parse(safeOutput);
      return output
    } catch (error) {
      if(error instanceof GenerationBlockedError){
        console.error("MovieFlow: GenerationBlockedError generating response:", error.message);
        return defaultOutput; 
      }
      else if(error instanceof Error && (error.message.includes('429') || error.message.includes('RESOURCE_EXHAUSTED'))){
        console.error("MovieFlow: There is a quota issue:", error.message);
        return defaultOutput;
        }
        else {
        console.error("MovieFlow: Error generating response:", error);
        throw error;
      }
    }
  }
);