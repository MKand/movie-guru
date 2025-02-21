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

import { ai, safetySettings } from './genkitConfig'
import { MovieFlowInputSchema, ChatFlowOutputSchema, ChatFlowOutput } from './movieFlowTypes'
import { MovieFlowPromptText } from './prompts';
import { GenerationBlockedError } from 'genkit';

export const MovieFlowPrompt = ai.definePrompt(
  {
    name: 'movieFlowPrompt',
    input: {
      schema: MovieFlowInputSchema,
    },
    output: {
      schema: ChatFlowOutputSchema,
      format: 'json',
    },  
    config:{
      safetySettings: safetySettings
      }
  }, 
 MovieFlowPromptText
)
export const MovieFlow = ai.defineFlow(
  {
    name: 'movieQAFlow',
    inputSchema: MovieFlowInputSchema,
    outputSchema: ChatFlowOutputSchema
  },
  async (input) => {
    const defaultOutput = ChatFlowOutputSchema.parse({})
    try {
      const response = await MovieFlowPrompt({ history: input.history, userPreferences: input.userPreferences, userMessage: input.userMessage, contextDocuments: input.contextDocuments });
      const safeOutput = response.output ?? defaultOutput;
      const output = ChatFlowOutputSchema.parse(safeOutput);
      return output
    } catch (error) {
      if(error instanceof GenerationBlockedError){
        console.error("MovieFlow: GenerationBlockedError generating response:", error.message);
        defaultOutput.modelOutputMetadata.safetyIssue = true;
        return defaultOutput; 
      }
      else if(error instanceof Error && (error.message.includes('429') || error.message.includes('RESOURCE_EXHAUSTED'))){
        console.error("MovieFlow: There is a quota issue:", error.message);
        defaultOutput.modelOutputMetadata.quotaIssue = true;
        return defaultOutput;
        }
        else {
        console.error("MovieFlow: Error generating response:", error);
        throw error;
      }
    }
  }
);