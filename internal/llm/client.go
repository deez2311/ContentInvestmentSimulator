package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Client communicates with the OpenAI ChatGPT API.
type Client struct {
	apiKey     string
	httpClient *http.Client
	model      string // e.g. "gpt-3.5-turbo"
	endpoint   string // "https://api.openai.com/v1/chat/completions"
}

// chatRequest represents the request body sent to the OpenAI API.
type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

// chatMessage represents a single message in the chat conversation.
type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatResponse represents the response body from the OpenAI API.
type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

// chatChoice represents a single choice in the API response.
type chatChoice struct {
	Message chatMessage `json:"message"`
}

// NewClient creates a Client by reading OPENAI_API_KEY from the environment.
// Returns an error if the key is not set.
func NewClient() (*Client, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		model:      "gpt-3.5-turbo",
		endpoint:   "https://api.openai.com/v1/chat/completions",
	}, nil
}

// Chat sends a prompt to ChatGPT and returns the text response.
// Returns an error with HTTP status code and body on API failure.
func (c *Client) Chat(prompt string) (string, error) {
	reqBody := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("ChatGPT API error (status %d): %s", resp.StatusCode, string(body))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("ChatGPT API returned no choices")
	}

	return chatResp.Choices[0].Message.Content, nil
}
