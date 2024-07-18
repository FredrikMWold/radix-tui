package pipelineForm

import (
	"github.com/FredrikMWold/radix-tui/commands"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - 34
	case commands.Application:
		var environments []string
		var branches []string
		for _, env := range msg.Environments {
			environments = append(environments, env.Name)
			if env.BranchMapping != "" {
				branches = append(branches, env.BranchMapping)
			}
		}
		m.form = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Key("environment").
					Options(huh.NewOptions(environments...)...).
					Title("Environment").
					Description("Select the environment you want to deploy to").
					WithTheme(huh.ThemeCatppuccin()),
				huh.NewSelect[string]().
					Key("branch").
					Options(huh.NewOptions(branches...)...).
					Title("Pipeline Type").
					Description("Select the type of pipeline you want to create").
					WithTheme(huh.ThemeCatppuccin()),
			),
		).WithWidth(m.width).WithShowHelp(false)
		m.SelectedApplication = msg.Name
		return m, m.form.Init()
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var formHeader = lipgloss.NewStyle().Padding(0, 0, 1, 0).Bold(true).Render("Build and deploy")
	return lipgloss.JoinVertical(lipgloss.Center, formHeader, m.form.View())
}
