package applicationTable

import (
	"fmt"

	"github.com/FredrikMWold/radix-tui/styles"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(tick(), m.spinner.Tick, getApplications)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "ctrl+r":
			if len(m.table.SelectedRow()) > 0 {
				m.selectedApp = m.table.SelectedRow()[0]
				if m.selectedApp == "" {
					cmds = append(cmds, tick())
				}
				cmds = append(cmds, getApplicationData(m.selectedApp), selectApplication(m.selectedApp))
				return m, tea.Batch(cmds...)
			}
		}
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 16)

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
	var section string
	var table string
	if m.isLoadingApplications {
		table = styles.LoadingSpinnerContainer(m.table.Height()+3, 30).
			Render(fmt.Sprintf("Loading applications " + m.spinner.View()))
	} else {
		table = m.table.View()
	}
	if m.table.Rows() != nil && m.isLoadingApplications {
		section = lipgloss.JoinVertical(lipgloss.Center, "Applications "+m.spinner.View(), table)
	} else {
		section = lipgloss.JoinVertical(lipgloss.Center, "Applications", table)
	}
	return section
}
