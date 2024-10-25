package web

import (
	"net/http"
	"time"

	m "github.com/movie-guru/pkg/metrics"
)

func createLogoutHandler(deps *Dependencies, meters *m.LogoutMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := r.Context()
		if r.Method == "GET" {
			meters.LogoutCounter.Add(ctx, 1)
			latency := PickLatencyValue(deps.CurrentProbMetrics.LoginLatencyMinMS, deps.CurrentProbMetrics.LoginLatencyMaxMS)
			startTime := time.Now()
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
