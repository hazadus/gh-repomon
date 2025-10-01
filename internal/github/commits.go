package github

import (
	"fmt"
	"sync"
	"time"

	"github.com/hazadus/gh-repomon/internal/types"
)

// commitResponse represents the GitHub API response for a commit
type commitResponse struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
	Author struct {
		Login   string `json:"login"`
		HTMLURL string `json:"html_url"`
		Type    string `json:"type"`
	} `json:"author"`
	HTMLURL string `json:"html_url"`
	Stats   struct {
		Additions int `json:"additions"`
		Deletions int `json:"deletions"`
	} `json:"stats"`
}

// GetCommits retrieves commits from a repository for the specified period
func (c *Client) GetCommits(repo, branch string, from, to time.Time) ([]types.Commit, error) {
	// Build API path with query parameters
	path := fmt.Sprintf("repos/%s/commits?since=%s&until=%s",
		repo,
		from.Format(time.RFC3339),
		to.Format(time.RFC3339))

	// Add branch parameter if specified
	if branch != "" {
		path += fmt.Sprintf("&sha=%s", branch)
	}

	// Make API request with retry
	var response []commitResponse
	err := c.doWithRetry("GET", path, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	// Convert to types.Commit with parallel stats fetching
	commits := make([]types.Commit, 0, len(response))
	commitsMutex := sync.Mutex{}

	// Worker pool for fetching commit stats
	maxWorkers := 10
	if len(response) < maxWorkers {
		maxWorkers = len(response)
	}

	commitChan := make(chan commitResponse, len(response))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for cr := range commitChan {
				// Check if we should filter out bots
				if c.excludeBots && cr.Author.Login != "" && c.isBot(cr.Author.Login) {
					continue
				}

				// Get detailed commit stats
				additions, deletions, err := c.GetCommitStats(repo, cr.SHA)
				if err != nil {
					// If we can't get stats, use zeros but don't fail
					additions = 0
					deletions = 0
				}

				// Create Author
				author := types.Author{
					Login:      cr.Author.Login,
					Name:       cr.Commit.Author.Name,
					ProfileURL: cr.Author.HTMLURL,
					IsBot:      c.isBot(cr.Author.Login),
				}

				// If author.Login is empty (deleted user), use name
				if author.Login == "" {
					author.Login = cr.Commit.Author.Name
				}

				commit := types.Commit{
					SHA:       cr.SHA,
					Message:   cr.Commit.Message,
					Author:    author,
					Date:      cr.Commit.Author.Date,
					Additions: additions,
					Deletions: deletions,
					URL:       cr.HTMLURL,
				}

				commitsMutex.Lock()
				commits = append(commits, commit)
				commitsMutex.Unlock()
			}
		}()
	}

	// Send commits to workers
	for _, cr := range response {
		commitChan <- cr
	}
	close(commitChan)

	// Wait for all workers to finish
	wg.Wait()

	return commits, nil
}

// GetCommitStats retrieves detailed statistics for a specific commit
func (c *Client) GetCommitStats(repo, sha string) (additions, deletions int, err error) {
	// Build API path
	path := fmt.Sprintf("repos/%s/commits/%s", repo, sha)

	// Make API request with retry
	var response commitResponse
	err = c.doWithRetry("GET", path, nil, &response)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get commit stats: %w", err)
	}

	return response.Stats.Additions, response.Stats.Deletions, nil
}
