package metrics

import (
	"log"

	"go.opentelemetry.io/otel/metric"
)

type PreferencesMeters struct {
	PreferencesGetCounter           metric.Int64Counter
	PreferencesUpdateCounter        metric.Int64Counter
	PreferencesGetSuccessCounter    metric.Int64Counter
	PreferencesUpdateSuccessCounter metric.Int64Counter

	PreferencesUpdateLatencyHistogram metric.Int64Histogram
	PreferencesGetLatencyHistogram    metric.Int64Histogram
}

func NewPreferencesMeters(meter metric.Meter) *PreferencesMeters {

	preferencesGetCounter, err := meter.Int64Counter("movieguru_prefGet_attempts_total", metric.WithDescription("Total number of pref get attempts"))
	if err != nil {
		log.Printf("Error creating preferencesGetCounter: %v", err)
	}
	preferencesUpdateCounter, err := meter.Int64Counter("movieguru_prefUpdate_attempts_total", metric.WithDescription("Total number of pref update attempts"))
	if err != nil {
		log.Printf("Error creating preferencesUpdateCounter: %v", err)
	}
	preferencesGetSuccessCounter, err := meter.Int64Counter("movieguru_prefGet_success_total", metric.WithDescription("Total number of successful pref get attempts"))
	if err != nil {
		log.Printf("Error creating preferencesGetSuccessCounter: %v", err)
	}
	preferencesUpdateSuccessCounter, err := meter.Int64Counter("movieguru_prefUpdate_success_total", metric.WithDescription("Total number of successful pref update attempts"))
	if err != nil {
		log.Printf("Error creating preferencesUpdateSuccessCounter: %v", err)
	}

	prefGetLatencyHistogram, err := meter.Int64Histogram("movieguru_prefGet_latency", metric.WithDescription("Histogram of pref get request latency"),
		metric.WithUnit("ms"),
		metric.WithExplicitBucketBoundaries(0.1, 0.5, 1, 1.5, 2, 3, 10),
	)
	if err != nil {
		log.Printf("Error creating prefGetLatencyHistogram: %v", err)
	}

	prefUpdateLatencyHistogram, err := meter.Int64Histogram("movieguru_prefUpdate_latency", metric.WithDescription("Histogram of pref update request latency"),
		metric.WithUnit("ms"),
		metric.WithExplicitBucketBoundaries(0.5, 1, 1.5, 2, 2.5, 3, 3.5, 10, 100),
	)
	if err != nil {
		log.Printf("Error creating prefUpdateLatencyHistogram: %v", err)
	}
	return &PreferencesMeters{
		PreferencesGetCounter:             preferencesGetCounter,
		PreferencesUpdateCounter:          preferencesUpdateCounter,
		PreferencesGetSuccessCounter:      preferencesGetSuccessCounter,
		PreferencesUpdateSuccessCounter:   preferencesUpdateSuccessCounter,
		PreferencesUpdateLatencyHistogram: prefUpdateLatencyHistogram,
		PreferencesGetLatencyHistogram:    prefGetLatencyHistogram,
	}
}
