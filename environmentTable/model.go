package environmenttable

import (
	"github.com/FredrikMWold/radix-tui/styles"
	"github.com/charmbracelet/bubbles/table"
)

type Model struct {
	table table.Model
}

func New() Model {
	return Model{
		table: table.New(
			table.WithFocused(false),
			table.WithHeight(5),
			table.WithStyles(styles.InfoTableStyles()),
			table.WithColumns([]table.Column{
				{Title: "name", Width: 7},
				{Title: "status", Width: 10},
				{Title: "branch", Width: 7},
			}),
		),
	}
}
