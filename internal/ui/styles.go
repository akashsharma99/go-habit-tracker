package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingLeft(2).
		PaddingRight(2)

	ListItemStyle = lipgloss.NewStyle().
		PaddingLeft(2)

	SelectedListItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		PaddingLeft(1)

	HelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))
)
