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

type Model struct {
	help             help.Model
	apps             list.Model
	templates        list.Model
	pane             Pane
	keys             KeyMap
	form             *huh.Form
	formActive       bool
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

func (m *Model) getTemplates() tea.Cmd {
	selectedApp := m.apps.SelectedItem().(Template).Name
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

func (m *Model) setFormValues(x interface{}) {
	switch v := x.(type) {
	case App:
		formName = v.Name
		formHook = v.Hook
		formApply = false
	case Template:
		formName = v.Name
		formApply = false
	default:
		formName = ""
		formHook = ""
		formApply = false
		m.selectedFile = ""
	}
}

func (m *Model) createNewFormItem() tea.Cmd {
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

		return UpdateAppList(newApp, m.apps.Items())

	case templatePane:
		return UpdateTemplateList(m.apps.SelectedItem().(App), m.form.GetString("name"))

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

	case updateTemplatesMsg:
		return m, m.getTemplates()

	case updateAppListMsg:
		return m, m.apps.SetItems(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.updatePane()
		case "n":
			if !m.formActive {
				m.formActive = true
				m.setFormValues("reset")
				m.form = newForm(m.pane)
				return m, m.form.Init()
			}
		case "e":
			if !m.formActive {
				m.formActive = true
				currApp := m.apps.SelectedItem().(App)
				m.setFormValues(currApp)
				m.selectedFile = currApp.Path
				m.form = newForm(m.pane)
				return m, m.form.Init()
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
				return m, m.createNewFormItem()
			}

			return m, cmd
		}

		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			cmds = append(cmds, cmd)
		}

		if m.form.State == huh.StateCompleted {
			if m.pane == appPane {
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
