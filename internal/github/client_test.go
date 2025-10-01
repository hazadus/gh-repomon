package github

import (
	"testing"

	"github.com/hazadus/gh-repomon/internal/types"
)

func TestIsBot(t *testing.T) {
	tests := []struct {
		name  string
		login string
		want  bool
	}{
		{
			name:  "Regular user",
			login: "octocat",
			want:  false,
		},
		{
			name:  "Bot with [bot] suffix",
			login: "github-actions[bot]",
			want:  true,
		},
		{
			name:  "Dependabot",
			login: "dependabot",
			want:  true,
		},
		{
			name:  "GitHub Actions",
			login: "github-actions",
			want:  true,
		},
		{
			name:  "Renovate bot",
			login: "renovate",
			want:  true,
		},
		{
			name:  "Custom bot with [bot] suffix",
			login: "my-custom-bot[bot]",
			want:  true,
		},
		{
			name:  "User with bot in name but not suffix",
			login: "botmaster",
			want:  false,
		},
		{
			name:  "Empty login",
			login: "",
			want:  false,
		},
		{
			name:  "User named bot",
			login: "bot",
			want:  false,
		},
	}

	// Create a client for testing
	client := &Client{
		excludeBots: false,
		userCache:   make(map[string]*types.Author),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := client.isBot(tt.login)
			if got != tt.want {
				t.Errorf("isBot(%q) = %v, want %v", tt.login, got, tt.want)
			}
		})
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name   string
		errMsg string
		want   bool
	}{
		{
			name:   "Nil error",
			errMsg: "",
			want:   false,
		},
		{
			name:   "Timeout error",
			errMsg: "request timeout",
			want:   true,
		},
		{
			name:   "Connection error",
			errMsg: "connection refused",
			want:   true,
		},
		{
			name:   "500 Internal Server Error",
			errMsg: "HTTP 500",
			want:   true,
		},
		{
			name:   "502 Bad Gateway",
			errMsg: "HTTP 502",
			want:   true,
		},
		{
			name:   "503 Service Unavailable",
			errMsg: "HTTP 503",
			want:   true,
		},
		{
			name:   "504 Gateway Timeout",
			errMsg: "HTTP 504",
			want:   true,
		},
		{
			name:   "404 Not Found",
			errMsg: "HTTP 404",
			want:   false,
		},
		{
			name:   "403 Forbidden",
			errMsg: "HTTP 403",
			want:   false,
		},
		{
			name:   "401 Unauthorized",
			errMsg: "HTTP 401",
			want:   false,
		},
		{
			name:   "Other error",
			errMsg: "something went wrong",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.errMsg != "" {
				err = &mockError{msg: tt.errMsg}
			}

			got := isRetryableError(err)
			if got != tt.want {
				t.Errorf("isRetryableError(%v) = %v, want %v", tt.errMsg, got, tt.want)
			}
		})
	}
}

// mockError is a simple error type for testing
type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}
