package applyconfig

import (
	"github.com/FredrikMWold/radix-tui/applicationDashboard/commands"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case commands.Application:
		m.selectedApplication = msg.Name
		m.form = NewForm()
	}

	if m.form.State == huh.StateCompleted && m.form.GetBool("applyConfig") {
		m.form.State = huh.StateAborted
		return m, ApplyConfig(m.selectedApplication)
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	formHeader := lipgloss.NewStyle().Padding(0, 0, 2, 0).Bold(true).Render("Apply config")
	form := lipgloss.NewStyle().MaxWidth(26).Render(m.form.View())
	return lipgloss.PlaceHorizontal(m.width-34, lipgloss.Center, formHeader+"\n"+form)
}
