package environmentTable

import (
	"fmt"

	"github.com/FredrikMWold/radix-tui/applicationTable"
	"github.com/FredrikMWold/radix-tui/styles"
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

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var section string
	var table string
	if m.isLoadingApplication && m.table.Rows() == nil {
		table = styles.LoadingSpinnerContainer(m.table.Height()+3, 30).
			Render(fmt.Sprintf("Loading applications " + m.spinner.View()))
	} else {
		table = m.table.View()
	}
	if m.table.Rows() != nil && m.isLoadingApplication {
		section = lipgloss.JoinVertical(lipgloss.Center, "Environments "+m.spinner.View(), table)
	} else {
		section = lipgloss.JoinVertical(lipgloss.Center, "Environments", table)
	}
	return styles.SectionContainer(false).Render(section)
}
