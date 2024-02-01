package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"sync"

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

		var wg sync.WaitGroup

		for _, app := range appsMap {
			go func(app App, wg *sync.WaitGroup) {
				basePath := "./config/templates/"
				wg.Add(1)

				if !app.Active || app.Path == "" || app.Template == "" {
					return
				}

				templates, err := os.ReadDir(basePath + app.Name)
				if err != nil {
					return
				}

				if len(templates) == 0 {
					return
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

				if activeTemplatePath == "" {
					return
				}

				template, err := os.ReadFile(activeTemplatePath)
				if err != nil {
					panic(err)
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
					insertTemplate(app.Path, "START_PIN_HERE", "END_PIN_HERE", data)
				}

				if app.Hook != "" {
					cmd := exec.Command("sh", "-c", app.Hook)
					err = cmd.Run()
					if err != nil {
						panic(err)
					}
				}

			}(app, &wg)
		}

		wg.Wait()

		err = os.WriteFile("./config/activeTheme", []byte(theme.Path), os.ModePerm)
		if err != nil {
			panic(err)
		}

		return UpdateActiveStyles()
	}
}

func insertTemplate(path, startString, endString, template string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")

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
		panic(errors.New("couldnt find start or end point"))
	}

	err = os.WriteFile(path, []byte(strings.TrimSpace(newData)), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
