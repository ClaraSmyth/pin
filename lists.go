package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

// App List

type App struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Template string `yaml:"template"`
	Hook     string `yaml:"hook"`
	Backup   string `yaml:"backup"`
	Active   bool   `yaml:"active"`
}

var defaultApp = App{
	Name:     "userdefined",
	Path:     "userdefined",
	Template: "empty",
	Hook:     "",
	Backup:   "active app config file",
	Active:   false,
}

func (a App) FilterValue() string { return a.Name }
func (a App) Test() App           { return a }

type AppDelegate struct{ styles ListStyles }

func (a AppDelegate) Height() int  { return 1 }
func (a AppDelegate) Spacing() int { return 0 }

func (a AppDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return func() tea.Msg {
		return updateTemplatesMsg(m.SelectedItem().(App))
	}
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

func GetAppListItems() []list.Item {
	data, err := os.ReadFile("./config/apps.yaml")
	if err != nil {
		panic(err)
	}

	apps := make(map[string]App)

	err = yaml.Unmarshal([]byte(data), &apps)
	if err != nil {
		panic(err)
	}

	//dirs, err := os.ReadDir("./config/templates/")
	//if err != nil {
	//	panic(err)
	//}

	// If apps from config file and number of app dirs dont align create new default apps for those missing from config file
	//if len(apps) != len(dirs) {
	//	newApp := []App{}

	//	for _, dir := range dirs {
	//		dirName := strings.ToLower(dir.Name())
	//		if _, ok := apps[dirName]; !ok {
	//			apps[dirName] = App{}
	//		}
	//	}
	//}

	var appListItems []list.Item

	for _, app := range apps {
		appListItems = append(appListItems, list.Item(app))
	}

	return appListItems
}

func UpdateAppList(newApp App, appList []list.Item) tea.Cmd {

	apps := make(map[string]App)

	for _, item := range appList {
		app := item.(App)
		apps[strings.ToLower(app.Name)] = app
	}

	apps[strings.ToLower(newApp.Name)] = newApp

	d, err := yaml.Marshal(&apps)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./config/apps.yaml", d, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(("./config/templates/" + strings.ToLower(newApp.Name)), os.ModePerm)
	if err != nil {
		panic(err)
	}

	var appListItems []list.Item

	for _, app := range apps {
		appListItems = append(appListItems, list.Item(app))
	}

	return func() tea.Msg {
		return updateAppListMsg(appListItems)
	}
}

// Template List

type Template struct {
	Name     string
	Filename string
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

func newLists() (list.Model, list.Model) {
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
	templateList.SetItems(GetTemplates(selectedApp))

	return appList, templateList
}
