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

// Defined in js/flows/prompts/movieSearch.prompt
// This has two variants. For more about how to use variants, see https://firebase.google.com/docs/genkit/dotprompt#prompt_variants
export const MovieFlowPrompt = ai.prompt('movie');

export const MovieFlow = ai.defineFlow(
  {
    name: 'movieQAFlow',
    inputSchema: MovieFlowInputSchema,
    outputSchema: MovieFlowOutputSchema,
    

  },
  async (input) => {
    const defaultOutput = MovieFlowOutputSchema.parse({})
    try {
      const response = await MovieFlowPrompt({ history: input.history, userPreferences: input.userPreferences, userMessage: input.userMessage, contextDocuments: input.contextDocuments });
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