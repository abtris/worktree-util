# Release Process

This document describes how to create a new release of worktree-util.

## Prerequisites

1. **Create a Homebrew Tap Repository** (one-time setup):
   ```bash
   # Create a new repository on GitHub named: homebrew-tap
   # Repository URL: https://github.com/abtris/homebrew-tap
   ```

2. **Set up GitHub Token for Homebrew** (one-time setup):
   - Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Generate new token with `repo` scope
   - Add it as a repository secret named `HOMEBREW_TAP_GITHUB_TOKEN`:
     - Go to repository Settings → Secrets and variables → Actions
     - Click "New repository secret"
     - Name: `HOMEBREW_TAP_GITHUB_TOKEN`
     - Value: Your personal access token

## Creating a Release

1. **Update version** (if needed):
   - The version is automatically determined from the git tag

2. **Create and push a tag**:
   ```bash
   # Make sure you're on main branch and up to date
   git checkout main
   git pull

   # Create a new tag (use semantic versioning)
   git tag -a v0.1.0 -m "Release v0.1.0"

   # Push the tag
   git push origin v0.1.0
   ```

3. **GitHub Actions will automatically**:
   - Build binaries for:
     - Linux (amd64, arm64)
     - macOS (amd64, arm64)
     - Windows (amd64)
   - Create a GitHub release with:
     - Release notes
     - Downloadable binaries
     - Checksums
   - Update the Homebrew tap with the new formula

4. **Verify the release**:
   - Check the GitHub Actions workflow: https://github.com/abtris/worktree-util/actions
   - Check the GitHub release: https://github.com/abtris/worktree-util/releases
   - Check the Homebrew tap: https://github.com/abtris/homebrew-tap

## Installing from Homebrew

After the release is published, users can install with:

```bash
# Add the tap (first time only)
brew tap abtris/tap

# Install
brew install worktree-util

# Or in one command
brew install abtris/tap/worktree-util
```

## Testing Locally

To test the release process locally without publishing:

```bash
# Install GoReleaser
brew install goreleaser

# Test the build
goreleaser build --snapshot --clean

# Test the full release process (without publishing)
goreleaser release --snapshot --clean
```

## Troubleshooting

### Homebrew formula not updating

If the Homebrew formula doesn't update automatically:
1. Check that `HOMEBREW_TAP_GITHUB_TOKEN` secret is set correctly
2. Check the GitHub Actions logs for errors
3. Manually update the formula in the homebrew-tap repository

### Build failures

If builds fail:
1. Check the GitHub Actions logs
2. Test locally with `goreleaser build --snapshot --clean`
3. Ensure all dependencies are in go.mod

