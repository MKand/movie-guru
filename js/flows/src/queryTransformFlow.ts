import { gemini15Flash } from '@genkit-ai/vertexai';
import {
  USERINTENT,
  QueryTransformFlowInputSchema,
  QueryTransformFlowOutputSchema,
} from './queryTransformTypes';
import { QueryTransformPromptText } from './prompts';
import { ai } from './genkitConfig';

export const QueryTransformPrompt = ai.definePrompt(
  {
    name: 'queryTransformFlowPrompt',
    model: gemini15Flash,
    input: {
      schema: QueryTransformFlowInputSchema,
    },
    output: {
      format: 'json',
    },
  },
  QueryTransformPromptText
);

export const QueryTransformFlow = ai.defineFlow(
  {
    name: 'queryTransformFlow',
    inputSchema: QueryTransformFlowInputSchema,
    outputSchema: QueryTransformFlowOutputSchema,
  },
  async (input) => {
    try {
      const response = await QueryTransformPrompt({
        history: input.history,
        userMessage: input.userMessage,
        userProfile: input.userProfile,
      });

      if (typeof response.text !== 'string') {
        throw new Error('Invalid response format: text property is not a string.');
      }

      const jsonResponse = JSON.parse(response.text)
      return {
        transformedQuery: jsonResponse.transformedQuery,
        userIntent: jsonResponse.userIntent || 'UNCLEAR',
        modelOutputMetadata: {
          justification: jsonResponse.justification,
          safetyIssue: jsonResponse.safetyIssue,
        },
      };
    } catch (error) {
      console.error('Error generating response:', {
        error,
        input,
      });

      // Return fallback response
      return {
        transformedQuery: '',
        userIntent: 'UNCLEAR',
        modelOutputMetadata: {
          justification: '',
          safetyIssue: false,
        },
      };
    }
  }
);
