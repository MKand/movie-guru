import {
    USERINTENT,
    QueryTransformFlowOutput
  } from './queryTransformTypes';
import { ChatFlowInputSchema, ChatFlowOutput, ChatOutputSchema } from './chatFlowTypes';
import { MovieContext, RelevantMovie, RelevantMovieSchema } from './movieFlowTypes';

import { ai } from './genkitConfig';
import { QueryTransformFlow } from './queryTransformFlow';
import { MovieDocFlow } from './docRetriever';
import { MovieFlow } from './movieFlow';

export const ChatFlow = ai.defineFlow(
    {
        name: "chatFlow",
        inputSchema: ChatFlowInputSchema,
        outputSchema: ChatOutputSchema,
    },
    async(input) => {
        try{
            const chatResponse: ChatFlowOutput = ChatOutputSchema.parse({})
            const qtResponse: QueryTransformFlowOutput = await QueryTransformFlow(input)
            if (qtResponse.modelOutputMetadata.safetyIssue || qtResponse.modelOutputMetadata.quotaIssue){
                 chatResponse.modelOutputMetadata = qtResponse.modelOutputMetadata
                 return chatResponse
            }
            var movieContexts: MovieContext[] = []
            if (qtResponse.userIntent==USERINTENT.parse("REQUEST" ) || qtResponse.userIntent==USERINTENT.parse("RESPONSE")){
                 movieContexts = await MovieDocFlow( {query: qtResponse.transformedQuery})
            }
            const movieFlowResponse = await MovieFlow({
                history: input.history,
                userPreferences: input.userPreferences,
                contextDocuments: movieContexts,
                userMessage: input.userMessage
            })
            
            chatResponse.answer = movieFlowResponse.answer;
            chatResponse.modelOutputMetadata = movieFlowResponse.modelOutputMetadata
            chatResponse.relevantMovies = movieFlowResponse.relevantMovies
            chatResponse.wrongQuery = movieFlowResponse.wrongQuery
            chatResponse.relevantMovies = movieFlowResponse.relevantMovies
            chatResponse.contextDocuments = parseContexts(movieFlowResponse.relevantMovies, movieContexts)
            
            return chatResponse
        }
        catch (error) {
            console.error("ChatFlow: Error generating response:", error);
            throw error;
        }
    }
)

function parseContexts(relevantMovies: RelevantMovie [], movieContexts:MovieContext[]): MovieContext[]{
   const relevantContext: MovieContext[] = []
   for (const r of relevantMovies){
        for(const c of movieContexts){
            if(r.title  == c.title){
                relevantContext.push(c)
            }
        }
   }
   return relevantContext

}
