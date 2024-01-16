package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := newModel()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
