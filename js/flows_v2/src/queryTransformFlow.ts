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
  QueryTransformFlowOutput
} from './queryTransformTypes';
import { ChatFlowInputSchema } from './chatFlowTypes';
import { QueryTransformPromptText } from './prompts';
import { ai, safetySettings } from './genkitConfig';
import { GenerationBlockedError } from 'genkit';

export const QueryTransformPrompt = ai.definePrompt(
  {
    name: 'queryTransformFlowPrompt',
    input: {
      schema: ChatFlowInputSchema,
    },
    output: {
      schema: QueryTransformFlowOutputSchema,
      format: 'json',
    },
    config:{
      safetySettings: safetySettings
      }
  },
  
  QueryTransformPromptText
);

export const QueryTransformFlow = ai.defineFlow(
  {
    name: 'queryTransformFlow',
    inputSchema: ChatFlowInputSchema,
    outputSchema: QueryTransformFlowOutputSchema,
  },
  async (input) => {
    const defaultOutput = QueryTransformFlowOutputSchema.parse({})
    try {
      const response = await QueryTransformPrompt({
        history: input.history,
        userMessage: input.userMessage,
        userPreferences: input.userPreferences,
      });
      const safeOutput = response.output?? defaultOutput;
      const output = QueryTransformFlowOutputSchema.parse(safeOutput)
      return output;
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
