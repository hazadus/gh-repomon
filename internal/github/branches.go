package github

import (
	"fmt"
	"sort"
	"time"

	"github.com/hazadus/gh-repomon/internal/types"
)

// branchResponse represents the GitHub API response for a branch
type branchResponse struct {
	Name   string `json:"name"`
	Commit struct {
		SHA string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
}

// GetBranches retrieves all branches from a repository
func (c *Client) GetBranches(repo string) ([]string, error) {
	// Build API path
	path := fmt.Sprintf("repos/%s/branches", repo)

	// Make API request
	var response []branchResponse
	err := c.client.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	// Extract branch names
	branches := make([]string, 0, len(response))
	for _, br := range response {
		branches = append(branches, br.Name)
	}

	return branches, nil
}

// GetActiveBranches retrieves branches that have commits during the specified period
func (c *Client) GetActiveBranches(repo string, from, to time.Time) ([]types.Branch, error) {
	// Get all branches
	branchNames, err := c.GetBranches(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	// Collect active branches
	activeBranches := make([]types.Branch, 0)

	for _, branchName := range branchNames {
		// Get commits for this branch during the period
		commits, err := c.GetCommits(repo, branchName, from, to)
		if err != nil {
			// If we can't get commits for a branch, skip it but don't fail
			continue
		}

		// Skip branches with no activity
		if len(commits) == 0 {
			continue
		}

		// Calculate statistics
		totalAdded := 0
		totalDeleted := 0
		for _, commit := range commits {
			totalAdded += commit.Additions
			totalDeleted += commit.Deletions
		}

		// Get unique authors
		authors := uniqueAuthors(commits)

		// Create branch object
		branch := types.Branch{
			Name:         branchName,
			Commits:      commits,
			PRs:          []types.PullRequest{}, // Will be populated later if needed
			TotalAdded:   totalAdded,
			TotalDeleted: totalDeleted,
			Authors:      authors,
		}

		activeBranches = append(activeBranches, branch)
	}

	return activeBranches, nil
}

// uniqueAuthors extracts unique author logins from commits and returns them sorted
func uniqueAuthors(commits []types.Commit) []string {
	// Use map to collect unique logins
	authorMap := make(map[string]bool)
	for _, commit := range commits {
		if commit.Author.Login != "" {
			authorMap[commit.Author.Login] = true
		}
	}

	// Convert to slice
	authors := make([]string, 0, len(authorMap))
	for login := range authorMap {
		authors = append(authors, login)
	}

	// Sort alphabetically
	sort.Strings(authors)

	return authors
}
