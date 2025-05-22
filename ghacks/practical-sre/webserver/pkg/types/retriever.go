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

package types

type MovieContext struct {
	Title          string   `json:"title"`
	RuntimeMinutes int      `json:"runtime_minutes"`
	Genres         []string `json:"genres"`
	Rating         float32  `json:"rating"`
	Plot           string   `json:"plot"`
	Released       int      `json:"released"`
	Director       string   `json:"director"`
	Actors         []string `json:"actors"`
	Poster         string   `json:"poster"`
	Tconst         string   `json:"tconst"`
}
