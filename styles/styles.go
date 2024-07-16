package styles

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func SectionContainer(isActive bool) lipgloss.Style {
	styles := lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
	if isActive {
		styles.BorderForeground(lipgloss.Color("26"))
	} else {
		styles.BorderForeground(lipgloss.Color("69"))
	}
	return styles
}

func LoadingSpinnerContainer(height int, width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
}

func TableStyles() table.Styles {
	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		BorderTop(true).
		Bold(false)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	return tableStyles
}

func InfoTableStyles() table.Styles {
	return table.Styles{
		Selected: lipgloss.NewStyle().Bold(false),
		Header: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			BorderTop(true).
			Bold(false).Padding(0, 1),
		Cell: lipgloss.NewStyle().Padding(0, 1),
	}
}
