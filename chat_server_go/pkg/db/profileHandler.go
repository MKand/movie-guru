package db

import (
	"context"
	"encoding/json"
	"time"

	types "github.com/movie-guru/pkg/types"
)

func (MovieDB *MovieDB) GetCurrentProfile(ctx context.Context, user string) (*types.UserProfile, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	preferences := types.NewUserProfile()
	rows := MovieDB.DB.QueryRowContext(dbCtx, `
	SELECT preferences FROM user_preferences 
	WHERE "user" = $1;`,
		user)
	var jsonData string
	err := rows.Scan(&jsonData)
	if err != nil {
		return preferences, nil
	}
	err = json.Unmarshal([]byte(jsonData), &preferences)
	if err != nil {
		return preferences, err
	}
	return preferences, nil
}

func (MovieDB *MovieDB) UpdateProfile(ctx context.Context, newPref *types.UserProfile, user string) error {
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
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
	_, err = MovieDB.DB.ExecContext(dbCtx, query, user, newPreferencesStr)
	if err != nil {
		return err
	}
	return nil
}

func (MovieDB *MovieDB) DeleteProfile(ctx context.Context, user string) error {
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	query := `
		DELETE FROM user_preferences
		WHERE "user" = %1;
	`
	_, err := MovieDB.DB.ExecContext(dbCtx, query, user)
	if err != nil {
		return err
	}
	return nil
}
