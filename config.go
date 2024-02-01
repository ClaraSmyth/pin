package main

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type Config struct {
	Home        string
	Apps        string
	Templates   string
	Schemes     string
	ActiveTheme string
}

var config = defaultConfig()

func defaultConfig() Config {
	homePath := os.Getenv("PIN_HOME")
	if homePath == "" {
		homePath = xdg.ConfigHome
	}

	return Config{
		Home:        filepath.Join(homePath, "pin"),
		Apps:        filepath.Join(homePath, "pin", "apps.yaml"),
		Templates:   filepath.Join(homePath, "pin", "templates"),
		Schemes:     filepath.Join(homePath, "pin", "schemes"),
		ActiveTheme: filepath.Join(homePath, "pin", "activeTheme"),
	}
}
