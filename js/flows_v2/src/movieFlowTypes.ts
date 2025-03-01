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
import { SimpleMessageSchema, UserProfileSchema } from './queryTransformTypes';
import { ModelOutputMetadata, ModelOutputMetadataSchema } from './modelOutputMetadataTypes';


// RelevantMovie schema
export const RelevantMovieSchema = z.object({
  title: z.string(),
  reason: z.string().optional(),
});
export type RelevantMovie = z.infer<typeof RelevantMovieSchema>

// MovieContext schema
export const MovieContextSchema = z.object({
  title: z.string(),
  runtime_minutes: z.number(),
  genres: z.array(z.string()),
  rating: z.number(),
  plot: z.string(),
  released: z.number(),
  director: z.string(),
  actors: z.array(z.string()),
  poster: z.string(),
  tconst: z.string().optional(),
});
export type MovieContext = z.infer<typeof MovieContextSchema>

// MovieFlowInput schema
export const MovieFlowInputSchema = z.object({
  history: z.array(SimpleMessageSchema),
  userPreferences: UserProfileSchema.optional().default(UserProfileSchema.parse({})),
  contextDocuments: z.array(MovieContextSchema).optional().default([]),
  userMessage: z.string(),
});
export type MovieFlowInput = z.infer<typeof MovieFlowInputSchema>


// MovieFlowOutput schema
export const MovieFlowOutputSchema = z.strictObject({
  answer: z.string().default(""),
  relevantMovies: z.array(RelevantMovieSchema).optional().default([]), // Changed to 'relevantMovies' for clarity
  wrongQuery: z.boolean().optional().default(false),
  modelOutputMetadata: ModelOutputMetadataSchema.default(ModelOutputMetadataSchema.parse({}))
});
