package types

import "time"

// Period represents a time period for the report.
type Period struct {
	// From is the start date of the period
	From time.Time
	// To is the end date of the period
	To time.Time
}

// ReportData contains all data collected for the report.
type ReportData struct {
	// Repository is the repository name (owner/repo)
	Repository string
	// RepositoryURL is the full URL to the repository
	RepositoryURL string
	// Period is the time period covered by this report
	Period Period
	// GeneratedAt is when this report was generated
	GeneratedAt time.Time
	// Branches is the list of branches with activity during the period
	Branches []Branch
	// OpenPRs is the list of currently open pull requests
	OpenPRs []PullRequest
	// UpdatedPRs is the list of pull requests updated during the period
	UpdatedPRs []PullRequest
	// OpenIssues is the list of currently open issues
	OpenIssues []Issue
	// ClosedIssues is the list of issues closed during the period
	ClosedIssues []Issue
	// AuthorStats is the statistics per author
	AuthorStats []AuthorStats
	// OverallStats is the overall statistics for the repository
	OverallStats OverallStats
}
