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

export const ChatFlow = ai.defineFlow(
    {
        name: "chatFlow",
        inputSchema: ChatFlowInputSchema,
        outputSchema: ChatOutputSchema,
    },
    async(input) => {
            const chatResponse: ChatFlowOutput = ChatOutputSchema.parse({});
            try{
            const safetyRawOutput =  await SafetyTransformPrompt({
                userMessage: input.userMessage,
              });
              const safeOutput = safetyRawOutput.output ??  SafetyPromptOutputSchema.parse({});
              const safetyOutput = SafetyPromptOutputSchema.parse(safeOutput);
            if(safetyOutput.safetyIssue == true || safetyOutput.wrongQuery == true){
                chatResponse.modelOutputMetadata.safetyIssue = safetyOutput.safetyIssue;
                chatResponse.wrongQuery = safetyOutput.wrongQuery
                return chatResponse;
            }


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
