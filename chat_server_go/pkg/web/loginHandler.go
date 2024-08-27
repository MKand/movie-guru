package web

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"github.com/movie-guru/pkg/db"
	utils "github.com/movie-guru/pkg/utils"
)

type AuthorizationError struct {
	Message string
}

func (e *AuthorizationError) Error() string {
	return e.Message
}

type UserLoginHandler struct {
	db            *db.MovieAgentDB
	tokenAudience string
}

func NewUserLoginHandler(tokenAudience string, db *db.MovieAgentDB) *UserLoginHandler {
	return &UserLoginHandler{
		db:            db,
		tokenAudience: tokenAudience,
	}
}

func (ulh *UserLoginHandler) handleLogin(authHeader, inviteCode string) (string, error) {
	token := ulh.getToken(authHeader)
	user, err := ulh.verifyGoogleToken(token)
	if err != nil {
		return "", err
	}

	if ulh.db.CheckUser(user) {
		return user, nil
	}

	inviteCodes, err := ulh.db.GetInviteCodes()
	if err != nil {
		return "", err
	}

	if utils.Contains(inviteCodes, inviteCode) {
		if err := ulh.db.CreateUser(user); err != nil {
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
