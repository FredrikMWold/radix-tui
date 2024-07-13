package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Application struct {
	Jobs         []Job         `json:"jobs"`
	Environments []Environment `json:"environments"`
}

type Job struct {
	AppName  string `json:"appName"`
	Branch   string `json:"branch"`
	Pipeline string `json:"pipeline"`
	Status   string `json:"status"`
	Created  string `json:"created"`
}

type Environment struct {
	BranchMapping string `json:"branchMapping"`
	Name          string `json:"name"`
	Status        string `json:"status"`
}

type model struct {
	applicationstable     table.Model
	pipelineTable         table.Model
	enviromentTable       table.Model
	spinner               spinner.Model
	selectedApp           string
	isLoadingApplications bool
	isLoadingApplication  bool
	height                int
}

func initialModel() model {
	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	infoTableStyles := table.Styles{
		Selected: lipgloss.NewStyle().Bold(false),
		Header: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false).Padding(0, 1),
		Cell: lipgloss.NewStyle().Padding(0, 1),
	}

	applicationColumns := []table.Column{
		{Title: "name", Width: 28},
	}
	applicationTable := table.New(
		table.WithColumns(applicationColumns),
		table.WithFocused(true),
		table.WithHeight(8),
		table.WithStyles(tableStyles),
	)

	pipelineTable := table.New(
		table.WithFocused(false),
		table.WithStyles(infoTableStyles),
	)

	enviromentTable := table.New(
		table.WithFocused(false),
		table.WithHeight(4),
		table.WithStyles(infoTableStyles),
		table.WithColumns([]table.Column{
			{Title: "name", Width: 7},
			{Title: "status", Width: 10},
			{Title: "branch", Width: 7},
		}),
	)

	spiner := spinner.New()
	spiner.Spinner = spinner.Meter
	spiner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		applicationstable:     applicationTable,
		pipelineTable:         pipelineTable,
		enviromentTable:       enviromentTable,
		spinner:               spiner,
		isLoadingApplications: true,
		isLoadingApplication:  false,
		height:                12,
	}
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("12"))

type applications []string

func getApplications() tea.Msg {
	result := exec.Command("rx", "get", "application")
	output, err := result.Output()
	if err != nil {
		fmt.Println(err)
	}
	trimmed := strings.TrimSpace(string(output))
	application_list := strings.Split(trimmed, "\n")
	return applications(application_list)
}

func getApplicationData(application string) tea.Cmd {
	return func() tea.Msg {
		result := exec.Command("rx", "get", "application", "-a", application)
		output, err := result.Output()
		if err != nil {
			fmt.Println(err)
		}
		var application Application
		err = json.Unmarshal(output, &application)
		if err != nil {
			fmt.Println(err)
		}
		return application
	}
}

