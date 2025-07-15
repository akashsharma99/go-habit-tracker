package ui

import (
	"fmt"
	"time"

	"habit-tracker/internal/model"
	"habit-tracker/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

type TuiModel struct {
	habits   []*model.Habit   // items to be displayed in the UI
	cursor   int              // current cursor position in the list
	selected map[int]struct{} // selected item
	storage  *storage.Storage
}

func InitializeModel(s *storage.Storage) TuiModel {
	return TuiModel{
		habits:   s.GetHabits(),
		cursor:   0,
		selected: make(map[int]struct{}),
		storage:  s,
	}
}

func (m TuiModel) Init() tea.Cmd {
	return nil
}

func (m TuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.habits)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.habits[m.cursor].ToggleToday()
			m.storage.UpdateCompletion(m.habits[m.cursor].ID, time.Now().Format("2006-01-02"), m.habits[m.cursor].IsCompletedToday())
		}
	}

	return m, nil
}

func (m TuiModel) View() string {
	s := "Habit Tracker\n\n"
	for i, habit := range m.habits {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "

		if habit.IsCompletedToday() {
			checked = "✓"
		}
		completionRate := habit.GetCompletionRate()
		s += fmt.Sprintf("%s [%s] %s (%.2f%%)\n", cursor, checked, habit.Name, completionRate)
	}
	s += "\nUse UP/DOWN arrow keys to navigate, and space or enter to toggle habit status. q to exit.\n"
	return s
}
