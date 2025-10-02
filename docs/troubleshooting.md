# Troubleshooting Guide

Common issues and their solutions for gh-repomon.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Authentication Issues](#authentication-issues)
- [Runtime Errors](#runtime-errors)
- [Performance Issues](#performance-issues)
- [AI/LLM Issues](#aillm-issues)
- [Output Issues](#output-issues)
- [FAQ](#faq)

## Installation Issues

### Command Not Found After Installation

**Problem:**
```bash
gh-repomon --help
# bash: command not found: gh-repomon
```

**Solutions:**

1. **Check if binary is in PATH:**
   ```bash
   which gh-repomon
   ```

2. **Add Go bin to PATH:**
   ```bash
   # Add to ~/.bashrc or ~/.zshrc
   export PATH="$PATH:$(go env GOPATH)/bin"

   # Reload shell
   source ~/.bashrc  # or source ~/.zshrc
   ```

3. **Verify installation location:**
   ```bash
   ls -la ~/go/bin/gh-repomon
   # or
   ls -la /usr/local/bin/gh-repomon
   ```

4. **Reinstall:**
   ```bash
   go install github.com/hazadus/gh-repomon/cmd/repomon@latest
   ```

### Go Version Too Old

**Problem:**
```
go: module github.com/hazadus/gh-repomon requires go >= 1.21
```

**Solution:**

Update Go to version 1.21 or higher:
```bash
# Check current version
go version

# Update via package manager (macOS)
brew upgrade go

# Or download from golang.org
# https://golang.org/dl/
```

### Build Errors

**Problem:**
```
error: cannot find package
```

**Solution:**

1. **Update dependencies:**
   ```bash
   go mod download
   go mod tidy
   ```

2. **Clear module cache:**
   ```bash
   go clean -modcache
   go mod download
   ```

3. **Verify Go environment:**
   ```bash
   go env
   ```

## Authentication Issues

### GitHub CLI Not Authenticated

**Problem:**
```
Error: authentication failed
```

**Solution:**

1. **Check authentication status:**
   ```bash
   gh auth status
   ```

2. **Login to GitHub:**
   ```bash
   gh auth login
   ```

3. **Refresh token:**
   ```bash
   gh auth refresh
   ```

4. **Verify token scopes:**
   ```bash
   gh auth status --show-token
   ```

### GitHub Models API Access Denied

**Problem:**
```
Error: access denied to GitHub Models API
```

**Solution:**

1. **Check GitHub Models access:**
   - Visit [GitHub Models](https://github.com/marketplace/models)
   - Ensure you have beta access

2. **Regenerate token:**
   ```bash
   gh auth refresh -h github.com -s read:org,read:user
   ```

3. **Use environment variable:**
   ```bash
   export GITHUB_TOKEN="your_token_here"
   gh-repomon --repo owner/repo --days 7
   ```

### Repository Access Denied

**Problem:**
```
Error: repository not found or access denied
```

**Solution:**

1. **Verify repository exists:**
   ```bash
   gh repo view owner/repository
   ```

2. **Check permissions:**
   - Ensure you have read access to the repository
   - For private repos, verify token has correct scopes

3. **Try with full URL:**
   ```bash
   gh-repomon --repo https://github.com/owner/repository --days 7
   ```

## Runtime Errors

### Rate Limit Exceeded

**Problem:**
```
Error: API rate limit exceeded
```

**Solution:**

1. **Check rate limit status:**
   ```bash
   gh api rate_limit
   ```

2. **Wait for reset:**
   - Rate limits reset hourly
   - Check `X-RateLimit-Reset` in the error message

3. **Use authenticated requests:**
   - Authenticated requests have higher limits (5000/hour vs 60/hour)
   - Verify you're authenticated: `gh auth status`

4. **Reduce scope:**
   ```bash
   # Shorter time period
   gh-repomon --repo owner/repo --days 1

   # Disable AI
   gh-repomon --repo owner/repo --days 7 --no-ai
   ```

### Network Timeout

**Problem:**
```
Error: request timeout
```

**Solution:**

1. **Check internet connection:**
   ```bash
   ping github.com
   ```

2. **Retry the command:**
   - Temporary network issues often resolve themselves

3. **Check GitHub status:**
   - Visit [GitHub Status](https://www.githubstatus.com/)

4. **Use proxy if needed:**
   ```bash
   export HTTPS_PROXY=http://proxy.example.com:8080
   ```

### Invalid Date Format

**Problem:**
```
Error: invalid date format
```

**Solution:**

Use correct date format (YYYY-MM-DD):
```bash
# Correct
gh-repomon --repo owner/repo --from 2025-09-01 --to 2025-09-30

# Incorrect
gh-repomon --repo owner/repo --from 09/01/2025 --to 09/30/2025
```

### Repository Not Found

**Problem:**
```
Error: repository not found
```

**Solution:**

1. **Verify repository name:**
   ```bash
   # Correct format
   gh-repomon --repo owner/repository

   # Not: repository only
   gh-repomon --repo repository
   ```

2. **Check for typos:**
   ```bash
   # List your repositories
   gh repo list owner
   ```

3. **Verify access:**
   ```bash
   gh repo view owner/repository
   ```

## Performance Issues

### Slow Report Generation

**Problem:**
Report takes too long to generate.

**Solutions:**

1. **Disable AI summaries:**
   ```bash
   gh-repomon --repo owner/repo --days 7 --no-ai
   ```

2. **Use faster AI model:**
   ```bash
   gh-repomon --repo owner/repo --days 7 --model openai/gpt-4o-mini
   ```

3. **Reduce time period:**
   ```bash
   gh-repomon --repo owner/repo --days 1
   ```

4. **Filter by user:**
   ```bash
   gh-repomon --repo owner/repo --days 7 --user username
   ```

5. **Check verbose output for bottlenecks:**
   ```bash
   gh-repomon --repo owner/repo --days 7 --verbose 2> debug.log
   ```

### High Memory Usage

**Problem:**
Process consumes too much memory.

**Solution:**

1. **Process smaller time periods:**
   ```bash
   # Instead of 30 days, do 7 days at a time
   gh-repomon --repo owner/repo --days 7
   ```

2. **Disable AI:**
   ```bash
   gh-repomon --repo owner/repo --days 7 --no-ai
   ```

3. **Report issue:**
   - Large repositories may require optimization
   - [Report an issue](https://github.com/hazadus/gh-repomon/issues) with details

## AI/LLM Issues

### Prompt Files Not Found (Legacy Issue - Fixed)

**Problem:**
```
Error: failed to load prompt: failed to read prompt file internal/llm/prompts/overall_summary.prompt.yml: no such file or directory
```

**Solution:**

This issue has been fixed in recent versions. Prompt files are now embedded directly into the binary using Go's `//go:embed` directive, making the binary fully self-contained.

**If you still see this error:**
1. Update to the latest version: `gh extension upgrade repomon`
2. Or rebuild from source: `just build`
3. The error should not occur with version 1.1.0 or later

### AI Summaries Not Generated

**Problem:**
Report shows placeholders instead of AI summaries.

**Solution:**

1. **Check GitHub Models access:**
   ```bash
   gh api https://models.inference.ai.azure.com/
   ```

2. **Verify model name:**
   ```bash
   # Use correct model name
   gh-repomon --repo owner/repo --days 7 --model openai/gpt-4o
   ```

3. **Check verbose output:**
   ```bash
   gh-repomon --repo owner/repo --days 7 --verbose 2>&1 | grep -i "llm\|ai\|model"
   ```

4. **Try without AI first:**
   ```bash
   gh-repomon --repo owner/repo --days 7 --no-ai
   ```

### LLM API Timeout

**Problem:**
```
Error: LLM API request timeout
```

**Solution:**

1. **Use faster model:**
   ```bash
   gh-repomon --repo owner/repo --days 7 --model openai/gpt-4o-mini
   ```

2. **Retry:**
   - Temporary API issues may resolve

3. **Disable AI:**
   ```bash
   gh-repomon --repo owner/repo --days 7 --no-ai
   ```

### Wrong Language in Summaries

**Problem:**
AI summaries are in English when another language was requested.

**Solution:**

1. **Verify language flag:**
   ```bash
   gh-repomon --repo owner/repo --days 7 --language russian
   ```

2. **Override prompt for testing:**
   - Create `internal/llm/prompts/overall_summary.prompt.yml` with modified language instructions
   - External prompts override embedded ones during development

3. **Try different model:**
   ```bash
   gh-repomon --repo owner/repo --days 7 --language russian --model openai/gpt-4o
   ```

## Output Issues

### Broken Markdown Formatting

**Problem:**
Output has formatting issues.

**Solution:**

1. **Save to file and check:**
   ```bash
   gh-repomon --repo owner/repo --days 7 > report.md
   ```

2. **Verify with markdown linter:**
   ```bash
   markdownlint report.md
   ```

3. **View with markdown viewer:**
   ```bash
   glow report.md
   # or
   mdcat report.md
   ```

### Special Characters Causing Issues

**Problem:**
Commit messages or descriptions with special characters break formatting.

**Solution:**

This should be handled automatically. If you encounter issues:

1. **Report the issue:**
   - Include sample commit message/description
   - [Create an issue](https://github.com/hazadus/gh-repomon/issues)

2. **Workaround:**
   - Use `--no-ai` to skip AI processing
   - Post-process the markdown file

### No Output Generated

**Problem:**
Command runs but produces no output.

**Solution:**

1. **Check if there's any activity:**
   ```bash
   # Try longer period
   gh-repomon --repo owner/repo --days 30
   ```

2. **Check error output:**
   ```bash
   gh-repomon --repo owner/repo --days 7 2>&1
   ```

3. **Use verbose mode:**
   ```bash
   gh-repomon --repo owner/repo --days 7 --verbose
   ```

4. **Verify repository has activity:**
   ```bash
   gh api repos/owner/repo/commits --jq '.[0].commit.author.date'
   ```

## FAQ

### Q: Can I use gh-repomon without GitHub Models API?

**A:** Yes, use the `--no-ai` flag:
```bash
gh-repomon --repo owner/repo --days 7 --no-ai
```

### Q: How do I reduce API calls?

**A:** Several options:
- Shorter time periods: `--days 1` instead of `--days 30`
- Disable AI: `--no-ai`
- Filter by user: `--user username`
- Exclude bots: `--exclude-bots` (reduces noise, not API calls)

### Q: Why are some commits missing?

**A:** Possible reasons:
- Commits are older than the specified period
- Commits are from excluded bot accounts
- User filter is active
- Check verbose output for details

### Q: Can I customize AI prompts?

**A:** Yes! Prompts are embedded in the binary but can be overridden:
- **For development:** Place modified prompts in `internal/llm/prompts/` - they take precedence
- **For production:** Rebuild the binary after modifying prompts to embed the changes
See [Prompts Guide](prompts.md) for details.

### Q: How do I report a bug?

**A:** [Create an issue](https://github.com/hazadus/gh-repomon/issues/new?labels=bug) with:
- Command you ran
- Expected behavior
- Actual behavior
- Error messages (if any)
- Output from `--verbose` mode

### Q: Can I use gh-repomon for private repositories?

**A:** Yes, as long as:
- You have read access to the repository
- Your GitHub token has appropriate scopes
- You're authenticated via `gh auth login`

### Q: How do I get support?

**A:**
- Check this troubleshooting guide
- Read the [documentation](../README.md)
- Search [existing issues](https://github.com/hazadus/gh-repomon/issues)
- Create a [new issue](https://github.com/hazadus/gh-repomon/issues/new)

### Q: Is there a way to debug LLM prompts?

**A:** Yes:
1. Use `--verbose` to see LLM interactions
2. Override prompts by placing modified files in `internal/llm/prompts/` (from source)
3. See [Prompts Guide](prompts.md) for customization
4. Prompts are embedded in the binary - external files take precedence during development

### Q: Why is the report different each time?

**A:** AI summaries are non-deterministic. For consistent results:
- Use `--no-ai` for raw data only
- AI summaries may vary slightly between runs

## Getting More Help

If your issue isn't covered here:

1. **Check Documentation:**
   - [Installation Guide](installation.md)
   - [Usage Guide](usage.md)
   - [Examples](examples.md)

2. **Search Issues:**
   - [Existing issues](https://github.com/hazadus/gh-repomon/issues)
   - Someone may have had the same problem

3. **Ask for Help:**
   - [Create an issue](https://github.com/hazadus/gh-repomon/issues/new)
   - Provide as much detail as possible
   - Include output from `--verbose` mode

4. **Contributing:**
   - Found a bug? Submit a PR!
   - See [Contributing Guide](contributing.md)

---

[Back to README](../README.md)
