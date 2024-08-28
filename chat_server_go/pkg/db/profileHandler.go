package db

import (
	"context"
	"encoding/json"

	types "github.com/movie-guru/pkg/types"
)

func (MovieDB *MovieDB) GetCurrentProfile(ctx context.Context, user string) (*types.UserProfile, error) {
	preferences := types.NewUserProfile()
	rows := MovieDB.DB.QueryRowContext(ctx, `
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

func (MovieDB *MovieDB) UpdateProfile(ctx context.Context, newPref *types.UserProfile, user string) error {
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
	_, err = MovieDB.DB.ExecContext(ctx, query, user, newPreferencesStr)
	if err != nil {
		return err
	}
	return nil
}

func (MovieDB *MovieDB) DeleteProfile(ctx context.Context, user string) error {
	query := `
		DELETE FROM user_preferences
		WHERE "user" = %1;
	`
	_, err := MovieDB.DB.ExecContext(ctx, query, user)
	if err != nil {
		return err
	}
	return nil
}
