package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hazadus/gh-repomon/internal/errors"
	"github.com/hazadus/gh-repomon/internal/github"
	"github.com/hazadus/gh-repomon/internal/llm"
	"github.com/hazadus/gh-repomon/internal/logger"
	"github.com/hazadus/gh-repomon/internal/report"
	"github.com/hazadus/gh-repomon/internal/types"
	"github.com/spf13/cobra"
)

var (
	repo        string
	days        int
	fromDate    string
	toDate      string
	user        string
	excludeBots bool
	model       string
	language    string
	noAI        bool
	verbose     bool
)

var rootCmd = &cobra.Command{
	Use:   "repomon",
	Short: "GitHub Repository Activity Monitor",
	Long:  `A CLI tool to monitor and report on GitHub repository activity`,
	RunE:  run,
}

func init() {
	// Required flags
	rootCmd.Flags().StringVarP(&repo, "repo", "r", "", "Repository name (owner/repo) (required)")
	_ = rootCmd.MarkFlagRequired("repo")

	// Optional flags
	rootCmd.Flags().IntVarP(&days, "days", "d", 1, "Number of days back from today")
	rootCmd.Flags().StringVar(&fromDate, "from", "", "Start date of the period (YYYY-MM-DD)")
	rootCmd.Flags().StringVar(&toDate, "to", "", "End date of the period (YYYY-MM-DD)")
	rootCmd.Flags().StringVarP(&user, "user", "u", "", "Filter by user")
	rootCmd.Flags().BoolVar(&excludeBots, "exclude-bots", false, "Exclude bot accounts")
	rootCmd.Flags().StringVarP(&model, "model", "m", "gpt-4o", "AI model to use")
	rootCmd.Flags().StringVarP(&language, "language", "l", "english", "Output language")
	rootCmd.Flags().BoolVar(&noAI, "no-ai", false, "Disable AI summary generation (faster)")
	rootCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")
}

func run(cmd *cobra.Command, args []string) error {
	// Create logger with verbose setting
	log := logger.New()
	if verbose {
		log.SetVerbose(true)
	}

	// Validate that --repo is provided
	if repo == "" {
		return errors.NewInvalidParamsError("repo", "repository flag is required")
	}

	// Calculate period
	var from, to time.Time
	var err error

	// If --from and --to are specified, use them and ignore --days
	if fromDate != "" && toDate != "" {
		from, err = parseDate(fromDate)
		if err != nil {
			return errors.NewInvalidParamsError("from", fmt.Sprintf("invalid date format: %v", err))
		}

		to, err = parseDate(toDate)
		if err != nil {
			return errors.NewInvalidParamsError("to", fmt.Sprintf("invalid date format: %v", err))
		}

		if from.After(to) {
			return errors.NewInvalidParamsError("from/to", "from date must be before to date")
		}
	} else if fromDate != "" || toDate != "" {
		return errors.NewInvalidParamsError("from/to", "both --from and --to must be specified together")
	} else {
		// Use --days
		to = time.Now().UTC()
		from = to.AddDate(0, 0, -days)
	}

	// Create GitHub client
	log.Info("Connecting to GitHub API...")
	ghClient, err := github.NewClient(excludeBots)
	if err != nil {
		return errors.NewGitHubAuthError("failed to create GitHub client", err)
	}
	log.Success("Connected to GitHub API")

	// Create LLM client (unless --no-ai is specified)
	var llmClient *llm.Client
	if noAI {
		log.Info("AI summary generation disabled (--no-ai)")
		llmClient = nil
	} else {
		log.Info("Connecting to LLM API...")
		llmClient, err = llm.NewClient()
		if err != nil {
			log.Warning("Failed to create LLM client, AI summaries will be disabled")
			llmClient = nil
		} else {
			log.Success("Connected to LLM API")
		}
	}

	log.Info(fmt.Sprintf("Analyzing repository %s (%s to %s)",
		repo,
		from.Format("2006-01-02"),
		to.Format("2006-01-02")))

	// Create report options
	opts := report.Options{
		Repository: repo,
		Period: types.Period{
			From: from,
			To:   to,
		},
		User:     user,
		Model:    model,
		Language: language,
	}

	// Create report generator
	generator := report.NewGenerator(ghClient, llmClient)

	// Generate report
	log.Progress("Collecting repository data...")
	reportText, err := generator.Generate(opts)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	log.Success("Report generated successfully!")

	// Output report to stdout
	fmt.Println(reportText)

	return nil
}

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(dateStr string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("date must be in YYYY-MM-DD format: %w", err)
	}
	return t, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		// Handle specific error types
		switch e := err.(type) {
		case *errors.ErrGitHubAuth:
			fmt.Fprintf(os.Stderr, "❌ Authentication Error: %v\n", e)
			fmt.Fprintf(os.Stderr, "Please ensure you are authenticated with GitHub CLI: gh auth login\n")
			os.Exit(1)
		case *errors.ErrGitHubAPI:
			fmt.Fprintf(os.Stderr, "❌ GitHub API Error: %v\n", e)
			os.Exit(1)
		case *errors.ErrRepoNotFound:
			fmt.Fprintf(os.Stderr, "❌ Repository Not Found: %v\n", e)
			os.Exit(1)
		case *errors.ErrInvalidParams:
			fmt.Fprintf(os.Stderr, "❌ Invalid Parameters: %v\n", e)
			os.Exit(1)
		case *errors.ErrLLMAPI:
			fmt.Fprintf(os.Stderr, "❌ LLM API Error: %v\n", e)
			os.Exit(1)
		default:
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}
	}
}
