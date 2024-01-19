package main

import (
	"cmp"
	"io/fs"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

func GetApps() []list.Item {
	rawData, err := os.ReadFile("./config/apps.yaml")
	if err != nil {
		panic(err)
	}

	appsMap := make(map[string]App)

	err = yaml.Unmarshal([]byte(rawData), &appsMap)
	if err != nil {
		panic(err)
	}

	dirs, err := os.ReadDir("./config/templates/")
	if err != nil {
		panic(err)
	}

	dirsMap := make(map[string]fs.DirEntry)

	for _, dir := range dirs {
		dirsMap[strings.ToLower(dir.Name())] = dir
	}

	appListItems := []list.Item{}

	for _, app := range appsMap {
		appListItems = append(appListItems, list.Item(app))
	}

	return appListItems
}

func CreateApp(newApp App, appList []list.Item) tea.Cmd {
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

	slices.SortFunc[[]list.Item, list.Item](appListItems, func(a, b list.Item) int {
		return cmp.Compare(a.(App).Name, b.(App).Name)
	})

	return func() tea.Msg {
		return updateAppListMsg(appListItems)
	}
}

func EditApp(newApp App, prevApp App, appList []list.Item) tea.Cmd {
	apps := make(map[string]App)

	for _, item := range appList {
		app := item.(App)

		if app.Name == prevApp.Name {
			continue
		}

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

	basePath := "./config/templates/"
	prevPath := basePath + strings.ToLower(prevApp.Name)
	newPath := basePath + strings.ToLower(newApp.Name)

	err = os.Rename(prevPath, newPath)

	var appListItems []list.Item

	for _, app := range apps {
		appListItems = append(appListItems, list.Item(app))
	}

	slices.SortFunc[[]list.Item, list.Item](appListItems, func(a, b list.Item) int {
		return cmp.Compare(a.(App).Name, b.(App).Name)
	})

	return func() tea.Msg {
		return updateAppListMsg(appListItems)
	}
}

func DeleteApp(prevApp App, appList []list.Item) tea.Cmd {
	apps := make(map[string]App)

	for _, item := range appList {
		app := item.(App)

		if app.Name == prevApp.Name {
			continue
		}

		apps[strings.ToLower(app.Name)] = app
	}

	d, err := yaml.Marshal(&apps)
	if err != nil {
		log.Fatal("1")
	}

	err = os.WriteFile("./config/apps.yaml", d, os.ModePerm)
	if err != nil {
		log.Fatal("2")
	}

	err = os.RemoveAll("./config/templates/" + strings.ToLower(prevApp.Name))
	if err != nil {
		log.Fatal("3")
	}

	var appListItems []list.Item

	for _, app := range apps {
		appListItems = append(appListItems, list.Item(app))
	}

	slices.SortFunc[[]list.Item, list.Item](appListItems, func(a, b list.Item) int {
		return cmp.Compare(a.(App).Name, b.(App).Name)
	})

	return func() tea.Msg {
		return updateAppListMsg(appListItems)
	}
}

func GetTemplates(appName string) []list.Item {
	path := "./config/templates/" + strings.ToLower(appName)

	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	templateList := []list.Item{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		template := Template{Name: strings.Split(filename, ".")[0], Filename: filename}
		templateList = append(templateList, list.Item(template))
	}

	return templateList
}

func CreateTemplate(app App, filename string) tea.Cmd {
	f, err := os.Create("./config/templates/" + strings.ToLower(app.Name) + "/" + filename)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	return func() tea.Msg {
		return updateTemplatesMsg(app)
	}
}

func EditTemplate(app App, prevFilename string, newFilename string) tea.Cmd {
	basePath := "./config/templates/"
	prevPath := basePath + app.Name + "/" + prevFilename
	newPath := basePath + app.Name + "/" + newFilename

	err := os.Rename(prevPath, newPath)
	if err != nil {
		panic(err)
	}

	return func() tea.Msg {
		return updateTemplatesMsg(app)
	}
}

func DeleteTemplate(app App, filename string) tea.Cmd {
	basePath := "./config/templates/"
	path := basePath + app.Name + "/" + filename
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}

	return func() tea.Msg {
		return updateTemplatesMsg(app)
	}
}
