package main

import (
	"fmt"
	"os"

	"context"

	appllicationDashboard "github.com/FredrikMWold/radix-tui/applicationDashboard"
	"github.com/FredrikMWold/radix-tui/internal/radix"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Helper mode: perform interactive login and exit
	if os.Getenv("RADIX_TUI_LOGIN_HELPER") == "1" {
		if client, err := radix.New(false); err == nil {
			_ = client.LoginInteractive(context.Background())
		}
		return
	}

	model := appllicationDashboard.New()
	// Enter alt-screen right away to minimize flicker and startup delay.
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
