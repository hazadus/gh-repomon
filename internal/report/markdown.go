package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/hazadus/gh-repomon/internal/types"
)

// generateHeader generates the report header with repository information
func generateHeader(data *types.ReportData) string {
	var sb strings.Builder

	// Repository name as main heading
	repoName := strings.Split(data.Repository, "/")
	displayName := data.Repository
	if len(repoName) == 2 {
		displayName = repoName[1]
	}

	sb.WriteString(fmt.Sprintf("# Repository Activity Report: %s\n\n", displayName))
	sb.WriteString(fmt.Sprintf("**Repository**: [%s](%s)\n\n", data.Repository, data.RepositoryURL))
	sb.WriteString(fmt.Sprintf("**Period**: %s to %s\n\n",
		data.Period.From.Format("2006-01-02"),
		data.Period.To.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Report Generated**: %s UTC\n\n",
		data.GeneratedAt.Format("2006-01-02 15:04:05")))

	return sb.String()
}

// calculateOverallStats calculates overall statistics from collected data
func calculateOverallStats(data *types.ReportData) types.OverallStats {
	stats := types.OverallStats{
		OpenPRCount:       len(data.OpenPRs),
		OpenIssuesCount:   len(data.OpenIssues),
		ClosedIssuesCount: len(data.ClosedIssues),
	}

	// Count total commits and collect unique authors
	authorSet := make(map[string]bool)
	for _, branch := range data.Branches {
		stats.TotalCommits += len(branch.Commits)
		for _, commit := range branch.Commits {
			authorSet[commit.Author.Login] = true
		}
	}
	stats.TotalAuthors = len(authorSet)

	// Count reviews (will be implemented when review functionality is added)
	stats.ReviewsCount = 0
	for _, pr := range data.OpenPRs {
		stats.ReviewsCount += pr.Reviews
	}
	for _, pr := range data.UpdatedPRs {
		// Avoid double counting PRs that are both open and updated
		isOpen := false
		for _, openPR := range data.OpenPRs {
			if openPR.Number == pr.Number {
				isOpen = true
				break
			}
		}
		if !isOpen {
			stats.ReviewsCount += pr.Reviews
		}
	}

	return stats
}

// generateSummaryStats generates the summary statistics section
func generateSummaryStats(stats types.OverallStats) string {
	var sb strings.Builder

	sb.WriteString("## Summary Statistics\n\n")
	sb.WriteString(fmt.Sprintf("- **Total Commits**: %d\n", stats.TotalCommits))
	sb.WriteString(fmt.Sprintf("- **Total Authors**: %d\n", stats.TotalAuthors))
	sb.WriteString(fmt.Sprintf("- **Open Pull Requests**: %d\n", stats.OpenPRCount))
	sb.WriteString(fmt.Sprintf("- **Open Issues**: %d\n", stats.OpenIssuesCount))
	sb.WriteString(fmt.Sprintf("- **Closed Issues**: %d\n", stats.ClosedIssuesCount))
	sb.WriteString(fmt.Sprintf("- **Code Reviews**: %d\n\n", stats.ReviewsCount))

	return sb.String()
}

// formatDate formats a time.Time to a readable string
func formatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

// formatAuthorLinks formats a list of author logins as markdown links
func formatAuthorLinks(authors []string) string {
	if len(authors) == 0 {
		return "none"
	}

	links := make([]string, len(authors))
	for i, author := range authors {
		links[i] = fmt.Sprintf("[%s](https://github.com/%s)", author, author)
	}
	return strings.Join(links, ", ")
}

// formatCommitMessage splits a commit message into short and full versions
func formatCommitMessage(message string) (short, full string) {
	lines := strings.Split(message, "\n")
	short = lines[0]

	// Truncate short message if too long
	if len(short) > 72 {
		short = short[:69] + "..."
	}

	full = message
	return
}

