package standaloneWrappers

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	agents "github.com/movie-guru/pkg/agents"
	db "github.com/movie-guru/pkg/db"
	types "github.com/movie-guru/pkg/types"
	utils "github.com/movie-guru/pkg/utils"
)

type MovieAgent struct {
	MovieAgentDB *db.MovieAgentDB
	Flow         *genkit.Flow[*types.MovieAgentInput, *types.MovieAgentOutput, struct{}]
}

func CreateMovieAgent(ctx context.Context, model ai.Model, db *db.MovieAgentDB) (*MovieAgent, error) {
	flow, err := agents.GetMovieAgentFlow(ctx, model)
	if err != nil {
		return nil, err
	}
	return &MovieAgent{
		MovieAgentDB: db,
		Flow:         flow,
	}, nil
}

func (m *MovieAgent) Run(movieDocs []*types.MovieContext, history []*types.SimpleMessage, userPreferences *types.UserProfile) (*types.AgentResponse, error) {
	input := &types.MovieAgentInput{
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
