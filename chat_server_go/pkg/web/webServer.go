package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/movie-guru/pkg/db"
	"github.com/movie-guru/pkg/types"
	"github.com/redis/go-redis/v9"

	"github.com/movie-guru/pkg/metrics"
)

var (
	redisStore *redis.Client
	metadata   *db.Metadata
	appConfig  = map[string]string{
		"CORS_HEADERS": "Content-Type",
	}
	corsOrigins []string
)

func StartServer(ulh *UserLoginHandler, m *db.Metadata, deps *Dependencies) error {
	metadata = m
	setupSessionStore()

	corsOrigins = strings.Split(metadata.CorsOrigin, ",")
	for i := range corsOrigins {
		corsOrigins[i] = strings.TrimSpace(corsOrigins[i])
	}

	loginMeters := metrics.NewLoginMeters()
	hcMeters := metrics.NewHCMeters()
	http.HandleFunc("/", createHealthCheckHandler(deps, hcMeters))
	http.HandleFunc("/chat", createChatHandler(deps))
	http.HandleFunc("/history", createHistoryHandler())
	http.HandleFunc("/preferences", createPreferencesHandler(deps.DB))
	http.HandleFunc("/startup", createStartupHandler(deps))
	http.HandleFunc("/login", createLoginHandler(ulh, loginMeters))
	http.HandleFunc("/logout", logoutHandler)
	return http.ListenAndServe(":8080", nil)
}

func setupSessionStore() {
	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	REDIS_PORT := os.Getenv("REDIS_PORT")

	redisStore = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
		Password: REDIS_PASSWORD,
		DB:       0,
	})
	if err := redisStore.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
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

func createHealthCheckHandler(deps *Dependencies, meters *metrics.HCMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")
		addResponseHeaders(w, origin)
		if r.Method == "GET" {
			meters.HCCounter.Add(r.Context(), 1)
			json.NewEncoder(w).Encode("OK")
			return

		}
	}
}