// generateBranchSection generates a detailed section for a single branch
func generateBranchSection(branch types.Branch) string {
	var sb strings.Builder

	// Branch header
	sb.WriteString(fmt.Sprintf("## 🌿 Branch: %s\n\n", branch.Name))

	// AI Summary
	if branch.AISummary != "" {
		sb.WriteString("### AI Summary\n\n")
		sb.WriteString(branch.AISummary)
		sb.WriteString("\n\n")
	}

	// Statistics subsection
	sb.WriteString("### Statistics\n\n")
	sb.WriteString(fmt.Sprintf("- **Total Commits**: %d\n", len(branch.Commits)))
	sb.WriteString(fmt.Sprintf("- **Lines Added**: +%d\n", branch.TotalAdded))
	sb.WriteString(fmt.Sprintf("- **Lines Deleted**: -%d\n", branch.TotalDeleted))
	sb.WriteString(fmt.Sprintf("- **Contributors**: %s\n\n", formatAuthorLinks(branch.Authors)))

	// Commits subsection
	sb.WriteString("### Commits\n\n")

	for _, commit := range branch.Commits {
		short, full := formatCommitMessage(commit.Message)

		sb.WriteString(fmt.Sprintf("#### [%s](%s)\n\n", short, commit.URL))
		sb.WriteString(fmt.Sprintf("**Author**: [%s](%s) | **Date**: %s\n\n",
			commit.Author.Login,
			commit.Author.ProfileURL,
			formatDate(commit.Date)))
		sb.WriteString(fmt.Sprintf("**Changes**: +%d / -%d lines\n\n",
			commit.Additions,
			commit.Deletions))

		// Add full message if it's multiline
		if strings.Contains(full, "\n") && full != short {
			// Escape any special markdown characters in the full message
			sb.WriteString("```\n")
			sb.WriteString(full)
			sb.WriteString("\n```\n\n")
		}

		sb.WriteString("---\n\n")
	}

	return sb.String()
}

// generateBranchesSection generates sections for all branches
func generateBranchesSection(branches []types.Branch) string {
	var sb strings.Builder

	for _, branch := range branches {
		sb.WriteString(generateBranchSection(branch))
	}

	return sb.String()
}

