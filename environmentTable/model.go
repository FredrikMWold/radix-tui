package environmentTable

import (
	"github.com/FredrikMWold/radix-tui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	table                table.Model
	spinner              spinner.Model
	isLoadingApplication bool
}

func New() Model {
	spiner := spinner.New()
	spiner.Spinner = spinner.Meter
	spiner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Model{
		table: table.New(
			table.WithFocused(false),
			table.WithHeight(4),
			table.WithStyles(styles.InfoTableStyles()),
			table.WithColumns([]table.Column{
				{Title: "name", Width: 7},
				{Title: "status", Width: 10},
				{Title: "branch", Width: 7},
			}),
		),
		spinner:              spiner,
		isLoadingApplication: false,
	}
}
