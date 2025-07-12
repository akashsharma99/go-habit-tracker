package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"habit-tracker/internal/model"
	"habit-tracker/internal/storage"
)

type item struct {
	habit *model.Habit
}

func (i item) Title() string {
	style := lipgloss.NewStyle().PaddingLeft(1)
	checkmark := "[ ]" // Empty checkbox
	if i.habit.IsCompletedToday() {
		checkmark = "[✅]" // Checked checkbox
		style = style.Foreground(lipgloss.Color("#fb5ae0c8"))
	}
	return style.Render(fmt.Sprintf("%s %s (%.1f%% completed)", checkmark, i.habit.Name, i.habit.GetCompletionRate()))
}

func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.habit.Name }

type Model struct {
	list     list.Model
	storage  *storage.Storage
	err      error
	adding   bool
	newHabit strings.Builder
}

func NewModel(s *storage.Storage) *Model {
	items := []list.Item{}
	habits := s.GetHabits()
	for _, habit := range habits {
		items = append(items, item{habit: habit})
	}

	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(2)
	delegate.SetSpacing(1)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#fb5ae0c8")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#fb5ae0c8")).
		Padding(0, 1)

	l := list.New(items, delegate, 0, 0)
	l.Title = "📋 Habit Tracker"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false) // Disable filtering to ensure all items are shown
	l.SetShowStatusBar(false)    // Hide the status bar
	l.SetShowPagination(true)    // Show pagination for better navigation

	return &Model{
		list:    l,
		storage: s,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.adding {
			return m.updateAddingState(msg)
		}

		switch msg.String() {
		case "k", "up":
			m.list.CursorUp()
			return m, nil
		case "j", "down":
			m.list.CursorDown()
			return m, nil
		case " ":
			return m.toggleSelectedHabit()
		default:
			return m.updateNormalState(msg)
		}

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) updateAddingState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		habitName := m.newHabit.String()
		if habitName != "" {
			habit := model.NewHabit(habitName)
			m.storage.AddHabit(habit)
			m.list.InsertItem(len(m.list.Items()), item{habit: habit})
		}
		m.adding = false
		m.newHabit.Reset()
		return m, nil

	case tea.KeyEsc:
		m.adding = false
		m.newHabit.Reset()
		return m, nil

	case tea.KeyBackspace:
		if m.newHabit.Len() > 0 {
			str := m.newHabit.String()
			m.newHabit.Reset()
			m.newHabit.WriteString(str[:len(str)-1])
		}
		return m, nil

	default:
		m.newHabit.WriteString(msg.String())
		return m, nil
	}
}

func (m Model) updateNormalState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "a":
		m.adding = true
		return m, nil

	case " ":
		i := m.list.SelectedItem()
		if i == nil {
			return m, nil
		}
		habit := i.(item).habit
		habit.ToggleToday()
		err := m.storage.UpdateHabit(habit)
		if err != nil {
			m.err = err
			return m, nil
		}

		// Refresh the list items to show the updated status
		items := []list.Item{}
		for _, h := range m.storage.GetHabits() {
			items = append(items, item{habit: h})
		}
		m.list.SetItems(items)
		return m, nil

	case "d":
		i := m.list.SelectedItem()
		if i == nil {
			return m, nil
		}
		habit := i.(item).habit
		err := m.storage.DeleteHabit(habit.ID)
		if err != nil {
			m.err = err
			return m, nil
		}

		// Refresh the list items
		items := []list.Item{}
		for _, h := range m.storage.GetHabits() {
			items = append(items, item{habit: h})
		}
		m.list.SetItems(items)
		return m, nil
	}

	return m, nil
}

func (m Model) toggleSelectedHabit() (tea.Model, tea.Cmd) {
	i := m.list.SelectedItem()
	if i == nil {
		return m, nil
	}

	habit := i.(item).habit
	habit.ToggleToday()
	err := m.storage.UpdateHabit(habit)
	if err != nil {
		m.err = err
		return m, nil
	}

	// Refresh the list items to show the updated status
	items := []list.Item{}
	for _, h := range m.storage.GetHabits() {
		items = append(items, item{habit: h})
	}
	m.list.SetItems(items)

	// Maintain the same selection
	if m.list.Index() >= len(items) {
		m.list.Select(len(items) - 1)
	}

	return m, nil
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	if m.adding {
		return fmt.Sprintf(
			"%s\n\nEnter habit name (press Enter to save, Esc to cancel):\n> %s",
			m.list.View(),
			m.newHabit.String(),
		)
	}

	controls := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render("\nControls: ↑/↓: select habit • space: toggle completion • a: add habit • d: delete habit • q: quit")

	return fmt.Sprintf("%s%s", m.list.View(), controls)
}
