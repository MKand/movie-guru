package flows

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/movie-guru/pkg/db"
	pgv "github.com/pgvector/pgvector-go"

	types "github.com/movie-guru/pkg/types"
)

func GetIndexerFlow(maxRetLength int, movieDB *db.MovieDB, embedder ai.Embedder) *genkit.Flow[*types.MovieContext, *ai.Document, struct{}] {
	indexerFlow := genkit.DefineFlow("movieDocFlow",
		func(ctx context.Context, doc *types.MovieContext) (*ai.Document, error) {
			time.Sleep(1 / 3 * time.Second) // reduce rate to rate limits on embedding model API
			content := createText(doc)
			aiDoc := ai.DocumentFromText(content, nil)
			eres, err := ai.Embed(ctx, embedder, ai.WithEmbedDocs(aiDoc))
			if err != nil {
				log.Println(err)
				return nil, err
			}

			genres := strings.Join(doc.Genres, ", ")
			actors := strings.Join(doc.Actors, ", ")

			query := `INSERT INTO movies (embedding, title, runtime_mins, genres, rating, released, actors, director, plot, poster, tconst, content)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
					ON CONFLICT (tconst) DO UPDATE
					SET embedding = EXCLUDED.embedding
					`
			dbCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			_, err = movieDB.DB.ExecContext(dbCtx, query,
				pgv.NewVector(eres.Embeddings[0].Embedding), doc.Title, doc.RuntimeMinutes, genres, doc.Rating, doc.Released, actors, doc.Director, doc.Plot, doc.Poster, doc.Tconst, content)
			if err != nil {
				return nil, err
			}
			return aiDoc, nil
		})
	return indexerFlow
}

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
				return movie.Plot
			}
			return ""
		}(),
	}

	jsonData, _ := json.Marshal(dataDict)
	stringData := string(jsonData)
	return stringData
}
