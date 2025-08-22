package appllicationdashboard

import (
	tea "github.com/charmbracelet/bubbletea"
	radixconfig "github.com/equinor/radix-cli/pkg/config"
)

type Context string

func getContext() tea.Msg {
	// Read from radix config (same as rx uses)
	cfg, err := radixconfig.GetRadixConfig()
	if err != nil || cfg == nil || cfg.CustomConfig == nil {
		return Context("")
	}
	return Context(cfg.CustomConfig.Context)
}
