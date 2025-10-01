package report

import (
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

	// Get open pull requests
	openPRs, err := g.githubClient.GetOpenPullRequests(opts.Repository)
	if err != nil {
		return nil, err
	}

	// Get updated pull requests
	updatedPRs, err := g.githubClient.GetUpdatedPullRequests(opts.Repository, fromISO, toISO)
	if err != nil {
		return nil, err
	}

	// Get open issues
	openIssues, err := g.githubClient.GetOpenIssues(opts.Repository)
	if err != nil {
		return nil, err
	}

	// Get closed issues
	closedIssues, err := g.githubClient.GetClosedIssues(opts.Repository, fromISO, toISO)
	if err != nil {
		return nil, err
	}

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
	// Placeholder for now - will be implemented in next steps
	return "# Report\n\nData collected successfully\n"
}
