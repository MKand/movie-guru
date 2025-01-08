package web

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/movie-guru/pkg/db"
	types "github.com/movie-guru/pkg/types"
)

func createPreferencesHandler(MovieDB *db.MovieDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := r.Context()
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

			json.NewEncoder(w).Encode(map[string]string{"update": "success"})
			return
		}
	}
}
