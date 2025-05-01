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

import { ChatFlowInputSchema, ChatFlowOutput, ChatOutputSchema } from './chatFlowTypes';
import { MovieContext, RelevantMovie } from './movieFlowTypes';

import { ai } from './genkitConfig';
import { GenerationBlockedError } from 'genkit';

import {  SafetyTransformPrompt, SafetyPromptOutputSchema } from './safetyFlow';
import { QueryTransformPrompt } from './queryTransformFlow';
import { QueryTransformFlowOutputSchema } from './queryTransformTypes';
import { MovieDocFlow } from './docRetriever';
import { MovieFlowPrompt } from './movieFlow';
import { MovieFlowOutputSchema } from './movieFlowTypes';

export const ChatFlow = ai.defineFlow(
    {
        name: "chatFlow",
        inputSchema: ChatFlowInputSchema,
        outputSchema: ChatOutputSchema,
    },
    async(input) => {
        const chatResponse: ChatFlowOutput = ChatOutputSchema.parse({});
        
        try {
        
            // Initial safety check
            const safetyRawOutput =  await SafetyTransformPrompt({
                userMessage: input.userMessage,
                });
                const defaultSafetyOutput = safetyRawOutput.output ??  SafetyPromptOutputSchema.parse({});
                const safetyOutput = SafetyPromptOutputSchema.parse(defaultSafetyOutput);
            
            if(safetyOutput.safetyIssue == true || safetyOutput.wrongQuery == true){
                chatResponse.modelOutputMetadata.safetyIssue = safetyOutput.safetyIssue;
                chatResponse.wrongQuery = safetyOutput.wrongQuery
                return chatResponse;
            }
            
            // Search Required Check

            // Add 1 min delay
            await new Promise(resolve => setTimeout(resolve, 5 * 1000));
            const qtRawOutput = await QueryTransformPrompt( ChatFlowInputSchema.parse({
                history: input.history,
                userPreferences: input.userPreferences,
                userMessage: input.userMessage
            }))
            const defaultQTOutput = qtRawOutput.output ??  QueryTransformFlowOutputSchema.parse({});
            const qtOutput = QueryTransformFlowOutputSchema.parse(defaultQTOutput);

            // Search if required

            // Add 1 min delay
            await new Promise(resolve => setTimeout(resolve, 5 * 1000));
            var movieContexts: MovieContext[] = []
            if(qtOutput.followupAction == "SEARCH_REQUIRED"){
                movieContexts = await MovieDocFlow( {query: qtOutput.searchQuery})
            }
            
            // Final RAG
            // Add 1 min delay
            await new Promise(resolve => setTimeout(resolve, 5 * 1000));
            const movieFlowRawOutput = await MovieFlowPrompt({
                history: input.history,
                userPreferences: input.userPreferences,
                contextDocuments: movieContexts,
                userMessage: input.userMessage
            })
            const defaultMovieQOutput = movieFlowRawOutput.output ??  MovieFlowOutputSchema.parse({});
            const movieQAOutput = MovieFlowOutputSchema.parse(defaultMovieQOutput);

            // Transform into chat Response
            chatResponse.answer = movieQAOutput.response;
            chatResponse.relevantMovies = movieQAOutput.relevantMovies;
            chatResponse.contextDocuments = parseContexts(movieQAOutput.relevantMovies, movieContexts);
            chatResponse.modelOutputMetadata.justification = movieQAOutput.justification;

            return chatResponse
        }
        catch (error) {
            if (error instanceof GenerationBlockedError){
        
                console.error("ChatFlow: GenerationBlockedError generating response:", error.message);
                chatResponse.modelOutputMetadata.safetyIssue = true;
                return chatResponse;
            }
            else if(error instanceof Error && (error.message.includes('429') || error.message.includes('RESOURCE_EXHAUSTED'))){
                console.error("ChatFlow: There is a quota issue:", error.message);
                chatResponse.modelOutputMetadata.quotaIssue = true;
                return chatResponse;
            }
            else {
                console.error("ChatFlow: Error generating response:", error);
                throw error;
            }
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
