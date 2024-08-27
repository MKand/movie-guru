package main

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type PreferencesAgent struct {
	Model ai.Model
	Flow  *genkit.Flow[*ProfileAgentInput, *UserProfileAgentOutput, struct{}]
	DB    *sql.DB
}

func CreatePreferencesAgent(ctx context.Context, model ai.Model, db *sql.DB) (*PreferencesAgent, error) {
	flow, err := GetUserPrefFlow(ctx, model)
	if err != nil {
		return nil, err
	}
	return &PreferencesAgent{
		Model: model,
		Flow:  flow,
		DB:    db,
	}, nil
}

func (p *PreferencesAgent) Run(ctx context.Context, history *ChatHistory, user string) (*UserProfileOutput, error) {
	userProfile, err := getCurrentProfile(ctx, user, p.DB)
	userProfileOutput := &UserProfileOutput{
		UserProfile: userProfile,
		ChangesMade: false,
		ModelOutputMetadata: &ModelOutputMetadata{
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
	if err != nil {
		return nil, err
	}
	lastUserMessage, err := history.GetLastMessage()
	if err != nil {
		return nil, err
	}

	prefInput := ProfileAgentInput{Query: lastUserMessage, AgentMessage: agentMessage}
	resp, err := p.Flow.Run(ctx, &prefInput)
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
		err = updatePreferences(ctx, updatedProfile, user, p.DB)
		if err != nil {
			return userProfileOutput, err
		}
		userProfileOutput.UserProfile = updatedProfile
	}
	return userProfileOutput, nil
}

func getCurrentProfile(ctx context.Context, user string, db *sql.DB) (*UserProfile, error) {
	preferences := NewUserProfile()
	rows := db.QueryRowContext(ctx, `
	SELECT preferences FROM user_preferences 
	WHERE "user" = $1;`,
		user)
	var jsonData string
	err := rows.Scan(&jsonData)
	if err != nil {
		return preferences, err
	}
	err = json.Unmarshal([]byte(jsonData), &preferences)
	if err != nil {
		return preferences, err
	}
	return preferences, nil
}

func updatePreferences(ctx context.Context, newPref *UserProfile, user string, db *sql.DB) error {
	newPreferencesStr, err := json.Marshal(newPref)
	if err != nil {
		return err
	}
	query := `
        INSERT INTO user_preferences ("user", preferences)
        VALUES ($1, $2)
        ON CONFLICT ("user") DO UPDATE
        SET preferences = EXCLUDED.preferences;
    `

	// Execute the query (replace with your actual execute_query function)
	_, err = db.ExecContext(ctx, query, user, newPreferencesStr)
	if err != nil {
		return err
	}
	return nil
}

func deletePreferences(ctx context.Context, user string, db *sql.DB) error {
	query := `
		DELETE FROM user_preferences
		WHERE "user" = %1;
	`
	_, err := db.ExecContext(ctx, query, user)
	if err != nil {
		return err
	}
	return nil
}
