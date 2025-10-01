# gh-repomon

[![CI Status](https://github.com/hazadus/gh-repomon/workflows/CI/badge.svg)](https://github.com/hazadus/gh-repomon/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**GitHub Repository Activity Monitor** - A CLI tool for generating comprehensive activity reports for GitHub repositories with AI-powered summaries.

## ✨ Key Features

- 📊 **Comprehensive Activity Reports** - Track commits, pull requests, issues, and code reviews
- 🤖 **AI-Powered Summaries** - Generate intelligent summaries using GitHub Models API
- 📈 **Author Statistics** - Detailed breakdown of contributions by author
- 🌿 **Branch Analysis** - Activity tracking across all active branches
- 🔍 **Flexible Filtering** - Filter by date range, user, and more
- 🌍 **Multi-language Support** - Generate reports in different languages
- ⚡ **Fast & Efficient** - Parallel processing for quick report generation
- 📝 **Markdown Output** - Readable reports ready to share

## 🚀 Requirements

- Go 1.21 or higher
- [GitHub CLI](https://cli.github.com/) (`gh`) installed and authenticated
- GitHub token with access to [GitHub Models](https://github.com/marketplace/models) (for AI summaries)

## 📦 Installation

### Via Go Install

```bash
go install github.com/hazadus/gh-repomon/cmd/repomon@latest
```

### As GitHub CLI Extension

```bash
gh extension install hazadus/gh-repomon
```

### From Source

```bash
git clone https://github.com/hazadus/gh-repomon.git
cd gh-repomon
just build
# Binary will be in ./bin/gh-repomon
```

## 🎯 Quick Start

Generate a report for the last 7 days:

```bash
gh-repomon --repo owner/repository --days 7
```

Or if installed as a GitHub CLI extension:

```bash
gh repomon --repo owner/repository --days 7
```

Generate a report for a specific date range:

```bash
gh-repomon --repo owner/repository --from 2025-09-01 --to 2025-09-30
```

Filter by specific user:

```bash
gh-repomon --repo owner/repository --days 7 --user username
```

Generate report in Russian:

```bash
gh-repomon --repo owner/repository --days 7 --language russian
```

Exclude bot activity:

```bash
gh-repomon --repo owner/repository --days 7 --exclude-bots
```

Save report to file:

```bash
gh-repomon --repo owner/repository --days 7 > report.md
```

## 📚 Documentation

- [Installation Guide](docs/installation.md) - Detailed setup instructions
- [Usage Guide](docs/usage.md) - Comprehensive flag descriptions and examples
- [Examples](docs/examples.md) - Real-world use cases
- [Troubleshooting](docs/troubleshooting.md) - Common issues and solutions
- [Architecture](docs/architecture.md) - System design overview
- [Contributing](docs/contributing.md) - How to contribute
- [Prompts Guide](docs/prompts.md) - Customizing AI prompts

## 🛠️ Development

### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/hazadus/gh-repomon.git
cd gh-repomon

# Install dependencies
go mod download

# Install Just (task runner)
# macOS
brew install just

# Linux
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin
```

### Available Commands

```bash
just run           # Run the application
just build         # Build binary
just test          # Run unit tests
just test-all      # Run all tests (unit + integration)
just format        # Format code
just clean         # Clean build artifacts
```

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](docs/contributing.md) for details on our code of conduct and the process for submitting pull requests.

## �� License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [gh-standup](https://github.com/sgoedecke/gh-standup)
- Built with [GitHub CLI Go library](https://github.com/cli/go-gh)
- Powered by [GitHub Models](https://github.com/marketplace/models)

## 📧 Support

- 🐛 [Report a Bug](https://github.com/hazadus/gh-repomon/issues/new?labels=bug)
- 💡 [Request a Feature](https://github.com/hazadus/gh-repomon/issues/new?labels=enhancement)
- 📖 [Read the Docs](docs/)

