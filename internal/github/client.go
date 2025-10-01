package github

import (
	"github.com/cli/go-gh/v2/pkg/api"
	"strings"
)

// Client is a GitHub API client wrapper
type Client struct {
	client      *api.RESTClient
	excludeBots bool
}

// NewClient creates a new GitHub API client
func NewClient(excludeBots bool) (*Client, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, err
	}

	return &Client{
		client:      client,
		excludeBots: excludeBots,
	}, nil
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
