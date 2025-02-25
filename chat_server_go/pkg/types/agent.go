// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

type RESULT string

const (
	UNDEFINED  RESULT = "UNDEFINED"
	SUCCESS    RESULT = "SUCCESS"
	BAD_QUERY  RESULT = "BAD_QUERY"
	UNSAFE     RESULT = "UNSAFE"
	TOO_LONG   RESULT = "TOO_LONG"
	ERROR      RESULT = "ERROR"
	QUOTALIMIT RESULT = "QUOTALIMIT"
)

type ModelOutputMetadata struct {
	Justification string `json:"justification" omitempty`
	SafetyIssue   bool   `json:"safetyIssue" omitempty`
	QuotaIssue    bool   `json: "quotaIssue" omitempty`
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
	r.Answer = "That was very naughty of you! I cannot answer that."
	r.Result = UNSAFE
	return r
}

func NewQuotaIssueAgentResponse() *AgentResponse {
	r := NewAgentResponse()
	r.Answer = "There are too many requests to the underlying LLM. Can you wait a bit and try again?"
	r.Result = QUOTALIMIT
	r.ErrorMessage = "Vertex Quota exceeded"
	return r
}

func NewErrorAgentResponse(errMessage string) *AgentResponse {
	r := NewAgentResponse()
	r.Result = ERROR
	r.ErrorMessage = errMessage
	r.Answer = "Something went wrong on my side. My apologies. Can you try that again?"
	return r
}
