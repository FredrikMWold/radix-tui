package main

import (
	"fmt"
	"os"

	appllicationDashboard "github.com/FredrikMWold/radix-tui/applicationDashboard"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := appllicationDashboard.New()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
