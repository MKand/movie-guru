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
			largeArray := createLargeArray()
			largeArray[100] = 10
			meters.HCCounter.Add(r.Context(), 1)
			json.NewEncoder(w).Encode("OK")
			return
		}
	}
}

func createLargeArray() *int[]{
	  try {
            // Attempt to create an array with a very large size.
            // Integer.MAX_VALUE is the maximum positive value for an int,
            // which is the theoretical max size for an array dimension in Java.
            // In practice, you'll run out of memory far before this.
            // Let's try a smaller, but still very large, number.
            // For example, 200 million integers.
            int size = 200 * 1000 * 1000; // 200 million integers
            System.out.println("Attempting to allocate an int array of size: " + size);
            int[] largeArray = new int[size];
            System.out.println("Successfully allocated int array of size: " + size);
			return largeArray
            // To prove it's allocated, you could try to access an element (optional)
            // largeArray[size - 1] = 1;
            // System.out.println("Accessed last element: " + largeArray[size-1]);

        } catch (OutOfMemoryError e) {
            System.err.println("OutOfMemoryError caught! Failed to allocate the large array.");
            e.printStackTrace();
        } catch (NegativeArraySizeException e) {
            System.err.println("NegativeArraySizeException caught! The requested size is too large and wrapped around to a negative number.");
            e.printStackTrace();
        }
    }
