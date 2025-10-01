package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/hazadus/gh-repomon/internal/github"
	"github.com/hazadus/gh-repomon/internal/llm"
	"github.com/hazadus/gh-repomon/internal/logger"
	"github.com/hazadus/gh-repomon/internal/types"
)

// Generator generates reports based on GitHub activity data
type Generator struct {
	githubClient *github.Client
	llmClient    *llm.Client
	logger       *logger.Logger
}

// GenerationStats holds statistics about the report generation process
type GenerationStats struct {
	TotalBranches       int
	TotalAISummaries    int
	SuccessfulSummaries int
	FailedSummaries     int
}

// Options contains configuration for report generation
type Options struct {
	// Repository is the repository name (owner/repo)
	Repository string
	// Period is the time period to analyze
	Period types.Period
	// User is an optional filter by user login
	User string
	// Model is the AI model to use for summaries
	Model string
	// Language is the output language for AI summaries
	Language string
}

// NewGenerator creates a new report generator
func NewGenerator(ghClient *github.Client, llmClient *llm.Client) *Generator {
	return &Generator{
		githubClient: ghClient,
		llmClient:    llmClient,
		logger:       logger.New(),
	}
}

// Generate generates a full report based on the provided options
func (g *Generator) Generate(opts Options) (string, error) {
	// Initialize statistics
	stats := &GenerationStats{}

	// Collect data from GitHub
	data, err := g.collectData(opts)
	if err != nil {
		return "", err
	}

	stats.TotalBranches = len(data.Branches)

	// Generate AI summary if LLM client is available
	var overallSummary string
	if g.llmClient != nil {
		g.logger.Info("Generating AI summaries...")
		summary, err := g.llmClient.GenerateOverallSummary(data, opts.Language, opts.Model)
		if err != nil {
			g.logger.Warning(fmt.Sprintf("Failed to generate overall summary: %v", err))
			overallSummary = "Summary generation failed. Please check the activity details below."
			stats.FailedSummaries++
		} else {
			overallSummary = summary
			g.logger.Success("Overall summary generated")
			stats.SuccessfulSummaries++
		}
		stats.TotalAISummaries++

		// Generate branch summaries
		for i := range data.Branches {
			branchSummary, err := g.llmClient.GenerateBranchSummary(&data.Branches[i], opts.Language, opts.Model)
			if err != nil {
				g.logger.Warning(fmt.Sprintf("Failed to generate summary for branch %s: %v", data.Branches[i].Name, err))
				data.Branches[i].AISummary = fmt.Sprintf("Development activity in branch %s", data.Branches[i].Name)
				stats.FailedSummaries++
			} else {
				data.Branches[i].AISummary = branchSummary
				stats.SuccessfulSummaries++
			}
			stats.TotalAISummaries++
		}
		g.logger.Success(fmt.Sprintf("Branch summaries generated (%d/%d)", stats.TotalBranches-stats.FailedSummaries, stats.TotalBranches))

		// Generate PR summaries for open PRs
		totalPRs := len(data.OpenPRs) + len(data.UpdatedPRs)
		prSuccessCount := 0
		for i := range data.OpenPRs {
			prSummary, err := g.llmClient.GeneratePRSummary(&data.OpenPRs[i], opts.Language, opts.Model)
			if err != nil {
				g.logger.Warning(fmt.Sprintf("Failed to generate summary for PR #%d: %v", data.OpenPRs[i].Number, err))
				data.OpenPRs[i].AISummary = fmt.Sprintf("Pull request: %s", data.OpenPRs[i].Title)
				stats.FailedSummaries++
			} else {
				data.OpenPRs[i].AISummary = prSummary
				prSuccessCount++
				stats.SuccessfulSummaries++
			}
			stats.TotalAISummaries++
		}

		// Generate PR summaries for updated PRs
		for i := range data.UpdatedPRs {
			prSummary, err := g.llmClient.GeneratePRSummary(&data.UpdatedPRs[i], opts.Language, opts.Model)
			if err != nil {
				g.logger.Warning(fmt.Sprintf("Failed to generate summary for PR #%d: %v", data.UpdatedPRs[i].Number, err))
				data.UpdatedPRs[i].AISummary = fmt.Sprintf("Pull request: %s", data.UpdatedPRs[i].Title)
				stats.FailedSummaries++
			} else {
				data.UpdatedPRs[i].AISummary = prSummary
				prSuccessCount++
				stats.SuccessfulSummaries++
			}
			stats.TotalAISummaries++
		}
		g.logger.Success(fmt.Sprintf("PR summaries generated (%d/%d)", prSuccessCount, totalPRs))
	} else {
		overallSummary = "[AI summary generation disabled]"
	}

	// Generate markdown report
	report := g.generateMarkdown(data, overallSummary, stats)

	return report, nil
}

