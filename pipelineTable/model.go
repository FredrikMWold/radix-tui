package pipelinetable

import (
	"github.com/FredrikMWold/radix-tui/commands"
	"github.com/FredrikMWold/radix-tui/styles"
	"github.com/charmbracelet/bubbles/table"
)

func New() Model {
	return Model{
		table: table.New(
			table.WithFocused(true),
			table.WithStyles(styles.TableStyles()),
		),
	}
}

type Model struct {
	table       table.Model
	application commands.Application
}
