package main

import "github.com/charmbracelet/lipgloss"

type ListStyles struct {
	Title      lipgloss.Style
	Selected   lipgloss.Style
	Unselected lipgloss.Style
	NoItems    lipgloss.Style
	StatusBar  lipgloss.Style
	TitleBar   lipgloss.Style
}

type Styles struct {
	BaseStyles    ListStyles
	FocusedStyles ListStyles
}

var styles = Styles{
	BaseStyles: ListStyles{
		Title:      lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Width(20),
		Selected:   lipgloss.NewStyle().Foreground(lipgloss.Color("#A46060")),
		Unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("237")),
	},
	FocusedStyles: ListStyles{
		Title:      lipgloss.NewStyle().Background(lipgloss.Color("#1B50CF")).Width(20),
		Selected:   lipgloss.NewStyle().Foreground(lipgloss.Color("#F44336")),
		Unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("237")),
	},
}
