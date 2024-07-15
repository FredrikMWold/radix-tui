package main

import (
	appTable "github.com/FredrikMWold/radix-tui/applicationTable"
	envTable "github.com/FredrikMWold/radix-tui/environmentTable"
	pipeTable "github.com/FredrikMWold/radix-tui/pipelineTable"
)

type Model struct {
	applicationsTable appTable.Model
	pipelineTable     pipeTable.Model
	enviromentTable   envTable.Model
}

func initialModel() Model {

	applicationTable := appTable.New()
	pipelineTable := pipeTable.New()
	enviromentTable := envTable.New()

	return Model{
		applicationsTable: applicationTable,
		pipelineTable:     pipelineTable,
		enviromentTable:   enviromentTable,
	}
}
