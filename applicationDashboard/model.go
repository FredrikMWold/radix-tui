package appllicationdashboard

import (
	applicationtable "github.com/FredrikMWold/radix-tui/applicationDashboard/applicationTable"
	"github.com/FredrikMWold/radix-tui/applicationDashboard/commands"
	contextswitcher "github.com/FredrikMWold/radix-tui/applicationDashboard/contextSwitcher"
	environmenttable "github.com/FredrikMWold/radix-tui/applicationDashboard/environmentTable"
	applyconfig "github.com/FredrikMWold/radix-tui/applicationDashboard/pipelineForms/applyConfig"
	buildanddeploy "github.com/FredrikMWold/radix-tui/applicationDashboard/pipelineForms/buildAndDeploy"
	pipelinetable "github.com/FredrikMWold/radix-tui/applicationDashboard/pipelineTable"
	"github.com/FredrikMWold/radix-tui/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Focused int

const (
	application Focused = iota
	pipeline
	buildAndDeploy
	applyConfig
	contextSwitch
)

type Model struct {
	applicationsTable    applicationtable.Model
	pipelineTable        pipelinetable.Model
	enviromentTable      tea.Model
	buildAndDeploy       tea.Model
	applyConfig          tea.Model
	contextSwitcher      contextswitcher.Model
	spinner              spinner.Model
	keys                 keyMap
	help                 help.Model
	context              string
	application          commands.Application
	focused              Focused
	previousFocused      Focused
	isLoadingApplication bool
	height               int
	width                int
	applications         []string
	hasAuthRedirect      bool
	loadError            string
	appsLoaded           bool
}

func New() Model {

	applicationTable := applicationtable.New()
	pipelineTable := pipelinetable.New()
	enviromentTable := environmenttable.New()
	buildAndDeploy := buildanddeploy.New()
	applyConfig := applyconfig.New()

	spiner := spinner.New()
	spiner.Spinner = spinner.Meter
	spiner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		applicationsTable:    applicationTable,
		pipelineTable:        pipelineTable,
		enviromentTable:      enviromentTable,
		buildAndDeploy:       buildAndDeploy,
		applyConfig:          applyConfig,
		help:                 help.New(),
		spinner:              spiner,
		focused:              application,
		isLoadingApplication: false,
		keys:                 ApplicationTableKeys,
	}
}

type keyMap struct {
	Enter         key.Binding
	Up            key.Binding
	Down          key.Binding
	Esc           key.Binding
	BuildDeploy   key.Binding
	ApplyConfig   key.Binding
	Refresh       key.Binding
	SwitchContext key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.BuildDeploy, k.ApplyConfig, k.Refresh, k.SwitchContext, k.Up, k.Down, k.Esc}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.BuildDeploy, k.ApplyConfig, k.Refresh, k.SwitchContext, k.Up, k.Down, k.Esc},
	}
}

func (m Model) getActivePageView() string {
	if m.isLoadingApplication {
		return styles.SectionContainer(true).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Width(m.width - 34).
			Height(m.height - 3).
			Render("Loading application data " + m.spinner.View())
	}
	if m.focused == pipeline {
		return styles.SectionContainer(true).
			Render(m.pipelineTable.View())
	}
	if m.focused == buildAndDeploy {
		return styles.SectionContainer(true).
			Width(m.width - 34).
			Height(m.height - 3).
			Render(m.buildAndDeploy.View())
	}
	if m.focused == applyConfig {
		return styles.SectionContainer(true).
			Width(m.width - 34).
			Height(m.height - 3).
			Render(m.applyConfig.View())
	}
	return styles.SectionContainer(false).
		Width(m.width - 34).
		Height(m.height - 3).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render("Select an application")
}

func (m Model) getEnvironemntTableView() string {
	if m.focused == application {
		return styles.SectionContainer(false).
			Width(30).
			Height(9).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Select an application")
	}
	if m.isLoadingApplication {
		return styles.SectionContainer(false).
			Width(30).
			Height(9).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Loading application data " + m.spinner.View())
	}
	return styles.SectionContainer(false).
		Render(m.enviromentTable.View())
}
