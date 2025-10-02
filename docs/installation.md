# Installation Guide

This guide provides detailed instructions for installing and setting up gh-repomon.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation Methods](#installation-methods)
- [Setting Up GitHub CLI](#setting-up-github-cli)
- [Configuring GitHub Models Access](#configuring-github-models-access)
- [Verification](#verification)
- [Updating](#updating)
- [Uninstallation](#uninstallation)

## Prerequisites

### Required

- **Go 1.21 or higher** (if building from source)
  - Check your version: `go version`
  - Install from [golang.org](https://golang.org/dl/)

- **GitHub CLI (`gh`)**
  - Check if installed: `gh --version`
  - Install instructions below

### Optional (for development)

- **Just** - Command runner (alternative to Make)
  - Install: `brew install just` (macOS) or see [just.systems](https://just.systems/)

- **golangci-lint** - Go linter
  - Install: `brew install golangci-lint` (macOS) or see [golangci-lint.run](https://golangci-lint.run/usage/install/)

## Installation Methods

### Method 1: Via Go Install (Recommended)

The simplest way to install gh-repomon:

```bash
go install github.com/hazadus/gh-repomon/cmd/repomon@latest
```

This will install the binary to `$GOPATH/bin` (usually `~/go/bin`). Make sure this directory is in your `PATH`:

```bash
# Add to your ~/.bashrc, ~/.zshrc, or equivalent
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Method 2: As GitHub CLI Extension

Install as a GitHub CLI extension for seamless integration:

```bash
gh extension install hazadus/gh-repomon
```

After installation, you can use it as:

```bash
gh repomon --repo owner/repository --days 7
```

**Note:** The binary is fully self-contained with all prompt files embedded, so it works immediately after installation without requiring any additional files or configuration.

### Method 3: From Source

Clone and build manually:

```bash
# Clone the repository
git clone https://github.com/hazadus/gh-repomon.git
cd gh-repomon

# Install dependencies
go mod download

# Build using Just
just build

# Or build manually
go build -o bin/gh-repomon ./cmd/repomon

# Add to PATH or move to a location in PATH
mv bin/gh-repomon /usr/local/bin/
```

### Method 4: Download Pre-built Binary

Download the latest release for your platform from the [releases page](https://github.com/hazadus/gh-repomon/releases):

```bash
# Example for macOS (Apple Silicon)
curl -L https://github.com/hazadus/gh-repomon/releases/latest/download/gh-repomon-darwin-arm64 -o gh-repomon
chmod +x gh-repomon
mv gh-repomon /usr/local/bin/
```

**Note:** Pre-built binaries are fully self-contained with all prompt files embedded. No additional files or dependencies required.

Available platforms:
- `gh-repomon-linux-amd64` - Linux (Intel/AMD)
- `gh-repomon-linux-arm64` - Linux (ARM)
- `gh-repomon-darwin-amd64` - macOS (Intel)
- `gh-repomon-darwin-arm64` - macOS (Apple Silicon)
- `gh-repomon-windows-amd64.exe` - Windows (Intel/AMD)
- `gh-repomon-windows-arm64.exe` - Windows (ARM)

## Setting Up GitHub CLI

### Installing GitHub CLI

**macOS:**
```bash
brew install gh
```

**Windows:**
```bash
winget install --id GitHub.cli
# or
scoop install gh
```

**Linux (Debian/Ubuntu):**
```bash
type -p curl >/dev/null || (sudo apt update && sudo apt install curl -y)
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update
sudo apt install gh -y
```

**Linux (Fedora/CentOS/RHEL):**
```bash
sudo dnf install gh
```

For other platforms, see [GitHub CLI installation guide](https://github.com/cli/cli#installation).

### Authenticating GitHub CLI

Authenticate with your GitHub account:

```bash
gh auth login
```

Follow the interactive prompts:
1. Choose "GitHub.com"
2. Choose "HTTPS" (recommended) or "SSH"
3. Authenticate via web browser or paste an authentication token

Verify authentication:

```bash
gh auth status
```

## Configuring GitHub Models Access

gh-repomon uses GitHub Models API for AI-powered summaries. Your GitHub token needs access to this service.

### Getting Access

1. Visit [GitHub Models Marketplace](https://github.com/marketplace/models)
2. GitHub Models is currently in beta and may require waitlist approval
3. Once you have access, your existing GitHub token should work

### Verify Token Access

Your token is automatically used from the GitHub CLI authentication. Test it:

```bash
gh api https://models.inference.ai.azure.com/chat/completions --method GET
```

If you see authorization errors, you may need to:

1. Generate a new token with appropriate scopes
2. Re-authenticate GitHub CLI: `gh auth refresh`

### Using Custom Token

If you need to use a different token, set the environment variable:

```bash
export GITHUB_TOKEN="your_token_here"
```

Add to your shell profile (`~/.bashrc`, `~/.zshrc`) to persist.

## Verification

Verify your installation:

### 1. Check Binary

```bash
# If installed via go install or built from source
gh-repomon --version

# If installed as GitHub CLI extension
gh repomon --version
```

### 2. Check Help

```bash
gh-repomon --help
```

You should see the help message with all available flags.

### 3. Test with a Repository

Run a simple test on a public repository:

```bash
gh-repomon --repo hazadus/gh-repomon --days 1
```

If everything is set up correctly, you should see:
- Progress messages in stderr
- A markdown report in stdout

## Updating

### Go Install Method

```bash
go install github.com/hazadus/gh-repomon/cmd/repomon@latest
```

### GitHub CLI Extension

```bash
gh extension upgrade repomon
```

### From Source

```bash
cd gh-repomon
git pull
just build
```

## Uninstallation

### Go Install Method

```bash
rm $(which gh-repomon)
# or
rm ~/go/bin/gh-repomon
```

### GitHub CLI Extension

```bash
gh extension remove repomon
```

### From Source

```bash
rm /usr/local/bin/gh-repomon
# And remove the cloned directory
rm -rf ~/path/to/gh-repomon
```

## Troubleshooting

### Command Not Found

If you get "command not found" after installation:

1. Make sure `$GOPATH/bin` is in your PATH
2. Reload your shell: `source ~/.bashrc` or `source ~/.zshrc`
3. Check the binary location: `which gh-repomon`

### Permission Denied

If you get permission errors when installing to `/usr/local/bin`:

```bash
sudo mv bin/gh-repomon /usr/local/bin/
```

Or install to a user directory:

```bash
mkdir -p ~/.local/bin
mv bin/gh-repomon ~/.local/bin/
export PATH="$PATH:$HOME/.local/bin"
```

### GitHub Authentication Issues

If you get authentication errors:

1. Check auth status: `gh auth status`
2. Re-authenticate: `gh auth login`
3. Refresh token: `gh auth refresh`

For more troubleshooting, see [Troubleshooting Guide](troubleshooting.md).

## Next Steps

- Read the [Usage Guide](usage.md) to learn about all available options
- Check out [Examples](examples.md) for common use cases
- Learn about [Customizing Prompts](prompts.md) for AI summaries

---

[Back to README](../README.md)
