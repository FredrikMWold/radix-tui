package buildanddeploy

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
		var options []string
		for _, env := range msg.Environments {
			if env.BranchMapping != "" {
				key := env.BranchMapping + " -> " + env.Name
				m.branchMapping[key] = env.BranchMapping
				options = append(options, key)
			}
		}
		m.form = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Key("environment").
					Options(huh.NewOptions(options...)...).
					Title("Environment").
					Description("Select the environemnt you want to deploy").
					WithTheme(huh.ThemeCatppuccin()),
			),
		)
		m.SelectedApplication = msg.Name
		return m, nil
	}

	if m.form.State == huh.StateCompleted {
		m.form.State = huh.StateAborted
		return m, BuildAndDeploy(m.SelectedApplication, m.branchMapping[m.form.GetString("environment")])
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	formHeader := lipgloss.NewStyle().Padding(0, 0, 2, 0).Bold(true).Render("Build and deploy")
	form := lipgloss.NewStyle().MaxWidth(50).Render(m.form.View())
	return lipgloss.PlaceHorizontal(m.width-34, lipgloss.Center, formHeader+"\n"+form)
}
