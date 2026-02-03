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

