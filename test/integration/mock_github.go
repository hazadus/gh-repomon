// +build integration

package integration

import (
	"time"

	"github.com/hazadus/gh-repomon/internal/types"
)

// MockGitHubClient is a mock implementation of GitHub API client for testing
type MockGitHubClient struct {
	activeBranches   []types.Branch
	openPRs          []types.PullRequest
	updatedPRs       []types.PullRequest
	openIssues       []types.Issue
	closedIssues     []types.Issue
	reviewsByAuthor  map[string]int
	totalReviewCount int
}

// NewMockGitHubClient creates a new mock GitHub client with predefined test data
func NewMockGitHubClient() *MockGitHubClient {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	return &MockGitHubClient{
		activeBranches: []types.Branch{
			{
				Name: "main",
				Commits: []types.Commit{
					{
						SHA:       "abc123",
						Message:   "feat: add new feature\n\nDetailed description of the feature",
						Author:    types.Author{Login: "developer1", Name: "Developer One", ProfileURL: "https://github.com/developer1", IsBot: false},
						Date:      yesterday,
						Additions: 100,
						Deletions: 20,
						URL:       "https://github.com/owner/repo/commit/abc123",
					},
					{
						SHA:       "def456",
						Message:   "fix: resolve bug in authentication",
						Author:    types.Author{Login: "developer2", Name: "Developer Two", ProfileURL: "https://github.com/developer2", IsBot: false},
						Date:      now,
						Additions: 50,
						Deletions: 10,
						URL:       "https://github.com/owner/repo/commit/def456",
					},
				},
				TotalAdded:   150,
				TotalDeleted: 30,
				Authors:      []string{"developer1", "developer2"},
			},
			{
				Name: "feature/new-ui",
				Commits: []types.Commit{
					{
						SHA:       "ghi789",
						Message:   "feat: redesign user interface",
						Author:    types.Author{Login: "developer1", Name: "Developer One", ProfileURL: "https://github.com/developer1", IsBot: false},
						Date:      yesterday,
						Additions: 200,
						Deletions: 50,
						URL:       "https://github.com/owner/repo/commit/ghi789",
					},
				},
				TotalAdded:   200,
				TotalDeleted: 50,
				Authors:      []string{"developer1"},
			},
		},
		openPRs: []types.PullRequest{
			{
				Number:    1,
				Title:     "Add authentication feature",
				Body:      "This PR adds OAuth authentication support",
				Author:    types.Author{Login: "developer1", Name: "Developer One", ProfileURL: "https://github.com/developer1", IsBot: false},
				State:     "open",
				CreatedAt: yesterday,
				UpdatedAt: now,
				Comments:  3,
				Reviews:   2,
				URL:       "https://github.com/owner/repo/pull/1",
			},
		},
		updatedPRs: []types.PullRequest{
			{
				Number:    2,
				Title:     "Refactor database layer",
				Body:      "This PR refactors the database access layer",
				Author:    types.Author{Login: "developer2", Name: "Developer Two", ProfileURL: "https://github.com/developer2", IsBot: false},
				State:     "open",
				CreatedAt: now.Add(-48 * time.Hour),
				UpdatedAt: yesterday,
				Comments:  5,
				Reviews:   3,
				URL:       "https://github.com/owner/repo/pull/2",
			},
		},
		openIssues: []types.Issue{
			{
				Number:    10,
				Title:     "Bug: Login fails on mobile",
				Body:      "Users report login failures on mobile devices",
				Author:    types.Author{Login: "user1", Name: "User One", ProfileURL: "https://github.com/user1", IsBot: false},
				State:     "open",
				CreatedAt: yesterday,
				Labels:    []string{"bug", "mobile"},
				Assignees: []types.Author{
					{Login: "developer1", Name: "Developer One", ProfileURL: "https://github.com/developer1", IsBot: false},
				},
				URL: "https://github.com/owner/repo/issues/10",
			},
		},
		closedIssues: []types.Issue{
			{
				Number:    9,
				Title:     "Feature request: Dark mode",
				Body:      "Add dark mode support to the application",
				Author:    types.Author{Login: "user2", Name: "User Two", ProfileURL: "https://github.com/user2", IsBot: false},
				State:     "closed",
				CreatedAt: now.Add(-72 * time.Hour),
				ClosedAt:  &yesterday,
				Labels:    []string{"enhancement"},
				Assignees: []types.Author{
					{Login: "developer2", Name: "Developer Two", ProfileURL: "https://github.com/developer2", IsBot: false},
				},
				URL: "https://github.com/owner/repo/issues/9",
			},
		},
		reviewsByAuthor: map[string]int{
			"developer1": 3,
			"developer2": 2,
		},
		totalReviewCount: 5,
	}
}

// GetActiveBranches returns mock active branches
func (m *MockGitHubClient) GetActiveBranches(repo string, from, to time.Time) ([]types.Branch, error) {
	return m.activeBranches, nil
}

// GetOpenPullRequests returns mock open PRs
func (m *MockGitHubClient) GetOpenPullRequests(repo string) ([]types.PullRequest, error) {
	return m.openPRs, nil
}

// GetUpdatedPullRequests returns mock updated PRs
func (m *MockGitHubClient) GetUpdatedPullRequests(repo, from, to string) ([]types.PullRequest, error) {
	return m.updatedPRs, nil
}

// GetOpenIssues returns mock open issues
func (m *MockGitHubClient) GetOpenIssues(repo string) ([]types.Issue, error) {
	return m.openIssues, nil
}

// GetClosedIssues returns mock closed issues
func (m *MockGitHubClient) GetClosedIssues(repo, from, to string) ([]types.Issue, error) {
	return m.closedIssues, nil
}

// GetReviewsByAuthor returns mock reviews by author
func (m *MockGitHubClient) GetReviewsByAuthor(repo string, prs []types.PullRequest) (map[string]int, error) {
	return m.reviewsByAuthor, nil
}

// GetAllReviews returns mock total review count
func (m *MockGitHubClient) GetAllReviews(repo string, prs []types.PullRequest) (int, error) {
	return m.totalReviewCount, nil
}
