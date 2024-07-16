package environmentTable

import (
	"fmt"

	"github.com/FredrikMWold/radix-tui/applicationTable"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return (m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		}
	case applicationTable.SelectedApplication:
		m.isLoadingApplication = true

	case applicationTable.Application:
		m.isLoadingApplication = false
		enviroments := make([]table.Row, len(msg.Environments))
		for i, env := range msg.Environments {
			enviroments[i] = table.Row([]string{env.Name, env.Status, env.BranchMapping})
		}
		m.table.SetRows(enviroments)

	}

	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	cmds = append(cmds, tableCmd)

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var enviromentTable string
	if m.isLoadingApplication && m.table.Rows() == nil {
		enviromentTable = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Height(m.table.Height() + 2).
			Width(30).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			BorderForeground(lipgloss.Color("12")).
			Render(fmt.Sprintf("Loading applications " + m.spinner.View()))
	} else {
		enviromentTable = baseStyle.Render(m.table.View())
	}
	if m.table.Rows() != nil && m.isLoadingApplication {
		enviromentTable = lipgloss.JoinVertical(lipgloss.Center, "Environments "+m.spinner.View(), enviromentTable)
	} else {
		enviromentTable = lipgloss.JoinVertical(lipgloss.Center, "Environments", enviromentTable)
	}
	return enviromentTable
}
