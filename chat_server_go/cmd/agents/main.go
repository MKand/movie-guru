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

	metadata, err := movieAgentDB.GetServerMetadata(os.Getenv("APP_VERSION"))
	if err != nil {
		log.Fatal(err)
	}

	GetDependencies(ctx, metadata, movieAgentDB.DB)

	if err := genkit.Init(ctx, nil); err != nil {
		log.Fatal(err)
	}

}
