package agents

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/vertexai"

	"github.com/movie-guru/pkg/db"
	types "github.com/movie-guru/pkg/types"
)

type AgentDependencies struct {
	QueryTransformFlow *genkit.Flow[*types.QueryTransformInput, *types.QueryTransformOutput, struct{}]
	PrefFlow           *genkit.Flow[*types.ProfileAgentInput, *types.UserProfileAgentOutput, struct{}]
	MovieFlow          *genkit.Flow[*types.MovieAgentInput, *types.MovieAgentOutput, struct{}]
	RetFlow            *genkit.Flow[*ai.RetrieverRequest, []*ai.Document, struct{}]
	Retriever          ai.Retriever
	DB                 *sql.DB
}

func GetDependencies(ctx context.Context, metadata *db.Metadata, db *sql.DB) *AgentDependencies {
	err := vertexai.Init(ctx, &vertexai.Config{ProjectID: os.Getenv("PROJECT_ID"), Location: os.Getenv("GCLOUD_LOCATION")})

	if err != nil {
		log.Fatal(err)
	}

	model := vertexai.Model(metadata.GoogleChatModelName)

	if model == nil {
		log.Fatal("Model not found")
	}

	queryTransformFlow, err := GetQueryTransformFlow(ctx, model)
	if err != nil {
		log.Fatal(err)
	}
	userProfileFlow, err := GetUserProfileFlow(ctx, model)
	if err != nil {
		log.Fatal(err)
	}

	embedder := GetEmbedder(metadata.GoogleEmbeddingModelName)
	if embedder == nil {
		log.Fatal("Embedder not found")
	}
	ret := DefineRetriever(metadata.RetrieverLength, db, embedder)
	retFlow := DefineRetFlow(ctx, ret)

	movieAgentFlow, err := GetMovieAgentFlow(ctx, model)
	if err != nil {
		log.Fatal(err)
	}

	deps := &AgentDependencies{
		QueryTransformFlow: queryTransformFlow,
		PrefFlow:           userProfileFlow,
		MovieFlow:          movieAgentFlow,
		Retriever:          ret,
		RetFlow:            retFlow,
		DB:                 db,
	}
	return deps
}
