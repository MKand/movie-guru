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

package web

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"

	m "github.com/movie-guru/pkg/metrics"
	"go.opentelemetry.io/otel/attribute"
	metric "go.opentelemetry.io/otel/metric"
)

type Feedback struct {
	Name               string `json:"name"`
	TraceId            string `json:"traceId"`
	SpanId             string `json:"spanId"`
	FeedbackExperience string `json:"feedbackExperience"`
	FeedbackText       string `json:"feedbackText" omitempty`
}

type Acceptance struct {
	Name     string `json:"name"`
	TraceId  string `json:"traceId"`
	SpanId   string `json:"spanId"`
	Accepted string `json:"accepted"`
}

func createFeedbackHandler(URL string, meters *m.ChatMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			ctx := r.Context()

			if URL == "" {
				slog.InfoContext(ctx, "No Feedback URL found. Not forwarding feedback to Genkit")

				json.NewEncoder(w).Encode("No Feedback URL found. Not forwarding feedback to Genkit")
				return
			}
			feedback := &Feedback{}
			err := json.NewDecoder(r.Body).Decode(feedback)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			inputJSON := map[string]interface{}{
				"name":    feedback.Name,
				"traceId": feedback.TraceId,
				"spanId":  feedback.SpanId,
				"feedback": map[string]interface{}{
					"value": feedback.FeedbackExperience,
				},
			}
			meters.CFeedbackCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Feedback", feedback.FeedbackExperience)))

			jsonData, err := json.Marshal(inputJSON)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))
			if err != nil {
				slog.ErrorContext(ctx, "Feedback: Error sending request to Feedback server", err.Error(), err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				slog.ErrorContext(ctx, "Feedback: Error sending request to Feedback server", err.Error(), err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				slog.ErrorContext(ctx, "Feedback: Genkit returned an Error", "errorCode", resp.StatusCode)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode("OK")
			return
		}
	}
}

func createAcceptanceHandler(URL string, meters *m.ChatMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			ctx := r.Context()

			if URL == "" {
				slog.InfoContext(ctx, "No Feedback URL found. Not forwarding feedback to Genkit")

				json.NewEncoder(w).Encode("No Feedback URL found. Not forwarding feedback to Genkit")
				return
			}
			acceptance := &Acceptance{}
			err := json.NewDecoder(r.Body).Decode(acceptance)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			meters.CAcceptanceCounter.Add(ctx, 1)

			inputJSON := map[string]interface{}{
				"name":       acceptance.Name,
				"traceId":    acceptance.TraceId,
				"spanId":     acceptance.SpanId,
				"acceptance": map[string]interface{}{"value": acceptance.Accepted},
			}
			jsonData, err := json.Marshal(inputJSON)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))
			if err != nil {
				slog.ErrorContext(ctx, "Acceptance: Error sending request to Acceptance server", err.Error(), err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				slog.ErrorContext(ctx, "Acceptance: Error sending request to Acceptance server", err.Error(), err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				slog.ErrorContext(ctx, "Acceptance: Genkit returned an Error", "errorCode", resp.StatusCode)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode("OK")
			return
		}
	}
}
