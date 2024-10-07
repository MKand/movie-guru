package types

type USERINTENT string

const (
	UNCLEAR          USERINTENT = "UNCLEAR"
	GREET            USERINTENT = "GREET"
	END_CONVERSATION USERINTENT = "END_CONVERSATION"
	REQUEST          USERINTENT = "REQUEST"
	RESPONSE         USERINTENT = "RESPONSE"
	ACKNOWLEDGE      USERINTENT = "ACKNOWLEDGE"
)

type QueryTransformFlowOutput struct {
	TransformedQuery string     `json:"transformedQuery, omitempty"`
	Intent           USERINTENT `json:"userIntent, omitempty"`
	*ModelOutputMetadata
}

type QueryTransformFlowInput struct {
	History     []*SimpleMessage `json:"history"`
	Profile     *UserProfile     `json:"userProfile"`
	UserMessage string           `json:"userMessage"`
}
