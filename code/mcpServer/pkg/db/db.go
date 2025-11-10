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

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"unsafe"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mkand/movieSearch/pkg/embedder"
	"github.com/mkand/movieSearch/pkg/types"
	"github.com/firebase/genkit/go/genkit"
)

type MovieDB struct {
	DB       *sql.DB
	Embedder *embedder.Embedder
}

func GetDB(ctx context.Context, g *genkit.Genkit) (*MovieDB, error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	// Initialize embedder
	emb, err := embedder.NewEmbedder(ctx, g)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedder: %w", err)
	}
	
	return &MovieDB{
		DB:       db,
		Embedder: emb,
	}, nil
}

func connectToDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "movieDB.db")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)

	err = db.PingContext(context.Background())
	if err != nil {
		log.Fatal("Error pinging database: %v", err)
	}
	slog.Log(context.Background(), slog.LevelInfo, "DB pinged successfully")
	return db, nil
}

func (db *MovieDB) PerformSearch(query types.UserQuery) ([]types.MovieContextSchema, error) {
	var results []types.MovieContextSchema
	var rows *sql.Rows
	var err error

	limit := 10
	if query.Total > 0 {
		limit = query.Total
	}

	switch query.SearchCategory {
	case types.Keyword:
		// Keyword search using FTS or LIKE
		sqlQuery := `
			SELECT title, runtime_mins, genres, rating, released, director, actors, poster, tconst
			FROM movies
			WHERE title LIKE ? OR plot LIKE ? OR director LIKE ? OR actors LIKE ?
			LIMIT ?`
		
		searchPattern := "%" + query.KeywordQuery + "%"
		rows, err = db.DB.Query(sqlQuery, searchPattern, searchPattern, searchPattern, searchPattern, limit)
		
	case types.Vector:
		// Vector search using sqlite-vec
		// Generate embedding from the query text
		embedding, err := db.Embedder.GenerateEmbedding(context.Background(), query.VectorQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding: %w", err)
		}
		
		sqlQuery := `
			SELECT m.title, m.runtime_mins, m.genres, m.rating, m.released, m.director, m.actors, m.poster, m.tconst
			FROM vec_movies v
			INNER JOIN movies m ON v.tconst = m.tconst
			WHERE v.embedding MATCH ?
			ORDER BY distance
			LIMIT ?`
		
		// Convert embedding to bytes for sqlite-vec
		// sqlite-vec expects a blob of float32 values
		embeddingBytes := (*[1 << 30]byte)(unsafe.Pointer(&embedding[0]))[:len(embedding)*4:len(embedding)*4]
		rows, err = db.DB.Query(sqlQuery, embeddingBytes, limit)
		
	case types.Mixed:
		// Mixed search - combine keyword and vector
		// This is a simplified version; you might want to merge and deduplicate results
		sqlQuery := `
			SELECT title, runtime_mins, genres, rating, released, director, actors, poster, tconst
			FROM movies
			WHERE title LIKE ? OR plot LIKE ?
			LIMIT ?`
		
		searchPattern := "%" + query.KeywordQuery + "%"
		rows, err = db.DB.Query(sqlQuery, searchPattern, searchPattern, limit)
		
	default:
		return nil, fmt.Errorf("invalid search category: %v", query.SearchCategory)
	}

	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Parse results
	for rows.Next() {
		var movie types.MovieContextSchema
		var genresJSON, actorsJSON string
		
		err := rows.Scan(
			&movie.Title,
			&movie.RuntimeMinutes,
			&genresJSON,
			&movie.Rating,
			&movie.Released,
			&movie.Director,
			&actorsJSON,
			&movie.Poster,
			&movie.Tconst,
		)
		if err != nil {
			slog.Error("failed to scan row:", "error", err)
			continue
		}

		// Parse JSON arrays for genres and actors
		if err := json.Unmarshal([]byte(genresJSON), &movie.Genres); err != nil {
			slog.Warn("failed to parse genres", "error", err)
			movie.Genres = []string{}
		}
		if err := json.Unmarshal([]byte(actorsJSON), &movie.Actors); err != nil {
			slog.Warn("failed to parse actors", "error", err)
			movie.Actors = []string{}
		}

		results = append(results, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}