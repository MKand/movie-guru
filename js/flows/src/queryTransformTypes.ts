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
import { ModelOutputMetadata, ModelOutputMetadataSchema } from './modelOutputMetadataTypes';

// USERINTENT as Zod Enum
export const USERINTENT = z.enum([
  'UNCLEAR',
  'GREET',
  'END_CONVERSATION',
  'REQUEST',
  'RESPONSE',
  'ACKNOWLEDGE',
]);


// ProfileCategories schema
const ProfileCategoriesSchema = z.object({
  actors: z.array(z.string()).optional(),
  directors: z.array(z.string()).optional(),
  genres: z.array(z.string()).optional(),
  others: z.array(z.string()).optional(),
});

export type ProfileCategories = z.infer<typeof ProfileCategoriesSchema>


// UserProfile schema
export const UserProfileSchema = z.object({
  likes: ProfileCategoriesSchema.optional(),
  dislikes: ProfileCategoriesSchema.optional(),
});

export type UserProfile = z.infer<typeof UserProfileSchema>


// SimpleMessage schema
export const SimpleMessageSchema = z.object({
  role: z.string(),
  content: z.string(),
});

export type SimpleMessage = z.infer<typeof SimpleMessageSchema>


// QueryTransformFlowInput schema
export const ChatFlowInputSchema = z.object({
  history: z.array(SimpleMessageSchema),
  userPreferences: UserProfileSchema.optional().default(UserProfileSchema.parse({})),
  userMessage: z.string(),
});

export type QueryTransformFlowInput = z.infer<typeof ChatFlowInputSchema>


// QueryTransformFlowOutput schema
export const QueryTransformFlowOutputSchema = z.strictObject({
  transformedQuery: z.string().optional().default(""),
  userIntent: USERINTENT.default("UNCLEAR"),
  modelOutputMetadata: ModelOutputMetadataSchema.default(ModelOutputMetadataSchema.parse({})),
});

export type QueryTransformFlowOutput = z.infer<typeof QueryTransformFlowOutputSchema>
