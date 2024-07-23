package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SelectedApplication string

func SelectApplication(application string) tea.Cmd {
	return func() tea.Msg {
		return SelectedApplication(application)
	}
}

type Applications []string

func GetApplications() tea.Msg {
	result := exec.Command("rx", "get", "application")
	output, err := result.Output()
	if err != nil {
		fmt.Println(err)
	}
	trimmed := strings.TrimSpace(string(output))
	application_list := strings.Split(trimmed, "\n")
	if strings.Contains(application_list[0], "login.microsoft") {
		application_list = GetApplications().(Applications)
	}
	return Applications(application_list)
}

type Application struct {
	Jobs         []Job         `json:"jobs"`
	Environments []Environment `json:"environments"`
	Name         string        `json:"name"`
}

type Job struct {
	Name         string   `json:"name"`
	TriggeredBy  string   `json:"triggeredBy"`
	Environments []string `json:"environments"`
	Pipeline     string   `json:"pipeline"`
	Status       string   `json:"status"`
	Created      string   `json:"created"`
}

type Environment struct {
	BranchMapping string `json:"branchMapping"`
	Name          string `json:"name"`
	Status        string `json:"status"`
}

func GetApplicationData(application string) tea.Cmd {
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
