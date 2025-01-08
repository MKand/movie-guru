import { z } from 'genkit';
import { SimpleMessageSchema, UserProfileSchema } from './queryTransformTypes'; 
import { ModelOutputMetadataSchema } from './modelOutputMetadataTypes';


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
  userPreferences: UserProfileSchema.optional(),
  contextDocuments: z.array(MovieContextSchema).optional(),
  userMessage: z.string(),
});
export type MovieFlowInput = z.infer<typeof MovieFlowInputSchema>


// MovieFlowOutput schema
export const MovieFlowOutputSchema = z.object({
  answer: z.string(),
  relevantMovies: z.array(RelevantMovieSchema).optional(), // Changed to 'relevantMovies' for clarity
  wrongQuery: z.boolean().optional(),
  modelOutputMetadata: ModelOutputMetadataSchema
});
export type MovieFlowOutput = z.infer<typeof MovieFlowOutputSchema>
