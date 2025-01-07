import { ai } from './genkitConfig';
import { 
  MovieFlowOutputSchema,
  MovieContext,
  MovieFlowOutput
} from './movieFlowTypes';
import { MovieFlowPrompt } from './movieFlow';
import {  SearchFlowPrompt , getDocuments} from './movieDocFlow';
import { z } from 'genkit';
import { SimpleMessageSchema } from './queryTransformTypes'; 

export const RAGFlowInputSchema = z.object({
    history: z.array(SimpleMessageSchema),
    userMessage: z.string(),
  });

export const MovieRAGFlow = ai.defineFlow(
  {
    name: 'RAGFlow',
    inputSchema: RAGFlowInputSchema,
    outputSchema: MovieFlowOutputSchema
  },

  async (input) => {
    try {

      const searchResponse = await SearchFlowPrompt( {
        query: input.userMessage
      })

      const searchResponseOutput = JSON.parse(searchResponse.text)
      const searchFlowOutput = {
        vectorQuery: searchResponseOutput.vectorQuery || "",
        keywordQuery: searchResponseOutput.keywordQuery || "",
        searchCategory: searchResponseOutput.searchCategory || 'NONE',
        modelOutputMetadata: {
          justification: searchResponseOutput.justification || "",
          safetyIssue: searchResponseOutput.safetyIssue || false,
        },
      }
  
      const movieContexts: MovieContext[] = await getDocuments(searchFlowOutput);

   
      let mfInput = {
        history: input.history,
        userProfile: {},
        userMessage: input.userMessage,
        contextDocuments: movieContexts
      }
      const response = await MovieFlowPrompt(mfInput)
      const jsonResponse =  JSON.parse(response.text);
      const output: MovieFlowOutput = {
        "answer":  jsonResponse.answer,
        "relevantMovies": jsonResponse.relevantMovies,
        "wrongQuery": jsonResponse.wrongQuery,
        "modelOutputMetadata": {
          "justification": jsonResponse.justification,
          "safetyIssue": jsonResponse.safetyIssue,
        }
      }
      return output    
    } catch (error) {
      console.error("Error generating response:", error);
      return { 
        relevantMovies: [],
        answer: "",
        modelOutputMetadata: {
            "justification": "",
            "safetyIssue": false,
          }      
        }; 
    }
  }
);