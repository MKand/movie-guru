// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import {
    USERINTENT,
    QueryTransformFlowOutput
  } from './queryTransformTypes';
import { ChatFlowInputSchema, ChatFlowOutput, ChatOutputSchema } from './chatFlowTypes';
import { MovieContext, RelevantMovie } from './movieFlowTypes';

import { ai } from './genkitConfig';
import { QueryTransformFlow } from './queryTransformFlow';
import { MovieDocFlow } from './docRetriever';
import { MovieFlow } from './movieFlow';

export const ChatFlow = ai.defineFlow(
    {
        name: "chatFlow",
        inputSchema: ChatFlowInputSchema,
        outputSchema: ChatOutputSchema,
    },
    async(input) => {
        try{
            const chatResponse: ChatFlowOutput = ChatOutputSchema.parse({})
            const qtResponse: QueryTransformFlowOutput = await QueryTransformFlow(input)
            if (qtResponse.modelOutputMetadata.safetyIssue || qtResponse.modelOutputMetadata.quotaIssue){
                 chatResponse.modelOutputMetadata = qtResponse.modelOutputMetadata
                 return chatResponse
            }
            var movieContexts: MovieContext[] = []
            if (qtResponse.userIntent==USERINTENT.parse("REQUEST" ) || qtResponse.userIntent==USERINTENT.parse("RESPONSE")){
                 movieContexts = await MovieDocFlow( {query: qtResponse.transformedQuery})
            }
            const movieFlowResponse = await MovieFlow({
                history: input.history,
                userPreferences: input.userPreferences,
                contextDocuments: movieContexts,
                userMessage: input.userMessage
            })
            
            chatResponse.answer = movieFlowResponse.answer;
            chatResponse.modelOutputMetadata = movieFlowResponse.modelOutputMetadata
            chatResponse.relevantMovies = movieFlowResponse.relevantMovies
            chatResponse.wrongQuery = movieFlowResponse.wrongQuery
            chatResponse.contextDocuments = parseContexts(movieFlowResponse.relevantMovies, movieContexts)
            
            return chatResponse
        }
        catch (error) {
            console.error("ChatFlow: Error generating response:", error);
            throw error;
        }
    }
)

function parseContexts(relevantMovies: RelevantMovie [], movieContexts:MovieContext[]): MovieContext[]{
   const relevantContext: MovieContext[] = []
   for (const r of relevantMovies){
        for(const c of movieContexts){
            if(r.title  == c.title){
                relevantContext.push(c)
            }
        }
   }
   return relevantContext

}
