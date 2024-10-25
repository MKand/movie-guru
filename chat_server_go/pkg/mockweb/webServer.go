package webmock

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"

	metrics "github.com/movie-guru/pkg/metrics"
	"go.opentelemetry.io/otel"
)

var (
	globalDeps *Dependencies
)

func StartServer(ctx context.Context, deps *Dependencies) error {

	globalDeps = deps
	podName := os.Getenv("POD_NAME")
	if podName == "" {
		podName = "local"
	}
	meter := otel.Meter(podName)

	loginMeters := metrics.NewLoginMeters(meter)
	logoutMeters := metrics.NewLogoutMeters(meter)
	hcMeters := metrics.NewHCMeters(meter)
	chatMeters := metrics.NewChatMeters(meter)
	prefMeters := metrics.NewPreferencesMeters(meter)
	startupMeters := metrics.NewStartupMeters(meter)

	http.HandleFunc("/login", createLoginHandler(globalDeps, loginMeters))
	http.HandleFunc("/logout", createLogoutHandler(globalDeps, logoutMeters))
	http.HandleFunc("/history", createHistoryHandler(globalDeps))
	http.HandleFunc("/chat", createChatHandler(globalDeps, chatMeters))
	http.HandleFunc("/", createHealthCheckHandler(globalDeps, hcMeters))
	http.HandleFunc("/preferences", createPreferencesHandler(globalDeps, prefMeters))
	http.HandleFunc("/startup", createStartupHandler(globalDeps, startupMeters))
	http.HandleFunc("/phase", createTestPhaseHandler())

	return http.ListenAndServe(":8080", nil)
}

func createTestPhaseHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(globalDeps)
			w.WriteHeader(http.StatusOK)
			return

		}
		if r.Method == "POST" {
			m := &MetricsProb{}
			err := json.NewDecoder(r.Body).Decode(m)

			if err != nil {
				slog.InfoContext(ctx, "Error while decoding request", slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			globalDeps.CurrentProbMetrics = m
			json.NewEncoder(w).Encode(globalDeps)
			w.WriteHeader(http.StatusAccepted)
			return

		}
		if r.Method == "OPTIONS" {
			return
		}

	}
}

func PickSuccess(rate float32) bool {

	probabilities := []float32{rate, 1 - rate} // Must add up to 1
	// Generate a random number between 0 and 1
	randomNumber := rand.Float32()

	// Cumulative probability
	cumulativeProbability := float32(0.0)

	// Iterate through the probabilities and check if the random number falls within the range
	for i, probability := range probabilities {
		cumulativeProbability += probability
		if randomNumber <= cumulativeProbability {
			switch i {
			case 0:
				return true
			case 1:
				return false
			}
		}
	}
	// This should not happen, but return a default value just in case
	return false
}

func PickLatencyValue(minMs int, maxMs int) int {
	// Convert milliseconds to time.Duration for easier calculations
	minDuration := time.Duration(minMs) * time.Millisecond
	maxDuration := time.Duration(maxMs) * time.Millisecond

	// Calculate the difference between max and min
	diff := maxDuration - minDuration

	// Generate a random duration between 0 and the difference
	randomDuration := time.Duration(rand.Int63n(int64(diff)))

	// Add the random duration to the minimum duration to get the final value
	result := minDuration + randomDuration

	// Return the result in milliseconds
	return int(result.Milliseconds())
}
