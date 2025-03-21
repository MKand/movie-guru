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
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/movie-guru/pkg/db"
	"github.com/movie-guru/pkg/types"

	m "github.com/movie-guru/pkg/metrics"
)

func createChatHandler(deps *Dependencies, meters *m.ChatMeters, metadata *db.Metadata) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := r.Context()
		sessionInfo := &SessionInfo{}
		if r.Method != "OPTIONS" {
			var shouldReturn bool
			sessionInfo, shouldReturn = authenticateAndGetSessionInfo(ctx, sessionInfo, err, r, w)
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
				slog.InfoContext(ctx, "Input message too long")
				agentResp := types.NewAgentResponse()
				agentResp.Result = types.TOO_LONG
				agentResp.Answer = "Sorry, I cannot process that. The message was too long."
				meters.CWrongQueryCounter.Add(ctx, 1)
				updateSuccessMeter(ctx, meters, agentResp.Result, "No trace Id from flows")
				json.NewEncoder(w).Encode(agentResp)
				return
			}
			ch, err := getHistory(ctx, user)
			if err != nil {
				slog.ErrorContext(ctx, "Error while fetching history")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			agentResp := chatSingleFlow(ctx, deps, metadata, ch, user, chatRequest.Content, meters)
			saveHistory(ctx, ch, user, metadata)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(agentResp)
			return

		}
	}
}
