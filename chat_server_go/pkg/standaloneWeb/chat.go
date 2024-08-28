package standaloneWeb

import (
	"context"

	db "github.com/movie-guru/pkg/db"
	types "github.com/movie-guru/pkg/types"
)

func chat(ctx context.Context, deps *Dependencies, metadata *db.Metadata, h *types.ChatHistory, user string, userMessage string) *types.AgentResponse {
	h.AddUserMessage(userMessage)
	simpleHistory, err := types.ParseRecentHistory(h.GetHistory(), metadata.HistoryLength)
	if agentResp, shouldReturn := processFlowOutput(nil, err, h); shouldReturn {
		return agentResp
	}

	pResp, err := deps.UserProfileFlow.Run(ctx, h, user)
	if agentResp, shouldReturn := processFlowOutput(pResp.ModelOutputMetadata, err, h); shouldReturn {
		return agentResp
	}

	qResp, err := deps.QueryTransformFlow.Run(ctx, simpleHistory, pResp.UserProfile)
	if agentResp, shouldReturn := processFlowOutput(qResp.ModelOutputMetadata, err, h); shouldReturn {
		return agentResp
	}

	movieContext := []*types.MovieContext{}
	if qResp.Intent == types.USERINTENT(types.REQUEST) || qResp.Intent == types.USERINTENT(types.RESPONSE) {
		movieContext, err = deps.MovieRetrieverFlow.RetriveDocuments(ctx, qResp.TransformedQuery)
		if agentResp, shouldReturn := processFlowOutput(nil, err, h); shouldReturn {
			return agentResp
		}
	}

	mAgentResp, err := deps.MovieFlow.Run(movieContext, simpleHistory, pResp.UserProfile)
	if agentResp, shouldReturn := processFlowOutput(nil, err, h); shouldReturn {
		return agentResp
	}
	h.AddAgentMessage(mAgentResp.Answer)
	return mAgentResp
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
