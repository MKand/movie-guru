package wrappers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/movie-guru/pkg/db"
	types "github.com/movie-guru/pkg/types"
)

type ProfileAgent struct {
	MovieAgentDB *db.MovieAgentDB
	URL          string
}

func CreateProfileAgent(db *db.MovieAgentDB, URL string) (*ProfileAgent, error) {
	return &ProfileAgent{
		MovieAgentDB: db,
		URL:          URL + "/userPreferencesFlow",
	}, nil
}

func (p *ProfileAgent) Run(ctx context.Context, history *types.ChatHistory, user string) (*types.UserProfileOutput, error) {
	userProfile, err := p.MovieAgentDB.GetCurrentProfile(ctx, user)
	userProfileOutput := &types.UserProfileOutput{
		UserProfile: userProfile,
		ChangesMade: false,
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

	prefInput := types.ProfileAgentInput{Query: lastUserMessage, AgentMessage: agentMessage}
	resp, err := p.runFlow(&prefInput)
	if err != nil {
		return userProfileOutput, err
	}
	userProfileOutput.ChangesMade = resp.ChangesMade
	userProfileOutput.ModelOutputMetadata.Justification = resp.ModelOutputMetadata.Justification
	userProfileOutput.ModelOutputMetadata.SafetyIssue = resp.ModelOutputMetadata.SafetyIssue

	if resp.ChangesMade {
		updatedProfile, err := processProfileChanges(userProfile, resp.ProfileChangeRecommendations)
		if err != nil {
			return userProfileOutput, err
		}
		err = p.MovieAgentDB.UpdateProfile(ctx, updatedProfile, user)
		if err != nil {
			return userProfileOutput, err
		}
		userProfileOutput.UserProfile = updatedProfile
	}
	return userProfileOutput, nil
}

func (agent *ProfileAgent) runFlow(input *types.ProfileAgentInput) (*types.UserProfileAgentOutput, error) {
	// Marshal the input struct to JSON
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input to JSON: %w", err)
	}
	req, err := http.NewRequest("POST", agent.URL, bytes.NewBuffer(inputJSON))
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
	defer resp.Body.Close()

	var output *types.UserProfileAgentOutput
	err = json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return nil, err
	}

	return output, nil
}
