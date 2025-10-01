package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hazadus/gh-repomon/internal/github"
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
	rootCmd.MarkFlagRequired("repo")

	// Optional flags
	rootCmd.Flags().IntVarP(&days, "days", "d", 1, "Number of days back from today")
	rootCmd.Flags().StringVar(&fromDate, "from", "", "Start date of the period (YYYY-MM-DD)")
	rootCmd.Flags().StringVar(&toDate, "to", "", "End date of the period (YYYY-MM-DD)")
	rootCmd.Flags().StringVarP(&user, "user", "u", "", "Filter by user")
	rootCmd.Flags().BoolVar(&excludeBots, "exclude-bots", false, "Exclude bot accounts")
	rootCmd.Flags().StringVarP(&model, "model", "m", "openai/gpt-4o", "AI model to use")
	rootCmd.Flags().StringVarP(&language, "language", "l", "english", "Output language")
}

func run(cmd *cobra.Command, args []string) error {
	// Validate that --repo is provided
	if repo == "" {
		return fmt.Errorf("--repo flag is required")
	}

	// Calculate period
	var from, to time.Time
	var err error

	// If --from and --to are specified, use them and ignore --days
	if fromDate != "" && toDate != "" {
		from, err = parseDate(fromDate)
		if err != nil {
			return fmt.Errorf("invalid --from date: %w", err)
		}

		to, err = parseDate(toDate)
		if err != nil {
			return fmt.Errorf("invalid --to date: %w", err)
		}

		if from.After(to) {
			return fmt.Errorf("--from date must be before --to date")
		}
	} else if fromDate != "" || toDate != "" {
		return fmt.Errorf("both --from and --to must be specified together")
	} else {
		// Use --days
		to = time.Now().UTC()
		from = to.AddDate(0, 0, -days)
	}

	// Create GitHub client
	ghClient, err := github.NewClient(excludeBots)
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// Output parameters to stderr for verification
	fmt.Fprintf(os.Stderr, "Connected to GitHub API\n")
	fmt.Fprintf(os.Stderr, "Repository: %s\n", repo)
	fmt.Fprintf(os.Stderr, "Period: %s to %s\n", from.Format("2006-01-02"), to.Format("2006-01-02"))
	fmt.Fprintf(os.Stderr, "Days: %d\n", days)
	if user != "" {
		fmt.Fprintf(os.Stderr, "Filter by user: %s\n", user)
	}
	fmt.Fprintf(os.Stderr, "Exclude bots: %v\n", excludeBots)
	fmt.Fprintf(os.Stderr, "Model: %s\n", model)
	fmt.Fprintf(os.Stderr, "Language: %s\n", language)
	fmt.Fprintf(os.Stderr, "\n")

	// Get active branches
	fmt.Fprintf(os.Stderr, "Fetching active branches...\n")
	branches, err := ghClient.GetActiveBranches(repo, from, to)
	if err != nil {
		return fmt.Errorf("failed to get active branches: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Found %d active branches\n\n", len(branches))

	// Display information for each branch
	for _, branch := range branches {
		fmt.Fprintf(os.Stderr, "Branch: %s\n", branch.Name)
		fmt.Fprintf(os.Stderr, "  Commits: %d\n", len(branch.Commits))
		fmt.Fprintf(os.Stderr, "  Total Added: +%d lines\n", branch.TotalAdded)
		fmt.Fprintf(os.Stderr, "  Total Deleted: -%d lines\n", branch.TotalDeleted)
		fmt.Fprintf(os.Stderr, "  Authors: %v\n", branch.Authors)
		fmt.Fprintf(os.Stderr, "\n")
	}

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
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
