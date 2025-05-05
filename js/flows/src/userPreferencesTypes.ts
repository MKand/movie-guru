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
import { ModelOutputMetadataSchema } from './modelOutputMetadataTypes';

// Enums as Zod Enums
const MovieFeatureCategory = z.enum(['OTHER', 'ACTOR', 'DIRECTOR', 'GENRE']);
const Sentiment = z.enum(['POSITIVE', 'NEGATIVE']);

// ProfileChangeRecommendation schema
export const ProfileChangeRecommendationSchema = z.object({
  item: z.string(),
  reason: z.string(),
  category: MovieFeatureCategory,
  sentiment: Sentiment,
});

export type ProfileChangeRecommendation = z.infer<typeof ProfileChangeRecommendationSchema>

// UserProfileFlowInput schema
export const UserPreferenceFlowInputSchema = z.object({
  query: z.string(),
  agentMessage: z.string(),
});

export type UserPreferenceFlowInput = z.infer<typeof UserPreferenceFlowInputSchema>

ai.defineSchema('UserPreferenceFlowInputSchema', UserPreferenceFlowInputSchema);


// HINT: Challenge 3. This is the output schema for the userProfile Flow. 
// Compare this to the schema in the error message.

export const UserPreferenceFlowOutputSchema = z.strictObject({
  profileChangeRecommendations: z.array(ProfileChangeRecommendationSchema).optional().default([]),
  modelOutputMetadata: z.object({
    justification: z.string().default("Unknown"),
    safetyIssue: z.boolean().optional().default(false),
    quotaIssue: z.boolean().optional().default(false)
  }).default({
    justification: "Unknown",
    safetyIssue: false,
    quotaIssue: false
  }),
});

export type UserPreferenceFlowOutput = z.infer<typeof UserPreferenceFlowOutputSchema>

ai.defineSchema('UserPreferenceFlowOutputSchema', UserPreferenceFlowOutputSchema);

