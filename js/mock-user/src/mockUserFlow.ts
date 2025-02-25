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
import { gemini15Flash } from '@genkit-ai/vertexai';
import {MockUserFlowInputSchema, MockUserFlowOutputSchema, MockUserFlowOutput} from './mockUserFlowTypes'
import { MockUserFlowPromptText } from './prompts';

export const MockUserPrompt = ai.definePrompt(
    {
      name: 'mockUserFlow',
      model: gemini15Flash,
      input: {
        schema: MockUserFlowInputSchema,
      },
      output: {
        format: 'json',
      },  
    }, 
    MockUserFlowPromptText
)
  export const MockUserFlow = ai.defineFlow(
    {
      name: 'mockUserFlow',
      inputSchema: MockUserFlowInputSchema,
      outputSchema: MockUserFlowOutputSchema
    },
    async (input) => {
      try {
        const response = await MockUserPrompt({ expert_answer: input.expert_answer, response_mood: input.response_mood, response_type: input.response_type });
        const jsonResponse =  JSON.parse(response.text);
        const output: MockUserFlowOutput = {
          "answer":  jsonResponse.answer,
        }
        return output
      } catch (error) {
        console.error("Error generating response:", error, input);
        return { 
          answer: "",
         }; 
      }
    }
  );
  