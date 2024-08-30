package standaloneWrappers

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	_ "github.com/lib/pq"
	db "github.com/movie-guru/pkg/db"
	flows "github.com/movie-guru/pkg/flows"
	types "github.com/movie-guru/pkg/types"
)

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
	return flows.ParseMovieContexts(rResp)
}
