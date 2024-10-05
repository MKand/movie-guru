package web

import (
	"context"

	_ "github.com/lib/pq"
	"github.com/movie-guru/pkg/db"
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
