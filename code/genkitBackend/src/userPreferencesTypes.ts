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
import { ai } from './genkitConfig'

// Enums as Zod Enums
const MovieFeatureCategory = z.enum(['OTHER', 'ACTOR', 'DIRECTOR', 'GENRE']);
const Sentiment = z.enum(['POSITIVE', 'NEGATIVE']);

// ProfileChangeRecommendation schema
export const PreferenceItemSchema = z.object({
  item: z.string(),
  reason: z.string(),
  category: MovieFeatureCategory,
  sentiment: Sentiment,
});

export type PreferenceItem = z.infer<typeof PreferenceItemSchema>


// UserProfileFlowInput schema
export const UserPreferenceInputSchema = z.object({
  query: z.string(),
  agentMessage: z.string(),
});

ai.defineSchema('UserPreferenceInputSchema', UserPreferenceInputSchema);

export const UserPreferenceOutputSchema = z.strictObject({
  profileItems: z.array(PreferenceItemSchema).optional().default([]),
  justification: z.string().default("No justification provided"),
});

export const UserPreferenceFlowSchema = z.strictObject({
  profileItems: z.array(PreferenceItemSchema).optional().default([]),
  safetyIssue: z.boolean().optional().default(false),
  quotaIssue: z.boolean().optional().default(false),
  justification: z.string().default("No justification provided"),

});

ai.defineSchema('UserPreferencePromptOutputSchema', UserPreferenceOutputSchema);
ai.defineSchema('UserPreferenceFlowOutputSchema', UserPreferenceFlowSchema);

