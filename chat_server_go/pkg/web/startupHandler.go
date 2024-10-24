package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	db "github.com/movie-guru/pkg/db"

	m "github.com/movie-guru/pkg/metrics"
	"github.com/movie-guru/pkg/types"
)

func createStartupHandler(deps *Dependencies, meters *m.StartupMeters, metadata *db.Metadata) http.HandlerFunc {
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
		if r.Method == "GET" {
			meters.StartupCounter.Add(ctx, 1)
			startTime := time.Now()
			defer func() {
				meters.StartupLatencyHistogram.Record(ctx, int64(time.Since(startTime).Milliseconds()))
			}()

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
			meters.StartupSuccessCounter.Add(ctx, 1)
			json.NewEncoder(w).Encode(agentResp)
			return

		}
	}
}
