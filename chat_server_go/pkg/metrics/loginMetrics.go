package metrics

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type LoginMeters struct {
	LoginCounter          metric.Int64Counter
	LoginSuccessCounter   metric.Int64Counter
	LoginLatencyHistogram metric.Int64Histogram
}

func NewLoginMeters() *LoginMeters {
	meter := otel.Meter("login-handler")

	loginCounter, err := meter.Int64Counter("movieguru_login_attempts_total", metric.WithDescription("Total number of login attempts"))
	if err != nil {
		log.Printf("Error creating login counter: %v", err)
	}
	loginSuccessCounter, err := meter.Int64Counter("movieguru_login_success_total", metric.WithDescription("Total number of successful logins"))
	if err != nil {
		log.Printf("Error creating login success counter: %v", err)
	}

	loginLatencyHistogram, err := meter.Int64Histogram("movieguru_login_latency", metric.WithDescription("Histogram of login request latency"))
	if err != nil {
		log.Printf("Error creating login latency histogram: %v", err)
	}
	return &LoginMeters{
		LoginCounter:          loginCounter,
		LoginSuccessCounter:   loginSuccessCounter,
		LoginLatencyHistogram: loginLatencyHistogram,
	}
}
