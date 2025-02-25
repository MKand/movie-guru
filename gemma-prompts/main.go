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

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const userPrompt = "Pretend to be a person between the ages 18 and 80 and ask for a type of movie you want to watch. Be creative."

var maxChatLen = 750
var limiter *rate.Limiter

type Response struct {
	Predictions []string `json:"predictions,omitempty"`
}

type ChatRequest struct {
	Content string `json:"content"`
}

var chatServer string

func main() {
	var url string

	if os.Getenv("VLLM_URL") != "" {
		url = os.Getenv("VLLM_URL")
	} else {
		url = "http://localhost:8000/generate"
	}

	if os.Getenv("CHAT_SERVER") != "" {
		chatServer = os.Getenv("CHAT_SERVER")
	} else {
		fmt.Println("CHAT_SERVER not set")
		return
	}

	if os.Getenv("RATE_LIMIT") != "" {
		if r, err := strconv.ParseFloat(os.Getenv("RATE_LIMIT"), 64); err != nil {
			slog.Log(context.Background(), slog.LevelWarn, "Error parsing RATE_LIMIT, using defaults", "warning", err)
			r = 5.0
		} else {
			limiter = rate.NewLimiter(rate.Limit(r/60.0), 1)
		}
	} else {
		// Rate limiter: 5 requests per minute
		limiter = rate.NewLimiter(rate.Limit(5.0/60.0), 1)
	}

	cookie, err := getCookie()
	if err != nil {
		slog.Log(context.Background(), slog.LevelError, "Error getting cookie", "error", err)
		return
	}

	ctx := context.Background()

	go func() {
		for { // Infinite loop

			data := map[string]interface{}{
				"prompt":      fmt.Sprintf("<start_of_turn>user\n%s<end_of_turn>\n", userPrompt),
				"temperature": 0.90,
				"top_p":       1.0,
				"max_tokens":  128,
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				break
			}

			err = limiter.Wait(ctx)
			if err != nil {
				fmt.Println("Error waiting for rate limit:", err)
				break
			}

			resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Error making request:", err)
				break
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				time.Sleep(1 * time.Second)
				break
			}

			resp.Body.Close()

			r := Response{}
			err = json.Unmarshal(body, &r)

			if err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				break
			}

			index := strings.Index(r.Predictions[0], "Output:")

			prompt := r.Predictions[0][index+8:]
			prompt = strings.ReplaceAll(prompt, "\n", "")
			prompt = strings.ReplaceAll(prompt, "/\"", "")
			prompt = strings.ReplaceAll(prompt, "*", "")
			prompt = removeSpecialCharacters(prompt)

			if len(prompt) > maxChatLen {
				prompt = prompt[:maxChatLen]
			}

			slog.Log(context.Background(), slog.LevelInfo, "Prompt", "prompt", prompt)

			chatRequest := ChatRequest{
				Content: prompt,
			}

			inputJSON, err := json.Marshal(chatRequest)
			if err != nil {
				fmt.Printf("error marshaling input to JSON: %w", err)
				break
			}

			err = invokeFlow(inputJSON, cookie)
			if err != nil {
				fmt.Printf("error invoking flow: %w", err)
				break
			}

			time.Sleep(1 * time.Second) // Add a delay between requests if needed.
		}
		os.Exit(1)
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	fmt.Println("Shutting down")

	os.Exit(0)
}

func invokeFlow(d []byte, cookie string) error {
	req, err := http.NewRequest("POST", chatServer+"/chat", bytes.NewBuffer(d))
	if err != nil {
		slog.Log(context.Background(), slog.LevelError, "Error creating request", "error", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Log(context.Background(), slog.LevelError, "Error sending request", "error", err)
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Server returned error: %s\n", string(bodyBytes))
		return fmt.Errorf("server returned error: %s (%d)", http.StatusText(resp.StatusCode), resp.StatusCode)
	}

	b, _ := io.ReadAll(resp.Body)
	slog.Log(context.Background(), slog.LevelInfo, string(b))

	defer resp.Body.Close()

	return nil
}

func getCookie() (cookie string, err error) {
	req, err := http.NewRequest("POST", chatServer+"/login", bytes.NewBuffer([]byte("{\"inviteCode\":\"\"}")))
	if err != nil {
		slog.Log(context.Background(), slog.LevelError, "Error creating request", "error", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ApiKey", "fake")
	req.Header.Set("User", "fake")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Log(context.Background(), slog.LevelError, "Error sending request", "error", err)
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Server returned error: %s\n", string(bodyBytes))
		return "", fmt.Errorf("server returned error: %s (%d)", http.StatusText(resp.StatusCode), resp.StatusCode)
	}

	b, _ := io.ReadAll(resp.Body)
	slog.Log(context.Background(), slog.LevelInfo, string(b))

	cookie = resp.Header.Get("Set-Cookie")
	cookie = strings.Split(cookie, ";")[0]

	defer resp.Body.Close()

	return cookie, nil
}

func removeSpecialCharacters(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	return re.ReplaceAllString(s, "")
}
