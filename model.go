package main

import (
	appTable "radix.go/applicationTable"
	envTable "radix.go/environmentTable"
	pipeTable "radix.go/pipelineTable"
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
