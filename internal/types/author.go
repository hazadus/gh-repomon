package types

// Author represents a GitHub user who contributed to the repository.
type Author struct {
	// Login is the GitHub username
	Login string
	// Name is the full name of the user (may be empty)
	Name string
	// ProfileURL is the link to the GitHub profile
	ProfileURL string
	// IsBot indicates whether this author is a bot account
	IsBot bool
}

// NewAuthor creates a new Author instance.
func NewAuthor(login, name, profileURL string, isBot bool) *Author {
	return &Author{
		Login:      login,
		Name:       name,
		ProfileURL: profileURL,
		IsBot:      isBot,
	}
}
