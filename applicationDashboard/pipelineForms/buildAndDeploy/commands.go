package buildanddeploy

import (
	"context"
	"fmt"

	"github.com/FredrikMWold/radix-tui/applicationDashboard/commands"
	"github.com/FredrikMWold/radix-tui/internal/radix"
	tea "github.com/charmbracelet/bubbletea"
)

func BuildAndDeploy(application string, branch string) tea.Cmd {
	return func() tea.Msg {
		client, err := radix.New(false)
		if err == nil {
			_, err = client.TriggerBuildDeploy(context.Background(), application, branch, "")
		}
		if err != nil {
			fmt.Println(err)
		}
		return commands.SelectedApplication(application)
	}
}
