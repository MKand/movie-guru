package main

import (
	"context"
	"log"
	"os"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/vertexai"
	"github.com/movie-guru/pkg/db"
	web "github.com/movie-guru/pkg/web"
	wrappers "github.com/movie-guru/pkg/wrappers"
)

func main() {
	ctx := context.Background()

	URL := os.Getenv("FLOWS_URL")

	MovieDB, err := db.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	defer MovieDB.DB.Close()

	metadata, err := MovieDB.GetServerMetadata(os.Getenv("APP_VERSION"))
	if err != nil {
		log.Fatal(err)
	}

	ulh := web.NewUserLoginHandler(metadata.TokenAudience, MovieDB)
	deps := getDependencies(ctx, metadata, MovieDB, URL)

	web.StartServer(ulh, metadata, deps)

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

	deps := &web.Dependencies{
		QueryTransformFlowClient: queryTransformFlowClient,
		UserProfileFlowClient:    userProfileFlowClient,
		MovieFlowClient:          movieFlowClient,
		MovieRetrieverFlowClient: movieRetrieverFlowClient,
		DB:                       db,
	}
	return deps
}
