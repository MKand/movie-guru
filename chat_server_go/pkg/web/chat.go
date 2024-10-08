package web

import (
	"context"
	"fmt"

	"github.com/movie-guru/pkg/db"
	"github.com/movie-guru/pkg/types"
)

func chat(ctx context.Context, deps *Dependencies, metadata *db.Metadata, h *types.ChatHistory, user string, userMessage string) (*types.AgentResponse, *types.ResponseQualityOutput) {
	h.AddUserMessage(userMessage)
	simpleHistory, err := types.ParseRecentHistory(h.GetHistory(), metadata.HistoryLength)
	respQuality := &types.ResponseQualityOutput{
		Outcome:       types.OutcomeUnknown,
		UserSentiment: types.SentimentUnknown,
	}
	respQualityChan := make(chan *types.ResponseQualityOutput)
	errChan := make(chan error)

	// Launch the goroutine
	go func() {
		pResp, err := deps.ResponseQualityFlowClient.Run(ctx, simpleHistory, user)
		if err != nil {
			errChan <- err
		} else {
			respQualityChan <- pResp
		}
	}()

	pResp, err := deps.UserProfileFlowClient.Run(ctx, h, user)
	if agentResp, shouldReturn := processFlowOutput(pResp.ModelOutputMetadata, err, h); shouldReturn {
		return agentResp, respQuality
	}

	qResp, err := deps.QueryTransformFlowClient.Run(simpleHistory, pResp.UserProfile)
	if agentResp, shouldReturn := processFlowOutput(qResp.ModelOutputMetadata, err, h); shouldReturn {
		return agentResp, respQuality
	}

	movieContext := []*types.MovieContext{}
	if qResp.Intent == types.USERINTENT(types.REQUEST) || qResp.Intent == types.USERINTENT(types.RESPONSE) {
		movieContext, err = deps.MovieRetrieverFlowClient.RetriveDocuments(ctx, qResp.TransformedQuery)
		if agentResp, shouldReturn := processFlowOutput(nil, err, h); shouldReturn {
			return agentResp, respQuality
		}
	}

	mAgentResp, err := deps.MovieFlowClient.Run(movieContext, simpleHistory, pResp.UserProfile)
	if agentResp, shouldReturn := processFlowOutput(nil, err, h); shouldReturn {
		return agentResp, respQuality
	}
	h.AddAgentMessage(mAgentResp.Answer)

	select {
	case respQuality = <-respQualityChan:
		fmt.Println("Response Quality: ", respQuality)

	case err := <-errChan:
		fmt.Println("Response Quality Error:", err)
	}

	return mAgentResp, respQuality
}

func processFlowOutput(metadata *types.ModelOutputMetadata, err error, h *types.ChatHistory) (*types.AgentResponse, bool) {
	if err != nil {
		h.AddAgentErrorMessage()
		return types.NewErrorAgentResponse(err.Error()), true
	}
	if metadata != nil && metadata.SafetyIssue {
		h.AddSafetyIssueErrorMessage()
		return types.NewSafetyIssueAgentResponse(), true
	}
	return types.NewAgentResponse(), false
}
