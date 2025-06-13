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
	"encoding/json"
	"net/http"
	"log/slog"
	"runtime"
	"context"
	"time"

	metrics "github.com/movie-guru/pkg/metrics"
)

func createHealthCheckHandler(deps *Dependencies, meters *metrics.HCMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method == "GET" {
			startTime := time.Now()
			defer func() {
				meters.HCLatency.Record(ctx, int64(time.Since(startTime).Milliseconds()))
			}()
			meters.HCCounter.Add(r.Context(), 1)

			// creating large slice and adding a value so it actually gets allocated.
			largeSlice := createLargeSlice(ctx)
			largeSlice[1000] = 100
			json.NewEncoder(w).Encode("OK")
			return
		}
	}
}

// createLargeSlice attempts to allocate a very large slice of integers.
// It will likely cause the program to crash with an out-of-memory error.
func createLargeSlice(ctx context.Context) []int {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "Recovered from a panic during memory allocation", "error", r)
			// It's good practice to explicitly free memory if possible,
			// though in an OOM situation, the GC has more work to do.
			runtime.GC()
		}
	}()

	// Attempt to create a slice with a very large size.
	// 200 million integers will require approximately 1.6 GB of memory
	// on a 64-bit system (200,000,000 * 8 bytes/int).
	const size = 200 * 1000 * 1000
	slog.InfoContext(ctx, "Attempting to allocate a large slice", "size", size)

	// In Go, slices are created with `make`. This allocation
	// will likely fail if sufficient memory is not available.
	largeSlice := make([]int, size)

	return largeSlice
}