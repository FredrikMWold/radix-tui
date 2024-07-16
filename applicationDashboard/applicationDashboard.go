package appllicationDashboard

import (
	"github.com/FredrikMWold/radix-tui/applicationTable"
	"github.com/FredrikMWold/radix-tui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.applicationsTable.Init(), m.pipelineTable.Init(), m.enviromentTable.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

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

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.ChangeFocus()
		case "ctrl+r":
			m.applicationsTable, cmd = m.applicationsTable.Update(msg)
			return m, cmd
		case "enter":
			m.focused = pipeline
		case "esc":
			m.focused = application
		}
	case tea.WindowSizeMsg, applicationTable.SelectedApplication, applicationTable.Application, spinner.TickMsg, applicationTable.UpdateApplicationDataTick:
		var appCmds, pipeCmds, envCmds tea.Cmd
		m.applicationsTable, appCmds = m.applicationsTable.Update(msg)
		m.pipelineTable, pipeCmds = m.pipelineTable.Update(msg)
		m.enviromentTable, envCmds = m.enviromentTable.Update(msg)
		cmds = append(cmds, appCmds, pipeCmds, envCmds)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Top,
			styles.SectionContainer(m.focused == application).Render(m.applicationsTable.View()),
			m.enviromentTable.View(),
		),
		styles.SectionContainer(m.focused == pipeline).Render(m.pipelineTable.View()),
	)
}

func (m *Model) ChangeFocus() {
	m.focused = (m.focused + 1) % (pipeline + 1)
}
