package main

import (
	"os"

	"github.com/ClaraSmyth/pin/filepicker"
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
	help             help.Model
	apps             list.Model
	templates        list.Model
	pane             Pane
	keys             KeyMap
	forms            map[Pane]*huh.Form
	formActive       bool
	filepicker       filepicker.Model
	filepickerActive bool
	selectedFile     string
	height           int
}

type updateTemplatesMsg App

type updateAppListMsg []list.Item

type toggleFormMsg bool

type resetFormsMsg bool

func newModel() *Model {
	appList, templateList := newLists()
	forms := newForms()

	filepicker := filepicker.New()
	filepicker.CurrentDirectory, _ = os.UserHomeDir()
	filepicker.ShowHidden = true

	return &Model{
		keys:             DefaultKeyMap,
		help:             help.New(),
		pane:             appPane,
		apps:             appList,
		templates:        templateList,
		forms:            forms,
		formActive:       false,
		filepicker:       filepicker,
		filepickerActive: false,
		selectedFile:     "",
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

func (m *Model) resetFormValues() tea.Cmd {
	formName = ""
	formHook = ""
	formApply = false
	m.selectedFile = ""
	m.formActive = false
	m.filepickerActive = false
	m.forms = newForms()

	m.filepicker = filepicker.New()
	m.filepicker.CurrentDirectory, _ = os.UserHomeDir()
	m.filepicker.ShowHidden = true

	return nil
}

func (m *Model) createNewFormItem() tea.Cmd {
	switch m.pane {
	case appPane:

		newApp := App{
			Name:     m.forms[m.pane].GetString("name"),
			Path:     m.selectedFile,
			Template: "Default",
			Hook:     "Default",
			Backup:   "Default",
			Active:   false,
		}

		return UpdateAppList(newApp, m.apps.Items())

	case templatePane:
		return UpdateTemplateList(m.apps.SelectedItem().(App), m.forms[m.pane].GetString("name"))

	default:
		return nil
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.apps.SetSize(msg.Width, msg.Height-2)
		m.templates.SetSize(msg.Width, msg.Height-2)
		m.height = msg.Height
		return m, nil

	case resetFormsMsg:
		return m, m.resetFormValues()

	case updateTemplatesMsg:
		return m, tea.Batch(m.getTemplates(), func() tea.Msg {
			return resetFormsMsg(true)
		})

	case updateAppListMsg:
		return m, tea.Batch(m.apps.SetItems(msg), func() tea.Msg {
			return resetFormsMsg(true)
		})

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.updatePane()
		case "n":
			if !m.formActive {
				m.formActive = true
				return m, m.forms[m.pane].Init()
			}
		}
	}

	if m.formActive {

		if m.filepickerActive {
			m.filepicker, cmd = m.filepicker.Update(msg)

			if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
				m.selectedFile = path
				return m, m.createNewFormItem()
			}

			return m, cmd
		}

		form, cmd := m.forms[m.pane].Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.forms[m.pane] = f
			cmds = append(cmds, cmd)
		}

		if m.forms[m.pane].State == huh.StateCompleted {

			if m.pane == appPane && m.selectedFile == "" {
				if !m.filepickerActive {
					m.filepickerActive = true
					m.filepicker.Height = m.height - 6
					return m, m.filepicker.Init()
				}

			} else {
				cmds = append(cmds, m.createNewFormItem())
			}
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

	filepickerTitle := lipgloss.NewStyle().SetString("Select Config File:").Width(20).Background(lipgloss.Color("#000000")).Render()

	appView := m.apps.View()
	if m.formActive && m.pane == appPane {
		appView = m.forms[m.pane].View()
		if m.filepickerActive {
			appView = lipgloss.JoinVertical(lipgloss.Top, filepickerTitle, "", lipgloss.NewStyle().Width(20).Render(m.filepicker.View()))
		}
	}

	templatesView := m.templates.View()
	if m.formActive && m.pane == templatePane {
		templatesView = m.forms[m.pane].View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left, appView, templatesView),
		"",
		m.help.View(m.keys),
	)
}
