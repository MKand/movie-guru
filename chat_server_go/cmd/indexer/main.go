package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/vertexai"

	"github.com/movie-guru/pkg/db"
	flows "github.com/movie-guru/pkg/flows"
	types "github.com/movie-guru/pkg/types"
)

func main() {
	ctx := context.Background()

	movieAgentDB, err := db.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	defer movieAgentDB.DB.Close()

	app_version := os.Getenv("APP_VERSION")
	app_version = "v1"
	log.Println("Getting metadata for app version: ", app_version)
	metadata, err := movieAgentDB.GetServerMetadata(app_version)
	log.Println(metadata)
	if err != nil {
		log.Fatal(err)
	}
	err = vertexai.Init(ctx, &vertexai.Config{ProjectID: os.Getenv("PROJECT_ID"), Location: os.Getenv("GCLOUD_LOCATION")})
	if err != nil {
		log.Fatal(err)
	}

	embedder := flows.GetEmbedder(metadata.GoogleEmbeddingModelName)
	indexerFlow := flows.GetIndexerFlow(metadata.RetrieverLength, movieAgentDB, embedder)
	fmt.Println(os.Getwd())
	file, err := os.Open("../../../dataset/movies_with_posters.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close() // Ensure the file is closed when done

	// Create a CSV reader
	reader := csv.NewReader(file)
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
			// End of file or error
			break
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

	if err := genkit.Init(ctx, &genkit.Options{FlowAddr: ":3402"}); err != nil {
		log.Fatal(err)
	}
}
