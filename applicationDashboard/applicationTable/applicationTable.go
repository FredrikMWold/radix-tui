package applicationtable

import (
	"github.com/FredrikMWold/radix-tui/applicationDashboard/commands"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if len(m.table.SelectedRow()) > 0 {
				return m, commands.SelectApplication(m.table.SelectedRow()[0])
			}
		}
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 21)

	case commands.Applications:
		m.isLoadingApplications = false
		rows := make([]table.Row, len(msg.Apps))
		for i, app := range msg.Apps {
			rows[i] = table.Row([]string{app})
		}
		m.table.SetRows(rows)
	}

	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	cmds = append(cmds, tableCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := lipgloss.NewStyle().Bold(true).Render("Applications")
	return lipgloss.JoinVertical(lipgloss.Center, header, m.table.View())
}
