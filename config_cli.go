package main

import (
	"fmt"
	"os"
)

// HandleConfigCommand handles all config-related CLI commands
func HandleConfigCommand(args []string) {
	if len(args) == 0 {
		// No subcommand - show current config
		showConfig()
		return
	}

	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "init":
		initConfig()
	case "set":
		if len(subArgs) < 2 {
			fmt.Println("Usage: worktree-util config set <key> <value>")
			fmt.Println("Available keys: worktree_dir")
			os.Exit(1)
		}
		setConfig(subArgs[0], subArgs[1])
	case "get":
		if len(subArgs) < 1 {
			fmt.Println("Usage: worktree-util config get <key>")
			fmt.Println("Available keys: worktree_dir, copy_files")
			os.Exit(1)
		}
		getConfig(subArgs[0])
	case "add-copy-file":
		if len(subArgs) < 1 {
			fmt.Println("Usage: worktree-util config add-copy-file <file>")
			os.Exit(1)
		}
		addCopyFile(subArgs[0])
	case "remove-copy-file":
		if len(subArgs) < 1 {
			fmt.Println("Usage: worktree-util config remove-copy-file <file>")
			os.Exit(1)
		}
		removeCopyFile(subArgs[0])
	default:
		fmt.Printf("Unknown config command: %s\n", subcommand)
		printConfigHelp()
		os.Exit(1)
	}
}

func printConfigHelp() {
	fmt.Println("\nAvailable config commands:")
	fmt.Println("  worktree-util config              Show current configuration")
	fmt.Println("  worktree-util config init         Create default config file")
	fmt.Println("  worktree-util config set <key> <value>")
	fmt.Println("                                    Set a configuration value")
	fmt.Println("  worktree-util config get <key>    Get a configuration value")
	fmt.Println("  worktree-util config add-copy-file <file>")
	fmt.Println("                                    Add a file to copy_files list")
	fmt.Println("  worktree-util config remove-copy-file <file>")
	fmt.Println("                                    Remove a file from copy_files list")
}

func showConfig() {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Current configuration:")
	fmt.Printf("  worktree_dir: %s\n", config.WorktreeDir)
	fmt.Printf("  copy_files: %v\n", config.CopyFiles)
	if len(config.CopyFiles) == 0 {
		fmt.Println("    (none)")
	} else {
		for _, file := range config.CopyFiles {
			fmt.Printf("    - %s\n", file)
		}
	}
}

func initConfig() {
	config := DefaultConfig()
	if err := SaveConfig(config); err != nil {
		fmt.Printf("Error creating config file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Config file created at ~/.config/worktree-util/config.yml")
	fmt.Println("Default configuration:")
	fmt.Printf("  worktree_dir: %s\n", config.WorktreeDir)
	fmt.Printf("  copy_files: [] (empty)\n")
}

func setConfig(key, value string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	switch key {
	case "worktree_dir":
		config.WorktreeDir = value
	default:
		fmt.Printf("Unknown config key: %s\n", key)
		fmt.Println("Available keys: worktree_dir")
		os.Exit(1)
	}

	if err := SaveConfig(config); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Set %s = %s\n", key, value)
}

func getConfig(key string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	switch key {
	case "worktree_dir":
		fmt.Println(config.WorktreeDir)
	case "copy_files":
		if len(config.CopyFiles) == 0 {
			fmt.Println("(none)")
		} else {
			for _, file := range config.CopyFiles {
				fmt.Println(file)
			}
		}
	default:
		fmt.Printf("Unknown config key: %s\n", key)
		fmt.Println("Available keys: worktree_dir, copy_files")
		os.Exit(1)
	}
}

func addCopyFile(file string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Check if file already exists in the list
	for _, f := range config.CopyFiles {
		if f == file {
			fmt.Printf("File '%s' is already in copy_files list\n", file)
			return
		}
	}

	// Add the file
	config.CopyFiles = append(config.CopyFiles, file)

	if err := SaveConfig(config); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Added '%s' to copy_files list\n", file)
}

func removeCopyFile(file string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Find and remove the file
	found := false
	newCopyFiles := []string{}
	for _, f := range config.CopyFiles {
		if f == file {
			found = true
		} else {
			newCopyFiles = append(newCopyFiles, f)
		}
	}

	if !found {
		fmt.Printf("File '%s' not found in copy_files list\n", file)
		return
	}

	config.CopyFiles = newCopyFiles

	if err := SaveConfig(config); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Removed '%s' from copy_files list\n", file)
}
