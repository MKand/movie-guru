import { embed } from '@genkit-ai/ai/embedder';
import { textEmbedding004 } from '@genkit-ai/vertexai';
import { defineFlow } from '@genkit-ai/flow';
import { toSql } from 'pgvector';
import { z } from 'zod';
import { MovieContextSchema, MovieContext } from './types';
import { openDB } from './db';


export const IndexerFlow = defineFlow(
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
        // Adding a delay to overcome rate limit issues
        await new Promise((resolve) => setTimeout(resolve, 300));

        //Create string for embedding
        const embeddedContent = createTextForEmbedding(doc);

        // Create Embedding
        const embedding = await embed({
          embedder: textEmbedding004,
          content: embeddedContent,
        });

        //Insert embedding and other columns into DB
        try {
          await db`
          INSERT INTO movies (embedding, title, runtime_mins, genres, rating, released, actors, director, plot, poster, tconst, content)
          VALUES (${toSql(embedding)}, ${doc.title}, ${doc.runtimeMinutes}, ${doc.genres}, ${doc.rating}, ${doc.released}, ${doc.actors}, ${doc.director}, ${doc.plot}, ${doc.poster}, ${doc.tconst}, ${embeddedContent})
          ON CONFLICT (tconst) DO UPDATE
          SET embedding = EXCLUDED.embedding
          `;
          return embeddedContent; 
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

   function createTextForEmbedding(movie: MovieContext): string {
    // What fields in the text are useful to create an embedding from?

    const dataDict = {
      title: movie.title,
      runtime_mins: movie.runtimeMinutes,
      genres: movie.genres.length > 0 ? movie.genres.join(', ') : '',
      rating: movie.rating > 0 ? movie.rating.toFixed(1) : '',
      released: movie.released > 0 ? movie.released : '',
      actors: movie.actors.length > 0 ? movie.actors.join(', ') : '',
      director: movie.director !== '' ? movie.director : '',
      plot: movie.plot !== '' ? movie.plot : '',
      tconst: movie.tconst,
      poster: movie.poster
    };
  
    const jsonData = JSON.stringify(dataDict);
    return jsonData;
  }