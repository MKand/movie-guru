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

type MovieFlowInput struct {
	History          []*SimpleMessage `json:"history"`
	UserPreferences  *UserProfile     `json:"userPreferences"`
	ContextDocuments []*MovieContext  `json:"contextDocuments"`
	UserMessage      string           `json:"userMessage"`
}

type MovieFlowOutput struct {
	Answer               string           `json:"answer"`
	RelevantMoviesTitles []*RelevantMovie `json:"relevantMovies"`
	WrongQuery           bool             `json:"wrongQuery,omitempty" `
	*ModelOutputMetadata `json:"modelOutputMetadata"`
}

type ExtendedMovieFlowOutput struct {
	Answer               string           `json:"answer"`
	RelevantMoviesTitles []*RelevantMovie `json:"relevantMovies"`
	WrongQuery           bool             `json:"wrongQuery,omitempty" `
	ContextDocuments     []*MovieContext  `json:"contextDocuments"`
	TraceId              string
	SpanId               string
	*ModelOutputMetadata `json:"modelOutputMetadata"`
}

type RelevantMovie struct {
	Title  string `json:"title"`
	Reason string `json:"reason"`
}
