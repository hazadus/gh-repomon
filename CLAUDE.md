# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**gh-repomon** is a CLI tool written in Go that generates comprehensive activity reports for GitHub repositories with AI-powered summaries. It uses the GitHub REST API (via GitHub CLI) to collect repository data and the GitHub Models API to generate intelligent summaries.

## Development Commands

### Building and Running

```bash
# Run the application with arguments
just run --repo owner/repo --days 7

# Build binary (output: ./bin/gh-repomon)
just build

# Build for all platforms
just build-all

# Clean build artifacts
just clean
```

### Testing

```bash
# Run unit tests with race detection
just test

# Run all tests (unit + integration)
just test-all

# Run integration tests only
just test-integration

# Run tests with coverage report
just test-coverage
```

### Code Quality

```bash
# Format code with gofmt
just format

# Run linter (requires golangci-lint)
just lint

# Install/update dependencies
just deps
```

### Release

```bash
# Create local release build with checksums
just release

# Test goreleaser configuration
just release-test
```

## Architecture Overview

### High-Level Flow

1. **CLI Layer** ([cmd/repomon/main.go](cmd/repomon/main.go)) - Parses arguments using Cobra, orchestrates workflow
2. **Data Collection** ([internal/github/](internal/github/)) - Fetches data from GitHub API in parallel using `go-gh` library
3. **AI Enhancement** ([internal/llm/](internal/llm/)) - Generates summaries using GitHub Models API with YAML-based prompts
4. **Report Generation** ([internal/report/](internal/report/)) - Aggregates data, calculates statistics, formats markdown output

### Key Architectural Patterns

- **Parallel Processing**: Uses worker pools and errgroups for concurrent API requests (branches, PRs, issues, commit stats, LLM summaries)
- **Graceful Degradation**: AI summaries fail gracefully without breaking report generation
- **Interface-Based Design**: GitHub and LLM clients use interfaces for easy testing with mocks
- **YAML-Based Prompts**: AI prompts stored in `internal/llm/prompts/*.prompt.yml` and embedded in binary via `//go:embed` for self-contained distribution
- **Prompt Override**: External prompt files in `internal/llm/prompts/` take precedence over embedded ones during development

### Module Structure

```
internal/
├── github/          # GitHub REST API client
│   ├── client.go    # Client initialization, retry logic, bot detection
│   ├── commits.go   # Fetch commits and stats
│   ├── branches.go  # Fetch branches
│   ├── pulls.go     # Fetch pull requests
│   ├── issues.go    # Fetch issues
│   └── reviews.go   # Fetch code reviews
├── llm/             # LLM client for GitHub Models API
│   ├── client.go    # HTTP client for chat completions
│   ├── generator.go # Summary generation methods
│   ├── prompts.go   # YAML prompt loading (with //go:embed) and variable rendering
│   └── prompts/     # YAML prompt templates (embedded in binary)
├── report/          # Report generation
│   ├── generator.go # Main generation logic, parallel data collection
│   └── markdown.go  # Markdown formatting functions
├── types/           # Data structures (Author, Commit, Branch, PullRequest, Issue, Stats, Report)
├── errors/          # Custom error types
├── logger/          # Structured logging to stderr
└── utils/           # Worker pool, helper functions
```

## Important Implementation Details

### Authentication

- Uses GitHub CLI (`gh`) for authentication via `github.com/cli/go-gh/v2`
- Automatically retrieves token with `gh auth token`
- Same token works for both GitHub REST API and GitHub Models API

### Error Handling

- **Fail Fast**: Critical errors (auth, repo not found) stop execution immediately
- **Graceful Degradation**: Non-critical errors (AI summaries) use fallback text and continue
- Custom error types in [internal/errors/errors.go](internal/errors/errors.go):
  - `ErrGitHubAuth` - Authentication failures
  - `ErrGitHubAPI` - API request failures
  - `ErrRepoNotFound` - Repository not found
  - `ErrInvalidParams` - Invalid parameters
  - `ErrLLMAPI` - LLM API failures

### Retry Logic

