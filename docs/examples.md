# Examples

Real-world use cases and examples for gh-repomon.

## Table of Contents

- [Daily Standup Reports](#daily-standup-reports)
- [Weekly Team Meetings](#weekly-team-meetings)
- [Sprint Reviews](#sprint-reviews)
- [Individual Performance Reviews](#individual-performance-reviews)
- [Open Source Project Maintenance](#open-source-project-maintenance)
- [Code Review Assistance](#code-review-assistance)
- [Integration Examples](#integration-examples)

## Daily Standup Reports

### Use Case: Engineering Team Daily Standup

Generate a concise report of yesterday's activity for team standup meetings.

**Command:**
```bash
gh-repomon --repo mycompany/api-service --days 1 --exclude-bots
```

**When to use:**
- Daily standup meetings
- Quick overview of recent activity
- Tracking daily progress

**Sample Output Structure:**
```markdown
# Repository Activity Report: api-service

**Period**: 2025-09-29 to 2025-09-30

## Summary Statistics
- Total Commits: 8
- Total Authors: 3
- Open Pull Requests: 2
- Open Issues: 5
- Closed Issues: 1

## 📊 Overall Summary
The team focused on implementing authentication improvements and fixing critical bugs...

## 🌿 Branch: feature/oauth-integration
### AI Summary
Implementation of OAuth 2.0 authentication flow with token refresh mechanism...
...
```

**Tips:**
- Run automatically via cron at 8:30 AM before standup
- Exclude bots to focus on team activity
- Save to a shared location or post to Slack

## Weekly Team Meetings

### Use Case: Sprint Progress Report

Comprehensive weekly report for team meetings and sprint planning.

**Command:**
```bash
gh-repomon --repo mycompany/frontend --days 7 --exclude-bots --language english
```

**When to use:**
- Weekly team meetings
- Sprint retrospectives
- Manager status reports
- Stakeholder updates

**Automation Script:**
```bash
#!/bin/bash
# weekly-report.sh

REPO="mycompany/frontend"
DATE=$(date +%Y-%m-%d)
OUTPUT_DIR="$HOME/reports/weekly"

mkdir -p "$OUTPUT_DIR"

echo "Generating weekly report for $REPO..."
gh-repomon --repo "$REPO" --days 7 --exclude-bots > "$OUTPUT_DIR/report-$DATE.md"

echo "Report saved to: $OUTPUT_DIR/report-$DATE.md"

# Optional: Post to Slack
# slack-cli -d "#team-updates" -f "$OUTPUT_DIR/report-$DATE.md"
```

**Cron Setup (Monday at 9 AM):**
```cron
0 9 * * 1 /home/user/scripts/weekly-report.sh
```

## Sprint Reviews

### Use Case: Two-Week Sprint Report

Generate a report for a specific sprint period.

**Command:**
```bash
gh-repomon \
  --repo mycompany/mobile-app \
  --from 2025-09-16 \
  --to 2025-09-29 \
  --exclude-bots
```

**When to use:**
- Sprint reviews
- Release planning
- Quarterly reports
- Project milestones

**Advanced: Compare Two Sprints**
```bash
#!/bin/bash
# compare-sprints.sh

REPO="mycompany/mobile-app"

echo "Generating Sprint 23 report..."
gh-repomon --repo "$REPO" --from 2025-09-02 --to 2025-09-15 > sprint-23.md

echo "Generating Sprint 24 report..."
gh-repomon --repo "$REPO" --from 2025-09-16 --to 2025-09-29 > sprint-24.md

echo "Reports generated. Compare manually or use your favorite diff tool."
```

## Individual Performance Reviews

### Use Case: Developer Contribution Report

Track individual developer contributions for performance reviews.

**Command:**
```bash
gh-repomon \
  --repo mycompany/backend \
  --days 30 \
  --user alice-developer \
  --exclude-bots
```

**When to use:**
- 1-on-1 meetings
- Performance reviews
- Contribution tracking
- Mentorship sessions

**Quarterly Review Script:**
```bash
#!/bin/bash
# quarterly-review.sh

REPO="mycompany/backend"
DEVELOPER="alice-developer"
QUARTER_START="2025-07-01"
QUARTER_END="2025-09-30"
OUTPUT="reviews/alice-q3-2025.md"

gh-repomon \
  --repo "$REPO" \
  --from "$QUARTER_START" \
  --to "$QUARTER_END" \
  --user "$DEVELOPER" \
  > "$OUTPUT"

echo "Quarterly review generated: $OUTPUT"
```

## Open Source Project Maintenance

### Use Case: Monthly Contributors Report

Generate reports for open source project maintainers.

**Command:**
```bash
gh-repomon \
  --repo facebook/react \
  --days 30 \
  --exclude-bots \
  --model openai/gpt-4o
```

**When to use:**
- Community updates
- Contributor recognition
- Project newsletters
- Sponsor reports

**Monthly Newsletter Script:**
```bash
#!/bin/bash
# monthly-newsletter.sh

REPO="yourorg/your-project"
MONTH=$(date +%B-%Y)
OUTPUT="newsletter-$MONTH.md"

# Generate base report
gh-repomon --repo "$REPO" --days 30 --exclude-bots > "$OUTPUT"

# Add newsletter header
cat > "newsletter-header.md" << EOF
# 📰 Project Newsletter - $MONTH

Welcome to this month's update for our project!

---

EOF

cat newsletter-header.md "$OUTPUT" > "final-newsletter-$MONTH.md"
rm newsletter-header.md "$OUTPUT"

echo "Newsletter ready: final-newsletter-$MONTH.md"
```

## Code Review Assistance

### Use Case: PR Context for Reviewers

Generate activity report to provide context for code reviews.

**Command:**
```bash
# Get activity on a specific branch
gh-repomon \
  --repo mycompany/api \
  --days 14 \
  --exclude-bots \
  --no-ai > branch-context.md
```

**When to use:**
- Before reviewing large PRs
- Understanding branch history
- Catching up after vacation
- Onboarding new reviewers

**PR Review Helper Script:**
```bash
#!/bin/bash
# pr-review-helper.sh

if [ $# -eq 0 ]; then
    echo "Usage: $0 <pr-number>"
    exit 1
fi

PR_NUMBER=$1
REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner)

echo "Fetching PR details for #$PR_NUMBER..."
PR_BRANCH=$(gh pr view $PR_NUMBER --json headRefName -q .headRefName)
PR_CREATED=$(gh pr view $PR_NUMBER --json createdAt -q .createdAt | cut -d'T' -f1)

echo "Generating activity report for branch $PR_BRANCH..."
gh-repomon \
  --repo "$REPO" \
  --from "$PR_CREATED" \
  --to $(date +%Y-%m-%d) \
  --no-ai \
  > "pr-$PR_NUMBER-context.md"

echo "Context saved to: pr-$PR_NUMBER-context.md"
```

## Integration Examples

### Slack Integration

Post daily reports to Slack channel.

**Script:**
```bash
#!/bin/bash
# slack-daily-report.sh

REPO="mycompany/backend"
SLACK_WEBHOOK="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"

# Generate report
REPORT=$(gh-repomon --repo "$REPO" --days 1 --exclude-bots 2>/dev/null)

# Post to Slack
curl -X POST "$SLACK_WEBHOOK" \
  -H 'Content-Type: application/json' \
  -d @- << EOF
{
  "text": "Daily Activity Report",
  "blocks": [
    {
      "type": "section",
      "text": {
        "type": "mrkdwn",
        "text": "📊 *Daily Activity Report for $REPO*\n\`\`\`\n${REPORT:0:2000}\n\`\`\`"
      }
    }
  ]
}
EOF
```

### Email Reports

Send reports via email using mail command.

**Script:**
```bash
#!/bin/bash
# email-weekly-report.sh

REPO="mycompany/frontend"
EMAIL="team@mycompany.com"
SUBJECT="Weekly Activity Report - $(date +%Y-%m-%d)"

# Generate report
gh-repomon --repo "$REPO" --days 7 --exclude-bots > /tmp/report.md

# Convert to HTML (requires pandoc)
pandoc /tmp/report.md -o /tmp/report.html

# Send email
mail -s "$SUBJECT" -a "Content-Type: text/html" "$EMAIL" < /tmp/report.html

# Cleanup
rm /tmp/report.md /tmp/report.html
```

### GitHub Actions Workflow

Automated weekly reports via GitHub Actions.

**`.github/workflows/weekly-report.yml`:**
```yaml
name: Weekly Activity Report

on:
  schedule:
    # Run every Monday at 9 AM UTC
    - cron: '0 9 * * 1'
  workflow_dispatch:  # Allow manual trigger

jobs:
  generate-report:
    runs-on: ubuntu-latest
    permissions:
      contents: write
      pull-requests: read
      issues: read

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install gh-repomon
        run: |
          go install github.com/hazadus/gh-repomon/cmd/repomon@latest

      - name: Generate Report
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          repomon --repo ${{ github.repository }} --days 7 --exclude-bots > weekly-report.md

      - name: Upload Report
        uses: actions/upload-artifact@v3
        with:
          name: weekly-report
          path: weekly-report.md
          retention-days: 90

      - name: Create Issue with Report
        uses: peter-evans/create-issue-from-file@v4
        with:
          title: Weekly Activity Report - ${{ env.REPORT_DATE }}
          content-filepath: weekly-report.md
          labels: report, automated
```

### Confluence/Wiki Integration

Upload reports to Confluence automatically.

**Script (requires `confluence-cli`):**
```bash
#!/bin/bash
# upload-to-confluence.sh

REPO="mycompany/backend"
CONFLUENCE_SPACE="TEAM"
CONFLUENCE_PARENT="Weekly Reports"

# Generate report
gh-repomon --repo "$REPO" --days 7 > /tmp/weekly-report.md

# Upload to Confluence
confluence \
  --action addPage \
  --space "$CONFLUENCE_SPACE" \
  --parent "$CONFLUENCE_PARENT" \
  --title "Activity Report $(date +%Y-%m-%d)" \
  --file /tmp/weekly-report.md

rm /tmp/weekly-report.md
```

### Dashboard Integration

Create a simple dashboard that displays recent reports.

**`generate-dashboard.sh`:**
```bash
#!/bin/bash
# generate-dashboard.sh

REPOS=("mycompany/backend" "mycompany/frontend" "mycompany/mobile")
OUTPUT_DIR="/var/www/reports"

mkdir -p "$OUTPUT_DIR"

# Generate HTML dashboard header
cat > "$OUTPUT_DIR/index.html" << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>Team Activity Dashboard</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/github-markdown-css@5/github-markdown.min.css">
</head>
<body class="markdown-body">
    <h1>Team Activity Dashboard</h1>
    <p>Generated: <strong>$(date)</strong></p>
EOF

# Generate reports for each repo
for repo in "${REPOS[@]}"; do
    repo_name=$(echo $repo | tr '/' '-')
    echo "<h2>$repo</h2>" >> "$OUTPUT_DIR/index.html"

    gh-repomon --repo "$repo" --days 7 --exclude-bots --no-ai | \
        pandoc -f markdown -t html >> "$OUTPUT_DIR/index.html"

    echo "<hr>" >> "$OUTPUT_DIR/index.html"
done

# Close HTML
echo "</body></html>" >> "$OUTPUT_DIR/index.html"

echo "Dashboard generated at: $OUTPUT_DIR/index.html"
```

## Tips for Different Scenarios

### For Managers

**Weekly Team Overview:**
```bash
gh-repomon --repo team/project --days 7 --exclude-bots
```

**Individual Check-ins:**
```bash
for dev in alice bob charlie; do
    gh-repomon --repo team/project --days 7 --user $dev > "reports/$dev-weekly.md"
done
```

### For Solo Developers

**Personal Activity Log:**
```bash
gh-repomon --repo myusername/project --days 7 --user myusername --no-ai
```

**Monthly Summary:**
```bash
gh-repomon --repo myusername/project --days 30 --user myusername
```

### For Open Source Maintainers

**Release Notes Helper:**
```bash
gh-repomon --repo org/project --from 2025-09-01 --to 2025-09-30 --exclude-bots
```

**Contributor Recognition:**
```bash
gh-repomon --repo org/project --days 30 | grep "### \[" | sort -u
```

### For DevOps/Release Engineers

**Pre-Release Activity Check:**
```bash
gh-repomon --repo company/product --from 2025-09-15 --to 2025-09-29 --exclude-bots --no-ai
```

**Change Log Generation:**
```bash
gh-repomon --repo company/product --from v1.0.0 --to v1.1.0 > CHANGELOG-1.1.0.md
```

## More Examples

For more examples and use cases, check out:

- [Usage Guide](usage.md) - Detailed flag documentation
- [Troubleshooting](troubleshooting.md) - Common issues and solutions
- [Contributing](contributing.md) - How to contribute examples

---

[Back to README](../README.md)
