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
	"encoding/json"
	"log/slog"
	"math"
	"math/rand"
	"net/http"

	types "github.com/movie-guru/pkg/types"
)

func randomisedFeaturedFilmsQuery() string {
	queries := []string{
		"great films", "best films", "new films", "high rated films",
	}
	return queries[rand.Intn(len(queries))]

}

func createStartupHandler(deps *Dependencies) http.HandlerFunc {
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
		if r.Method == "GET" {
			user := sessionInfo.User
			pref, err := deps.DB.GetCurrentProfile(ctx, user)
			if err != nil {
				slog.ErrorContext(ctx, "Cannot get preferences", slog.String("user", user), slog.Any("error", err.Error()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			context, err := deps.MovieRetrieverFlowClient.RetriveDocuments(ctx, randomisedFeaturedFilmsQuery())
			if err != nil {
				slog.ErrorContext(ctx, "Error getting movie recommendations", slog.Any("error", err.Error()))
				http.Error(w, "Error get movie recommendations", http.StatusInternalServerError)
				return
			}
			agentResp := types.NewAgentResponse()
			agentResp.Context = truncateContext(ctx, context)
			agentResp.Preferences = pref
			agentResp.Result = types.SUCCESS

			json.NewEncoder(w).Encode(agentResp)
			return

		}
	}
}

func truncateContext(ctx context.Context, context []*types.MovieContext) []*types.MovieContext {
	if len(context) == 0 {
		slog.InfoContext(ctx, "0 context documents recieved at startup")
		return context
	}
	return context[:int(math.Min(float64(len(context)), 5))]
}
