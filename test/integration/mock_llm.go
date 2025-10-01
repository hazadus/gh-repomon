// +build integration

package integration

import (
	"fmt"
	"sync"

	"github.com/hazadus/gh-repomon/internal/types"
)

// MockLLMClient is a mock implementation of LLM client for testing
type MockLLMClient struct {
	summaryCounter int
	mu             sync.Mutex
}

// NewMockLLMClient creates a new mock LLM client
func NewMockLLMClient() *MockLLMClient {
	return &MockLLMClient{
		summaryCounter: 0,
	}
}

// GenerateOverallSummary generates a mock overall summary
func (m *MockLLMClient) GenerateOverallSummary(data *types.ReportData, language, model string) (string, error) {
	m.mu.Lock()
	m.summaryCounter++
	m.mu.Unlock()
	return fmt.Sprintf("This is a comprehensive summary of repository activity for %s. "+
		"During the reporting period, the team made significant progress with %d commits across %d active branches. "+
		"The development focused on new features and bug fixes, with active collaboration evident from pull request reviews and issue discussions. "+
		"Overall, the repository shows healthy development activity with consistent contributions from multiple developers.",
		data.Repository, data.OverallStats.TotalCommits, len(data.Branches)), nil
}

// GenerateBranchSummary generates a mock branch summary
func (m *MockLLMClient) GenerateBranchSummary(branch *types.Branch, language, model string) (string, error) {
	m.mu.Lock()
	m.summaryCounter++
	m.mu.Unlock()

	if len(branch.Commits) == 0 {
		return fmt.Sprintf("Branch %s shows minimal activity during this period.", branch.Name), nil
	}

	firstCommit := branch.Commits[0]
	return fmt.Sprintf("The %s branch received %d commits focusing on %s. "+
		"Development work includes enhancements and improvements with %d lines added and %d lines removed.",
		branch.Name, len(branch.Commits), firstCommit.Message, branch.TotalAdded, branch.TotalDeleted), nil
}

// GeneratePRSummary generates a mock PR summary
func (m *MockLLMClient) GeneratePRSummary(pr *types.PullRequest, language, model string) (string, error) {
	m.mu.Lock()
	m.summaryCounter++
	m.mu.Unlock()

	return fmt.Sprintf("This pull request introduces changes related to: %s. "+
		"The PR has received %d reviews and %d comments, indicating active collaboration and code review.",
		pr.Title, pr.Reviews, pr.Comments), nil
}

// GetSummaryCount returns the number of summaries generated (for testing)
func (m *MockLLMClient) GetSummaryCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.summaryCounter
}
