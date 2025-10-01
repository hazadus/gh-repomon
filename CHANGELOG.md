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

## [1.0.0] - TBD

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

[Unreleased]: https://github.com/hazadus/gh-repomon/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/hazadus/gh-repomon/releases/tag/v1.0.0
