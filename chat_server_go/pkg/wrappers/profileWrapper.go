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

package wrappers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	db "github.com/movie-guru/pkg/db"
	types "github.com/movie-guru/pkg/types"
	utils "github.com/movie-guru/pkg/utils"
)

type UserProfileFlowClient struct {
	MovieDB *db.MovieDB
	URL     string
}

type UserProfileFlowOutput struct {
	ProfileChangeRecommendations []*types.ProfileChangeRecommendation `json:"profileChangeRecommendations"`
	Justification                string                               `json:"justification"`
	BadQuery                     bool                                 `json:"badQuery,omitempty" `
	SafetyIssue                  bool                                 `json:"safetyIssue,omitempty"`
	QuotaIssue                   bool                                 `json:"quotaIssue,omitempty"`
}

func CreateUserProfileFlowClient(db *db.MovieDB, URL string) (*UserProfileFlowClient, error) {
	return &UserProfileFlowClient{
		MovieDB: db,
		URL:     URL + "/userPreferenceFlow",
	}, nil
}

func (flowClient *UserProfileFlowClient) Run(ctx context.Context, history *types.ChatHistory, user string, userProfile *types.UserProfile) (*types.UserProfileWrapperOutput, error) {
	userProfileOutput := &types.UserProfileWrapperOutput{
		UserProfile: userProfile,
		ModelOutputMetadata: &types.ModelOutputMetadata{
			SafetyIssue:   false,
			Justification: "",
			BadQuery:      false,
		},
	}

	agentMessage := ""
	if len(history.History) > 1 {
		agentMessage = history.History[len(history.History)-2].Content[0].Text
	}
	lastUserMessage, err := history.GetLastMessage()
	if err != nil {
		return nil, err
	}

	userProfileFlowInput := types.UserProfileFlowInput{Query: lastUserMessage, AgentMessage: agentMessage}
	resp, err := flowClient.runFlow(&userProfileFlowInput)
	if err != nil {
		return nil, err
	}
	userProfileOutput.ModelOutputMetadata.Justification = resp.Justification
	userProfileOutput.ModelOutputMetadata.SafetyIssue = resp.SafetyIssue
	userProfileOutput.ModelOutputMetadata.QuotaIssue = resp.QuotaIssue

	if len(resp.ProfileChangeRecommendations) > 0 {
		updatedProfile, err := utils.ProcessProfileChanges(userProfile, resp.ProfileChangeRecommendations)
		if err != nil {
			return userProfileOutput, err
		}
		err = flowClient.MovieDB.UpdateProfile(ctx, updatedProfile, user)
		if err != nil {
			slog.ErrorContext(ctx, "DB Update error", err.Error(), err)
			return userProfileOutput, err
		}
		userProfileOutput.UserProfile = updatedProfile
	}
	return userProfileOutput, nil
}

func (flowClient *UserProfileFlowClient) runFlow(input *types.UserProfileFlowInput) (*UserProfileFlowOutput, error) {
	// Marshal the input struct to JSON
	dataInput := DataInput{
		Data: input,
	}

	inputJSON, err := json.Marshal(dataInput)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input to JSON: %w", err)
	}

	req, err := http.NewRequest("POST", flowClient.URL, bytes.NewBuffer(inputJSON))
	if err != nil {
		slog.Log(context.Background(), slog.LevelError, "Error creating request", "error", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	ctx := context.Background()

	if err != nil {
		slog.ErrorContext(ctx, "QualityFlow: Error sending request to Flows", err.Error(), err)
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.ErrorContext(ctx, "QualityFlow: Genkit returned an Error", "errorCode", resp.StatusCode)
		return nil, fmt.Errorf("genkit server returned error: %s (%d)", http.StatusText(resp.StatusCode), resp.StatusCode)
	}
	var result struct {
		Result *UserProfileFlowOutput `json:"result"`
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(b, &result)
	if err != nil {
		slog.Log(context.Background(), slog.LevelError, "Error unmarshaling JSON response", "error", err)
		return nil, err
	}

	/*b = bytes.TrimSpace(b)
	resp.Body = ioutil.NopCloser(bytes.NewReader(b))

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		slog.Log(context.Background(), slog.LevelError, "Error decoding JSON response", "error", err)
		return nil, err
	}*/

	return result.Result, nil
}
