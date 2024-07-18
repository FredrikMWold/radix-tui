package appllicationDashboard

import (
	"github.com/FredrikMWold/radix-tui/applicationTable"
	"github.com/FredrikMWold/radix-tui/commands"
	"github.com/FredrikMWold/radix-tui/environmentTable"
	"github.com/FredrikMWold/radix-tui/pipelineForm"
	"github.com/FredrikMWold/radix-tui/pipelineTable"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Focused int

const (
	application Focused = iota
	pipeline
	form
)

type Model struct {
	applicationsTable    applicationTable.Model
	pipelineTable        pipelineTable.Model
	spinner              spinner.Model
	pipelineForm         tea.Model
	enviromentTable      tea.Model
	focused              Focused
	applications         []string
	isLoadingApplication bool
	height               int
	width                int
	application          commands.Application
}

func New() Model {

	applicationTable := applicationTable.New()
	pipelineTable := pipelineTable.New()
	enviromentTable := environmentTable.New()
	pipelineForm := pipelineForm.New()

	spiner := spinner.New()
	spiner.Spinner = spinner.Meter
	spiner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		applicationsTable:    applicationTable,
		pipelineTable:        pipelineTable,
		enviromentTable:      enviromentTable,
		pipelineForm:         pipelineForm,
		spinner:              spiner,
		focused:              application,
		isLoadingApplication: false,
	}
}