- GitHub client implements exponential backoff for transient errors (see [internal/github/client.go:31](internal/github/client.go#L31))
- Retries network errors and 5xx status codes
- Does NOT retry 403 (rate limit) or 404 (not found)

### Bot Detection

- Filters bot accounts when `--exclude-bots` flag is used
- Checks for `[bot]` suffix and known bot names (github-actions, dependabot, renovate)
- Bot detection in [internal/github/client.go:96](internal/github/client.go#L96)

### YAML Prompt System

- Prompts located in [internal/llm/prompts/](internal/llm/prompts/)
- **Embedded in binary** using Go's `//go:embed` directive for self-contained distribution
- **External override support**: Files in `internal/llm/prompts/` take precedence over embedded ones (development/testing)
- Variable substitution using `{{variable_name}}` syntax
- Supports multi-language output via `{{language}}` variable
- Three main prompts:
  - `overall_summary.prompt.yml` - Overall repository activity
  - `branch_summary.prompt.yml` - Per-branch activity
  - `pr_summary.prompt.yml` - Pull request summaries

### Parallel Processing

- Worker pools limit concurrent LLM requests (max 5 workers by default)
- GitHub API calls use errgroups for parallel collection
- Helper function `utils.ProcessInParallel` in [internal/utils/pool.go](internal/utils/pool.go)

## Testing

### Unit Tests

- Standard Go test files (`*_test.go`)
- Example: [internal/types/author_test.go](internal/types/author_test.go), [internal/llm/prompts_test.go](internal/llm/prompts_test.go)

### Integration Tests

- Located in [test/integration/](test/integration/)
- Use mock clients: `mock_github.go`, `mock_llm.go`
- End-to-end test in [e2e_test.go](test/integration/e2e_test.go)
- Run with `just test-integration` or `just test-all`

### Testing with Real Data

```bash
# Test with actual repository
./bin/gh-repomon --repo hazadus/gh-repomon --days 7

# Test without AI summaries (faster)
./bin/gh-repomon --repo hazadus/gh-repomon --days 7 --no-ai

# Test language support
./bin/gh-repomon --repo hazadus/gh-repomon --days 7 --language russian
```

## Common Development Tasks

### Adding a New Command-Line Flag

1. Add variable in [cmd/repomon/main.go](cmd/repomon/main.go) (around line 17)
2. Register flag in `init()` function (around line 37)
3. Pass to `report.Options` struct (around line 123)
4. Use in report generation logic

### Modifying AI Prompts

1. Edit YAML files in [internal/llm/prompts/](internal/llm/prompts/)
2. **For development**: Changes take effect immediately (external files override embedded ones)
3. **For production**: Rebuild binary to embed changes: `just build`
4. Adjust `modelParameters.temperature` for creativity vs consistency
5. Add/modify variables in prompt content
6. Update variable map in [internal/llm/generator.go](internal/llm/generator.go)

**Testing prompt changes:**
```bash
# Edit prompt
vim internal/llm/prompts/overall_summary.prompt.yml

# Test immediately without rebuilding
go run ./cmd/repomon --repo test/repo --days 1

# Rebuild to embed changes
just build

# Verify embedded version works
cp bin/gh-repomon /tmp/ && cd /tmp/ && ./gh-repomon --repo test/repo --days 1
```

### Adding New Statistics

1. Define fields in [internal/types/stats.go](internal/types/stats.go)
2. Calculate in `calculateOverallStats()` or `calculateAuthorStats()` in [internal/report/markdown.go](internal/report/markdown.go)
3. Format in appropriate markdown section function

### Adding New Data Sources

1. Add method to `GitHubClient` interface in [internal/report/generator.go:18](internal/report/generator.go#L18)
2. Implement in [internal/github/](internal/github/)
3. Add parallel fetch in `collectData()` method (around line 222)
4. Update `types.ReportData` struct
5. Add markdown formatting section

## Dependencies

Key dependencies (see [go.mod](go.mod)):
- `github.com/cli/go-gh/v2` - GitHub CLI Go library
- `github.com/spf13/cobra` - CLI framework
- `golang.org/x/sync/errgroup` - Parallel execution with error handling
- `gopkg.in/yaml.v3` - YAML parsing for prompts

## Release Process

The project uses GoReleaser for automated releases:
1. Update [CHANGELOG.md](CHANGELOG.md)
2. Create and push git tag: `git tag vX.Y.Z && git push origin vX.Y.Z`
3. GitHub Actions workflow automatically builds and creates release
4. Binaries built for Linux, macOS, Windows (amd64, arm64)

## Code Conventions

- Follow standard Go conventions (gofmt, golint)
- Use descriptive error messages with context
- Log progress to stderr, output report to stdout
- Write tests for new functionality
- Document complex logic with comments
- Keep functions focused and small
