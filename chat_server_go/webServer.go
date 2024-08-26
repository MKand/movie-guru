package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	key        = os.Getenv("SECRET_KEY")
	redisStore *redis.Client
	metadata   *Metadata
	appConfig  = map[string]string{
		"CORS_HEADERS": "Content-Type",
	}
	corsOrigins []string
)

type SessionInfo struct {
	ID            string
	User          string
	Authenticated bool
}

type LoginBody struct {
	InviteCode string `json:"inviteCode" omitempty`
}

type PrefBody struct {
	Content *UserProfile `json:"content"`
}

type ChatRequest struct {
	Content string `json:"content"`
}

func startServer(ulh UserLoginHandler, m *Metadata, chatDeps *ChatDependencies) {
	metadata = m
	setupSessionStore()

	corsOrigins = strings.Split(metadata.CorsOrigin, ",")
	for i := range corsOrigins {
		corsOrigins[i] = strings.TrimSpace(corsOrigins[i])
	}
	http.HandleFunc("/chat", createChatHandler(chatDeps))
	http.HandleFunc("/history", createHistoryHandler())
	http.HandleFunc("/preferences", createPreferencesHandler(chatDeps.DB))
	http.HandleFunc("/startup", createStartupHandler(chatDeps))
	http.HandleFunc("/login", createLoginHandler(ulh))
	http.HandleFunc("/logout", logoutHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))
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
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "user, Origin, Cookie, Accept, Content-Type, Content-Length, Accept-Encoding,Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func createLoginHandler(ulh UserLoginHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		errLogPrefix := "Error: LoginHandler: "
		if r.Method == "POST" {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println(errLogPrefix, "No auth header")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			var loginBody LoginBody
			err := json.NewDecoder(r.Body).Decode(&loginBody)
			if err != nil {
				log.Println(errLogPrefix, "Bad Request at login", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			user, err := ulh.handleLogin(authHeader, loginBody.InviteCode)
			if err != nil {
				if _, ok := err.(*AuthorizationError); ok {
					log.Println(errLogPrefix, "Unauthorized. ", err.Error())
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				log.Println(errLogPrefix, "Cannot get user from db", err)
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
				log.Println(errLogPrefix, "error decoding session to json", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			err = redisStore.Set(r.Context(), sessionID, sessionJSON, 0).Err()
			if err != nil {
				log.Println(errLogPrefix, "error setting context in redis", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			setCookieHeader := fmt.Sprintf("session=%s; HttpOnly; Secure; SameSite=None; Path=/; Domain=%s; Max-Age=86400", sessionID, metadata.FrontEndDomain)
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

func getHistory(ctx context.Context, user string) (*ChatHistory, error) {
	historyJson, err := redisStore.Get(ctx, user).Result()
	ch := NewChatHistory()
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

func saveHistory(ctx context.Context, history *ChatHistory, user string) error {
	history.Trim(metadata.HistoryLength)

	err := redisStore.Set(ctx, user, history, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func deleteHistory(ctx context.Context, user string) error {
	_, err := redisStore.Del(ctx, user).Result()
	if err != nil {
		return err
	}
	return nil
}

func createStartupHandler(deps *ChatDependencies) http.HandlerFunc {
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
			pref, err := getCurrentProfile(ctx, user, deps.DB)
			if err != nil {
				log.Println(errLogPrefix, "Cannot get preferences for user:", user, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			context, err := deps.Retriever.RetriveDocuments(ctx, randomisedFeaturedFilmsQuery())
			agentResp := NewAgentResponse()
			agentResp.Context = context[0:5]
			agentResp.Preferences = pref
			agentResp.Result = SUCCESS
			if err != nil {
				log.Println(errLogPrefix, err.Error())
				http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(agentResp)
			return

		}
	}
}

func createPreferencesHandler(db *sql.DB) http.HandlerFunc {
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
			pref, err := getCurrentProfile(ctx, user, db)
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
				Content: NewUserProfile(),
			}
			err := json.NewDecoder(r.Body).Decode(pref)
			if err != nil {
				log.Println(errLogPrefix, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = updatePreferences(ctx, pref.Content, sessionInfo.User, db)
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
			simpleHistory, err := ParseRecentHistory(ch.GetHistory(), metadata.HistoryLength)
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

func createChatHandler(deps *ChatDependencies) http.HandlerFunc {
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
