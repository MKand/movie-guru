package standaloneWrappers

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	agents "github.com/movie-guru/pkg/agents"
	db "github.com/movie-guru/pkg/db"
	types "github.com/movie-guru/pkg/types"
)

type QueryTransformAgent struct {
	MovieAgentDB *db.MovieAgentDB
	Flow         *genkit.Flow[*types.QueryTransformInput, *types.QueryTransformOutput, struct{}]
}

func CreateQueryTransformAgent(ctx context.Context, model ai.Model, db *db.MovieAgentDB) (*QueryTransformAgent, error) {
	flow, err := agents.GetQueryTransformFlow(ctx, model)
	if err != nil {
		return nil, err
	}
	return &QueryTransformAgent{
		MovieAgentDB: db,
		Flow:         flow,
	}, nil
}

func (q *QueryTransformAgent) Run(ctx context.Context, history []*types.SimpleMessage, preferences *types.UserProfile) (*types.QueryTransformOutput, error) {
	queryTransformInput := types.QueryTransformInput{Profile: preferences, History: history, UserMessage: history[len(history)-1].Content}
	resp, err := q.Flow.Run(ctx, &queryTransformInput)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
