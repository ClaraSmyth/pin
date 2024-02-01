package main

import (
	"cmp"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ClaraSmyth/pin/builder"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

func WriteAppData(appsMap map[string]App) {
	d, err := yaml.Marshal(&appsMap)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./config/apps.yaml", d, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

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

	entries, err := os.ReadDir("./config/templates/")
	if err != nil {
		panic(err)
	}

	appListItems := []list.Item{}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		app, exists := appsMap[name]

		if exists {
			appListItems = append(appListItems, app)
			delete(appsMap, name)
		} else {
			newApp := App{Name: name}
			appListItems = append(appListItems, newApp)
		}
	}

	// Append any remaining apps that have a missing dir + create a dir
	for _, app := range appsMap {
		err = os.Mkdir(("./config/templates/" + app.Name), os.ModePerm)
		if err != nil {
			panic(err)
		}

		appListItems = append(appListItems, app)
	}

	slices.SortFunc[[]list.Item, list.Item](appListItems, func(a, b list.Item) int {
		return cmp.Compare(a.(App).Name, b.(App).Name)
	})

	return appListItems
}

func CreateApp(newApp App, appList []list.Item) tea.Cmd {
	return func() tea.Msg {
		if newApp.Name == "" {
			return nil
		}

		appList = append(appList, newApp)
		appsMap := make(map[string]App)

		for _, item := range appList {
			app := item.(App)
			appsMap[app.Name] = app
		}

		WriteAppData(appsMap)

		err := os.Mkdir(("./config/templates/" + newApp.Name), os.ModePerm)
		if err != nil {
			panic(err)
		}

		slices.SortFunc[[]list.Item, list.Item](appList, func(a, b list.Item) int {
			return cmp.Compare(a.(App).Name, b.(App).Name)
		})

		templateList := []list.Item{}

		return updateAppListMsg{
			appListItems:      appList,
			templateListItems: templateList,
		}
	}
}

func EditApp(newApp App, prevApp App, prevList []list.Item) tea.Cmd {
	return func() tea.Msg {
		newList := []list.Item{newApp}
		appsMap := make(map[string]App)
		appsMap[newApp.Name] = newApp

		for _, item := range prevList {
			app := item.(App)

			if app.Name == prevApp.Name {
				continue
			}

			appsMap[app.Name] = app
			newList = append(newList, app)
		}

		WriteAppData(appsMap)

		if newApp.Name != prevApp.Name {
			basePath := "./config/templates/"
			prevPath := basePath + prevApp.Name
			newPath := basePath + newApp.Name

			err := os.Rename(prevPath, newPath)
			if err != nil {
				panic(err)
			}
		}

		newAppTemplates := GetTemplates(newApp)

		slices.SortFunc[[]list.Item, list.Item](newList, func(a, b list.Item) int {
			return cmp.Compare(a.(App).Name, b.(App).Name)
		})

		return updateAppListMsg{
			appListItems:      newList,
			templateListItems: newAppTemplates,
		}
	}
}

func DeleteApp(prevApp App, prevList []list.Item) tea.Cmd {
	return func() tea.Msg {
		newList := []list.Item{}
		appsMap := make(map[string]App)

		for _, item := range prevList {
			app := item.(App)

			if app.Name == prevApp.Name {
				continue
			}

			appsMap[app.Name] = app
			newList = append(newList, app)
		}

		WriteAppData(appsMap)

		err := os.RemoveAll("./config/templates/" + prevApp.Name)
		if err != nil {
			panic(err)
		}

		newAppTemplates := []list.Item{}

		slices.SortFunc[[]list.Item, list.Item](newList, func(a, b list.Item) int {
			return cmp.Compare(a.(App).Name, b.(App).Name)
		})
		return updateAppListMsg{
			appListItems:      newList,
			templateListItems: newAppTemplates,
		}
	}
}

func GetTemplates(app App) []list.Item {
	path := "./config/templates/" + app.Name

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
		active := false
		if strings.Contains(app.Template, filename) {
			active = true
		}

		template := Template{Name: strings.Split(filename, ".")[0], Path: path + "/" + filename, AppPath: path, Active: active}
		templateList = append(templateList, list.Item(template))
	}

	return templateList
}

func UpdateTemplates(app App) tea.Cmd {
	return func() tea.Msg {
		templates := GetTemplates(app)
		return updateTemplateListMsg(templates)
	}
}

func CreateTemplate(app App, filename string) tea.Cmd {
	return func() tea.Msg {
		if filename == "" {
			return nil
		}

		_, err := os.Create("./config/templates/" + app.Name + "/" + filename + ".mustache")
		if err != nil {
			panic(err)
		}

		templates := GetTemplates(app)
		return updateTemplateListMsg(templates)
	}
}

