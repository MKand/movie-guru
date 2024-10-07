package main

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/invopop/jsonschema"
)

// GetResponseQualityAnalysisFlow creates a GenKit flow for analyzing response quality.
func GetResponseQualityAnalysisFlow(ctx context.Context, model ai.Model) (*genkit.Flow[*types.ResponseQualityFlowInput, *types.ResponseQualityFlowOutput, struct{}], error) {

	// Define the prompt using dotprompt
	responseQualityPrompt, err := dotprompt.Define("responseQualityAnalysisFlow",
		`
		You are an AI assistant designed to analyze conversations between users and a movie expert agent. 
		Your task is to objectively assess the flow of the conversation and determine the outcome of the agent's response based solely on the user's reaction to it.
		You also need to determine the user's sentiment based on their last message (it can be positive, negative, neutral, or ambiguous).
		You only get a truncated version of the conversation history.

		Here's how to analyze the conversation:

		1. Read the conversation history carefully, paying attention to the sequence of messages and the topics discussed.
		2. Focus on the agent's response and how the user reacts to it.

		Guidelines for classification of the conversation outcome:

		*   OutcomeIrrelevant: The agent's response is not connected to the user's previous turn or doesn't address the user's query or request.
		*   OutcomeAcknowledged: The user acknowledges the agent's response with neutral remarks like "Okay," "Got it," or a simple "Thanks" without indicating further interest or engagement.
		*   OutcomeRejected: The user responds negatively to the agent's response like "No," "I don't like it," or a simple "No thanks" without indicating further interest or engagement.
		*   OutcomeEngaged: The user shows interest in the agent's response and wants to delve deeper into the topic. This could be through follow-up questions, requests for more details, or expressing a desire to learn more about the movie or topic mentioned by the agent.
		*   OutcomeTopicChange: The user shifts the conversation to a new topic unrelated to the agent's response.
		*   OutcomeAmbiguous: The user's response is too vague or open-ended to determine the outcome with certainty.
		*   OutcomeOther: The user's response doesn't fit into any of the above categories. You can use this if the user's message is the only one in the history.

		Examples:

		User: "I'm looking for a movie with strong female characters."
		Agent: "Have you seen 'Alien'?"
		User: "Tell me more about it."
		Outcome: OutcomeEngaged (The user shows interest in the agent's suggestion and wants to learn more.)

		Agent: "Let me tell you about the movie 'Alien'?"
		User: "I hate that film"
		Outcome: OutcomeRejected (The user rejects the agent's suggestion.)

		Agent: "Have you seen 'Alien'?"
		User: "No. Tell me about 'Princess diaries'"
		Outcome: OutcomeTopicChange (The user shows no interest in the agent's suggestion and changes the topic.)

		Agent: "Have you seen 'Alien'?"
		User: "I told you I am not interested in sci-fi."
		Outcome: OutcomeIrrelevant (The agent made a wrong suggestion.)

		Provide a brief explanation for your classification based solely on the user's following turn.

		Guidelines for classification of the user sentiment:
		* Positive: If the user expresses excitement, joy etc. Simply rejecting an agent's suggestion is not negative.
		* Negative: If the user expresses frustration, irritation, anger etc. Simply rejecting an agent's suggestion is not negative.
		* Neutral: If the user expresses no specific emotion
		* Ambiguous: If the user sentiment is not clear.


		Remember:

		*   Do not make assumptions about the user's satisfaction or perception of helpfulness.
		*   Focus only on the objective flow of the conversation and how the user's response relates to the agent's previous turn.
		*   If the outcome is unclear based on the user's response, use OutcomeAmbiguous.

		Here are the inputs:
		* Conversation History (this is a truncated version and also may only have a single message if the user just started the conversation):
			{{messageHistory}}
		`,

		dotprompt.Config{
			Model:        model,
			InputSchema:  jsonschema.Reflect(types.ResponseQualityFlowInput{}),
			OutputSchema: jsonschema.Reflect(types.ResponseQualityFlowOutput{}),
			OutputFormat: ai.OutputFormatJSON,
			GenerationConfig: &ai.GenerationCommonConfig{
				Temperature: 0.5,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to define prompt: %w", err)
	}

	// Define the GenKit flow
	responseQualityFlow := genkit.DefineFlow("responseQualityFlow", func(ctx context.Context, input *ResponseQualityFlowInput) (*ResponseQualityFlowOutput, error) {
		output := types.NewResponseQualityFlowOutput() // More concise variable name

		resp, err := responseQualityPrompt.Generate(ctx,
			&dotprompt.PromptRequest{
				Variables: input,
			},
			nil,
		)
		if err != nil {
			if blockedErr, ok := err.(*genai.BlockedError); ok {
				fmt.Println("Request was blocked:", blockedErr)
				return output, nil
			}
			return nil, fmt.Errorf("failed to generate response: %w", err)
		}

		err = json.Unmarshal([]byte(resp.Text()), &output)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return output, nil
	})

	return responseQualityFlow, nil
}
