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
	MovieAgentDB, err := db.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	defer MovieAgentDB.DB.Close()

	metadata, err := MovieAgentDB.GetServerMetadata(os.Getenv("APP_VERSION"))
	if err != nil {
		log.Fatal(err)
	}

	ulh := web.NewUserLoginHandler(metadata.TokenAudience, MovieAgentDB)
	deps := getDependencies(ctx, metadata, MovieAgentDB)

	standaloneWeb.StartServer(ulh, metadata, deps)

	if err := genkit.Init(ctx, nil); err != nil {
		log.Fatal(err)
	}

}

func getDependencies(ctx context.Context, metadata *db.Metadata, db *db.MovieAgentDB) *standaloneWeb.Dependencies {
	err := vertexai.Init(ctx, &vertexai.Config{ProjectID: os.Getenv("PROJECT_ID"), Location: os.Getenv("GCLOUD_LOCATION")})

	if err != nil {
		log.Fatal(err)
	}
	model := vertexai.Model(metadata.GoogleChatModelName)

	if model == nil {
		log.Fatal("Model not found")
	}
	queryTransformAgent, err := standaloneWrappers.CreateQueryTransformAgent(ctx, model, db)
	if err != nil {
		log.Fatal(err)
	}
	prefAgent, err := standaloneWrappers.CreateProfileAgent(ctx, model, db)
	if err != nil {
		log.Fatal(err)
	}

	ret := standaloneWrappers.CreateMovieRetriever(ctx, metadata.GoogleEmbeddingModelName, metadata.RetrieverLength, db)

	movieAgent, err := standaloneWrappers.CreateMovieAgent(ctx, model, db)
	if err != nil {
		log.Fatal(err)
	}

	deps := &standaloneWeb.Dependencies{
		QueryTransformAgent: queryTransformAgent,
		PrefAgent:           prefAgent,
		MovieAgent:          movieAgent,
		Retriever:           ret,
		DB:                  db,
	}
	return deps
}
