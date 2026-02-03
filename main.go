package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Handle version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("worktree-util version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built at: %s\n", date)
		os.Exit(0)
	}

	// Handle config command
	if len(os.Args) > 1 && os.Args[1] == "config" {
		HandleConfigCommand(os.Args[2:])
		os.Exit(0)
	}

	// Handle help flag
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help") {
		printHelp()
		os.Exit(0)
	}

	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
		fmt.Println("Using default configuration")
	}

	// Set global config
	appConfig = config

	// Start TUI
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf("worktree-util version %s\n\n", version)
	fmt.Println("A TUI for managing Git worktrees")
	fmt.Println("\nUsage:")
	fmt.Println("  worktree-util              Start the TUI")
	fmt.Println("  worktree-util config       Manage configuration")
	fmt.Println("  worktree-util --version    Show version information")
	fmt.Println("  worktree-util --help       Show this help message")
	fmt.Println("\nConfig commands:")
	fmt.Println("  worktree-util config              Show current configuration")
	fmt.Println("  worktree-util config init         Create default config file")
	fmt.Println("  worktree-util config set <key> <value>")
	fmt.Println("                                    Set a configuration value")
	fmt.Println("  worktree-util config get <key>    Get a configuration value")
	fmt.Println("  worktree-util config add-copy-file <file>")
	fmt.Println("                                    Add a file to copy_files list")
	fmt.Println("  worktree-util config remove-copy-file <file>")
	fmt.Println("                                    Remove a file from copy_files list")
	fmt.Println("\nFor more information, visit: https://github.com/abtris/worktree-util")
}

