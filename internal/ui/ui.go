package ui

import (
	"fmt"
	"log"
	"strings"
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
				inputStr := m.input.Value()
				trimmed := strings.TrimSpace(inputStr)
				if len(trimmed) > 0 {
					newHabit := &model.Habit{
						ID:        uuid.New().String(),
						Name:      trimmed,
						CreatedAt: time.Now(),
					}
					m.storage.AddHabit(newHabit)
					m.habits = m.storage.GetHabits()
				}
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
		case "d":
			if len(m.habits) > 0 {
				habitToDelete := m.habits[m.cursor]
				if err := m.storage.DeleteHabit(habitToDelete.ID); err != nil {
					log.Printf("Error deleting habit: %v\n", err)
				} else {
					m.habits = m.storage.GetHabits()
					if m.cursor >= len(m.habits) {
						m.cursor = len(m.habits) - 1
					}
				}
			}
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
			TextInputPromptStyle.Render(m.input.View()),
		)
	}

	s := TitleStyle.Render("Habit Tracker") + "\n\n"
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

		line := fmt.Sprintf("%s [%s] %s (%.2f%%)", cursor, checked, habit.Name, completionRate)
		if m.cursor == i {
			s += SelectedListItemStyle.Render(line)
		} else {
			s += ListItemStyle.Render(line)
		}
		s += "\n"
	}
	s += HelpStyle.Render("\nUse UP/DOWN arrow keys to navigate.\nspace or enter to toggle habit status.\n'a' to add a new habit.\n'd' to delete the selected habit.\n'q' to exit.\n")
	return s
}
