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

type UserProfileFlowInput struct {
	Query        string `json:"query"`
	AgentMessage string `json:"agentMessage"`
}

type UserProfileFlowOutput struct {
	ProfileChangeRecommendations []*ProfileChangeRecommendation `json:"profileChangeRecommendations"`
	*ModelOutputMetadata         `json:"modelOutputMetadata"`
}

type ProfileChangeRecommendation struct {
	Item     string               `json:"item"`
	Reason   string               `json:"reason"`
	Category MovieFeatureCategory `json:"category"`
	Sentiment
}

type UserProfileOutput struct {
	UserProfile *UserProfile `json:"userProfile"`
	*ModelOutputMetadata
}

func NewUserProfileFlowOuput() *UserProfileFlowOutput {
	return &UserProfileFlowOutput{
		ProfileChangeRecommendations: make([]*ProfileChangeRecommendation, 5),
		ModelOutputMetadata: &ModelOutputMetadata{
			Justification: "",
			SafetyIssue:   false,
		},
	}
}

type MovieFeatureCategory string
type ProfileAction string
type Sentiment string

const (
	OTHER    MovieFeatureCategory = "OTHER"
	ACTOR    MovieFeatureCategory = "ACTOR"
	DIRECTOR MovieFeatureCategory = "DIRECTOR"
	GENRE    MovieFeatureCategory = "GENRE"
)

const (
	UNSURE ProfileAction = "UNSURE"
	ADD    ProfileAction = "ADD"
	REMOVE ProfileAction = "REMOVE"
)

const (
	POSITIVE Sentiment = "POSITIVE"
	NEGATIVE Sentiment = "NEGATIVE"
)

type UserProfile struct {
	Likes    ProfileCategories `json:"likes, omitempty"`
	Dislikes ProfileCategories `json:"dislikes, omitempty"`
}
type ProfileCategories struct {
	Actors    []string `json:"actors, omitempty"`
	Directors []string `json:"directors, omitempty"`
	Genres    []string `json:"genres, omitempty"`
	Others    []string `json:"others, omitempty"`
}

func NewUserProfile() *UserProfile {
	return &UserProfile{
		Likes: ProfileCategories{
			Actors:    make([]string, 0),
			Directors: make([]string, 0),
			Genres:    make([]string, 0),
			Others:    make([]string, 0),
		},
		Dislikes: ProfileCategories{
			Actors:    make([]string, 0),
			Directors: make([]string, 0),
			Genres:    make([]string, 0),
			Others:    make([]string, 0),
		},
	}
}
