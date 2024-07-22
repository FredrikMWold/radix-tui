package appllicationdashboard

import (
	"github.com/FredrikMWold/radix-tui/commands"
	"github.com/FredrikMWold/radix-tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.applicationsTable.Init(),
		m.pipelineTable.Init(),
		m.enviromentTable.Init(),
		m.pipelineForm.Init(),
		commands.GetApplications,
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+r":
			if m.application.Name != "" {
				m.isLoadingApplication = true
				return m, commands.GetApplicationData(m.application.Name)
			}
		case "ctrl+n":
			if m.focused == pipeline {
				m.focused = form
				m.keys = BuildDeployFormKeys
			}
		case "esc":
			if m.focused == pipeline {
				m.focused = application
				m.keys = ApplicationTableKeys
			}
			if m.focused != application {
				m.focused = pipeline
				m.keys = PipelineTableKeys
			}
		}

	case commands.Applications:
		m.applications = msg

	case commands.SelectedApplication:
		m.isLoadingApplication = true
		m.focused = pipeline
		m.keys = PipelineTableKeys
		return m, commands.GetApplicationData(string(msg))

	case tea.WindowSizeMsg, commands.Application:
		var appCmds, pipeCmds, envCmds, formCmds tea.Cmd
		if _, ok := msg.(commands.Application); ok {
			m.isLoadingApplication = false
			m.application = msg.(commands.Application)
		}
		if _, ok := msg.(tea.WindowSizeMsg); ok {
			m.height = msg.(tea.WindowSizeMsg).Height
			m.width = msg.(tea.WindowSizeMsg).Width
			m.help.Width = msg.(tea.WindowSizeMsg).Width
		}
		m.applicationsTable, appCmds = m.applicationsTable.Update(msg)
		m.pipelineTable, pipeCmds = m.pipelineTable.Update(msg)
		m.enviromentTable, envCmds = m.enviromentTable.Update(msg)
		m.pipelineForm, formCmds = m.pipelineForm.Update(msg)
		cmds = append(cmds, appCmds, pipeCmds, envCmds, formCmds)
	}

	if m.focused == form {
		var formCmd tea.Cmd
		m.pipelineForm, formCmd = m.pipelineForm.Update(msg)
		cmds = append(cmds, formCmd)
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

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if len(m.applications) == 0 {
		return lipgloss.NewStyle().
			Height(m.height).
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Loading applications " + m.spinner.View())
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Top,
			styles.SectionContainer(m.focused == application).Render(m.applicationsTable.View()),
			m.getEnvironemntTableView(),
		),
		m.getActivePageView(),
	) + "\n" + m.help.View(m.keys)
}