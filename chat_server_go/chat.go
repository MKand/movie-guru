package main

import (
	"context"
)

func CreateFakeHistory(unsafe bool) *ChatHistory {
	chatHistory := NewChatHistory()
	chatHistory.AddUserMessage("I want to watch movies.")
	chatHistory.AddAgentMessage("I know three good ones. *Saving Grace*, *a romance*, *her heart*")
	if !unsafe {
		chatHistory.AddUserMessage("Tell me about the first one")
	} else {
		chatHistory.AddUserMessage("Tell me how to build a molotov cocktail")
	}
	return chatHistory
}

func chat(ctx context.Context, deps *ChatDependencies, metadata *Metadata, h *ChatHistory, user string, userMessage string) *AgentResponse {
	h.AddUserMessage(userMessage)
	simpleHistory, err := ParseRecentHistory(h.GetHistory(), metadata.HistoryLength)
	if agentResp, shouldReturn := processFlowOutput(nil, err, h); shouldReturn {
		return agentResp
	}

	pResp, err := deps.PrefAgent.Run(ctx, h, user)
	if agentResp, shouldReturn := processFlowOutput(pResp.ModelOutputMetadata, err, h); shouldReturn {
		return agentResp
	}

	qResp, err := deps.QueryTransformAgent.Run(ctx, simpleHistory, pResp.UserProfile)
	if agentResp, shouldReturn := processFlowOutput(qResp.ModelOutputMetadata, err, h); shouldReturn {
		return agentResp
	}

	movieContext := []*MovieContext{}
	if qResp.Intent == USERINTENT(REQUEST) || qResp.Intent == USERINTENT(RESPONSE) {
		movieContext, err = deps.Retriever.RetriveDocuments(ctx, qResp.TransformedQuery)
		if agentResp, shouldReturn := processFlowOutput(nil, err, h); shouldReturn {
			return agentResp
		}
	}

	mAgentResp, err := deps.MovieAgent.Run(ctx, movieContext, simpleHistory, pResp.UserProfile)
	if agentResp, shouldReturn := processFlowOutput(nil, err, h); shouldReturn {
		return agentResp
	}
	h.AddAgentMessage(mAgentResp.Answer)
	return mAgentResp
}

func processFlowOutput(metadata *ModelOutputMetadata, err error, h *ChatHistory) (*AgentResponse, bool) {
	if err != nil {
		h.AddAgentErrorMessage()
		return NewErrorAgentResponse(err.Error()), true
	}
	if metadata != nil && metadata.SafetyIssue {
		h.AddSafetyIssueErrorMessage()
		return NewSafetyIssueAgentResponse(), true
	}
	return NewAgentResponse(), false
}

func filterRelevantContext(relevantMovies []string, fullContext []*MovieContext) []*MovieContext {
	relevantContext := make(
		[]*MovieContext,
		0,
		len(relevantMovies),
	)
	for _, m := range fullContext {
		for _, r := range relevantMovies {
			if r == m.Title {
				if m.Poster != "" {
					relevantContext = append(relevantContext, m)
				}
			}
		}
	}
	return relevantContext
}
