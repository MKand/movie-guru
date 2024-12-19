package wrappers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/movie-guru/pkg/db"
	types "github.com/movie-guru/pkg/types"
	utils "github.com/movie-guru/pkg/utils"
)

type MovieFlowClient struct {
	MovieDB *db.MovieDB
	URL     string
}

func CreateMovieFlowClient(db *db.MovieDB, URL string) (*MovieFlowClient, error) {
	return &MovieFlowClient{
		MovieDB: db,
		URL:     URL + "/movieQAFlow",
	}, nil
}

func (flowClient *MovieFlowClient) Run(movieDocs []*types.MovieContext, history []*types.SimpleMessage, userPreferences *types.UserProfile) (*types.AgentResponse, error) {
	input := &types.MovieFlowInput{
		History:          history,
		UserPreferences:  userPreferences,
		ContextDocuments: movieDocs,
		UserMessage:      history[len(history)-1].Content,
	}
	resp, err := flowClient.runFlow(input)
	if err != nil {
		return nil, err
	}

	relevantMovies := make([]string, 0, len(resp.RelevantMoviesTitles))
	for _, r := range resp.RelevantMoviesTitles {
		relevantMovies = append(relevantMovies, r.Title)
	}

	agentResponse := &types.AgentResponse{
		Answer:         resp.Answer,
		RelevantMovies: relevantMovies,
		Context:        utils.FilterRelevantContext(relevantMovies, movieDocs),
		ErrorMessage:   "",
		Result:         types.SUCCESS,
		Preferences:    userPreferences,
	}
	return agentResponse, nil
}

func (flowClient *MovieFlowClient) runFlow(input *types.MovieFlowInput) (*types.MovieFlowOutput, error) {
	// Marshal the input struct to JSON
	dataInput := DataInput{
		Data: input,
	}

	inputJSON, err := json.Marshal(dataInput)
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
		Result *types.MovieFlowOutput `json:"result"`
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return nil, err
	}

	return result.Result, nil
}
