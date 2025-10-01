package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

const GitHubModelsEndpoint = "https://models.inference.ai.azure.com"

// Client represents an LLM client for GitHub Models API
type Client struct {
	token    string
	endpoint string
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest represents a request to the chat completions API
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// Choice represents a completion choice
type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// ChatCompletionResponse represents a response from the chat completions API
type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}

// NewClient creates a new LLM client
func NewClient() (*Client, error) {
	// Get GitHub token using gh CLI
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub token: %w", err)
	}

	token := strings.TrimSpace(string(output))
	if token == "" {
		return nil, fmt.Errorf("GitHub token is empty")
	}

	return &Client{
		token:    token,
		endpoint: GitHubModelsEndpoint,
	}, nil
}

// Complete sends a chat completion request and returns the response text
func (c *Client) Complete(request ChatCompletionRequest) (string, error) {
	// Prepare request body
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := c.endpoint + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check if response has choices
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}
