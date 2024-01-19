package main

import (
	"os"
	"strings"

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

type FormAction int

const (
	formActionCreate FormAction = iota
	formActionEdit
	formActionDelete
)

type Model struct {
	help             help.Model
	apps             list.Model
	templates        list.Model
	pane             Pane
	keys             KeyMap
	form             *huh.Form
	formActive       bool
	formAction       FormAction
	filepicker       filepicker.Model
	filepickerActive bool
	selectedFile     string
	height           int
}

type updateTemplatesMsg App

type updateAppListMsg []list.Item

func newModel() *Model {
	appList, templateList := newLists()

	return &Model{
		keys:             DefaultKeyMap,
		help:             help.New(),
		pane:             appPane,
		apps:             appList,
		templates:        templateList,
		formActive:       false,
		filepickerActive: false,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) updateTemplates() tea.Cmd {
	selectedApp := m.apps.SelectedItem().(App).Name
	return m.templates.SetItems(GetTemplates(selectedApp))
}

func (m *Model) updatePane() {
	switch m.pane {
	case appPane:
		m.pane = templatePane
	case templatePane:
		m.pane = appPane
	}

}

func (m *Model) triggerForm(formAction FormAction) (tea.Model, tea.Cmd) {

	// Set form to active
	m.formActive = true

	// Reset Form Values
	formName = ""
	formHook = ""
	formApply = false
	m.selectedFile = ""

	switch formAction {
	case formActionCreate:
		m.formAction = formAction
		m.form = newForm(m.pane)
		return m, m.form.Init()

	case formActionDelete:
		m.formAction = formAction

		switch m.pane {
		case appPane:
			if len(m.apps.Items()) == 0 {
				m.formActive = false
				m.formAction = formActionCreate
				return m, nil
			}

		case templatePane:
			if len(m.templates.Items()) == 0 {
				m.formActive = false
				m.formAction = formActionCreate
				return m, nil
			}
		}

		m.form = deleteForm()
		return m, m.form.Init()

	case formActionEdit:
		m.formAction = formAction

		switch m.pane {
		case appPane:
			if len(m.apps.Items()) == 0 {
				m.formActive = false
				return m, nil
			}
			currApp := m.apps.SelectedItem().(App)
			formName = currApp.Name
			formHook = currApp.Hook
			m.selectedFile = currApp.Path

		case templatePane:
			if len(m.templates.Items()) == 0 {
				m.formActive = false
				return m, nil
			}
			currTemplate := m.templates.SelectedItem().(Template)
			formName = currTemplate.Name
		}

		m.form = newForm(m.pane)
		return m, m.form.Init()

	default:
		return m, nil
	}
}

func (m *Model) handleFormSubmit() tea.Cmd {
	switch m.pane {
	case appPane:

		newApp := App{
			Name:     m.form.GetString("name"),
			Path:     m.selectedFile,
			Template: "Default",
			Hook:     "Default",
			Backup:   "Default",
			Active:   false,
		}

		switch m.formAction {
		case formActionCreate:
			return CreateApp(newApp, m.apps.Items())

		case formActionEdit:
			return EditApp(newApp, m.apps.SelectedItem().(App), m.apps.Items())

		case formActionDelete:
			return DeleteApp(m.apps.SelectedItem().(App), m.apps.Items())
		}

	case templatePane:
		switch m.formAction {
		case formActionCreate:
			return CreateTemplate(m.apps.SelectedItem().(App), m.form.GetString("name"))

		case formActionEdit:
			return EditTemplate(m.apps.SelectedItem().(App), m.templates.SelectedItem().(Template).Filename, m.form.GetString("name"))

		case formActionDelete:
			return DeleteTemplate(m.apps.SelectedItem().(App), m.templates.SelectedItem().(Template).Filename)
		}
	}

	return nil
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

	case updateTemplatesMsg:
		return m, m.updateTemplates()

	case updateAppListMsg:
		return m, m.apps.SetItems(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if !m.formActive {
				m.updatePane()
			}
		case "n":
			if !m.formActive {
				return m.triggerForm(formActionCreate)
			}
		case "e":
			if !m.formActive {
				return m.triggerForm(formActionEdit)
			}
		case "x":
			if !m.formActive {
				return m.triggerForm(formActionDelete)
			}
		}
	}

	if m.formActive {
		if m.filepickerActive {
			m.filepicker, cmd = m.filepicker.Update(msg)

			if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
				m.selectedFile = path
				m.formActive = false
				m.filepickerActive = false
				return m, m.handleFormSubmit()
			}

			return m, cmd
		}

		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			cmds = append(cmds, cmd)
		}

		if m.form.State == huh.StateCompleted {
			if m.pane == appPane && m.formAction != formActionDelete {
				if !m.filepickerActive {
					m.filepickerActive = true
					m.filepicker = filepicker.New()
					m.filepicker.CurrentDirectory, _ = os.UserHomeDir()
					if m.selectedFile != "" {
						dir := strings.TrimRightFunc(m.selectedFile, func(r rune) bool {
							return !strings.ContainsRune("/", r)
						})
						m.filepicker.CurrentDirectory = dir
					}
					m.filepicker.ShowHidden = true
					m.filepicker.Height = m.height - 6
					return m, m.filepicker.Init()
				}

			} else {
				m.formActive = false
				cmds = append(cmds, m.handleFormSubmit())
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

	appView := m.apps.View()
	if m.formActive && m.pane == appPane {
		appView = m.form.View()
		if m.filepickerActive {
			filepickerTitle := lipgloss.NewStyle().SetString("Select Config File:").Width(20).Background(lipgloss.Color("#000000")).Render()
			appView = lipgloss.JoinVertical(lipgloss.Top, filepickerTitle, "", lipgloss.NewStyle().Width(20).Render(m.filepicker.View()))
		}
	}

	templatesView := m.templates.View()
	if m.formActive && m.pane == templatePane {
		templatesView = m.form.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left, appView, templatesView),
		"",
		m.help.View(m.keys),
	)
}
