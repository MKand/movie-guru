// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
