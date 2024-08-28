package standaloneWrappers

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	db "github.com/movie-guru/pkg/db"
	flows "github.com/movie-guru/pkg/flows"
	types "github.com/movie-guru/pkg/types"
	utils "github.com/movie-guru/pkg/utils"
)

type MovieFlow struct {
	MovieDB *db.MovieDB
	Flow    *genkit.Flow[*types.MovieFlowInput, *types.MovieFlowOutput, struct{}]
}

func CreateMovieFlow(ctx context.Context, model ai.Model, db *db.MovieDB) (*MovieFlow, error) {
	flow, err := flows.GetMovieFlow(ctx, model)
	if err != nil {
		return nil, err
	}
	return &MovieFlow{
		MovieDB: db,
		Flow:    flow,
	}, nil
}

func (m *MovieFlow) Run(movieDocs []*types.MovieContext, history []*types.SimpleMessage, userPreferences *types.UserProfile) (*types.AgentResponse, error) {
	input := &types.MovieFlowInput{
		History:          history,
		UserPreferences:  userPreferences,
		ContextDocuments: movieDocs,
		UserMessage:      history[len(history)-1].Content,
	}
	resp, err := m.Flow.Run(context.Background(), input)
	if err != nil {
		return nil, err
	}

	relevantMovies := make([]string, 0, len(resp.RelevantMoviesTitles))
	for _, r := range resp.RelevantMoviesTitles {
		relevantMovies = append(relevantMovies, r.Title)
	}

	agentResponse := &types.AgentResponse{
		Answer:         resp.Answer,
		RelevantMovies: relevantMovies,
		Context:        utils.FilterRelevantContext(relevantMovies, movieDocs),
		ErrorMessage:   "",
		Result:         types.SUCCESS,
		Preferences:    userPreferences,
	}
	return agentResponse, nil
}
