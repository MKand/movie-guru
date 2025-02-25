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

package web

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/movie-guru/pkg/db"
	m "github.com/movie-guru/pkg/metrics"
	utils "github.com/movie-guru/pkg/utils"
)

type AuthorizationError struct {
	Message string
}

func (e *AuthorizationError) Error() string {
	return e.Message
}

type UserLoginHandler struct {
	db            *db.MovieDB
	tokenAudience string
}

func NewUserLoginHandler(tokenAudience string, db *db.MovieDB) *UserLoginHandler {
	return &UserLoginHandler{
		db:            db,
		tokenAudience: tokenAudience,
	}
}

func (ulh *UserLoginHandler) HandleAPILogin(ctx context.Context, authHeader, inviteCode string) (string, error) {
	token := ulh.getToken(authHeader)
	user, err := ulh.verifyGoogleToken(token)
	if err != nil {
		return "", err
	}

	if ulh.db.CheckUser(ctx, user) {
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

func (ulh *UserLoginHandler) HandleLogin(ctx context.Context, authHeader, inviteCode string) (string, error) {
	token := ulh.getToken(authHeader)
	user, err := ulh.verifyGoogleToken(token)
	if err != nil {
		return "", err
	}

	if ulh.db.CheckUser(ctx, user) {
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

func (ulh *UserLoginHandler) HandleApiKeyLogin(ctx context.Context, apiKey, user string) (string, error) {
	// simple implementation for now
	if user == "" {
		return "", &AuthorizationError{"Invalid invite code"}
	}
	return user, nil
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

func createLoginHandler(ulh *UserLoginHandler, meters *m.LoginMeters, metadata *db.Metadata) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method == "POST" {
			startTime := time.Now()
			defer func() {
				meters.LoginLatencyHistogram.Record(ctx, int64(time.Since(startTime).Milliseconds()))
			}()

			meters.LoginCounter.Add(ctx, 1)

			authHeader := r.Header.Get("Authorization")
			apiKey := r.Header.Get("ApiKey")
			user := r.Header.Get("User")

			if authHeader == "" && apiKey == "" {
				slog.InfoContext(ctx, "No auth header")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			var loginBody LoginBody
			err := json.NewDecoder(r.Body).Decode(&loginBody)
			if err != nil {
				slog.ErrorContext(ctx, "Bad Request at login", slog.Any("error", err.Error()))

				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if authHeader != "" {
				user, err = ulh.HandleLogin(ctx, authHeader, loginBody.InviteCode)
			}
			if apiKey != "" {
				user, err = ulh.HandleApiKeyLogin(ctx, apiKey, user)
			}

			if err != nil {
				if _, ok := err.(*AuthorizationError); ok {
					slog.InfoContext(ctx, "Unauthorized")
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				slog.ErrorContext(ctx, "Error while getting user from db", slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			sessionID := uuid.New().String()
			session := &SessionInfo{
				ID:            sessionID,
				User:          user,
				Authenticated: true,
			}
			sessionJSON, err := json.Marshal(session)
			if err != nil {
				slog.ErrorContext(ctx, "Error while decoding session info", slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			err = redisStore.Set(r.Context(), sessionID, sessionJSON, 0).Err()
			if err != nil {
				slog.ErrorContext(ctx, "Error while setting context in redis", slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			meters.LoginSuccessCounter.Add(ctx, 1)
			cookie := http.Cookie{
				Name:     "movie-guru-sid",
				Value:    sessionID,
				Path:     "/",
				MaxAge:   86400,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, &cookie)
			w.Header().Set("Vary", "Cookie, Origin")
			json.NewEncoder(w).Encode(map[string]string{"login": "success"})
		}
	}
}

// get_token extracts the token from the authorization header
func (ulh *UserLoginHandler) getToken(authHeader string) string {
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) == 2 && strings.ToLower(tokenParts[0]) == "bearer" {
		return tokenParts[1]
	}
	return ""
}
