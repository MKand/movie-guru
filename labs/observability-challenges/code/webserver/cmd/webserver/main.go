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
	"github.com/movie-guru/pkg/types"
	web "github.com/movie-guru/pkg/web"
	wrappers "github.com/movie-guru/pkg/wrappers"
)

func main() {

	ctx := context.Background()

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
	metadata, err := getMetadata(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting metadata", slog.Any("error", err))
		os.Exit(1)
	}

	// Set up dependencies
	ulh := web.NewUserLoginHandler(metadata.TokenAudience, movieAgentDB)
	deps := getDependencies(ctx, metadata, movieAgentDB)

	// Start telemetry if metrics are enabled
	if metadata.EnableMetrics {
		if shutdown, err := met.SetupOpenTelemetry(ctx); err != nil {
			slog.ErrorContext(ctx, "Error setting up OpenTelemetry", slog.Any("error", err))
			os.Exit(1)
		} else {
			defer shutdown(ctx)
		}
	}

	getLargeArray();

	// Start the server
	if err := web.StartServer(ctx, ulh, metadata, deps); err != nil {
		slog.ErrorContext(ctx, "Server exited with error", slog.Any("error", err))
		os.Exit(1)
	}
}

func getDependencies(ctx context.Context, metadata *types.Metadata, db *db.MovieDB) *web.Dependencies {

	userProfileFlowClient, err := wrappers.CreateUserProfileFlowClient(db, metadata.FlowsURL)
	if err != nil {
		slog.ErrorContext(ctx, "error setting up userProfileFlowClient client")
	}

	movieRetrieverFlowClient := wrappers.CreateMovieRetrieverFlowClient(metadata.FlowsURL, metadata.PosterBucketName)

	responseQualityFlowClient, err := wrappers.CreateResponseQualityFlowClient(metadata.FlowsURL)
	if err != nil {
		slog.ErrorContext(ctx, "error setting up responseQualityFlowClient client")
	}

	chatFlowClient, _ := wrappers.CreateChatFlowClient(metadata.FlowsURL, metadata.PosterBucketName)

	deps := &web.Dependencies{
		UserProfileFlowClient:     userProfileFlowClient,
		MovieRetrieverFlowClient:  movieRetrieverFlowClient,
		ResponseQualityFlowClient: responseQualityFlowClient,
		ChatFlowClient:            chatFlowClient,
		DB:                        db,
	}
	return deps
}

func getMetadata(ctx context.Context) (*types.Metadata, error) {
	metadata := &types.Metadata{
		UseAuth:              false,
		EnableMetrics:        false,
		StrictCors:           false,
		MaxUserMessageLength: 500,
		HistoryLength:        10,
		PosterBucketName:     "generated_posters",
	}

	posterBucketName := os.Getenv("POSTER_BUCKET_NAME")
	if posterBucketName != "" {
		metadata.PosterBucketName = posterBucketName
	}
	flowsURL := os.Getenv("FLOWS_URL")
	if flowsURL == "" {
		slog.ErrorContext(ctx, "No FlowsURL found in environment variables")
		os.Exit(1)
	}
	metadata.FlowsURL = flowsURL

	feedbackURL := os.Getenv("FEEDBACK_URL")
	if feedbackURL == "" {
		metadata.FeedbackURL = "NONE"
		slog.WarnContext(ctx, "No FeedbackURL found in environment variables. No feedback will be sent to Genkit.")
	} else {
		metadata.FeedbackURL = feedbackURL
	}

	enableMetrics, err := strconv.ParseBool(os.Getenv("ENABLE_METRICS"))
	if err == nil && enableMetrics {
		metadata.EnableMetrics = true
	}

	useAuth, err := strconv.ParseBool(os.Getenv("USE_AUTH"))
	if err == nil && useAuth {
		metadata.UseAuth = true
		tokenAudience := os.Getenv("TOKEN_AUDIENCE")
		if tokenAudience == "" {
			slog.ErrorContext(ctx, "No TokenAudience found in environment variables")
			os.Exit(1)
		}
		metadata.TokenAudience = tokenAudience

		corsOrigins := os.Getenv("CORS_ORIGINS")
		strictCors, err := strconv.ParseBool(os.Getenv("STRICT_CORS"))
		if err != nil && strictCors {
			metadata.StrictCors = true
		}
		if corsOrigins == "" && metadata.StrictCors {
			slog.ErrorContext(ctx, "No CORS_ORIGINS found in environment variables")
			os.Exit(1)
		}
		metadata.CorsOrigins = corsOrigins

		maxUserMessageLength := os.Getenv("MAX_USER_MESSAGE_LENGTH")
		if maxUserMessageLength != "" {
			metadata.MaxUserMessageLength, err = strconv.Atoi(maxUserMessageLength)
			if err != nil {
				slog.ErrorContext(ctx, "Error parsing MAX_USER_MESSAGE_LENGTH", slog.Any("error", err))
				os.Exit(1)
			}
		}

		historyLength := os.Getenv("HISTORY_LENGTH")
		if historyLength != "" {
			metadata.HistoryLength, err = strconv.Atoi(historyLength)
			if err != nil {
				slog.ErrorContext(ctx, "Error parsing HISTORY_LENGTH", slog.Any("error", err))
				os.Exit(1)
			}
		}
	}

	slog.InfoContext(ctx, "Metadata",
		slog.String("FlowsURL", metadata.FlowsURL),
		slog.String("PosterBucketName", metadata.PosterBucketName),
		slog.String("FeedbackURL", metadata.FeedbackURL),
		slog.Bool("EnableMetrics", metadata.EnableMetrics),
		slog.Int("MaxUserMessageLength", metadata.MaxUserMessageLength),
		slog.Int("HistoryLength", metadata.HistoryLength),
	)

	slog.InfoContext(ctx, "Auth Metadata",
		slog.Bool("UseAuth", metadata.UseAuth),
		slog.String("TokenAudience", metadata.TokenAudience),
		slog.String("CorsOrigins", metadata.CorsOrigins),
		slog.Bool("StrictCors", metadata.StrictCors),
	)
	return metadata, nil
}

func getLargeArray(){
	var largeSlice [][]byte
	totalAllocatedMB := 0

	// Loop indefinitely, allocating memory in chunks
	for i := 0; ; i++ {
	
		chunkSizeMB := 100
		chunk := make([]byte, chunkSizeMB*1024*1024) // 100 MB

		for j := 0; j < len(chunk); j++ {
			chunk[j] = byte(j % 256)
		}

		largeSlice = append(largeSlice, chunk)
		totalAllocatedMB += chunkSizeMB

		// Print memory usage periodically
		if i%5 == 0 { // Print every 5 chunks (500 MB)
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("Iteration %d: Allocated ~%d MB (HeapSys: %v MB, HeapAlloc: %v MB)\n",
				i, totalAllocatedMB, m.HeapSys/(1024*1024), m.HeapAlloc/(1024*1024))
		}

		time.Sleep(10 * time.Millisecond)
	}
}
