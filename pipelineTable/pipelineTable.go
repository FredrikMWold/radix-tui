package pipelineTable

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/FredrikMWold/radix-tui/applicationTable"
	"github.com/FredrikMWold/radix-tui/styles"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.openJobInBrowser()
		}
	case tea.WindowSizeMsg:
		var columns = []table.Column{
			{Title: "Triggered by", Width: (msg.Width - 44) / 5},
			{Title: "Environment", Width: (msg.Width - 44) / 5},
			{Title: "pipeline", Width: (msg.Width - 44) / 5},
			{Title: "status", Width: (msg.Width - 44) / 5},
			{Title: "created", Width: (msg.Width - 44) / 5},
		}
		m.table.SetHeight(msg.Height - 6)
		m.table.SetWidth(msg.Width - 34)
		m.table.SetColumns(columns)

	case applicationTable.SelectedApplication:
		m.isLoadingApplication = true
		m.selectedApplication = string(msg)

	case applicationTable.UpdateApplicationDataTick:
		m.isLoadingApplication = true

	case applicationTable.Application:
		m.isLoadingApplication = false
		m.loadApplication(msg)

	}

	var tableCmd, spinnerCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	cmds = append(cmds, tableCmd)

	m.spinner, spinnerCmd = m.spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var section string
	var table string
	if m.isLoadingApplication && m.table.Rows() == nil {
		table = styles.LoadingSpinnerContainer(m.table.Height()+3, m.table.Width()).
			Render(fmt.Sprintf("Loading pipeline jobs " + m.spinner.View()))
	} else {
		table = m.table.View()
	}
	if m.table.Rows() != nil && m.isLoadingApplication {
		section = lipgloss.JoinVertical(lipgloss.Center, m.selectedApplication+" "+m.spinner.View(), table)
	} else {
		section = lipgloss.JoinVertical(lipgloss.Center, m.selectedApplication, table)
	}

	return section
}

func (m *Model) loadApplication(application applicationTable.Application) {
	rows := make([]table.Row, len(application.Jobs))
	jobs := make([]string, len(application.Jobs))
	for i, job := range application.Jobs {
		parsedTime, err := time.Parse(time.RFC3339, job.Created)
		if err != nil {
			fmt.Println(err)
		}
		job.Created = parsedTime.Format("02.01.2006 15:04:05")
		environment := "No environment"
		if len(job.Environments) > 0 {
			environment = job.Environments[0]
		}
		rows[i] = table.Row([]string{job.TriggeredBy, environment, job.Pipeline, job.Status, job.Created})
		jobs[i] = job.Name
	}
	m.table.SetRows(rows)
	m.jobs = jobs
}

func (m Model) openJobInBrowser() {
	if len(m.jobs) == 0 {
		return
	}
	tableCursor := m.table.Cursor()
	url := fmt.Sprintf("https://console.radix.equinor.com/applications/%s/jobs/view/%s", m.selectedApplication, m.jobs[tableCursor])
	exec.Command("open", url).Start()
}
