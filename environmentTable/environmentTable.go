package environmentTable

import (
	"github.com/FredrikMWold/radix-tui/commands"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	case commands.Application:
		enviroments := make([]table.Row, len(msg.Environments))
		for i, env := range msg.Environments {
			enviroments[i] = table.Row([]string{env.Name, env.Status, env.BranchMapping})
		}
		m.table.SetRows(enviroments)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := lipgloss.NewStyle().Bold(true).Render("Environments")
	return lipgloss.JoinVertical(lipgloss.Center, header, m.table.View())
}
