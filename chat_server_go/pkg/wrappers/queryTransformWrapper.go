package wrappers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/movie-guru/pkg/db"
	types "github.com/movie-guru/pkg/types"
)

type QueryTransformAgent struct {
	MovieAgentDB *db.MovieAgentDB
	URL          string
}

func CreateQueryTransformAgent(db *db.MovieAgentDB, URL string) (*QueryTransformAgent, error) {
	return &QueryTransformAgent{
		MovieAgentDB: db,
		URL:          URL + "/queryTransformFlow",
	}, nil
}

func (q *QueryTransformAgent) Run(history []*types.SimpleMessage, preferences *types.UserProfile) (*types.QueryTransformOutput, error) {
	queryTransformInput := types.QueryTransformInput{Profile: preferences, History: history, UserMessage: history[len(history)-1].Content}
	resp, err := q.runFlow(&queryTransformInput)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (agent *QueryTransformAgent) runFlow(input *types.QueryTransformInput) (*types.QueryTransformOutput, error) {
	// Marshal the input struct to JSON
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input to JSON: %w", err)
	}
	req, err := http.NewRequest("POST", agent.URL, bytes.NewBuffer(inputJSON))
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
		Result *types.QueryTransformOutput `json:"result"`
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return nil, err
	}

	return result.Result, nil
}
