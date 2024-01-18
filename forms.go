package main

import (
	"github.com/charmbracelet/huh"
)

var (
	formName  string
	formHook  string
	formApply bool
)

func resetFormValues() {
	formName = ""
	formHook = ""
	formApply = false
}

func newForms() map[Pane]*huh.Form {

	appForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("App name").
				Value(&formName),

			huh.NewInput().
				Key("hook").
				Title("Hook").
				Value(&formHook),

			huh.NewConfirm().
				Key("apply").
				Title("Apply?").
				Value(&formApply),
		),
	).WithShowHelp(false).WithWidth(20)

	templateForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Template name").
				Value(&formName),

			huh.NewConfirm().
				Key("apply").
				Title("Apply?").
				Value(&formApply),
		),
	).WithShowHelp(false).WithWidth(20)

	forms := map[Pane]*huh.Form{
		appPane:      appForm,
		templatePane: templateForm,
	}

	return forms
}
