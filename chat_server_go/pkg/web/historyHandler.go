package web

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/movie-guru/pkg/types"
	"github.com/redis/go-redis/v9"
)

func getHistory(ctx context.Context, user string) (*types.ChatHistory, error) {
	redisContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	historyJson, err := redisStore.Get(redisContext, user).Result()
	ch := types.NewChatHistory()
	if err == redis.Nil {
		return ch, nil
	} else if err != nil {
		return ch, err
	}
	err = json.Unmarshal([]byte(historyJson), ch)
	if err != nil {
		return ch, err
	}
	return ch, nil
}

func saveHistory(ctx context.Context, history *types.ChatHistory, user string) error {
	history.Trim(metadata.HistoryLength)
	redisContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	err := redisStore.Set(redisContext, user, history, 0).Err()
	if err != nil {
		return err
	}
	return nil
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

func createHistoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		origin := r.Header.Get("Origin")
		var err error
		addResponseHeaders(w, origin)
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
			addResponseHeaders(w, origin)
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
			addResponseHeaders(w, origin)
			err := deleteHistory(ctx, sessionInfo.User)
			if err != nil {
				slog.ErrorContext(ctx, "Error while deleting history", slog.String("user", user), slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method == "OPTIONS" {
			addResponseHeaders(w, origin)
			handleOptions(w, origin)
			return
		}
	}
}
