import { embed } from '@genkit-ai/ai/embedder';
import { Document, defineRetriever, retrieve } from '@genkit-ai/ai/retriever';
import { defineFlow } from '@genkit-ai/flow';
import { textEmbedding004 } from '@genkit-ai/vertexai';
import { toSql } from 'pgvector';
import { z } from 'zod';
import { MovieContextSchema, MovieContext } from './movieFlowTypes';
import { openDB } from './db';


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
    const db = await openDB();
    if (!db) {
      throw new Error('Database connection failed');
    }
    const embedding = await embed({
      embedder: textEmbedding004,
      content: input,
    });
    const results = await db`
      SELECT title, poster, content, released, runtime_mins, rating, genres, director, actors, plot
	    FROM movies
        ORDER BY embedding <#> ${toSql(embedding)}
        LIMIT ${options.k ?? 10}
      `;
    return {
      documents: results.map((row) => {
        const { content, ...metadata } = row;
        return Document.fromText(content, metadata);
      }),
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