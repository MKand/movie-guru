package main

import (
	"context"
	"log"
	"os"

	"github.com/firebase/genkit/go/genkit"

	. "github.com/movie-guru/pkg/agents"
	"github.com/movie-guru/pkg/db"
)

func main() {
	ctx := context.Background()

	movieAgentDB, err := db.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	defer movieAgentDB.DB.Close()

	app_version := os.Getenv("APP_VERSION")
	metadata, err := movieAgentDB.GetServerMetadata(app_version)
	if err != nil {
		log.Fatal(err)
	}

	GetDependencies(ctx, metadata, movieAgentDB.DB)

	if err := genkit.Init(ctx, &genkit.Options{FlowAddr: ":3401"}); err != nil {
		log.Fatal(err)
	}

}
