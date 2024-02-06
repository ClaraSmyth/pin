package main

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type ListStyles struct {
	Title            lipgloss.Style
	Selected         lipgloss.Style
	Unselected       lipgloss.Style
	NoItems          lipgloss.Style
	StatusBar        lipgloss.Style
	TitleBar         lipgloss.Style
	FilterTextStyle  lipgloss.Style
	FilterPrompt     lipgloss.Style
	FilterCursor     lipgloss.Style
	FilterCursorText lipgloss.Style
}

type HelpStyles struct {
	Key  lipgloss.Style
	Desc lipgloss.Style
}

type FilePickerStyles struct {
	DisabledCursor   lipgloss.Style
	Cursor           lipgloss.Style
	Symlink          lipgloss.Style
	Directory        lipgloss.Style
	File             lipgloss.Style
	DisabledFile     lipgloss.Style
	Permission       lipgloss.Style
	Selected         lipgloss.Style
	DisabledSelected lipgloss.Style
	FileSize         lipgloss.Style
	EmptyDirectory   lipgloss.Style
}

type Styles struct {
	BaseStyles       ListStyles
	FocusedStyles    ListStyles
	HelpStyles       HelpStyles
	FilePickerStyles FilePickerStyles
	FormStyles       *huh.Theme
}

type Colors struct {
	Base00 lipgloss.TerminalColor
	Base01 lipgloss.TerminalColor
	Base02 lipgloss.TerminalColor
	Base03 lipgloss.TerminalColor
	Base04 lipgloss.TerminalColor
	Base05 lipgloss.TerminalColor
	Base06 lipgloss.TerminalColor
	Base07 lipgloss.TerminalColor
	Base08 lipgloss.TerminalColor
	Base09 lipgloss.TerminalColor
	Base0A lipgloss.TerminalColor
	Base0B lipgloss.TerminalColor
	Base0C lipgloss.TerminalColor
	Base0D lipgloss.TerminalColor
	Base0E lipgloss.TerminalColor
	Base0F lipgloss.TerminalColor
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

func FormStyles(colors Colors) *huh.Theme {
	t := huh.ThemeBase16()

	t.Focused.Base.BorderForeground(colors.Base0D)
	t.Focused.Title.Foreground(colors.Base0D)
	t.Focused.Description.Foreground(colors.Base03)
	t.Focused.ErrorIndicator.Foreground(colors.Base08)
	t.Focused.ErrorMessage.Foreground(colors.Base08)
	t.Focused.SelectSelector.Foreground(colors.Base0D)
	t.Focused.Option.Foreground(colors.Base03)
	t.Focused.MultiSelectSelector.Foreground(colors.Base0D)
	t.Focused.SelectedOption.Foreground(colors.Base0D)
	t.Focused.SelectedPrefix.Foreground(colors.Base0D)
	t.Focused.UnselectedOption.Foreground(colors.Base03)
	t.Focused.FocusedButton.Foreground(colors.Base00).Background(colors.Base0D)
	t.Focused.BlurredButton.Foreground(colors.Base03).Background(colors.Base00)

	t.Focused.TextInput.Cursor.Foreground(colors.Base04)
	t.Focused.TextInput.Placeholder.Foreground(colors.Base03)
	t.Focused.TextInput.Prompt.Foreground(colors.Base0D)

	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Title.Foreground(colors.Base03)
	t.Blurred.TextInput.Prompt.Foreground(colors.Base03)
	t.Blurred.TextInput.Text.Foreground(colors.Base03)
	t.Blurred.FocusedButton.Foreground(colors.Base00).Background(colors.Base03)
	t.Blurred.BlurredButton.Foreground(colors.Base03).Background(colors.Base00)

	return t
}

func DefaultStyles(colors Colors) Styles {
	return Styles{
		BaseStyles: ListStyles{
			Title:      lipgloss.NewStyle().Foreground(colors.Base00).Background(colors.Base02),
			Selected:   lipgloss.NewStyle().Foreground(colors.Base03),
			Unselected: lipgloss.NewStyle().Foreground(colors.Base03),
			TitleBar:   lipgloss.NewStyle().Foreground(colors.Base00).Background(colors.Base02).Width(25).Padding(0, 2).MarginBottom(1).MarginRight(2),
			NoItems:    lipgloss.NewStyle().Foreground(colors.Base03).Margin(0, 2),
			StatusBar:  lipgloss.NewStyle().Foreground(colors.Base05).Width(25).Padding(0, 2).MarginBottom(1),

			// Filter Styles
			FilterTextStyle:  lipgloss.NewStyle().Inline(true).Background(colors.Base0D).Foreground(colors.Base00),
			FilterPrompt:     lipgloss.NewStyle().Foreground(colors.Base00),
			FilterCursor:     lipgloss.NewStyle().Foreground(colors.Base03).Background(colors.Base0D),
			FilterCursorText: lipgloss.NewStyle().Foreground(colors.Base00).Background(colors.Base0D),
		},
		FocusedStyles: ListStyles{
			Title:      lipgloss.NewStyle().Foreground(colors.Base00).Background(colors.Base0D),
			Selected:   lipgloss.NewStyle().Foreground(colors.Base0D),
			Unselected: lipgloss.NewStyle().Foreground(colors.Base03),
			TitleBar:   lipgloss.NewStyle().Foreground(colors.Base00).Background(colors.Base0D).Width(25).Padding(0, 2).MarginRight(2).MaxHeight(1),
			NoItems:    lipgloss.NewStyle().Foreground(colors.Base03).Margin(0, 2),
			StatusBar:  lipgloss.NewStyle().Foreground(colors.Base05).Width(25).Padding(0, 2).Margin(1),

			// Filter Styles
			FilterTextStyle:  lipgloss.NewStyle().Inline(true).Background(colors.Base0D).Foreground(colors.Base00),
			FilterPrompt:     lipgloss.NewStyle().Foreground(colors.Base00),
			FilterCursor:     lipgloss.NewStyle().Foreground(colors.Base03).Background(colors.Base00),
			FilterCursorText: lipgloss.NewStyle().Foreground(colors.Base00).Background(colors.Base0D),
		},
		HelpStyles: HelpStyles{
			Key:  lipgloss.NewStyle().Foreground(colors.Base0D),
			Desc: lipgloss.NewStyle().Foreground(colors.Base03),
		},
		FilePickerStyles: FilePickerStyles{
			DisabledCursor:   lipgloss.NewStyle().Foreground(colors.Base02),
			Cursor:           lipgloss.NewStyle().Foreground(colors.Base04),
			Symlink:          lipgloss.NewStyle().Foreground(colors.Base0E),
			Directory:        lipgloss.NewStyle().Foreground(colors.Base0C),
			File:             lipgloss.NewStyle().Foreground(colors.Base03),
			DisabledFile:     lipgloss.NewStyle().Foreground(colors.Base02),
			Permission:       lipgloss.NewStyle().Foreground(colors.Base00),
			Selected:         lipgloss.NewStyle().Foreground(colors.Base0D),
			DisabledSelected: lipgloss.NewStyle().Foreground(colors.Base02),
			FileSize:         lipgloss.NewStyle().Foreground(colors.Base00),
			EmptyDirectory:   lipgloss.NewStyle().Foreground(colors.Base02),
		},
		FormStyles: FormStyles(colors),
	}
}
