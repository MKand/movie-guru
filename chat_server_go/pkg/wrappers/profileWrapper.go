package wrappers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/movie-guru/pkg/db"
	types "github.com/movie-guru/pkg/types"
	utils "github.com/movie-guru/pkg/utils"
)

type UserProfileFlowClient struct {
	MovieDB *db.MovieDB
	URL     string
}

func CreateUserProfileFlowClient(db *db.MovieDB, URL string) (*UserProfileFlowClient, error) {
	return &UserProfileFlowClient{
		MovieDB: db,
		URL:     URL + "/userProfileFlow",
	}, nil
}

func (flowClient *UserProfileFlowClient) Run(ctx context.Context, history *types.ChatHistory, user string) (*types.UserProfileOutput, error) {
	userProfile, err := flowClient.MovieDB.GetCurrentProfile(ctx, user)
	userProfileOutput := &types.UserProfileOutput{
		UserProfile: userProfile,
		ModelOutputMetadata: &types.ModelOutputMetadata{
			SafetyIssue:   false,
			Justification: "",
		},
	}
	if err != nil {
		return nil, err
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
		return userProfileOutput, err
	}
	userProfileOutput.ModelOutputMetadata.Justification = resp.ModelOutputMetadata.Justification
	userProfileOutput.ModelOutputMetadata.SafetyIssue = resp.ModelOutputMetadata.SafetyIssue

	if len(resp.ProfileChangeRecommendations) > 0 {
		updatedProfile, err := utils.ProcessProfileChanges(userProfile, resp.ProfileChangeRecommendations)
		if err != nil {
			return userProfileOutput, err
		}
		err = flowClient.MovieDB.UpdateProfile(ctx, updatedProfile, user)
		if err != nil {
			return userProfileOutput, err
		}
		userProfileOutput.UserProfile = updatedProfile
	}
	return userProfileOutput, nil
}

func (flowClient *UserProfileFlowClient) runFlow(input *types.UserProfileFlowInput) (*types.UserProfileFlowOutput, error) {
	// Marshal the input struct to JSON
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input to JSON: %w", err)
	}
	req, err := http.NewRequest("POST", flowClient.URL, bytes.NewBuffer(inputJSON))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}

	var result struct {
		Result *types.UserProfileFlowOutput `json:"result"`
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return nil, err
	}
	return result.Result, nil
}
