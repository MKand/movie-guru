package types

import "github.com/firebase/genkit/go/ai"

type SimpleMessage struct {
	Role    string `json:"sender"`
	Content string `json:"message"`
}

type ChatHistory struct {
	History []*ai.Message
}
