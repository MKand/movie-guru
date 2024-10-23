import { defineFlow } from '@genkit-ai/flow';
import { 
  MovieFlowInputSchema, 
  MovieFlowOutputSchema,
  MovieContext 
} from './movieFlowTypes';
import { MovieFlowPrompt } from './movieFlow';
import { Document, retrieve } from '@genkit-ai/ai/retriever';
import { MixedSearchFlowPrompt, mixedRetriever } from './mixedSearchFlow';

import {  QueryTransformPrompt } from './queryTransformFlow';

export const MovieRAGFlow = defineFlow(
  {
    name: 'RAGFlow',
    inputSchema: MovieFlowInputSchema,
    outputSchema: MovieFlowOutputSchema
  },
  async (input) => {
    try {
      let qtInput = {
        history: input.history,
        userProfile: {},
        userMessage: input.userMessage,
      }
     
      const qtResponse = await QueryTransformPrompt.generate({input: qtInput});
      
      const searchResponse = await MixedSearchFlowPrompt.generate({ 
        input: {
          inputQuery: qtResponse.output(0).transformedQuery}
      });

      const searchResponseOutput = searchResponse.output(0)
      console.log("searchResponse.output(0) ", searchResponseOutput)
      let docs: Document[] = []
      const movieContexts: MovieContext[] = [];
      if(searchResponseOutput.searchCategory!= "NONE"){
        docs = await retrieve({
          retriever: mixedRetriever,
          query: searchResponseOutput.outputQuery,
          options: {
            k: 10,
            searchCategory: searchResponseOutput.searchCategory
          },
        });
  
        for (const doc of docs) {
          if (doc.metadata) {
            const movieContext: MovieContext = {
              title: doc.metadata.title,
              runtime_minutes: doc.metadata.runtime_mins,
              genres: doc.metadata.genres.split(","),
              rating: parseFloat(doc.metadata.rating),
              plot: doc.metadata.plot,
              released: parseFloat(doc.metadata.released),
              director: doc.metadata.director,
              actors: doc.metadata.actors.split(","),
              poster: doc.metadata.poster,
              tconst: doc.metadata.tconst,
            };
            movieContexts.push(movieContext);
          } 
        }
      }
      let mfInput = {
        history: input.history,
        userProfile: {},
        userMessage: input.userMessage,
        contextDocuments: movieContexts
      }
      const response = await MovieFlowPrompt.generate({ input : mfInput})
      return response.output(0);
    } catch (error) {
      console.error("Error generating response:", error);
      return { 
        relevantMovies: [],
        answer: "",
        justification: ""
      }; 
    }
  }
);