package appllicationdashboard

import (
	"fmt"
	"os/exec"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
)

type Context string

func getContext() tea.Msg {
	result := exec.Command("rx", "get", "context")
	regPattern := regexp.MustCompile(`'([^']*)'`)
	output, err := result.CombinedOutput()
	platform := regPattern.FindStringSubmatch(string(output))
	if err != nil {
		fmt.Println(err)
	}
	return Context(platform[1])
}
