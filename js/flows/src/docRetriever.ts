import { Document } from '@genkit-ai/ai/retriever';
import { textEmbedding004 } from '@genkit-ai/vertexai';
import { toSql } from 'pgvector';
import { openDB } from './db';
import { ai } from './genkitConfig'
import {RetrieverOptionsSchema} from './movieDocTypes'

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

    if(options.searchCategory == "NONE"){
      return {
        documents: [],
      };
    }

    if(options.searchCategory == "KEYWORD"){
      results =  await db`SELECT content, title, poster, released, runtime_mins, rating, genres, director, actors, plot, tconst
      FROM movies
      WHERE ${db.unsafe(options.keywordQuery)} 
      LIMIT ${options.k ?? 10}`
    }

     //Vector Query
     if(options.searchCategory == "VECTOR"){
      const embedding = await ai.embed({
        embedder: textEmbedding004,
        content: options.vectorQuery,
      });  
        results = await db`
        SELECT content, title, poster, released, runtime_mins, rating, genres, director, actors, plot, tconst
       FROM movies
          ORDER BY embedding <#> ${toSql(embedding)}
          LIMIT ${options.k ?? 10}
        ;`
    }

    //Mixed Query
    if (options.searchCategory === "MIXED") {
      // Generate the vector embedding for the vector query
      const embedding = await ai.embed({
        embedder: textEmbedding004,
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

