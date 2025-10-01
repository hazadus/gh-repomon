// Package errors provides custom error types for gh-repomon
package errors

import (
	"fmt"
)

// ErrGitHubAuth represents GitHub authentication errors
type ErrGitHubAuth struct {
	Message string
	Cause   error
}

func (e *ErrGitHubAuth) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("GitHub authentication error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("GitHub authentication error: %s", e.Message)
}

// NewGitHubAuthError creates a new GitHub authentication error
func NewGitHubAuthError(message string, cause error) *ErrGitHubAuth {
	return &ErrGitHubAuth{
		Message: message,
		Cause:   cause,
	}
}

// ErrGitHubAPI represents GitHub API errors
type ErrGitHubAPI struct {
	Message    string
	StatusCode int
	Cause      error
}

func (e *ErrGitHubAPI) Error() string {
	if e.StatusCode > 0 {
		if e.Cause != nil {
			return fmt.Sprintf("GitHub API error (status %d): %s: %v", e.StatusCode, e.Message, e.Cause)
		}
		return fmt.Sprintf("GitHub API error (status %d): %s", e.StatusCode, e.Message)
	}
	if e.Cause != nil {
		return fmt.Sprintf("GitHub API error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("GitHub API error: %s", e.Message)
}

// NewGitHubAPIError creates a new GitHub API error
func NewGitHubAPIError(message string, statusCode int, cause error) *ErrGitHubAPI {
	return &ErrGitHubAPI{
		Message:    message,
		StatusCode: statusCode,
		Cause:      cause,
	}
}

// ErrRepoNotFound represents repository not found errors
type ErrRepoNotFound struct {
	Repository string
}

func (e *ErrRepoNotFound) Error() string {
	return fmt.Sprintf("repository not found: %s", e.Repository)
}

// NewRepoNotFoundError creates a new repository not found error
func NewRepoNotFoundError(repository string) *ErrRepoNotFound {
	return &ErrRepoNotFound{
		Repository: repository,
	}
}

// ErrInvalidParams represents invalid parameter errors
type ErrInvalidParams struct {
	Parameter string
	Reason    string
}

func (e *ErrInvalidParams) Error() string {
	return fmt.Sprintf("invalid parameter '%s': %s", e.Parameter, e.Reason)
}

// NewInvalidParamsError creates a new invalid parameters error
func NewInvalidParamsError(parameter, reason string) *ErrInvalidParams {
	return &ErrInvalidParams{
		Parameter: parameter,
		Reason:    reason,
	}
}

// ErrLLMAPI represents LLM API errors
type ErrLLMAPI struct {
	Message    string
	StatusCode int
	Cause      error
}

func (e *ErrLLMAPI) Error() string {
	if e.StatusCode > 0 {
		if e.Cause != nil {
			return fmt.Sprintf("LLM API error (status %d): %s: %v", e.StatusCode, e.Message, e.Cause)
		}
		return fmt.Sprintf("LLM API error (status %d): %s", e.StatusCode, e.Message)
	}
	if e.Cause != nil {
		return fmt.Sprintf("LLM API error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("LLM API error: %s", e.Message)
}

// NewLLMAPIError creates a new LLM API error
func NewLLMAPIError(message string, statusCode int, cause error) *ErrLLMAPI {
	return &ErrLLMAPI{
		Message:    message,
		StatusCode: statusCode,
		Cause:      cause,
	}
}
