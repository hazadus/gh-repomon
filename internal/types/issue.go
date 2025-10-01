package types

import "time"

// Issue represents a GitHub issue.
type Issue struct {
	// Number is the issue number
	Number int
	// Title is the issue title
	Title string
	// Body is the issue description/body
	Body string
	// Author is the author of the issue
	Author Author
	// State is the current state (open, closed)
	State string
	// CreatedAt is when the issue was created
	CreatedAt time.Time
	// ClosedAt is when the issue was closed (nil if still open)
	ClosedAt *time.Time
	// Labels is the list of labels attached to the issue
	Labels []string
	// Assignees is the list of users assigned to the issue
	Assignees []Author
	// URL is the link to the issue on GitHub
	URL string
}
