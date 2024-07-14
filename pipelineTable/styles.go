package pipelineTable

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("12"))

func getInfoTableStyles() table.Styles {
	return table.Styles{
		Selected: lipgloss.NewStyle().Bold(false),
		Header: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false).Padding(0, 1),
		Cell: lipgloss.NewStyle().Padding(0, 1),
	}
}
