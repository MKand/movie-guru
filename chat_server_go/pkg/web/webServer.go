package web

import (
	"context"

	"net/http"

	"strings"

	"github.com/movie-guru/pkg/db"
	metrics "github.com/movie-guru/pkg/metrics"
)

func enableCORS(allowedOrigins []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if the origin is in the allowed list
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Set other CORS headers
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Allow credentials

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}

func StartServer(ctx context.Context, ulh *UserLoginHandler, metadata *db.Metadata, deps *Dependencies) error {
	setupSessionStore(ctx)

	corsOrigins := strings.Split(metadata.CorsOrigin, ",")
	for i := range corsOrigins {
		corsOrigins[i] = strings.TrimSpace(corsOrigins[i])
	}

	loginMeters := metrics.NewLoginMeters()
	hcMeters := metrics.NewHCMeters()
	chatMeters := metrics.NewChatMeters()

	mux := http.NewServeMux()

	http.HandleFunc("/", createHealthCheckHandler(deps, hcMeters))
	mux.HandleFunc("/chat", createChatHandler(deps, chatMeters, metadata))
	mux.HandleFunc("/history", createHistoryHandler(metadata))
	mux.HandleFunc("/preferences", createPreferencesHandler(deps.DB))
	mux.HandleFunc("/startup", createStartupHandler(deps))
	mux.HandleFunc("/login", createLoginHandler(ulh, loginMeters, metadata))
	mux.HandleFunc("/logout", logoutHandler)
	return http.ListenAndServe(":8080", enableCORS(corsOrigins, mux))
}
