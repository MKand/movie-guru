package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	metrics "github.com/movie-guru/pkg/metrics"
	"github.com/movie-guru/pkg/types"
)

func createChatHandler(deps *Dependencies, meters *metrics.ChatMeters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errLogPrefix := "Error: ChatHandler: "
		var err error
		ctx := r.Context()
		origin := r.Header.Get("Origin")
		addResponseHeaders(w, origin)
		sessionInfo := &SessionInfo{}
		if r.Method != "OPTIONS" {
			sessionInfo, err = getSessionInfo(ctx, r)
			if err != nil {
				if err, ok := err.(*AuthorizationError); ok {
					log.Println(errLogPrefix, "Unauthorized")
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				log.Println(errLogPrefix, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !sessionInfo.Authenticated {
				log.Println(errLogPrefix, "Unauthenticated")
				http.Error(w, "Forbidden", http.StatusForbidden)
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
				log.Println(errLogPrefix, err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if len(chatRequest.Content) > metadata.MaxUserMessageLen {
				log.Println(errLogPrefix, "Message too long")
				http.Error(w, "Message too long", http.StatusBadRequest)
				return
			}
			ch, err := getHistory(ctx, user)
			if err != nil {
				log.Println(errLogPrefix, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			agentResp, respQuality := chat(ctx, deps, metadata, ch, user, chatRequest.Content)
			updateChatMeters(ctx, agentResp, meters, respQuality)

			saveHistory(ctx, ch, user)
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

func updateChatMeters(ctx context.Context, agentResp *types.AgentResponse, meters *metrics.ChatMeters, respQuality *types.ResponseQualityOutput) {
	if agentResp.Result == types.UNSAFE {
		meters.CSafetyIssueCounter.Add(ctx, 1)
	}
	if agentResp.Result == types.SUCCESS {
		meters.CSuccessCounter.Add(ctx, 1)
	}
	switch strings.ToUpper(string(respQuality.UserSentiment)) {
	case strings.ToUpper(string(types.SentimentPositive)):
		meters.CSentimentPositiveCounter.Add(ctx, 1)
	case strings.ToUpper(string(types.SentimentNegative)):
		meters.CSentimentNegativeCounter.Add(ctx, 1)
	case strings.ToUpper(string(types.SentimentNeutral)):
		meters.CSentimentNeutralCounter.Add(ctx, 1)
	default:
		meters.CSentimentUnclassifiedCounter.Add(ctx, 1)
		fmt.Println("unclassified sentiment added")
	}
	switch strings.ToUpper(string(respQuality.Outcome)) {
	case strings.ToUpper(string(types.OutcomeAcknowledged)):
		meters.COutcomeAcknowledgedCounter.Add(ctx, 1)
	case strings.ToUpper(string(types.OutcomeEngaged)):
		meters.COutcomeEngagedCounter.Add(ctx, 1)
	case strings.ToUpper(string(types.OutcomeIrrelevant)):
		meters.COutcomeIrrelevantCounter.Add(ctx, 1)
	case strings.ToUpper(string(types.OutcomeRejected)):
		meters.COutcomeRejectedCounter.Add(ctx, 1)
	default:
		meters.COutcomeUnclassifiedCounter.Add(ctx, 1)
		fmt.Println("unclassified outcome added")
	}
}
