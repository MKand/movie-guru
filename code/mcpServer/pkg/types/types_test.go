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

import (
	"encoding/json"
	"testing"
)

func TestUserQueryJSON(t *testing.T) {
	query := UserQuery{
		VectorQuery:    "Find action movies",
		KeywordQuery:   "Matrix",
		Total:          5,
		SearchCategory: Vector,
	}

	// Test marshaling
	data, err := json.Marshal(query)
	if err != nil {
		t.Fatalf("Failed to marshal UserQuery: %v", err)
	}

	// Test unmarshaling
	var decoded UserQuery
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal UserQuery: %v", err)
	}

	if decoded.VectorQuery != query.VectorQuery {
		t.Errorf("Expected VectorQuery %s, got %s", query.VectorQuery, decoded.VectorQuery)
	}
	if decoded.SearchCategory != query.SearchCategory {
		t.Errorf("Expected SearchCategory %s, got %s", query.SearchCategory, decoded.SearchCategory)
	}
}

func TestMovieContextSchemaJSON(t *testing.T) {
	movie := MovieContextSchema{
		Title:          "The Matrix",
		RuntimeMinutes: 136,
		Genres:         []string{"Action", "Sci-Fi"},
		Rating:         8.7,
		Released:       1999,
		Director:       "Wachowski",
		Actors:         []string{"Keanu Reeves", "Laurence Fishburne"},
		Poster:         "matrix.jpg",
		Tconst:         "tt0133093",
	}

	// Test marshaling
	data, err := json.Marshal(movie)
	if err != nil {
		t.Fatalf("Failed to marshal MovieContextSchema: %v", err)
	}

	// Test unmarshaling
	var decoded MovieContextSchema
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal MovieContextSchema: %v", err)
	}

	if decoded.Title != movie.Title {
		t.Errorf("Expected Title %s, got %s", movie.Title, decoded.Title)
	}
	if len(decoded.Genres) != len(movie.Genres) {
		t.Errorf("Expected %d genres, got %d", len(movie.Genres), len(decoded.Genres))
	}
}

func TestSearchTypeConstants(t *testing.T) {
	// Ensure constants are defined correctly
	if Vector != "Vector" {
		t.Errorf("Expected Vector constant to be 'Vector', got '%s'", Vector)
	}
	if Keyword != "Keyword" {
		t.Errorf("Expected Keyword constant to be 'Keyword', got '%s'", Keyword)
	}
	if Mixed != "Mixed" {
		t.Errorf("Expected Mixed constant to be 'Mixed', got '%s'", Mixed)
	}
}
