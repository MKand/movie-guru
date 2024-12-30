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
	"time"

	"github.com/movie-guru/pkg/db"
	metrics "github.com/movie-guru/pkg/metrics"
	types "github.com/movie-guru/pkg/types"
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

	http.HandleFunc("/", createHealthCheckHandler(deps, hcMeters))
	http.HandleFunc("/chat", createChatHandler(deps, chatMeters))
	http.HandleFunc("/history", createHistoryHandler())
	http.HandleFunc("/preferences", createPreferencesHandler(deps.DB))
	http.HandleFunc("/startup", createStartupHandler(deps))
	http.HandleFunc("/login", createLoginHandler(ulh, loginMeters))
	http.HandleFunc("/logout", logoutHandler)
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
		"great films", "cool films", "best films", "new films", "high rated films", "classic films",
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
	sessionID := strings.Split(r.Header.Get("Cookie"), "session=")[1]
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

func createStartupHandler(deps *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := r.Context()
		origin := r.Header.Get("Origin")
		addResponseHeaders(w, origin)
		sessionInfo := &SessionInfo{}
		if r.Method != "OPTIONS" {
			var shouldReturn bool
			sessionInfo, shouldReturn = authenticateAndGetSessionInfo(ctx, sessionInfo, err, r, w)
			if shouldReturn {
				return
			}
		}
		if r.Method == "GET" {
			addResponseHeaders(w, origin)
			user := sessionInfo.User
			pref, err := deps.DB.GetCurrentProfile(ctx, user)
			if err != nil {
				slog.ErrorContext(ctx, "Cannot get preferences", slog.String("user", user), slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			context, err := deps.MovieRetrieverFlowClient.RetriveDocuments(ctx, randomisedFeaturedFilmsQuery())
			if err != nil {
				slog.ErrorContext(ctx, "Error getting movie recommendations", slog.Any("error", err.Error()))
				http.Error(w, "Error get movie recommendations", http.StatusInternalServerError)
				return
			}
			agentResp := types.NewAgentResponse()
			agentResp.Context = context[0:5]
			agentResp.Preferences = pref
			agentResp.Result = types.SUCCESS

			json.NewEncoder(w).Encode(agentResp)
			return

		}
	}
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

func createPreferencesHandler(MovieDB *db.MovieDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := r.Context()
		origin := r.Header.Get("Origin")
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
			pref, err := MovieDB.GetCurrentProfile(ctx, user)
			if err != nil {
				slog.ErrorContext(ctx, "Cannot get preferences", slog.String("user", user), slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(pref)
			return
		}
		if r.Method == "POST" {
			pref := &PrefBody{
				Content: types.NewUserProfile(),
			}
			err := json.NewDecoder(r.Body).Decode(pref)
			if err != nil {
				slog.InfoContext(ctx, "Error while decoding request", slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = MovieDB.UpdateProfile(ctx, pref.Content, sessionInfo.User)
			if err != nil {
				slog.ErrorContext(ctx, "Error while fetching preferences", slog.String("user", user), slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			addResponseHeaders(w, origin)
			json.NewEncoder(w).Encode(map[string]string{"update": "success"})
			return
		}
		if r.Method == "OPTIONS" {
			addResponseHeaders(w, origin)
			handleOptions(w, origin)
			return
		}
	}
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

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	origin := r.Header.Get("Origin")
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
		err := deleteSessionInfo(ctx, sessionInfo.ID)
		if err != nil {
			slog.ErrorContext(ctx, "Error while deleting session info", slog.String("user", user), slog.Any("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"logout": "success"})

		return
	}
	if r.Method == "OPTIONS" {
		addResponseHeaders(w, origin)
		handleOptions(w, origin)
		return
	}

}
