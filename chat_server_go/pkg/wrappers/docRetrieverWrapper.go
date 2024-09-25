package wrappers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/firebase/genkit/go/ai"
	_ "github.com/lib/pq"
	_ "github.com/movie-guru/pkg/types"
	types "github.com/movie-guru/pkg/types"
)

func parseMovieContexts(docs []*ai.Document) ([]*types.MovieContext, error) {
	movies := make([]*types.MovieContext, 0, len(docs))

	for _, doc := range docs {
		var intermediate struct {
			Title       string `json:"title"`
			RuntimeMins int    `json:"runtime_mins"`
			Genres      string `json:"genres"`
			Rating      string `json:"rating"`
			Released    int    `json:"released"`
			Actors      string `json:"actors"`
			Director    string `json:"director"`
			Plot        string `json:"plot"`
		}

		err := json.Unmarshal([]byte(doc.Content[0].Text), &intermediate)
		if err != nil {
			return nil, err
		}

		rating, err := strconv.ParseFloat(intermediate.Rating, 32)
		if err != nil {
			rating = 0
		}

		movies = append(movies, &types.MovieContext{
			Title:          intermediate.Title,
			RuntimeMinutes: intermediate.RuntimeMins,
			Genres:         strings.Split(intermediate.Genres, ", "),
			Rating:         float32(rating),
			Plot:           intermediate.Plot,
			Released:       intermediate.Released,
			Director:       intermediate.Director,
			Actors:         strings.Split(intermediate.Actors, ", "),
			Poster:         doc.Metadata["poster"].(string),
		})
	}

	return movies, nil
}

type MovieRetrieverFlowClient struct {
	RetrieverLength int
	URL             string
}

func CreateMovieRetrieverFlowClient(retrieverLength int, url string) *MovieRetrieverFlowClient {
	return &MovieRetrieverFlowClient{
		RetrieverLength: retrieverLength,
		URL:             url + "/movieDocFlow",
	}
}

func (flowClient *MovieRetrieverFlowClient) RetriveDocuments(ctx context.Context, query string) ([]*types.MovieContext, error) {
	doc := ai.DocumentFromText(query, nil)
	retDoc := ai.RetrieverRequest{
		Document: doc,
		Options:  flowClient.RetrieverLength,
	}
	rResp, err := flowClient.runFlow(retDoc)
	if err != nil {
		return nil, err
	}
	return parseMovieContexts(rResp)
}

func (flowClient *MovieRetrieverFlowClient) runFlow(retRequest ai.RetrieverRequest) ([]*ai.Document, error) {
	// Marshal the input struct to JSON
	inputJSON, err := json.Marshal(retRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input to JSON: %w", err)
	}
	req, err := http.NewRequest("POST", flowClient.URL, bytes.NewBuffer(inputJSON))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	var result struct {
		Result []*ai.Document `json:"result"`
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return nil, err
	}

	return result.Result, nil
}
