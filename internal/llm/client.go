package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/hazadus/gh-repomon/internal/errors"
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
		return nil, errors.NewLLMAPIError("failed to get GitHub token", 0, err)
	}

	token := strings.TrimSpace(string(output))
	if token == "" {
		return nil, errors.NewLLMAPIError("GitHub token is empty", 0, nil)
	}

	return &Client{
		token:    token,
		endpoint: GitHubModelsEndpoint,
	}, nil
}

// Complete sends a chat completion request and returns the response text
func (c *Client) Complete(request ChatCompletionRequest) (string, error) {
	// Create context with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Prepare request body
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", errors.NewLLMAPIError("failed to marshal request", 0, err)
	}

	// Create HTTP request with context
	url := c.endpoint + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", errors.NewLLMAPIError("failed to create request", 0, err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Check if it's a timeout error
		if ctx.Err() == context.DeadlineExceeded {
			return "", errors.NewLLMAPIError("request timeout", 0, err)
		}
		return "", errors.NewLLMAPIError("failed to send request", 0, err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.NewLLMAPIError("failed to read response", resp.StatusCode, err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", errors.NewLLMAPIError(fmt.Sprintf("API request failed: %s", string(body)), resp.StatusCode, nil)
	}

	// Parse response
	var response ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", errors.NewLLMAPIError("failed to unmarshal response", resp.StatusCode, err)
	}

	// Check if response has choices
	if len(response.Choices) == 0 {
		return "", errors.NewLLMAPIError("no choices in response", resp.StatusCode, nil)
	}

	return response.Choices[0].Message.Content, nil
}
