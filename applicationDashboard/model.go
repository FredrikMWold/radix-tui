package appllicationDashboard

import (
	"github.com/FredrikMWold/radix-tui/applicationTable"
	"github.com/FredrikMWold/radix-tui/environmentTable"
	"github.com/FredrikMWold/radix-tui/pipelineTable"
	tea "github.com/charmbracelet/bubbletea"
)

type Focused int

const (
	application Focused = iota
	pipeline
)

type Model struct {
	applicationsTable applicationTable.Model
	pipelineTable     pipelineTable.Model
	enviromentTable   tea.Model
	focused           Focused
}

func New() Model {

	applicationTable := applicationTable.New()
	pipelineTable := pipelineTable.New()
	enviromentTable := environmentTable.New()

	return Model{
		applicationsTable: applicationTable,
		pipelineTable:     pipelineTable,
		enviromentTable:   enviromentTable,
		focused:           application,
	}
}
