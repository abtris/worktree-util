package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	WorktreeDir string `yaml:"worktree_dir"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		WorktreeDir: ".worktrees",
	}
}

// LoadConfig loads configuration from ~/.config/worktree-util/config.yml
// If the file doesn't exist, it returns the default configuration
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// If we can't get home dir, just use defaults
		return config, nil
	}

	// Construct config file path
	configPath := filepath.Join(homeDir, ".config", "worktree-util", "config.yml")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config file doesn't exist, use defaults
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		// If we can't read the file, use defaults
		return config, nil
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		// If parsing fails, use defaults
		return config, nil
	}

	// Validate and set defaults for empty values
	if config.WorktreeDir == "" {
		config.WorktreeDir = ".worktrees"
	}

	return config, nil
}

// SaveConfig saves the configuration to ~/.config/worktree-util/config.yml
func SaveConfig(config *Config) error {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Construct config directory path
	configDir := filepath.Join(homeDir, ".config", "worktree-util")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Construct config file path
	configPath := filepath.Join(configDir, "config.yml")

	// Marshal config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return err
	}

	return nil
}

