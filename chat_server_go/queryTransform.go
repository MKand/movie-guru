package main

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type QueryTransformAgent struct {
	Model ai.Model
	Flow  *genkit.Flow[*QueryTransformInput, *QueryTransformOutput, struct{}]
}

func CreateQueryTransformAgent(ctx context.Context, model ai.Model) (*QueryTransformAgent, error) {
	flow, err := GetQueryTransformFlow(ctx, model)
	if err != nil {
		return nil, err
	}
	return &QueryTransformAgent{
		Model: model,
		Flow:  flow,
	}, nil
}

func (q *QueryTransformAgent) Run(ctx context.Context, history []*SimpleMessage, preferences *UserProfile) (*QueryTransformOutput, error) {
	queryTransformInput := QueryTransformInput{Profile: preferences, History: history, UserMessage: history[len(history)-1].Content}
	resp, err := q.Flow.Run(ctx, &queryTransformInput)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
