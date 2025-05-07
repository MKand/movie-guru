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

import { Document } from '@genkit-ai/ai/retriever';
import { textEmbedding005 } from '@genkit-ai/vertexai';
import { toSql } from 'pgvector';
import { openDB } from './db';
import { ai } from './genkitConfig'
import { z } from 'genkit';
import { MovieContextSchema, MovieContext } from './movieFlowTypes';
import {  ModelOutputMetadataSchema } from './modelOutputMetadataTypes';

const SearchTypeCategory = z.enum(['KEYWORD', 'VECTOR', 'MIXED', 'NONE']);

export const RetrieverOptionsSchema = z.object({
  k: z.number().optional().default(10),
  searchCategory: SearchTypeCategory.default("VECTOR"),
  keywordQuery: z.string().default(""),
  vectorQuery: z.string().default(""),
});

export const QuerySchema = z.object({
  query: z.string(),
});

ai.defineSchema('QuerySchema', QuerySchema);

export const SearchFlowOutputSchema = z.strictObject({
  keywordQuery: z.string().optional().default(""),
  vectorQuery: z.string().optional().default(""),
  searchCategory: SearchTypeCategory.default("NONE"),
  modelOutputMetadata: ModelOutputMetadataSchema.default(ModelOutputMetadataSchema.parse({})),
});

ai.defineSchema('SearchFlowOutputSchema', SearchFlowOutputSchema);

/**
 * Prompt file: js/flows/prompts/docSearch.prompt
 * 
 * This prompt takes the generated search terms from the queryTransformFlow and uses that to retrieve relevant documents
 * from the database.
 * 
 * Input schema: QuerySchema
 * Output schema: SearchFlowOutputSchema
 * 
 * The MovieGuru development team, uses a variant system to version our prompts. The "v2" variant corresponds to the docSearch.v2.prompt file. 
 * To use the default variant -- ai.prompt('docSearch')
 * To use a variant -- ai.prompt('docSearch', {variant: 'v2'})
 * 
 * ATTENTION: SREs will test variant 'v2' to see if we can eliminate keyword search.
 */
export const searchForRelevantMovies = ai.prompt('docSearch');

export const DocSearchFlow = ai.defineFlow(
  {
    name: 'docSearchFlow',
    inputSchema: QuerySchema,
    outputSchema: z.array(MovieContextSchema),
  },
  async (input) => {
    const movieContexts: MovieContext[] = [];
    const searchFlowOutput = await createSearchObject(input);
  try{
    if (searchFlowOutput.searchCategory == "NONE" && (searchFlowOutput.keywordQuery == "" && searchFlowOutput.vectorQuery == "")){
      return movieContexts;
    }
    const docs = await ai.retrieve({
      retriever: sqlRetriever,
      query: {
        content: [{ text: "" }],
      },
      options: {
        k: 10,
        searchCategory: searchFlowOutput.searchCategory,
        keywordQuery: searchFlowOutput.keywordQuery,
        vectorQuery: searchFlowOutput.vectorQuery
      },
    });

    for (const doc of docs) {
      if (doc.metadata) {
        const movieContext: MovieContext = {
          title: doc.metadata.title,
          runtime_minutes: doc.metadata.runtime_mins,
          genres: doc.metadata.genres.split(","),
          rating: parseFloat(parseFloat(doc.metadata.rating).toFixed(1)),
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
  catch(e){
    console.error(`Retriever: Unable to get documents: ${e instanceof Error ? e.message : e}`)
    return movieContexts;
  }
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
    if(options.searchCategory == "KEYWORD" || options.keywordQuery != ""){
      results =  await db`SELECT content, title, poster, released, runtime_mins, rating, genres, director, actors, plot, tconst
      FROM movies
      WHERE ${db.unsafe(options.keywordQuery)} 
      LIMIT ${options.k ?? 10}`
    }

     //Vector Query
     if(options.searchCategory == "VECTOR" || options.vectorQuery != ""){
      const embedding = await ai.embed({
        embedder: textEmbedding005,
        content: options.vectorQuery,
      });  
        results = await db`
        SELECT content, title, poster, released, runtime_mins, rating, genres, director, actors, plot, tconst
       FROM movies
          ORDER BY embedding <#> ${toSql(embedding[0].embedding)}
          LIMIT ${options.k ?? 10}
        ;`
    }

    //Mixed Query
    if (options.searchCategory === "MIXED") {
      // Generate the vector embedding for the vector query
      const embedding = await ai.embed({
        embedder: textEmbedding005,
        content: options.vectorQuery,
      });
    
      // Execute the database query with both keyword and vector search components
      results = await db`
        SELECT 
          content, 
          title, 
          poster, 
          released, 
          runtime_mins, 
          rating, 
          genres, 
          director, 
          actors, 
          plot, 
          tconst
        FROM 
          movies
        WHERE 
        ${db.unsafe(options.keywordQuery)} 
        ORDER BY 
          embedding <#> ${toSql(embedding)}
        LIMIT 
          ${options.k ?? 10}
      ;`;
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
async function createSearchObject(input: { query: string; }) {
  const defaultOutput = SearchFlowOutputSchema.parse({});
  try {
    const response = await searchForRelevantMovies({
      query: input.query
    });
    const safeOutput = response.output ?? SearchFlowOutputSchema.parse({});
    return SearchFlowOutputSchema.parse(safeOutput);
  }
  catch (error) {
    console.error('MovieDocFlow: Error generating response:', {
      error,
      input,
    });
    return defaultOutput;
  }
}

