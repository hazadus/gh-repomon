# Usage Guide

Complete guide to using gh-repomon with detailed flag descriptions and examples.

## Table of Contents

- [Basic Usage](#basic-usage)
- [Command-Line Flags](#command-line-flags)
- [Common Scenarios](#common-scenarios)
- [Output Management](#output-management)
- [Advanced Usage](#advanced-usage)

## Basic Usage

The basic syntax for gh-repomon:

```bash
gh-repomon --repo OWNER/REPOSITORY [flags]
```

Or if installed as a GitHub CLI extension:

```bash
gh repomon --repo OWNER/REPOSITORY [flags]
```

## Command-Line Flags

### Required Flags

#### `--repo`, `-r` (string)

The GitHub repository to analyze in the format `owner/repository`.

```bash
gh-repomon --repo microsoft/vscode --days 7
```

### Time Period Flags

You can specify the time period in two ways: relative (days) or absolute (from/to dates).

#### `--days`, `-d` (int, default: 1)

Number of days to look back from today.

```bash
# Report for the last 7 days
gh-repomon --repo owner/repo --days 7

# Report for yesterday only
gh-repomon --repo owner/repo --days 1

# Report for the last 30 days
gh-repomon --repo owner/repo --days 30
```

#### `--from` (string, format: YYYY-MM-DD)

Start date of the reporting period.

```bash
gh-repomon --repo owner/repo --from 2025-09-01 --to 2025-09-30
```

#### `--to` (string, format: YYYY-MM-DD)

End date of the reporting period.

```bash
gh-repomon --repo owner/repo --from 2025-09-01 --to 2025-09-30
```

**Note:** When both `--from`/`--to` and `--days` are specified, `--from`/`--to` takes precedence.

### Filtering Flags

#### `--user`, `-u` (string)

Filter activity by a specific GitHub username.

```bash
# Only show activity from user "hazadus"
gh-repomon --repo owner/repo --days 7 --user hazadus
```

This filters:
- Commits by the user
- Pull requests created by the user
- Issues created by the user
- Code reviews by the user

#### `--exclude-bots` (boolean, default: false)

Exclude activity from bot accounts.

```bash
gh-repomon --repo owner/repo --days 7 --exclude-bots
```

Automatically excludes:
- Accounts ending with `[bot]`
- `github-actions`
- `dependabot`
- `renovate`
- Other known bot accounts

### AI Configuration Flags

#### `--model`, `-m` (string, default: "openai/gpt-4o")

Specify the AI model to use for generating summaries.

```bash
# Use GPT-4o (default)
gh-repomon --repo owner/repo --days 7 --model openai/gpt-4o

# Use GPT-4o-mini for faster, cheaper summaries
gh-repomon --repo owner/repo --days 7 --model openai/gpt-4o-mini

# Use other available models
gh-repomon --repo owner/repo --days 7 --model meta/llama-3.1-405b-instruct
```

Available models (as of October 2025):
- `openai/gpt-4o` - Most capable, recommended
- `openai/gpt-4o-mini` - Faster and cheaper
- `meta/llama-3.1-405b-instruct` - Open source alternative
- `anthropic/claude-3-5-sonnet` - High-quality alternative

Check [GitHub Models Marketplace](https://github.com/marketplace/models) for the latest available models.

#### `--language`, `-l` (string, default: "english")

Language for AI-generated summaries.

```bash
# English (default)
gh-repomon --repo owner/repo --days 7 --language english

# Russian
gh-repomon --repo owner/repo --days 7 --language russian

# Spanish
gh-repomon --repo owner/repo --days 7 --language spanish

# Any language
gh-repomon --repo owner/repo --days 7 --language french
```

### Performance Flags

#### `--no-ai` (boolean, default: false)

Disable AI summary generation for faster report generation.

```bash
gh-repomon --repo owner/repo --days 7 --no-ai
```

Use this when:
- You only need raw data without summaries
- You want faster generation
- GitHub Models API is unavailable
- You're generating many reports in batch

#### `--verbose`, `-v` (boolean, default: false)

Enable detailed logging for debugging.

```bash
gh-repomon --repo owner/repo --days 7 --verbose
```

Shows:
- Detailed API calls
- Token usage statistics
- Processing timings
- Debug information

## Common Scenarios

### Daily Standup Report

Generate a report for yesterday's activity:

```bash
gh-repomon --repo mycompany/backend --days 1
```

### Weekly Team Report

Generate a comprehensive weekly report:

```bash
gh-repomon --repo mycompany/backend --days 7
```

### Sprint Report

Generate a report for a specific sprint period:

```bash
gh-repomon --repo mycompany/backend --from 2025-09-15 --to 2025-09-29
```

### Individual Developer Report

Track a specific developer's contributions:

```bash
gh-repomon --repo mycompany/backend --days 7 --user john-doe
```

### Clean Report (No Bots)

Generate a report excluding bot activity:

```bash
gh-repomon --repo mycompany/backend --days 7 --exclude-bots
```

### Quick Report (No AI)

Generate a fast report without AI summaries:

```bash
gh-repomon --repo mycompany/backend --days 7 --no-ai
```

### Multilingual Report

Generate a report in Russian for the team:

```bash
gh-repomon --repo mycompany/backend --days 7 --language russian
```

### Cost-Optimized Report

Use a cheaper model for regular reports:

```bash
gh-repomon --repo mycompany/backend --days 7 --model openai/gpt-4o-mini
```

## Output Management

### Saving to File

Save the report to a markdown file:

```bash
gh-repomon --repo owner/repo --days 7 > weekly-report.md
```

### Viewing in Terminal

Use a markdown viewer for better readability:

```bash
# Using glow
gh-repomon --repo owner/repo --days 7 | glow -

# Using mdcat
gh-repomon --repo owner/repo --days 7 | mdcat

# Using bat
gh-repomon --repo owner/repo --days 7 | bat -l md
```

### Separating Progress Messages

Progress messages go to stderr, report goes to stdout:

```bash
# Save only the report, show progress
gh-repomon --repo owner/repo --days 7 > report.md

# Save both report and progress
gh-repomon --repo owner/repo --days 7 > report.md 2> progress.log

# Suppress progress messages
gh-repomon --repo owner/repo --days 7 2>/dev/null > report.md
```

### Converting to Other Formats

Convert markdown to other formats using pandoc:

```bash
# Convert to HTML
gh-repomon --repo owner/repo --days 7 | pandoc -f markdown -t html -o report.html

# Convert to PDF
gh-repomon --repo owner/repo --days 7 | pandoc -f markdown -o report.pdf

# Convert to Word
gh-repomon --repo owner/repo --days 7 | pandoc -f markdown -o report.docx
```

## Advanced Usage

### Automated Daily Reports

Set up a cron job for daily reports:

```bash
# Add to crontab (run daily at 9 AM)
0 9 * * * /usr/local/bin/gh-repomon --repo mycompany/backend --days 1 > ~/reports/daily-$(date +\%Y-\%m-\%d).md 2>&1
```

### Weekly Reports via GitHub Actions

Create `.github/workflows/weekly-report.yml`:

```yaml
name: Weekly Report

on:
  schedule:
    - cron: '0 9 * * 1'  # Every Monday at 9 AM
  workflow_dispatch:

jobs:
  generate-report:
    runs-on: ubuntu-latest
    steps:
      - name: Generate Report
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          gh extension install hazadus/gh-repomon
          gh repomon --repo ${{ github.repository }} --days 7 > weekly-report.md

      - name: Upload Report
        uses: actions/upload-artifact@v3
        with:
          name: weekly-report
          path: weekly-report.md
```

### Batch Reports for Multiple Repositories

Generate reports for multiple repositories:

```bash
#!/bin/bash
REPOS=(
  "mycompany/backend"
  "mycompany/frontend"
  "mycompany/mobile"
)

for repo in "${REPOS[@]}"; do
  echo "Generating report for $repo..."
  gh-repomon --repo "$repo" --days 7 > "reports/$(echo $repo | tr '/' '-')-report.md"
done
```

### Compare Periods

Generate and compare reports for different periods:

```bash
# Current week
gh-repomon --repo owner/repo --days 7 > this-week.md

# Previous week
gh-repomon --repo owner/repo --from 2025-09-08 --to 2025-09-14 > last-week.md

# Compare manually or use diff tools
```

### Custom AI Model with Environment Variable

If you need to use a custom endpoint:

```bash
export GITHUB_MODELS_ENDPOINT="https://custom-endpoint.com"
gh-repomon --repo owner/repo --days 7
```

### Debugging Issues

Enable verbose mode and capture logs:

```bash
gh-repomon --repo owner/repo --days 7 --verbose 2> debug.log
```

### Performance Tuning

For large repositories, optimize performance:

```bash
# Disable AI for speed
gh-repomon --repo large/repo --days 7 --no-ai

# Use faster model
gh-repomon --repo large/repo --days 7 --model openai/gpt-4o-mini

# Shorter time period
gh-repomon --repo large/repo --days 1
```

## Tips and Best Practices

### 1. Start Small

When analyzing a new repository, start with a small time period:

```bash
gh-repomon --repo new/repo --days 1
```

### 2. Exclude Bots for Cleaner Reports

For team reports, always exclude bots:

```bash
gh-repomon --repo team/repo --days 7 --exclude-bots
```

### 3. Use Appropriate Models

- **gpt-4o**: Best quality, use for important reports
- **gpt-4o-mini**: Good balance, use for regular reports
- **--no-ai**: Fastest, use for quick checks

### 4. Filter by User for 1-on-1s

When preparing for 1-on-1 meetings:

```bash
gh-repomon --repo team/repo --days 14 --user developer-name
```

### 5. Automate Regular Reports

Set up automation for consistency:
- Daily reports for active development
- Weekly reports for team meetings
- Monthly reports for stakeholders

### 6. Combine with Other Tools

Integrate with your workflow:
- Slack notifications
- Email reports
- Dashboard integrations
- Wiki updates

## Getting Help

For more information:

```bash
gh-repomon --help
```

For troubleshooting, see [Troubleshooting Guide](troubleshooting.md).

For examples, see [Examples Guide](examples.md).

---

[Back to README](../README.md)
