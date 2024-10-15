package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/movie-guru/pkg/db"
	metrics "github.com/movie-guru/pkg/metrics"
	"github.com/redis/go-redis/v9"
)

var (
	redisStore *redis.Client
	metadata   *db.Metadata
	appConfig  = map[string]string{
		"CORS_HEADERS": "Content-Type",
	}
	corsOrigins []string
)

func StartServer(ctx context.Context, ulh *UserLoginHandler, m *db.Metadata, deps *Dependencies) error {
	metadata = m
	setupSessionStore(ctx)

	corsOrigins = strings.Split(metadata.CorsOrigin, ",")
	for i := range corsOrigins {
		corsOrigins[i] = strings.TrimSpace(corsOrigins[i])
	}

	loginMeters := metrics.NewLoginMeters()
	hcMeters := metrics.NewHCMeters()
	chatMeters := metrics.NewChatMeters()
	prefMeters := metrics.NewPreferencesMeters()
	startupMeters := metrics.NewStartupMeters()
	logoutMeters := metrics.NewLogoutMeters()

	http.HandleFunc("/", createHealthCheckHandler(deps, hcMeters))
	http.HandleFunc("/chat", createChatHandler(deps, chatMeters))
	http.HandleFunc("/history", createHistoryHandler())
	http.HandleFunc("/preferences", createPreferencesHandler(deps.DB, prefMeters))
	http.HandleFunc("/startup", createStartupHandler(deps, startupMeters))
	http.HandleFunc("/login", createLoginHandler(ulh, loginMeters))
	http.HandleFunc("/logout", createLogoutHandler(logoutMeters))
	return http.ListenAndServe(":8080", nil)
}

func setupSessionStore(ctx context.Context) {
	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	REDIS_PORT := os.Getenv("REDIS_PORT")

	redisStore = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
		Password: REDIS_PASSWORD,
		DB:       0,
	})
	if err := redisStore.Ping(ctx).Err(); err != nil {
		slog.ErrorContext(ctx, "error connecting to redis", slog.Any("error", err))
	}
}

func randomisedFeaturedFilmsQuery() string {
	queries := []string{
		"top films", "cool films", "best films", "new films", "top rated films", "classic films",
	}
	return queries[rand.Intn(len(queries))]

}

func addResponseHeaders(w http.ResponseWriter, origin string) {
	isAllowed := true
	for _, allowedOrigin := range corsOrigins {
		if origin == allowedOrigin {
			isAllowed = true
			break
		}
	}
	if isAllowed {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "user, Origin, Cookie, Accept, Content-Type, Content-Length, Accept-Encoding,Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func getSessionID(r *http.Request) (string, error) {
	if r.Header.Get("Cookie") == "" {
		return "", errors.New("No cookie found")
	}
	sessionID := strings.Split(r.Header.Get("Cookie"), "movieguru=")[1]
	return sessionID, nil
}

func handleOptions(w http.ResponseWriter, origin string) {
	isAllowed := false
	for _, allowedOrigin := range corsOrigins {
		if origin == allowedOrigin {
			isAllowed = true
			break
		}
	}
	if isAllowed {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,OPTIONS,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "user, Origin, Cookie, Accept, Content-Type, Content-Length, Accept-Encoding,Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.WriteHeader(http.StatusOK)
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
	session := &SessionInfo{}
	sessionID, err := getSessionID(r)
	if err != nil {
		return session, &AuthorizationError{err.Error()}
	}
	s, err := redisStore.Get(ctx, sessionID).Result()
	err = json.Unmarshal([]byte(s), session)
	if err != nil {
		return nil, err
	}
	return session, err
}

func deleteSessionInfo(ctx context.Context, sessionID string) error {
	_, err := redisStore.Del(ctx, sessionID).Result()
	if err != nil {
		return err
	}
	return nil
}
