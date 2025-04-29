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

import {ResponseQualityFlowInputSchema, ResponseQualityFlowOutputSchema} from './verifyQualityTypes'
import { ai } from './genkitConfig'


export const QualityFlowPrompt = ai.prompt( 'verifyQuality');
  
export const QualityFlow = ai.defineFlow(
  {
    name: 'qualityFlow',
    inputSchema: ResponseQualityFlowInputSchema,
    outputSchema: ResponseQualityFlowOutputSchema
  },
  async (input) => {
    const defaultOutput = ResponseQualityFlowOutputSchema.parse({})
    try {
      const response = await QualityFlowPrompt({ history: input.history });
      const safeOutput = response.output?? defaultOutput

      console.log("quality response:", response.output)
      const output = ResponseQualityFlowOutputSchema.parse(safeOutput);
      return output
    } catch (error) {
      console.error("Quality Flow: Error generating response:", error);
      return defaultOutput
    }
  }
  );
  