package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/movie-guru/pkg/db"
	m "github.com/movie-guru/pkg/metrics"
	"github.com/movie-guru/pkg/types"
)

func createPreferencesHandler(MovieDB *db.MovieDB, meters *m.PreferencesMeters) http.HandlerFunc {
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
			meters.PreferencesGetCounter.Add(ctx, 1)
			startTime := time.Now()
			defer func() {
				meters.PreferencesGetLatencyHistogram.Record(ctx, int64(time.Since(startTime).Milliseconds()))
			}()

			pref, err := MovieDB.GetCurrentProfile(ctx, user)
			if err != nil {
				slog.ErrorContext(ctx, "Cannot get preferences", slog.String("user", user), slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			meters.PreferencesGetSuccessCounter.Add(ctx, 1)
			addResponseHeaders(w, origin)
			json.NewEncoder(w).Encode(pref)
			return
		}
		if r.Method == "POST" {
			meters.PreferencesUpdateCounter.Add(ctx, 1)
			startTime := time.Now()
			defer func() {
				meters.PreferencesUpdateLatencyHistogram.Record(ctx, int64(time.Since(startTime).Milliseconds()))
			}()
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

			meters.PreferencesGetSuccessCounter.Add(ctx, 1)
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
