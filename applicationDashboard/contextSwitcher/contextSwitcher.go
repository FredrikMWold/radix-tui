package contextswitcher

import (
	"github.com/FredrikMWold/radix-tui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	radixconfig "github.com/equinor/radix-cli/pkg/config"
)

// ContextSelected is sent when a context is selected
type ContextSelected string

// ContextCancelled is sent when the context switcher is cancelled
type ContextCancelled struct{}

type contextItem struct {
	name        string
	description string
}

func (i contextItem) Title() string       { return i.name }
func (i contextItem) Description() string { return i.description }
func (i contextItem) FilterValue() string { return i.name }

type Model struct {
	list           list.Model
	currentContext string
	width          int
	height         int
}

func New(currentContext string) Model {
	items := []list.Item{
		contextItem{name: radixconfig.ContextPlatform, description: "Production platform (default)"},
		contextItem{name: radixconfig.ContextPlatform2, description: "Platform 2 (c2)"},
		contextItem{name: radixconfig.ContextPlayground, description: "Playground environment"},
		contextItem{name: radixconfig.ContextDevelopment, description: "Development environment"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#cba6f7")).
		BorderLeftForeground(lipgloss.Color("#cba6f7"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("241")).
		BorderLeftForeground(lipgloss.Color("#cba6f7"))

	l := list.New(items, delegate, 40, 17)
	l.Title = "Switch Context"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true)
	l.Styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#cba6f7")).
		Padding(0, 1)

	// Set focus on current context
	for i, item := range items {
		if item.(contextItem).name == currentContext {
			l.Select(i)
			break
		}
	}

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
		}
	}

	return Model{
		list:           l,
		currentContext: currentContext,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if item, ok := m.list.SelectedItem().(contextItem); ok {
				return m, func() tea.Msg { return ContextSelected(item.name) }
			}
		case "esc", "q":
			return m, func() tea.Msg { return ContextCancelled{} }
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(min(50, msg.Width-4), min(17, msg.Height-4))
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	modalStyle := styles.SectionContainer(true).
		Width(52).
		Height(16).
		Align(lipgloss.Center)

	return modalStyle.Render(m.list.View())
}
