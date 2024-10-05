import { embed } from '@genkit-ai/ai/embedder';
import { Document, defineRetriever, retrieve } from '@genkit-ai/ai/retriever';
import { defineFlow } from '@genkit-ai/flow';
import { textEmbedding004 } from '@genkit-ai/vertexai';
import { toSql } from 'pgvector';
import postgres from 'postgres';
import { z } from 'zod';
import { MovieContextSchema, MovieContext } from './movieFlowTypes';
import { OpenDB } from './db';

const sql = postgres({ ssl: false, database: 'recaps' });

const QueryOptionsSchema = z.object({
  query: z.string(),
  k: z.number().optional(),
});


const sqlRetriever = defineRetriever(
  {
    name: 'movies',
    configSchema: QueryOptionsSchema,
  },
  async (input, options) => {
    const db = await OpenDB();
    if (!db) {
      throw new Error('Database connection failed');
    }
    // INSTRUCTIONS:
    // 1. Generate an embedding from the query.
    // 2. Search for the relevant documents in the vector db based on the embedding
    // 3. Convert the model output to type RetrieverFlowOutput
    // HINT: https://firebase.google.com/docs/genkit/templates/pgvector
    return {
      documents: [],
    };
  }
);

export const movieDocFlow = defineFlow(
    {
      name: 'movieDocFlow',
      inputSchema: QueryOptionsSchema,
      outputSchema: z.array(MovieContextSchema), // Array of MovieContextSchema
    },
    async (inputQuestion) => {
      const docs = await retrieve({
        retriever: sqlRetriever,
        query: {
          content: [{ text: inputQuestion.query }], 
        },
        options: {
          k: inputQuestion.k ?? 10,
          query: inputQuestion.query,
        },
      });
  
      console.log(docs);
  
      const movieContexts: MovieContext[] = [];
  
      for (const doc of docs) {
        if (doc.metadata) {
          const movieContext: MovieContext = {
            title: doc.metadata.title,
            runtimeMinutes: doc.metadata.runtimeMinutes,
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
        }
      }
  
      return movieContexts; 
    }
  );