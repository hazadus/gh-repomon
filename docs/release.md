# Release Process

This document describes the release process for gh-repomon maintainers.

## Overview

gh-repomon uses [Semantic Versioning](https://semver.org/) for version numbers:
- **MAJOR** version when you make incompatible API changes
- **MINOR** version when you add functionality in a backwards compatible manner
- **PATCH** version when you make backwards compatible bug fixes

## Prerequisites

Before creating a release, ensure:

1. All tests pass:
   ```bash
   just test-all
   ```

2. Code is properly formatted:
   ```bash
   just format
   ```

3. Linter passes:
   ```bash
   just lint
   ```

4. Documentation is up to date

5. CHANGELOG.md is updated with changes for the new version

## Release Steps

### 1. Update CHANGELOG.md

Edit `CHANGELOG.md` to:
- Move changes from `[Unreleased]` to a new version section
- Add the release date
- Update the comparison links at the bottom

Example:
```markdown
## [1.1.0] - 2025-02-01

### Added
- New feature X
- New feature Y

### Fixed
- Bug Z

[Unreleased]: https://github.com/hazadus/gh-repomon/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/hazadus/gh-repomon/compare/v1.0.0...v1.1.0
```

### 2. Update Version in Code

Update the version in `extension.yml`:
```yaml
version: 1.1.0
```

### 3. Commit Changes

Commit the version bump:
```bash
git add CHANGELOG.md extension.yml
git commit -m "chore: prepare release v1.1.0"
git push origin main
```

### 4. Create and Push Tag

Create an annotated tag:
```bash
git tag -a v1.1.0 -m "Release version 1.1.0"
git push origin v1.1.0
```

### 5. Automated Release

Once the tag is pushed, GitHub Actions will automatically:
- Run all tests
- Build binaries for all platforms:
  - Linux (AMD64, ARM64)
  - macOS (Intel, Apple Silicon)
  - Windows (AMD64, ARM64)
- Generate checksums
- Create a GitHub Release with:
  - Release notes extracted from CHANGELOG.md
  - All platform binaries
  - Checksums file

### 6. Verify Release

After the GitHub Actions workflow completes:

1. Go to the [Releases page](https://github.com/hazadus/gh-repomon/releases)
2. Verify the new release is published
3. Check that all binaries are attached
4. Review the release notes

### 7. Test Installation

Test the release by installing it as a GitHub CLI extension:

```bash
# Remove existing installation if any
gh extension remove repomon

# Install the new release
gh extension install hazadus/gh-repomon

# Verify it works
gh repomon --version
gh repomon --help
```

### 8. Announce Release

Consider announcing the release:
- GitHub Discussions (if enabled)
- Twitter/Social media
- Relevant community forums

## Local Release Testing

Before creating an actual release, you can test the build process locally:

### Using Just

```bash
# Test complete release build with checksums
just release

# Test goreleaser configuration (if goreleaser is installed)
just release-test
```

### Manual Testing

```bash
# Build for all platforms
just build-all

# Test a specific binary
./bin/gh-repomon-linux-amd64 --help
```

## Using GoReleaser (Alternative)

If you prefer to use GoReleaser for releases:

### Prerequisites

Install goreleaser:
```bash
brew install goreleaser
```

### Test Release

Test the release process without publishing:
```bash
goreleaser release --snapshot --clean
```

### Create Release

With a tag checked out:
```bash
goreleaser release --clean
```

Set `GITHUB_TOKEN` environment variable if needed:
```bash
export GITHUB_TOKEN="your_github_token"
goreleaser release --clean
```

## Troubleshooting

### Release Build Fails

If the GitHub Actions release workflow fails:

1. Check the [Actions tab](https://github.com/hazadus/gh-repomon/actions) for error details
2. Fix the issue in a new commit
3. Delete the failed tag:
   ```bash
   git tag -d v1.1.0
   git push origin :refs/tags/v1.1.0
   ```
4. Restart the release process from step 3

### Binaries Not Building

If a specific platform binary fails to build:

1. Test the build locally:
   ```bash
   GOOS=linux GOARCH=amd64 go build -o test-binary ./cmd/repomon
   ```
2. Fix any platform-specific issues
3. Update the code and create a new release

### Release Notes Missing

If release notes don't appear in the GitHub Release:

1. Verify CHANGELOG.md has the correct format
2. Check that the version section exists: `## [1.1.0]`
3. The workflow extracts notes between version headers
4. You can manually edit the release notes on GitHub after creation

## Hotfix Releases

For urgent bug fixes:

1. Create a hotfix branch from the release tag:
   ```bash
   git checkout -b hotfix/1.0.1 v1.0.0
   ```

2. Make the fix and commit:
   ```bash
   git commit -am "fix: critical bug in X"
   ```

3. Update CHANGELOG.md with the hotfix

4. Merge to main:
   ```bash
   git checkout main
   git merge hotfix/1.0.1
   ```

5. Follow the normal release process for version 1.0.1

## Version Numbering Guidelines

### When to Bump MAJOR Version (X.0.0)

- Breaking changes to command-line interface
- Incompatible changes to output format
- Removal of features
- Major architectural changes

### When to Bump MINOR Version (0.X.0)

- New features added
- New command-line flags
- New output formats (if backwards compatible)
- Significant performance improvements

### When to Bump PATCH Version (0.0.X)

- Bug fixes
- Documentation updates
- Performance improvements (minor)
- Dependency updates (security fixes)

## Post-Release Tasks

After a successful release:

1. Update the [Unreleased] section in CHANGELOG.md:
   ```markdown
   ## [Unreleased]

   ### Added
   ### Changed
   ### Deprecated
   ### Removed
   ### Fixed
   ### Security
   ```

2. Close any milestone associated with the release on GitHub

3. Review and close related issues

4. Update project documentation if needed

## Emergency Rollback

If a critical issue is discovered after release:

1. Create a new hotfix release with the fix (preferred)

2. OR, if necessary, delete the release:
   ```bash
   # Delete the tag locally and remotely
   git tag -d v1.1.0
   git push origin :refs/tags/v1.1.0

   # Delete the GitHub Release via the web interface
   ```

3. Notify users about the issue and resolution

## Release Checklist

Use this checklist for each release:

- [ ] All tests pass (`just test-all`)
- [ ] Code is formatted (`just format`)
- [ ] Linter passes (`just lint`)
- [ ] CHANGELOG.md is updated
- [ ] Version bumped in extension.yml
- [ ] Changes committed and pushed
- [ ] Tag created and pushed
- [ ] GitHub Actions workflow completes successfully
- [ ] Release appears on GitHub with all binaries
- [ ] Installation tested via `gh extension install`
- [ ] Basic functionality tested
- [ ] Release announced (if applicable)
- [ ] Post-release tasks completed

## Support

For questions about the release process:
- Open an issue on GitHub
- Check existing releases for reference
- Review GitHub Actions workflow logs
