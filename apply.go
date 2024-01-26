package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/ClaraSmyth/pin/builder"
	"gopkg.in/yaml.v3"
)

func ApplyTheme(theme Theme) {
	rawData, err := os.ReadFile("./config/apps.yaml")
	if err != nil {
		panic(err)
	}

	appsMap := make(map[string]App)

	err = yaml.Unmarshal([]byte(rawData), &appsMap)
	if err != nil {
		panic(err)
	}

	for _, app := range appsMap {
		go func(app App) {

			basePath := "./config/templates/"

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

			var activeTemplate string

			for _, template := range templates {
				if basePath+app.Name+"/"+template.Name() == app.Template {
					activeTemplate = app.Template
				}

				if strings.Split(template.Name(), ".")[0] == theme.Name {
					activeTemplate = basePath + app.Name + "/" + template.Name()
					break
				}
			}

			if activeTemplate == "" {
				return
			}

			data := builder.BuildTemplate(theme.Path, activeTemplate)

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
		}(app)
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
