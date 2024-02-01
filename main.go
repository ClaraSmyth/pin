package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		themes := GetThemes()

		if len(themes) == 0 {
			return
		}

		for _, v := range themes {
			theme := v.(Theme)

			if theme.Name == args[0] {
				applyTheme(theme)
				return
			}
		}

		return
	}

	m := newModel()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
