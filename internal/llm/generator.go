package llm

import (
	"fmt"
	"strings"

	"github.com/hazadus/gh-repomon/internal/types"
)

// GenerateOverallSummary generates an AI summary of overall repository activity
func (c *Client) GenerateOverallSummary(data *types.ReportData, language, model string) (string, error) {
	// Load prompt
	config, err := LoadPrompt("overall_summary")
	if err != nil {
		return "Summary generation failed. Please check the activity details below.", fmt.Errorf("failed to load prompt: %w", err)
	}

	// Prepare variables
	vars := map[string]string{
		"language":      language,
		"repo_name":     data.Repository,
		"period":        formatPeriod(data.Period),
		"total_commits": fmt.Sprintf("%d", data.OverallStats.TotalCommits),
		"total_authors": fmt.Sprintf("%d", data.OverallStats.TotalAuthors),
		"branches":      formatBranchesForPrompt(data.Branches),
		"prs":           formatPRsForPrompt(data.OpenPRs, data.UpdatedPRs),
		"issues":        formatIssuesForPrompt(data.OpenIssues, data.ClosedIssues),
	}

	// Render prompt
	rendered, err := RenderPrompt(config, vars)
	if err != nil {
		return "Summary generation failed. Please check the activity details below.", fmt.Errorf("failed to render prompt: %w", err)
	}

	// Convert prompt messages to chat messages
	messages := make([]Message, len(rendered.Messages))
	for i, msg := range rendered.Messages {
		messages[i] = Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Create request
	request := ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: rendered.ModelParameters.Temperature,
	}

	// Send request
	response, err := c.Complete(request)
	if err != nil {
		return "Summary generation failed. Please check the activity details below.", fmt.Errorf("failed to complete request: %w", err)
	}

	return response, nil
}

// formatPeriod formats a period for display
func formatPeriod(period types.Period) string {
	return fmt.Sprintf("%s to %s", period.From.Format("2006-01-02"), period.To.Format("2006-01-02"))
}

// formatBranchesForPrompt formats branches for inclusion in prompt
func formatBranchesForPrompt(branches []types.Branch) string {
	if len(branches) == 0 {
		return "No active branches"
	}

	var parts []string
	for _, branch := range branches {
		authors := strings.Join(branch.Authors, ", ")
		parts = append(parts, fmt.Sprintf("- %s: %d commits by %s",
			branch.Name, len(branch.Commits), authors))
	}
	return strings.Join(parts, "\n")
}

// formatPRsForPrompt formats PRs for inclusion in prompt
func formatPRsForPrompt(openPRs, updatedPRs []types.PullRequest) string {
	var parts []string

	if len(openPRs) > 0 {
		parts = append(parts, "Open Pull Requests:")
		for _, pr := range openPRs {
			parts = append(parts, fmt.Sprintf("- #%d: %s (by %s)", pr.Number, pr.Title, pr.Author.Login))
		}
	}

	if len(updatedPRs) > 0 {
		if len(parts) > 0 {
			parts = append(parts, "")
		}
		parts = append(parts, "Updated Pull Requests:")
		for _, pr := range updatedPRs {
			parts = append(parts, fmt.Sprintf("- #%d: %s (by %s)", pr.Number, pr.Title, pr.Author.Login))
		}
	}

	if len(parts) == 0 {
		return "No pull requests"
	}

	return strings.Join(parts, "\n")
}

// formatIssuesForPrompt formats issues for inclusion in prompt
func formatIssuesForPrompt(openIssues, closedIssues []types.Issue) string {
	var parts []string

	if len(openIssues) > 0 {
		parts = append(parts, "Open Issues:")
		for _, issue := range openIssues {
			parts = append(parts, fmt.Sprintf("- #%d: %s (by %s)", issue.Number, issue.Title, issue.Author.Login))
		}
	}

	if len(closedIssues) > 0 {
		if len(parts) > 0 {
			parts = append(parts, "")
		}
		parts = append(parts, "Closed Issues:")
		for _, issue := range closedIssues {
			parts = append(parts, fmt.Sprintf("- #%d: %s (by %s)", issue.Number, issue.Title, issue.Author.Login))
		}
	}

	if len(parts) == 0 {
		return "No issues"
	}

	return strings.Join(parts, "\n")
}
