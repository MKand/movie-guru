package web

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/movie-guru/pkg/db"
	types "github.com/movie-guru/pkg/types"
)

func createHistoryHandler(metadata *db.Metadata) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var err error
		sessionInfo := &SessionInfo{}
		if r.Method != "OPTIONS" {
			var shouldReturn bool
			sessionInfo, shouldReturn = authenticateAndGetSessionInfo(ctx, sessionInfo, err, r, w)
			if shouldReturn {
				return
			}
		}
		user := sessionInfo.User
		if r.Method == "GET" {
			ch, err := getHistory(ctx, user)
			if err != nil {
				slog.ErrorContext(ctx, "Error while fetching history", slog.String("user", user), slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			simpleHistory, err := types.ParseRecentHistory(ch.GetHistory(), metadata.HistoryLength)
			if err != nil {
				slog.ErrorContext(ctx, "Error while parsing history", slog.String("user", user), slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(simpleHistory)
		}
		if r.Method == "DELETE" {
			err := deleteHistory(ctx, sessionInfo.User)
			if err != nil {
				slog.ErrorContext(ctx, "Error while deleting history", slog.String("user", user), slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}

func deleteHistory(ctx context.Context, user string) error {
	redisContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := redisStore.Del(redisContext, user).Result()
	if err != nil {
		return err
	}
	return nil
}
