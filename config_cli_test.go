package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddCopyFile(t *testing.T) {
	// Create a temporary home directory
	tempHome := t.TempDir()

	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set temporary HOME
	os.Setenv("HOME", tempHome)

	// Create initial config
	config := DefaultConfig()
	if err := SaveConfig(config); err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	// Add a file
	addCopyFile(".env")

	// Load config and verify
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(loadedConfig.CopyFiles) != 1 {
		t.Errorf("Expected 1 file in copy_files, got %d", len(loadedConfig.CopyFiles))
	}

	if loadedConfig.CopyFiles[0] != ".env" {
		t.Errorf("Expected .env in copy_files, got %s", loadedConfig.CopyFiles[0])
	}

	// Try adding the same file again (should not duplicate)
	addCopyFile(".env")
	loadedConfig, _ = LoadConfig()
	if len(loadedConfig.CopyFiles) != 1 {
		t.Errorf("Expected 1 file in copy_files after duplicate add, got %d", len(loadedConfig.CopyFiles))
	}
}

func TestRemoveCopyFile(t *testing.T) {
	// Create a temporary home directory
	tempHome := t.TempDir()

	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set temporary HOME
	os.Setenv("HOME", tempHome)

	// Create initial config with files
	config := &Config{
		WorktreeDir: ".worktrees",
		CopyFiles:   []string{".env", ".env.local"},
	}
	if err := SaveConfig(config); err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	// Remove a file
	removeCopyFile(".env")

	// Load config and verify
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(loadedConfig.CopyFiles) != 1 {
		t.Errorf("Expected 1 file in copy_files, got %d", len(loadedConfig.CopyFiles))
	}

	if loadedConfig.CopyFiles[0] != ".env.local" {
		t.Errorf("Expected .env.local in copy_files, got %s", loadedConfig.CopyFiles[0])
	}
}

func TestSetConfig(t *testing.T) {
	// Create a temporary home directory
	tempHome := t.TempDir()

	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set temporary HOME
	os.Setenv("HOME", tempHome)

	// Create initial config
	config := DefaultConfig()
	if err := SaveConfig(config); err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	// Set worktree_dir
	setConfig("worktree_dir", "custom-worktrees")

	// Load config and verify
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedConfig.WorktreeDir != "custom-worktrees" {
		t.Errorf("Expected worktree_dir to be 'custom-worktrees', got '%s'", loadedConfig.WorktreeDir)
	}
}

func TestInitConfig(t *testing.T) {
	// Create a temporary home directory
	tempHome := t.TempDir()

	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set temporary HOME
	os.Setenv("HOME", tempHome)

	// Initialize config
	initConfig()

	// Verify config file was created
	configPath := filepath.Join(tempHome, ".config", "worktree-util", "config.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created")
	}

	// Load and verify default values
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.WorktreeDir != ".worktrees" {
		t.Errorf("Expected default worktree_dir '.worktrees', got '%s'", config.WorktreeDir)
	}
}
