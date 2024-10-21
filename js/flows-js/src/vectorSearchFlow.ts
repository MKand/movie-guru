
import { z } from 'zod';
import { MovieContextSchema, MovieContext } from './movieFlowTypes';
import { defineFlow } from '@genkit-ai/flow';
import { retrieve } from '@genkit-ai/ai/retriever';
import { vectorRetriever } from './mixedSearchFlow';

export const SearchFlowInputSchema = z.object({
    inputQuery: z.string(),
  });
  
  

export const VectorSearchFlow = defineFlow(
    {
      name: 'VectorSearchFlow',
      inputSchema: SearchFlowInputSchema,
      outputSchema: z.array(MovieContextSchema)
    },
    async (input) => {
      const movieContexts: MovieContext[] = [];
      const docs = await retrieve({
        retriever: vectorRetriever,
        query: {
          content: [{ text: input.inputQuery }],
        },
        options: {
          k: 10,
        },
      });
      for (const doc of docs) {
        if (doc.metadata) {
          const movieContext: MovieContext = {
            title: doc.metadata.title,
            runtime_minutes: doc.metadata.runtime_mins,
            genres: doc.metadata.genres,
            rating: doc.metadata.rating,
            plot: doc.metadata.plot,
            released: doc.metadata.released,
            director: doc.metadata.director,
            actors: doc.metadata.actors,
            poster: doc.metadata.poster,
            tconst: doc.metadata.tconst,
          };
          movieContexts.push(movieContext);
        } else {
          console.warn('Movie metadata is missing for a document.');
          return []
        }
      }
      return movieContexts
  
    }
  );
  