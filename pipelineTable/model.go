package pipelineTable

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func New() Model {
	spiner := spinner.New()
	spiner.Spinner = spinner.Meter
	spiner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Model{
		table: table.New(
			table.WithFocused(false),
			table.WithStyles(getInfoTableStyles()),
		),
		spinner:             spiner,
		selectedApplication: "No application selected",
	}
}

type Model struct {
	table                table.Model
	isLoadingApplication bool
	selectedApplication  string
	spinner              spinner.Model
}
