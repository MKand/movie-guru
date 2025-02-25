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

package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/movie-guru/pkg/dataset"
	"github.com/movie-guru/pkg/flows"

	"github.com/firebase/genkit/go/plugins/vertexai"

	"github.com/movie-guru/pkg/db"
	types "github.com/movie-guru/pkg/types"
)

func main() {
	ctx := context.Background()

	var appVersion string

	if err := validate(); err != nil {
		log.Fatal(err)
	}

	movieAgentDB, err := db.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	if movieAgentDB != nil {
		defer movieAgentDB.DB.Close()
	}

	if os.Getenv("APP_VERSION") != "" {
		appVersion = os.Getenv("APP_VERSION")
	} else {
		appVersion = "v1"
	}

	log.Println("Getting metadata for app version: ", appVersion)
	metadata, err := movieAgentDB.GetServerMetadata(appVersion)
	log.Println(metadata)
	if err != nil {
		log.Fatal(err)
	}
	err = vertexai.Init(ctx, &vertexai.Config{ProjectID: os.Getenv("PROJECT_ID"), Location: os.Getenv("LOCATION")})
	if err != nil {
		log.Fatal(err)
	}

	embedder := flows.GetEmbedder(metadata.GoogleEmbeddingModelName)
	indexerFlow := flows.GetIndexerFlow(metadata.RetrieverLength, movieAgentDB, embedder)

	data, err := dataset.Asset("../dataset/movies_with_posters.csv")

	// Create a CSV reader
	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = '\t' // Set the delimiter to tab

	// Read the header row (if present)
	header, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading header:", err)
		return
	}
	fmt.Println("Header:", header)
	index := 0
	for {
		record, err := reader.Read()
		if err != nil {
			return
		}
		// Process the record (row)
		year, _ := strconv.ParseFloat(record[1], 32)
		rating, _ := strconv.ParseFloat(record[5], 32)
		runtime, _ := strconv.ParseFloat(record[6], 32)
		movieContext := &types.MovieContext{
			Title:          record[0],
			RuntimeMinutes: int(runtime),
			Genres:         strings.Split(record[7], ", "),
			Rating:         float32(rating),
			Plot:           record[4],
			Released:       int(year),
			Director:       record[3],
			Actors:         strings.Split(record[2], ", "),
			Poster:         record[9],
			Tconst:         strconv.Itoa(index),
		}
		indexerFlow.Run(ctx, movieContext)
		index += 1
	}

	fmt.Println("Done indexing")
}

func validate() error {
	if os.Getenv("PROJECT_ID") == "" {
		return fmt.Errorf("PROJECT_ID is not set")
	}
	if os.Getenv("LOCATION") == "" {
		return fmt.Errorf("LOCATION is not set")
	}
	return nil
}
