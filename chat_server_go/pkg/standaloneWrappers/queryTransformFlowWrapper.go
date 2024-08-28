package standaloneWrappers

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	db "github.com/movie-guru/pkg/db"
	flows "github.com/movie-guru/pkg/flows"
	types "github.com/movie-guru/pkg/types"
)

type QueryTransformFlow struct {
	MovieDB *db.MovieDB
	Flow    *genkit.Flow[*types.QueryTransformFlowInput, *types.QueryTransformFlowOutput, struct{}]
}

func CreateQueryTransformFlow(ctx context.Context, model ai.Model, db *db.MovieDB) (*QueryTransformFlow, error) {
	flow, err := flows.GetQueryTransformFlow(ctx, model)
	if err != nil {
		return nil, err
	}
	return &QueryTransformFlow{
		MovieDB: db,
		Flow:    flow,
	}, nil
}

func (q *QueryTransformFlow) Run(ctx context.Context, history []*types.SimpleMessage, preferences *types.UserProfile) (*types.QueryTransformFlowOutput, error) {
	queryTransformInput := types.QueryTransformFlowInput{Profile: preferences, History: history, UserMessage: history[len(history)-1].Content}
	resp, err := q.Flow.Run(ctx, &queryTransformInput)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
