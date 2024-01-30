package main

import (
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/ClaraSmyth/pin/builder"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

func ApplyTheme(theme Theme) tea.Cmd {
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

	err = os.WriteFile("./config/activeTheme.yaml", []byte(themeData), os.ModePerm)
	if err != nil {
		panic(err)
	}

	return func() tea.Msg {
		return UpdateActiveStyles()
	}
}

func insertTemplate(path, startString, endString, template string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	dataString := string(data)

	startIndex := len(startString) + strings.Index(dataString, startString)
	endIndex := strings.Index(dataString, endString)

	if startIndex == -1 || endIndex == -1 {
		panic(err)
	}

	endLineStart := strings.LastIndex(dataString[:endIndex], "\n") + 1

	newData := dataString[:startIndex] + "\n" + template + "\n" + dataString[endLineStart:endIndex] + dataString[endIndex:]

	err = os.WriteFile(path, []byte(newData), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
