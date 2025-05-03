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

/**
 * Prompt file: js/flows/prompts/verifyQuality.prompt
 * 
 * This prompt instructs the LLM to assess the quality of the response based on the user's reaction to it.
 * 
 * Input schema: ResponseQualityFlowInputSchema
 * Output schema: ResponseQualityFlowOutputSchema
 * 
 * This uses a variant system to version our prompts. The "v2" variant corresponds to the verifyQuality.v2.prompt file. 
 * To use the default variant -- ai.prompt('verifyQuality')
 * To use a variant -- ai.prompt('verifyQuality', {variant: 'v2'})
 */
export const verifyResponseQuality = ai.prompt( 'verifyQuality');
  
export const QualityFlow = ai.defineFlow(
  {
    name: 'qualityFlow',
    inputSchema: ResponseQualityFlowInputSchema,
    outputSchema: ResponseQualityFlowOutputSchema
  },
  async (input) => {
    const defaultOutput = ResponseQualityFlowOutputSchema.parse({})
    try {
      const response = await verifyResponseQuality({ history: input.history });
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
  