package report

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hazadus/gh-repomon/internal/github"
	"github.com/hazadus/gh-repomon/internal/types"
)

// Generator generates reports based on GitHub activity data
type Generator struct {
	githubClient *github.Client
	llmClient    interface{} // Will be implemented later
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
func NewGenerator(ghClient *github.Client) *Generator {
	return &Generator{
		githubClient: ghClient,
		llmClient:    nil, // Will be set later
	}
}

// Generate generates a full report based on the provided options
func (g *Generator) Generate(opts Options) (string, error) {
	// Collect data from GitHub
	data, err := g.collectData(opts)
	if err != nil {
		return "", err
	}

	// Generate markdown report
	report := g.generateMarkdown(data)

	return report, nil
}

// collectData collects all necessary data from GitHub API
func (g *Generator) collectData(opts Options) (*types.ReportData, error) {
	// Convert times to ISO8601 format for API calls
	fromISO := opts.Period.From.Format(time.RFC3339)
	toISO := opts.Period.To.Format(time.RFC3339)

	// Get active branches
	branches, err := g.githubClient.GetActiveBranches(opts.Repository, opts.Period.From, opts.Period.To)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "  🔍 Found %d active branches\n", len(branches))

	// Count total commits
	totalCommits := 0
	for _, branch := range branches {
		totalCommits += len(branch.Commits)
	}
	fmt.Fprintf(os.Stderr, "  🔍 Collecting commits... Found %d commits\n", totalCommits)

	// Get open pull requests
	fmt.Fprintf(os.Stderr, "  🔍 Collecting pull requests...\n")
	openPRs, err := g.githubClient.GetOpenPullRequests(opts.Repository)
	if err != nil {
		return nil, err
	}

	// Get updated pull requests
	updatedPRs, err := g.githubClient.GetUpdatedPullRequests(opts.Repository, fromISO, toISO)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "  🔍 Found %d open PRs, %d updated PRs\n", len(openPRs), len(updatedPRs))

	// Get open issues
	fmt.Fprintf(os.Stderr, "  🔍 Collecting issues...\n")
	openIssues, err := g.githubClient.GetOpenIssues(opts.Repository)
	if err != nil {
		return nil, err
	}

	// Get closed issues
	closedIssues, err := g.githubClient.GetClosedIssues(opts.Repository, fromISO, toISO)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "  🔍 Found %d open issues, %d closed issues\n", len(openIssues), len(closedIssues))

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
func (g *Generator) generateMarkdown(data *types.ReportData) string {
	var sb strings.Builder

	// Calculate overall statistics
	data.OverallStats = calculateOverallStats(data)

	// Generate header
	sb.WriteString(generateHeader(data))

	// Generate summary statistics
	sb.WriteString(generateSummaryStats(data.OverallStats))

	// Placeholder for Overall AI Summary
	sb.WriteString("## 📊 Overall Summary\n\n")
	sb.WriteString("[AI-generated overall summary will be here]\n\n")

	// Generate branches section
	if len(data.Branches) > 0 {
		sb.WriteString(generateBranchesSection(data.Branches))
	}

	// Placeholder for remaining sections
	sb.WriteString("---\n\n")
	sb.WriteString("*More sections (PRs, Issues, Reviews, Author Stats) will be added in the next steps*\n")

	return sb.String()
}