// collectData collects all necessary data from GitHub API
func (g *Generator) collectData(opts Options) (*types.ReportData, error) {
	// Convert times to ISO8601 format for API calls
	fromISO := opts.Period.From.Format(time.RFC3339)
	toISO := opts.Period.To.Format(time.RFC3339)

	// Get active branches
	g.logger.Progress("Collecting branches...")
	branches, err := g.githubClient.GetActiveBranches(opts.Repository, opts.Period.From, opts.Period.To)
	if err != nil {
		return nil, err
	}
	g.logger.Success(fmt.Sprintf("Found %d active branches", len(branches)))

	// Count total commits
	totalCommits := 0
	for _, branch := range branches {
		totalCommits += len(branch.Commits)
	}
	g.logger.Success(fmt.Sprintf("Collected %d commits", totalCommits))

	// Get open pull requests
	g.logger.Progress("Collecting pull requests...")
	openPRs, err := g.githubClient.GetOpenPullRequests(opts.Repository)
	if err != nil {
		return nil, err
	}

	// Get updated pull requests
	updatedPRs, err := g.githubClient.GetUpdatedPullRequests(opts.Repository, fromISO, toISO)
	if err != nil {
		return nil, err
	}
	g.logger.Success(fmt.Sprintf("Found %d open PRs, %d updated PRs", len(openPRs), len(updatedPRs)))

	// Get open issues
	g.logger.Progress("Collecting issues...")
	openIssues, err := g.githubClient.GetOpenIssues(opts.Repository)
	if err != nil {
		return nil, err
	}

	// Get closed issues
	closedIssues, err := g.githubClient.GetClosedIssues(opts.Repository, fromISO, toISO)
	if err != nil {
		return nil, err
	}
	g.logger.Success(fmt.Sprintf("Found %d open issues, %d closed issues", len(openIssues), len(closedIssues)))

	// Create report data
	data := &types.ReportData{
		Repository:    opts.Repository,
		RepositoryURL: "https://github.com/" + opts.Repository,
		Period:        opts.Period,
		GeneratedAt:   time.Now().UTC(),
		Branches:      branches,
		OpenPRs:       openPRs,
		UpdatedPRs:    updatedPRs,
		OpenIssues:    openIssues,
		ClosedIssues:  closedIssues,
	}

	return data, nil
}

// generateMarkdown generates a markdown report from collected data
func (g *Generator) generateMarkdown(data *types.ReportData, overallSummary string, stats *GenerationStats) string {
	var sb strings.Builder

	// Calculate overall statistics
	data.OverallStats = calculateOverallStats(data)

	// Calculate author statistics
	data.AuthorStats = calculateAuthorStats(data)

	// Generate header
	sb.WriteString(generateHeader(data))

	// Generate summary statistics
	sb.WriteString(generateSummaryStats(data.OverallStats))

	// Overall AI Summary
	sb.WriteString("## 📊 Overall Summary\n\n")
	sb.WriteString(overallSummary)
	sb.WriteString("\n\n")

	// Generate branches section
	if len(data.Branches) > 0 {
		sb.WriteString(generateBranchesSection(data.Branches))
	}

	// Generate pull requests sections
	sb.WriteString(generateOpenPRsSection(data.OpenPRs))
	sb.WriteString(generateUpdatedPRsSection(data.UpdatedPRs))

	// Generate issues sections
	sb.WriteString(generateOpenIssuesSection(data.OpenIssues))
	sb.WriteString(generateClosedIssuesSection(data.ClosedIssues))

	// Generate code reviews section
	sb.WriteString(generateCodeReviewsSection(data))

	// Generate author statistics section
	sb.WriteString(generateAuthorStatsSection(data.AuthorStats))

	// Generate footer with statistics
	sb.WriteString(generateFooter(stats))

	return sb.String()
}

// generateFooter generates the report footer with generation statistics
func generateFooter(stats *GenerationStats) string {
	var sb strings.Builder

	sb.WriteString("\n---\n\n")
	sb.WriteString("*Generated by [gh-repomon](https://github.com/hazadus/gh-repomon)*\n\n")

	if stats != nil && stats.TotalAISummaries > 0 {
		sb.WriteString(fmt.Sprintf("*AI Summaries: %d successful, %d failed (total %d)*\n",
			stats.SuccessfulSummaries,
			stats.FailedSummaries,
			stats.TotalAISummaries))
	}

	return sb.String()
}
