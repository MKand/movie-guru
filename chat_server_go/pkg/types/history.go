package types

import (
	"encoding/json"
	"errors"

	"github.com/firebase/genkit/go/ai"
)

type SimpleMessage struct {
	Role    string `json:"sender"`
	Content string `json:"message"`
}

type ChatHistory struct {
	History []*ai.Message
}

func (ch *ChatHistory) MarshalBinary() ([]byte, error) {
	// Logic to convert your ChatHistory object into a byte slice
	// You can use JSON, Gob, or any other serialization method you prefer

	// Example using JSON:
	data, err := json.Marshal(ch)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (ch *ChatHistory) UnmarshalBinary(data []byte) error {
	// Logic to convert the byte slice back into a ChatHistory object

	// Example using JSON:
	return json.Unmarshal(data, ch)
}

func (m *ChatHistory) Trim(maxLength int) int {
	startIndex := 0

	if len(m.History) >= maxLength {
		startIndex = len(m.History) - maxLength
	}
	recentMessages := m.History[startIndex:]
	m.History = recentMessages
	return len(m.History)
}

func NewChatHistory() *ChatHistory {
	return &ChatHistory{
		History: []*ai.Message{},
	}
}

func ParseRecentHistory(aiMessages []*ai.Message, maxLength int) ([]*SimpleMessage, error) {
	startIndex := 0
	if len(aiMessages) >= maxLength {
		startIndex = len(aiMessages) - maxLength
	}
	recentMessages := aiMessages[startIndex:]
	messages := make([]*SimpleMessage, 0, maxLength)
	for _, aiMessage := range recentMessages {

		role := ""
		if aiMessage.Role == "user" {
			role = "user"
		} else {
			role = "agent"
		}
		if aiMessage.Role == "system" {
			role = "system"
		}
		message := &SimpleMessage{
			Role:    role,
			Content: aiMessage.Content[0].Text,
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (m *ChatHistory) GetLastMessage() (string, error) {
	if len(m.History) > 0 {
		message := m.History[len(m.History)-1]
		return message.Content[0].Text, nil
	}
	return "", errors.New("No messages found")
}

func (m *ChatHistory) AddUserMessage(message string) {
	m.History = append(m.History, ai.NewUserTextMessage(message))
}

func (m *ChatHistory) AddAgentMessage(message string) {
	m.History = append(m.History, ai.NewModelTextMessage(message))
}

func (m *ChatHistory) AddAgentErrorMessage() {
	m.History = append(m.History, ai.NewModelTextMessage("Something went wrong. Try again."))
}

func (m *ChatHistory) AddSafetyIssueErrorMessage() {
	m.History = append(m.History, ai.NewModelTextMessage("That was a naughty request. I cannot proces it."))
}

func (m ChatHistory) GetHistory() []*ai.Message {
	return m.History
}
