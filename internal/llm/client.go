package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hazadus/gh-repomon/internal/errors"
	"github.com/hazadus/gh-repomon/internal/logger"
)

const GitHubModelsEndpoint = "https://models.inference.ai.azure.com"

const (
	maxRetries = 3
	baseDelay  = 1 * time.Second
)

// Client represents an LLM client for GitHub Models API
type Client struct {
	token    string
	endpoint string
}

// rateLimitError represents a parsed rate limit error response
type rateLimitError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details string `json:"details"`
	} `json:"error"`
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
// Automatically retries on rate limit errors with exponential backoff
func (c *Client) Complete(request ChatCompletionRequest) (string, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Create context with 30 second timeout for each attempt
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		// Prepare request body
		requestBody, err := json.Marshal(request)
		if err != nil {
			cancel()
			return "", errors.NewLLMAPIError("failed to marshal request", 0, err)
		}

		// Create HTTP request with context
		url := c.endpoint + "/chat/completions"
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			cancel()
			return "", errors.NewLLMAPIError("failed to create request", 0, err)
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.token)

		// Send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			cancel()
			// Check if it's a timeout error
			if ctx.Err() == context.DeadlineExceeded {
				return "", errors.NewLLMAPIError("request timeout", 0, err)
			}
			return "", errors.NewLLMAPIError("failed to send request", 0, err)
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		cancel()

		if err != nil {
			return "", errors.NewLLMAPIError("failed to read response", resp.StatusCode, err)
		}

		// Check status code
		if resp.StatusCode != http.StatusOK {
			// Handle rate limit error (429)
			if resp.StatusCode == http.StatusTooManyRequests {
				if waitTime, ok := extractWaitTime(string(body)); ok {
					// Add small buffer to wait time
					waitTime += 2 * time.Second

					if attempt < maxRetries {
						logger.Warningf("Rate limit reached, waiting %v before retry (attempt %d/%d)", waitTime, attempt+1, maxRetries)
						time.Sleep(waitTime)
						continue
					}
				}
			}

			// For other errors or if we've exhausted retries, return error
			lastErr = errors.NewLLMAPIError(fmt.Sprintf("API request failed: %s", string(body)), resp.StatusCode, nil)

			// Retry with exponential backoff for server errors (5xx)
			if resp.StatusCode >= 500 && attempt < maxRetries {
				delay := baseDelay * time.Duration(1<<uint(attempt))
				logger.Warningf("Server error (status %d), retrying in %v (attempt %d/%d)", resp.StatusCode, delay, attempt+1, maxRetries)
				time.Sleep(delay)
				continue
			}

			return "", lastErr
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

		// Success!
		if attempt > 0 {
			logger.Infof("Request succeeded after %d retries", attempt)
		}
		return response.Choices[0].Message.Content, nil
	}

	// Should not reach here, but return last error if we do
	if lastErr != nil {
		return "", lastErr
	}
	return "", errors.NewLLMAPIError("max retries exceeded", 0, nil)
}

// extractWaitTime extracts the wait time in seconds from a rate limit error message
// Expected format: "Please wait X seconds before retrying"
// Returns wait time and whether it's acceptable to wait (< 1 hour)
func extractWaitTime(errorBody string) (time.Duration, bool) {
	// Try to parse as JSON first
	var rateLimitErr rateLimitError
	if err := json.Unmarshal([]byte(errorBody), &rateLimitErr); err == nil {
		if rateLimitErr.Error.Code == "RateLimitReached" {
			// Extract seconds from message like "Please wait 29 seconds before retrying"
			// Use more specific regex to match "Please wait X second" pattern
			re := regexp.MustCompile(`Please wait (\d+) second`)
			matches := re.FindStringSubmatch(rateLimitErr.Error.Message)
			if len(matches) >= 2 {
				if seconds, err := strconv.Atoi(matches[1]); err == nil {
					waitTime := time.Duration(seconds) * time.Second

					// If wait time is more than 1 hour, it's likely a daily limit
					// Don't retry automatically - let the user know
					if seconds > 3600 {
						logger.Warningf("Rate limit exceeded: need to wait %v (%d seconds)", waitTime, seconds)
						logger.Warningf("This appears to be a daily rate limit. Consider using --no-ai flag to continue.")
						return 0, false
					}

					// For shorter waits (< 1 hour), we can retry
					if seconds > 0 {
						return waitTime, true
					}
				}
			}
		}
	}
	return 0, false
}
