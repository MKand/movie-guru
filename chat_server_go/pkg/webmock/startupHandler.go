package webmock

import (
	"net/http"

	m "github.com/movie-guru/pkg/metrics"
)

func createStartupHandler(deps *Dependencies, meters *m.StartupMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method == "GET" {
			meters.StartupCounter.Add(ctx, 1)
			success := PickSuccess(deps.CurrentProbMetrics.StartupSuccess)
			latency := PickLatencyValue(deps.CurrentProbMetrics.StartupLatencyMinMS, deps.CurrentProbMetrics.StartupLatencyMaxMS)
			defer func() {
				meters.StartupLatencyHistogram.Record(ctx, int64(latency))
			}()

			if success {
				meters.StartupSuccessCounter.Add(ctx, 1)
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
