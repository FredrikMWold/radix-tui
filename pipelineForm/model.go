package pipelineForm

import (
	"github.com/charmbracelet/huh"
)

type Model struct {
	form                *huh.Form
	SelectedApplication string
	environments        []string
	branches            []string
	width               int
}

func New() Model {
	return Model{
		environments: []string{},
		branches:     []string{},
		form:         huh.NewForm(huh.NewGroup(huh.NewInput().Key("name"))),
	}
}
