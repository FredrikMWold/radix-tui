package appllicationDashboard

import (
	"github.com/FredrikMWold/radix-tui/applicationTable"
	"github.com/FredrikMWold/radix-tui/commands"
	"github.com/FredrikMWold/radix-tui/environmentTable"
	"github.com/FredrikMWold/radix-tui/pipelineForm"
	"github.com/FredrikMWold/radix-tui/pipelineTable"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	keys                 keyMap
	pipelineTable        pipelineTable.Model
	spinner              spinner.Model
	help                 help.Model
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
		help:                 help.New(),
		spinner:              spiner,
		focused:              application,
		isLoadingApplication: false,
		keys:                 ApplicationTableKeys,
	}
}

type keyMap struct {
	Enter       key.Binding
	Up          key.Binding
	Down        key.Binding
	Esc         key.Binding
	BuildDeploy key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Up, k.Down, k.Esc, k.BuildDeploy}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Up, k.Down, k.Esc, k.BuildDeploy},
	}
}

var ApplicationTableKeys = keyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select application"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("down", "move down"),
	),
}

var PipelineTableKeys = keyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("down", "move down"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	BuildDeploy: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "build-deploy"),
	),
}

var BuildDeployFormKeys = keyMap{
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}
