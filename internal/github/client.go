package github

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/hazadus/gh-repomon/internal/errors"
	"github.com/hazadus/gh-repomon/internal/types"
)

// Client is a GitHub API client wrapper
type Client struct {
	client      *api.RESTClient
	excludeBots bool
	userCache   map[string]*types.Author
	cacheMutex  sync.RWMutex
}

// NewClient creates a new GitHub API client
func NewClient(excludeBots bool) (*Client, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, errors.NewGitHubAuthError("failed to create GitHub REST client", err)
	}

	return &Client{
		client:      client,
		excludeBots: excludeBots,
		userCache:   make(map[string]*types.Author),
	}, nil
}

// getUserInfo retrieves user information from cache or API
func (c *Client) getUserInfo(login string) (*types.Author, error) {
	// Check cache first
	c.cacheMutex.RLock()
	if author, found := c.userCache[login]; found {
		c.cacheMutex.RUnlock()
		return author, nil
	}
	c.cacheMutex.RUnlock()

	// Fetch from API
	var response struct {
		Login      string `json:"login"`
		Name       string `json:"name"`
		ProfileURL string `json:"html_url"`
		Type       string `json:"type"`
	}

	err := c.doWithRetry("GET", "users/"+login, nil, &response)
	if err != nil {
		// Return basic author info on error
		return &types.Author{
			Login:      login,
			Name:       login,
			ProfileURL: "https://github.com/" + login,
			IsBot:      c.isBot(login),
		}, nil
	}

	author := &types.Author{
		Login:      response.Login,
		Name:       response.Name,
		ProfileURL: response.ProfileURL,
		IsBot:      response.Type == "Bot" || c.isBot(login),
	}

	// Cache the result
	c.cacheMutex.Lock()
	c.userCache[login] = author
	c.cacheMutex.Unlock()

	return author, nil
}

// doWithRetry performs a GET request with retry logic for transient errors
func (c *Client) doWithRetry(method, path string, body interface{}, response interface{}) error {
	maxRetries := 3
	retryDelay := time.Second

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			time.Sleep(retryDelay)
			retryDelay *= 2 // Exponential backoff
		}

		// We only support GET requests with retry for now
		var err error
		if method == "GET" {
			err = c.client.Get(path, response)
		} else {
			return errors.NewGitHubAPIError("unsupported method for retry", 0, nil)
		}

		if err == nil {
			return nil
		}

		lastErr = err

		// Check if it's a rate limit error (status 403)
		if strings.Contains(err.Error(), "403") {
			return errors.NewGitHubAPIError("rate limit exceeded", http.StatusForbidden, err)
		}

		// Check if it's a 404 error (not found)
		if strings.Contains(err.Error(), "404") {
			return errors.NewGitHubAPIError("resource not found", http.StatusNotFound, err)
		}

		// Retry on network errors or 5xx errors
		if isRetryableError(err) {
			continue
		}

		// Non-retryable error, return immediately
		return errors.NewGitHubAPIError("API request failed", 0, err)
	}

	return errors.NewGitHubAPIError("API request failed after retries", 0, lastErr)
}

// isRetryableError checks if an error is worth retrying
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	// Network errors and 5xx errors are retryable
	return strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "connection") ||
		strings.Contains(errStr, "500") ||
		strings.Contains(errStr, "502") ||
		strings.Contains(errStr, "503") ||
		strings.Contains(errStr, "504")
}

// isBot checks if a login belongs to a bot account
func (c *Client) isBot(login string) bool {
	// Check if login ends with [bot]
	if strings.HasSuffix(login, "[bot]") {
		return true
	}

	// Check known bot accounts
	knownBots := []string{"github-actions", "dependabot", "renovate"}
	for _, bot := range knownBots {
		if login == bot {
			return true
		}
	}

	return false
}
