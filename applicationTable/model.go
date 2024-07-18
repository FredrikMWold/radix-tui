package applicationTable

import (
	"github.com/FredrikMWold/radix-tui/styles"
	"github.com/charmbracelet/bubbles/table"
)

func New() Model {
	applicationColumns := []table.Column{
		{Title: "name", Width: 28},
	}

	return Model{
		table: table.New(
			table.WithColumns(applicationColumns),
			table.WithFocused(true),
			table.WithHeight(8),
			table.WithStyles(styles.TableStyles()),
		),
		isLoadingApplications: true,
	}
}

type Model struct {
	table                 table.Model
	isLoadingApplications bool
}
