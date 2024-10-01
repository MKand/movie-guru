package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"cloud.google.com/go/vertexai/genai"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/vertexai"
	_ "github.com/lib/pq"
	_ "github.com/movie-guru/pkg/types"
	types "github.com/movie-guru/pkg/types"
	pgv "github.com/pgvector/pgvector-go"
)

func ParseMovieContexts(docs []*ai.Document) ([]*types.MovieContext, error) {
	movies := make([]*types.MovieContext, 0, len(docs))

	for _, doc := range docs {
		var intermediate struct {
			Title    string `json:"title"`
			Genres   string `json:"genres"`
			Actors   string `json:"actors"`
			Director string `json:"director"`
			Plot     string `json:"plot"`
			Poster   string `json:"poster, omitempty`
		}

		err := json.Unmarshal([]byte(doc.Content[0].Text), &intermediate)
		if err != nil {
			return nil, err
		}

		rating, _ := doc.Metadata["rating"].(float32)
		runTimeMins, _ := doc.Metadata["runtime_minutes"].(int)
		released, _ := doc.Metadata["releases"].(int)
		poster := doc.Metadata["poster"].(string)
		movies = append(movies, &types.MovieContext{
			Title:          intermediate.Title,
			RuntimeMinutes: runTimeMins,
			Genres:         strings.Split(intermediate.Genres, ", "),
			Rating:         rating,
			Plot:           intermediate.Plot,
			Released:       released,
			Director:       intermediate.Director,
			Actors:         strings.Split(intermediate.Actors, ", "),
			Poster:         poster,
		})
	}

	return movies, nil
}

type MovieContextList struct {
	Movies []*types.MovieContext `json:"movies"`
}

type MovieRetriever struct {
	DB              *sql.DB
	RetrieverLength int
	Retriever       ai.Retriever
}

func (m *MovieRetriever) RetriveDocuments(ctx context.Context, query string) ([]*types.MovieContext, error) {
	doc := ai.DocumentFromText(query, nil)
	retDoc := ai.RetrieverRequest{
		Document: doc,
		Options:  m.RetrieverLength,
	}
	rResp, err := m.Retriever.Retrieve(ctx, &retDoc)
	if err != nil {
		return nil, err
	}
	return ParseMovieContexts(rResp.Documents)
}

func GetEmbedder(embeddingModelName string) ai.Embedder {
	embedder := vertexai.Embedder(embeddingModelName)
	return embedder
}

func CreateMovieRetriever(embeddingModelName string, maxRetLength int, db *sql.DB) *MovieRetriever {
	embedder := GetEmbedder(embeddingModelName)
	ret := DefineRetriever(maxRetLength, db, embedder)
	return &MovieRetriever{
		DB:              db,
		RetrieverLength: maxRetLength,
		Retriever:       ret,
	}
}

func GetRetrieverFlow(ctx context.Context, ret ai.Retriever) *genkit.Flow[string, []*ai.Document, struct{}] {
	retFlow := genkit.DefineFlow("movieDocFlow",
		func(ctx context.Context, query string) ([]*ai.Document, error) {
			doc := ai.DocumentFromText(query, nil)
			input := &ai.RetrieverRequest{
				Document: doc,
				Options:  10,
			}
			retOutput := make([]*ai.Document, 0, 10)
			resp, err := ret.Retrieve(ctx, input)
			if err != nil {
				if blockedErr, ok := err.(*genai.BlockedError); ok {
					fmt.Println("Request was blocked:", blockedErr)
					return retOutput, nil
				} else {
					return nil, err
				}
			}
			t := resp.Documents
			return t, nil
		})
	return retFlow
}

func DefineRetriever(maxRetLength int, db *sql.DB, embedder ai.Embedder) ai.Retriever {
	f := func(ctx context.Context, req *ai.RetrieverRequest) (*ai.RetrieverResponse, error) {
		eres, err := ai.Embed(ctx, embedder, ai.WithEmbedDocs(req.Document))
		if err != nil {
			return nil, err
		}

		rows, err := db.QueryContext(ctx, `
					SELECT title, poster, content, released, runtime_mins, rating
					FROM movies
					ORDER BY embedding <-> $1
					LIMIT $2`,
			pgv.NewVector(eres.Embeddings[0].Embedding), maxRetLength)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		retrieverResponse := &ai.RetrieverResponse{}
		for rows.Next() {
			var title, poster, content string
			var released, runtime_mins int
			var rating float32
			if err := rows.Scan(&title, &poster, &content, &released, &runtime_mins, &rating); err != nil {
				return nil, err
			}
			meta := map[string]any{
				"title":        title,
				"poster":       poster,
				"released":     released,
				"rating":       rating,
				"runtime_mins": runtime_mins,
			}
			doc := &ai.Document{
				Content:  []*ai.Part{ai.NewTextPart(content)},
				Metadata: meta,
			}
			retrieverResponse.Documents = append(retrieverResponse.Documents, doc)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return retrieverResponse, nil
	}
	return ai.DefineRetriever("pgvector", "movieRetriever", f)
}
