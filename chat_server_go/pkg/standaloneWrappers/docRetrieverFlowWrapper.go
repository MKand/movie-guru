package standaloneWrappers

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	_ "github.com/lib/pq"
	db "github.com/movie-guru/pkg/db"
	flows "github.com/movie-guru/pkg/flows"
	types "github.com/movie-guru/pkg/types"
)

func parseMovieContexts(docs []*ai.Document) ([]*types.MovieContext, error) {
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

type MovieRetrieverFlow struct {
	RetrieverLength int
	Flow            *genkit.Flow[*ai.RetrieverRequest, []*ai.Document, struct{}]
}

func CreateMovieRetrieverFlow(ctx context.Context, embeddingModelName string, maxRetLength int, db *db.MovieDB) *MovieRetrieverFlow {
	ret := flows.CreateMovieRetriever(embeddingModelName, maxRetLength, db.DB)
	flow := flows.GetRetrieverFlow(ctx, ret.Retriever)
	return &MovieRetrieverFlow{
		RetrieverLength: maxRetLength,
		Flow:            flow,
	}
}

func (r *MovieRetrieverFlow) RetriveDocuments(ctx context.Context, query string) ([]*types.MovieContext, error) {
	doc := ai.DocumentFromText(query, nil)
	retDoc := ai.RetrieverRequest{
		Document: doc,
		Options:  r.RetrieverLength,
	}
	rResp, err := r.Flow.Run(ctx, &retDoc)
	if err != nil {
		return nil, err
	}
	return parseMovieContexts(rResp)
}
