# Contributing to gh-repomon

Thank you for your interest in contributing to gh-repomon! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Development Workflow](#development-workflow)
- [Code Style](#code-style)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)

## Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inclusive experience for everyone. We expect all contributors to:

- Be respectful and considerate
- Welcome newcomers and help them learn
- Focus on what is best for the community
- Show empathy towards others

### Unacceptable Behavior

- Harassment, discrimination, or offensive comments
- Personal attacks or trolling
- Publishing others' private information
- Other conduct inappropriate in a professional setting

## Getting Started

### Prerequisites

Before contributing, make sure you have:

- Go 1.21 or higher installed
- GitHub CLI (`gh`) installed and authenticated
- Git configured with your name and email
- (Optional) Just task runner installed
- (Optional) golangci-lint for code quality checks

### Fork and Clone

1. **Fork the repository** on GitHub
2. **Clone your fork:**
   ```bash
   git clone https://github.com/YOUR-USERNAME/gh-repomon.git
   cd gh-repomon
   ```
3. **Add upstream remote:**
   ```bash
   git remote add upstream https://github.com/hazadus/gh-repomon.git
   ```

## Development Setup

### Install Dependencies

```bash
# Download Go dependencies
go mod download

# Install Just (optional but recommended)
# macOS
brew install just

# Linux
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash

# Install golangci-lint (optional)
# macOS
brew install golangci-lint

# Linux/macOS via script
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

### Verify Setup

```bash
# Build the project
just build

# Run tests
just test

# Run the binary
./bin/gh-repomon --help
```

## Development Workflow

### 1. Create a Branch

Create a descriptive branch for your work:

```bash
# Update main branch
git checkout main
git pull upstream main

# Create feature branch
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/issue-number-description
```

### 2. Make Changes

- Write clean, well-documented code
- Follow the existing code style
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run unit tests
just test

# Run all tests (including integration)
just test-all

# Check code formatting
go fmt ./...

# Run linter (if installed)
golangci-lint run
```

### 4. Commit Changes

Write clear, concise commit messages:

```bash
git add .
git commit -m "feat: add support for custom date ranges

- Added --from and --to flags
- Updated documentation
- Added tests for date parsing"
```

**Commit Message Format:**
```
<type>: <subject>

<body>

<footer>
```

**Types:**
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `test:` Adding or updating tests
- `refactor:` Code refactoring
- `perf:` Performance improvements
- `chore:` Maintenance tasks

### 5. Push Changes

```bash
git push origin feature/your-feature-name
```

### 6. Create Pull Request

1. Go to your fork on GitHub
2. Click "New Pull Request"
3. Select your branch
4. Fill out the PR template
5. Submit the PR

## Code Style

### Go Code Style

Follow standard Go conventions:

- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use meaningful variable names
- Add comments for exported functions
- Keep functions small and focused

**Example:**

```go
// GenerateReport creates a repository activity report for the specified period.
// It collects data from GitHub, optionally generates AI summaries, and returns
// a formatted markdown report.
func (g *Generator) GenerateReport(opts Options) (string, error) {
    // Implementation
}
```

### File Organization

- One package per directory
- Related functionality grouped together
- Clear import organization:
  ```go
  import (
      // Standard library
      "fmt"
      "time"

      // External packages
      "github.com/spf13/cobra"

      // Internal packages
      "github.com/hazadus/gh-repomon/internal/types"
  )
  ```

### Error Handling

- Always check and handle errors
- Provide context in error messages
- Use custom error types where appropriate

```go
if err != nil {
    return fmt.Errorf("failed to fetch commits: %w", err)
}
```

### Documentation

- Add godoc comments for all exported items
- Include examples where helpful
- Update user documentation for user-facing changes

## Testing

### Writing Tests

- Place tests in `*_test.go` files
- Use table-driven tests where appropriate
- Test both success and failure cases
- Mock external dependencies

**Example:**

```go
func TestFormatDate(t *testing.T) {
    tests := []struct {
        name     string
        input    time.Time
        expected string
    }{
        {
            name:     "standard date",
            input:    time.Date(2025, 9, 30, 15, 30, 0, 0, time.UTC),
            expected: "2025-09-30 15:30",
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := formatDate(tt.input)
            if result != tt.expected {
                t.Errorf("got %q, want %q", result, tt.expected)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
just test

# Run with coverage
just test-coverage

# Run integration tests
just test-integration

# Run specific test
go test -v ./internal/report -run TestFormatDate

# Run with race detector
go test -race ./...
```

### Test Coverage

Aim for:
- 70%+ overall coverage
- 90%+ for critical paths
- 100% for utility functions

View coverage:
```bash
just test-coverage
# Opens coverage.html in browser
```

## Pull Request Process

### Before Submitting

- [ ] Tests pass locally
- [ ] Code is formatted (`go fmt`)
- [ ] Linter passes (if using golangci-lint)
- [ ] Documentation is updated
- [ ] Commit messages are clear
- [ ] Branch is up to date with main

### PR Checklist

When creating a PR, ensure:

- [ ] **Title** clearly describes the change
- [ ] **Description** explains what and why
- [ ] **Tests** are included and passing
- [ ] **Documentation** is updated if needed
- [ ] **Breaking changes** are noted (if any)
- [ ] **Issue** is linked (if applicable)

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
Describe testing done

## Checklist
- [ ] Tests pass
- [ ] Code formatted
- [ ] Documentation updated
- [ ] No breaking changes (or documented)

## Related Issues
Fixes #123
```

### Review Process

1. Maintainer will review your PR
2. Address any feedback
3. Make requested changes
4. Push updates
5. PR will be merged when approved

### After Merge

- Delete your feature branch
- Update your local main:
  ```bash
  git checkout main
  git pull upstream main
  ```

## Reporting Bugs

### Before Reporting

- Check [existing issues](https://github.com/hazadus/gh-repomon/issues)
- Try the latest version
- Gather reproduction steps
- Collect error messages and logs

### Bug Report Template

```markdown
**Describe the bug**
Clear description of the bug

**To Reproduce**
Steps to reproduce:
1. Run command: `gh-repomon --repo ... --days 7`
2. See error: ...

**Expected behavior**
What should happen

**Actual behavior**
What actually happens

**Environment**
- OS: macOS 14.0
- Go version: 1.21
- gh-repomon version: 1.0.0
- gh CLI version: 2.35.0

**Additional context**
Any other relevant information
```

## Suggesting Features

### Feature Request Template

```markdown
**Is your feature request related to a problem?**
Description of the problem

**Describe the solution you'd like**
Clear description of desired functionality

**Describe alternatives you've considered**
Other approaches you've thought about

**Additional context**
Mockups, examples, or other details
```

### Discussion First

For major features:
1. Open an issue for discussion
2. Wait for maintainer feedback
3. Get approval before implementing
4. Proceed with implementation

## Development Tips

### Using Just Commands

```bash
# See all available commands
just --list

# Build the project
just build

# Run tests
just test

# Format code
just format

# Clean build artifacts
just clean

# Run the application
just run -- --repo owner/repo --days 7
```

### Debugging

```bash
# Run with verbose output
./bin/gh-repomon --repo owner/repo --days 7 --verbose

# Use delve debugger
dlv debug ./cmd/repomon -- --repo owner/repo --days 7

# Print debug logs
go run ./cmd/repomon --repo owner/repo --days 7 2> debug.log
```

### Local Testing

Test with a real repository:

```bash
# Build and test
just build
./bin/gh-repomon --repo hazadus/gh-repomon --days 1

# Compare with expected output
./bin/gh-repomon --repo hazadus/gh-repomon --days 1 > actual.md
diff expected.md actual.md
```

## Project Structure

Understanding the codebase:

```
gh-repomon/
├── cmd/repomon/          # CLI entry point
├── internal/             # Private code
│   ├── github/           # GitHub API client
│   ├── llm/              # AI/LLM client
│   ├── report/           # Report generation
│   ├── types/            # Data types
│   ├── logger/           # Logging
│   ├── errors/           # Error types
│   └── utils/            # Utilities
├── test/                 # Tests
├── docs/                 # Documentation
└── Justfile              # Task runner
```

See [Architecture](architecture.md) for detailed design.

## Questions?

- 📖 Read the [documentation](../README.md)
- 🐛 Check [existing issues](https://github.com/hazadus/gh-repomon/issues)
- 💬 [Open a discussion](https://github.com/hazadus/gh-repomon/discussions)
- ✉️ Contact maintainers

## Recognition

Contributors will be:
- Listed in release notes
- Mentioned in CHANGELOG.md
- Credited in commit history

Thank you for contributing! 🎉

---

[Back to README](../README.md)
