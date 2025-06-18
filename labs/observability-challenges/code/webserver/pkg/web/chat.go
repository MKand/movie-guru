// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	m "github.com/movie-guru/pkg/metrics"
	"github.com/movie-guru/pkg/types"
	"go.opentelemetry.io/otel/attribute"
	metric "go.opentelemetry.io/otel/metric"
)

func chatSingleFlow(ctx context.Context, deps *Dependencies, metadata *types.Metadata, h *types.ChatHistory, user string, userMessage string, meters *m.ChatMeters) *types.AgentResponse {
	h.AddUserMessage(userMessage)

	userProfile, err := deps.DB.GetCurrentProfile(ctx, user)
	if err != nil {
		slog.ErrorContext(ctx, "Unable to get profile info for user", err.Error(), err)
	}

	simpleHistory, err := types.ParseRecentHistory(h.GetHistory(), metadata.HistoryLength)
	if err != nil {
		return types.NewErrorAgentResponse(fmt.Sprintf("Error getting user history: %w", err))
	}

	var wg sync.WaitGroup

	userProfileChan := make(chan *types.UserProfileWrapperOutput, 1)
	errChanProfile := make(chan error, 1)

	// Launch the goroutines
	// Independant goroutine with seperate context
	go func() {
		qualityContext := context.Background()
		qualityResp, err := deps.ResponseQualityFlowClient.Run(qualityContext, simpleHistory, user)
		if qualityResp != nil {
			updateChatQualityMeters(qualityContext, meters, qualityResp)
		}
		if err != nil {
			slog.ErrorContext(qualityContext, "Error updating quality meters", err.Error(), err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		pResp, err := deps.UserProfileFlowClient.Run(ctx, h, user, userProfile)
		if err != nil {
			errChanProfile <- err
			close(errChanProfile)
			return
		}
		userProfileChan <- pResp
		close(userProfileChan)
	}()

	// This is in the main thread, not async
	agentResp := types.NewAgentResponse()
	chatResp, err := deps.ChatFlowClient.Run(simpleHistory, userProfile)
	if agentResp, shouldReturn := processFlowOutput(chatResp.ModelOutputMetadata, err, h, "chatFlow"); shouldReturn {
		updateFailureChatMeters(ctx, agentResp, meters)
		return agentResp
	}

	relevantMovies := make([]string, 0, len(chatResp.RelevantMoviesTitles))
	for _, r := range chatResp.RelevantMoviesTitles {
		relevantMovies = append(relevantMovies, r.Title)
	}
	h.AddAgentMessage(chatResp.Answer)
	agentResp.Answer = chatResp.Answer
	agentResp.RelevantMovies = relevantMovies
	agentResp.Context = chatResp.ContextDocuments
	agentResp.Result = types.SUCCESS
	agentResp.TraceId = chatResp.TraceId
	agentResp.SpanId = chatResp.SpanId
	// If the user made a bad query, update it
	if chatResp.BadQuery {
		agentResp.Result = types.BAD_QUERY
		agentResp.Answer = "I cannot answer that question. Please ask me about movies or movie related information."
	}
	// Wait for goroutines to complete
	wg.Wait()

	select {
	case userProfileOutput := <-userProfileChan:
		agentResp.Preferences = userProfileOutput.UserProfile
		// Finished processing
	case err := <-errChanProfile:
		slog.ErrorContext(ctx, "UserProfileFlowClient failed", err.Error(), err)
	}

	meters.CSuccessCounter.Add(ctx, 1)
	return agentResp
}

func processFlowOutput(metadata *types.ModelOutputMetadata, err error, h *types.ChatHistory, caller string) (*types.AgentResponse, bool) {
	if err != nil {
		h.RemoveLastMessage()
		slog.ErrorContext(context.Background(), fmt.Sprintf("Error from Genkit Server: %s", caller), err.Error(), err)
		return types.NewErrorAgentResponse(err.Error()), true
	}
	if metadata != nil && metadata.SafetyIssue {
		h.RemoveLastMessage()
		return types.NewSafetyIssueAgentResponse(), true
	}
	if metadata != nil && metadata.QuotaIssue {
		h.RemoveLastMessage()
		return types.NewQuotaIssueAgentResponse(), true
	}
	return types.NewAgentResponse(), false
}

func updateChatQualityMeters(ctx context.Context, meters *m.ChatMeters, respQuality *types.ResponseQualityOutput) {
	switch strings.ToUpper(string(respQuality.UserSentiment)) {
	case strings.ToUpper(string(types.SentimentPositive)):
		meters.CSentimentCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Sentiment", "Positive")))
	case strings.ToUpper(string(types.SentimentNegative)):
		meters.CSentimentCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Sentiment", "Negative")))
	case strings.ToUpper(string(types.SentimentNeutral)):
		meters.CSentimentCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Sentiment", "Neutral")))
	default:
		meters.CSentimentCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Sentiment", "Unclassified")))
	}
	switch strings.ToUpper(string(respQuality.Outcome)) {
	case strings.ToUpper(string(types.OutcomeAcknowledged)):
		meters.COutcomeCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Outcome", "Acknowledged")))
	case strings.ToUpper(string(types.OutcomeEngaged)):
		meters.COutcomeCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Outcome", "Engaged")))
	case strings.ToUpper(string(types.OutcomeIrrelevant)):
		meters.COutcomeCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Outcome", "Irrelevant")))
	case strings.ToUpper(string(types.OutcomeRejected)):
		meters.COutcomeCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Outcome", "Rejected")))
	default:
		meters.COutcomeCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Outcome", "Unclassified")))
	}
}

func updateFailureChatMeters(ctx context.Context, agentResp *types.AgentResponse, meters *m.ChatMeters) {
	if agentResp.Result == types.UNSAFE {
		meters.CSafetyIssueCounter.Add(ctx, 1)
		slog.Error("Flows returned error: User gave UNSAFE input")
	}
	if agentResp.Result == types.QUOTALIMIT {
		meters.CQuotaLimitCounter.Add(ctx, 1)
	}
	if agentResp.Result == types.BAD_QUERY {
		meters.CWrongQueryCounter.Add(ctx, 1)
	}
}
