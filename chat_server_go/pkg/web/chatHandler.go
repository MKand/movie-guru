package web

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	db "github.com/movie-guru/pkg/db"

	m "github.com/movie-guru/pkg/metrics"
	"github.com/movie-guru/pkg/types"
	"go.opentelemetry.io/otel/attribute"
	metric "go.opentelemetry.io/otel/metric"
	"golang.org/x/exp/slog"
)

func createChatHandler(deps *Dependencies, meters *m.ChatMeters, metadata *db.Metadata) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := r.Context()
		origin := r.Header.Get("Origin")
		addResponseHeaders(w, origin)
		sessionInfo := &SessionInfo{}
		if r.Method != "OPTIONS" {
			var shouldReturn bool
			sessionInfo, shouldReturn = authenticateAndGetSessionInfo(ctx, sessionInfo, err, r, w, metadata)
			if shouldReturn {
				return
			}
		}
		if r.Method == "POST" {
			meters.CCounter.Add(ctx, 1)
			startTime := time.Now()
			defer func() {
				meters.CLatencyHistogram.Record(ctx, int64(time.Since(startTime).Milliseconds()))
			}()
			addResponseHeaders(w, origin)
			user := sessionInfo.User
			chatRequest := &ChatRequest{
				Content: "",
			}
			err := json.NewDecoder(r.Body).Decode(chatRequest)
			if err != nil {
				slog.InfoContext(ctx, "Error while decoding request", slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if len(chatRequest.Content) > metadata.MaxUserMessageLen {
				slog.InfoContext(ctx, "Input message too long", slog.String("user", user), slog.Any("error", err.Error()))
				http.Error(w, "Message too long", http.StatusBadRequest)
				return
			}
			ch, err := getHistory(ctx, user)
			if err != nil {
				slog.ErrorContext(ctx, "Error while fetching history", slog.String("user", user), slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			agentResp, respQuality := chat(ctx, deps, metadata, ch, user, chatRequest.Content)
			updateChatMeters(ctx, agentResp, meters, respQuality)

			saveHistory(ctx, ch, user, metadata)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(agentResp)
			return

		}
		if r.Method == "OPTIONS" {
			addResponseHeaders(w, origin)
			handleOptions(w, origin)
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
	case strings.ToUpper(string(types.OutcomeIrrelevant)):
		meters.COutcomeCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Outcome", "Irrelevant")))
	case strings.ToUpper(string(types.OutcomeRejected)):
		meters.COutcomeCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Outcome", "Rejected")))
	default:
		meters.COutcomeCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("Outcome", "Unclassified")))
	}
}
