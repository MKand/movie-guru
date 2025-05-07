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
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	redis "github.com/redis/go-redis/v9"
)

var (
	redisStore redis.Cmdable
)

func TestRedis() {
	ctx := context.Background()
	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	REDIS_PORT := os.Getenv("REDIS_PORT")

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
		Password: REDIS_PASSWORD,
		DB:       0,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		slog.ErrorContext(ctx, "error connecting to redis", slog.Any("error", err))
		return
	}
}

func setupSessionStore(ctx context.Context) {
	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_MODE := os.Getenv("REDIS_MODE")
	if REDIS_MODE == "" {
		REDIS_MODE = "SINGLE"
	}

	//REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	REDIS_PORT := os.Getenv("REDIS_PORT")
	if REDIS_MODE == "SINGLE" {
		redisStore = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
			DB:   0,
		})
	} else if REDIS_MODE == "CLUSTER" {
		redisStore = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: []string{fmt.Sprintf("%s:%s", REDIS_HOST, "6379")},
		})
	}

	if err := redisStore.Ping(ctx).Err(); err != nil {
		slog.ErrorContext(ctx, "error connecting to redis", slog.Any("error", err))
		os.Exit(1)
	}
}

func getSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("movie-guru-sid")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return "", errors.New("No cookie found")
		default:
			log.Println(err)
			return "", err
		}
	}
	sessionID := cookie.Value
	if sessionID == "" {
		return "", errors.New("None or malformed cookie found")
	}
	return sessionID, nil
}

func authenticateAndGetSessionInfo(ctx context.Context, sessionInfo *SessionInfo, err error, r *http.Request, w http.ResponseWriter) (*SessionInfo, bool) {
	sessionInfo, err = getSessionInfo(ctx, r)
	if err != nil {
		if err, ok := err.(*AuthorizationError); ok {
			slog.InfoContext(ctx, "Unauthorized", slog.Any("error", err.Error()))
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return nil, true
		}
		slog.ErrorContext(ctx, "Error while getting session info", slog.Any("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, true
	}
	if !sessionInfo.Authenticated {
		slog.InfoContext(ctx, "Forbidden")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return nil, true
	}
	return sessionInfo, false
}

func getSessionInfo(ctx context.Context, r *http.Request) (*SessionInfo, error) {
	sessionID, err := getSessionID(r)
	if err != nil {
		return nil, &AuthorizationError{err.Error()}
	}
	session := &SessionInfo{}
	s, err := redisStore.Get(ctx, sessionID).Result()
	if err != nil {
		return nil, &AuthorizationError{fmt.Sprintf("Unknown session: %s", sessionID)}
	}
	err = json.Unmarshal([]byte(s), session)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to retrieve session info. %s", err))
	}
	return session, nil
}

func deleteSessionInfo(ctx context.Context, sessionID string) error {
	_, err := redisStore.Del(ctx, sessionID).Result()
	if err != nil {
		return err
	}
	return nil
}
