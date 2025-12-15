package appllicationdashboard

import (
	"fmt"

	"github.com/FredrikMWold/radix-tui/applicationDashboard/commands"
	contextswitcher "github.com/FredrikMWold/radix-tui/applicationDashboard/contextSwitcher"
	"github.com/FredrikMWold/radix-tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.applicationsTable.Init(),
		m.pipelineTable.Init(),
		m.enviromentTable.Init(),
		commands.CheckAuth(),
		// Emit cached applications immediately if available for fast startup
		commands.LoadCachedApplications(),
		m.spinner.Tick, // spinner will be visible only when we choose to render it
		getContext,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case commands.AuthWaiting:
		// Show waiting for auth (no spinner), and trigger interactive login
		m.hasAuthRedirect = true
		return m, commands.LoginInteractive()

	case commands.LoggedIn, commands.AuthOK:
		// After login or when already authenticated, start loading applications
		m.hasAuthRedirect = false
		return m, commands.GetApplications

	case contextswitcher.ContextSelected:
		// Context was selected, save it and reload
		m.focused = m.previousFocused
		return m, commands.SetContext(string(msg))

	case contextswitcher.ContextCancelled:
		// Context switcher was cancelled, restore previous focus
		m.focused = m.previousFocused
		return m, nil

	case commands.ContextChanged:
		// Context changed successfully, reload everything
		m.context = string(msg)
		m.applications = []string{}
		m.application = commands.Application{}
		m.isLoadingApplication = false
		m.loadError = ""
		m.appsLoaded = false
		m.focused = application
		m.keys = ApplicationTableKeys
		// Reinitialize context switcher with new context
		m.contextSwitcher = contextswitcher.New(m.context)
		return m, tea.Batch(commands.GetApplications, m.applicationsTable.Init())

	case tea.KeyMsg:
		// Handle context switch key globally (except when in context switcher)
		if msg.String() == "ctrl+s" && m.focused != contextSwitch {
			m.previousFocused = m.focused
			m.focused = contextSwitch
			m.contextSwitcher = contextswitcher.New(m.context)
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			if m.focused == contextSwitch {
				m.focused = m.previousFocused
				return m, nil
			}
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
			if m.focused == contextSwitch {
				m.focused = m.previousFocused
				return m, nil
			}
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
		// Accept applications if:
		// 1. Context matches current context, OR
		// 2. We don't have a context yet (initial load - accept cache immediately)
		if msg.Context == m.context || m.context == "" {
			m.applications = msg.Apps
			m.hasAuthRedirect = false
			m.loadError = ""
			m.appsLoaded = true
			// If we didn't have context yet, set it from the message
			if m.context == "" {
				m.context = msg.Context
			}
		}

	case commands.ApplicationsError:
		// Only show error if it's for the current context (or no context set yet)
		if msg.Context == m.context || m.context == "" {
			m.hasAuthRedirect = false
			m.loadError = fmt.Sprintf("Failed to load applications for context '%s':\n%s\n\nPress ctrl+s to switch context or q to quit", msg.Context, msg.Err.Error())
			m.appsLoaded = true
			if m.context == "" {
				m.context = msg.Context
			}
		}

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
	if m.focused == contextSwitch {
		m.contextSwitcher, pageCmd = m.contextSwitcher.Update(msg)
	}
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
	// Before we get the first WindowSizeMsg, width/height can be 0; render a minimal view instead of blank.
	if m.width == 0 || m.height == 0 {
		if m.hasAuthRedirect {
			return "Waiting for authentication"
		}
		return "Loading applications " + m.spinner.View()
	}
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
	if !m.appsLoaded {
		if m.hasAuthRedirect {
			// During redirect, show text without spinner
			return lipgloss.NewStyle().
				Height(m.height).
				Width(m.width).
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				Render("Waiting for authentication")
		}
		if m.loadError != "" {
			// Show error message with option to switch context
			errorStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Bold(true)
			content := errorStyle.Render(m.loadError)
			// If context switcher is open, show it as overlay
			if m.focused == contextSwitch {
				return m.renderContextSwitcherOverlay(content)
			}
			return lipgloss.NewStyle().
				Height(m.height).
				Width(m.width).
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				Render(content)
		}
		// Otherwise, show normal loading with spinner
		return lipgloss.NewStyle().
			Height(m.height).
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Loading applications " + m.spinner.View())
	}

	// Handle empty apps list (loaded successfully but no apps)
	if len(m.applications) == 0 {
		content := fmt.Sprintf("No applications found in context '%s'\n\nPress ctrl+s to switch context or q to quit", m.context)
		// If context switcher is open, show it as overlay
		if m.focused == contextSwitch {
			return m.renderContextSwitcherOverlay(content)
		}
		return lipgloss.NewStyle().
			Height(m.height).
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(content)
	}

	// Build the main view
	mainView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Top,
			styles.SectionContainer(m.focused == application).Render(m.applicationsTable.View()),
			m.getEnvironemntTableView(),
			styles.SectionContainer(false).Width(30).Align(lipgloss.Center).Render("Context: "+m.context),
		),
		m.getActivePageView(),
	) + "\n" + m.help.View(m.keys)

	// If context switcher is active, show it as a modal overlay
	if m.focused == contextSwitch {
		return m.renderContextSwitcherOverlay(mainView)
	}

	return mainView
}

func (m Model) renderContextSwitcherOverlay(_ string) string {
	// Create the modal
	modal := m.contextSwitcher.View()

	// Create a dimmed background effect by placing the modal on top
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
}
