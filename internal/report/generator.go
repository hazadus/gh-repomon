package report

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/hazadus/gh-repomon/internal/github"
	"github.com/hazadus/gh-repomon/internal/llm"
	"github.com/hazadus/gh-repomon/internal/logger"
	"github.com/hazadus/gh-repomon/internal/types"
	"github.com/hazadus/gh-repomon/internal/utils"
)

// GitHubClient defines the interface for GitHub API operations
type GitHubClient interface {
	GetActiveBranches(repo string, from, to time.Time) ([]types.Branch, error)
	GetOpenPullRequests(repo string) ([]types.PullRequest, error)
	GetUpdatedPullRequests(repo, from, to string) ([]types.PullRequest, error)
	GetOpenIssues(repo string) ([]types.Issue, error)
	GetClosedIssues(repo, from, to string) ([]types.Issue, error)
}

// LLMClient defines the interface for LLM operations
type LLMClient interface {
	GenerateOverallSummary(data *types.ReportData, language, model string) (string, error)
	GenerateBranchSummary(branch *types.Branch, language, model string) (string, error)
	GeneratePRSummary(pr *types.PullRequest, language, model string) (string, error)
}

// Generator generates reports based on GitHub activity data
type Generator struct {
	githubClient GitHubClient
	llmClient    LLMClient
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

// NewGeneratorWithClients creates a new report generator with custom client implementations
// This is useful for testing with mock clients
func NewGeneratorWithClients(ghClient GitHubClient, llmClient LLMClient) *Generator {
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

		// Generate branch summaries in parallel with rate limiting
		maxWorkers := 5 // Limit concurrent LLM requests
		branchSummaryErrors := 0
		var branchMu sync.Mutex

		err = utils.ProcessInParallel(data.Branches, maxWorkers, func(branch types.Branch) error {
			branchSummary, err := g.llmClient.GenerateBranchSummary(&branch, opts.Language, opts.Model)

			// Find the branch in data.Branches and update it
			branchMu.Lock()
			for i := range data.Branches {
				if data.Branches[i].Name == branch.Name {
					if err != nil {
						g.logger.Warning(fmt.Sprintf("Failed to generate summary for branch %s: %v", branch.Name, err))
						data.Branches[i].AISummary = fmt.Sprintf("Development activity in branch %s", branch.Name)
						branchSummaryErrors++
					} else {
						data.Branches[i].AISummary = branchSummary
					}
					stats.TotalAISummaries++
					break
				}
			}
			branchMu.Unlock()

			// Don't fail the entire process if one summary fails
			return nil
		})

		if err != nil {
			g.logger.Warning(fmt.Sprintf("Error generating branch summaries: %v", err))
		}

		stats.SuccessfulSummaries += (stats.TotalBranches - branchSummaryErrors)
		stats.FailedSummaries += branchSummaryErrors
		g.logger.Success(fmt.Sprintf("Branch summaries generated (%d/%d)", stats.TotalBranches-branchSummaryErrors, stats.TotalBranches))

		// Generate PR summaries in parallel
		totalPRs := len(data.OpenPRs) + len(data.UpdatedPRs)
		prSuccessCount := 0
		prErrors := 0
		var prMu sync.Mutex

		// Generate summaries for open PRs
		err = utils.ProcessInParallel(data.OpenPRs, maxWorkers, func(pr types.PullRequest) error {
			prSummary, err := g.llmClient.GeneratePRSummary(&pr, opts.Language, opts.Model)

			prMu.Lock()
			for i := range data.OpenPRs {
				if data.OpenPRs[i].Number == pr.Number {
					if err != nil {
						g.logger.Warning(fmt.Sprintf("Failed to generate summary for PR #%d: %v", pr.Number, err))
						data.OpenPRs[i].AISummary = fmt.Sprintf("Pull request: %s", pr.Title)
						prErrors++
					} else {
						data.OpenPRs[i].AISummary = prSummary
						prSuccessCount++
					}
					stats.TotalAISummaries++
					break
				}
			}
			prMu.Unlock()
			return nil
		})

		if err != nil {
			g.logger.Warning(fmt.Sprintf("Error generating open PR summaries: %v", err))
		}

		// Generate summaries for updated PRs
		err = utils.ProcessInParallel(data.UpdatedPRs, maxWorkers, func(pr types.PullRequest) error {
			prSummary, err := g.llmClient.GeneratePRSummary(&pr, opts.Language, opts.Model)

			prMu.Lock()
			for i := range data.UpdatedPRs {
				if data.UpdatedPRs[i].Number == pr.Number {
					if err != nil {
						g.logger.Warning(fmt.Sprintf("Failed to generate summary for PR #%d: %v", pr.Number, err))
						data.UpdatedPRs[i].AISummary = fmt.Sprintf("Pull request: %s", pr.Title)
						prErrors++
					} else {
						data.UpdatedPRs[i].AISummary = prSummary
						prSuccessCount++
					}
					stats.TotalAISummaries++
					break
				}
			}
			prMu.Unlock()
			return nil
		})

		if err != nil {
			g.logger.Warning(fmt.Sprintf("Error generating updated PR summaries: %v", err))
		}

		stats.SuccessfulSummaries += prSuccessCount
		stats.FailedSummaries += prErrors
		g.logger.Success(fmt.Sprintf("PR summaries generated (%d/%d)", prSuccessCount, totalPRs))
	} else {
		overallSummary = "[AI summary generation disabled]"
	}

	// Generate markdown report
	report := g.generateMarkdown(data, overallSummary, stats)

	return report, nil
}

