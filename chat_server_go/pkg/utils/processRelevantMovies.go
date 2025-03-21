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

package utils

import (
	"fmt"
	"os"

	types "github.com/movie-guru/pkg/types"
)

func FilterRelevantContext(relevantMovies []string, fullContext []*types.MovieContext) []*types.MovieContext {
	relevantContext := make(
		[]*types.MovieContext,
		0,
		len(relevantMovies),
	)
	for _, m := range fullContext {
		for _, r := range relevantMovies {
			if r == m.Title {
				if m.Poster != "" {
					relevantContext = append(relevantContext, m)
				}

			}
		}
	}
	return relevantContext
}

func AddPosterURLs(contextDocuments []*types.MovieContext) error {
	// make this defensive
	projectId := os.Getenv("PROJECT_ID")
	region := os.Getenv("LOCATION")

	for _, c := range contextDocuments {
		if c.Poster != "" {
			c.Poster = fmt.Sprintf("https://storage.googleapis.com/%s_%s_posters/%s", projectId, region, c.Poster)
		}
		if os.Getenv("USE_SIGNED_URL") != "" {
			var err error
			c.Poster, err = GetSignedURL(c.Poster)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
