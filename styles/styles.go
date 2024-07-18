package styles

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func SectionContainer(isActive bool) lipgloss.Style {
	if isActive {
		return lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#b4befe"))
	} else {
		return lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#6c7086"))
	}
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
		Foreground(lipgloss.Color("#11111b")).
		Background(lipgloss.Color("#cba6f7"))
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
