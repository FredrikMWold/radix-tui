package applicationTable

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(tick(), m.spinner.Tick, getApplications)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "ctrl+r":
			m.selectedApp = m.table.SelectedRow()[0]
			return m, tea.Batch(
				getApplicationData(m.selectedApp),
				selectApplication(m.selectedApp),
			)

		case "tab":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		}
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 14)

	case UpdateApplicationDataTick:
		if m.selectedApp != "" {
			return m, tea.Batch(
				getApplicationData(m.selectedApp),
				selectApplication(m.selectedApp),
				tick(),
			)
		}
		cmds = append(cmds, tick())
	case Applications:
		m.isLoadingApplications = false
		rows := make([]table.Row, len(msg))
		for i, app := range msg {
			rows[i] = table.Row([]string{app})
		}
		m.table.SetRows(rows)
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
	var applicationsTable string
	if m.isLoadingApplications {
		applicationsTable = loadingStyles(m.table.Height() + 2).
			Render(fmt.Sprintf("Loading applications " + m.spinner.View()))
	} else {
		applicationsTable = baseStyle.Render(m.table.View())
	}
	if m.table.Rows() != nil && m.isLoadingApplications {
		applicationsTable = lipgloss.JoinVertical(lipgloss.Center, "Applications "+m.spinner.View(), applicationsTable)
	} else {
		applicationsTable = lipgloss.JoinVertical(lipgloss.Center, "Applications", applicationsTable)
	}
	return applicationsTable
}
