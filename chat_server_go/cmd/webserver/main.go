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

package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/movie-guru/pkg/db"
	met "github.com/movie-guru/pkg/metrics"
	web "github.com/movie-guru/pkg/web"
	wrappers "github.com/movie-guru/pkg/wrappers"
)

func main() {

	ctx := context.Background()

	// Load environment variables
	URL := os.Getenv("FLOWS_URL")
	metricsEnabled, err := strconv.ParseBool(os.Getenv("ENABLE_METRICS"))

	if err != nil {
		slog.WarnContext(ctx, "Error getting ENABLE_METRICS, setting to false.", slog.Any("error", err))
		metricsEnabled = false
	}
	// Set up database
	movieAgentDB, err := db.GetDB()
	if err != nil {
		slog.ErrorContext(ctx, "Error setting up DB", slog.Any("error", err))
		os.Exit(1)
	}
	defer movieAgentDB.DB.Close()

	// test redis connection
	web.TestRedis()

	// Fetch metadata
	metadata, err := movieAgentDB.GetMetadata(ctx, os.Getenv("APP_VERSION"))
	if err != nil {
		slog.ErrorContext(ctx, "Error getting metadata", slog.Any("error", err))
		os.Exit(1)
	}

	// Set up dependencies
	ulh := web.NewUserLoginHandler(metadata.TokenAudience, movieAgentDB)
	deps := getDependencies(ctx, metadata, movieAgentDB, URL)

	// Start telemetry if metrics are enabled
	if metricsEnabled {
		if shutdown, err := met.SetupOpenTelemetry(ctx); err != nil {
			slog.ErrorContext(ctx, "Error setting up OpenTelemetry", slog.Any("error", err))
			os.Exit(1)
		} else {
			defer shutdown(ctx)
		}
	}

	// Start the server
	if err := web.StartServer(ctx, ulh, metadata, deps); err != nil {
		slog.ErrorContext(ctx, "Server exited with error", slog.Any("error", err))
		os.Exit(1)
	}
}

func getDependencies(ctx context.Context, metadata *db.Metadata, db *db.MovieDB, url string) *web.Dependencies {

	queryTransformFlowClient, err := wrappers.CreateQueryTransformFlowClient(db, url)
	if err != nil {
		slog.ErrorContext(ctx, "error setting up queryTransformFlowClient client")

	}
	userProfileFlowClient, err := wrappers.CreateUserProfileFlowClient(db, url)
	if err != nil {
		slog.ErrorContext(ctx, "error setting up userProfileFlowClient client")
	}

	movieRetrieverFlowClient := wrappers.CreateMovieRetrieverFlowClient(metadata.RetrieverLength, url)

	movieFlowClient, err := wrappers.CreateMovieFlowClient(db, url)
	if err != nil {
		slog.ErrorContext(ctx, "error setting up movieFlowClient client")
	}

	responseQualityFlowClient, err := wrappers.CreateResponseQualityFlowClient(url)
	if err != nil {
		slog.ErrorContext(ctx, "error setting up responseQualityFlowClient client")
	}

	chatFlowClient, err := wrappers.CreateChatFlowClient(url)

	deps := &web.Dependencies{
		QueryTransformFlowClient:  queryTransformFlowClient,
		UserProfileFlowClient:     userProfileFlowClient,
		MovieFlowClient:           movieFlowClient,
		MovieRetrieverFlowClient:  movieRetrieverFlowClient,
		ResponseQualityFlowClient: responseQualityFlowClient,
		ChatFlowClient:            chatFlowClient,
		DB:                        db,
	}
	return deps
}
