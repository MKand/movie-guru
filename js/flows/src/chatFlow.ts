import {
    USERINTENT,
    ChatFlowInputSchema,
    QueryTransformFlowOutputSchema,
    QueryTransformFlowOutput
  } from './queryTransformTypes';
import { ChatFlowOutputSchema, ChatFlowOutput, MovieContextSchema, MovieContext } from './movieFlowTypes';

import { ai, safetySettings } from './genkitConfig';
import { GenerationBlockedError } from 'genkit';
import { QueryTransformFlow, QueryTransformPrompt } from './queryTransformFlow';
import { MovieDocFlow } from './docRetriever';
import { MovieFlow } from './movieFlow';

export const ChatFlow = ai.defineFlow(
    {
        name: "ChatFlow",
        inputSchema: ChatFlowInputSchema,
        outputSchema: ChatFlowOutputSchema,
    },
    async(input) => {
        try{
            var chatResponse: ChatFlowOutput = ChatFlowOutputSchema.parse({})
            const qtResponse: QueryTransformFlowOutput = await QueryTransformFlow(input)
            if (qtResponse.modelOutputMetadata.safetyIssue || qtResponse.modelOutputMetadata.quotaIssue){
                 chatResponse.modelOutputMetadata = qtResponse.modelOutputMetadata
                 return chatResponse
            }
            var movieContexts: MovieContext[] = []
            if (qtResponse.userIntent=="REQUEST" || qtResponse.userIntent=="RESPONSE"){
                 movieContexts = await MovieDocFlow( {query: qtResponse.transformedQuery})
            }
            chatResponse = await MovieFlow({
                history: input.history,
                userPreferences: input.userPreferences,
                contextDocuments: movieContexts,
                userMessage: input.userMessage
            })
            return chatResponse
        }
        catch (error) {
            console.error("ChatFlow: Error generating response:", error);
            throw error;
        }
    }
)