// collectData collects all necessary data from GitHub API in parallel
func (g *Generator) collectData(opts Options) (*types.ReportData, error) {
	// Convert times to ISO8601 format for API calls
	fromISO := opts.Period.From.Format(time.RFC3339)
	toISO := opts.Period.To.Format(time.RFC3339)

	// Use errgroup for parallel data collection
	var eg errgroup.Group
	var branches []types.Branch
	var openPRs, updatedPRs []types.PullRequest
	var openIssues, closedIssues []types.Issue
	var mu sync.Mutex

	// Get active branches
	g.logger.Progress("Collecting branches...")
	eg.Go(func() error {
		b, err := g.githubClient.GetActiveBranches(opts.Repository, opts.Period.From, opts.Period.To)
		if err != nil {
			return err
		}
		mu.Lock()
		branches = b
		mu.Unlock()
		return nil
	})

	// Get open and updated pull requests
	g.logger.Progress("Collecting pull requests...")
	eg.Go(func() error {
		prs, err := g.githubClient.GetOpenPullRequests(opts.Repository)
		if err != nil {
			return err
		}
		mu.Lock()
		openPRs = prs
		mu.Unlock()
		return nil
	})

	eg.Go(func() error {
		prs, err := g.githubClient.GetUpdatedPullRequests(opts.Repository, fromISO, toISO)
		if err != nil {
			return err
		}
		mu.Lock()
		updatedPRs = prs
		mu.Unlock()
		return nil
	})

	// Get open and closed issues
	g.logger.Progress("Collecting issues...")
	eg.Go(func() error {
		issues, err := g.githubClient.GetOpenIssues(opts.Repository)
		if err != nil {
			return err
		}
		mu.Lock()
		openIssues = issues
		mu.Unlock()
		return nil
	})

	eg.Go(func() error {
		issues, err := g.githubClient.GetClosedIssues(opts.Repository, fromISO, toISO)
		if err != nil {
			return err
		}
		mu.Lock()
		closedIssues = issues
		mu.Unlock()
		return nil
	})

	// Wait for all goroutines to complete
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// Log results
	g.logger.Success(fmt.Sprintf("Found %d active branches", len(branches)))

	totalCommits := 0
	for _, branch := range branches {
		totalCommits += len(branch.Commits)
	}
	g.logger.Success(fmt.Sprintf("Collected %d commits", totalCommits))
	g.logger.Success(fmt.Sprintf("Found %d open PRs, %d updated PRs", len(openPRs), len(updatedPRs)))
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
