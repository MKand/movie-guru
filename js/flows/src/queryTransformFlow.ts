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

import {
  QueryTransformFlowOutputSchema,
  SearchQueryOutputSchema,
  SearchRequiredOutputSchema
} from './queryTransformTypes';
import { ChatFlowInputSchema } from './chatFlowTypes';
import { ai } from './genkitConfig';
import { GenerationBlockedError } from 'genkit';

/**
 * Prompt file: js/flows/prompts/searchRequired.prompt
 * 
 * This prompt asks the LLM to determine if we require a database search for this query.
 * 
 * Input schema: ChatFlowInputSchema
 * Output schema: SearchRequiredOutputSchema
 * 
 * This team, uses a variant system to version our prompts. The "v2" variant corresponds to the searchRequired.v2.prompt file. 
 * To use the default variant -- ai.prompt('searchRequired')
 * To use a variant -- ai.prompt('searchRequired', {variant: 'v2'})
 */
export const isDbSearchRequired = ai.prompt('searchRequired');

/**
 * Prompt file: js/flows/prompts/searchQuery.prompt
 * 
 * This prompt asks the LLM to extract the relevant phrases for a vector search.
 * 
 * Input schema: ChatFlowInputSchema
 * Output schema: SearchQueryOutputSchema
 * 
 * This team, uses a variant system to version our prompts. The "v2" variant corresponds to the searchQuery.v2.prompt file. 
 * To use the default variant -- ai.prompt('searchRequired')
 * To use a variant -- ai.prompt('searchRequired', {variant: 'v2'})
 */
export const createSearchQuery = ai.prompt('searchQuery');

export const QueryTransformFlow = ai.defineFlow(
  {
    name: 'queryTransformFlow',
    inputSchema: ChatFlowInputSchema,
    outputSchema: QueryTransformFlowOutputSchema,
  },
  async (input) => {
    const defaultOutput = QueryTransformFlowOutputSchema.parse({})
    try {
      const response = await isDbSearchRequired(input);
      const safeOutput = response.output?? SearchRequiredOutputSchema.parse({});
      const searchRequiredOutput = SearchRequiredOutputSchema.parse(safeOutput);

      if(searchRequiredOutput.followupAction === 'SEARCH_REQUIRED') {
        const searchQueryResponse = await createSearchQuery(input);
        const safeSearchQueryResponse = searchQueryResponse.output?? SearchQueryOutputSchema.parse({});
        const searchQueryOutput = SearchQueryOutputSchema.parse(safeSearchQueryResponse);

        return {
          searchQuery: searchQueryOutput.searchQuery,
          followupAction: searchRequiredOutput.followupAction,
          justification: searchRequiredOutput.justification + " " + searchQueryOutput.justification
        };
      } else {
        return {
          searchQuery: "",
          followupAction: searchRequiredOutput.followupAction,
          justification: searchRequiredOutput.justification
        };
      }
    } catch (error) {
      if (error instanceof GenerationBlockedError){
        
        console.error("QTFlow: GenerationBlockedError generating response:", error.message);
        return defaultOutput;
      }
      else if(error instanceof Error && (error.message.includes('429') || error.message.includes('RESOURCE_EXHAUSTED'))){
        console.error("QTFlow: There is a quota issue:", error.message);
        return defaultOutput;
        }
        else {
        console.error("QTFlow: Error generating response:", error);
        throw error;
      }
      
    }
  }
);
