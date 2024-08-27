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

	dbase, err := db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer dbase.Close()

	metadata, err := db.GetServerMetadata(os.Getenv("APP_VERSION"), dbase)
	if err != nil {
		log.Fatal(err)
	}

	GetDependencies(ctx, metadata, dbase)

	if err := genkit.Init(ctx, nil); err != nil {
		log.Fatal(err)
	}

}
