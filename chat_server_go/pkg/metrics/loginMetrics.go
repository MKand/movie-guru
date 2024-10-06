package metrics

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type LoginMeters struct {
	LoginCounter          metric.Int64Counter
	LoginSuccessCounter   metric.Int64Counter
	LoginErrorCounter     metric.Int64Counter
	LoginLatencyHistogram metric.Int64Histogram
}

func NewLoginMeters() *LoginMeters {
	meter := otel.Meter("login-handler")

	loginCounter, err := meter.Int64Counter("login_attempts_total", metric.WithDescription("Total number of login attempts"))
	if err != nil {
		log.Printf("Error creating login counter: %v", err)
	}
	loginSuccessCounter, err := meter.Int64Counter("login_success_total", metric.WithDescription("Total number of successful logins"))
	if err != nil {
		log.Printf("Error creating login success counter: %v", err)
	}
	loginErrorCounter, err := meter.Int64Counter("login_errors_total", metric.WithDescription("Total number of login errors"))
	if err != nil {
		log.Printf("Error creating login error counter: %v", err)
	}
	loginLatencyHistogram, err := meter.Int64Histogram("login.latency", metric.WithDescription("Histogram of login request latency"))
	if err != nil {
		log.Printf("Error creating login latency histogram: %v", err)
	}
	return &LoginMeters{
		LoginCounter:          loginCounter,
		LoginSuccessCounter:   loginSuccessCounter,
		LoginErrorCounter:     loginErrorCounter,
		LoginLatencyHistogram: loginLatencyHistogram,
	}
}
