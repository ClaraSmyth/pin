package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit        key.Binding
	Back        key.Binding
	NextPane    key.Binding
	PrevPane    key.Binding
	Select      key.Binding
	New         key.Binding
	Edit        key.Binding
	Open        key.Binding
	Delete      key.Binding
	Search      key.Binding
	FetchThemes key.Binding
	ToggleHelp  key.Binding
}

var DefaultKeyMap = KeyMap{
	Quit:        key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "exit")),
	Back:        key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	NextPane:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "navigate")),
	PrevPane:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "navigate")),
	Select:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	New:         key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
	Edit:        key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Open:        key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open")),
	Delete:      key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "delete")),
	Search:      key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	FetchThemes: key.NewBinding(key.WithKeys("P"), key.WithHelp("P", "fetch themes"), key.WithDisabled()),
	ToggleHelp:  key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
		k.NextPane,
		k.Select,
		k.New,
		k.Delete,
		k.Open,
		k.Edit,
		k.FetchThemes,
		k.ToggleHelp,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Open, k.Edit},
		{k.New, k.Delete},
		{k.Quit, k.Back},
		{k.NextPane, k.PrevPane},
		{k.ToggleHelp, k.FetchThemes},
	}
}
