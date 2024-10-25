package webmock

import (
	"net/http"

	m "github.com/movie-guru/pkg/metrics"
)

func createPreferencesHandler(deps *Dependencies, meters *m.PreferencesMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if r.Method == "GET" {
			meters.PreferencesGetCounter.Add(ctx, 1)
			success := PickSuccess(deps.CurrentProbMetrics.PrefGetSuccess)
			latency := PickLatencyValue(deps.CurrentProbMetrics.PrefGetLatencyMinMS, deps.CurrentProbMetrics.PrefGetLatencyMaxMS)

			defer func() {
				meters.PreferencesGetLatencyHistogram.Record(ctx, int64(latency))
			}()

			if success {
				meters.PreferencesGetSuccessCounter.Add(ctx, 1)
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.Method == "POST" {
			meters.PreferencesUpdateCounter.Add(ctx, 1)
			success := PickSuccess(deps.CurrentProbMetrics.PrefUpdateSuccess)
			latency := PickLatencyValue(deps.CurrentProbMetrics.PrefUpdateLatencyMinMS, deps.CurrentProbMetrics.PrefUpdateLatencyMaxMS)

			defer func() {
				meters.PreferencesUpdateLatencyHistogram.Record(ctx, int64(latency))
			}()

			if success {
				meters.PreferencesUpdateSuccessCounter.Add(ctx, 1)
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
