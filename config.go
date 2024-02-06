package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DefaultShell string `yaml:"DefaultShell"`
	InsertStart  string `yaml:"InsertStart"`
	InsertEnd    string `yaml:"InsertEnd"`
	Paths        Paths  `yaml:"-"`
}

type Paths struct {
	Home          string
	Apps          string
	Templates     string
	ActiveTheme   string
	CustomSchemes string
	BaseSchemes   string
}

var config = readConfig()

func readConfig() Config {
	homePath := os.Getenv("PIN_HOME")
	if homePath == "" {
		homePath = xdg.ConfigHome
	}

	dataPath := os.Getenv("PIN_DATA")
	if dataPath == "" {
		dataPath = xdg.DataHome
	}

	configFile, err := os.ReadFile(filepath.Join(homePath, "pin", "config.yaml"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.MkdirAll(filepath.Join(homePath, "pin"), 0777)
			if err != nil {
				panic(err)
			}

			err = os.WriteFile(filepath.Join(homePath, "pin", "config.yaml"), []byte(strings.TrimSpace(defaultConfigFile)), 0666)
			if err != nil {
				panic(err)

			}
		} else {
			panic(err)
		}
	}

	configYaml := Config{}
	err = yaml.Unmarshal(configFile, &configYaml)
	if err != nil {
		panic(err)
	}

	if configYaml.DefaultShell == "" {
		configYaml.DefaultShell = "sh -c"
	}

	if configYaml.InsertStart == "" {
		configYaml.InsertStart = "START_PIN_HERE"
	}

	if configYaml.InsertEnd == "" {
		configYaml.InsertEnd = "END_PIN_HERE"
	}

	configYaml.Paths = Paths{
		Home:          filepath.Join(homePath, "pin"),
		Apps:          filepath.Join(homePath, "pin", "apps.yaml"),
		Templates:     filepath.Join(homePath, "pin", "templates"),
		ActiveTheme:   filepath.Join(homePath, "pin", "activeTheme"),
		CustomSchemes: filepath.Join(homePath, "pin", "schemes"),
		BaseSchemes:   filepath.Join(dataPath, "pin", "schemes"),
	}

	return configYaml
}

var defaultConfigFile = `
# Change the default shell and any required args
DefaultShell: sh -c

# Change the default start string to search for when inserting templates
InsertStart: START_PIN_HERE

#Change the default end string to search for when inserting templates
InsertEnd: END_PIN_HERE
`
