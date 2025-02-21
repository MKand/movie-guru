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

import {ResponseQualityFlowInputSchema, ResponseQualityFlowOutputSchema, OUTCOME, USERSENTIMENT, ResponseQualityFlowOutput} from './verifyQualityTypes'
import { ConversationQualityAnalysisPromptText } from './prompts';
import { ai } from './genkitConfig'


export const QualityFlowPrompt = ai.definePrompt(
    {
      name: 'qualityFlowPrompt',
      input: {
        schema: ResponseQualityFlowInputSchema,
      },
      output: {
        schema: ResponseQualityFlowOutputSchema,
        format: 'json',
      },  
    }, 
    ConversationQualityAnalysisPromptText
  )
  
  export const QualityFlow = ai.defineFlow(
    {
      name: 'qualityFlow',
      inputSchema: ResponseQualityFlowInputSchema,
      outputSchema: ResponseQualityFlowOutputSchema
    },
    async (input) => {
      try {
        const response = await QualityFlowPrompt({ history: input.history });
        const safeOutput = response.output?? {
          outcome: OUTCOME.parse('OUTCOMEUNKNOWN'), 
          userSentiment: USERSENTIMENT.parse('SENTIMENTUNKNOWN') 
       };

        console.log("quality response:", response.output)
        const output = ResponseQualityFlowOutputSchema.parse(safeOutput);
        return output;
      } catch (error) {
        console.error("Error generating response:", error);
        const output: ResponseQualityFlowOutput = {  
          outcome: OUTCOME.parse('OUTCOMEUNKNOWN'),
          userSentiment: USERSENTIMENT.parse('SENTIMENTUNKNOWN'),         
         }; 
         return output
      }
    }
  );
  