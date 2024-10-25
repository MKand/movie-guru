package metrics

import (
	"log"

	"go.opentelemetry.io/otel/metric"
)

type LoginMeters struct {
	LoginCounter          metric.Int64Counter
	LoginSuccessCounter   metric.Int64Counter
	LoginLatencyHistogram metric.Int64Histogram
}

func NewLoginMeters(meter metric.Meter) *LoginMeters {
	loginCounter, err := meter.Int64Counter("movieguru_login_attempts_total", metric.WithDescription("Total number of login attempts"))
	if err != nil {
		log.Printf("Error creating login counter: %v", err)
	}
	loginSuccessCounter, err := meter.Int64Counter("movieguru_login_success_total", metric.WithDescription("Total number of successful logins"))
	if err != nil {
		log.Printf("Error creating login success counter: %v", err)
	}

	loginLatencyHistogram, err := meter.Int64Histogram("movieguru_login_latency", metric.WithDescription("Histogram of login request latency"),
		metric.WithUnit("ms"),
		metric.WithExplicitBucketBoundaries(0.05, 0.1, 0.5, 1, 10, 50, 100, 200, 500, 1000, 5000),
	)
	if err != nil {
		log.Printf("Error creating login latency histogram: %v", err)
	}
	return &LoginMeters{
		LoginCounter:          loginCounter,
		LoginSuccessCounter:   loginSuccessCounter,
		LoginLatencyHistogram: loginLatencyHistogram,
	}
}
