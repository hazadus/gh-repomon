package types

import "time"

// PullRequest represents a GitHub pull request.
type PullRequest struct {
	// Number is the PR number
	Number int
	// Title is the PR title
	Title string
	// Body is the PR description/body
	Body string
	// Author is the author of the PR
	Author Author
	// State is the current state (open, closed, merged)
	State string
	// CreatedAt is when the PR was created
	CreatedAt time.Time
	// UpdatedAt is when the PR was last updated
	UpdatedAt time.Time
	// Comments is the number of comments on the PR
	Comments int
	// Reviews is the number of reviews on the PR
	Reviews int
	// URL is the link to the PR on GitHub
	URL string
}
