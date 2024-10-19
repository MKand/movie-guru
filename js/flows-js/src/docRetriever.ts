import { embed } from '@genkit-ai/ai/embedder';
import { Document, defineRetriever, retrieve } from '@genkit-ai/ai/retriever';
import { defineFlow } from '@genkit-ai/flow';
import { textEmbedding004 } from '@genkit-ai/vertexai';
import { toSql } from 'pgvector';
import { z } from 'zod';
import { MovieContextSchema, MovieContext } from './movieFlowTypes';
import { openDB } from './db';

const RetrieverOptionsSchema = z.object({
  k: z.number().optional().default(10),
});

const QuerySchema = z.object({
  query: z.string(),
});

// Defining the Retriever
const sqlRetriever = defineRetriever(
  {
    name: 'movies',
    configSchema: RetrieverOptionsSchema,
  },
  async (query, options) => {
    const db = await openDB();
    if (!db) {
      throw new Error('Database connection failed');
    }
    const embedding = await embed({
      embedder: textEmbedding004,
      content: query,
    });
    const results = await db`
      SELECT content, title, poster, released, runtime_mins, rating, genres, director, actors, plot, tconst
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
    name: 'SemanticSearchFlow',
    inputSchema: QuerySchema,
    outputSchema: z.array(MovieContextSchema), // Array of MovieContextSchema
  },
  async (input) => {
    const docs = await retrieve({
      retriever: sqlRetriever,
      query: {
        content: [{ text: input.query }],
      },
      options: {
        k: 10,
      },
    });
    const movieContexts: MovieContext[] = [];

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
      }
    }

    return movieContexts;
  }
);