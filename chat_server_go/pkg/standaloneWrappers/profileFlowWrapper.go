package standaloneWrappers

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	db "github.com/movie-guru/pkg/db"
	flows "github.com/movie-guru/pkg/flows"
	types "github.com/movie-guru/pkg/types"
	utils "github.com/movie-guru/pkg/utils"
)

type ProfileFlow struct {
	MovieDB *db.MovieDB
	Flow    *genkit.Flow[*types.UserProfileFlowInput, *types.UserProfileFlowOutput, struct{}]
}

func CreateProfileFlow(ctx context.Context, model ai.Model, db *db.MovieDB) (*ProfileFlow, error) {
	flow, err := flows.GetUserProfileFlow(ctx, model)
	if err != nil {
		return nil, err
	}
	return &ProfileFlow{
		MovieDB: db,
		Flow:    flow,
	}, nil
}

func (p *ProfileFlow) Run(ctx context.Context, history *types.ChatHistory, user string) (*types.UserProfileOutput, error) {
	userProfile, err := p.MovieDB.GetCurrentProfile(ctx, user)
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

	prefInput := types.UserProfileFlowInput{Query: lastUserMessage, AgentMessage: agentMessage}
	resp, err := p.Flow.Run(ctx, &prefInput)
	if err != nil {
		return userProfileOutput, err
	}
	userProfileOutput.ChangesMade = resp.ChangesMade
	userProfileOutput.ModelOutputMetadata.Justification = resp.ModelOutputMetadata.Justification
	userProfileOutput.ModelOutputMetadata.SafetyIssue = resp.ModelOutputMetadata.SafetyIssue

	if len(resp.ProfileChangeRecommendations) > 0 {
		updatedProfile, err := utils.ProcessProfileChanges(userProfile, resp.ProfileChangeRecommendations)
		if err != nil {
			return userProfileOutput, err
		}
		err = p.MovieDB.UpdateProfile(ctx, updatedProfile, user)
		if err != nil {
			return userProfileOutput, err
		}
		userProfileOutput.UserProfile = updatedProfile
	}
	return userProfileOutput, nil
}
