package agents

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/firebase/genkit/go/ai"
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
			Title       string  `json:"title"`
			RuntimeMins int     `json:"runtime_mins"`
			Genres      string  `json:"genres"`
			Rating      float32 `json:"rating"`
			Released    float64 `json:"released"`
			Actors      string  `json:"actors"`
			Director    string  `json:"director"`
			Plot        string  `json:"plot"`
			Poster      string  `json:"poster"`
		}

		err := json.Unmarshal([]byte(doc.Content[0].Text), &intermediate)
		if err != nil {
			return nil, err
		}

		movies = append(movies, &types.MovieContext{
			Title:          intermediate.Title,
			RuntimeMinutes: intermediate.RuntimeMins,
			Genres:         strings.Split(intermediate.Genres, ", "),
			Rating:         intermediate.Rating,
			Plot:           intermediate.Plot,
			Released:       int(intermediate.Released),
			Director:       intermediate.Director,
			Actors:         strings.Split(intermediate.Actors, ", "),
			Poster:         intermediate.Poster,
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
	ret := defineRetriever(maxRetLength, db, embedder)
	return &MovieRetriever{
		DB:              db,
		RetrieverLength: maxRetLength,
		Retriever:       ret,
	}
}

func defineRetriever(maxRetLength int, db *sql.DB, embedder ai.Embedder) ai.Retriever {
	f := func(ctx context.Context, req *ai.RetrieverRequest) (*ai.RetrieverResponse, error) {
		eres, err := ai.Embed(ctx, embedder, ai.WithEmbedDocs(req.Document))
		if err != nil {
			return nil, err
		}

		rows, err := db.QueryContext(ctx, `
					SELECT title, poster, content
					FROM fake_movies_table
					ORDER BY embedding <-> $1
					LIMIT $2`,
			pgv.NewVector(eres.Embeddings[0].Embedding), maxRetLength)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		res := &ai.RetrieverResponse{}
		for rows.Next() {
			var title, poster, content string
			if err := rows.Scan(&title, &poster, &content); err != nil {
				return nil, err
			}
			meta := map[string]any{
				"title":  title,
				"poster": poster,
			}
			doc := &ai.Document{
				Content:  []*ai.Part{ai.NewTextPart(content)},
				Metadata: meta,
			}
			res.Documents = append(res.Documents, doc)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return res, nil
	}
	return ai.DefineRetriever("pgvector", "movies", f)
}
