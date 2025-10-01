package github

import (
	"fmt"
	"time"

	"github.com/hazadus/gh-repomon/internal/types"
)

// GetOpenPullRequests retrieves all open pull requests for a repository
func (c *Client) GetOpenPullRequests(repo string) ([]types.PullRequest, error) {
	path := fmt.Sprintf("repos/%s/pulls", repo)

	var response []map[string]interface{}
	err := c.doWithRetry("GET", path+"?state=open&per_page=100", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get open pull requests: %w", err)
	}

	var prs []types.PullRequest
	for _, prData := range response {
		pr, err := c.parsePullRequest(prData)
		if err != nil {
			continue // Skip invalid PRs
		}

		// Filter bots if requested
		if c.excludeBots && c.isBot(pr.Author.Login) {
			continue
		}

		prs = append(prs, pr)
	}

	return prs, nil
}

// GetUpdatedPullRequests retrieves pull requests updated during the specified period
func (c *Client) GetUpdatedPullRequests(repo, from, to string) ([]types.PullRequest, error) {
	path := fmt.Sprintf("repos/%s/pulls", repo)

	var response []map[string]interface{}
	err := c.doWithRetry("GET", path+"?state=all&sort=updated&direction=desc&per_page=100", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated pull requests: %w", err)
	}

	fromTime, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return nil, fmt.Errorf("invalid from date: %w", err)
	}

	toTime, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return nil, fmt.Errorf("invalid to date: %w", err)
	}

	var prs []types.PullRequest
	for _, prData := range response {
		pr, err := c.parsePullRequest(prData)
		if err != nil {
			continue
		}

		// Filter by updated_at within the period
		if pr.UpdatedAt.Before(fromTime) || pr.UpdatedAt.After(toTime) {
			continue
		}

		// Filter bots if requested
		if c.excludeBots && c.isBot(pr.Author.Login) {
			continue
		}

		prs = append(prs, pr)
	}

	return prs, nil
}

// GetPullRequestComments retrieves the number of comments on a pull request
func (c *Client) GetPullRequestComments(repo string, prNumber int) (int, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/comments", repo, prNumber)

	var response []map[string]interface{}
	err := c.doWithRetry("GET", path, nil, &response)
	if err != nil {
		return 0, fmt.Errorf("failed to get PR comments: %w", err)
	}

	return len(response), nil
}

// parsePullRequest converts GitHub API response to types.PullRequest
func (c *Client) parsePullRequest(data map[string]interface{}) (types.PullRequest, error) {
	pr := types.PullRequest{}

	// Parse number
	if num, ok := data["number"].(float64); ok {
		pr.Number = int(num)
	}

	// Parse title
	if title, ok := data["title"].(string); ok {
		pr.Title = title
	}

	// Parse body
	if body, ok := data["body"].(string); ok {
		pr.Body = body
	}

	// Parse state
	if state, ok := data["state"].(string); ok {
		pr.State = state
	}

	// Parse URL
	if htmlURL, ok := data["html_url"].(string); ok {
		pr.URL = htmlURL
	}

	// Parse author
	if user, ok := data["user"].(map[string]interface{}); ok {
		pr.Author = c.parseAuthor(user)
	}

	// Parse created_at
	if createdAt, ok := data["created_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			pr.CreatedAt = t
		}
	}

	// Parse updated_at
	if updatedAt, ok := data["updated_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
			pr.UpdatedAt = t
		}
	}

	// Parse comments count
	if comments, ok := data["comments"].(float64); ok {
		pr.Comments = int(comments)
	}

	// Reviews count will be set separately if needed
	pr.Reviews = 0

	return pr, nil
}

// parseAuthor converts GitHub API user object to types.Author
func (c *Client) parseAuthor(data map[string]interface{}) types.Author {
	author := types.Author{}

	if login, ok := data["login"].(string); ok {
		author.Login = login
		author.IsBot = c.isBot(login)
	}

	if name, ok := data["name"].(string); ok {
		author.Name = name
	}

	if htmlURL, ok := data["html_url"].(string); ok {
		author.ProfileURL = htmlURL
	}

	return author
}
