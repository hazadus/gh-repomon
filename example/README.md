# Example Reports

This directory contains sample outputs from gh-repomon to demonstrate what generated reports look like.

## Files

### [sample_report.md](sample_report.md)

A complete example report showing all sections and features:

- Repository overview and statistics
- AI-generated overall summary
- Branch-by-branch analysis with AI summaries
- Pull requests (open and updated) with AI descriptions
- Issues (open and closed)
- Code review statistics
- Author activity breakdown by branch

This example demonstrates:
- ✅ Comprehensive markdown formatting
- ✅ AI-powered summaries at multiple levels
- ✅ Proper GitHub links for all entities
- ✅ Statistics and metrics
- ✅ Multi-author collaboration tracking
- ✅ Professional report structure

## Generating Your Own Reports

To generate similar reports for your repositories:

### Basic Example

```bash
gh-repomon --repo owner/repository --days 7
```

### With Custom Options

```bash
gh-repomon \
  --repo owner/repository \
  --from 2025-09-23 \
  --to 2025-09-30 \
  --exclude-bots \
  --language english \
  > my-report.md
```

### Quick Report (No AI)

```bash
gh-repomon --repo owner/repository --days 7 --no-ai > quick-report.md
```

## Report Sections Explained

### Header
Contains repository name, analysis period, and generation timestamp.

### Summary Statistics
High-level metrics: commits, authors, PRs, issues, reviews.

### Overall Summary (AI)
3-4 paragraph AI-generated overview of all activity during the period.

### Branches
For each active branch:
- AI summary of the branch's purpose
- Statistics (commits, lines changed, contributors)
- Detailed commit list with messages and stats

### Pull Requests
- **Open PRs**: Currently open pull requests
- **Updated PRs**: PRs updated during the period

Each PR includes:
- Title, author, dates
- Comment and review counts
- AI-generated summary of changes

### Issues
- **Open Issues**: Currently open
- **Closed Issues**: Closed during the period

Each issue includes:
- Title, author, dates
- Labels and assignees
- State information

### Code Reviews
List of PRs with reviews and total review count.

### Author Activity
For each contributor:
- Overall statistics (commits, lines, PRs, issues, reviews)
- Breakdown by branch

## Customizing Report Language

Generate reports in different languages:

```bash
# English (default)
gh-repomon --repo owner/repo --days 7 --language english

# Russian
gh-repomon --repo owner/repo --days 7 --language russian

# Spanish
gh-repomon --repo owner/repo --days 7 --language spanish
```

## Converting to Other Formats

### HTML

```bash
gh-repomon --repo owner/repo --days 7 | pandoc -f markdown -t html -o report.html
```

### PDF

```bash
gh-repomon --repo owner/repo --days 7 | pandoc -f markdown -o report.pdf
```

### Word Document

```bash
gh-repomon --repo owner/repo --days 7 | pandoc -f markdown -o report.docx
```

## Viewing in Terminal

Use markdown viewers for better readability:

```bash
# Using glow
gh-repomon --repo owner/repo --days 7 | glow -

# Using mdcat
gh-repomon --repo owner/repo --days 7 | mdcat

# Using bat
gh-repomon --repo owner/repo --days 7 | bat -l md
```

## Tips for Better Reports

1. **Use appropriate time periods**
   - Daily reports: `--days 1`
   - Weekly reports: `--days 7`
   - Sprint reports: `--from YYYY-MM-DD --to YYYY-MM-DD`

2. **Exclude bots** for cleaner team reports
   ```bash
   --exclude-bots
   ```

3. **Filter by user** for individual reports
   ```bash
   --user username
   ```

4. **Choose the right AI model**
   - `openai/gpt-4o` - Best quality
   - `openai/gpt-4o-mini` - Faster, good quality
   - `--no-ai` - No AI summaries (fastest)

## More Examples

For detailed usage examples, see:
- [Usage Guide](../docs/usage.md)
- [Examples Guide](../docs/examples.md)

---

[Back to README](../README.md)
