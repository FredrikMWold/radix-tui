package applicationTable

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SelectedApplication string

func selectApplication(application string) tea.Cmd {
	return func() tea.Msg {
		return SelectedApplication(application)
	}
}
