import { Document } from '@genkit-ai/ai/retriever';
import { textEmbedding004 } from '@genkit-ai/vertexai';
import { toSql } from 'pgvector';
import { openDB } from './db';
import { ai } from './genkitConfig'
import { z } from 'genkit';
import { MovieContextSchema, MovieContext } from './movieFlowTypes';

const RetrieverOptionsSchema = z.object({
  k: z.number().optional().default(10),
});

const QuerySchema = z.object({
  query: z.string(),
});

export const sqlRetriever = ai.defineRetriever(
  {
    name: 'movies',
    configSchema: RetrieverOptionsSchema,
  },
  async (query, options) => {
    const db = await openDB();
    if (!db) {
      throw new Error('Database connection failed');
    }
    const embedding = await ai.embed({
      embedder: textEmbedding004,
      content: query,
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

export const MovieDocFlow = ai.defineFlow(
  {
    name: 'movieDocFlow',
    inputSchema: QuerySchema,
    outputSchema: z.array(MovieContextSchema), // Array of MovieContextSchema
  },
  async (input) => {
    const docs = await ai.retrieve({
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
          genres: doc.metadata.genres.split(","),
          rating: parseInt(doc.metadata.rating, 10),
          plot: doc.metadata.plot,
          released: parseInt(doc.metadata.released,10),
          director: doc.metadata.director,
          actors: doc.metadata.actors.split(","),
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