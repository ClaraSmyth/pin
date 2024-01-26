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
	themePane
)

type FormAction int

const (
	formActionCreate FormAction = iota
	formActionEdit
	formActionDelete
)

type Model struct {
	help             help.Model
	lists            map[Pane]*list.Model
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

type updateThemeListMsg []list.Item

type updateTemplateListMsg []list.Item

type updateAppListMsg struct {
	appListItems      []list.Item
	templateListItems []list.Item
}

func newModel() *Model {
	listMap := newLists()

	return &Model{
		lists:            listMap,
		keys:             DefaultKeyMap,
		help:             help.New(),
		pane:             appPane,
		formActive:       false,
		filepickerActive: false,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) updatePane() {
	switch m.pane {
	case appPane:
		m.pane = templatePane
	case templatePane:
		m.pane = themePane
	case themePane:
		m.pane = appPane
	}
}

func (m *Model) selectItem() tea.Cmd {
	switch selectedItem := m.lists[m.pane].SelectedItem().(type) {
	case Template:
		app := m.lists[appPane].SelectedItem().(App)
		newApp := app
		newApp.Template = selectedItem.Path
		return EditApp(newApp, app, m.lists[appPane].Items())
	case Theme:
		return ApplyTheme(selectedItem)
	default:
		return nil
	}
}

func (m *Model) openFileEditor() tea.Cmd {
	var path string

	switch selectedItem := m.lists[m.pane].SelectedItem().(type) {
	case App:
		path = selectedItem.Path
	case Template:
		path = selectedItem.Path
	case Theme:
		path = selectedItem.Path
	}

	return tea.ExecProcess(editorCmd(path), func(err error) tea.Msg {
		return nil
	})

}

func (m *Model) triggerFilepicker() tea.Cmd {
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
	return m.filepicker.Init()
}

func (m *Model) triggerForm(formAction FormAction) tea.Cmd {
	// Set form to active and set formAction
	m.formActive = true
	m.formAction = formAction

	// Reset Form Values
	formEdit = false
	formName = ""
	formHook = ""
	formFilepicker = false
	formRewrite = false
	formApply = false
	m.selectedFile = ""

	selectedItem := m.lists[m.pane].SelectedItem()
	if selectedItem == nil && formAction != formActionCreate {
		m.formActive = false
		m.formAction = formActionCreate
		return nil
	}

	switch formAction {
	case formActionCreate:
		m.form = newForm(m.pane, m.lists[m.pane].Items())

	case formActionEdit:
		switch item := selectedItem.(type) {
		case App:
			formEdit = true
			formName = item.Name
			formHook = item.Hook
			formRewrite = item.Rewrite
			m.selectedFile = item.Path
		case Template:
			formEdit = true
			formName = item.Name
		case Theme:
			m.formActive = false
			return nil
		}
		m.form = newForm(m.pane, m.lists[m.pane].Items())

	case formActionDelete:
		m.form = deleteForm()
	}

	return m.form.Init()
}

func (m *Model) handleFormSubmit() tea.Cmd {

	m.formActive = false
	m.filepickerActive = false

	if m.form.State != huh.StateCompleted || formApply == false {
		return nil
	}

	switch m.pane {
	case appPane:
		newApp := App{
			Name:    m.form.GetString("name"),
			Path:    m.selectedFile,
			Rewrite: !m.form.GetBool("rewrite"),
			Active:  false,
		}

		switch m.formAction {
		case formActionCreate:
			return CreateApp(newApp, m.lists[appPane].Items())

		case formActionEdit:
			return EditApp(newApp, m.lists[appPane].SelectedItem().(App), m.lists[appPane].Items())

		case formActionDelete:
			return DeleteApp(m.lists[appPane].SelectedItem().(App), m.lists[appPane].Items())
		}

	case templatePane:
		switch m.formAction {
		case formActionCreate:
			return CreateTemplate(m.lists[appPane].SelectedItem().(App), m.form.GetString("name"))

		case formActionEdit:
			return EditTemplate(m.lists[appPane].SelectedItem().(App), m.lists[templatePane].SelectedItem().(Template).Name, m.form.GetString("name"))

		case formActionDelete:
			return DeleteTemplate(m.lists[appPane].SelectedItem().(App), m.lists[templatePane].SelectedItem().(Template).Name)
		}

	case themePane:
		switch m.formAction {
		case formActionCreate:
			return CreateTheme(m.form.GetString("name"), m.lists[themePane].Items())

		case formActionDelete:
			return DeleteTheme(m.lists[themePane].SelectedItem().(Theme))
		}
	}

	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.lists[appPane].SetSize(msg.Width, msg.Height-2)
		m.lists[templatePane].SetSize(msg.Width, msg.Height-2)
		m.lists[themePane].SetSize(msg.Width, msg.Height-2)
		m.height = msg.Height
		return m, nil

	case updateThemeListMsg:
		m.lists[themePane].StopSpinner()
		return m, m.lists[themePane].SetItems(msg)

	case updateTemplateListMsg:
		return m, m.lists[templatePane].SetItems(msg)

	case updateAppListMsg:
		return m, tea.Batch(m.lists[appPane].SetItems(msg.appListItems), m.lists[templatePane].SetItems(msg.templateListItems))

	case tea.KeyMsg:

		if m.formActive {
			switch msg.String() {
			case "esc":
				return m, m.handleFormSubmit()
			}
		}

		if !m.formActive {
			switch msg.String() {
			case "tab":
				m.updatePane()
			case "n":
				return m, m.triggerForm(formActionCreate)
			case "e":
				return m, m.triggerForm(formActionEdit)
			case "x":
				return m, m.triggerForm(formActionDelete)
			case "o":
				return m, m.openFileEditor()
			case "enter":
				return m, m.selectItem()
			case "P":
				if m.pane == themePane {
					cmd := m.lists[themePane].StartSpinner()
					return m, tea.Batch(GitCloneSchemes(), cmd)
				}
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
				formFilepicker = false
				formApply = true
				return m, m.handleFormSubmit()
			}

			return m, cmd
		}

		var cmds []tea.Cmd

		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			cmds = append(cmds, cmd)
		}

		if m.form.State == huh.StateAborted {
			return m, m.handleFormSubmit()
		}

		if m.form.State == huh.StateCompleted {
			if m.pane == appPane && m.formAction != formActionDelete && m.form.GetBool("filepicker") == true {
				return m, m.triggerFilepicker()
			} else {
				return m, m.handleFormSubmit()
			}
		}

		return m, tea.Batch(cmds...)
	}

	*m.lists[m.pane], cmd = m.lists[m.pane].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	appView := m.lists[appPane].View()
	if m.formActive && m.pane == appPane {
		appView = m.form.View()
		if m.filepickerActive {
			filepickerTitle := lipgloss.NewStyle().SetString("Select Config File:").Width(20).Background(lipgloss.Color("#000000")).Render()
			appView = lipgloss.JoinVertical(lipgloss.Top, filepickerTitle, "", lipgloss.NewStyle().Width(20).Render(m.filepicker.View()))
		}
	}

	templatesView := m.lists[templatePane].View()
	if m.formActive && m.pane == templatePane {
		templatesView = m.form.View()
	}

	themeView := m.lists[themePane].View()
	if m.formActive && m.pane == themePane {
		themeView = m.form.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(m.height-1).Render(lipgloss.JoinHorizontal(lipgloss.Left, appView, templatesView, themeView)),
		m.help.View(m.keys),
	)
}
