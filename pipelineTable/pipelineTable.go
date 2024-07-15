package pipelineTable

import (
	"fmt"
	"time"

	appTable "github.com/FredrikMWold/radix-tui/applicationTable"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var columns = []table.Column{
			{Title: "name", Width: (msg.Width - 44) / 5},
			{Title: "branch", Width: (msg.Width - 44) / 5},
			{Title: "pipeline", Width: (msg.Width - 44) / 5},
			{Title: "status", Width: (msg.Width - 44) / 5},
			{Title: "created", Width: (msg.Width - 44) / 5},
		}
		m.table.SetHeight(msg.Height - 5)
		m.table.SetWidth(msg.Width - 34)
		m.table.SetColumns(columns)

	case appTable.SelectedApplication:
		m.isLoadingApplication = true

	case appTable.Application:
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
		m.table.SetRows(jobs)
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var pipelineJobs string
	if m.isLoadingApplication && m.table.Rows() == nil {
		pipelineJobs = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12")).
			Height(m.table.Height() + 2).
			Width(m.table.Width()).
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Render(fmt.Sprintf("Loading pipeline jobs " + m.spinner.View()))
	} else {
		pipelineJobs = baseStyle.Render(m.table.View())
	}

	if m.table.Rows() != nil && m.isLoadingApplication {
		pipelineJobs = lipgloss.JoinVertical(lipgloss.Center, "Pipeline jobs "+m.spinner.View(), pipelineJobs)
	} else {
		pipelineJobs = lipgloss.JoinVertical(lipgloss.Center, "Pipeline jobs", pipelineJobs)
	}

	return pipelineJobs
}
