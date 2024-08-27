package types

import (
	"fmt"
	"regexp"
)

type RESULT string

const (
	UNDEFINED RESULT = "UNDEFINED"
	SUCCESS   RESULT = "SUCCESS"
	BAD_QUERY RESULT = "BAD_QUERY"
	UNSAFE    RESULT = "UNSAFE"
	TOO_LONG  RESULT = "TOO_LONG"
	ERROR     RESULT = "ERROR"
)

type ModelOutputMetadata struct {
	Justification string `json:"justification" omitempty`
	SafetyIssue   bool   `json:"safetyIssue" omitempty`
}

type AgentResponse struct {
	Answer         string          `json:"answer"`
	RelevantMovies []string        `json:"relevant_movies"`
	Context        []*MovieContext `json:"context"`
	ErrorMessage   string          `json:"error_message"`
	Result         RESULT          `json:"result"`
	Preferences    *UserProfile    `json:"preferences"`
}

func NewAgentResponse() *AgentResponse {
	return &AgentResponse{
		RelevantMovies: make([]string, 0),
		Context:        make([]*MovieContext, 0),
		Preferences:    NewUserProfile(),
		Result:         UNDEFINED,
	}
}

func NewSafetyIssueAgentResponse() *AgentResponse {
	r := NewAgentResponse()
	r.Result = UNSAFE
	return r
}

func NewErrorAgentResponse(errMessage string) *AgentResponse {
	r := NewAgentResponse()
	r.Result = ERROR
	r.ErrorMessage = errMessage
	return r
}

func makeJsonMarshallable(input string) (string, error) {
	// Regex to extract JSON content from Markdown code block
	re := regexp.MustCompile("```(json)?((\n|.)*?)```")
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		return input, fmt.Errorf("no JSON content found in the input")
	}

	jsonContent := matches[2]
	return jsonContent, nil
}
