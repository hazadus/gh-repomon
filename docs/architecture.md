# Architecture

System design and architecture overview for gh-repomon.

## Table of Contents

- [Overview](#overview)
- [High-Level Architecture](#high-level-architecture)
- [Core Components](#core-components)
- [Data Flow](#data-flow)
- [Module Structure](#module-structure)
- [Design Decisions](#design-decisions)
- [Extension Points](#extension-points)

## Overview

gh-repomon is a CLI tool built in Go that analyzes GitHub repository activity and generates comprehensive markdown reports with AI-powered summaries. The architecture is designed to be:

- **Modular**: Clear separation of concerns
- **Extensible**: Easy to add new features
- **Testable**: Components can be tested independently
- **Performant**: Parallel processing where possible
- **Maintainable**: Clean code structure

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                            │
│                   (cmd/repomon/main.go)                      │
│  - Parse arguments (Cobra)                                   │
│  - Orchestrate workflow                                      │
│  - Handle output                                             │
└────────────┬────────────────────────────────┬───────────────┘
             │                                │
             ▼                                ▼
┌────────────────────────┐        ┌──────────────────────────┐
│   GitHub API Client    │        │     LLM Client           │
│  (internal/github)     │        │   (internal/llm)         │
│                        │        │                          │
│  - REST API calls      │        │  - GitHub Models API     │
│  - Data fetching       │        │  - Prompt management     │
│  - Rate limiting       │        │  - Summary generation    │
│  - Bot detection       │        │  - YAML prompts          │
└────────┬───────────────┘        └──────────┬───────────────┘
         │                                   │
         │         ┌─────────────────────────┘
         │         │
         ▼         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Report Generator                          │
│                  (internal/report)                           │
│                                                              │
│  - Data aggregation                                          │
│  - Statistics calculation                                    │
│  - Markdown formatting                                       │
│  - AI summary integration                                    │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
            ┌────────────────┐
            │   Type System   │
            │ (internal/types)│
            │                 │
            │  - Data models  │
            │  - Structures   │
            └─────────────────┘
```

## Core Components

### 1. CLI Layer (`cmd/repomon/`)

**Responsibility:** User interface and workflow orchestration

**Key Files:**
- `main.go` - Entry point, argument parsing, main execution flow

**Technologies:**
- [Cobra](https://github.com/spf13/cobra) - CLI framework

**Flow:**
1. Parse command-line arguments
2. Validate inputs
3. Initialize GitHub and LLM clients
4. Create report generator
5. Generate and output report
6. Handle errors

### 2. GitHub API Client (`internal/github/`)

**Responsibility:** Interact with GitHub REST API

**Key Files:**
- `client.go` - Client initialization, authentication
- `commits.go` - Fetch commits and stats
- `branches.go` - Fetch branches and activity
- `pulls.go` - Fetch pull requests
- `issues.go` - Fetch issues
- `reviews.go` - Fetch code reviews

**Key Features:**
- Uses [go-gh](https://github.com/cli/go-gh) library
- Automatic token management via GitHub CLI
- Bot detection and filtering
- Parallel request processing
- User caching

**Design Patterns:**
- Client pattern for API abstraction
- Worker pool for parallel processing

### 3. LLM Client (`internal/llm/`)

**Responsibility:** Generate AI summaries via GitHub Models API

**Key Files:**
- `client.go` - Client initialization, API calls
- `generator.go` - Summary generation methods
- `prompts.go` - YAML prompt loading and rendering
- `prompts/*.yml` - Prompt templates

**Key Features:**
- YAML-based prompt system
- Template variable substitution
- Graceful fallback on errors
- **Automatic retry with exponential backoff**
- **Intelligent rate limit handling**
- Timeout handling (30s per request)

**Retry Logic:**
- Detects rate limit errors (HTTP 429)
- Extracts wait time from error response
- Automatically waits and retries (up to 3 attempts)
- Exponential backoff for server errors (5xx)
- Logs retry progress and success

**Design Patterns:**
- Template pattern for prompts
- Strategy pattern for different summary types
- Retry pattern with intelligent backoff

### 4. Report Generator (`internal/report/`)

**Responsibility:** Aggregate data and generate markdown reports

**Key Files:**
- `generator.go` - Main generation logic, data collection
- `markdown.go` - Markdown formatting functions

**Key Features:**
- Statistics calculation
- Author activity aggregation
- Branch analysis
- Markdown formatting
- AI summary integration

**Design Patterns:**
- Builder pattern for report construction
- Visitor pattern for data aggregation

### 5. Type System (`internal/types/`)

**Responsibility:** Define data structures

**Key Files:**
- `author.go` - Author/user information
- `commit.go` - Commit data
- `branch.go` - Branch information
- `pull_request.go` - PR data
- `issue.go` - Issue data
- `stats.go` - Statistics structures
- `report.go` - Report data structure

**Design Principles:**
- Immutable where possible
- Rich domain models
- Clear field naming

### 6. Supporting Modules

#### Logger (`internal/logger/`)
- Structured logging to stderr
- Progress indicators
- Error reporting

#### Errors (`internal/errors/`)
- Custom error types
- Error wrapping
- Context preservation

#### Utils (`internal/utils/`)
- Worker pool implementation
- Helper functions
- Common utilities

## Data Flow

### 1. Initialization Phase

```
User Input → CLI Parsing → Validation → Client Creation
```

### 2. Data Collection Phase

```
GitHub Client → [Parallel Requests]
  ├─ Fetch Branches
  ├─ Fetch Commits (per branch)
  ├─ Fetch Commit Stats
  ├─ Fetch Pull Requests
  ├─ Fetch Issues
  └─ Fetch Reviews

        ↓

  Aggregation

        ↓

  ReportData Structure
```

### 3. AI Enhancement Phase (Optional)

```
ReportData → LLM Client
  ├─ Load Prompts (YAML)
  ├─ Render with Variables
  ├─ [Parallel Requests]
  │   ├─ Overall Summary
  │   ├─ Branch Summaries
  │   └─ PR Summaries
  └─ Integrate Results

        ↓

  Enhanced ReportData
```

### 4. Report Generation Phase

```
Enhanced ReportData → Report Generator
  ├─ Calculate Statistics
  ├─ Aggregate Author Data
  ├─ Format Markdown Sections
  │   ├─ Header
  │   ├─ Summary Stats
  │   ├─ Overall Summary (AI)
  │   ├─ Branches (with AI)
  │   ├─ Pull Requests (with AI)
  │   ├─ Issues
  │   ├─ Code Reviews
  │   └─ Author Activity
  └─ Combine Sections

        ↓

  Markdown Report → stdout
```

## Module Structure

```
gh-repomon/
├── cmd/
│   └── repomon/          # CLI entry point
│       └── main.go
│
├── internal/             # Private application code
│   ├── github/           # GitHub API client
│   │   ├── client.go
│   │   ├── commits.go
│   │   ├── branches.go
│   │   ├── pulls.go
│   │   ├── issues.go
│   │   └── reviews.go
│   │
│   ├── llm/              # LLM client
│   │   ├── client.go
│   │   ├── generator.go
│   │   ├── prompts.go    # Embeds prompts via //go:embed
│   │   └── prompts/      # YAML prompt templates (embedded in binary)
│   │       ├── overall_summary.prompt.yml
│   │       ├── branch_summary.prompt.yml
│   │       └── pr_summary.prompt.yml
│   │
│   ├── report/           # Report generation
│   │   ├── generator.go
│   │   └── markdown.go
│   │
│   ├── types/            # Data structures
│   │   ├── author.go
│   │   ├── commit.go
│   │   ├── branch.go
│   │   ├── pull_request.go
│   │   ├── issue.go
│   │   ├── stats.go
│   │   └── report.go
│   │
│   ├── logger/           # Logging
│   │   └── logger.go
│   │
│   ├── errors/           # Error types
│   │   └── errors.go
│   │
│   └── utils/            # Utilities
│       └── pool.go
│
├── test/                 # Tests
│   └── integration/      # Integration tests
│       ├── e2e_test.go
│       ├── mock_github.go
│       └── mock_llm.go
│
├── docs/                 # Documentation
├── bin/                  # Compiled binaries (gitignored)
├── go.mod                # Go module definition
├── go.sum                # Dependency checksums
└── Justfile              # Task runner commands
```

## Design Decisions

### 1. Why Go?

**Reasons:**
- Fast compilation and execution
- Easy cross-platform builds
- Strong stdlib for HTTP/JSON
- Good concurrency support (goroutines)
- Single binary distribution
- GitHub CLI ecosystem

### 2. Why GitHub CLI Library?

**Reasons:**
- Automatic authentication
- Token management
- Consistent with GitHub CLI ecosystem
- Well-maintained by GitHub
- Easy for users (already have `gh` installed)

### 3. Why YAML Prompts?

**Reasons:**
- Easy to read and edit
- Non-developers can customize
- Version controllable
- Supports multi-line text naturally
- Template variable support
- Embedded in binary via `//go:embed` for self-contained distribution
- External override support for development and customization

**Alternatives Considered:**
- Hardcoded prompts (too rigid)
- JSON (less readable)
- Separate .txt files (harder to manage metadata)

**Implementation:**
Prompts are embedded directly into the binary using Go's `//go:embed` directive. This makes the binary fully self-contained while still allowing developers to override prompts by placing external files in `internal/llm/prompts/` during development.

### 4. Why Markdown Output?

**Reasons:**
- Universal format
- Easy to read
- Easy to convert (pandoc)
- GitHub-native
- Can be rendered in browsers, editors, etc.

**Alternatives Considered:**
- HTML (not universal, harder to read raw)
- JSON (not human-friendly)
- PDF (not easily editable)

### 5. Parallel Processing

**Where Applied:**
- GitHub API requests (branches, PRs, issues)
- Commit stats fetching
- LLM summary generation

**Why:**
- Significantly faster report generation
- Better resource utilization
- GitHub API can handle concurrent requests

**Trade-offs:**
- Complexity
- Rate limit considerations
- Memory usage

### 6. Error Handling Strategy

**Approach:**
- Fail fast on critical errors (auth, repo not found)
- Graceful degradation on non-critical errors (AI summaries)
- Clear error messages with context
- Continue processing when possible

**Example:**
- If AI fails: Use fallback text, continue report generation
- If GitHub auth fails: Stop immediately with clear message

### 7. Separation of Concerns

**Principle:**
Each module has a single, well-defined responsibility.

**Benefits:**
- Easier testing
- Easier to understand
- Easier to modify
- Reusable components

## Extension Points

### 1. New Output Formats

To add HTML, JSON, or other formats:

1. Create `internal/report/html.go` (or similar)
2. Implement formatting functions
3. Add `--format` flag to CLI
4. Call appropriate formatter

### 2. New AI Providers

To support OpenAI API, Anthropic, etc.:

1. Create interface in `internal/llm/client.go`
2. Implement provider-specific client
3. Add provider selection logic
4. Update prompts if needed

### 3. New Data Sources

To support GitLab, Bitbucket, etc.:

1. Create `internal/gitlab/` (or similar)
2. Implement common interface
3. Abstract report generator
4. Add provider selection

### 4. Custom Metrics

To add custom statistics:

1. Define new fields in `internal/types/stats.go`
2. Add calculation in `internal/report/generator.go`
3. Add formatting in `internal/report/markdown.go`

### 5. Plugins

Future consideration for plugin system:

- Custom data collectors
- Custom formatters
- Custom AI providers
- Pre/post-processing hooks

## Performance Considerations

### 1. API Rate Limits

**GitHub API:**
- 5000 requests/hour (authenticated)
- Use conditional requests where possible
- Cache user data
- Parallel requests within limits

**GitHub Models API:**
- Rate limits vary by model (typically 10 requests per 60 seconds)
- **Automatic retry with intelligent backoff** (implemented in v1.1.0)
- Extracts wait time from error response
- Retries up to 3 times with proper delays
- Use worker pools to limit concurrency (max 5 concurrent requests)

### 2. Memory Usage

**Optimizations:**
- Stream large responses where possible
- Process data incrementally
- Limit concurrent goroutines
- Clear temporary data

### 3. Speed

**Optimizations:**
- Parallel API requests
- Concurrent AI summary generation
- Efficient data structures
- Minimal allocations

## Testing Strategy

### Unit Tests
- Individual functions
- Data transformations
- Formatting logic
- Prompt rendering

### Integration Tests
- Full pipeline with mocks
- End-to-end scenarios
- Error handling

### Manual Testing
- Real repositories
- Various sizes
- Different configurations
- Edge cases

## Security Considerations

1. **Token Handling:**
   - Never log tokens
   - Use GitHub CLI token management
   - Support environment variables

2. **Input Validation:**
   - Validate date formats
   - Sanitize repository names
   - Validate user inputs

3. **API Safety:**
   - Respect rate limits
   - Timeout protection
   - Error handling

4. **Output Safety:**
   - Escape markdown special characters
   - Handle malicious commit messages
   - Sanitize user-generated content

## Future Enhancements

Potential architecture improvements:

1. **Caching Layer:**
   - Cache GitHub API responses
   - Persist between runs
   - Faster repeated queries

2. **Configuration File:**
   - `.gh-repomon.yml` for defaults
   - Per-repository settings
   - Team preferences

3. **Plugin System:**
   - Custom data sources
   - Custom formatters
   - Custom AI providers

4. **Web UI:**
   - Browser-based interface
   - Interactive reports
   - Visualization

5. **Database Support:**
   - Store historical data
   - Trend analysis
   - Comparison over time

---

[Back to README](../README.md)
