// Package ai provides a small, stateless client for OpenAI-compatible chat
// completion APIs such as Kimi / Moonshot (https://api.moonshot.ai/v1).
package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Message is one entry in an OpenAI-compatible chat conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

const (
	// baseURL is the Kimi/Moonshot OpenAI-compatible API root.
	baseURL = "https://api.moonshot.ai/v1"
	// model is the Kimi model used for completions.
	model = "kimi-k2.7-code"
)

// client has an explicit timeout so a hung upstream call can't leak the
// goroutine handling the Discord interaction.
var client = &http.Client{Timeout: 60 * time.Second}

// Ask sends the conversation to the Kimi/Moonshot chat completions endpoint and
// returns the assistant's reply.
func Ask(apiKey string, messages []Message) (string, error) {
	payload, err := json.Marshal(chatRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   1024,
		Temperature: 1, // kimi-k2.7-code only permits temperature 1
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var parsed chatResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("ai: decoding response (status %s): %w", res.Status, err)
	}

	if res.StatusCode != http.StatusOK {
		if parsed.Error != nil && parsed.Error.Message != "" {
			return "", fmt.Errorf("ai: API error (%s): %s", res.Status, parsed.Error.Message)
		}
		return "", fmt.Errorf("ai: unexpected status %s: %s", res.Status, string(body))
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("ai: model returned no choices")
	}
	return parsed.Choices[0].Message.Content, nil
}
