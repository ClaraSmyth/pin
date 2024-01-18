package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type Pane int

const (
	appPane Pane = iota
	templatePane
)

type Model struct {
	help       help.Model
	apps       list.Model
	templates  list.Model
	pane       Pane
	keys       KeyMap
	forms      map[Pane]*huh.Form
	formActive bool
}

type updateTemplatesMsg App

type toggleFormMsg bool

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

	forms := newForms()

	return &Model{
		keys:       DefaultKeyMap,
		help:       help.New(),
		pane:       appPane,
		apps:       appList,
		templates:  templateList,
		forms:      forms,
		formActive: false,
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

func (m *Model) updatePane() {
	if !m.formActive {
		switch m.pane {
		case appPane:
			m.pane = templatePane
		case templatePane:
			m.pane = appPane
		}
	}
}

func (m *Model) triggerForm() tea.Cmd {
	if !m.formActive {
		m.formActive = true
		return m.forms[m.pane].Init()
	}
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

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
			m.updatePane()
		case "n":
			return m, m.triggerForm()
		}
	}

	if m.formActive {
		form, cmd := m.forms[m.pane].Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.forms[m.pane] = f
			cmds = append(cmds, cmd)
		}

		if m.forms[m.pane].State == huh.StateCompleted {
			m.formActive = false
			resetFormValues()
			m.forms = newForms()
		}

		return m, tea.Batch(cmds...)
	}

	if m.pane == appPane {
		m.apps, cmd = m.apps.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.pane == templatePane {
		m.templates, cmd = m.templates.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {

	appView := m.apps.View()
	if m.formActive && m.pane == appPane {
		appView = m.forms[m.pane].View()
	}

	templatesView := m.templates.View()
	if m.formActive && m.pane == templatePane {
		templatesView = m.forms[m.pane].View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left, appView, templatesView),
		m.help.View(m.keys),
	)
}
