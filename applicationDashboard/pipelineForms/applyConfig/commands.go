package applyconfig

import (
	"fmt"
	"os/exec"

	"github.com/FredrikMWold/radix-tui/applicationDashboard/commands"
	tea "github.com/charmbracelet/bubbletea"
)

func ApplyConfig(application string) tea.Cmd {
	return func() tea.Msg {
		result := exec.Command("rx", "create", "pipeline-job", "apply-config", "-a", application)
		_, err := result.Output()
		if err != nil {
			fmt.Println(err)
		}
		return commands.SelectedApplication(application)
	}

}
