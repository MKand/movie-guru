package flows

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/movie-guru/pkg/db"

	types "github.com/movie-guru/pkg/types"
)

func GetIndexerFlow(maxRetLength int, movieDB *db.MovieDB, embedder ai.Embedder) *genkit.Flow[*types.MovieContext, *ai.Document, struct{}] {
	indexerFlow := genkit.DefineFlow("movieDocFlow",
		// Uploading one entry (document) at a time
		func(ctx context.Context, doc *types.MovieContext) (*ai.Document, error) {
			time.Sleep(1 / 3 * time.Second)            // reduce rate at which operation is performed to avoid hitting VertexAI rate limits
			content := createText(doc)                 // creates a JSON string representation of the important fields in a MovieContext object.
			aiDoc := ai.DocumentFromText(content, nil) // create an object of type AIDocument from  the content

			// INSTRUCTIONS: Write code that generates an embedding
			// - Step 1: Create an embedding from the aiDoc
			// - Step 2: Write a SQL statement to insert the embedding along with the other fields in the table.
			// - HINT: Look at the schema for the table to understand what fields are required.
			// - Take inspiration from the indexer here: https://github.com/firebase/genkit/blob/main/go/samples/pgvector/main.go

			return aiDoc, nil
		})
	return indexerFlow
}

// createText creates a JSON string representation of the relevant fields in a MovieContext object.
// This string is used as the content for the AI document from which the vector embedding is created.
func createText(movie *types.MovieContext) string {
	dataDict := map[string]interface{}{
		// INSTRUCTIONS: Write code that populates dataDict with relevant fields from raw data.
		// 1. Which other fields from the raw data should the dict contain?
		// 1. Are there any fields in the orginal data that need to be reformatted?

		// Here are two freebies to help you get started.
		"title": movie.Title,
		"genres": func() string {
			if len(movie.Genres) > 0 {
				return strings.Join(movie.Genres, ", ") // Assuming you want to join genres with commas
			}
			return ""
		}(),
	}

	jsonData, _ := json.Marshal(dataDict)
	stringData := string(jsonData)
	return stringData
}
