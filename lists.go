package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func newLists() map[Pane]*list.Model {
	appList := list.New(GetApps(), AppDelegate{styles}, 0, 0)
	templateList := list.New([]list.Item{}, TemplateDelegate{styles}, 0, 0)

	appList.Title = "Apps"
	appList.Styles.NoItems = lipgloss.NewStyle().Margin(0, 2)
	appList.Styles.Title = styles.Title
	appList.SetShowHelp(false)
	appList.SetShowFilter(false)

	templateList.Title = "Templates"
	templateList.Styles.NoItems = lipgloss.NewStyle().Margin(0, 2)
	templateList.Styles.Title = styles.Title
	templateList.SetShowHelp(false)
	templateList.SetShowFilter(false)

	if len(appList.Items()) != 0 {
		selectedApp := appList.SelectedItem().(App)
		templateList.SetItems(GetTemplates(selectedApp.Name))
	}

	listMap := make(map[Pane]*list.Model)

	listMap[appPane] = &appList
	listMap[templatePane] = &templateList

	return listMap
}
