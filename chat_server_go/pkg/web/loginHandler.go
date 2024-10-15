package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/movie-guru/pkg/db"
	m "github.com/movie-guru/pkg/metrics"
	"golang.org/x/exp/slog"
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

func (ulh *UserLoginHandler) HandleLogin(ctx context.Context, user string) (string, error) {
	// Minimal login logic for simplicity. Accepts any email and just returns it.
	return user, nil
}

func createLoginHandler(ulh *UserLoginHandler, meters *m.LoginMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		origin := r.Header.Get("Origin")
		if r.Method == "POST" {
			startTime := time.Now()
			defer func() {
				meters.LoginLatencyHistogram.Record(ctx, int64(time.Since(startTime).Milliseconds()))
			}()

			meters.LoginCounter.Add(ctx, 1)

			user := r.Header.Get("user")
			if user == "" {
				slog.InfoContext(ctx, "No user header")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := ulh.HandleLogin(ctx, user)
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
			sessionID := fmt.Sprintf("session_%s", user)
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
			setCookieHeader := ""
			if os.Getenv("LOCAL") == "true" {
				setCookieHeader = fmt.Sprintf("movieguru=%s; HttpOnly; SameSite=Lax; Path=/; Domain=localhost; Max-Age=86400", sessionID)
			} else {
				setCookieHeader = fmt.Sprintf("movieguru=%s; HttpOnly; Secure; SameSite=None; Path=/; Domain=%s; Max-Age=86400", sessionID, metadata.FrontEndDomain)
			}
			w.Header().Set("Set-Cookie", setCookieHeader)
			w.Header().Set("Vary", "Cookie, Origin")
			addResponseHeaders(w, origin)
			json.NewEncoder(w).Encode(map[string]string{"login": "success"})
		}
		if r.Method == "OPTIONS" {
			handleOptions(w, origin)
			return
		}
	}
}
