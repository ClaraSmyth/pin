package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// App List

type App struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Template string `yaml:"template"`
	Hook     string `yaml:"hook"`
	Active   bool   `yaml:"active"`
	Rewrite  bool   `yaml:"rewrite"`
}

func (a App) FilterValue() string { return a.Name }

type AppDelegate struct{ styles ListStyles }

func (a AppDelegate) Height() int  { return 1 }
func (a AppDelegate) Spacing() int { return 0 }

func (a AppDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	selectedItem := m.SelectedItem()
	if selectedItem != nil {
		return UpdateTemplates(selectedItem.(App))
	}

	return nil
}

func (a AppDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	app, ok := item.(App)
	if !ok {
		return
	}

	statusDot := "● "
	if !app.Active {
		statusDot = "○ "
	}

	if app.Path == "" || app.Template == "" {
		statusDot = "✗ "
	}

	if index == m.Index() {
		fmt.Fprint(w, a.styles.Selected.Render("❯ "+statusDot+app.Name))
		return
	}
	fmt.Fprint(w, a.styles.Unselected.Render("  "+statusDot+app.Name))
}

// Template List

type Template struct {
	Name    string
	Path    string
	AppPath string
	Active  bool
}

func (t Template) FilterValue() string { return t.Name }

type TemplateDelegate struct{ styles ListStyles }

func (t TemplateDelegate) Height() int                               { return 1 }
func (t TemplateDelegate) Spacing() int                              { return 0 }
func (t TemplateDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (t TemplateDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	template, ok := item.(Template)
	if !ok {
		return
	}

	statusDot := "● "
	if !template.Active {
		statusDot = "○ "
	}

	if index == m.Index() {
		fmt.Fprint(w, t.styles.Selected.Render("❯ "+statusDot+template.Name))
		return
	}
	fmt.Fprint(w, t.styles.Unselected.Render("  "+statusDot+template.Name))
}

// Theme List

type Theme struct {
	Name   string
	Path   string
	Active bool
	Err    bool
}

func (t Theme) FilterValue() string { return t.Name }

type ThemeDelegate struct{ styles ListStyles }

func (t ThemeDelegate) Height() int                               { return 1 }
func (t ThemeDelegate) Spacing() int                              { return 0 }
func (t ThemeDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (t ThemeDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	theme, ok := item.(Theme)
	if !ok {
		return
	}

	statusDot := "● "
	if !theme.Active {
		statusDot = "○ "
	}

	if theme.Err {
		statusDot = "✗ "
	}

	if index == m.Index() {
		fmt.Fprint(w, t.styles.Selected.Render("❯ "+statusDot+theme.Name))
		return
	}
	fmt.Fprint(w, t.styles.Unselected.Render("  "+statusDot+theme.Name))
}

func newLists(styles Styles) map[Pane]*list.Model {
	appList := list.New(GetApps(), AppDelegate{styles.FocusedStyles}, 0, 0)
	appList.Title = "Apps"
	appList.FilterInput.Prompt = "Find: "
	appList.FilterInput.CharLimit = 14
	appList.SetShowHelp(false)
	UpdateListStyles(&appList, styles.FocusedStyles)

	templateList := list.New([]list.Item{}, TemplateDelegate{styles.BaseStyles}, 0, 0)
	templateList.Title = "Templates"
	templateList.FilterInput.Prompt = "Find: "
	templateList.FilterInput.CharLimit = 14
	templateList.SetShowHelp(false)
	UpdateListStyles(&templateList, styles.BaseStyles)

	if len(appList.Items()) != 0 {
		selectedApp := appList.SelectedItem().(App)
		templateList.SetItems(GetTemplates(selectedApp))
	}

	themeList := list.New(GetThemes(), ThemeDelegate{styles.BaseStyles}, 0, 0)
	themeList.Title = "Themes"
	themeList.FilterInput.Prompt = "Find: "
	themeList.FilterInput.CharLimit = 14
	themeList.SetShowHelp(false)
	themeList.SetSpinner(spinner.MiniDot)
	UpdateListStyles(&themeList, styles.BaseStyles)
	themeList.Styles.StatusBar = themeList.Styles.StatusBar.Copy().UnsetWidth()
	themeList.Styles.Spinner.Foreground(lipgloss.ANSIColor(0))

	listMap := make(map[Pane]*list.Model)

	listMap[appPane] = &appList
	listMap[templatePane] = &templateList
	listMap[themePane] = &themeList

	return listMap
}

func UpdateListStyles(list *list.Model, styles ListStyles) {
	list.Styles.Title = styles.Title
	list.Styles.TitleBar = styles.TitleBar
	list.Styles.NoItems = styles.NoItems
	list.Styles.StatusBar = styles.StatusBar
	list.Styles.StatusEmpty = styles.StatusEmpty
	list.Styles.StatusBarFilterCount = styles.StatusBarFilterCount
	list.Styles.StatusBarActiveFilter = styles.StatusBarActiveFilter

	list.FilterInput.TextStyle = styles.FilterTextStyle
	list.FilterInput.PromptStyle = styles.FilterPrompt
	list.FilterInput.Cursor.Style = styles.FilterCursor
	list.FilterInput.Cursor.TextStyle = styles.FilterCursorText
}
