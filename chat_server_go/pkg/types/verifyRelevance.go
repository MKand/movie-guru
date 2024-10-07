package types

// ConversationTurnOutcome represents the outcome of an agent's response in a conversation.
type ConversationTurnOutcome string

const (
	// OutcomeIrrelevant indicates the agent's response was not relevant to the user's previous turn.
	OutcomeIrrelevant ConversationTurnOutcome = "IRRELEVANT"

	// OutcomeAcknowledged indicates the user acknowledged the agent's response,
	// but it's unclear if they found it helpful or relevant.
	OutcomeAcknowledged ConversationTurnOutcome = "ACKNOWLEDGED"

	// OutcomeEngaged indicates the user is engaged and wants to continue the conversation
	// on the same topic (e.g., asking follow-up questions, requesting more information).
	OutcomeEngaged ConversationTurnOutcome = "ENGAGED"

	// OutcomeTopicChange indicates the user changed the topic of the conversation.
	OutcomeTopicChange ConversationTurnOutcome = "TOPIC_CHANGE"

	// OutcomeAmbiguous indicates the user's response is ambiguous and the outcome cannot be clearly determined.
	OutcomeAmbiguous ConversationTurnOutcome = "AMBIGUOUS"

	// OutcomeRejected indicates the user explicitly rejected the agent's response.
	OutcomeRejected ConversationTurnOutcome = "REJECTED"

	// OutcomeOther indicates an outcome that doesn't fit the other categories.
	OutcomeOther ConversationTurnOutcome = "OTHER"
)

// UserSentiment represents the sentiment expressed in the user's message.
type UserSentiment string

const (
	// SentimentPositive indicates a positive sentiment expressed by the user.
	SentimentPositive UserSentiment = "POSITIVE"

	// SentimentNegative indicates a negative sentiment expressed by the user.
	SentimentNegative UserSentiment = "NEGATIVE"

	// SentimentNeutral indicates a neutral sentiment expressed by the user.
	SentimentNeutral UserSentiment = "NEUTRAL"

	// SentimentAmbiguous indicates the sentiment in the user's message is ambiguous.
	SentimentAmbiguous UserSentiment = "AMBIGUOUS"
)

// ResponseQualityFlowInput represents the input to the response quality analysis flow.
type ResponseQualityFlowInput struct {
	MessageHistory []*types.SimpleMessage `json:"messageHistory"`
}

// ResponseQualityFlowOutput represents the output of the response quality analysis flow.
type ResponseQualityFlowOutput struct {
	Outcome       ConversationTurnOutcome `json:"outcome"`
	UserSentiment UserSentiment           `json:"userSentiment"` // Now using the UserSentiment type
	*ModelOutputMetadata
}

type ResponseQualityOutput struct {
	Outcome       ConversationTurnOutcome `json:"outcome"`
	UserSentiment UserSentiment           `json:"userSentiment"` // Now using the UserSentiment type
}

// NewResponseQualityFlowOutput creates a new ResponseQualityFlowOutput with default values.
func NewResponseQualityFlowOutput() *ResponseQualityFlowOutput {
	return &ResponseQualityFlowOutput{
		Outcome:       OutcomeOther,       // Default outcome
		UserSentiment: SentimentAmbiguous, // Default sentiment
	}
}
