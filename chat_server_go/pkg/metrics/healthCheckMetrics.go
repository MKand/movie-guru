package metrics

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type HCMeters struct {
	HCCounter metric.Int64Counter
}

func NewHCMeters() *HCMeters {
	meter := otel.Meter("hc-handler")

	hcCounter, err := meter.Int64Counter("hc.attempts_total", metric.WithDescription("Total number of hc attempts"))
	if err != nil {
		log.Printf("Error creating hc counter: %v", err)
	}
	return &HCMeters{
		HCCounter: hcCounter,
	}
}
