package main

import (
	"fmt"
	"os"

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
		}
	}
	var applicationsTableCmd tea.Cmd
	m.applicationsTable, applicationsTableCmd = m.applicationsTable.Update(msg)
	cmds = append(cmds, applicationsTableCmd)

	var pipelineTableCmd tea.Cmd
	m.pipelineTable, pipelineTableCmd = m.pipelineTable.Update(msg)
	cmds = append(cmds, pipelineTableCmd)

	var envComd tea.Cmd
	m.enviromentTable, envComd = m.enviromentTable.Update(msg)
	cmds = append(cmds, envComd)

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

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
