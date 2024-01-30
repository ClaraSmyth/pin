package main

import (
	"github.com/charmbracelet/lipgloss"
)

type ListStyles struct {
	Title      lipgloss.Style
	Selected   lipgloss.Style
	Unselected lipgloss.Style
	NoItems    lipgloss.Style
	StatusBar  lipgloss.Style
	TitleBar   lipgloss.Style
}

type HelpStyles struct {
	Key  lipgloss.Style
	Desc lipgloss.Style
}

type Styles struct {
	BaseStyles    ListStyles
	FocusedStyles ListStyles
	HelpStyles    HelpStyles
}

type Colors struct {
	Base00 lipgloss.TerminalColor `yaml:"base00"`
	Base01 lipgloss.TerminalColor `yaml:"base01"`
	Base02 lipgloss.TerminalColor `yaml:"base02"`
	Base03 lipgloss.TerminalColor `yaml:"base03"`
	Base04 lipgloss.TerminalColor `yaml:"base04"`
	Base05 lipgloss.TerminalColor `yaml:"base05"`
	Base06 lipgloss.TerminalColor `yaml:"base06"`
	Base07 lipgloss.TerminalColor `yaml:"base07"`
	Base08 lipgloss.TerminalColor `yaml:"base08"`
	Base09 lipgloss.TerminalColor `yaml:"base09"`
	Base0A lipgloss.TerminalColor `yaml:"base0A"`
	Base0B lipgloss.TerminalColor `yaml:"base0B"`
	Base0C lipgloss.TerminalColor `yaml:"base0C"`
	Base0D lipgloss.TerminalColor `yaml:"base0D"`
	Base0E lipgloss.TerminalColor `yaml:"base0E"`
	Base0F lipgloss.TerminalColor `yaml:"base0F"`
}

// Uses the base16 rose-pine theme for Dark and rose-pine-dawn theme for Light
func DefaultColors() Colors {
	return Colors{
		Base00: lipgloss.AdaptiveColor{Light: "#faf4ed", Dark: "#191724"},
		Base01: lipgloss.AdaptiveColor{Light: "#fffaf3", Dark: "#1f1d2e"},
		Base02: lipgloss.AdaptiveColor{Light: "#f2e9de", Dark: "#26233a"},
		Base03: lipgloss.AdaptiveColor{Light: "#9893a5", Dark: "#6e6a86"},
		Base04: lipgloss.AdaptiveColor{Light: "#797593", Dark: "#908caa"},
		Base05: lipgloss.AdaptiveColor{Light: "#575279", Dark: "#e0def4"},
		Base06: lipgloss.AdaptiveColor{Light: "#575279", Dark: "#e0def4"},
		Base07: lipgloss.AdaptiveColor{Light: "#cecacd", Dark: "#524f67"},
		Base08: lipgloss.AdaptiveColor{Light: "#b4637a", Dark: "#eb6f92"},
		Base09: lipgloss.AdaptiveColor{Light: "#ea9d34", Dark: "#f6c177"},
		Base0A: lipgloss.AdaptiveColor{Light: "#d7827e", Dark: "#ebbcba"},
		Base0B: lipgloss.AdaptiveColor{Light: "#286983", Dark: "#31748f"},
		Base0C: lipgloss.AdaptiveColor{Light: "#56949f", Dark: "#9ccfd8"},
		Base0D: lipgloss.AdaptiveColor{Light: "#907aa9", Dark: "#c4a7e7"},
		Base0E: lipgloss.AdaptiveColor{Light: "#ea9d34", Dark: "#f6c177"},
		Base0F: lipgloss.AdaptiveColor{Light: "#cecacd", Dark: "#524f67"},
	}
}

func DefaultStyles(colors Colors) Styles {
	return Styles{
		BaseStyles: ListStyles{
			Title:      lipgloss.NewStyle().Background(colors.Base03).Foreground(colors.Base00).Width(20),
			Selected:   lipgloss.NewStyle().Foreground(colors.Base03),
			Unselected: lipgloss.NewStyle().Foreground(colors.Base03),
		},
		FocusedStyles: ListStyles{
			Title:      lipgloss.NewStyle().Background(colors.Base0D).Foreground(colors.Base00).Width(20),
			Selected:   lipgloss.NewStyle().Foreground(colors.Base0D),
			Unselected: lipgloss.NewStyle().Foreground(colors.Base03),
		},
		HelpStyles: HelpStyles{
			Key:  lipgloss.NewStyle().Foreground(colors.Base0D),
			Desc: lipgloss.NewStyle().Foreground(colors.Base03),
		},
	}
}
