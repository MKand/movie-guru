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
	deps := getDependencies(ctx, metadata, MovieAgentDB, URL)

	web.StartServer(ulh, metadata, deps)

	if err := genkit.Init(ctx, nil); err != nil {
		log.Fatal(err)
	}

}

func getDependencies(ctx context.Context, metadata *db.Metadata, db *db.MovieAgentDB, url string) *web.Dependencies {
	model := vertexai.Model(metadata.GoogleChatModelName)

	if model == nil {
		log.Fatal("Model not found")
	}
	queryTransformAgent, err := wrappers.CreateQueryTransformAgent(db, url)
	if err != nil {
		log.Fatal(err)
	}
	prefAgent, err := wrappers.CreateProfileAgent(db, url)
	if err != nil {
		log.Fatal(err)
	}

	ret := wrappers.CreateMovieRetriever(metadata.RetrieverLength, url)

	movieAgent, err := wrappers.CreateMovieAgent(db, url)
	if err != nil {
		log.Fatal(err)
	}

	deps := &web.Dependencies{
		QueryTransformAgent: queryTransformAgent,
		PrefAgent:           prefAgent,
		MovieAgent:          movieAgent,
		Retriever:           ret,
		DB:                  db,
	}
	return deps
}
