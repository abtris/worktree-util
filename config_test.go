package main

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.WorktreeDir != ".worktrees" {
		t.Errorf("DefaultConfig().WorktreeDir = %v, want .worktrees", config.WorktreeDir)
	}

	if config.CopyFiles != nil && len(config.CopyFiles) != 0 {
		t.Errorf("DefaultConfig().CopyFiles should be empty, got %v", config.CopyFiles)
	}
}

func TestLoadConfig_NoFile(t *testing.T) {
	// Create a temporary home directory
	tempHome := t.TempDir()
	
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Set temporary HOME
	os.Setenv("HOME", tempHome)

	config, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() should not error when file doesn't exist, got: %v", err)
	}

	// Should return default config
	if config.WorktreeDir != ".worktrees" {
		t.Errorf("LoadConfig() with no file should return default, got WorktreeDir = %v", config.WorktreeDir)
	}
}

func TestLoadConfig_WithCopyFiles(t *testing.T) {
	// Create a temporary home directory
	tempHome := t.TempDir()
	
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Set temporary HOME
	os.Setenv("HOME", tempHome)

	// Create config directory
	configDir := filepath.Join(tempHome, ".config", "worktree-util")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create config file with copy_files
	configPath := filepath.Join(configDir, "config.yml")
	configContent := `worktree_dir: my-worktrees
copy_files:
  - .env
  - .env.local
  - config/local.yml
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	config, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() failed: %v", err)
	}

	if config.WorktreeDir != "my-worktrees" {
		t.Errorf("LoadConfig().WorktreeDir = %v, want my-worktrees", config.WorktreeDir)
	}

	expectedFiles := []string{".env", ".env.local", "config/local.yml"}
	if len(config.CopyFiles) != len(expectedFiles) {
		t.Errorf("LoadConfig().CopyFiles length = %v, want %v", len(config.CopyFiles), len(expectedFiles))
	}

	for i, file := range expectedFiles {
		if i >= len(config.CopyFiles) || config.CopyFiles[i] != file {
			t.Errorf("LoadConfig().CopyFiles[%d] = %v, want %v", i, config.CopyFiles[i], file)
		}
	}
}

func TestSaveConfig(t *testing.T) {
	// Create a temporary home directory
	tempHome := t.TempDir()
	
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Set temporary HOME
	os.Setenv("HOME", tempHome)

	config := &Config{
		WorktreeDir: "custom-worktrees",
		CopyFiles:   []string{".env", "config.yml"},
	}

	if err := SaveConfig(config); err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// Verify file was created
	configPath := filepath.Join(tempHome, ".config", "worktree-util", "config.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created")
	}

	// Read and verify content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var loadedConfig Config
	if err := yaml.Unmarshal(data, &loadedConfig); err != nil {
		t.Fatalf("Failed to parse config file: %v", err)
	}

	if loadedConfig.WorktreeDir != config.WorktreeDir {
		t.Errorf("Saved WorktreeDir = %v, want %v", loadedConfig.WorktreeDir, config.WorktreeDir)
	}

	if len(loadedConfig.CopyFiles) != len(config.CopyFiles) {
		t.Errorf("Saved CopyFiles length = %v, want %v", len(loadedConfig.CopyFiles), len(config.CopyFiles))
	}
}

