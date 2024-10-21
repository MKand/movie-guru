import { embed } from '@genkit-ai/ai/embedder';
import { Document, defineRetriever, retrieve } from '@genkit-ai/ai/retriever';
import { defineFlow } from '@genkit-ai/flow';
import { textEmbedding004 } from '@genkit-ai/vertexai';
import { toSql } from 'pgvector';
import { z } from 'zod';
import { MovieContextSchema, MovieContext } from './movieFlowTypes';
import { openDB } from './db';
import { gemini15Flash } from '@genkit-ai/vertexai';
import { defineDotprompt } from '@genkit-ai/dotprompt'
import { SearchFlowInputSchema } from './vectorSearchFlow';

const SearchTypeCategory = z.enum(['KEYWORD', 'VECTOR', 'NONE']);

const RetrieverOptionsSchema = z.object({
  k: z.number().optional().default(10),
});

export const SearchFlowOutputSchema = z.object({
  outputQuery: z.string().optional(),
  searchCategory: SearchTypeCategory,
  Justification: z.string(),
});

export const MixedSearchFlowPrompt = defineDotprompt(
  {
    name: 'SearchFlowPrompt',
    model: gemini15Flash,
    input: {
      schema: SearchFlowInputSchema,
    },
    output: {
      format: 'json',
      schema: SearchFlowOutputSchema,
    },  
  }, 
    
    `Analyze the query string: "{{inputQuery}}" with respect to a movie database with these fields:

    *  embedding: Vector representation of the movie's title, plot, and genres.
    *  genres: List of genres (e.g., "Action", "Comedy", "Drama").
    *  title: Title of the movie.
    *  plot: A textual summary of the movie's plot.
    *  runtime_mins: Duration of the movie in minutes.
    *  released: Release year of the movie.
    *  actors: List of actors in the movie.
    *  director: Director of the movie.
    *  rating: Numerical rating from 1 to 5.

    Determine if the query is best satisfied by a **KEYWORD** or **VECTOR** search.

    **KEYWORD search:** Use for queries that can be expressed with simple SQL operators (=, !=, >, <, IN) on the title, actors, director, genres, runtime_mins, or released fields.

    *   Do not include the WHERE keyword in the output query.
    *   Some user queries might need to be transformed (for KEYWORD search). Where this is necessary is:
       - User is asking for movies based on their lengths, ratings, quality or recency. 
            Before classifying, apply these transformations to the inputQuery:
            *   Movie quality:
                *   Bad: rating < 2
                *   Good: rating > 3.5
                *   Great: rating > 4.5
                *   Terrible: rating < 1
                *   Average: rating BETWEEN 2 AND 3.5
            *   Movie length:
                *   short: runtime_mins < 20
                *   long: runtime_mins > 120
                *   very long: runtime_mins > 150
            *   Movie year:
                *   recent: released > 2020
                *   old: released < 2005
            Examples:
              inputQuery: "great movie that is short" 
              outputQuery: "rating > 4.5 AND runtime_mins < 20" 
    *   Examples:
        *   inputQuery: "movie with a rating higher than 3". outputQuery: "rating > 3"
        *   inputQuery: "movies with actress Tilda Swinton". outputQuery: "'Tilda Swinton' IN actors" 
        *   inputQuery: "movies released after 2000". outputQuery: "released > 2000"

    **VECTOR search:** Use for queries requiring semantic understanding of the title, plot, or genres fields. This includes:

    *   Queries about concepts, emotions, or themes.
    *   Queries matching analogies or metaphors (e.g., "movies that make you cry").
    *   Important: Any query that would require the LIKE operator in SQL on the title, plot, or genres fields should be classified as a **VECTOR** search.
    *   Return a more concise form of the inputQuery as the outputQuery
    *   Examples:
        *   inputQuery: "movies with strong female leads" . outputQuery: "strong female leads" 
        *   inputQuery: "movies with location names in their titles". outputQuery: "location names in their titles"
        *   inputQuery: "find movies like The Matrix". outputQuery: "like The Matrix"
   `
)


export const mixedSearchFlow = defineFlow(
  {
    name: 'MixedSearchFlow',
    inputSchema: SearchFlowInputSchema,
    outputSchema: z.array(MovieContextSchema)
  },
  async (input) => {
    const response = await MixedSearchFlowPrompt.generate({input : input})
    const searchFlowOutput = response.output(0)
    const movieContexts: MovieContext[] = [];
    // if (searchFlowOutput.outputQuery == "" ||  null || "null"){
    //   searchFlowOutput.outputQuery = input.inputQuery
    // }
    let docs: Document[] = []
      if (searchFlowOutput.searchCategory == "VECTOR"){
      docs = await retrieve({
        retriever: vectorRetriever,
        query: {
          content: [{ text: searchFlowOutput.outputQuery }],
        },
        options: {
          k: 10,
        },
      });
    }
    if (searchFlowOutput.searchCategory == "KEYWORD"){
      console.log("doing keyword search ", searchFlowOutput)
      docs = await retrieve({
        retriever: keywordRetriever,
        query: {
          content: [{ text: searchFlowOutput.outputQuery }],
        },
        options: {
          k: 10,
        },
      });
    }
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
    console.log("query text", query.text())
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

// Defining the keyword Retriever
export const keywordRetriever = defineRetriever(
  {
    name: 'keywordRetriever',
    configSchema: RetrieverOptionsSchema,
  },
  async (query, options) => {
    const db = await openDB();
    if (!db) {
      throw new Error('Database connection failed');
    }
    const sqlQuery = `SELECT content, title, poster, released, runtime_mins, rating, genres, director, actors, plot, tconst
       FROM movies
       WHERE ${query.text()}
       LIMIT ${options.k ?? 10}`
    
    const results = await db.unsafe(sqlQuery)
    return {
      documents: results.map((row) => {
        const { content, ...metadata } = row;
        return Document.fromText(content, metadata);
      }),
    };
  }
);