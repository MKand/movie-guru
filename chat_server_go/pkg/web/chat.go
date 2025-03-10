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

	"github.com/movie-guru/pkg/db"
	m "github.com/movie-guru/pkg/metrics"
	"github.com/movie-guru/pkg/types"
	"go.opentelemetry.io/otel/attribute"
	metric "go.opentelemetry.io/otel/metric"
)

func chatSingleFlow(ctx context.Context, deps *Dependencies, metadata *db.Metadata, h *types.ChatHistory, user string, userMessage string, meters *m.ChatMeters) *types.AgentResponse {
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

	userProfileChan := make(chan *types.UserProfileOutput, 1)
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

	// Process anamolous responses such as quota issues, bad queries, safety issues etc.
	if agentResp, shouldReturn := processFlowOutput(chatResp.ModelOutputMetadata, err, h, "chatFlow"); shouldReturn {
		agentResp.TraceId = chatResp.TraceId
		agentResp.SpanId = chatResp.SpanId
		updateAnamolyChatMeters(ctx, agentResp, meters)
		if agentResp.Result == types.UNSAFE { // Unsafe is still a response indicating a successful flow
			updateSuccessMeter(ctx, meters, agentResp.Result, agentResp.TraceId)
		} else {
			slog.InfoContext(ctx, fmt.Sprintf("NOT updating success meter, result: %s, traceId: %s", agentResp.Result, agentResp.TraceId))
		}
		return agentResp
	}
	if chatResp.WrongQuery {
		agentResp.Result = types.BAD_QUERY
		updateSuccessMeter(ctx, meters, agentResp.Result, agentResp.TraceId)
		updateAnamolyChatMeters(ctx, agentResp, meters)
		h.RemoveLastMessage()
		agentResp.Answer = "I cannot answer that question. Please ask me about movies or movie related information."
		return agentResp
	}

	// Response is successful
	relevantMovies := make([]string, 0, len(chatResp.RelevantMoviesTitles))
	for _, r := range chatResp.RelevantMoviesTitles {
		relevantMovies = append(relevantMovies, r.Title)
	}
	agentResp.Answer = chatResp.Answer
	agentResp.RelevantMovies = relevantMovies
	agentResp.Context = chatResp.ContextDocuments
	agentResp.Result = types.SUCCESS
	agentResp.TraceId = chatResp.TraceId
	agentResp.SpanId = chatResp.SpanId
	h.AddAgentMessage(chatResp.Answer)

	// Wait for goroutines to complete
	wg.Wait()

	select {
	case userProfileOutput := <-userProfileChan:
		agentResp.Preferences = userProfileOutput.UserProfile
		// Finished processing
	case err := <-errChanProfile:
		slog.ErrorContext(ctx, "UserProfileFlowClient failed", err.Error(), err)
	}
	updateSuccessMeter(ctx, meters, agentResp.Result, agentResp.TraceId)
	return agentResp
}

func updateSuccessMeter(ctx context.Context, meters *m.ChatMeters, result types.RESULT, traceId string) {
	slog.InfoContext(ctx, fmt.Sprintf("Updating success meter, result: %s, traceId: %s", result, traceId))
	meters.CSuccessCounter.Add(ctx, 1) // The agent behaved as expected.

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

func updateAnamolyChatMeters(ctx context.Context, agentResp *types.AgentResponse, meters *m.ChatMeters) {

	if agentResp.Result == types.UNSAFE {
		slog.InfoContext(ctx, "Updating UNSAFE counter")
		meters.CSafetyIssueCounter.Add(ctx, 1)
	}
	if agentResp.Result == types.QUOTALIMIT {
		slog.InfoContext(ctx, "Updating QUOTALIMIT counter")
		meters.CQuotaLimitCounter.Add(ctx, 1)
	}
	if agentResp.Result == types.BAD_QUERY {
		slog.InfoContext(ctx, "Updating BADQUERY counter")
		meters.CWrongQueryCounter.Add(ctx, 1)
	}
}
