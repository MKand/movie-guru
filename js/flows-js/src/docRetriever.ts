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
    //INTRUCTIONS:
    // 1. Create an embedding for the query
    // 2. Query the database 
    // 3. Create a document for each row (of type Document) 
    // 4. Return the content field from the row as content in the document and the remaining fields as metadata
    // 5. Return a list of documents.

		// Why content and metadata?
		// We separate movie data into 'content' and 'metadata' to accommodate varying approaches to data handling in GenAI frameworks.
		// Some frameworks, particularly those focused on RAG and utilizing a 'Document' object,
		// primarily use the 'content' field during RAG, potentially ignoring 'metadata'.

		// This separation is partly rooted in the historical context of these frameworks, which were often initially designed
		// to work with document-style databases rather than relational databases.
		// In document dbs, all the informational content is contained in the content of the document and not its metadata.
		// But in a relational db, the information may be spread across different columns.

		// In our application (using Genkit), we have the flexibility to pass a custom 'MovieContext' object into the RAG flow (next challenge) 
    // (and not restricted to document.content).
		// However, when interacting with other frameworks, especially those relying on a 'Document' structure,
		// it's crucial to be mindful of how metadata is utilized or if adjustments are needed to ensure all essential information is included.

    // Actually if you look at how MovieContext is constructed, we even throw away the content and only process the data in the metadata fields while constructing the MovieContext.

    return {
        // Return empty document list
        documents: [] as Document[],
    };
  }
);

export const movieDocFlow = defineFlow(
  {
    name: 'movieDocFlow',
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
          runtime_minutes: doc.metadata.runtimeMinutes,
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