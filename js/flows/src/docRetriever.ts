import { Document } from '@genkit-ai/ai/retriever';
import { textEmbedding004 } from '@genkit-ai/vertexai';
import { toSql } from 'pgvector';
import { openDB } from './db';
import { ai } from './genkitConfig'
import { z } from 'genkit';
import { MovieContextSchema, MovieContext } from './movieFlowTypes';
import { gemini15Flash } from '@genkit-ai/vertexai';
import { DocSearchFlowPromptText } from './prompts';
import { ModelOutputMetadata, ModelOutputMetadataSchema } from './modelOutputMetadataTypes';

const SearchTypeCategory = z.enum(['KEYWORD', 'VECTOR', 'NONE']);


export const RetrieverOptionsSchema = z.object({
  k: z.number().optional().default(10),
  searchCategory: SearchTypeCategory.optional().default("VECTOR")
});

export const QuerySchema = z.object({
  query: z.string(),
});

export const SearchFlowOutputSchema = z.object({
  outputQuery: z.string().optional(),
  searchCategory: SearchTypeCategory,
  modelOutputMetadata: ModelOutputMetadataSchema,
});

export const SearchFlowPrompt = ai.definePrompt(
  {
    name: 'MixedSearchFlowPrompt',
    model: gemini15Flash,
    input: {
      schema: QuerySchema,
    },
    output: {
      format: 'json',
      schema: SearchFlowOutputSchema,
    },  
  }, 
  DocSearchFlowPromptText
)

export const MovieDocFlow = ai.defineFlow(
  {
    name: 'movieDocFlow',
    inputSchema: QuerySchema,
    outputSchema: z.array(MovieContextSchema), // Array of MovieContextSchema
  },
  async (input) => {
    const response = await SearchFlowPrompt( {
      query: input.query
    })
    if (typeof response.text !== 'string') {
      throw new Error('Invalid response format: text property is not a string.');
    }
    const jsonResponse = JSON.parse(response.text)
    const searchFlowOutput = {
      outputQuery: jsonResponse.outputQuery || "",
      searchCategory: jsonResponse.searchCategory || 'VECTOR',
      modelOutputMetadata: {
        justification: jsonResponse.justification || "",
        safetyIssue: jsonResponse.safetyIssue || false,
      },
    }

    const docs = await ai.retrieve({
      retriever: sqlRetriever,
      query: {
        content: [{ text: searchFlowOutput.outputQuery }],
      },
      options: {
        k: 10,
        searchCategory: searchFlowOutput.searchCategory
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

    let results;

    if(options.searchCategory == "KEYWORD"){
      const sqlQuery = `SELECT content, title, poster, released, runtime_mins, rating, genres, director, actors, plot, tconst
      FROM movies
      WHERE ${query.text}
      LIMIT ${options.k ?? 10}`
      results = await db.unsafe(sqlQuery)
    }

     //Vector Query
     if(options.searchCategory == "VECTOR"){
      const embedding = await ai.embed({
        embedder: textEmbedding004,
        content: query,
      });  
        results = await db`
        SELECT content, title, poster, released, runtime_mins, rating, genres, director, actors, plot, tconst
       FROM movies
          ORDER BY embedding <#> ${toSql(embedding)}
          LIMIT ${options.k ?? 10}
        ;`
    }
    if (!results) {
      throw new Error('No results found.'); 
    }  
    return {
      documents: results.map((row) => {
        const { content, ...metadata } = row;
        return Document.fromText(content, metadata);
      }),
    };
  }
);

