import { embed } from '@genkit-ai/ai/embedder';
import { Document, defineRetriever, retrieve } from '@genkit-ai/ai/retriever';
import { defineFlow } from '@genkit-ai/flow';
import { textEmbedding004 } from '@genkit-ai/vertexai';
import { toSql } from 'pgvector';
import { z } from 'zod';
import { openDB } from './db';
import { MovieContextSchema, MovieContext } from './movieFlowTypes';

const RetrieverOptionsSchema = z.object({
  k: z.number().optional().default(10),
});

const SearchFlowInputSchema = z.object({
    inputQuery: z.string(),
});
  
export const vectorSearchFlow = defineFlow(
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

// Defining the vector Retriever
export const vectorRetriever = defineRetriever(
  {
    name: 'vectorRetriever',
    configSchema: RetrieverOptionsSchema,
  },
  async (query, options) => {
    const db = await openDB();
    if (!db) {
      throw new Error('Database connection failed');
    }
    const embedding = await embed({
      embedder: textEmbedding004,
      content: query.text(),
    });
    const results = await db`
      SELECT content, title, poster, released, runtime_mins, rating, genres, director, actors, plot, tconst
     FROM movies
        ORDER BY embedding <#> ${toSql(embedding)}
        LIMIT ${options.k ?? 10}
      ;`
    return {
      documents: results.map((row) => {
        const { content, ...metadata } = row;
        return Document.fromText(content, metadata);
      }),
    };
  }
);