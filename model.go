package main

import (
	"github.com/FredrikMWold/radix-tui/applicationTable"
	"github.com/FredrikMWold/radix-tui/environmentTable"
	"github.com/FredrikMWold/radix-tui/pipelineTable"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	applicationsTable tea.Model
	pipelineTable     tea.Model
	enviromentTable   tea.Model
}

func initialModel() Model {

	applicationTable := applicationTable.New()
	pipelineTable := pipelineTable.New()
	enviromentTable := environmentTable.New()

	return Model{
		applicationsTable: applicationTable,
		pipelineTable:     pipelineTable,
		enviromentTable:   enviromentTable,
	}
}
