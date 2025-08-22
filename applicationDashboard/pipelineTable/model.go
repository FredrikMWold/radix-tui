package pipelinetable

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/FredrikMWold/radix-tui/applicationDashboard/commands"
	"github.com/FredrikMWold/radix-tui/styles"
	"github.com/charmbracelet/bubbles/table"
)

func New() Model {
	return Model{
		table: table.New(
			table.WithFocused(true),
			table.WithStyles(styles.TableStyles()),
		),
	}
}

type Model struct {
	table       table.Model
	application commands.Application
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
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
