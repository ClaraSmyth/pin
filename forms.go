package main

import (
	"errors"
	"slices"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/huh"
)

var (
	formName       string
	formHook       string
	formFilepicker bool
	formApply      bool
)

func newForm(pane Pane, items []list.Item) *huh.Form {
	switch pane {
	case appPane:
		return huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("name").
					Title("App name").
					Value(&formName).
					Validate(func(str string) error {
						if str == "" {
							return errors.New("Cant be empty!")
						}

						if !validateFilename(str) {
							return errors.New("Invalid name!")
						}

						if slices.ContainsFunc[[]list.Item, list.Item](items, func(v list.Item) bool {
							return strings.ToLower(v.FilterValue()) == strings.ToLower(str)
						}) {
							return errors.New("Already Exists!")
						}

						return nil
					}),

				huh.NewInput().
					Key("hook").
					Title("Hook").
					Value(&formHook),

				huh.NewConfirm().
					Key("filepicker").
					Title("Select config file").
					Affirmative("Select").
					Negative("Cancel").
					Value(&formFilepicker),
			),
		).WithShowHelp(false).WithWidth(20)

	case templatePane:
		return huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("name").
					Title("Template name").
					Value(&formName).
					Validate(func(str string) error {
						if str == "" {
							return errors.New("Cant be empty!")
						}

						if !validateFilename(str) {
							return errors.New("Invalid name!")
						}

						if slices.ContainsFunc[[]list.Item, list.Item](items, func(v list.Item) bool {
							return strings.ToLower(v.FilterValue()) == strings.ToLower(str)
						}) {
							return errors.New("Already Exists!")
						}

						return nil
					}),

				huh.NewConfirm().
					Key("apply").
					Title("Apply?").
					Value(&formApply),
			),
		).WithShowHelp(false).WithShowErrors(false).WithWidth(20)

	default:
		return nil
	}
}

func deleteForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("apply").
				Title("Are you sure?").
				Value(&formApply),
		),
	).WithShowHelp(false).WithShowErrors(false).WithWidth(20)
}

func validateFilename(filename string) bool {
	for _, v := range filename {
		if !unicode.IsLetter(v) && !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}
