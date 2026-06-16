package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LoadingModel provides a spinner component for loading states.
type LoadingModel struct {
	Spinner spinner.Model
	Message string
}

// NewLoading creates a new loading spinner with the given message.
func NewLoading(message string) LoadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))
	return LoadingModel{
		Spinner: s,
		Message: message,
	}
}

// Init starts the spinner.
func (m LoadingModel) Init() tea.Cmd {
	return m.Spinner.Tick
}

// Update handles spinner messages.
func (m LoadingModel) Update(msg tea.Msg) (LoadingModel, tea.Cmd) {
	var cmd tea.Cmd
	m.Spinner, cmd = m.Spinner.Update(msg)
	return m, cmd
}

// View renders the spinner with its message.
func (m LoadingModel) View() string {
	return m.Spinner.View() + " " + m.Message
}
