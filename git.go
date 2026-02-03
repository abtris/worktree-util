package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

const defaultWorktreeDir = ".worktrees"

// Worktree represents a git worktree
type Worktree struct {
	Path   string
	Branch string
	Commit string
	IsMain bool
}

// Branch represents a git branch (local or remote)
type Branch struct {
	Name     string
	IsRemote bool
}

// GetRepoRoot returns the root directory of the git repository
func GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if strings.Contains(errMsg, "not a git repository") {
			return "", fmt.Errorf("not a git repository")
		}
		return "", fmt.Errorf("git error: %s", errMsg)
	}

	return strings.TrimSpace(out.String()), nil
}

// GenerateWorktreePath generates a path for a worktree based on branch name
func GenerateWorktreePath(branch string) (string, error) {
	repoRoot, err := GetRepoRoot()
	if err != nil {
		return "", err
	}

	// Sanitize branch name for use as directory name
	// Replace slashes and other problematic characters
	sanitized := strings.ReplaceAll(branch, "/", "-")
	sanitized = strings.ReplaceAll(sanitized, " ", "-")
	sanitized = strings.ReplaceAll(sanitized, "\\", "-")

	// Create path: <repo-root>/.worktrees/<sanitized-branch-name>
	worktreePath := filepath.Join(repoRoot, defaultWorktreeDir, sanitized)

	return worktreePath, nil
}

// ListWorktrees returns all git worktrees in the current repository
func ListWorktrees() ([]Worktree, error) {
	// First check if we're in a git repository
	_, err := GetRepoRoot()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return nil, fmt.Errorf("failed to list worktrees: %s", errMsg)
		}
		return nil, fmt.Errorf("failed to list worktrees")
	}

	worktrees := parseWorktrees(out.String())

	// Even a regular git repo has at least one worktree (the main one)
	// If we get here with no worktrees, something is wrong
	if len(worktrees) == 0 {
		return []Worktree{}, nil
	}

	return worktrees, nil
}

// parseWorktrees parses the porcelain output of git worktree list
func parseWorktrees(output string) []Worktree {
	var worktrees []Worktree
	var current Worktree
	
	lines := strings.Split(strings.TrimSpace(output), "\n")
	
	for _, line := range lines {
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = Worktree{}
			}
			continue
		}
		
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}
		
		key := parts[0]
		value := parts[1]
		
		switch key {
		case "worktree":
			current.Path = value
		case "HEAD":
			current.Commit = value
		case "branch":
			// Remove refs/heads/ prefix
			current.Branch = strings.TrimPrefix(value, "refs/heads/")
		case "bare":
			// Skip bare repositories
		case "detached":
			current.Branch = "detached"
		}
	}
	
	// Add the last worktree if exists
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}
	
	// Mark the first one as main
	if len(worktrees) > 0 {
		worktrees[0].IsMain = true
	}
	
	return worktrees
}

// AddWorktree creates a new worktree
func AddWorktree(path, branch string, createBranch bool) error {
	args := []string{"worktree", "add"}
	
	if createBranch {
		args = append(args, "-b", branch)
	}
	
	args = append(args, path)
	
	if !createBranch && branch != "" {
		args = append(args, branch)
	}
	
	cmd := exec.Command("git", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add worktree: %s", stderr.String())
	}
	
	return nil
}

// RemoveWorktree removes a worktree
func RemoveWorktree(path string, force bool) error {
	args := []string{"worktree", "remove"}
	
	if force {
		args = append(args, "--force")
	}
	
	args = append(args, path)
	
	cmd := exec.Command("git", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove worktree: %s", stderr.String())
	}
	
	return nil
}

// Title returns the title for the list item
func (w Worktree) Title() string {
	if w.IsMain {
		return fmt.Sprintf("📁 %s (main)", w.Path)
	}
	return fmt.Sprintf("📁 %s", w.Path)
}

// Description returns the description for the list item
func (w Worktree) Description() string {
	if w.Branch != "" {
		return fmt.Sprintf("Branch: %s | Commit: %.7s", w.Branch, w.Commit)
	}
	return fmt.Sprintf("Commit: %.7s", w.Commit)
}

// FilterValue returns the value to filter on
func (w Worktree) FilterValue() string {
	return w.Path + " " + w.Branch
}

