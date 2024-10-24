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
	"go.opentelemetry.io/otel"
)

var (
	redisStore *redis.Client
	metadata   *db.Metadata
	appConfig  = map[string]string{
		"CORS_HEADERS": "Content-Type",
	}
)

func StartServer(ctx context.Context, ulh *UserLoginHandler, m *db.Metadata, deps *Dependencies) error {
	metadata = m
	setupSessionStore(ctx)

	podName := os.Getenv("POD_NAME")
	if podName == "" {
		podName = "local"
	}
	meter := otel.Meter(podName)

	loginMeters := metrics.NewLoginMeters(meter)

	hcMeters := metrics.NewHCMeters(meter)
	chatMeters := metrics.NewChatMeters(meter)
	prefMeters := metrics.NewPreferencesMeters(meter)
	startupMeters := metrics.NewStartupMeters(meter)
	logoutMeters := metrics.NewLogoutMeters(meter)

	http.HandleFunc("/", createHealthCheckHandler(deps, hcMeters))
	http.HandleFunc("/chat", createChatHandler(deps, chatMeters))
	http.HandleFunc("/history", createHistoryHandler())
	http.HandleFunc("/preferences", createPreferencesHandler(deps.DB, prefMeters))
	http.HandleFunc("/startup", createStartupHandler(deps, startupMeters))
	http.HandleFunc("/login", createLoginHandler(ulh, loginMeters, m))
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
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "user, Origin, Cookie, Accept, Content-Type, Content-Length, Accept-Encoding,Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func handleOptions(w http.ResponseWriter, origin string) {
	addResponseHeaders(w, origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,OPTIONS,PUT")
	w.WriteHeader(http.StatusOK)
}

func getSessionID(r *http.Request) (string, error) {
	sessionID := ""
	if os.Getenv("SIMPLE") == "true" {
		user := r.Header.Get("user")
		sessionID = createSessionID(user)
		return sessionID, nil
	}
	if r.Header.Get("Cookie") == "" {
		return "", errors.New("No cookie found")
	}
	slog.InfoContext(r.Context(), "Cookies", r.Header.Get("Cookie"))
	sessionID = strings.Split(r.Header.Get("Cookie"), "movieguru=")[1]
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
