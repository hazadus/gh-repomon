// +build integration

package integration

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hazadus/gh-repomon/internal/report"
	"github.com/hazadus/gh-repomon/internal/types"
)

// TestGenerateReport tests the full end-to-end report generation pipeline
func TestGenerateReport(t *testing.T) {
	// Setup mock clients
	mockGitHub := NewMockGitHubClient()
	mockLLM := NewMockLLMClient()

	// Create generator with mock clients
	gen := report.NewGeneratorWithClients(mockGitHub, mockLLM)

	// Setup test options
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	opts := report.Options{
		Repository: "owner/test-repo",
		Period: types.Period{
			From: yesterday,
			To:   now,
		},
		Model:    "openai/gpt-4o",
		Language: "english",
	}

	// Generate report
	reportText, err := gen.Generate(opts)
	if err != nil {
		t.Fatalf("Failed to generate report: %v", err)
	}

	// Verify report is not empty
	if len(reportText) == 0 {
		t.Error("Generated report is empty")
	}

	// Test: Report contains header
	t.Run("Contains header", func(t *testing.T) {
		if !strings.Contains(reportText, "# Repository Activity Report:") {
			t.Error("Report missing main header")
		}
		if !strings.Contains(reportText, "owner/test-repo") {
			t.Error("Report missing repository name")
		}
		if !strings.Contains(reportText, "**Repository**:") {
			t.Error("Report missing repository field")
		}
		if !strings.Contains(reportText, "**Period**:") {
			t.Error("Report missing period field")
		}
		if !strings.Contains(reportText, "**Report Generated**:") {
			t.Error("Report missing generation timestamp")
		}
	})

	// Test: Report contains Summary Statistics
	t.Run("Contains Summary Statistics", func(t *testing.T) {
		if !strings.Contains(reportText, "## Summary Statistics") {
			t.Error("Report missing Summary Statistics section")
		}
		if !strings.Contains(reportText, "- **Total Commits**:") {
			t.Error("Report missing Total Commits statistic")
		}
		if !strings.Contains(reportText, "- **Total Authors**:") {
			t.Error("Report missing Total Authors statistic")
		}
		if !strings.Contains(reportText, "- **Open Pull Requests**:") {
			t.Error("Report missing Open Pull Requests statistic")
		}
		if !strings.Contains(reportText, "- **Open Issues**:") {
			t.Error("Report missing Open Issues statistic")
		}
		if !strings.Contains(reportText, "- **Closed Issues**:") {
			t.Error("Report missing Closed Issues statistic")
		}
	})

	// Test: Report contains Overall AI Summary
	t.Run("Contains Overall AI Summary", func(t *testing.T) {
		if !strings.Contains(reportText, "## 📊 Overall Summary") {
			t.Error("Report missing Overall Summary section")
		}
		// Verify AI-generated content is present
		if !strings.Contains(reportText, "comprehensive summary") {
			t.Error("Report missing AI-generated overall summary content")
		}
	})

	// Test: Report contains branch sections
	t.Run("Contains branch sections", func(t *testing.T) {
		if !strings.Contains(reportText, "## 🌿 Branch:") {
			t.Error("Report missing branch sections")
		}
		if !strings.Contains(reportText, "### Statistics") {
			t.Error("Report missing branch statistics subsection")
		}
		if !strings.Contains(reportText, "### Commits") {
			t.Error("Report missing commits subsection")
		}
		if !strings.Contains(reportText, "### AI Summary") {
			t.Error("Report missing branch AI summary")
		}

		// Check for specific branches
		if !strings.Contains(reportText, "main") {
			t.Error("Report missing 'main' branch")
		}
		if !strings.Contains(reportText, "feature/new-ui") {
			t.Error("Report missing 'feature/new-ui' branch")
		}
	})

	// Test: Report contains commit information
	t.Run("Contains commit information", func(t *testing.T) {
		if !strings.Contains(reportText, "**Author**:") {
			t.Error("Report missing commit author information")
		}
		if !strings.Contains(reportText, "**Date**:") {
			t.Error("Report missing commit date information")
		}
		if !strings.Contains(reportText, "**Changes**:") {
			t.Error("Report missing commit changes information")
		}

		// Check for specific commits
		if !strings.Contains(reportText, "feat: add new feature") {
			t.Error("Report missing specific commit message")
		}
	})

	// Test: Report contains Pull Requests section
	t.Run("Contains Pull Requests section", func(t *testing.T) {
		if !strings.Contains(reportText, "## 🔀 Open Pull Requests") {
			t.Error("Report missing Open Pull Requests section")
		}
		if !strings.Contains(reportText, "## 🔄 Updated Pull Requests") {
			t.Error("Report missing Updated Pull Requests section")
		}
		if !strings.Contains(reportText, "Add authentication feature") {
			t.Error("Report missing specific PR title")
		}
		if !strings.Contains(reportText, "#### AI Summary") {
			t.Error("Report missing PR AI summary")
		}
	})

	// Test: Report contains Issues section
	t.Run("Contains Issues section", func(t *testing.T) {
		if !strings.Contains(reportText, "## 📋 Open Issues") {
			t.Error("Report missing Open Issues section")
		}
		if !strings.Contains(reportText, "## ✅ Closed Issues") {
			t.Error("Report missing Closed Issues section")
		}
		if !strings.Contains(reportText, "Bug: Login fails on mobile") {
			t.Error("Report missing specific issue title")
		}
	})

	// Test: Report contains Code Reviews section
	t.Run("Contains Code Reviews section", func(t *testing.T) {
		if !strings.Contains(reportText, "## 👀 Code Reviews") {
			t.Error("Report missing Code Reviews section")
		}
		if !strings.Contains(reportText, "**Total Reviews**:") {
			t.Error("Report missing total reviews count")
		}
	})

	// Test: Report contains Author Activity section
	t.Run("Contains Author Activity section", func(t *testing.T) {
		if !strings.Contains(reportText, "## 👥 Author Activity") {
			t.Error("Report missing Author Activity section")
		}
		if !strings.Contains(reportText, "#### Overall Statistics") {
			t.Error("Report missing author overall statistics")
		}
		if !strings.Contains(reportText, "#### Activity by Branch") {
			t.Error("Report missing author activity by branch")
		}
	})

	// Test: Report contains footer
	t.Run("Contains footer", func(t *testing.T) {
		if !strings.Contains(reportText, "gh-repomon") {
			t.Error("Report missing footer")
		}
	})

	// Test: Verify markdown link format
	t.Run("Uses correct markdown link format", func(t *testing.T) {
		// Check for markdown links [text](url)
		linkPattern := regexp.MustCompile(`\[.+\]\(https?://[^\)]+\)`)
		if !linkPattern.MatchString(reportText) {
			t.Error("Report missing proper markdown links")
		}
	})

	// Test: Verify date formatting
	t.Run("Uses correct date format", func(t *testing.T) {
		// Check for date format YYYY-MM-DD HH:MM
		datePattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}`)
		if !datePattern.MatchString(reportText) {
			t.Error("Report missing or using incorrect date format")
		}
	})

	// Test: Verify statistics are numeric
	t.Run("Contains numeric statistics", func(t *testing.T) {
		// Check that commits count is present and numeric
		commitPattern := regexp.MustCompile(`\*\*Total Commits\*\*: \d+`)
		if !commitPattern.MatchString(reportText) {
			t.Error("Report missing numeric commit count")
		}

		// Check for lines added/deleted format
		changesPattern := regexp.MustCompile(`\+\d+ / -\d+ lines`)
		if !changesPattern.MatchString(reportText) {
			t.Error("Report missing proper changes format (+X / -Y lines)")
		}
	})

	// Test: Verify AI summaries were generated
	t.Run("AI summaries generated", func(t *testing.T) {
		summaryCount := mockLLM.GetSummaryCount()
		expectedMinimum := 1 + 2 + 2 // 1 overall + 2 branches + 2 PRs
		if summaryCount < expectedMinimum {
			t.Errorf("Expected at least %d AI summaries, got %d", expectedMinimum, summaryCount)
		}
	})

	// Test: Verify no placeholder text remains
	t.Run("No placeholder text", func(t *testing.T) {
		placeholders := []string{
			"[AI summary will be generated here]",
			"TODO",
			"FIXME",
			"placeholder",
		}
		for _, placeholder := range placeholders {
			if strings.Contains(strings.ToLower(reportText), strings.ToLower(placeholder)) {
				t.Errorf("Report contains placeholder text: %s", placeholder)
			}
		}
	})
}

// TestGenerateReportWithNoActivity tests report generation when there's no activity
func TestGenerateReportWithNoActivity(t *testing.T) {
	// Create empty mock client
	mockGitHub := &MockGitHubClient{
		activeBranches:   []types.Branch{},
		openPRs:          []types.PullRequest{},
		updatedPRs:       []types.PullRequest{},
		openIssues:       []types.Issue{},
		closedIssues:     []types.Issue{},
		reviewsByAuthor:  map[string]int{},
		totalReviewCount: 0,
	}
	mockLLM := NewMockLLMClient()

	gen := report.NewGeneratorWithClients(mockGitHub, mockLLM)

	now := time.Now()
	opts := report.Options{
		Repository: "owner/empty-repo",
		Period: types.Period{
			From: now.Add(-24 * time.Hour),
			To:   now,
		},
		Model:    "openai/gpt-4o",
		Language: "english",
	}

	reportText, err := gen.Generate(opts)
	if err != nil {
		t.Fatalf("Failed to generate report for empty repository: %v", err)
	}

	// Should still contain header and sections, just with zero counts
	if !strings.Contains(reportText, "# Repository Activity Report:") {
		t.Error("Empty report missing header")
	}

	if !strings.Contains(reportText, "**Total Commits**: 0") {
		t.Error("Empty report should show 0 commits")
	}
}

// TestGenerateReportErrorHandling tests error handling during report generation
func TestGenerateReportErrorHandling(t *testing.T) {
	// This test would verify graceful error handling
	// For now, we test that the generator doesn't panic with edge cases

	mockGitHub := NewMockGitHubClient()
	mockLLM := NewMockLLMClient()
	gen := report.NewGeneratorWithClients(mockGitHub, mockLLM)

	// Test with invalid period (from > to)
	now := time.Now()
	opts := report.Options{
		Repository: "owner/test-repo",
		Period: types.Period{
			From: now,
			To:   now.Add(-24 * time.Hour), // Invalid: from is after to
		},
		Model:    "openai/gpt-4o",
		Language: "english",
	}

	// Should not panic, even with invalid input
	_, err := gen.Generate(opts)
	// The actual behavior (error or proceeding) depends on implementation
	// For now, we just verify it doesn't panic
	_ = err
}
