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

import { z } from 'genkit';
import { SimpleMessageSchema } from './queryTransformTypes'; 

export const OUTCOME = z.enum([
    'OUTCOMEIRRELEVANT',
    'OUTCOMEACKNOWLEDGED',
    'OUTCOMEENGAGED',
    'OUTCOMETOPICCHANGE',
    'OUTCOMEAMBIGUOUS',
    'OUTCOMEREJECTED',
    'OUTCOMEUNKNOWN'
  ]);

  export const USERSENTIMENT = z.enum([
    'SENTIMENTPOSITIVE',
    'SENTIMENTNEGATIVE',
    'SENTIMENTNEUTRAL',
    'SENTIMENTUNKNOWN',
  ]);

  
// ResponseQualityFlowInput represents the input to the response quality analysis flow.
export const ResponseQualityFlowInputSchema = z.object({
  history: z.array(SimpleMessageSchema),
})

export type ResponseQualityFlowInput = z.infer<typeof ResponseQualityFlowInputSchema>

// ResponseQualityFlowOutput represents the output of the response quality analysis flow.
export const ResponseQualityFlowOutputSchema = z.strictObject({
	outcome: OUTCOME.default('OUTCOMEUNKNOWN'),
	userSentiment: USERSENTIMENT.default('SENTIMENTUNKNOWN'),
})

export type ResponseQualityFlowOutput = z.infer<typeof ResponseQualityFlowOutputSchema>
