package applyconfig

import (
	"github.com/charmbracelet/huh"
)

type Model struct {
	form                *huh.Form
	selectedApplication string
	width               int
}

func NewForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("applyConfig").
				Affirmative("Create job").
				Negative("Cancel").
				WithTheme(huh.ThemeCatppuccin()),
		),
	)
}

func New() Model {
	return Model{
		form: NewForm(),
	}
}
