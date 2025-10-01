package types

import "time"

// Commit represents a single commit in the repository.
type Commit struct {
	// SHA is the unique identifier of the commit
	SHA string
	// Message is the commit message
	Message string
	// Author is the author of the commit
	Author Author
	// Date is when the commit was created
	Date time.Time
	// Additions is the number of lines added in this commit
	Additions int
	// Deletions is the number of lines deleted in this commit
	Deletions int
	// URL is the link to the commit on GitHub
	URL string
}
