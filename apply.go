package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/ClaraSmyth/pin/builder"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

func ApplyTheme(theme Theme) tea.Cmd {
	return func() tea.Msg {
		rawData, err := os.ReadFile("./config/apps.yaml")
		if err != nil {
			panic(err)
		}

		appsMap := make(map[string]App)

		err = yaml.Unmarshal([]byte(rawData), &appsMap)
		if err != nil {
			panic(err)
		}

		themeData, err := os.ReadFile(theme.Path)
		if err != nil {
			panic(err)
		}

		scheme := builder.Scheme{}

		err = yaml.Unmarshal([]byte(themeData), &scheme)
		if err != nil {
			panic(err)
		}

		for key, app := range appsMap {
			basePath := "./config/templates/"

			if !app.Active || app.Path == "" || app.Template == "" {
				app.Active = false
				appsMap[key] = app
				continue
			}

			templates, err := os.ReadDir(basePath + app.Name)
			if err != nil {
				continue
			}

			var activeTemplatePath string

			for _, template := range templates {
				if basePath+app.Name+"/"+template.Name() == app.Template {
					activeTemplatePath = app.Template
				}

				if strings.Split(template.Name(), ".")[0] == theme.Name {
					activeTemplatePath = basePath + app.Name + "/" + template.Name()
					break
				}
			}

			template, err := os.ReadFile(activeTemplatePath)
			if err != nil {
				app.Active = false
				app.Template = ""
				appsMap[key] = app
				continue
			}

			data, err := builder.BuildTemplate(scheme, template)
			if err != nil {
				panic(err)
			}

			if app.Rewrite {
				err = os.WriteFile(app.Path, []byte(data), os.ModePerm)
				if err != nil {
					panic(err)
				}
			}

			if !app.Rewrite {
				configFileData, err := os.ReadFile(app.Path)
				if err != nil {
					app.Active = false
					app.Path = ""
					appsMap[key] = app
					continue
				}

				updatedData := insertTemplate(string(configFileData), "START_PIN_HERE", "END_PIN_HERE", data)

				err = os.WriteFile(app.Path, []byte(strings.TrimSpace(updatedData)), os.ModePerm)
				if err != nil {
					panic(err)
				}
			}

			if app.Hook != "" {
				cmd := exec.Command("sh", "-c", app.Hook)
				err = cmd.Run()
				if err != nil {
					panic(err)
				}
			}

		}

		err = os.WriteFile("./config/activeTheme", []byte(theme.Path), os.ModePerm)
		if err != nil {
			panic(err)
		}

		WriteAppData(appsMap)

		return UpdateActiveStyles()
	}
}

func insertTemplate(fileData, startString, endString, template string) string {
	lines := strings.Split(fileData, "\n")

	newData := ""
	startFound := false
	endFound := false
	canInsert := true

	for _, line := range lines {
		if strings.Contains(line, startString) {
			startFound = true
			canInsert = false
			newData += line + "\n"
			newData += template + "\n"
			continue
		}

		if strings.Contains(line, endString) {
			endFound = true
			canInsert = true
		}

		if canInsert {
			newData += line + "\n"
		}
	}

	if !startFound || !endFound {
		return fileData
	}

	return newData
}
