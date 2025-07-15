package ui

import (
	"fmt"
	"time"

	"habit-tracker/internal/model"
	"habit-tracker/internal/storage"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type TuiModel struct {
	habits      []*model.Habit   // items to be displayed in the UI
	cursor      int              // current cursor position in the list
	selected    map[int]struct{} // selected item
	storage     *storage.Storage
	input       textinput.Model
	addingHabit bool
}

func InitializeModel(s *storage.Storage) TuiModel {
	input := textinput.New()
	input.Placeholder = "New Habit"
	input.Focus()

	return TuiModel{
		habits:   s.GetHabits(),
		cursor:   0,
		selected: make(map[int]struct{}),
		storage:  s,
		input:    input,
	}
}

func (m TuiModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.addingHabit {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				newHabit := &model.Habit{
					ID:        uuid.New().String(),
					Name:      m.input.Value(),
					CreatedAt: time.Now(),
				}
				m.storage.AddHabit(newHabit)
				m.habits = m.storage.GetHabits()
				m.addingHabit = false
				m.input.Reset()
			case "esc":
				m.addingHabit = false
				m.input.Reset()
			}
		}
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "a":
			m.addingHabit = true
			return m, nil
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
	if m.addingHabit {
		return fmt.Sprintf(
			"Enter new habit name:\n%s\n\n(press enter to save, esc to cancel)",
			m.input.View(),
		)
	}

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
	s += "\nUse UP/DOWN arrow keys to navigate, and 'space' or 'enter' to toggle habit status. 'a' to add a new habit. 'q' to exit.\n"
	return s
}
