package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/vertexai"
	_ "github.com/lib/pq"
	_ "github.com/movie-guru/pkg/types"
)

type MovieContext struct {
	Title          string   `json:"title"`
	RuntimeMinutes int      `json:"runtime_minutes"`
	Genres         []string `json:"genres"`
	Rating         float32  `json:"rating"`
	Plot           string   `json:"plot"`
	Released       int      `json:"released"`
	Director       string   `json:"director"`
	Actors         []string `json:"actors"`
	Poster         string   `json:"poster"`
	Tconst         string   `json:"tconst"`
}

type RetrieverFlowInput struct {
	Query string `json:"query"`
}
type RetrieverFlowOutput struct {
	Documents []*ai.Document `json:"documents"`
}

func ParseMovieContexts(docs []*ai.Document) ([]*MovieContext, error) {
	movies := make([]*MovieContext, 0, len(docs))

	for _, doc := range docs {
		var intermediate struct {
			Title    string `json:"title"`
			Genres   string `json:"genres"`
			Actors   string `json:"actors"`
			Director string `json:"director"`
			Plot     string `json:"plot"`
			Poster   string `json:"poster"`
		}

		err := json.Unmarshal([]byte(doc.Content[0].Text), &intermediate)
		if err != nil {
			return nil, err
		}

		rating, _ := doc.Metadata["rating"].(float32)
		runTimeMins, _ := doc.Metadata["runtime_minutes"].(int)
		released, _ := doc.Metadata["releases"].(int)
		poster := doc.Metadata["poster"].(string)
		movies = append(movies, &MovieContext{
			Title:          intermediate.Title,
			RuntimeMinutes: runTimeMins,
			Genres:         strings.Split(intermediate.Genres, ", "),
			Rating:         rating,
			Plot:           intermediate.Plot,
			Released:       released,
			Director:       intermediate.Director,
			Actors:         strings.Split(intermediate.Actors, ", "),
			Poster:         poster,
		})
	}

	return movies, nil
}

type MovieContextList struct {
	Movies []*MovieContext `json:"movies"`
}

type MovieRetriever struct {
	DB              *sql.DB
	RetrieverLength int
	Retriever       ai.Retriever
}

func (m *MovieRetriever) RetriveDocuments(ctx context.Context, query string) ([]*MovieContext, error) {
	doc := ai.DocumentFromText(query, nil)
	retDoc := ai.RetrieverRequest{
		Document: doc,
		Options:  m.RetrieverLength,
	}
	rResp, err := m.Retriever.Retrieve(ctx, &retDoc)
	if err != nil {
		return nil, err
	}
	return ParseMovieContexts(rResp.Documents)
}

func GetEmbedder(embeddingModelName string) ai.Embedder {
	embedder := vertexai.Embedder(embeddingModelName)
	return embedder
}

func CreateMovieRetriever(embeddingModelName string, maxRetLength int, db *sql.DB) *MovieRetriever {
	embedder := GetEmbedder(embeddingModelName)
	ret := DefineRetriever(maxRetLength, db, embedder)
	return &MovieRetriever{
		DB:              db,
		RetrieverLength: maxRetLength,
		Retriever:       ret,
	}
}

func GetRetrieverFlow(ctx context.Context, ret ai.Retriever) *genkit.Flow[*RetrieverFlowInput, *RetrieverFlowOutput, struct{}] {
	retFlow := genkit.DefineFlow("movieDocFlow",
		func(ctx context.Context, input *RetrieverFlowInput) (*RetrieverFlowOutput, error) {
			doc := ai.DocumentFromText(input.Query, nil)
			query := &ai.RetrieverRequest{
				Document: doc,
				Options:  10,
			}
			retOutput := make([]*ai.Document, 0, 10)
			retFlowOutput := &RetrieverFlowOutput{
				Documents: retOutput,
			}
			resp, err := ret.Retrieve(ctx, query)
			if err != nil {
				if blockedErr, ok := err.(*genai.BlockedError); ok {
					fmt.Println("Request was blocked:", blockedErr)
					return retFlowOutput, nil
				} else {
					return nil, err
				}
			}
			t := resp.Documents
			retFlowOutput.Documents = t
			return retFlowOutput, nil
		})
	return retFlow
}

func DefineRetriever(maxRetLength int, db *sql.DB, embedder ai.Embedder) ai.Retriever {
	// Defining the Retriever
	f := func(ctx context.Context, req *ai.RetrieverRequest) (*ai.RetrieverResponse, error) {
		// Create an empty default response
		retrieverResponse := &ai.RetrieverResponse{}
		retrieverResponse.Documents = make([]*ai.Document, 0, maxRetLength)

		// Steps
		// 1. Get the embedding for the query which is in req.Document
		// 2. Query the db for the relevant rows based on the embedding
		// 3. Iterate through the rows and create a list of type ai.Document
		// 4. Return the content field from the row as content in the document and the remaining fields as metadata
		// 5. Store the list as documents in the retrieverResponse
		//
		// HINT: https://github.com/firebase/genkit/blob/main/go/samples/pgvector/main.go

		// Why content and metadata?
		// We separate movie data into 'content' and 'metadata' to accommodate varying approaches to data handling in GenAI frameworks.
		// Some frameworks, particularly those focused on RAG and utilizing a 'Document' object,
		// primarily use the 'content' field during RAG, potentially ignoring 'metadata'.

		// This separation is partly rooted in the historical context of these frameworks, which were often initially designed
		// to work with document-style databases rather than relational databases.
		// In document dbs, all the informational content is contained in the content of the document and not its metadata.
		// But in a relational db, the information may be spread across different columns.

		// In our application (using Genkit), we have the flexibility to pass a custom 'MovieContext' object into the RAG flow (next challenge)
		// (and not restricted to document.content)
		// However, when interacting with other frameworks, especially those relying on a 'Document' structure,
		// it's crucial to be mindful of how metadata is utilized or if adjustments are needed to ensure all essential information is included.

		// Actually if you look at RetriveDocuments and ParseMovieContexts, we even throw away the content and only process the data in the metadata fields while constructing the MovieContext.
		return retrieverResponse, nil
	}
	// Return the Retriever
	return ai.DefineRetriever("pgvector", "movieRetriever", f)
}
