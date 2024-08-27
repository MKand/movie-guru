package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/firebase/genkit/go/ai"
	_ "github.com/lib/pq"
	_ "github.com/movie-guru/pkg/types"
)

func ParseMovieContexts(docs []*ai.Document) ([]*MovieContext, error) {
	movies := make([]*MovieContext, 0, len(docs))

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

		movies = append(movies, &MovieContext{
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
