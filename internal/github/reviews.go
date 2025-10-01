package github

import (
	"fmt"
	"os"

	"github.com/hazadus/gh-repomon/internal/types"
)

// Review represents a simplified code review
type Review struct {
	User struct {
		Login string `json:"login"`
	} `json:"user"`
	State       string `json:"state"`
	SubmittedAt string `json:"submitted_at"`
}

// GetReviews retrieves all reviews for a specific pull request
func (c *Client) GetReviews(repo string, prNumber int) ([]Review, error) {
	path := fmt.Sprintf("repos/%s/pulls/%d/reviews", repo, prNumber)

	var reviews []Review
	err := c.doWithRetry("GET", path, nil, &reviews)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews for PR #%d: %w", prNumber, err)
	}

	// Filter out bot reviews if excludeBots is enabled
	if c.excludeBots {
		filtered := make([]Review, 0)
		for _, review := range reviews {
			if !c.isBot(review.User.Login) {
				filtered = append(filtered, review)
			}
		}
		return filtered, nil
	}

	return reviews, nil
}

// GetAllReviews counts the total number of reviews across all provided PRs
func (c *Client) GetAllReviews(repo string, prs []types.PullRequest) (int, error) {
	totalReviews := 0

	for _, pr := range prs {
		reviews, err := c.GetReviews(repo, pr.Number)
		if err != nil {
			// Log error but continue counting other PRs
			fmt.Fprintf(os.Stderr, "Warning: failed to get reviews for PR #%d: %v\n", pr.Number, err)
			continue
		}
		totalReviews += len(reviews)
	}

	return totalReviews, nil
}

// GetReviewsByAuthor groups reviews by reviewer login across all provided PRs
func (c *Client) GetReviewsByAuthor(repo string, prs []types.PullRequest) (map[string]int, error) {
	reviewsByAuthor := make(map[string]int)

	for _, pr := range prs {
		reviews, err := c.GetReviews(repo, pr.Number)
		if err != nil {
			// Log error but continue with other PRs
			fmt.Fprintf(os.Stderr, "Warning: failed to get reviews for PR #%d: %v\n", pr.Number, err)
			continue
		}

		for _, review := range reviews {
			reviewsByAuthor[review.User.Login]++
		}
	}

	return reviewsByAuthor, nil
}

// GetReviewsForPR retrieves reviews for a PR and returns the count
// This is a helper method that also stores review info in the PR object
func (c *Client) GetReviewsForPR(repo string, pr *types.PullRequest) error {
	path := fmt.Sprintf("repos/%s/pulls/%d/reviews", repo, pr.Number)

	var reviewsData []map[string]interface{}
	err := c.doWithRetry("GET", path, nil, &reviewsData)
	if err != nil {
		return fmt.Errorf("failed to get reviews for PR #%d: %w", pr.Number, err)
	}

	// Count non-bot reviews
	reviewCount := 0
	for _, reviewData := range reviewsData {
		if user, ok := reviewData["user"].(map[string]interface{}); ok {
			if login, ok := user["login"].(string); ok {
				if !c.excludeBots || !c.isBot(login) {
					reviewCount++
				}
			}
		}
	}

	pr.Reviews = reviewCount
	return nil
}
