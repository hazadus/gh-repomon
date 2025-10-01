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

	// AI Summary placeholder
	sb.WriteString("### AI Summary\n\n")
	sb.WriteString("[AI summary will be generated here]\n\n")

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
