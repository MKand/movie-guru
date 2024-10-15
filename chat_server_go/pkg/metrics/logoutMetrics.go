package metrics

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type LogoutMeters struct {
	LogoutCounter          metric.Int64Counter
	LogoutSuccessCounter   metric.Int64Counter
	LogoutLatencyHistogram metric.Int64Histogram
}

func NewLogoutMeters() *LogoutMeters {
	meter := otel.Meter("Logout-handler")

	logoutCounter, err := meter.Int64Counter("movieguru_Logout_attempts_total", metric.WithDescription("Total number of Logout attempts"))
	if err != nil {
		log.Printf("Error creating Logout counter: %v", err)
	}
	logoutSuccessCounter, err := meter.Int64Counter("movieguru_Logout_success_total", metric.WithDescription("Total number of successful Logouts"))
	if err != nil {
		log.Printf("Error creating Logout success counter: %v", err)
	}

	logoutLatencyHistogram, err := meter.Int64Histogram("movieguru_Logout_latency", metric.WithDescription("Histogram of Logout request latency"),
		metric.WithUnit("ms"),
		metric.WithExplicitBucketBoundaries(0.1, 0.5, 1, 1.5, 2, 3, 10),
	)
	if err != nil {
		log.Printf("Error creating Logout latency histogram: %v", err)
	}
	return &LogoutMeters{
		LogoutCounter:          logoutCounter,
		LogoutSuccessCounter:   logoutSuccessCounter,
		LogoutLatencyHistogram: logoutLatencyHistogram,
	}
}
