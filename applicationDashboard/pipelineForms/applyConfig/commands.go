package applyconfig

import (
	"context"
	"fmt"

	"github.com/FredrikMWold/radix-tui/applicationDashboard/commands"
	"github.com/FredrikMWold/radix-tui/internal/radix"
	tea "github.com/charmbracelet/bubbletea"
)

func ApplyConfig(application string) tea.Cmd {
	return func() tea.Msg {
		client, err := radix.New(false)
		if err == nil {
			_, err = client.TriggerApplyConfig(context.Background(), application)
		}
		if err != nil {
			fmt.Println(err)
		}
		return commands.SelectedApplication(application)
	}

}
