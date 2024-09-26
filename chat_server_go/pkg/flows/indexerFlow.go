package flows

import (
	"context"
	"encoding/json"
	"fmt"
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

			// Write code that generates an embedding
			// - Step 1: Create an embedding from the aiDoc
			// - Step 2: Write a SQL statement to insert the embedding along with the other fields in the table.
			// - HINT: Look at the schema for the table to understand what fields are required.
			// - Take inspiration from the indexer here: https://github.com/firebase/genkit/blob/main/go/samples/pgvector/main.go
			return aiDoc, nil
		})
	return indexerFlow
}

// createText creates a JSON string representation of the important fields in a MovieContext object.
// This string is used as the content for the AI document that is used for embedding.
func createText(movie *types.MovieContext) string {
	dataDict := map[string]interface{}{
		"title":        movie.Title,
		"runtime_mins": movie.RuntimeMinutes,
		"genres": func() string {
			if len(movie.Genres) > 0 {
				return strings.Join(movie.Genres, ", ") // Assuming you want to join genres with commas
			}
			return ""
		}(),
		"rating": func() interface{} {
			if movie.Rating > 0 {
				return fmt.Sprintf("%.1f", movie.Rating)
			}
			return ""
		}(),
		"released": func() interface{} {
			if movie.Released > 0 {
				return movie.Released
			}
			return ""
		}(),
		"actors": func() string {
			if len(movie.Actors) > 0 {
				return strings.Join(movie.Actors, ", ") // Assuming you want to join actors with commas
			}
			return ""
		}(),
		"director": func() string {
			if movie.Director != "" {
				return movie.Director
			}
			return ""
		}(),
		"plot": func() string {
			if movie.Plot != "" {
				return strings.ReplaceAll(movie.Plot, "\n", "")
			}
			return ""
		}(),
	}

	jsonData, _ := json.Marshal(dataDict)
	stringData := string(jsonData)
	return stringData
}
