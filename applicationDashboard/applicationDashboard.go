package appllicationdashboard

import (
	"github.com/FredrikMWold/radix-tui/applicationDashboard/commands"
	"github.com/FredrikMWold/radix-tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.applicationsTable.Init(),
		m.pipelineTable.Init(),
		m.enviromentTable.Init(),
		commands.GetApplications,
		m.spinner.Tick,
		getContext,
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
				m.focused = buildAndDeploy
				m.keys = BuildDeployFormKeys
			}
		case "ctrl+a":
			if m.focused == pipeline {
				m.focused = applyConfig
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

	case Context:
		m.context = string(msg)
	case commands.SelectedApplication:
		m.isLoadingApplication = true
		m.focused = pipeline
		m.keys = PipelineTableKeys
		return m, commands.GetApplicationData(string(msg))

	case tea.WindowSizeMsg, commands.Application:
		if _, ok := msg.(commands.Application); ok {
			m.isLoadingApplication = false
			m.application = msg.(commands.Application)
		}
		if _, ok := msg.(tea.WindowSizeMsg); ok {
			m.height = msg.(tea.WindowSizeMsg).Height
			m.width = msg.(tea.WindowSizeMsg).Width
			m.help.Width = msg.(tea.WindowSizeMsg).Width
		}

		m.applicationsTable, _ = m.applicationsTable.Update(msg)
		m.pipelineTable, _ = m.pipelineTable.Update(msg)
		m.enviromentTable, _ = m.enviromentTable.Update(msg)
		m.buildAndDeploy, _ = m.buildAndDeploy.Update(msg)
		m.applyConfig, _ = m.applyConfig.Update(msg)
	}

	var pageCmd tea.Cmd
	if m.focused == buildAndDeploy {
		m.buildAndDeploy, pageCmd = m.buildAndDeploy.Update(msg)
	}
	if m.focused == application {
		m.applicationsTable, pageCmd = m.applicationsTable.Update(msg)
	}
	if m.focused == pipeline {
		m.pipelineTable, pageCmd = m.pipelineTable.Update(msg)
	}
	if m.focused == applyConfig {
		m.applyConfig, pageCmd = m.applyConfig.Update(msg)
	}
	cmds = append(cmds, pageCmd)

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width <= 110 || m.height <= 25 {
		return lipgloss.NewStyle().
			Height(m.height).
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					"Terminal window is too small \n",
					"Please resize the terminal window to at least 110x25"),
			)

	}
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
			styles.SectionContainer(false).Width(30).Align(lipgloss.Center).Render("Context: "+m.context),
		),
		m.getActivePageView(),
	) + "\n" + m.help.View(m.keys)
}