func createLoginHandler(ulh *UserLoginHandler, meters *metrics.LoginMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		origin := r.Header.Get("Origin")
		errLogPrefix := "Error: LoginHandler: "
		if r.Method == "POST" {
			startTime := time.Now()
			defer func() {
				meters.LoginLatencyHistogram.Record(ctx, int64(time.Since(startTime).Milliseconds()))
			}()

			meters.LoginCounter.Add(ctx, 1)

			user := r.Header.Get("user")
			if user == "" {
				log.Println(errLogPrefix, "No auth header")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := ulh.HandleLogin(ctx, user)
			if err != nil {
				if _, ok := err.(*AuthorizationError); ok {
					log.Println(errLogPrefix, "Unauthorized. ", err.Error())
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				log.Println(errLogPrefix, "Cannot get user from db", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				meters.LoginErrorCounter.Add(ctx, 1)
				return
			}
			meters.LoginSuccessCounter.Add(ctx, 1)
			sessionID := uuid.New().String()
			session := &SessionInfo{
				ID:            sessionID,
				User:          user,
				Authenticated: true,
			}
			sessionJSON, err := json.Marshal(session)
			if err != nil {
				log.Println(errLogPrefix, "error decoding session to json", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			err = redisStore.Set(r.Context(), sessionID, sessionJSON, 0).Err()
			if err != nil {
				log.Println(errLogPrefix, "error setting context in redis", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			setCookieHeader := ""
			if os.Getenv("LOCAL") == "true" {
				setCookieHeader = fmt.Sprintf("session=%s; HttpOnly; SameSite=Lax; Path=/; Domain=localhost; Max-Age=86400", sessionID)
			} else {
				setCookieHeader = fmt.Sprintf("session=%s; HttpOnly; Secure; SameSite=None; Path=/; Domain=%s; Max-Age=86400", sessionID, metadata.FrontEndDomain)
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
	log.Println(w.Header())
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
		errLogPrefix := "Error: StartupHandler: "
		var err error
		ctx := r.Context()
		origin := r.Header.Get("Origin")
		addResponseHeaders(w, origin)
		sessionInfo := &SessionInfo{}
		if r.Method != "OPTIONS" {
			sessionInfo, err = getSessionInfo(ctx, r)
			if err != nil {
				if err, ok := err.(*AuthorizationError); ok {
					log.Println(errLogPrefix, "Unauthorized")
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				log.Println(errLogPrefix, "Cannot get session info ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !sessionInfo.Authenticated {
				log.Println(errLogPrefix, "Unauthenticated")
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}
		if r.Method == "GET" {
			addResponseHeaders(w, origin)
			user := sessionInfo.User
			pref, err := deps.DB.GetCurrentProfile(ctx, user)
			if err != nil {
				log.Println(errLogPrefix, "Cannot get preferences for user:", user, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			context, err := deps.MovieRetrieverFlowClient.RetriveDocuments(ctx, randomisedFeaturedFilmsQuery())
			if err != nil {
				log.Println(errLogPrefix, err.Error())
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

func createPreferencesHandler(MovieDB *db.MovieDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errLogPrefix := "Error: PreferencesHandler: "
		var err error
		ctx := r.Context()
		origin := r.Header.Get("Origin")
		addResponseHeaders(w, origin)
		sessionInfo := &SessionInfo{}
		if r.Method != "OPTIONS" {
			sessionInfo, err = getSessionInfo(ctx, r)
			if err != nil {
				if err, ok := err.(*AuthorizationError); ok {
					log.Println(errLogPrefix, "Unauthorized")
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				log.Println(errLogPrefix, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !sessionInfo.Authenticated {
				log.Println(errLogPrefix, "Unauthenticated")
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}
		if r.Method == "GET" {
			addResponseHeaders(w, origin)
			user := sessionInfo.User
			pref, err := MovieDB.GetCurrentProfile(ctx, user)
			if err != nil {
				log.Println(errLogPrefix, "Cannot get preferences for user:", user, err.Error())
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
				log.Println(errLogPrefix, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = MovieDB.UpdateProfile(ctx, pref.Content, sessionInfo.User)
			if err != nil {
				log.Println(errLogPrefix, err.Error())
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
		errLogPrefix := "Error: HistoryHandler: "
		ctx := r.Context()
		origin := r.Header.Get("Origin")
		var err error
		addResponseHeaders(w, origin)
		sessionInfo := &SessionInfo{}
		if r.Method != "OPTIONS" {
			sessionInfo, err = getSessionInfo(ctx, r)
			if err != nil {
				if err, ok := err.(*AuthorizationError); ok {
					log.Println(errLogPrefix, "Unauthorized")
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				log.Println(errLogPrefix, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !sessionInfo.Authenticated {
				log.Println(errLogPrefix, "Unauthenticated")
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}
		if r.Method == "GET" {
			addResponseHeaders(w, origin)
			user := sessionInfo.User
			ch, err := getHistory(ctx, user)
			if err != nil {
				log.Println(errLogPrefix, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			simpleHistory, err := types.ParseRecentHistory(ch.GetHistory(), metadata.HistoryLength)
			if err != nil {
				log.Println(errLogPrefix, err.Error())
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
				log.Println(errLogPrefix, err.Error())
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

func createChatHandler(deps *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errLogPrefix := "Error: ChatHandler: "
		var err error
		ctx := r.Context()
		origin := r.Header.Get("Origin")
		addResponseHeaders(w, origin)
		sessionInfo := &SessionInfo{}
		if r.Method != "OPTIONS" {
			sessionInfo, err = getSessionInfo(ctx, r)
			if err != nil {
				if err, ok := err.(*AuthorizationError); ok {
					log.Println(errLogPrefix, "Unauthorized")
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				log.Println(errLogPrefix, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !sessionInfo.Authenticated {
				log.Println(errLogPrefix, "Unauthenticated")
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}
		if r.Method == "POST" {
			addResponseHeaders(w, origin)
			user := sessionInfo.User
			chatRequest := &ChatRequest{
				Content: "",
			}
			err := json.NewDecoder(r.Body).Decode(chatRequest)
			if err != nil {
				log.Println(errLogPrefix, err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if len(chatRequest.Content) > metadata.MaxUserMessageLen {
				log.Println(errLogPrefix, "Message too long")
				http.Error(w, "Message too long", http.StatusBadRequest)
				return
			}
			ch, err := getHistory(ctx, user)
			if err != nil {
				log.Println(errLogPrefix, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			agentResp := chat(ctx, deps, metadata, ch, user, chatRequest.Content)
			saveHistory(ctx, ch, user)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(agentResp)
			return

		}
		if r.Method == "OPTIONS" {
			addResponseHeaders(w, origin)
			handleOptions(w, origin)
			return
		}
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	errLogPrefix := "Error: LogoutHandler: "
	var err error
	ctx := r.Context()
	origin := r.Header.Get("Origin")
	addResponseHeaders(w, origin)
	sessionInfo := &SessionInfo{}
	if r.Method != "OPTIONS" {
		sessionInfo, err = getSessionInfo(ctx, r)
		if err != nil {
			if err, ok := err.(*AuthorizationError); ok {
				log.Println(errLogPrefix, "Unauthorized")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			log.Println(errLogPrefix, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !sessionInfo.Authenticated {
			log.Println(errLogPrefix, "Unauthenticated")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}
	if r.Method == "GET" {
		addResponseHeaders(w, origin)
		err := deleteSessionInfo(ctx, sessionInfo.ID)
		if err != nil {
			log.Println(errLogPrefix, err.Error())
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
