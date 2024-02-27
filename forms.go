package main

import (
	"errors"
	"slices"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/huh"
)

type FormType int

const (
	appForm FormType = iota
	appEditForm
	templateForm
	themeForm
)

var (
	formEdit       bool
	formName       string
	formHook       string
	formRewrite    bool
	formFilepicker bool
	formApply      bool
)

func newForm(formType FormType, items []list.Item, theme *huh.Theme) *huh.Form {

	confirmNegativeText := "Cancel"

	if formType == appEditForm {
		confirmNegativeText = "Skip"
	}

	switch formType {
	case appForm, appEditForm:
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
							return strings.EqualFold(v.FilterValue(), str)
						}) && !formEdit {
							return errors.New("Already Exists!")
						}

						return nil
					}),

				huh.NewInput().
					Key("hook").
					Title("Hook").
					Value(&formHook),

				huh.NewConfirm().
					Key("rewrite").
					Title("Write method").
					Affirmative("Insert").
					Negative("Rewrite").
					Value(&formRewrite),

				huh.NewConfirm().
					Key("filepicker").
					Title("Select config file").
					Affirmative("Select").
					Negative(confirmNegativeText).
					Value(&formFilepicker),
			),
		).WithShowHelp(false).WithWidth(25).WithTheme(theme)

	case templateForm:
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
							return strings.EqualFold(v.FilterValue(), str)
						}) && !formEdit {
							return errors.New("Already Exists!")
						}

						return nil
					}),

				huh.NewConfirm().
					Key("apply").
					Title("Confirm?").
					Value(&formApply),
			),
		).WithShowHelp(false).WithWidth(25).WithTheme(theme)

	case themeForm:
		return huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("name").
					Title("Theme name").
					Value(&formName).
					Validate(func(str string) error {
						if str == "" {
							return errors.New("Cant be empty!")
						}

						if !validateFilename(str) {
							return errors.New("Invalid name!")
						}

						if slices.ContainsFunc[[]list.Item, list.Item](items, func(v list.Item) bool {
							return strings.EqualFold(v.FilterValue(), str)
						}) && !formEdit {
							return errors.New("Already Exists!")
						}

						return nil
					}),

				huh.NewInput().
					Key("hook").
					Title("Hook").
					Value(&formHook),

				huh.NewConfirm().
					Key("apply").
					Title("Confirm?").
					Value(&formApply),
			),
		).WithShowHelp(false).WithWidth(25).WithTheme(theme)

	default:
		return nil
	}
}

func deleteForm(theme *huh.Theme) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("apply").
				Title("Are you sure?").
				Value(&formApply),
		),
	).WithShowHelp(false).WithWidth(25).WithTheme(theme)
}

func validateFilename(filename string) bool {
	for _, v := range filename {
		if !unicode.IsLetter(v) && !unicode.IsDigit(v) && string(v) != "-" {
			return false
		}
	}
	return true
}
