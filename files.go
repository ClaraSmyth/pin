package main

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

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
