package applicationtable

import (
	"github.com/FredrikMWold/radix-tui/commands"
	tea "github.com/charmbracelet/bubbletea"
)

func selectApplication(application string) tea.Cmd {
	return func() tea.Msg {
		return commands.SelectedApplication(application)
	}
}
