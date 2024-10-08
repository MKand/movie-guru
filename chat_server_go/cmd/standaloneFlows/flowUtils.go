package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/vertexai"

	"github.com/movie-guru/pkg/db"
)

type FlowDependencies struct {
	QueryTransformFlow  *genkit.Flow[*QueryTransformFlowInput, *QueryTransformFlowOutput, struct{}]
	PrefFlow            *genkit.Flow[*UserProfileFlowInput, *UserProfileFlowOutput, struct{}]
	MovieFlow           *genkit.Flow[*MovieFlowInput, *MovieFlowOutput, struct{}]
	RetFlow             *genkit.Flow[*RetrieverFlowInput, *RetrieverFlowOutput, struct{}]
	ResponseQualityFlow *genkit.Flow[*ResponseQualityFlowInput, *ResponseQualityFlowOutput, struct{}]
	Retriever           ai.Retriever
	DB                  *sql.DB
}

type Prompts struct {
	UserPrefPrompt       string
	MovieFlowPrompt      string
	QueryTransformPrompt string
}

func GetDependencies(ctx context.Context, metadata *db.Metadata, db *sql.DB, prompts *Prompts) *FlowDependencies {
	err := vertexai.Init(ctx, &vertexai.Config{ProjectID: os.Getenv("PROJECT_ID"), Location: os.Getenv("LOCATION")})

	if err != nil {
		log.Fatal(err)
	}

	model := vertexai.Model(metadata.GoogleChatModelName)

	if model == nil {
		log.Fatal("Model not found")
	}

	userProfileFlow, err := GetUserProfileFlow(ctx, model, prompts.UserPrefPrompt)
	if err != nil {
		log.Fatal(err)
	}
	queryTransformFlow, err := GetQueryTransformFlow(ctx, model, prompts.QueryTransformPrompt)
	if err != nil {
		log.Fatal(err)
	}
	embedder := GetEmbedder(metadata.GoogleEmbeddingModelName)
	if embedder == nil {
		log.Fatal("Embedder not found")
	}
	ret := DefineRetriever(metadata.RetrieverLength, db, embedder)
	retFlow := GetRetrieverFlow(ctx, ret)

	movieAgentFlow, err := GetMovieFlow(ctx, model, prompts.MovieFlowPrompt)
	if err != nil {
		log.Fatal(err)
	}

	responseQualityFlow, err := GetResponseQualityAnalysisFlow(ctx, model, "")
	if err != nil {
		log.Fatal(err)
	}

	deps := &FlowDependencies{
		QueryTransformFlow:  queryTransformFlow,
		PrefFlow:            userProfileFlow,
		MovieFlow:           movieAgentFlow,
		ResponseQualityFlow: responseQualityFlow,
		Retriever:           ret,
		RetFlow:             retFlow,
		DB:                  db,
	}
	return deps
}

func makeJsonMarshallable(input string) (string, error) {
	// Regex to extract JSON content from Markdown code block
	re := regexp.MustCompile("```(json)?((\n|.)*?)```")
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		return input, fmt.Errorf("no JSON content found in the input")
	}

	jsonContent := matches[2]
	return jsonContent, nil
}
