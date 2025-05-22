package webmock

import (
	"context"
	"math/rand"
	"net/http"
	"strings"

	m "github.com/movie-guru/pkg/metrics"
	"github.com/movie-guru/pkg/types"
	"go.opentelemetry.io/otel/attribute"
	metric "go.opentelemetry.io/otel/metric"
)

func createChatHandler(deps *Dependencies, meters *m.ChatMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method == "POST" {
			meters.CCounter.Add(ctx, 1)

			success := PickSuccess(deps.CurrentProbMetrics.ChatSuccess)
			latency := PickLatencyValue(deps.CurrentProbMetrics.ChatLatencyMinMS, deps.CurrentProbMetrics.ChatLatencyMaxMS)

			defer func() {
				meters.CLatencyHistogram.Record(ctx, int64(latency))
			}()

			if success {
				meters.CSuccessCounter.Add(ctx, 1)
				agentResp := types.NewAgentResponse()
				respQuality := &types.ResponseQualityOutput{
					Outcome:       types.OutcomeUnknown,
					UserSentiment: types.SentimentUnknown,
				}

				if PickSuccess(deps.CurrentProbMetrics.ChatSafetyIssue) {
					agentResp = types.NewSafetyIssueAgentResponse()
				}
				sentimentProbabilities := []float32{deps.CurrentProbMetrics.ChatSPositive, deps.CurrentProbMetrics.ChatSNegative, deps.CurrentProbMetrics.ChatSNeutral, deps.CurrentProbMetrics.ChatSUnclassified}
				engagementProbabilities := []float32{deps.CurrentProbMetrics.ChatEngaged, deps.CurrentProbMetrics.ChatAcknowledged, deps.CurrentProbMetrics.ChatRejected, deps.CurrentProbMetrics.ChatUnclassified}
				chatSentiment := pickChatQuality(sentimentProbabilities)
				switch chatSentiment {
				case 0:
					respQuality.UserSentiment = types.SentimentPositive
				case 1:
					respQuality.UserSentiment = types.SentimentNegative
				case 2:
					respQuality.UserSentiment = types.SentimentNeutral
				case 3:
					respQuality.UserSentiment = types.SentimentUnknown
				}
				chatEngagement := pickChatQuality(engagementProbabilities)
				switch chatEngagement {
				case 0:
					respQuality.Outcome = types.OutcomeEngaged
				case 1:
					respQuality.Outcome = types.OutcomeAcknowledged
				case 2:
					respQuality.Outcome = types.OutcomeRejected
				case 3:
					respQuality.Outcome = types.OutcomeUnknown
				}

				updateChatMeters(ctx, agentResp, meters, respQuality)
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		if r.Method == "OPTIONS" {
			return
		}
	}
}

func updateChatMeters(ctx context.Context, agentResp *types.AgentResponse, meters *m.ChatMeters, respQuality *types.ResponseQualityOutput) {
	if agentResp.Result == types.UNSAFE {
		meters.CSafetyIssueCounter.Add(ctx, 1)
	}
	if agentResp.Result == types.SUCCESS {
		meters.CSuccessCounter.Add(ctx, 1)
	}
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
	case strings.ToUpper(string(types.OutcomeRejected)):
		meters.COutcomeCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Outcome", "Rejected")))
	default:
		meters.COutcomeCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Outcome", "Unclassified")))
	}
}

func pickChatQuality(probabilities []float32) int {

	// Generate a random number between 0 and 1
	randomNumber := rand.Float32()

	// Cumulative probability
	cumulativeProbability := float32(0.0)

	// Iterate through the probabilities and check if the random number falls within the range
	for i, probability := range probabilities {
		cumulativeProbability += probability
		if randomNumber <= cumulativeProbability {
			switch i {
			case 0:
				return 0
			case 1:
				return 1
			case 2:
				return 2
			case 3:
				return 3
			}
		}
	}
	// This should not happen, but return a default value just in case
	return 3
}
