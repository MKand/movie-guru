package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func logoutHandler(w http.ResponseWriter, r *http.Request) {
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
		err := deleteSessionInfo(ctx, sessionInfo.ID)
		if err != nil {
			slog.ErrorContext(ctx, "Error while deleting session info", slog.String("user", user), slog.Any("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"logout": "success"})

		return
	}

}
