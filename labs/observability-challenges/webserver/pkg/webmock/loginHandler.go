package webmock

import (
	"net/http"

	_ "github.com/lib/pq"
	m "github.com/movie-guru/pkg/metrics"
)

func createLoginHandler(deps *Dependencies, meters *m.LoginMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method == "POST" {
			meters.LoginCounter.Add(ctx, 1)

			success := PickSuccess(deps.CurrentProbMetrics.LoginSuccess)
			latency := PickLatencyValue(deps.CurrentProbMetrics.LoginLatencyMinMS, deps.CurrentProbMetrics.LoginLatencyMaxMS)

			defer func() {
				meters.LoginLatencyHistogram.Record(ctx, int64(latency))
			}()

			if success {
				meters.LoginSuccessCounter.Add(ctx, 1)
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		if r.Method == "OPTIONS" {
			return
		}
	}
}
