package main

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	_ "github.com/movie-guru/pkg/types"
)

type MovieAgent struct {
	Model ai.Model
	Flow  *genkit.Flow[*MovieAgentInput, *MovieAgentOutput, struct{}]
}

func CreateMovieAgent(ctx context.Context, model ai.Model) (*MovieAgent, error) {
	flow, err := GetMovieAgentFlow(ctx, model)
	if err != nil {
		return nil, err
	}
	return &MovieAgent{
		Model: model,
		Flow:  flow,
	}, nil
}

func (m *MovieAgent) Run(ctx context.Context, movieDocs []*MovieContext, history []*SimpleMessage, userPreferences *UserProfile) (*AgentResponse, error) {

	input := &MovieAgentInput{
		History:          history,
		UserPreferences:  userPreferences,
		ContextDocuments: movieDocs,
		UserMessage:      history[len(history)-1].Content,
	}
	resp, err := m.Flow.Run(ctx, input)
	if err != nil {
		return nil, err
	}

	relevantMovies := make([]string, 0, len(resp.RelevantMoviesTitles))
	for _, r := range resp.RelevantMoviesTitles {
		relevantMovies = append(relevantMovies, r.Title)
	}

	agentResponse := &AgentResponse{
		Answer:         resp.Answer,
		RelevantMovies: relevantMovies,
		Context:        filterRelevantContext(relevantMovies, movieDocs),
		ErrorMessage:   "",
		Result:         SUCCESS,
		Preferences:    userPreferences,
	}
	return agentResponse, nil
}