// Title returns the title for the branch list item
func (b Branch) Title() string {
	if b.IsRemote {
		return fmt.Sprintf("🌐 %s", b.Name)
	}
	return fmt.Sprintf("📌 %s", b.Name)
}

// Description returns the description for the branch list item
func (b Branch) Description() string {
	if b.IsRemote {
		return "Remote branch"
	}
	return "Local branch"
}

// FilterValue returns the value to filter on
func (b Branch) FilterValue() string {
	return b.Name
}

// GetLocalBranches returns a list of all local branches
func GetLocalBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "--format=%(refname:short)")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		return nil, fmt.Errorf("failed to list local branches: %s", errMsg)
	}

	output := strings.TrimSpace(out.String())
	if output == "" {
		return []string{}, nil
	}

	branches := strings.Split(output, "\n")
	return branches, nil
}

// GetRemoteBranches returns a list of all remote branches
func GetRemoteBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "-r", "--format=%(refname:short)")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		return nil, fmt.Errorf("failed to list remote branches: %s", errMsg)
	}

	output := strings.TrimSpace(out.String())
	if output == "" {
		return []string{}, nil
	}

	branches := strings.Split(output, "\n")
	// Filter out HEAD references
	var filtered []string
	for _, branch := range branches {
		if !strings.Contains(branch, "HEAD") {
			filtered = append(filtered, branch)
		}
	}
	return filtered, nil
}

// GetAllBranches returns a combined list of all local and remote branches
func GetAllBranches() ([]Branch, error) {
	var allBranches []Branch

	// Get local branches
	localBranches, err := GetLocalBranches()
	if err != nil {
		return nil, err
	}

	for _, name := range localBranches {
		allBranches = append(allBranches, Branch{
			Name:     name,
			IsRemote: false,
		})
	}

	// Get remote branches
	remoteBranches, err := GetRemoteBranches()
	if err != nil {
		return nil, err
	}

	for _, name := range remoteBranches {
		allBranches = append(allBranches, Branch{
			Name:     name,
			IsRemote: true,
		})
	}

	return allBranches, nil
}

// CreateWorktreeFromBranch creates a new worktree from an existing local or remote branch
// branchName can be a local branch name (e.g., "feature") or a remote branch (e.g., "origin/feature")
// Returns the path where the worktree was created
func CreateWorktreeFromBranch(branchName string) (string, error) {
	branchName = strings.TrimSpace(branchName)
	if branchName == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}

	// Get local and remote branches
	localBranches, err := GetLocalBranches()
	if err != nil {
		return "", err
	}

	remoteBranches, err := GetRemoteBranches()
	if err != nil {
		return "", err
	}

	// Check if branch exists locally
	isLocal := false
	for _, branch := range localBranches {
		if branch == branchName {
			isLocal = true
			break
		}
	}

	// Check if branch exists remotely
	isRemote := false
	remoteBranchName := ""
	localBranchName := branchName

	// If branchName contains a slash, it might be a remote branch reference
	if strings.Contains(branchName, "/") {
		for _, branch := range remoteBranches {
			if branch == branchName {
				isRemote = true
				remoteBranchName = branchName
				// Extract local branch name (e.g., "origin/feature" -> "feature")
				parts := strings.SplitN(branchName, "/", 2)
				if len(parts) == 2 {
					localBranchName = parts[1]
				}
				break
			}
		}
	} else {
		// Check if there's a remote branch with this name
		for _, branch := range remoteBranches {
			// Check for origin/<branchName> or any other remote
			if strings.HasSuffix(branch, "/"+branchName) {
				isRemote = true
				remoteBranchName = branch
				break
			}
		}
	}

	if !isLocal && !isRemote {
		return "", fmt.Errorf("branch '%s' not found in local or remote branches", branchName)
	}

	// Generate path for the worktree
	// Use the local branch name for path generation
	path, err := GenerateWorktreePath(localBranchName)
	if err != nil {
		return "", err
	}

	// Create the worktree
	var cmd *exec.Cmd
	if isLocal {
		// Use existing local branch
		cmd = exec.Command("git", "worktree", "add", path, branchName)
	} else {
		// Create local tracking branch from remote
		// git worktree add <path> -b <local-name> --track <remote-branch>
		cmd = exec.Command("git", "worktree", "add", path, "-b", localBranchName, "--track", remoteBranchName)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to create worktree: %s", stderr.String())
	}

	return path, nil
}

