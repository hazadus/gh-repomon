# YAML Prompts Guide

Complete guide to understanding and customizing AI prompts in gh-repomon.

## Table of Contents

- [Overview](#overview)
- [Prompt Structure](#prompt-structure)
- [Available Prompts](#available-prompts)
- [Template Variables](#template-variables)
- [Creating Custom Prompts](#creating-custom-prompts)
- [Best Practices](#best-practices)
- [Examples](#examples)

## Overview

gh-repomon uses YAML-based prompts to generate AI summaries. This approach offers several advantages:

- **Easy to Edit** - Simple text files, no code changes needed
- **Version Control** - Track prompt changes over time
- **Customizable** - Adapt prompts to your needs
- **Maintainable** - Non-developers can modify prompts
- **Transparent** - See exactly what's being sent to the AI

All prompts are located in `internal/llm/prompts/` with the `.prompt.yml` extension.

### Embedded vs External Prompts

**Embedded Prompts (Production):**
- All prompt files are embedded directly into the compiled binary using Go's `//go:embed` directive
- The binary is fully self-contained and works anywhere without requiring external files
- Perfect for CI/CD environments, `gh extension install`, and distribution
- No need to copy or deploy prompt files separately

**External Prompts (Development/Customization):**
- You can override embedded prompts by placing files in `internal/llm/prompts/`
- External files take precedence over embedded ones
- Useful for development, testing, and customization
- Changes take effect immediately without recompilation

**Load Priority:**
1. First: External file in `internal/llm/prompts/` (if exists)
2. Fallback: Embedded file in the binary

This design ensures the tool works out-of-the-box while allowing easy customization when needed.

## Prompt Structure

### Basic YAML Format

```yaml
name: Prompt Name
description: Brief description of what this prompt does
model: openai/gpt-4o
modelParameters:
  temperature: 0.7
  topP: 0.9
messages:
  - role: system
    content: >
      System instructions go here.
      They define the AI's behavior and role.

  - role: user
    content: |
      User message goes here.
      This contains the actual task and data.

      Variables are inserted like: {{variable_name}}
```

### Fields Explained

#### `name` (string)
Human-readable name for the prompt. Used for identification and logging.

```yaml
name: Overall Repository Summary
```

#### `description` (string)
Brief description of the prompt's purpose. Helps developers understand when to use it.

```yaml
description: Generates overall summary of repository activity
```

#### `model` (string)
Default AI model to use for this prompt. Can be overridden by CLI flag.

```yaml
model: openai/gpt-4o
```

**Available Models:**
- `openai/gpt-4o` - Most capable (recommended)
- `openai/gpt-4o-mini` - Faster, cheaper
- `meta/llama-3.1-405b-instruct` - Open source
- `anthropic/claude-3-5-sonnet` - Alternative

#### `modelParameters` (object)
Model-specific parameters for fine-tuning output.

```yaml
modelParameters:
  temperature: 0.7    # Creativity (0.0 = deterministic, 1.0 = creative)
  topP: 0.9          # Nucleus sampling parameter
```

**Parameter Guide:**
- **temperature**:
  - `0.0-0.3`: Focused, deterministic
  - `0.4-0.7`: Balanced (recommended)
  - `0.8-1.0`: Creative, varied
- **topP**: Usually keep at `0.9` for good results

#### `messages` (array)
Array of message objects that form the conversation.

**Message Types:**

1. **System Message** - Sets AI behavior
   ```yaml
   - role: system
     content: You are a helpful assistant analyzing code.
   ```

2. **User Message** - Contains the task and data
   ```yaml
   - role: user
     content: Analyze this commit: {{commit_message}}
   ```

### Content Formatting

**Multi-line with `>`** - Folds newlines into spaces:
```yaml
content: >
  This will be
  a single line
  of text.
```
Result: "This will be a single line of text."

**Multi-line with `|`** - Preserves newlines:
```yaml
content: |
  Line 1
  Line 2
  Line 3
```
Result: Three separate lines.

## Available Prompts

### 1. Overall Summary (`overall_summary.prompt.yml`)

**Purpose:** Generates comprehensive summary of entire repository activity

**Location:** `internal/llm/prompts/overall_summary.prompt.yml`

**Variables:**
- `{{language}}` - Output language
- `{{repo_name}}` - Repository name
- `{{period}}` - Time period
- `{{total_commits}}` - Number of commits
- `{{total_authors}}` - Number of authors
- `{{branches}}` - List of branches
- `{{prs}}` - Pull requests
- `{{issues}}` - Issues

**Output:** 3-4 paragraph summary of overall activity

### 2. Branch Summary (`branch_summary.prompt.yml`)

**Purpose:** Generates summary for a specific branch's activity

**Location:** `internal/llm/prompts/branch_summary.prompt.yml`

**Variables:**
- `{{language}}` - Output language
- `{{branch_name}}` - Name of the branch
- `{{commit_count}}` - Number of commits
- `{{commit_messages}}` - List of commit messages
- `{{authors}}` - Contributing authors

**Output:** 2-3 sentences describing branch functionality

### 3. PR Summary (`pr_summary.prompt.yml`)

**Purpose:** Generates summary for a pull request

**Location:** `internal/llm/prompts/pr_summary.prompt.yml`

**Variables:**
- `{{language}}` - Output language
- `{{pr_title}}` - PR title
- `{{pr_description}}` - PR description
- `{{commit_messages}}` - Commits in the PR

**Output:** Brief description of PR purpose and changes

## Template Variables

### Variable Syntax

Variables use double curly braces:
```yaml
content: |
  Repository: {{repo_name}}
  Period: {{period}}
```

### Common Variables

| Variable | Type | Description | Example |
|----------|------|-------------|---------|
| `{{language}}` | string | Output language | `"english"`, `"russian"` |
| `{{repo_name}}` | string | Repository name | `"owner/repo"` |
| `{{period}}` | string | Date range | `"2025-09-01 to 2025-09-30"` |
| `{{total_commits}}` | string | Commit count | `"42"` |
| `{{total_authors}}` | string | Author count | `"5"` |
| `{{branch_name}}` | string | Branch name | `"feature/auth"` |
| `{{commit_count}}` | string | Number of commits | `"15"` |
| `{{commit_messages}}` | string | Formatted list | `"- feat: add login\n- fix: bug"` |
| `{{authors}}` | string | Author list | `"alice, bob, charlie"` |

### Variable Processing

Variables are replaced during prompt rendering:

```go
vars := map[string]string{
    "repo_name": "mycompany/backend",
    "period": "2025-09-01 to 2025-09-30",
    "total_commits": "156",
}

// Template: "Repository: {{repo_name}}"
// Result: "Repository: mycompany/backend"
```

## Creating Custom Prompts

### Step 1: Create YAML File

Create a new file in `internal/llm/prompts/`:

```bash
touch internal/llm/prompts/my_custom.prompt.yml
```

**Note:** When developing from source, external files in `internal/llm/prompts/` automatically override embedded prompts. To use custom prompts in production, you need to rebuild the binary after adding the file (the new prompt will be embedded).

### Step 2: Define Prompt Structure

```yaml
name: Custom Summary
description: My custom prompt for specific use case
model: openai/gpt-4o
modelParameters:
  temperature: 0.7
  topP: 0.9
messages:
  - role: system
    content: >
      You are an expert software analyst.
      Your task is to analyze code changes and provide insights.

      Output language: {{language}}

  - role: user
    content: |
      Analyze the following data:

      Repository: {{repo_name}}
      Changes: {{changes_description}}

      Provide a detailed analysis.
```

### Step 3: Define Variables

List variables needed:
```go
vars := map[string]string{
    "language": language,
    "repo_name": data.Repository,
    "changes_description": formatChanges(data),
}
```

### Step 4: Use in Code

```go
// Load and render prompt
config, err := LoadPrompt("my_custom")
if err != nil {
    return "", err
}

rendered, err := RenderPrompt(config, vars)
if err != nil {
    return "", err
}

// Generate summary
summary, err := client.Complete(createRequest(rendered))
```

## Best Practices

### 1. Clear Instructions

Be specific about what you want:

**Bad:**
```yaml
content: Summarize this.
```

**Good:**
```yaml
content: >
  Analyze the repository activity and create a 3-4 paragraph summary.
  Include: main features developed, bug fixes, and team collaboration patterns.
  Focus on the business value of changes.
```

### 2. Provide Context

Give the AI necessary background:

```yaml
- role: system
  content: >
    You are analyzing a software development repository.
    Users are technical team members who need concise, actionable summaries.
    Focus on what changed and why it matters.
```

### 3. Structured Output

Request specific format:

```yaml
content: |
  Create a summary with these sections:
  1. Overview (2-3 sentences)
  2. Key Changes (bullet points)
  3. Impact (1-2 sentences)
```

### 4. Examples in Prompts

Show desired output format:

```yaml
content: |
  Format your response like this example:

  "The team focused on authentication improvements, implementing OAuth2
  integration and fixing several security issues. This enhances user
  security and enables single sign-on capabilities."

  Now analyze: {{data}}
```

### 5. Temperature Settings

Choose based on use case:

```yaml
# For factual summaries (use low temperature)
modelParameters:
  temperature: 0.3

# For creative descriptions (use higher temperature)
modelParameters:
  temperature: 0.7
```

### 6. Language Handling

Always include language variable:

```yaml
- role: system
  content: >
    Generate your response in {{language}}.
    If {{language}} is "english", use clear technical English.
    If {{language}} is "russian", use Russian with technical terms in English.
```

### 7. Error Handling

Prepare for missing variables:

```yaml
# Good: Provide defaults or handle gracefully
Repository: {{repo_name|default:"unknown"}}

# Better: Use code to validate before rendering
if repo_name == "" {
    return errors.New("repo_name is required")
}
```

### 8. Testing Prompts

Test with various inputs:

```bash
# Test with real data
gh-repomon --repo small/repo --days 1

# Test with large data
gh-repomon --repo large/repo --days 30

# Test different languages
gh-repomon --repo test/repo --days 7 --language russian
```

## Examples

### Example 1: Concise Technical Summary

```yaml
name: Concise Technical Summary
description: Short, technical summary for developers
model: openai/gpt-4o-mini
modelParameters:
  temperature: 0.3
messages:
  - role: system
    content: >
      You are a senior software engineer reviewing code changes.
      Provide concise, technical summaries focused on implementation details.
      Use technical terminology. Be direct and factual.

      Output language: {{language}}

  - role: user
    content: |
      Branch: {{branch_name}}
      Commits: {{commit_count}}

      Changes:
      {{commit_messages}}

      Provide a 1-2 sentence technical summary of changes.
```

### Example 2: Executive Summary

```yaml
name: Executive Summary
description: High-level summary for non-technical stakeholders
model: openai/gpt-4o
modelParameters:
  temperature: 0.5
messages:
  - role: system
    content: >
      You are a technical project manager explaining changes to executives.
      Use clear, non-technical language. Focus on business value and impact.
      Avoid jargon and implementation details.

      Output language: {{language}}

  - role: user
    content: |
      Repository: {{repo_name}}
      Period: {{period}}
      Team Size: {{total_authors}} developers

      Activity Summary:
      - Commits: {{total_commits}}
      - Pull Requests: {{pr_count}}
      - Issues Resolved: {{issues_closed}}

      Create a 2-3 paragraph executive summary explaining:
      1. What the team accomplished
      2. Business value delivered
      3. Any notable challenges or achievements
```

### Example 3: Release Notes Generator

```yaml
name: Release Notes
description: Generate release notes from changes
model: openai/gpt-4o
modelParameters:
  temperature: 0.4
messages:
  - role: system
    content: >
      You are creating release notes for a software product.
      Format: Markdown with sections for Features, Fixes, and Breaking Changes.
      Be clear and user-focused.

      Output language: {{language}}

  - role: user
    content: |
      Version: {{version}}
      Date: {{release_date}}

      Changes:
      {{changes}}

      Generate release notes with these sections:

      ## 🎉 New Features
      ## 🐛 Bug Fixes
      ## ⚠️ Breaking Changes
      ## 📝 Other Changes

      Use bullet points. Include commit references where relevant.
```

## Troubleshooting

### Variables Not Replaced

**Problem:** Variables appear as `{{variable}}` in output

**Solution:**
- Check variable name spelling
- Ensure variable is provided in `vars` map
- Verify curly braces are doubled: `{{var}}` not `{var}`

### Unexpected AI Output

**Problem:** AI generates wrong format or content

**Solution:**
- Lower temperature for more consistent output
- Add examples to the prompt
- Be more specific in instructions
- Test with different models

### Language Not Respected

**Problem:** AI responds in wrong language

**Solution:**
- Ensure `{{language}}` variable is in system message
- Be explicit: "You MUST respond in {{language}}"
- Test that language variable is passed correctly

## Advanced Topics

### Conditional Content

For advanced use cases, handle conditionals in code:

```go
promptContent := basePrompt
if includeStats {
    promptContent += "\n\nStatistics: {{stats}}"
}
```

### Prompt Chaining

Use multiple prompts in sequence:

```go
// First: Generate summary
summary := generateSummary(data)

// Second: Refine summary
refined := refineSummary(summary)

// Third: Translate
translated := translate(refined, language)
```

### A/B Testing Prompts

Test different prompts:

```go
prompts := []string{"prompt_v1", "prompt_v2"}
for _, promptName := range prompts {
    result := testPrompt(promptName, testData)
    evaluateQuality(result)
}
```

## Further Reading

- [OpenAI Prompt Engineering Guide](https://platform.openai.com/docs/guides/prompt-engineering)
- [Architecture Documentation](architecture.md)
- [Contributing Guide](contributing.md)

---

[Back to README](../README.md)
