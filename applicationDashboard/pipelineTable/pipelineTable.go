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
		// Compute right-pane width and distribute across columns to fill fully
		w := msg.Width - 34
		m.table.SetHeight(msg.Height - 7)
		m.table.SetWidth(w)
		// Bubble table adds cell padding (default is 1 left + 1 right per column)
		const cols = 5
		const cellPad = 2 // total horizontal padding per cell (left+right)
		contentW := w - (cols * cellPad)
		if contentW < cols { // ensure non-negative widths
			contentW = cols
		}
		// Proportional widths based on content width; last column takes remainder
		w1 := contentW * 2 / 7               // Triggered by
		w2 := contentW * 1 / 7               // Environment
		w3 := contentW * 1 / 7               // pipeline
		w4 := contentW * 1 / 7               // status
		w5 := contentW - (w1 + w2 + w3 + w4) // created gets the rest
		if w1 < 1 {
			w1 = 1
		}
		if w2 < 1 {
			w2 = 1
		}
		if w3 < 1 {
			w3 = 1
		}
		if w4 < 1 {
			w4 = 1
		}
		if w5 < 1 {
			w5 = 1
		}
		var columns = []table.Column{
			{Title: "Triggered by", Width: w1},
			{Title: "Environment", Width: w2},
			{Title: "pipeline", Width: w3},
			{Title: "status", Width: w4},
			{Title: "created", Width: w5},
		}
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
		spin = " " + m.spinner.View()
	}
	// Reserve a wider spinner slot so centering is stable across frames
	spinBox := lipgloss.NewStyle().Width(4).Render(spin)
	headerLine := lipgloss.JoinHorizontal(lipgloss.Left, titleText, spinBox)
	// Make header exactly as wide as the table and center the content for symmetric padding
	header := lipgloss.NewStyle().
		Width(lipgloss.Width(m.table.View())).
		AlignHorizontal(lipgloss.Center).
		Render(headerLine)
	return lipgloss.JoinVertical(lipgloss.Left, header, m.table.View())
}

// pollTick is an internal message used to schedule periodic refreshes.
type pollTick struct{ seq int }
