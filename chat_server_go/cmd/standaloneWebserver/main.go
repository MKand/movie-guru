package main

import (
	"context"
	"log"
	"os"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/vertexai"
	db "github.com/movie-guru/pkg/db"
	standaloneWeb "github.com/movie-guru/pkg/standaloneWeb"
	standaloneWrappers "github.com/movie-guru/pkg/standaloneWrappers"
	web "github.com/movie-guru/pkg/web"
)

func main() {
	ctx := context.Background()
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
	deps := getDependencies(ctx, metadata, MovieDB)

	standaloneWeb.StartServer(ulh, metadata, deps)

	if err := genkit.Init(ctx, nil); err != nil {
		log.Fatal(err)
	}

}

func getDependencies(ctx context.Context, metadata *db.Metadata, db *db.MovieDB) *standaloneWeb.Dependencies {
	err := vertexai.Init(ctx, &vertexai.Config{ProjectID: os.Getenv("PROJECT_ID"), Location: os.Getenv("GCLOUD_LOCATION")})

	if err != nil {
		log.Fatal(err)
	}
	model := vertexai.Model(metadata.GoogleChatModelName)

	if model == nil {
		log.Fatal("Model not found")
	}
	queryTransformFlow, err := standaloneWrappers.CreateQueryTransformFlow(ctx, model, db)
	if err != nil {
		log.Fatal(err)
	}
	userProfileFlow, err := standaloneWrappers.CreateProfileFlow(ctx, model, db)
	if err != nil {
		log.Fatal(err)
	}

	ret := standaloneWrappers.CreateMovieRetrieverFlow(ctx, metadata.GoogleEmbeddingModelName, metadata.RetrieverLength, db)

	movieFlow, err := standaloneWrappers.CreateMovieFlow(ctx, model, db)
	if err != nil {
		log.Fatal(err)
	}

	deps := &standaloneWeb.Dependencies{
		QueryTransformFlow: queryTransformFlow,
		UserProfileFlow:    userProfileFlow,
		MovieFlow:          movieFlow,
		MovieRetrieverFlow: ret,
		DB:                 db,
	}
	return deps
}
