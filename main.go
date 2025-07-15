package main

import (
	"fmt"
	"os"
	"path/filepath"

	"habit-tracker/internal/storage"
	"habit-tracker/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	dataDir := filepath.Join(homeDir, ".habit-tracker/app.log")
	f, err := tea.LogToFile(dataDir, "info")
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	s, err := storage.NewStorage()
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.InitializeModel(s), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
