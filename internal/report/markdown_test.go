package report

import (
	"strings"
	"testing"
	"time"

	"github.com/hazadus/gh-repomon/internal/types"
)

func TestFormatCommitMessage(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		wantShort string
		wantFull  string
	}{
		{
			name:      "Single line message",
			message:   "Fix bug in parser",
			wantShort: "Fix bug in parser",
			wantFull:  "Fix bug in parser",
		},
		{
			name:      "Multi-line message",
			message:   "Add new feature\n\nThis adds a new feature with detailed description",
			wantShort: "Add new feature",
			wantFull:  "Add new feature\n\nThis adds a new feature with detailed description",
		},
		{
			name:      "Long single line message",
			message:   "This is a very long commit message that exceeds the 72 character limit and should be truncated",
			wantShort: "This is a very long commit message that exceeds the 72 character limi...",
			wantFull:  "This is a very long commit message that exceeds the 72 character limit and should be truncated",
		},
		{
			name:      "Empty message",
			message:   "",
			wantShort: "",
			wantFull:  "",
		},
		{
			name:      "Message with newlines",
			message:   "feat: add tests\n\n- Add unit tests\n- Add integration tests\n\nCloses #123",
			wantShort: "feat: add tests",
			wantFull:  "feat: add tests\n\n- Add unit tests\n- Add integration tests\n\nCloses #123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotShort, gotFull := formatCommitMessage(tt.message)
			if gotShort != tt.wantShort {
				t.Errorf("formatCommitMessage() short = %v, want %v", gotShort, tt.wantShort)
			}
			if gotFull != tt.wantFull {
				t.Errorf("formatCommitMessage() full = %v, want %v", gotFull, tt.wantFull)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			name: "Standard date",
			time: time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
			want: "2024-01-15 14:30",
		},
		{
			name: "Single digit month and day",
			time: time.Date(2024, 3, 5, 9, 5, 0, 0, time.UTC),
			want: "2024-03-05 09:05",
		},
		{
			name: "End of year",
			time: time.Date(2023, 12, 31, 23, 59, 0, 0, time.UTC),
			want: "2023-12-31 23:59",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDate(tt.time)
			if got != tt.want {
				t.Errorf("formatDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatAuthorLinks(t *testing.T) {
	tests := []struct {
		name    string
		authors []string
		want    string
	}{
		{
			name:    "Single author",
			authors: []string{"octocat"},
			want:    "[octocat](https://github.com/octocat)",
		},
		{
			name:    "Multiple authors",
			authors: []string{"alice", "bob", "charlie"},
			want:    "[alice](https://github.com/alice), [bob](https://github.com/bob), [charlie](https://github.com/charlie)",
		},
		{
			name:    "Empty list",
			authors: []string{},
			want:    "none",
		},
		{
			name:    "Nil list",
			authors: nil,
			want:    "none",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAuthorLinks(tt.authors)
			if got != tt.want {
				t.Errorf("formatAuthorLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateHeader(t *testing.T) {
	data := &types.ReportData{
		Repository:    "owner/repo",
		RepositoryURL: "https://github.com/owner/repo",
		Period: types.Period{
			From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			To:   time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
		},
		GeneratedAt: time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
	}

	got := generateHeader(data)

	// Check that header contains expected elements
	if !strings.Contains(got, "# Repository Activity Report: repo") {
		t.Errorf("generateHeader() missing main heading")
	}
	if !strings.Contains(got, "[owner/repo](https://github.com/owner/repo)") {
		t.Errorf("generateHeader() missing repository link")
	}
	if !strings.Contains(got, "2024-01-01 to 2024-01-07") {
		t.Errorf("generateHeader() missing period")
	}
	if !strings.Contains(got, "2024-01-08 12:00:00 UTC") {
		t.Errorf("generateHeader() missing generation timestamp")
	}
}

func TestGenerateSummaryStats(t *testing.T) {
	stats := types.OverallStats{
		TotalCommits:      42,
		TotalAuthors:      5,
		OpenPRCount:       3,
		OpenIssuesCount:   7,
		ClosedIssuesCount: 2,
		ReviewsCount:      15,
	}

	got := generateSummaryStats(stats)

	// Check that summary contains all stats
	if !strings.Contains(got, "## Summary Statistics") {
		t.Errorf("generateSummaryStats() missing heading")
	}
	if !strings.Contains(got, "**Total Commits**: 42") {
		t.Errorf("generateSummaryStats() missing or incorrect commits count")
	}
	if !strings.Contains(got, "**Total Authors**: 5") {
		t.Errorf("generateSummaryStats() missing or incorrect authors count")
	}
	if !strings.Contains(got, "**Open Pull Requests**: 3") {
		t.Errorf("generateSummaryStats() missing or incorrect open PRs count")
	}
	if !strings.Contains(got, "**Open Issues**: 7") {
		t.Errorf("generateSummaryStats() missing or incorrect open issues count")
	}
	if !strings.Contains(got, "**Closed Issues**: 2") {
		t.Errorf("generateSummaryStats() missing or incorrect closed issues count")
	}
	if !strings.Contains(got, "**Code Reviews**: 15") {
		t.Errorf("generateSummaryStats() missing or incorrect reviews count")
	}
}

func TestCalculateOverallStats(t *testing.T) {
	author1 := types.Author{Login: "alice", ProfileURL: "https://github.com/alice"}
	author2 := types.Author{Login: "bob", ProfileURL: "https://github.com/bob"}

	data := &types.ReportData{
		Branches: []types.Branch{
			{
				Name: "main",
				Commits: []types.Commit{
					{Author: author1},
					{Author: author2},
				},
			},
			{
				Name: "feature",
				Commits: []types.Commit{
					{Author: author1},
				},
			},
		},
		OpenPRs: []types.PullRequest{
			{Number: 1, Reviews: 2},
			{Number: 2, Reviews: 3},
		},
		UpdatedPRs: []types.PullRequest{
			{Number: 1, Reviews: 2}, // Same as open PR, should not double count
			{Number: 3, Reviews: 1}, // Different PR
		},
		OpenIssues: []types.Issue{
			{Number: 10},
			{Number: 11},
		},
		ClosedIssues: []types.Issue{
			{Number: 5},
		},
	}

	got := calculateOverallStats(data)

	if got.TotalCommits != 3 {
		t.Errorf("calculateOverallStats() TotalCommits = %d, want 3", got.TotalCommits)
	}
	if got.TotalAuthors != 2 {
		t.Errorf("calculateOverallStats() TotalAuthors = %d, want 2", got.TotalAuthors)
	}
	if got.OpenPRCount != 2 {
		t.Errorf("calculateOverallStats() OpenPRCount = %d, want 2", got.OpenPRCount)
	}
	if got.OpenIssuesCount != 2 {
		t.Errorf("calculateOverallStats() OpenIssuesCount = %d, want 2", got.OpenIssuesCount)
	}
	if got.ClosedIssuesCount != 1 {
		t.Errorf("calculateOverallStats() ClosedIssuesCount = %d, want 1", got.ClosedIssuesCount)
	}
	// Reviews: 2 + 3 (from open) + 1 (from updated, not open) = 6
	if got.ReviewsCount != 6 {
		t.Errorf("calculateOverallStats() ReviewsCount = %d, want 6", got.ReviewsCount)
	}
}

func TestCalculateAuthorStats(t *testing.T) {
	author1 := types.Author{Login: "alice", ProfileURL: "https://github.com/alice"}
	author2 := types.Author{Login: "bob", ProfileURL: "https://github.com/bob"}

	data := &types.ReportData{
		Branches: []types.Branch{
			{
				Name: "main",
				Commits: []types.Commit{
					{Author: author1, Additions: 100, Deletions: 20},
					{Author: author1, Additions: 50, Deletions: 10},
					{Author: author2, Additions: 30, Deletions: 5},
				},
			},
			{
				Name: "feature",
				Commits: []types.Commit{
					{Author: author1, Additions: 200, Deletions: 50},
				},
			},
		},
		OpenPRs: []types.PullRequest{
			{Author: author1},
			{Author: author2},
		},
		OpenIssues: []types.Issue{
			{Author: author1},
		},
	}

	got := calculateAuthorStats(data)

	if len(got) != 2 {
		t.Fatalf("calculateAuthorStats() returned %d authors, want 2", len(got))
	}

	// Should be sorted by commit count, alice has 3, bob has 1
	if got[0].Author.Login != "alice" {
		t.Errorf("calculateAuthorStats() first author = %s, want alice", got[0].Author.Login)
	}
	if got[0].TotalCommits != 3 {
		t.Errorf("calculateAuthorStats() alice TotalCommits = %d, want 3", got[0].TotalCommits)
	}
	if got[0].TotalAdded != 350 {
		t.Errorf("calculateAuthorStats() alice TotalAdded = %d, want 350", got[0].TotalAdded)
	}
	if got[0].TotalDeleted != 80 {
		t.Errorf("calculateAuthorStats() alice TotalDeleted = %d, want 80", got[0].TotalDeleted)
	}
	if got[0].PRsCreated != 1 {
		t.Errorf("calculateAuthorStats() alice PRsCreated = %d, want 1", got[0].PRsCreated)
	}
	if got[0].IssuesCreated != 1 {
		t.Errorf("calculateAuthorStats() alice IssuesCreated = %d, want 1", got[0].IssuesCreated)
	}

	// Check branch activity for alice
	if len(got[0].BranchActivity) != 2 {
		t.Errorf("calculateAuthorStats() alice has %d branch activities, want 2", len(got[0].BranchActivity))
	}
	if got[0].BranchActivity["main"].Commits != 2 {
		t.Errorf("calculateAuthorStats() alice main branch commits = %d, want 2", got[0].BranchActivity["main"].Commits)
	}
	if got[0].BranchActivity["feature"].Commits != 1 {
		t.Errorf("calculateAuthorStats() alice feature branch commits = %d, want 1", got[0].BranchActivity["feature"].Commits)
	}

	// Check bob
	if got[1].Author.Login != "bob" {
		t.Errorf("calculateAuthorStats() second author = %s, want bob", got[1].Author.Login)
	}
	if got[1].TotalCommits != 1 {
		t.Errorf("calculateAuthorStats() bob TotalCommits = %d, want 1", got[1].TotalCommits)
	}
}
