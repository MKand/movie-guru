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
	pgv "github.com/pgvector/pgvector-go"
)

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

type RetrieverFlowInput struct {
	Query string `json:"query"`
}
type RetrieverFlowOutput struct {
	Documents []*ai.Document `json:"documents"`
}

func ParseMovieContexts(docs []*ai.Document) ([]*MovieContext, error) {
	movies := make([]*MovieContext, 0, len(docs))

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
		movies = append(movies, &MovieContext{
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
	Movies []*MovieContext `json:"movies"`
}

type MovieRetriever struct {
	DB              *sql.DB
	RetrieverLength int
	Retriever       ai.Retriever
}

func (m *MovieRetriever) RetriveDocuments(ctx context.Context, query string) ([]*MovieContext, error) {
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

func GetRetrieverFlow(ctx context.Context, ret ai.Retriever) *genkit.Flow[*RetrieverFlowInput, *RetrieverFlowOutput, struct{}] {
	retFlow := genkit.DefineFlow("movieDocFlow",
		func(ctx context.Context, input *RetrieverFlowInput) (*RetrieverFlowOutput, error) {
			doc := ai.DocumentFromText(input.Query, nil)
			query := &ai.RetrieverRequest{
				Document: doc,
				Options:  10,
			}
			retOutput := make([]*ai.Document, 0, 10)
			retFlowOutput := &RetrieverFlowOutput{
				Documents: retOutput,
			}
			resp, err := ret.Retrieve(ctx, query)
			if err != nil {
				if blockedErr, ok := err.(*genai.BlockedError); ok {
					fmt.Println("Request was blocked:", blockedErr)
					return retFlowOutput, nil
				} else {
					return nil, err
				}
			}
			t := resp.Documents
			retFlowOutput.Documents = t
			return retFlowOutput, nil
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
