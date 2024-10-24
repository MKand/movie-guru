package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	db "github.com/movie-guru/pkg/db"

	m "github.com/movie-guru/pkg/metrics"
)

func createLogoutHandler(meters *m.LogoutMeters, metadata *db.Metadata) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := r.Context()
		origin := r.Header.Get("Origin")
		addResponseHeaders(w, origin)
		sessionInfo := &SessionInfo{}
		if r.Method != "OPTIONS" {
			var shouldReturn bool
			sessionInfo, shouldReturn = authenticateAndGetSessionInfo(ctx, sessionInfo, err, r, w, metadata)
			if shouldReturn {
				return
			}
		}
		user := sessionInfo.User
		if r.Method == "GET" {
			meters.LogoutCounter.Add(ctx, 1)
			startTime := time.Now()
			defer func() {
				meters.LogoutLatencyHistogram.Record(ctx, int64(time.Since(startTime).Milliseconds()))
			}()
			addResponseHeaders(w, origin)
			err := deleteSessionInfo(ctx, sessionInfo.ID)
			if err != nil {
				slog.ErrorContext(ctx, "Error while deleting session info", slog.String("user", user), slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			meters.LogoutSuccessCounter.Add(ctx, 1)
			json.NewEncoder(w).Encode(map[string]string{"logout": "success"})
			return
		}
		if r.Method == "OPTIONS" {
			addResponseHeaders(w, origin)
			handleOptions(w, origin)
			return
		}
	}
}
