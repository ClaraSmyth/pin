package main

import "github.com/charmbracelet/lipgloss"

type ListStyles struct {
	Title      lipgloss.Style
	Selected   lipgloss.Style
	Unselected lipgloss.Style
}

var styles = ListStyles{
	Title:      lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Width(20),
	Selected:   lipgloss.NewStyle().Foreground(lipgloss.Color("#A46060")),
	Unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("237")),
}
