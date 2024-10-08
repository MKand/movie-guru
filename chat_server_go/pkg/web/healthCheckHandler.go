package web

import (
	"encoding/json"
	"net/http"
	"time"

	metrics "github.com/movie-guru/pkg/metrics"
)

func createHealthCheckHandler(deps *Dependencies, meters *metrics.HCMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		origin := r.Header.Get("Origin")
		addResponseHeaders(w, origin)
		if r.Method == "GET" {
			startTime := time.Now()
			defer func() {
				meters.HCLatency.Record(ctx, int64(time.Since(startTime).Milliseconds()))
			}()

			meters.HCCounter.Add(r.Context(), 1)
			json.NewEncoder(w).Encode("OK")
			return

		}
	}
}
