package main

import (
	"context"
	"database/sql"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/plugins/vertexai"
	_ "github.com/lib/pq"
	_ "github.com/movie-guru/pkg/types"
	pgv "github.com/pgvector/pgvector-go"
)

func GetEmbedder(embeddingModelName string) *ai.Embedder {
	embedder := vertexai.Embedder(embeddingModelName)
	return &embedder
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
