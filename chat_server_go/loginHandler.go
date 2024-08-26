package main

import (
	"context"
	"database/sql"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

// AuthorizationError represents an authorization failure
type AuthorizationError struct {
	Message string
}

func (e *AuthorizationError) Error() string {
	return e.Message
}

type UserLoginHandler struct {
	db            *sql.DB
	tokenAudience string
}

func (ulh *UserLoginHandler) handleLogin(authHeader, inviteCode string) (string, error) {
	token := ulh.getToken(authHeader)
	user, err := ulh.verifyGoogleToken(token)
	if err != nil {
		return "", err
	}

	if ulh.checkUser(user) {
		return user, nil
	}

	inviteCodes, err := ulh.getInviteCodes()
	if err != nil {
		return "", err
	}

	if contains(inviteCodes, inviteCode) {
		if err := ulh.createUser(user); err != nil {
			return "", err
		}
		return user, nil
	}

	return "", &AuthorizationError{"Invalid invite code"}
}

// verify_google_token verifies the Google token and extracts the user email
func (ulh *UserLoginHandler) verifyGoogleToken(tokenString string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", &AuthorizationError{"Invalid token"}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", &AuthorizationError{"Invalid token claims"}
	}

	aud, ok := claims["aud"].(string)
	if !ok || aud != ulh.tokenAudience {
		return "", &AuthorizationError{"Invalid token audience"}
	}

	emailVerified, ok := claims["email_verified"].(bool)
	if !ok || !emailVerified {
		return "", &AuthorizationError{"Email not verified"}
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", &AuthorizationError{"Email not found in token"}
	}

	return email, nil
}

// get_token extracts the token from the authorization header
func (ulh *UserLoginHandler) getToken(authHeader string) string {
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) == 2 && strings.ToLower(tokenParts[0]) == "bearer" {
		return tokenParts[1]
	}
	return ""
}

// create_user creates a new user in the database
func (ulh *UserLoginHandler) createUser(user string) error {
	query := `
        INSERT INTO user_logins (email) VALUES ($1)
        ON CONFLICT (email) DO UPDATE
        SET login_count = user_logins.login_count + 1;
    `
	_, err := ulh.db.ExecContext(context.Background(), query, user)
	return err
}

// check_user checks if the user exists in the database
func (ulh *UserLoginHandler) checkUser(user string) bool {
	query := `SELECT email FROM user_logins WHERE "email" = $1;`
	var email string
	err := ulh.db.QueryRowContext(context.Background(), query, user).Scan(&email)
	return err == nil && email == user
}

// get_invite_codes retrieves valid invite codes from the database
func (ulh *UserLoginHandler) getInviteCodes() ([]string, error) {
	query := `SELECT code FROM invite_codes WHERE valid = true`
	rows, err := ulh.db.QueryContext(context.Background(), query)
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
