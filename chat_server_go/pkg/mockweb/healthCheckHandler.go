package webmock

import (
	"encoding/json"
	"net/http"
	"time"

	m "github.com/movie-guru/pkg/metrics"
)

func createHealthCheckHandler(deps *Dependencies, meters *m.HCMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method == "GET" {
			meters.HCCounter.Add(r.Context(), 1)
			startTime := time.Now()
			defer func() {
				meters.HCLatency.Record(ctx, int64(time.Since(startTime).Milliseconds()))
			}()
			json.NewEncoder(w).Encode("OK")
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}
