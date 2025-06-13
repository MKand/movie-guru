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
	"time"
	"log/slog"

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
			// creating large array and adding a value so it actually gets allocated.
			largeArray := createLargeArray(ctx)
			largeArray[1000] = 100
			meters.HCCounter.Add(r.Context(), 1)
			json.NewEncoder(w).Encode("OK")
			return
		}
	}
}

func createLargeArray(ctx context.Context) *int[]{
	  try {
            // Attempt to create an array with a very large size.
            int size = 200 * 1000 * 1000; 
            int[] largeArray = new int[size];
			return largeArray
            largeArray[size - 1] = 1;

        } catch (OutOfMemoryError e) {
            slog.ErrorContext(ctx, "OutOfMemoryError caught! Failed to allocate the large array.");
        } catch (NegativeArraySizeException e) {
            slog.ErrorContext(ctx, "NegativeArraySizeException caught! The requested size is too large and wrapped around to a negative number.");
        }
    }
