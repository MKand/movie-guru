package webmock

import (
	"net/http"

	m "github.com/movie-guru/pkg/metrics"
)

func createLogoutHandler(deps *Dependencies, meters *m.LogoutMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method == "GET" {
			meters.LogoutCounter.Add(ctx, 1)
			latency := PickLatencyValue(deps.CurrentProbMetrics.LoginLatencyMinMS, deps.CurrentProbMetrics.LoginLatencyMaxMS)
			defer func() {
				meters.LogoutLatencyHistogram.Record(ctx, int64(latency))
			}()
			meters.LogoutSuccessCounter.Add(ctx, 1)
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == "OPTIONS" {
			return
		}
	}
}
