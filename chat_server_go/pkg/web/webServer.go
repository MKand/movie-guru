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

	"net/http"

	"strings"

	"github.com/movie-guru/pkg/db"
	metrics "github.com/movie-guru/pkg/metrics"
	"golang.org/x/exp/slog"
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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, ApiKey, User")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
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
		slog.InfoContext(ctx, "Setting cors origin", slog.Any("origin", corsOrigins[i]))

	}

	loginMeters := metrics.NewLoginMeters()
	hcMeters := metrics.NewHCMeters()
	chatMeters := metrics.NewChatMeters()

	mux := http.NewServeMux()

	mux.HandleFunc("/", createHealthCheckHandler(deps, hcMeters))
	mux.HandleFunc("/chat", createChatHandler(deps, chatMeters, metadata))
	mux.HandleFunc("/history", createHistoryHandler(metadata))
	mux.HandleFunc("/preferences", createPreferencesHandler(deps.DB))
	mux.HandleFunc("/startup", createStartupHandler(deps))
	mux.HandleFunc("/login", createLoginHandler(ulh, loginMeters, metadata))
	mux.HandleFunc("/logout", logoutHandler)
	return http.ListenAndServe(":8080", enableCORS(corsOrigins, mux))
}