func EditTemplate(app App, prevFilename string, newFilename string) tea.Cmd {
	return func() tea.Msg {
		basePath := "./config/templates/"
		prevPath := basePath + app.Name + "/" + prevFilename + ".mustache"
		newPath := basePath + app.Name + "/" + newFilename + ".mustache"

		err := os.Rename(prevPath, newPath)
		if err != nil {
			panic(err)
		}

		templates := GetTemplates(app)
		return updateTemplateListMsg(templates)
	}
}

func DeleteTemplate(app App, filename string) tea.Cmd {
	return func() tea.Msg {
		basePath := "./config/templates/"
		path := basePath + app.Name + "/" + filename + ".mustache"
		err := os.Remove(path)
		if err != nil {
			panic(err)
		}

		templates := GetTemplates(app)
		return updateTemplateListMsg(templates)
	}
}

func GetThemes() []list.Item {
	basePath := "./config/schemes/"

	activeThemePath, _ := os.ReadFile("./config/activeTheme")

	themeList := []list.Item{}

	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if strings.Contains(d.Name(), ".yaml") && strings.Contains(path, "base16") {

			if strings.Contains(string(activeThemePath), d.Name()) {
				themeList = append(themeList, Theme{Name: strings.Split(d.Name(), ".")[0], Path: path, Active: true})
			} else {
				themeList = append(themeList, Theme{Name: strings.Split(d.Name(), ".")[0], Path: path, Active: false})
			}

		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return themeList
}

func CreateTheme(themeName string, themeList []list.Item) tea.Cmd {
	return func() tea.Msg {
		basePath := "./config/schemes/custom/"

		if themeName == "" {
			return nil
		}

		f, err := os.Create(basePath + themeName + ".yaml")
		if err != nil {
			panic(err)
		}

		defer f.Close()

		themeList := GetThemes()

		return updateThemeListMsg(themeList)
	}
}

func DeleteTheme(theme Theme) tea.Cmd {
	return func() tea.Msg {
		err := os.Remove(theme.Path)
		if err != nil {
			panic(err)
		}

		themeList := GetThemes()

		return updateThemeListMsg(themeList)
	}
}

func GitCloneSchemes() tea.Cmd {
	return func() tea.Msg {
		repo := "https://github.com/tinted-theming/schemes.git"
		target := "./config/schemes/tinted-theming/"

		err := os.RemoveAll(target)
		if err != nil {
			panic(err)
		}

		cmd := exec.Command("git", "clone", repo, target)
		err = cmd.Run()
		if err != nil {
			panic(err)
		}

		themeList := GetThemes()

		return updateThemeListMsg(themeList)
	}
}

func GetActiveColors() Colors {
	activeTheme, err := os.ReadFile("./config/activeTheme")
	if err != nil {
		return DefaultColors()
	}

	file, err := os.ReadFile(string(activeTheme))
	if err != nil {
		return DefaultColors()
	}

	scheme := builder.Scheme{}

	err = yaml.Unmarshal([]byte(file), &scheme)
	if err != nil {
		return DefaultColors()
	}

	if len(scheme.Palette) != 16 {
		return DefaultColors()
	}

	for i, clr := range scheme.Palette {
		c, err := builder.ParseHexColor(clr)
		if err != nil {
			return DefaultColors()
		}
		scheme.Palette[i] = fmt.Sprintf("%02x%02x%02x", c.R, c.G, c.B)
	}

	return Colors{
		Base00: lipgloss.Color("#" + scheme.Palette["base00"]),
		Base01: lipgloss.Color("#" + scheme.Palette["base01"]),
		Base02: lipgloss.Color("#" + scheme.Palette["base02"]),
		Base03: lipgloss.Color("#" + scheme.Palette["base03"]),
		Base04: lipgloss.Color("#" + scheme.Palette["base04"]),
		Base05: lipgloss.Color("#" + scheme.Palette["base05"]),
		Base06: lipgloss.Color("#" + scheme.Palette["base06"]),
		Base07: lipgloss.Color("#" + scheme.Palette["base07"]),
		Base08: lipgloss.Color("#" + scheme.Palette["base08"]),
		Base09: lipgloss.Color("#" + scheme.Palette["base09"]),
		Base0A: lipgloss.Color("#" + scheme.Palette["base0A"]),
		Base0B: lipgloss.Color("#" + scheme.Palette["base0B"]),
		Base0C: lipgloss.Color("#" + scheme.Palette["base0C"]),
		Base0D: lipgloss.Color("#" + scheme.Palette["base0D"]),
		Base0E: lipgloss.Color("#" + scheme.Palette["base0E"]),
		Base0F: lipgloss.Color("#" + scheme.Palette["base0F"]),
	}
}

func UpdateActiveStyles() tea.Msg {
	colors := GetActiveColors()
	styles := DefaultStyles(colors)
	return updateStylesMsg(styles)
}
