package ui

import (
	"fmt"
	"strings"
	"time"

	"habit-tracker/internal/model"
	"habit-tracker/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	listHabits state = iota
	promptAddHabit
)

type Model struct {
	storage  *storage.Storage
	state    state
	input    string
	habits   []*model.Habit
	selected int
}

func NewModel(s *storage.Storage) *Model {
	return &Model{
		storage:  s,
		state:    listHabits,
		habits:   s.GetHabits(),
		selected: 0,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Quit on 'q' in any state
		if msg.String() == "q" {
			return m, tea.Quit
		}
		switch m.state {
		case listHabits:
			switch msg.String() {
			case "a":
				m.state = promptAddHabit
				m.input = ""
				return m, nil
			case "up":
				if m.selected > 0 {
					m.selected--
				}
				return m, nil
			case "down":
				if m.selected < len(m.habits)-1 {
					m.selected++
				}
				return m, nil
			case " ":
				if len(m.habits) > 0 {
					h := m.habits[m.selected]
					h.ToggleToday()
					today := time.Now().Format("2006-01-02")
					completed := h.IsCompletedToday()
					m.storage.UpdateCompletion(h.ID, today, completed)
					m.habits = m.storage.GetHabits()
				}
				return m, nil
			}
		case promptAddHabit:
			switch msg.Type {
			case tea.KeyEnter:
				if strings.TrimSpace(m.input) != "" {
					habit := model.NewHabit(m.input)
					m.storage.AddHabit(habit)
					m.habits = m.storage.GetHabits()
					m.state = listHabits
					m.selected = len(m.habits) - 1
				}
				return m, nil
			case tea.KeyEsc:
				m.state = listHabits
				return m, nil
			case tea.KeyBackspace:
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
				return m, nil
			default:
				if msg.Type == tea.KeyRunes {
					m.input += msg.String()
				}
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		return m, nil
	}
	return m, nil
}

func (m *Model) View() string {
	switch m.state {
	case promptAddHabit:
		return "Enter new habit name: " + m.input
	case listHabits:
		if len(m.habits) == 0 {
			return "No habits yet. Press 'a' to add a habit."
		}
		var b strings.Builder
		b.WriteString("Habits:\n")
		for i, h := range m.habits {
			checked := "[ ]"
			if h.IsCompletedToday() {
				checked = "[x]"
			}
			percent := h.GetCompletionRate()
			selector := "  "
			if i == m.selected {
				selector = "> "
			}
			b.WriteString(fmt.Sprintf("%s%s %s (%.0f%%)\n", selector, checked, h.Name, percent))
		}
		b.WriteString("\nPress 'a' to add, space to toggle completion, up/down to select.")
		return b.String()
	}
	return ""
}
