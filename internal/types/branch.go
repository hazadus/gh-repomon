package types

// Branch represents a branch with its activity.
type Branch struct {
	// Name is the branch name
	Name string
	// Commits is the list of commits in this branch during the period
	Commits []Commit
	// PRs is the list of pull requests associated with this branch
	PRs []PullRequest
	// TotalAdded is the total number of lines added across all commits
	TotalAdded int
	// TotalDeleted is the total number of lines deleted across all commits
	TotalDeleted int
	// Authors is the list of unique author logins who contributed to this branch
	Authors []string
	// AISummary is the AI-generated summary of branch activity
	AISummary string
}
