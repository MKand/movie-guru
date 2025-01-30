package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/movie-guru/pkg/db"
	met "github.com/movie-guru/pkg/metrics"
	web "github.com/movie-guru/pkg/web"
	wrappers "github.com/movie-guru/pkg/wrappers"
)

func main() {
	ctx := context.Background()

	URL := os.Getenv("FLOWS_URL")

	shutdown, err := met.SetupOpenTelemetry(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error setting up OpenTelemetry", slog.Any("error", err))
		os.Exit(1)
	}

	movieAgentDB, err := db.GetDB()
	if err != nil {
		slog.ErrorContext(ctx, "error setting up DB", slog.Any("error", err))
	}
	defer movieAgentDB.DB.Close()

	metadata, err := movieAgentDB.GetMetadata(ctx, os.Getenv("APP_VERSION"))
	if err != nil {
		slog.ErrorContext(ctx, "error getting metadata", slog.Any("error", err))
	}

	ulh := web.NewUserLoginHandler(metadata.TokenAudience, movieAgentDB)
	deps := getDependencies(ctx, movieAgentDB, URL)

	if err = errors.Join(web.StartServer(ctx, ulh, metadata, deps), shutdown(ctx)); err != nil {
		slog.ErrorContext(ctx, "server exited with error", slog.Any("error", err))
		os.Exit(1)
	}

}

func getDependencies(ctx context.Context, db *db.MovieDB, url string) *web.Dependencies {
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
