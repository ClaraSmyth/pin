package main

import (
	"os"
	"os/exec"
	"strings"
)

func editorCmd(path string) *exec.Cmd {
	editor, args := getEditor()
	return exec.Command(editor, append(args, path)...)
}

func getEditor() (string, []string) {
	editorConfig := strings.Fields(config.DefaultEditor)
	if len(editorConfig) > 0 {
		return editorConfig[0], editorConfig[1:]
	}

	editorEnv := strings.Fields(os.Getenv("EDITOR"))
	if len(editorEnv) > 0 {
		return editorEnv[0], editorEnv[1:]
	}

	return "nano", nil
}
