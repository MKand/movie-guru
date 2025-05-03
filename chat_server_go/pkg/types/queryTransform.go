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
	TransformedQuery     string     `json:"transformedQuery, omitempty"`
	Intent               USERINTENT `json:"userIntent, omitempty"`
	*ModelOutputMetadata `json:"modelOutputMetadata"`
}

type QueryTransformFlowInput struct {
	History     []*SimpleMessage `json:"history"`
	Profile     *UserProfile     `json:"userProfile"`
	UserMessage string           `json:"userMessage"`
}