type updateApplicationDataTick time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second*10, func(t time.Time) tea.Msg {
		return updateApplicationDataTick(t)
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(getApplications, m.spinner.Tick, tick())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	case table.KeyMap:
		fmt.Println(msg)

	case tea.WindowSizeMsg:
		var columns = []table.Column{
			{Title: "name", Width: (msg.Width - 44) / 5},
			{Title: "branch", Width: (msg.Width - 44) / 5},
			{Title: "pipeline", Width: (msg.Width - 44) / 5},
			{Title: "status", Width: (msg.Width - 44) / 5},
			{Title: "created", Width: (msg.Width - 44) / 5},
		}
		m.pipelineTable.SetHeight(msg.Height - 5)
		m.pipelineTable.SetWidth(msg.Width - 34)
		m.pipelineTable.SetColumns(columns)

		m.applicationstable.SetHeight(msg.Height - 14)
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter", "ctrl+r":
			m.selectedApp = m.applicationstable.SelectedRow()[0]
			m.isLoadingApplication = true
			return m, getApplicationData(m.selectedApp)
		case "left", "right":
			if m.applicationstable.Focused() {
				m.applicationstable.Blur()
				m.pipelineTable.Focus()
			} else {
				m.pipelineTable.Blur()
				m.applicationstable.Focus()
			}
		}

	case updateApplicationDataTick:
		if !m.isLoadingApplication && m.selectedApp != "" {
			m.isLoadingApplication = true
			return m, tea.Batch(getApplicationData(m.selectedApp), tick())
		}
		cmds = append(cmds, tick())

	case applications:
		rows := make([]table.Row, len(msg))
		for i, app := range msg {
			rows[i] = table.Row([]string{app})
		}
		m.applicationstable.SetRows(rows)
		m.isLoadingApplications = false

	case Application:
		m.isLoadingApplication = false
		jobs := make([]table.Row, len(msg.Jobs))
		for i, job := range msg.Jobs {
			parsedTime, err := time.Parse(time.RFC3339, job.Created)
			if err != nil {
				fmt.Println(err)
			}
			job.Created = parsedTime.Format("02.01.2006 15:04:05")
			jobs[i] = table.Row([]string{job.AppName, job.Branch, job.Pipeline, job.Status, job.Created})
		}

		enviroments := make([]table.Row, len(msg.Environments))
		for i, env := range msg.Environments {
			enviroments[i] = table.Row([]string{env.Name, env.Status, env.BranchMapping})
		}
		m.enviromentTable.SetRows(enviroments)
		m.pipelineTable.SetRows(jobs)

	}

	var tableCmd tea.Cmd
	m.applicationstable, tableCmd = m.applicationstable.Update(msg)
	cmds = append(cmds, tableCmd)

	var pipelineTableCmd tea.Cmd
	m.pipelineTable, pipelineTableCmd = m.pipelineTable.Update(msg)
	cmds = append(cmds, pipelineTableCmd)

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var applicationPicker string
	if m.isLoadingApplications {
		applicationPicker = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(30).
			Height(m.applicationstable.Height() + 2).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			BorderForeground(lipgloss.Color("12")).
			Render(fmt.Sprintf("Loading applications " + m.spinner.View()))
	} else {
		applicationPicker = baseStyle.Render(m.applicationstable.View())
	}

	var pipelineJobs string
	if m.isLoadingApplication && m.pipelineTable.Rows() == nil {
		pipelineJobs = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12")).
			Height(m.pipelineTable.Height() + 2).
			Width(m.pipelineTable.Width()).
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Render(fmt.Sprintf("Loading pipeline jobs " + m.spinner.View()))
	} else {
		pipelineJobs = baseStyle.Render(m.pipelineTable.View())
	}

	if m.pipelineTable.Rows() != nil && m.isLoadingApplication {
		pipelineJobs = lipgloss.JoinVertical(lipgloss.Center, "Pipeline jobs "+m.spinner.View(), pipelineJobs)
	} else {
		pipelineJobs = lipgloss.JoinVertical(lipgloss.Center, "Pipeline jobs", pipelineJobs)
	}

	var enviromentTable string
	if m.isLoadingApplication && m.enviromentTable.Rows() == nil {
		enviromentTable = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Height(m.enviromentTable.Height() + 2).
			Width(30).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			BorderForeground(lipgloss.Color("12")).
			Render(fmt.Sprintf("Loading applications " + m.spinner.View()))
	} else {
		enviromentTable = baseStyle.Render(m.enviromentTable.View())
	}
	if m.enviromentTable.Rows() != nil && m.isLoadingApplication {
		enviromentTable = lipgloss.JoinVertical(lipgloss.Center, "Environments "+m.spinner.View(), enviromentTable)
	} else {
		enviromentTable = lipgloss.JoinVertical(lipgloss.Center, "Environments", enviromentTable)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Top,
			lipgloss.JoinVertical(lipgloss.Center, "Applications", applicationPicker), enviromentTable),
		pipelineJobs)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
