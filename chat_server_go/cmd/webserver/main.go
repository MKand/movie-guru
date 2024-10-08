package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/vertexai"
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
		log.Fatal(err)
	}
	defer movieAgentDB.DB.Close()

	metadata, err := movieAgentDB.GetMetadata(ctx, os.Getenv("APP_VERSION"))
	if err != nil {
		log.Fatal(err)
	}

	ulh := web.NewUserLoginHandler(metadata.TokenAudience, movieAgentDB)
	deps := getDependencies(ctx, metadata, movieAgentDB, URL)

	if err = errors.Join(web.StartServer(ulh, metadata, deps), shutdown(ctx)); err != nil {
		slog.ErrorContext(ctx, "server exited with error", slog.Any("error", err))
		os.Exit(1)
	}

	if err := genkit.Init(ctx, nil); err != nil {
		log.Fatal(err)
	}

}

func getDependencies(ctx context.Context, metadata *db.Metadata, db *db.MovieDB, url string) *web.Dependencies {
	model := vertexai.Model(metadata.GoogleChatModelName)

	if model == nil {
		log.Fatal("Model not found")
	}
	queryTransformFlowClient, err := wrappers.CreateQueryTransformFlowClient(db, url)
	if err != nil {
		log.Fatal(err)
	}
	userProfileFlowClient, err := wrappers.CreateUserProfileFlowClient(db, url)
	if err != nil {
		log.Fatal(err)
	}

	movieRetrieverFlowClient := wrappers.CreateMovieRetrieverFlowClient(metadata.RetrieverLength, url)

	movieFlowClient, err := wrappers.CreateMovieFlowClient(db, url)
	if err != nil {
		log.Fatal(err)
	}

	responseQualityFlowClient, err := wrappers.CreateResponseQualityFlowClient(url)
	if err != nil {
		log.Fatal(err)
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
