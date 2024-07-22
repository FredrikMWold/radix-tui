package buildanddeploy

import (
	"github.com/charmbracelet/huh"
)

type Model struct {
	form                *huh.Form
	SelectedApplication string
	branchMapping       map[string]string
	width               int
}

func New() Model {
	return Model{
		branchMapping: make(map[string]string),
		form:          huh.NewForm(huh.NewGroup(huh.NewInput().Key("name"))),
	}
}
