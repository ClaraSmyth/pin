package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// App List

type App struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Template string `yaml:"template"`
	Hook     string `yaml:"hook"`
	Backup   string `yaml:"backup"`
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
		return UpdateTemplates(selectedItem.(App).Name)
	}

	return nil
}

func (a AppDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	f, ok := item.(App)
	if !ok {
		return
	}

	fmt.Fprint(w, "  ")
	if index == m.Index() {
		fmt.Fprint(w, a.styles.Selected.Render("❯ "+f.Name))
		return
	}
	fmt.Fprint(w, a.styles.Unselected.Render("  "+f.Name))
}

// Template List

type Template struct {
	Name    string
	Path    string
	AppPath string
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

	fmt.Fprint(w, "  ")
	if index == m.Index() {
		fmt.Fprint(w, t.styles.Selected.Render("❯ "+template.Name))
		return
	}
	fmt.Fprint(w, t.styles.Unselected.Render("  "+template.Name))
}

// Theme List

type Theme struct {
	Name string
	Path string
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

	fmt.Fprint(w, "  ")
	if index == m.Index() {
		fmt.Fprint(w, t.styles.Selected.Render("❯ "+theme.Name))
		return
	}
	fmt.Fprint(w, t.styles.Unselected.Render("  "+theme.Name))
}

func newLists(styles Styles) map[Pane]*list.Model {
	appList := list.New(GetApps(), AppDelegate{styles.FocusedStyles}, 0, 0)
	appList.Title = "Apps"
	appList.Styles.Title = styles.FocusedStyles.Title
	appList.Styles.TitleBar = styles.FocusedStyles.TitleBar
	appList.Styles.NoItems = styles.FocusedStyles.NoItems
	appList.Styles.StatusBar = styles.FocusedStyles.StatusBar
	appList.FilterInput.CharLimit = 12
	appList.SetShowHelp(false)
	appList.SetShowFilter(false)

	templateList := list.New([]list.Item{}, TemplateDelegate{styles.BaseStyles}, 0, 0)
	templateList.Title = "Templates"
	templateList.Styles.Title = styles.BaseStyles.Title
	templateList.Styles.TitleBar = styles.BaseStyles.TitleBar
	templateList.Styles.NoItems = styles.BaseStyles.NoItems
	templateList.Styles.StatusBar = styles.BaseStyles.StatusBar
	templateList.FilterInput.CharLimit = 12
	templateList.SetShowHelp(false)
	templateList.SetShowFilter(false)

	if len(appList.Items()) != 0 {
		selectedApp := appList.SelectedItem().(App)
		templateList.SetItems(GetTemplates(selectedApp.Name))
	}

	themeList := list.New(GetThemes(), ThemeDelegate{styles.BaseStyles}, 0, 0)
	themeList.Title = "Themes"
	themeList.Styles.Title = styles.BaseStyles.Title
	themeList.Styles.TitleBar = styles.BaseStyles.TitleBar
	themeList.Styles.NoItems = styles.BaseStyles.NoItems
	themeList.Styles.StatusBar = styles.BaseStyles.StatusBar
	themeList.FilterInput.CharLimit = 12
	themeList.SetShowHelp(false)
	themeList.SetShowFilter(false)
	themeList.SetSpinner(spinner.MiniDot)

	listMap := make(map[Pane]*list.Model)

	listMap[appPane] = &appList
	listMap[templatePane] = &templateList
	listMap[themePane] = &themeList

	return listMap
}
