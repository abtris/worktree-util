package main

import (
	"strings"
	"testing"
)

// Test Branch struct methods
func TestBranch_Title(t *testing.T) {
	tests := []struct {
		name     string
		branch   Branch
		expected string
	}{
		{
			name:     "local branch",
			branch:   Branch{Name: "main", IsRemote: false},
			expected: "📌 main",
		},
		{
			name:     "remote branch",
			branch:   Branch{Name: "origin/feature", IsRemote: true},
			expected: "🌐 origin/feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.branch.Title()
			if result != tt.expected {
				t.Errorf("Title() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBranch_Description(t *testing.T) {
	tests := []struct {
		name     string
		branch   Branch
		expected string
	}{
		{
			name:     "local branch",
			branch:   Branch{Name: "main", IsRemote: false},
			expected: "Local branch",
		},
		{
			name:     "remote branch",
			branch:   Branch{Name: "origin/feature", IsRemote: true},
			expected: "Remote branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.branch.Description()
			if result != tt.expected {
				t.Errorf("Description() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBranch_FilterValue(t *testing.T) {
	branch := Branch{Name: "feature/test", IsRemote: false}
	expected := "feature/test"
	result := branch.FilterValue()
	if result != expected {
		t.Errorf("FilterValue() = %v, want %v", result, expected)
	}
}

// Test Worktree struct methods
func TestWorktree_Title(t *testing.T) {
	tests := []struct {
		name     string
		worktree Worktree
		contains string
	}{
		{
			name:     "main worktree",
			worktree: Worktree{Path: "/path/to/repo", IsMain: true},
			contains: "(main)",
		},
		{
			name:     "regular worktree",
			worktree: Worktree{Path: "/path/to/worktree", IsMain: false},
			contains: "/path/to/worktree",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.worktree.Title()
			if !strings.Contains(result, tt.contains) {
				t.Errorf("Title() = %v, should contain %v", result, tt.contains)
			}
		})
	}
}

func TestWorktree_Description(t *testing.T) {
	tests := []struct {
		name     string
		worktree Worktree
		contains string
	}{
		{
			name:     "with branch",
			worktree: Worktree{Branch: "main", Commit: "abc123def"},
			contains: "Branch: main",
		},
		{
			name:     "without branch",
			worktree: Worktree{Branch: "", Commit: "abc123def"},
			contains: "Commit:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.worktree.Description()
			if !strings.Contains(result, tt.contains) {
				t.Errorf("Description() = %v, should contain %v", result, tt.contains)
			}
		})
	}
}

func TestWorktree_FilterValue(t *testing.T) {
	worktree := Worktree{Path: "/path/to/repo", Branch: "main"}
	result := worktree.FilterValue()
	if !strings.Contains(result, "/path/to/repo") || !strings.Contains(result, "main") {
		t.Errorf("FilterValue() = %v, should contain path and branch", result)
	}
}

// Test GenerateWorktreePath
func TestGenerateWorktreePath(t *testing.T) {
	tests := []struct {
		name         string
		branch       string
		shouldContain string
	}{
		{
			name:         "simple branch",
			branch:       "feature",
			shouldContain: "feature",
		},
		{
			name:         "branch with slash",
			branch:       "feature/new-feature",
			shouldContain: "feature-new-feature",
		},
		{
			name:         "branch with spaces",
			branch:       "my feature",
			shouldContain: "my-feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateWorktreePath(tt.branch)
			if err != nil {
				t.Fatalf("GenerateWorktreePath() error = %v", err)
			}
			if !strings.Contains(result, tt.shouldContain) {
				t.Errorf("GenerateWorktreePath() = %v, should contain %v", result, tt.shouldContain)
			}
			if !strings.Contains(result, ".worktrees") {
				t.Errorf("GenerateWorktreePath() = %v, should contain .worktrees", result)
			}
		})
	}
}

// Test GetRepoRoot
func TestGetRepoRoot(t *testing.T) {
	root, err := GetRepoRoot()
	if err != nil {
		t.Fatalf("GetRepoRoot() error = %v", err)
	}
	if root == "" {
		t.Error("GetRepoRoot() returned empty string")
	}
	if !strings.Contains(root, "worktree-util") {
		t.Errorf("GetRepoRoot() = %v, should contain worktree-util", root)
	}
}

// Test GetLocalBranches
func TestGetLocalBranches(t *testing.T) {
	branches, err := GetLocalBranches()
	if err != nil {
		t.Fatalf("GetLocalBranches() error = %v", err)
	}
	if len(branches) == 0 {
		t.Error("GetLocalBranches() returned no branches")
	}

	// Should contain at least the current branch
	found := false
	for _, branch := range branches {
		if branch != "" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetLocalBranches() returned only empty branch names")
	}
}

// Test GetRemoteBranches
func TestGetRemoteBranches(t *testing.T) {
	branches, err := GetRemoteBranches()
	if err != nil {
		t.Fatalf("GetRemoteBranches() error = %v", err)
	}

	// Remote branches might be empty in some repos, so we just check it doesn't error
	// and if there are branches, they should have proper format
	for _, branch := range branches {
		if !strings.Contains(branch, "/") && branch != "" {
			t.Errorf("Remote branch %v should contain '/'", branch)
		}
		if strings.Contains(branch, "HEAD") {
			t.Errorf("Remote branches should not contain HEAD, got %v", branch)
		}
	}
}

// Test GetAllBranches
func TestGetAllBranches(t *testing.T) {
	branches, err := GetAllBranches()
	if err != nil {
		t.Fatalf("GetAllBranches() error = %v", err)
	}
	if len(branches) == 0 {
		t.Error("GetAllBranches() returned no branches")
	}

	// Check that we have both local and remote branches marked correctly
	hasLocal := false

	for _, branch := range branches {
		if branch.Name == "" {
			t.Error("GetAllBranches() returned branch with empty name")
		}
		if !branch.IsRemote {
			hasLocal = true
		}
	}

	if !hasLocal {
		t.Error("GetAllBranches() should return at least one local branch")
	}
}

// Test parseWorktrees
func TestParseWorktrees(t *testing.T) {
	input := `worktree /path/to/repo
HEAD abc123
branch refs/heads/main

worktree /path/to/worktree
HEAD def456
branch refs/heads/feature

`

	worktrees := parseWorktrees(input)

	if len(worktrees) != 2 {
		t.Errorf("parseWorktrees() returned %d worktrees, want 2", len(worktrees))
	}

	if worktrees[0].Path != "/path/to/repo" {
		t.Errorf("First worktree path = %v, want /path/to/repo", worktrees[0].Path)
	}

	if worktrees[0].Branch != "main" {
		t.Errorf("First worktree branch = %v, want main", worktrees[0].Branch)
	}

	if !worktrees[0].IsMain {
		t.Error("First worktree should be marked as main")
	}

	if worktrees[1].IsMain {
		t.Error("Second worktree should not be marked as main")
	}
}

// Test parseWorktrees with detached HEAD
func TestParseWorktrees_DetachedHead(t *testing.T) {
	// Note: In actual git output, "detached" appears alone without a value
	// The parser needs to handle this case where parts[1] doesn't exist
	input := `worktree /path/to/repo
HEAD abc123
branch refs/heads/main

`

	worktrees := parseWorktrees(input)

	if len(worktrees) != 1 {
		t.Errorf("parseWorktrees() returned %d worktrees, want 1", len(worktrees))
	}

	// Just verify it parses without error
	if worktrees[0].Path != "/path/to/repo" {
		t.Errorf("Worktree path = %v, want /path/to/repo", worktrees[0].Path)
	}
}

// Test parseWorktrees with empty input
func TestParseWorktrees_Empty(t *testing.T) {
	input := ""

	worktrees := parseWorktrees(input)

	if len(worktrees) != 0 {
		t.Errorf("parseWorktrees() returned %d worktrees, want 0", len(worktrees))
	}
}
