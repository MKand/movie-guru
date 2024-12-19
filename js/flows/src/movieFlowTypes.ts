import { z } from 'genkit';

import { SimpleMessageSchema, UserProfileSchema } from './queryTransformTypes'; 


// RelevantMovie schema
export const RelevantMovieSchema = z.object({
  title: z.string(),
  reason: z.string(),
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
  userPreferences: UserProfileSchema,
  contextDocuments: z.array(MovieContextSchema),
  userMessage: z.string(),
});
export type MovieFlowInput = z.infer<typeof MovieFlowInputSchema>


// MovieFlowOutput schema
export const MovieFlowOutputSchema = z.object({
  answer: z.string(),
  relevantMovies: z.array(RelevantMovieSchema), // Changed to 'relevantMovies' for clarity
  wrongQuery: z.boolean().optional(),
  justification: z.string().optional(),
});
export type MovieFlowOutput = z.infer<typeof MovieFlowOutputSchema>