// generatePRSection generates a section for a single pull request
func generatePRSection(pr types.PullRequest) string {
	var sb strings.Builder

	// PR header with link
	sb.WriteString(fmt.Sprintf("### [PR #%d: %s](%s)\n\n", pr.Number, pr.Title, pr.URL))

	// Metadata
	sb.WriteString(fmt.Sprintf("- **Author**: [%s](%s)\n", pr.Author.Login, pr.Author.ProfileURL))
	sb.WriteString(fmt.Sprintf("- **Created**: %s\n", pr.CreatedAt.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("- **Status**: %s\n", pr.State))
	sb.WriteString(fmt.Sprintf("- **Comments**: %d\n", pr.Comments))
	sb.WriteString(fmt.Sprintf("- **Reviews**: %d\n\n", pr.Reviews))

	// AI Summary placeholder
	sb.WriteString("#### AI Summary\n\n")
	sb.WriteString("[Brief AI summary will be here]\n\n")

	sb.WriteString("---\n\n")

	return sb.String()
}

// generateOpenPRsSection generates the section for open pull requests
func generateOpenPRsSection(prs []types.PullRequest) string {
	var sb strings.Builder

	sb.WriteString("## 🔀 Open Pull Requests\n\n")

	if len(prs) == 0 {
		sb.WriteString("No open pull requests\n\n")
		return sb.String()
	}

	for _, pr := range prs {
		sb.WriteString(generatePRSection(pr))
	}

	return sb.String()
}

// generateUpdatedPRsSection generates the section for updated pull requests
func generateUpdatedPRsSection(prs []types.PullRequest) string {
	var sb strings.Builder

	sb.WriteString("## 🔄 Updated Pull Requests\n\n")

	if len(prs) == 0 {
		sb.WriteString("No pull requests were updated during this period\n\n")
		return sb.String()
	}

	for _, pr := range prs {
		sb.WriteString(generatePRSection(pr))
	}

	return sb.String()
}

// generateIssueSection generates a section for a single issue
func generateIssueSection(issue types.Issue) string {
	var sb strings.Builder

	// Issue header with link
	sb.WriteString(fmt.Sprintf("### [Issue #%d: %s](%s)\n\n", issue.Number, issue.Title, issue.URL))

	// Metadata
	sb.WriteString(fmt.Sprintf("- **Author**: [%s](%s)\n", issue.Author.Login, issue.Author.ProfileURL))
	sb.WriteString(fmt.Sprintf("- **Created**: %s\n", issue.CreatedAt.Format("2006-01-02")))

	// Labels (if any)
	if len(issue.Labels) > 0 {
		sb.WriteString(fmt.Sprintf("- **Labels**: %s\n", strings.Join(issue.Labels, ", ")))
	}

	// Assignees (if any)
	if len(issue.Assignees) > 0 {
		assigneeLogins := make([]string, len(issue.Assignees))
		for i, assignee := range issue.Assignees {
			assigneeLogins[i] = assignee.Login
		}
		sb.WriteString(fmt.Sprintf("- **Assignees**: %s\n", formatAuthorLinks(assigneeLogins)))
	}

	sb.WriteString("\n---\n\n")

	return sb.String()
}

// generateOpenIssuesSection generates the section for open issues
func generateOpenIssuesSection(issues []types.Issue) string {
	var sb strings.Builder

	sb.WriteString("## 📋 Open Issues\n\n")

	if len(issues) == 0 {
		sb.WriteString("No open issues\n\n")
		return sb.String()
	}

	for _, issue := range issues {
		sb.WriteString(generateIssueSection(issue))
	}

	return sb.String()
}

// generateClosedIssuesSection generates the section for closed issues
func generateClosedIssuesSection(issues []types.Issue) string {
	var sb strings.Builder

	sb.WriteString("## ✅ Closed Issues\n\n")

	if len(issues) == 0 {
		sb.WriteString("No issues were closed during this period\n\n")
		return sb.String()
	}

	for _, issue := range issues {
		sb.WriteString(generateIssueSection(issue))
	}

	return sb.String()
}

// calculateAuthorStats calculates statistics per author from collected data
func calculateAuthorStats(data *types.ReportData) []types.AuthorStats {
	// Map to accumulate statistics by author login
	authorMap := make(map[string]*types.AuthorStats)

	// Process all commits in all branches
	for _, branch := range data.Branches {
		for _, commit := range branch.Commits {
			login := commit.Author.Login

			// Initialize author stats if not exists
			if _, exists := authorMap[login]; !exists {
				authorMap[login] = &types.AuthorStats{
					Author:         commit.Author,
					BranchActivity: make(map[string]types.BranchActivity),
				}
			}

			stats := authorMap[login]

			// Update overall stats
			stats.TotalCommits++
			stats.TotalAdded += commit.Additions
			stats.TotalDeleted += commit.Deletions

			// Update branch activity
			branchActivity := stats.BranchActivity[branch.Name]
			branchActivity.Commits++
			branchActivity.Added += commit.Additions
			branchActivity.Deleted += commit.Deletions
			stats.BranchActivity[branch.Name] = branchActivity
		}
	}

	// Process PRs
	for _, pr := range data.OpenPRs {
		login := pr.Author.Login
		if stats, exists := authorMap[login]; exists {
			stats.PRsCreated++
		} else {
			authorMap[login] = &types.AuthorStats{
				Author:         pr.Author,
				PRsCreated:     1,
				BranchActivity: make(map[string]types.BranchActivity),
			}
		}
	}

	// Process updated PRs (avoid double counting open PRs)
	for _, pr := range data.UpdatedPRs {
		isOpen := false
		for _, openPR := range data.OpenPRs {
			if openPR.Number == pr.Number {
				isOpen = true
				break
			}
		}
		if !isOpen {
			login := pr.Author.Login
			if stats, exists := authorMap[login]; exists {
				stats.PRsCreated++
			} else {
				authorMap[login] = &types.AuthorStats{
					Author:         pr.Author,
					PRsCreated:     1,
					BranchActivity: make(map[string]types.BranchActivity),
				}
			}
		}
	}

	// Process Issues
	for _, issue := range data.OpenIssues {
		login := issue.Author.Login
		if stats, exists := authorMap[login]; exists {
			stats.IssuesCreated++
		} else {
			authorMap[login] = &types.AuthorStats{
				Author:         issue.Author,
				IssuesCreated:  1,
				BranchActivity: make(map[string]types.BranchActivity),
			}
		}
	}

	for _, issue := range data.ClosedIssues {
		login := issue.Author.Login
		if stats, exists := authorMap[login]; exists {
			stats.IssuesCreated++
		} else {
			authorMap[login] = &types.AuthorStats{
				Author:         issue.Author,
				IssuesCreated:  1,
				BranchActivity: make(map[string]types.BranchActivity),
			}
		}
	}

	// Note: ReviewsCount will be populated when review functionality is added
	// In a real implementation, we would query actual reviewers from GitHub API
	// For now, this is just structural preparation for future enhancement

	// Convert map to slice and sort by total commits (descending)
	result := make([]types.AuthorStats, 0, len(authorMap))
	for _, stats := range authorMap {
		result = append(result, *stats)
	}

	// Sort by total commits (descending)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].TotalCommits > result[i].TotalCommits {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

// generateCodeReviewsSection generates the code reviews section
func generateCodeReviewsSection(data *types.ReportData) string {
	var sb strings.Builder

	sb.WriteString("## 👀 Code Reviews\n\n")

	// Collect PRs with reviews
	prsWithReviews := []types.PullRequest{}
	for _, pr := range data.OpenPRs {
		if pr.Reviews > 0 {
			prsWithReviews = append(prsWithReviews, pr)
		}
	}
	for _, pr := range data.UpdatedPRs {
		// Check if not already counted in open PRs
		isOpen := false
		for _, openPR := range data.OpenPRs {
			if openPR.Number == pr.Number {
				isOpen = true
				break
			}
		}
		if !isOpen && pr.Reviews > 0 {
			prsWithReviews = append(prsWithReviews, pr)
		}
	}

	if len(prsWithReviews) == 0 {
		sb.WriteString("No code reviews found during this period\n\n")
		return sb.String()
	}

	sb.WriteString("### Pull Requests Reviewed\n\n")

	totalReviews := 0
	for _, pr := range prsWithReviews {
		sb.WriteString(fmt.Sprintf("- [PR #%d: %s](%s) - %d reviews\n",
			pr.Number, pr.Title, pr.URL, pr.Reviews))
		totalReviews += pr.Reviews
	}

	sb.WriteString(fmt.Sprintf("\n**Total Reviews**: %d\n\n", totalReviews))

	return sb.String()
}

// generateAuthorSection generates a detailed section for a single author
func generateAuthorSection(stats types.AuthorStats) string {
	var sb strings.Builder

	// Author header
	sb.WriteString(fmt.Sprintf("### [%s](%s)\n\n", stats.Author.Login, stats.Author.ProfileURL))

	// Overall statistics
	sb.WriteString("#### Overall Statistics\n\n")
	sb.WriteString(fmt.Sprintf("- **Total Commits**: %d\n", stats.TotalCommits))
	sb.WriteString(fmt.Sprintf("- **Total Lines Added**: +%d\n", stats.TotalAdded))
	sb.WriteString(fmt.Sprintf("- **Total Lines Deleted**: -%d\n", stats.TotalDeleted))
	sb.WriteString(fmt.Sprintf("- **Pull Requests Created**: %d\n", stats.PRsCreated))
	sb.WriteString(fmt.Sprintf("- **Issues Created**: %d\n", stats.IssuesCreated))
	sb.WriteString(fmt.Sprintf("- **Code Reviews**: %d\n\n", stats.ReviewsCount))

	// Activity by branch
	if len(stats.BranchActivity) > 0 {
		sb.WriteString("#### Activity by Branch\n\n")

		for branchName, activity := range stats.BranchActivity {
			sb.WriteString(fmt.Sprintf("##### Branch: %s\n\n", branchName))
			sb.WriteString(fmt.Sprintf("- **Commits**: %d\n", activity.Commits))
			sb.WriteString(fmt.Sprintf("- **Lines Added**: +%d\n", activity.Added))
			sb.WriteString(fmt.Sprintf("- **Lines Deleted**: -%d\n\n", activity.Deleted))
		}
	}

	sb.WriteString("---\n\n")

	return sb.String()
}

// generateAuthorStatsSection generates the author activity section
func generateAuthorStatsSection(authorStats []types.AuthorStats) string {
	var sb strings.Builder

	sb.WriteString("## 👥 Author Activity\n\n")

	if len(authorStats) == 0 {
		sb.WriteString("No author activity found during this period\n\n")
		return sb.String()
	}

	for _, stats := range authorStats {
		sb.WriteString(generateAuthorSection(stats))
	}

	return sb.String()
}
