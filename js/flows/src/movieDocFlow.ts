
import { ai } from './genkitConfig'
import { z } from 'genkit';
import { MovieContextSchema, MovieContext } from './movieFlowTypes';
import { gemini15Flash } from '@genkit-ai/vertexai';
import { DocSearchFlowPromptText } from './prompts';
import { QuerySchema, SearchFlowOutputSchema, SearchFlowOutput} from './movieDocTypes';
import { sqlRetriever } from './docRetriever';

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
      vectorQuery: jsonResponse.vectorQuery || "",
      keywordQuery: jsonResponse.keywordQuery || "",
      searchCategory: jsonResponse.searchCategory || 'NONE',
      modelOutputMetadata: {
        justification: jsonResponse.justification || "",
        safetyIssue: jsonResponse.safetyIssue || false,
      },
    }

    return await getDocuments(searchFlowOutput);
  }
);
export async function getDocuments(searchFlowOutput: SearchFlowOutput) {
  const docs = await ai.retrieve({
    retriever: sqlRetriever,
    query: {
      content: [{ text: "" }],
    },
    options: {
      k: 10,
      searchCategory: searchFlowOutput.searchCategory,
      keywordQuery: searchFlowOutput.keywordQuery || "",
      vectorQuery: searchFlowOutput.vectorQuery || ""
    },
  });
  const movieContexts: MovieContext[] = [];

  for (const doc of docs) {
    if (doc.metadata) {
      const movieContext: MovieContext = {
        title: doc.metadata.title,
        runtime_minutes: doc.metadata.runtime_mins,
        genres: doc.metadata.genres.split(","),
        rating: parseFloat(parseFloat(doc.metadata.rating).toFixed(1)),
        plot: doc.metadata.plot,
        released: parseInt(doc.metadata.released, 10),
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

