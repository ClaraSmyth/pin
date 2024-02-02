package main

import (
	"os"
	"strings"

	"github.com/ClaraSmyth/pin/filepicker"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	styles           Styles
}

type updateThemeListMsg []list.Item

type updateTemplateListMsg []list.Item

type updateAppListMsg struct {
	appListItems      []list.Item
	templateListItems []list.Item
}

type updateStylesMsg Styles

func newModel() *Model {
	colors := GetActiveColors()
	styles := DefaultStyles(colors)

	listMap := newLists(styles)

	help := help.New()
	help.Styles.ShortKey = styles.HelpStyles.Key
	help.Styles.FullKey = styles.HelpStyles.Key
	help.Styles.ShortDesc = styles.HelpStyles.Desc
	help.Styles.FullDesc = styles.HelpStyles.Desc

	return &Model{
		lists:            listMap,
		keys:             DefaultKeyMap,
		help:             help,
		pane:             appPane,
		formActive:       false,
		filepickerActive: false,
		styles:           styles,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) updateStyles() tea.Cmd {
	for pane, list := range m.lists {
		if pane == m.pane {
			list.Styles = UpdateListStyles(list.Styles, m.styles, true)
		} else {
			list.Styles = UpdateListStyles(list.Styles, m.styles, false)
		}
	}

	switch m.pane {
	case appPane:
		m.lists[appPane].SetDelegate(AppDelegate{m.styles.FocusedStyles})
		m.lists[templatePane].SetDelegate(TemplateDelegate{m.styles.BaseStyles})
		m.lists[themePane].SetDelegate(ThemeDelegate{m.styles.BaseStyles})

	case templatePane:
		m.lists[appPane].SetDelegate(AppDelegate{m.styles.BaseStyles})
		m.lists[templatePane].SetDelegate(TemplateDelegate{m.styles.FocusedStyles})
		m.lists[themePane].SetDelegate(ThemeDelegate{m.styles.BaseStyles})
	case themePane:
		m.lists[appPane].SetDelegate(AppDelegate{m.styles.BaseStyles})
		m.lists[templatePane].SetDelegate(TemplateDelegate{m.styles.BaseStyles})
		m.lists[themePane].SetDelegate(ThemeDelegate{m.styles.FocusedStyles})
	}

	m.lists[themePane].Styles.StatusBar.UnsetWidth()

	m.help.Styles.ShortKey = m.styles.HelpStyles.Key
	m.help.Styles.FullKey = m.styles.HelpStyles.Key
	m.help.Styles.ShortDesc = m.styles.HelpStyles.Desc
	m.help.Styles.FullDesc = m.styles.HelpStyles.Desc
	return nil
}

func (m *Model) updateKeys() tea.Cmd {
	switch m.pane {
	case themePane:
		m.keys.FetchThemes.SetEnabled(true)
		m.keys.Edit.SetEnabled(false)
		return nil
	default:
		m.keys.FetchThemes.SetEnabled(false)
		m.keys.Edit.SetEnabled(true)
		return nil
	}
}

func (m *Model) selectItem() tea.Cmd {
	switch selectedItem := m.lists[m.pane].SelectedItem().(type) {
	case App:
		if selectedItem.Path == "" {
			return nil
		}
		app := m.lists[appPane].SelectedItem().(App)
		newApp := app
		newApp.Active = !app.Active
		return EditApp(newApp, app, m.lists[appPane].Items())

	case Template:
		app := m.lists[appPane].SelectedItem().(App)
		newApp := app
		newApp.Template = selectedItem.Path
		return EditApp(newApp, app, m.lists[appPane].Items())

	case Theme:
		items := m.lists[themePane].Items()
		for i, item := range items {
			newItem := item.(Theme)
			newItem.Active = false
			if selectedItem.Name == newItem.Name {
				newItem.Active = true
			}

			items[i] = newItem
		}

		return tea.Batch(ApplyThemeCmd(selectedItem), m.lists[themePane].SetItems(items))
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
	m.filepicker.Styles = filepicker.Styles(m.styles.FilePickerStyles)
	m.filepicker.Cursor = "❯"
	m.filepicker.CurrentDirectory, _ = os.UserHomeDir()
	if m.selectedFile != "" {
		dir := strings.TrimRightFunc(m.selectedFile, func(r rune) bool {
			return !strings.ContainsRune("/", r)
		})
		m.filepicker.CurrentDirectory = dir
	}
	m.filepicker.ShowHidden = true
	m.filepicker.Height = m.height - 4
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
		if m.pane == templatePane {
			selectedApp := m.lists[appPane].SelectedItem()
			if selectedApp == nil {
				m.formActive = false
				m.formAction = formActionCreate
				return nil
			}
		}
		m.form = newForm(m.pane, m.lists[m.pane].Items(), m.styles.FormStyles)

	case formActionEdit:
		switch item := selectedItem.(type) {
		case App:
			formEdit = true
			formName = item.Name
			formHook = item.Hook
			formRewrite = !item.Rewrite
			m.selectedFile = item.Path
		case Template:
			formEdit = true
			formName = item.Name
		case Theme:
			m.formActive = false
			return nil
		}
		m.form = newForm(m.pane, m.lists[m.pane].Items(), m.styles.FormStyles)

	case formActionDelete:
		m.form = deleteForm(m.styles.FormStyles)
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
			Hook:    m.form.GetString("hook"),
			Rewrite: !m.form.GetBool("rewrite"),
			Active:  false,
		}

		switch m.formAction {
		case formActionCreate:
			return CreateApp(newApp, m.lists[appPane].Items())

		case formActionEdit:
			return EditApp(newApp, m.lists[appPane].SelectedItem().(App), m.lists[appPane].Items())

		case formActionDelete:
			return DeleteApp(m.lists[appPane].SelectedItem().(App), m.lists[appPane].Index(), m.lists[appPane].Items())
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

	case updateStylesMsg:
		m.styles = Styles(msg)
		return m, m.updateStyles()

	case tea.KeyMsg:
		if m.formActive {
			switch {
			case key.Matches(msg, m.keys.Back):
				return m, m.handleFormSubmit()
			}
		}

		if !m.formActive && !m.lists[m.pane].SettingFilter() {
			switch {
			case key.Matches(msg, m.keys.NextPane):
				m.pane = (m.pane + 1) % 3
				return m, tea.Batch(m.updateKeys(), m.updateStyles())

			case key.Matches(msg, m.keys.PrevPane):
				m.pane = (m.pane + 2) % 3
				return m, tea.Batch(m.updateKeys(), m.updateStyles())

			case key.Matches(msg, m.keys.New):
				return m, m.triggerForm(formActionCreate)

			case key.Matches(msg, m.keys.Edit):
				return m, m.triggerForm(formActionEdit)

			case key.Matches(msg, m.keys.Delete):
				return m, m.triggerForm(formActionDelete)

			case key.Matches(msg, m.keys.Open):
				return m, m.openFileEditor()

			case key.Matches(msg, m.keys.Select):
				return m, m.selectItem()

			case key.Matches(msg, m.keys.ToggleHelp):
				m.help.ShowAll = !m.help.ShowAll
				m.lists[appPane].SetHeight(m.height - lipgloss.Height(m.help.View(m.keys)))
				m.lists[templatePane].SetHeight(m.height - lipgloss.Height(m.help.View(m.keys)))
				m.lists[themePane].SetHeight(m.height - lipgloss.Height(m.help.View(m.keys)))

			case key.Matches(msg, m.keys.FetchThemes):
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

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {

	currentList := m.lists[m.pane]
	if currentList.SettingFilter() {
		currentList.Title = "Find: " + currentList.FilterValue()
	} else {
		switch m.pane {
		case appPane:
			currentList.Title = "Apps"
		case templatePane:
			currentList.Title = "Templates"
		case themePane:
			currentList.Title = "Themes"
		}
	}

	appView := m.lists[appPane].View()
	templatesView := m.lists[templatePane].View()
	themeView := m.lists[themePane].View()

	if m.formActive {
		formTitleText := ""

		switch m.formAction {
		case formActionCreate:
			formTitleText = "New "
		case formActionEdit:
			formTitleText = "Edit "
		case formActionDelete:
			formTitleText = "Delete "
		}

		formTitleStyles := m.styles.FocusedStyles.Title

		switch m.pane {
		case appPane:
			appView = lipgloss.JoinVertical(lipgloss.Top, formTitleStyles.Render(formTitleText+"App"), "", m.form.View())

			if m.filepickerActive {
				formTitleText = "Select file:"
				appView = lipgloss.JoinVertical(lipgloss.Top, formTitleStyles.Render(formTitleText), "", lipgloss.NewStyle().Width(25).Render(m.filepicker.View()))
			}

		case templatePane:
			templatesView = lipgloss.JoinVertical(lipgloss.Top, formTitleStyles.Render(formTitleText+"Template"), "", m.form.View())
		case themePane:
			themeView = lipgloss.JoinVertical(lipgloss.Top, formTitleStyles.Render(formTitleText+"Theme"), "", m.form.View())
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.
			NewStyle().
			Height(m.height-lipgloss.Height(m.help.View(m.keys))).
			Render(lipgloss.JoinHorizontal(lipgloss.Left, appView, templatesView, themeView)),
		m.help.View(m.keys),
	)
}
