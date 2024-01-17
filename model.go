package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	appState state = iota
	templateState
)

type Model struct {
	help      help.Model
	apps      list.Model
	templates list.Model
	state     state
	keys      KeyMap
}

type updateTemplatesMsg App

func newModel() *Model {
	appList := list.New(GetAppListItems(), AppDelegate{styles}, 0, 0)
	templateList := list.New([]list.Item{}, TemplateDelegate{styles}, 0, 0)

	appList.Title = "Apps"
	appList.Styles.NoItems = lipgloss.NewStyle().Margin(0, 2)
	appList.Styles.Title = styles.Title
	appList.SetShowHelp(false)
	appList.SetShowFilter(false)

	selectedApp := appList.SelectedItem().FilterValue()

	templateList.Title = selectedApp + " Templates"
	templateList.Styles.NoItems = lipgloss.NewStyle().Margin(0, 2)
	templateList.Styles.Title = styles.Title
	templateList.SetShowHelp(false)
	templateList.SetShowFilter(false)
	templateList.SetItems(GetTemplateListItems(selectedApp))

	return &Model{
		keys:      DefaultKeyMap,
		help:      help.New(),
		state:     appState,
		apps:      appList,
		templates: templateList,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) getTemplates() tea.Cmd {
	selectedApp := m.apps.SelectedItem().FilterValue()
	m.templates.Title = selectedApp + " Templates"
	return m.templates.SetItems(GetTemplateListItems(selectedApp))
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.apps.SetSize(msg.Width, msg.Height-1)
		m.templates.SetSize(msg.Width, msg.Height-1)
		return m, nil

	case updateTemplatesMsg:
		return m, m.getTemplates()

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.state == appState {
				m.state = templateState
				return m, nil
			} else {
				m.state = appState
				return m, nil
			}
		}
	}

	if m.state == templateState {
		m.templates, cmd = m.templates.Update(msg)
		return m, cmd
	}

	m.apps, cmd = m.apps.Update(msg)

	return m, cmd
}

func (m *Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left, m.apps.View(), m.templates.View()),
		m.help.View(m.keys),
	)
}
