package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ClaraSmyth/pin/builder"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

type mapData struct {
	Data map[string]App
	mu   sync.Mutex
}

func ApplyThemeCmd(theme Theme, themeList []list.Item) tea.Cmd {
	return func() tea.Msg {
		err := applyTheme(theme)

		if err != nil {
			for i, item := range themeList {
				if item.(Theme).Path == theme.Path {
					theme.Err = true
					themeList[i] = theme
					break
				}
			}
			return updateThemeListMsg(themeList)
		}

		return UpdateActiveStyles()
	}
}

func applyTheme(theme Theme) error {
	rawData, err := os.ReadFile(config.Paths.Apps)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			dir := filepath.Dir(config.Paths.ActiveTheme)
			err := os.MkdirAll(dir, 0777)
			if err != nil {
				return err
			}

			err = os.WriteFile(config.Paths.ActiveTheme, []byte(theme.Path), 0666)
			if err != nil {
				return err
			}

			return nil
		}
		return err
	}

	appsMap := make(map[string]App)

	err = yaml.Unmarshal([]byte(rawData), &appsMap)
	if err != nil {
		return err
	}

	themeData, err := os.ReadFile(theme.Path)
	if err != nil {
		return err
	}

	scheme := builder.Scheme{}

	err = yaml.Unmarshal([]byte(themeData), &scheme)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	data := mapData{Data: appsMap}

	for key, app := range data.Data {

		wg.Add(1)

		go func(key string, app App, data *mapData, wg *sync.WaitGroup) {

			defer wg.Done()

			if !app.Active || app.Path == "" || app.Template == "" {
				app.Active = false
				data.mu.Lock()
				defer data.mu.Unlock()
				data.Data[key] = app
				return
			}

			templates, err := os.ReadDir(filepath.Join(config.Paths.Templates, app.Name))
			if err != nil {
				return
			}

			var activeTemplatePath string

			for _, template := range templates {
				if filepath.Join(config.Paths.Templates, app.Name, template.Name()) == app.Template {
					activeTemplatePath = app.Template
				}

				if strings.Split(template.Name(), ".")[0] == theme.Name {
					activeTemplatePath = filepath.Join(config.Paths.Templates, app.Name, template.Name())
					break
				}
			}

			template, err := os.ReadFile(activeTemplatePath)
			if err != nil {
				app.Active = false
				app.Template = ""
				data.mu.Lock()
				defer data.mu.Unlock()
				data.Data[key] = app
				return
			}

			completeTemplate, err := builder.BuildTemplate(scheme, template)
			if err != nil {
				return
			}

			if app.Rewrite {
				err = os.WriteFile(app.Path, []byte(completeTemplate), 0666)
				if err != nil {
					return
				}
			}

			if !app.Rewrite {
				configFileData, err := os.ReadFile(app.Path)
				if err != nil {
					app.Active = false
					app.Path = ""
					data.mu.Lock()
					defer data.mu.Unlock()
					data.Data[key] = app
					return
				}

				updatedData := insertTemplate(string(configFileData), config.InsertStart, config.InsertEnd, completeTemplate)

				err = os.WriteFile(app.Path, []byte(strings.TrimSpace(updatedData)), 0666)
				if err != nil {
					return
				}
			}
		}(key, app, &data, &wg)
	}

	wg.Wait()

	wg2 := sync.WaitGroup{}

	if theme.Hook != "" {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			shellArgs := strings.Fields(config.DefaultShell)
			cmdArgs := append(shellArgs, theme.Hook)

			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
			_ = cmd.Run()
		}()
	}

	for _, app := range appsMap {
		if app.Hook != "" {
			wg2.Add(1)

			go func(hook string) {
				defer wg2.Done()

				shellArgs := strings.Fields(config.DefaultShell)
				cmdArgs := append(shellArgs, hook)

				cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
				_ = cmd.Run()
			}(app.Hook)
		}
	}

	wg2.Wait()

	err = os.WriteFile(config.Paths.ActiveTheme, []byte(theme.Path), 0666)
	if err != nil {
		return err
	}

	WriteAppData(appsMap)

	return nil
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
