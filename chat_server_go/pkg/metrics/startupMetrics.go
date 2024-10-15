package metrics

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type StartupMeters struct {
	StartupCounter          metric.Int64Counter
	StartupSuccessCounter   metric.Int64Counter
	StartupLatencyHistogram metric.Int64Histogram
}

func NewStartupMeters() *StartupMeters {
	meter := otel.Meter("startup-handler")

	startupCounter, err := meter.Int64Counter("movieguru_startup_attempts_total", metric.WithDescription("Total number of startup attempts"))
	if err != nil {
		log.Printf("Error creating startupCounter: %v", err)
	}
	startupSuccessCounter, err := meter.Int64Counter("movieguru_startup_success_total", metric.WithDescription("Total number of successful startups"))
	if err != nil {
		log.Printf("Error creating startup success counter: %v", err)
	}

	startupLatencyHistogram, err := meter.Int64Histogram("movieguru_startup_latency", metric.WithDescription("Histogram of startup request latency"),
		metric.WithUnit("ms"),
		metric.WithExplicitBucketBoundaries(1, 2, 5, 10, 20, 50, 100, 200, 500, 1000, 10000),
	)
	if err != nil {
		log.Printf("Error creating startup latency histogram: %v", err)
	}
	return &StartupMeters{
		StartupCounter:          startupCounter,
		StartupSuccessCounter:   startupSuccessCounter,
		StartupLatencyHistogram: startupLatencyHistogram,
	}
}
