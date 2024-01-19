package main

import (
	"github.com/charmbracelet/huh"
)

var (
	formName  string
	formHook  string
	formApply bool
)

func newForm(pane Pane) *huh.Form {
	switch pane {

	case appPane:
		return huh.NewForm(
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

	case templatePane:
		return huh.NewForm(
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

	default:
		return nil
	}
}
