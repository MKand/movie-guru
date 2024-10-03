import { embed } from '@genkit-ai/ai/embedder';
import { textEmbeddingGecko } from '@genkit-ai/vertexai';
import { defineFlow } from '@genkit-ai/flow';
import { toSql } from 'pgvector';
import { z } from 'zod';
import { MovieContextSchema, MovieContext } from './movieFlowTypes';
import { OpenDB } from './db';


export const IndexerFlow = defineFlow(
{
    name: 'indexerFlow',
    inputSchema: MovieContextSchema,
    outputSchema: z.string(),
},
  async (doc) => {
    const db = await OpenDB();
    if (!db) {
      throw new Error('Database connection failed');
    }
    try {
      // sleep for 300 ms
      await new Promise((resolve) => setTimeout(resolve, 300));
      const contentString = createText(doc);
      const eres = await embed({
        embedder: textEmbeddingGecko,
        content: contentString,
      });
      try {
        await db`
        INSERT INTO test_movies (embedding, title, runtime_mins, genres, rating, released, actors, director, plot, poster, tconst, content)
        VALUES (${toSql(eres)}, ${doc.title}, ${doc.runtimeMinutes}, ${doc.genres.join(', ')}, ${doc.rating}, ${doc.released}, ${doc.actors.join(', ')}, ${doc.director}, ${doc.plot}, ${doc.poster}, ${doc.tconst}, ${contentString})
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
    plot: movie.plot !== '' ? movie.plot.replace(/\n/g, '') : '',
  };

  const jsonData = JSON.stringify(dataDict);
  return jsonData;
}
