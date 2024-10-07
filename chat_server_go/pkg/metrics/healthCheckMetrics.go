package metrics

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type HCMeters struct {
	HCCounter metric.Int64Counter
	HCLatency metric.Int64Histogram
}

func NewHCMeters() *HCMeters {
	meter := otel.Meter("healthcheck-handler")

	hcCounter, err := meter.Int64Counter("movieguru_healthcheck_attempts_total", metric.WithDescription("Total number of healthcheck attempts"))
	if err != nil {
		log.Printf("Error creating hc counter: %v", err)
	}
	hcLatencyHistogram, err := meter.Int64Histogram("movieguru_healthcheck_latency", metric.WithDescription("Histogram of healthcheck request latency"))
	if err != nil {
		log.Printf("Error creating hc latency histogram: %v", err)
	}
	return &HCMeters{
		HCCounter: hcCounter,
		HCLatency: hcLatencyHistogram,
	}
}
