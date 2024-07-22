package pipelinetable

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/FredrikMWold/radix-tui/commands"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return nil
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
		m.table.SetHeight(msg.Height - 7)
		m.table.SetWidth(msg.Width - 34)
		m.table.SetColumns(columns)

	case commands.Application:
		m.loadApplication(msg)
	}

	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	cmds = append(cmds, tableCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := lipgloss.NewStyle().Bold(true).Render("Pipeline jobs")
	return lipgloss.JoinVertical(lipgloss.Center, header, m.table.View())
}

func (m *Model) loadApplication(application commands.Application) {
	m.application = application
	rows := make([]table.Row, len(application.Jobs))
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
	}
	m.table.SetRows(rows)
}

func (m Model) openJobInBrowser() {
	if len(m.application.Jobs) == 0 {
		return
	}
	tableCursor := m.table.Cursor()
	url := fmt.Sprintf("https://console.radix.equinor.com/applications/%s/jobs/view/%s", m.application.Name, m.application.Jobs[tableCursor].Name)
	exec.Command("open", url).Start()
}
