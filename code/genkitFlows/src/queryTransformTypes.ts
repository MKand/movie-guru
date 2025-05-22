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
import { ai } from './genkitConfig';
import { ChatFlowInputSchema } from './chatFlowTypes';

// FOLLOWUP_ACTION as Zod Enum
export const FOLLOWUP_ACTION = z.enum([
  'SEARCH_REQUIRED',
  'SEARCH_NOT_REQUIRED',
]);


// ProfileCategories schema
const ProfileCategoriesSchema = z.object({
  actors: z.array(z.string()).optional(),
  directors: z.array(z.string()).optional(),
  genres: z.array(z.string()).optional(),
  others: z.array(z.string()).optional(),
});

export type ProfileCategories = z.infer<typeof ProfileCategoriesSchema>

ai.defineSchema('ProfileCategoriesSchema', ProfileCategoriesSchema);

// UserProfile schema
export const UserProfileSchema = z.object({
  likes: ProfileCategoriesSchema.optional(),
  dislikes: ProfileCategoriesSchema.optional(),
});

export type UserProfile = z.infer<typeof UserProfileSchema>

ai.defineSchema('UserProfileSchema', UserProfileSchema);

// SimpleMessage schema
export const SimpleMessageSchema = z.object({
  role: z.string(),
  content: z.string(),
});

export type SimpleMessage = z.infer<typeof SimpleMessageSchema>

ai.defineSchema('SimpleMessageSchema', SimpleMessageSchema);

export type QueryTransformFlowInput = z.infer<typeof ChatFlowInputSchema>

// Schema is used as the output schema for the searchRequired prompt.
// See js/flows/prompts/searchRequired.prompt
export const SearchRequiredOutputSchema = z.strictObject({
  followupAction: FOLLOWUP_ACTION.default("SEARCH_NOT_REQUIRED"),
  justification: z.string().default("No justification provided")
});

export type SearchRequiredOutput = z.infer<typeof SearchRequiredOutputSchema>

ai.defineSchema('SearchRequiredOutputSchema', SearchRequiredOutputSchema);

// Schema is used as the output schema for the searchRequired prompt.
// See js/flows/prompts/searchQuery.prompt
export const SearchQueryOutputSchema = z.strictObject({
  searchQuery: z.string().optional().default(""),
  justification: z.string().default("No justification provided")
});

export type SearchQueryOutput = z.infer<typeof SearchQueryOutputSchema>

ai.defineSchema('SearchQueryOutputSchema', SearchQueryOutputSchema);

// QueryTransformFlowOutput schema
export const QueryTransformFlowOutputSchema = z.strictObject({
  searchQuery: z.string().optional().default(""),
  followupAction: FOLLOWUP_ACTION.default("SEARCH_NOT_REQUIRED"),
  justification: z.string().default("No justification provided")
});

export type QueryTransformFlowOutput = z.infer<typeof QueryTransformFlowOutputSchema>

ai.defineSchema('QueryTransformFlowOutputSchema', QueryTransformFlowOutputSchema);

