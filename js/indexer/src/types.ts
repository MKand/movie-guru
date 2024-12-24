import { z } from 'genkit';

export type MovieContext = z.infer<typeof MovieContextSchema>;

export const MovieContextSchema = z.object({
  title: z.string(),
  runtimeMinutes: z.number(),
  genres: z.array(z.string()),
  rating: z.number(),
  plot: z.string(),
  released: z.number(),
  director: z.string(),
  actors: z.array(z.string()),
  poster: z.string(),
  tconst: z.string(),
});
