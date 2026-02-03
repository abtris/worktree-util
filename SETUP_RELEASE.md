# Setting Up Releases - Quick Start Guide

Follow these steps to set up automated releases for worktree-util.

## Step 1: Create Homebrew Tap Repository

1. Go to GitHub and create a new repository:
   - Name: `homebrew-tap`
   - Owner: `abtris`
   - Public repository
   - Don't initialize with README (GoReleaser will create it)

2. The repository URL should be: `https://github.com/abtris/homebrew-tap`

## Step 2: Create GitHub Personal Access Token

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Give it a descriptive name: "Homebrew Tap Token for worktree-util"
4. Select scopes:
   - ✅ `repo` (Full control of private repositories)
   - ✅ `workflow` (Update GitHub Action workflows)
5. Click "Generate token"
6. **Copy the token** (you won't be able to see it again!)

## Step 3: Add Token to Repository Secrets

1. Go to your worktree-util repository on GitHub
2. Navigate to: Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Add the secret:
   - Name: `HOMEBREW_TAP_GITHUB_TOKEN`
   - Value: Paste the token you copied in Step 2
5. Click "Add secret"

## Step 4: Create Your First Release

1. Make sure all your changes are committed and pushed to main:
   ```bash
   git add .
   git commit -m "Prepare for first release"
   git push origin main
   ```

2. Create and push a tag:
   ```bash
   git tag -a v0.1.0 -m "First release"
   git push origin v0.1.0
   ```

3. Watch the magic happen! 🎉
   - Go to: https://github.com/abtris/worktree-util/actions
   - You should see the "Release" workflow running
   - Wait for it to complete (usually 2-3 minutes)

## Step 5: Verify the Release

1. **Check GitHub Release:**
   - Go to: https://github.com/abtris/worktree-util/releases
   - You should see v0.1.0 with binaries for multiple platforms

2. **Check Homebrew Tap:**
   - Go to: https://github.com/abtris/homebrew-tap
   - You should see a `Formula/worktree-util.rb` file

3. **Test Installation:**
   ```bash
   brew install abtris/tap/worktree-util
   worktree-util --version
   ```

## Future Releases

For subsequent releases, just create and push a new tag:

```bash
# Update version
git tag -a v0.2.0 -m "Release v0.2.0 - Added new features"
git push origin v0.2.0
```

That's it! The GitHub Actions workflow will handle everything automatically.

## Troubleshooting

### "Resource not accessible by integration" error

This means the `HOMEBREW_TAP_GITHUB_TOKEN` is not set or doesn't have the right permissions.
- Verify the token has `repo` and `workflow` scopes
- Make sure the secret name is exactly `HOMEBREW_TAP_GITHUB_TOKEN`

### Homebrew tap repository not found

- Make sure you created the `homebrew-tap` repository
- Verify the repository is public
- Check that the owner name in `.goreleaser.yml` matches your GitHub username

### Build fails

- Check the GitHub Actions logs for specific errors
- Test locally: `goreleaser build --snapshot --clean`
- Make sure `go.mod` is up to date: `go mod tidy`

## Need Help?

- Check the [RELEASE.md](RELEASE.md) file for detailed release process
- Review GitHub Actions logs for error messages
- Test locally with GoReleaser before pushing tags

