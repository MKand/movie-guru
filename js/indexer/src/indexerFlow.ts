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

import { textEmbedding004 } from '@genkit-ai/vertexai';
import { toSql } from 'pgvector';
import { z } from 'genkit';
import { MovieContextSchema, MovieContext } from './types';
import { openDB } from './db';
import { ai } from './genkitConfig'

export const IndexerFlow = ai.defineFlow(
  {
    name: 'indexerFlow',
    inputSchema: MovieContextSchema,
    outputSchema: z.string(),
  },
  async (doc) => {
    const db = await openDB();
    if (!db) {
      throw new Error('Database connection failed');
    }
    try {
      // Reduce rate at which operation is performed to avoid hitting VertexAI rate limits
      await new Promise((resolve) => setTimeout(resolve, 300));
      const contentString = createText(doc);
      const eres = await ai.embed({
        embedder: textEmbedding004,
        content: contentString,
      });
      try {
        await db`
          INSERT INTO movies (embedding, title, runtime_mins, genres, rating, released, actors, director, plot, poster, tconst, content)
          VALUES (${toSql(eres)}, ${doc.title}, ${doc.runtimeMinutes}, ${doc.genres}, ${doc.rating}, ${doc.released}, ${doc.actors}, ${doc.director}, ${doc.plot}, ${doc.poster}, ${doc.tconst}, ${contentString})
          ON CONFLICT (tconst) DO UPDATE
          SET embedding = EXCLUDED.embedding
        `;
        return contentString;
      } catch (error) {
        console.error('Error inserting or updating movie:', error);
        throw error; // Re-throw the error to be handled by the outer try...catch
      }
    } catch (error) {
      console.error('Error indexing movie:', error);
      return 'Error indexing movie'; // Return an error message
    }
  }
);


function createText(movie: MovieContext): string {
  const dataDict = {
    title: movie.title,
    runtime_mins: movie.runtimeMinutes,
    genres: movie.genres.length > 0 ? movie.genres.join(', ') : '',
    rating: movie.rating > 0 ? movie.rating.toFixed(1) : '',
    released: movie.released > 0 ? movie.released : '',
    actors: movie.actors.length > 0 ? movie.actors.join(', ') : '',
    director: movie.director !== '' ? movie.director : '',
    plot: movie.plot !== '' ? movie.plot : '',
  };

  const jsonData = JSON.stringify(dataDict);
  return jsonData;
}
