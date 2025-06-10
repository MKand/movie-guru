package main

import (
	"context"
	"log/slog"
	"os"

	met "github.com/movie-guru/pkg/metrics"
	web "github.com/movie-guru/pkg/webmock"
)

func main() {
	ctx := context.Background()

	if shutdown, err := met.SetupOpenTelemetry(ctx); err != nil {
		slog.ErrorContext(ctx, "Error setting up OpenTelemetry", slog.Any("error", err))
		os.Exit(1)
	} else {
		defer shutdown(ctx)
	}

	deps := getDependencies()

	if err := web.StartServer(ctx, deps); err != nil {
		slog.ErrorContext(ctx, "Server exited with error", slog.Any("error", err))
		os.Exit(1)
	}

}

func getDependencies() *web.Dependencies {
	deps := &web.Dependencies{}

	initialPhase := &web.MetricsProb{
		ChatSuccess:       0.7,
		ChatSafetyIssue:   0.2,
		ChatEngaged:       0.5,
		ChatAcknowledged:  0.15,
		ChatRejected:      0.25,
		ChatUnclassified:  0.1,
		ChatSPositive:     0.4,
		ChatSNegative:     0.3,
		ChatSNeutral:      0.1,
		ChatSUnclassified: 0.2,
		LoginSuccess:      0.99,
		StartupSuccess:    0.75,
		PrefUpdateSuccess: 0.84,
		PrefGetSuccess:    0.99,

		LoginLatencyMinMS: 10,
		LoginLatencyMaxMS: 200,

		ChatLatencyMinMS: 1607,
		ChatLatencyMaxMS: 7683,

		StartupLatencyMinMS: 456,
		StartupLatencyMaxMS: 1634,

		PrefGetLatencyMinMS: 153,
		PrefGetLatencyMaxMS: 348,

		PrefUpdateLatencyMinMS: 463,
		PrefUpdateLatencyMaxMS: 745,
	}

	deps.CurrentProbMetrics = initialPhase
	return deps
}
