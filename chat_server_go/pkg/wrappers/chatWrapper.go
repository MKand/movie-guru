// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wrappers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	types "github.com/movie-guru/pkg/types"
	utils "github.com/movie-guru/pkg/utils"
)

type ChatFlowClient struct {
	URL              string
	PosterBucketName string
}

type ChatFlowOutput struct {
	Answer               string                 `json:"answer"`
	RelevantMoviesTitles []*types.RelevantMovie `json:"relevantMovies"`
	ContextDocuments     []*types.MovieContext  `json:"contextDocuments"`
	Justification        string                 `json:"justification"`
	BadQuery             bool                   `json:"badQuery,omitempty" `
	SafetyIssue          bool                   `json:"safetyIssue,omitempty"`
	QuotaIssue           bool                   `json:"quotaIssue,omitempty"`
}

func CreateChatFlowClient(URL string, posterBucketName string) (*ChatFlowClient, error) {
	return &ChatFlowClient{
		URL:              URL + "/chatFlow",
		PosterBucketName: posterBucketName,
	}, nil
}

func (flowClient *ChatFlowClient) Run(history []*types.SimpleMessage, preferences *types.UserProfile) (*types.ChatWrapperOutput, error) {
	chatInput := types.ChatFlowInput{Profile: preferences, History: history, UserMessage: history[len(history)-1].Content}
	resp, err := flowClient.runFlow(&chatInput)
	if err != nil {
		return nil, err
	}

	err = utils.AddPosterURLs(resp.ContextDocuments, flowClient.PosterBucketName)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (flowClient *ChatFlowClient) runFlow(input *types.ChatFlowInput) (*types.ChatWrapperOutput, error) {
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
		slog.Log(context.Background(), slog.LevelError, "Error creating request", "error", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	ctx := context.Background()
	resp, err := client.Do(req)

	if err != nil {
		slog.ErrorContext(ctx, "ChatFlow: Error sending request to Flows", err.Error(), err)
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.ErrorContext(ctx, "QualityFlow: Genkit returned an Error", "errorCode", resp.StatusCode, "error: ", string(b))
		return nil, fmt.Errorf("genkit server returned error: %s (%d)", http.StatusText(resp.StatusCode), resp.StatusCode)
	}

	traceId := ""
	spanId := ""

	traceId = resp.Header.Get("x-genkit-trace-id")
	spanId = resp.Header.Get("x-genkit-span-id")

	slog.InfoContext(ctx, fmt.Sprintf("trace id: %s, span id: %s", traceId, spanId))

	var flowOutput struct {
		Result *ChatFlowOutput `json:"result"`
	}

	wrapperOutput := types.ChatWrapperOutput{
		ModelOutputMetadata: &types.ModelOutputMetadata{
			SafetyIssue:   false,
			Justification: "None Provided",
			BadQuery:      false,
		},
	}

	err = json.Unmarshal(b, &flowOutput)
	if err != nil {
		slog.Log(context.Background(), slog.LevelError, "Error unmarshaling JSON response", "error", err)
		return nil, err
	}

	wrapperOutput.Answer = flowOutput.Result.Answer
	wrapperOutput.RelevantMoviesTitles = flowOutput.Result.RelevantMoviesTitles
	wrapperOutput.ContextDocuments = flowOutput.Result.ContextDocuments
	wrapperOutput.ModelOutputMetadata.Justification = flowOutput.Result.Justification
	wrapperOutput.ModelOutputMetadata.BadQuery = flowOutput.Result.BadQuery
	wrapperOutput.ModelOutputMetadata.SafetyIssue = flowOutput.Result.SafetyIssue
	wrapperOutput.ModelOutputMetadata.QuotaIssue = flowOutput.Result.QuotaIssue
	wrapperOutput.TraceId = traceId
	wrapperOutput.SpanId = spanId

	return &wrapperOutput, nil

}
