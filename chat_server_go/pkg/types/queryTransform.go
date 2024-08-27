package types

type USERINTENT string

const (
	UNCLEAR          RESULT = "UNCLEAR"
	GREET            RESULT = "GREET"
	END_CONVERSATION RESULT = "END_CONVERSATION"
	REQUEST          RESULT = "REQUEST"
	RESPONSE         RESULT = "RESPONSE"
	ACKNOWLEDGE      RESULT = "ACKNOWLEDGE"
)

type QueryTransformOutput struct {
	TransformedQuery string     `json:"transformedQuery, omitempty"`
	Intent           USERINTENT `json:"userIntent, omitempty"`
	*ModelOutputMetadata
}

type QueryTransformInput struct {
	History     []*SimpleMessage `json:"history"`
	Profile     *UserProfile     `json:"userProfile"`
	UserMessage string           `json:"userMessage"`
}
