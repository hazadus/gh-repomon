package github

import (
	"fmt"
	"time"

	"github.com/hazadus/gh-repomon/internal/types"
)

// GetOpenIssues retrieves all open issues from the repository.
// It filters out pull requests (which GitHub API returns as issues).
func (c *Client) GetOpenIssues(repo string) ([]types.Issue, error) {
	var issues []types.Issue
	page := 1
	perPage := 100

	for {
		path := fmt.Sprintf("repos/%s/issues?state=open&per_page=%d&page=%d", repo, perPage, page)
		var response []map[string]interface{}

		err := c.client.Get(path, &response)
		if err != nil {
			return nil, fmt.Errorf("failed to get open issues: %w", err)
		}

		if len(response) == 0 {
			break
		}

		for _, item := range response {
			// Skip pull requests
			if isPullRequest(item) {
				continue
			}

			issue, err := parseIssue(item)
			if err != nil {
				continue // Skip malformed issues
			}

			// Filter bots if needed
			if c.excludeBots && c.isBot(issue.Author.Login) {
				continue
			}

			issues = append(issues, issue)
		}

		// Check if there are more pages
		if len(response) < perPage {
			break
		}
		page++
	}

	return issues, nil
}

// GetClosedIssues retrieves issues closed during the specified period.
// It filters out pull requests and only returns issues closed between from and to dates.
func (c *Client) GetClosedIssues(repo, from, to string) ([]types.Issue, error) {
	var issues []types.Issue
	page := 1
	perPage := 100

	fromTime, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return nil, fmt.Errorf("invalid from date: %w", err)
	}

	toTime, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return nil, fmt.Errorf("invalid to date: %w", err)
	}

	for {
		path := fmt.Sprintf("repos/%s/issues?state=closed&per_page=%d&page=%d&sort=updated&direction=desc", repo, perPage, page)
		var response []map[string]interface{}

		err := c.client.Get(path, &response)
		if err != nil {
			return nil, fmt.Errorf("failed to get closed issues: %w", err)
		}

		if len(response) == 0 {
			break
		}

		foundOlder := false
		for _, item := range response {
			// Skip pull requests
			if isPullRequest(item) {
				continue
			}

			issue, err := parseIssue(item)
			if err != nil {
				continue // Skip malformed issues
			}

			// Check if issue was closed in the period
			if issue.ClosedAt == nil {
				continue
			}

			if issue.ClosedAt.Before(fromTime) {
				foundOlder = true
				continue
			}

			if issue.ClosedAt.After(toTime) {
				continue
			}

			// Filter bots if needed
			if c.excludeBots && c.isBot(issue.Author.Login) {
				continue
			}

			issues = append(issues, issue)
		}

		// If we found issues older than our period, we can stop
		if foundOlder || len(response) < perPage {
			break
		}
		page++
	}

	return issues, nil
}

// isPullRequest checks if an issue is actually a pull request
func isPullRequest(issueData map[string]interface{}) bool {
	_, hasPR := issueData["pull_request"]
	return hasPR
}

// parseIssue converts GitHub API response to types.Issue
func parseIssue(data map[string]interface{}) (types.Issue, error) {
	issue := types.Issue{}

	// Parse number
	if num, ok := data["number"].(float64); ok {
		issue.Number = int(num)
	} else {
		return issue, fmt.Errorf("invalid issue number")
	}

	// Parse title
	if title, ok := data["title"].(string); ok {
		issue.Title = title
	}

	// Parse body
	if body, ok := data["body"].(string); ok {
		issue.Body = body
	}

	// Parse author
	if user, ok := data["user"].(map[string]interface{}); ok {
		if login, ok := user["login"].(string); ok {
			issue.Author = types.Author{
				Login:      login,
				ProfileURL: fmt.Sprintf("https://github.com/%s", login),
			}
			if name, ok := user["name"].(string); ok {
				issue.Author.Name = name
			}
		}
	}

	// Parse state
	if state, ok := data["state"].(string); ok {
		issue.State = state
	}

	// Parse created_at
	if createdStr, ok := data["created_at"].(string); ok {
		if created, err := time.Parse(time.RFC3339, createdStr); err == nil {
			issue.CreatedAt = created
		}
	}

	// Parse closed_at (nullable)
	if closedStr, ok := data["closed_at"].(string); ok && closedStr != "" {
		if closed, err := time.Parse(time.RFC3339, closedStr); err == nil {
			issue.ClosedAt = &closed
		}
	}

	// Parse labels
	if labelsData, ok := data["labels"].([]interface{}); ok {
		for _, labelItem := range labelsData {
			if labelMap, ok := labelItem.(map[string]interface{}); ok {
				if labelName, ok := labelMap["name"].(string); ok {
					issue.Labels = append(issue.Labels, labelName)
				}
			}
		}
	}

	// Parse assignees
	if assigneesData, ok := data["assignees"].([]interface{}); ok {
		for _, assigneeItem := range assigneesData {
			if assigneeMap, ok := assigneeItem.(map[string]interface{}); ok {
				if login, ok := assigneeMap["login"].(string); ok {
					assignee := types.Author{
						Login:      login,
						ProfileURL: fmt.Sprintf("https://github.com/%s", login),
					}
					if name, ok := assigneeMap["name"].(string); ok {
						assignee.Name = name
					}
					issue.Assignees = append(issue.Assignees, assignee)
				}
			}
		}
	}

	// Parse URL
	if url, ok := data["html_url"].(string); ok {
		issue.URL = url
	}

	return issue, nil
}
