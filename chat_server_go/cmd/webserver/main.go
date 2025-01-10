package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/vertexai"
	"github.com/movie-guru/pkg/db"
	met "github.com/movie-guru/pkg/metrics"
	web "github.com/movie-guru/pkg/web"
	wrappers "github.com/movie-guru/pkg/wrappers"
)

func main() {
	ctx := context.Background()

	// Load environment variables
	URL := os.Getenv("FLOWS_URL")
	metricsEnabled, _ := strconv.ParseBool(os.Getenv("ENABLE_METRICS"))

	// Set up database
	movieAgentDB, err := db.GetDB()
	if err != nil {
		slog.ErrorContext(ctx, "Error setting up DB", slog.Any("error", err))
		os.Exit(1)
	}
	defer movieAgentDB.DB.Close()

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

	// Initialize genkit
	if err := genkit.Init(ctx, nil); err != nil {
		slog.ErrorContext(ctx, "Error setting up genkit", slog.Any("error", err))
		os.Exit(1)
	}
}

func getDependencies(ctx context.Context, metadata *db.Metadata, db *db.MovieDB, url string) *web.Dependencies {
	model := vertexai.Model(metadata.GoogleChatModelName)

	if model == nil {
		slog.ErrorContext(ctx, "error getting model", slog.Any("model name", metadata.GoogleChatModelName))
	}
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

	deps := &web.Dependencies{
		QueryTransformFlowClient:  queryTransformFlowClient,
		UserProfileFlowClient:     userProfileFlowClient,
		MovieFlowClient:           movieFlowClient,
		MovieRetrieverFlowClient:  movieRetrieverFlowClient,
		ResponseQualityFlowClient: responseQualityFlowClient,
		DB:                        db,
	}
	return deps
}
