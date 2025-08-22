package pipelinetable

import (
	"time"

	"github.com/FredrikMWold/radix-tui/applicationDashboard/commands"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd { return nil }

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
			{Title: "Pipeline", Width: (msg.Width - 44) / 5},
			{Title: "Status", Width: (msg.Width - 44) / 5},
			{Title: "Created", Width: (msg.Width - 44) / 5},
		}
		m.table.SetHeight(msg.Height - 7)
		m.table.SetWidth(msg.Width - 34)
		m.table.SetColumns(columns)

	case commands.Application:
		m.loadApplication(msg)
		m.refreshing = false
		// Schedule a poll in 15s tied to current sequence
		seq := m.pollSeq
		cmds = append(cmds,
			tea.Tick(15*time.Second, func(t time.Time) tea.Msg { return pollTick{seq: seq} }),
		)

	case pollTick:
		// Only handle the latest scheduled tick
		if msg.seq != m.pollSeq {
			break
		}
		// Only poll when we have a selected application and not already refreshing
		if m.application.Name != "" && !m.refreshing {
			m.refreshing = true
			// start spinner ticking only during refresh
			cmds = append(cmds, m.spinner.Tick, commands.GetApplicationData(m.application.Name))
		}
		// Bump seq so any older scheduled ticks are ignored until next schedule
		m.pollSeq++
	}

	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	cmds = append(cmds, tableCmd)

	// Update spinner only when refreshing
	if m.refreshing {
		var spinCmd tea.Cmd
		m.spinner, spinCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// Title with a fixed-width spinner slot to avoid layout shift
	titleText := lipgloss.NewStyle().Bold(true).Render("Pipeline jobs")
	spin := ""
	if m.refreshing {
		spin = m.spinner.View()
	}
	// Reserve two columns for spinner output so the title width stays constant
	spinBox := lipgloss.NewStyle().Width(3).Render(spin)
	header := lipgloss.JoinHorizontal(lipgloss.Left, titleText, " ", spinBox)
	return lipgloss.JoinVertical(lipgloss.Center, header, m.table.View())
}

// pollTick is an internal message used to schedule periodic refreshes.
type pollTick struct{ seq int }
