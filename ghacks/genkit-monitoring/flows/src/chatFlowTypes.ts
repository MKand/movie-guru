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
import { SimpleMessageSchema } from './queryTransformTypes';
import { UserProfileSchema } from './queryTransformTypes';
import { RelevantMovieSchema } from './movieFlowTypes';
import { MovieContextSchema } from './movieFlowTypes';


// ChatFlowInput schema
export const ChatFlowInputSchema = z.object({
    history: z.array(SimpleMessageSchema),
    userPreferences: UserProfileSchema.optional().default(UserProfileSchema.parse({})),
    userMessage: z.string(),
  });

ai.defineSchema('ChatFlowInputSchema', ChatFlowInputSchema);


// ChatFlowOutput schema
export const ChatFlowOutputSchema = z.strictObject({
  answer: z.string().optional().default(""),
  relevantMovies: z.array(RelevantMovieSchema).optional().default([]),
  contextDocuments: z.array(MovieContextSchema).optional().default([]),
  badQuery: z.boolean().optional().default(false),
  safetyIssue: z.boolean().optional().default(false),
  quotaIssue: z.boolean().optional().default(false),
  justification: z.string().default("No justification provided"),
});

export type ChatFlowOutput = z.infer<typeof ChatFlowOutputSchema>;
ai.defineSchema('ChatOutputSchema', ChatFlowOutputSchema);

  