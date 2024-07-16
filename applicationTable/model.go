package applicationTable

import (
	"github.com/FredrikMWold/radix-tui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func New() Model {
	applicationColumns := []table.Column{
		{Title: "name", Width: 28},
	}
	spiner := spinner.New()
	spiner.Spinner = spinner.Meter
	spiner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		table: table.New(
			table.WithColumns(applicationColumns),
			table.WithFocused(true),
			table.WithHeight(8),
			table.WithStyles(styles.TableStyles()),
		),
		spinner:               spiner,
		isLoadingApplications: true,
		focused:               true,
	}
}

type Model struct {
	table                 table.Model
	spinner               spinner.Model
	isLoadingApplications bool
	selectedApp           string
	focused               bool
}

type Application struct {
	Jobs         []Job         `json:"jobs"`
	Environments []Environment `json:"environments"`
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
