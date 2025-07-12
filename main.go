package main

import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"
    "habit-tracker/internal/storage"
    "habit-tracker/internal/ui"
)

func main() {
    s, err := storage.NewStorage()
    if err != nil {
        fmt.Printf("Error initializing storage: %v\n", err)
        os.Exit(1)
    }

    p := tea.NewProgram(ui.NewModel(s), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error running program: %v\n", err)
        os.Exit(1)
    }
}
