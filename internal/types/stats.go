package types

// BranchActivity represents activity statistics for a specific branch.
type BranchActivity struct {
	// Commits is the number of commits in this branch
	Commits int
	// Added is the number of lines added in this branch
	Added int
	// Deleted is the number of lines deleted in this branch
	Deleted int
}

// AuthorStats represents statistics for a single author.
type AuthorStats struct {
	// Author is the author information
	Author Author
	// TotalCommits is the total number of commits by this author
	TotalCommits int
	// TotalAdded is the total number of lines added by this author
	TotalAdded int
	// TotalDeleted is the total number of lines deleted by this author
	TotalDeleted int
	// PRsCreated is the number of pull requests created by this author
	PRsCreated int
	// IssuesCreated is the number of issues created by this author
	IssuesCreated int
	// ReviewsCount is the number of code reviews performed by this author
	ReviewsCount int
	// BranchActivity maps branch names to activity statistics
	BranchActivity map[string]BranchActivity
}

// OverallStats represents overall statistics for the repository activity.
type OverallStats struct {
	// TotalCommits is the total number of commits across all branches
	TotalCommits int
	// TotalAuthors is the total number of unique authors
	TotalAuthors int
	// OpenPRCount is the number of currently open pull requests
	OpenPRCount int
	// OpenIssuesCount is the number of currently open issues
	OpenIssuesCount int
	// ClosedIssuesCount is the number of issues closed during the period
	ClosedIssuesCount int
	// ReviewsCount is the total number of code reviews
	ReviewsCount int
}
