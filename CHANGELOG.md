# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- JSON output format
- Verbose mode with detailed logging
- Configuration file support
- Caching of GitHub API responses
- HTML output with interactive charts
- Support for multiple repositories
- Comparison between different time periods

## [1.1.0] - 2025-10-02

### Added
- Automatic retry with exponential backoff for LLM API requests
- Rate limit detection (HTTP 429) with intelligent wait time extraction from error messages
- Up to 3 automatic retries for server errors (5xx)
- Enhanced logging with `Infof` and `Warningf` methods in logger

### Changed
- Improved LLM client error handling and reliability
- Better API failure recovery with automatic retry logic

### Documentation
- Updated [architecture.md](docs/architecture.md) with retry pattern and intelligent backoff details
- Revised [troubleshooting.md](docs/troubleshooting.md) with automatic retry behavior information

### Tests
- Added unit tests for wait time extraction from error messages
- Added tests for daily rate limit and shorter wait time handling

## [1.0.1] - 2025-10-02

### Changed
- Embedded YAML prompt files directly into binary using Go's `//go:embed` directive
- Binary is now fully self-contained and works without external prompt files
- External prompt files now override embedded ones during development for easier testing

### Added
- Unit tests for embedded prompt loading and external file override mechanism
- Documentation for prompt embedding approach in multiple docs files

### Fixed
- Improved nil pointer handling in LLM client initialization

### Documentation
- Updated [architecture.md](docs/architecture.md) with prompt embedding details
- Updated [installation.md](docs/installation.md) to highlight self-contained binary
- Updated [troubleshooting.md](docs/troubleshooting.md) with prompt embedding info
- Updated [contributing.md](docs/contributing.md) with prompt testing instructions
- Enhanced [prompts.md](docs/prompts.md) with embedded vs external prompts explanation

### Chore
- Updated `.clocignore` to exclude coverage files (`coverage.html`, `coverage.out`)

## [1.0.0] - 2025-10-01

### Added
- Initial release of gh-repomon
- GitHub API integration for repository activity monitoring
- Comprehensive activity reports including:
  - Branch activity with commits and statistics
  - Pull requests (open and updated)
  - Issues (open and closed)
  - Code reviews
  - Author statistics
- AI-powered summaries using GitHub Models API:
  - Overall repository summary
  - Branch-specific summaries
  - Pull request summaries
- Markdown report generation
- Command-line interface with Cobra
- Flexible time period selection:
  - Relative time (--days flag)
  - Absolute time range (--from and --to flags)
- Filtering options:
  - Filter by user
  - Exclude bot accounts
- Multi-language support for AI summaries
- Configurable AI model selection
- YAML-based prompt templates for AI generation
- GitHub CLI extension support
- Cross-platform binaries (Linux, macOS, Windows for AMD64 and ARM64)
- Comprehensive test suite:
  - Unit tests for core components
  - Integration tests for end-to-end scenarios
- Structured logging with progress indicators
- Error handling with graceful degradation
- Parallel data collection for improved performance
- Rate limiting for API calls

### Documentation
- Complete user guide
- Installation instructions
- Usage examples
- Architecture documentation
- Contributing guidelines
- Troubleshooting guide
- Prompt customization guide

[Unreleased]: https://github.com/hazadus/gh-repomon/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/hazadus/gh-repomon/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/hazadus/gh-repomon/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/hazadus/gh-repomon/releases/tag/v1.0.0
