package applicationTable

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func tick() tea.Cmd {
	return tea.Tick(time.Second*10, func(t time.Time) tea.Msg {
		return UpdateApplicationDataTick(t)
	})
}

func getApplicationData(application string) tea.Cmd {
	return func() tea.Msg {
		result := exec.Command("rx", "get", "application", "-a", application)
		output, err := result.Output()
		if err != nil {
			fmt.Println(err)
		}
		var application Application
		err = json.Unmarshal(output, &application)
		if err != nil {
			fmt.Println(err)
		}
		return application
	}
}

func selectApplication(application string) tea.Cmd {
	return func() tea.Msg {
		return SelectedApplication(application)
	}
}

func getApplications() tea.Msg {
	result := exec.Command("rx", "get", "application")
	output, err := result.Output()
	if err != nil {
		fmt.Println(err)
	}
	trimmed := strings.TrimSpace(string(output))
	application_list := strings.Split(trimmed, "\n")
	if strings.Contains(application_list[0], "login.microsoft") {
		application_list = getApplications().(Applications)
	}
	return Applications(application_list)
}
