package pipelinetable

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
			m.openJobInBrowser()
		}
	case tea.WindowSizeMsg:
		var columns = []table.Column{
			{Title: "Triggered by", Width: (msg.Width - 44) / 5},
			{Title: "Environment", Width: (msg.Width - 44) / 5},
			{Title: "pipeline", Width: (msg.Width - 44) / 5},
			{Title: "status", Width: (msg.Width - 44) / 5},
			{Title: "created", Width: (msg.Width - 44) / 5},
		}
		m.table.SetHeight(msg.Height - 7)
		m.table.SetWidth(msg.Width - 34)
		m.table.SetColumns(columns)

	case commands.Application:
		m.loadApplication(msg)
	}

	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	cmds = append(cmds, tableCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := lipgloss.NewStyle().Bold(true).Render("Pipeline jobs")
	return lipgloss.JoinVertical(lipgloss.Center, header, m.table.View())
}
