package pipelinetable

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/FredrikMWold/radix-tui/applicationDashboard/commands"
	"github.com/FredrikMWold/radix-tui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func New() Model {
	sp := spinner.New()
	sp.Spinner = spinner.Meter
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Model{
		table: table.New(
			table.WithFocused(true),
			table.WithStyles(styles.TableStyles()),
		),
		spinner: sp,
	}
}

type Model struct {
	table       table.Model
	application commands.Application
	spinner     spinner.Model
	refreshing  bool
	pollSeq     int
}

func (m *Model) loadApplication(application commands.Application) {
	m.application = application
	m.refreshing = false
	rows := make([]table.Row, len(application.Jobs))
	for i, job := range application.Jobs {
		if t, ok := parseCreated(job.Created); ok {
			// Render in local time to avoid confusion with UTC offsets
			job.Created = t.In(time.Local).Format("02.01.2006 15:04:05")
		}
		environment := "No environment"
		if len(job.Environments) > 0 {
			environment = job.Environments[0]
		}
		rows[i] = table.Row([]string{job.TriggeredBy, environment, job.Pipeline, job.Status, job.Created})
	}
	m.table.SetRows(rows)
}

// parseCreated tries common layouts from the API and returns a parsed time if successful.
func parseCreated(s string) (time.Time, bool) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		// Additional explicit forms sometimes used by generators
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05.999999Z07:00",
		"2006-01-02T15:04:05Z07:00",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	// Keep original string on failure to avoid showing year 0001
	fmt.Println("failed to parse job.Created:", s)
	return time.Time{}, false
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
