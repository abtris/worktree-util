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

	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
		fmt.Println("Using default configuration")
	}

	// Set global config
	appConfig = config

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

