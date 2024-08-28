package db

import (
	"context"

	_ "github.com/lib/pq"
)

// create_user creates a new user in the database
func (db *MovieDB) CreateUser(user string) error {
	query := `
        INSERT INTO user_logins (email) VALUES ($1)
        ON CONFLICT (email) DO UPDATE
        SET login_count = user_logins.login_count + 1;
    `
	_, err := db.DB.ExecContext(context.Background(), query, user)
	return err
}

// check_user checks if the user exists in the database
func (db *MovieDB) CheckUser(user string) bool {
	query := `SELECT email FROM user_logins WHERE "email" = $1;`
	var email string
	err := db.DB.QueryRowContext(context.Background(), query, user).Scan(&email)
	return err == nil && email == user
}

// get_invite_codes retrieves valid invite codes from the database
func (db *MovieDB) GetInviteCodes() ([]string, error) {
	query := `SELECT code FROM invite_codes WHERE valid = true`
	rows, err := db.DB.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inviteCodes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		inviteCodes = append(inviteCodes, code)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return inviteCodes, nil
}
