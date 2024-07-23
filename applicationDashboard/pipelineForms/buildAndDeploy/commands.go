package buildanddeploy

import (
	"fmt"
	"os/exec"

	"github.com/FredrikMWold/radix-tui/applicationDashboard/commands"
	tea "github.com/charmbracelet/bubbletea"
)

func BuildAndDeploy(application string, branch string) tea.Cmd {
	return func() tea.Msg {
		result := exec.Command("rx", "create", "pipeline-job", "build-deploy", "-a", application, "-b", branch)
		_, err := result.Output()
		if err != nil {
			fmt.Println(err)
		}
		return commands.SelectedApplication(application)
	}
}
