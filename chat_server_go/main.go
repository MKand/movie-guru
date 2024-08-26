package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/vertexai"
)

func main() {
	ctx := context.Background()
	err := vertexai.Init(ctx, &vertexai.Config{ProjectID: os.Getenv("PROJECT_ID"), Location: os.Getenv("GCLOUD_LOCATION")})

	if err != nil {
		log.Fatal(err)
	}

	db, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	metadata, err := getServerMetadata(os.Getenv("APP_VERSION"), db)
	if err != nil {
		log.Fatal(err)
	}

	deps := getDependencies(ctx, metadata, db)

	ulh := UserLoginHandler{db: db, tokenAudience: metadata.TokenAudience}
	startServer(ulh, metadata, deps)

	if err := genkit.Init(ctx, nil); err != nil {
		log.Fatal(err)
	}

}

func getDependencies(ctx context.Context, metadata *Metadata, db *sql.DB) *ChatDependencies {
	model := vertexai.Model(metadata.GoogleChatModelName)

	if model == nil {
		log.Fatal("Model not found")
	}
	queryTransformAgent, err := CreateQueryTransformAgent(ctx, model)
	if err != nil {
		log.Fatal(err)
	}
	prefAgent, err := CreatePreferencesAgent(ctx, model, db)
	if err != nil {
		log.Fatal(err)
	}

	ret := CreateMovieRetriever(metadata.GoogleEmbeddingModelName, metadata.RetrieverLength, db)

	movieAgent, err := CreateMovieAgent(ctx, model)
	if err != nil {
		log.Fatal(err)
	}

	deps := &ChatDependencies{
		QueryTransformAgent: queryTransformAgent,
		PrefAgent:           prefAgent,
		MovieAgent:          movieAgent,
		Retriever:           ret,
		DB:                  db,
	}
	return deps
}
