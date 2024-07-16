package appllicationDashboard

import (
	"github.com/FredrikMWold/radix-tui/applicationTable"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.applicationsTable.Init(), m.pipelineTable.Init(), m.enviromentTable.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.ChangeFocus()
		case "ctrl+r":
			var cmd tea.Cmd
			m.applicationsTable, cmd = m.applicationsTable.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg, applicationTable.SelectedApplication, applicationTable.Application, spinner.TickMsg, applicationTable.UpdateApplicationDataTick:
		var appCmds, pipeCmds, envCmds tea.Cmd
		m.applicationsTable, appCmds = m.applicationsTable.Update(msg)
		m.pipelineTable, pipeCmds = m.pipelineTable.Update(msg)
		m.enviromentTable, envCmds = m.enviromentTable.Update(msg)
		return m, tea.Batch(appCmds, pipeCmds, envCmds)
	}

	if m.focused == application {
		var applicationsTableCmd tea.Cmd
		m.applicationsTable, applicationsTableCmd = m.applicationsTable.Update(msg)
		cmds = append(cmds, applicationsTableCmd)
	}

	if m.focused == pipeline {
		var pipelineTableCmd tea.Cmd
		m.pipelineTable, pipelineTableCmd = m.pipelineTable.Update(msg)
		cmds = append(cmds, pipelineTableCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Top,
			m.applicationsTable.View(),
			m.enviromentTable.View(),
		),
		m.pipelineTable.View(),
	)
}

func (m *Model) ChangeFocus() {
	m.focused = (m.focused + 1) % (pipeline + 1)
	if m.focused == pipeline {
		m.pipelineTable.Focus()
		m.applicationsTable.Blur()
	} else {
		m.applicationsTable.Focus()
		m.pipelineTable.Blur()
	}
}
