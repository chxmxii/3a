package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#7C3AED")
	successColor   = lipgloss.Color("#10B981")
	warningColor   = lipgloss.Color("#F59E0B")
	dangerColor    = lipgloss.Color("#EF4444")
	mutedColor     = lipgloss.Color("#6B7280")
	bgColor        = lipgloss.Color("#1F2937")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F9FAFB"))

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D5DB"))

	severityCriticalStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(dangerColor)

	severityHighStyle = lipgloss.NewStyle().
				Foreground(dangerColor)

	severityMediumStyle = lipgloss.NewStyle().
				Foreground(warningColor)

	severityLowStyle = lipgloss.NewStyle().
				Foreground(successColor)

	passStyle = lipgloss.NewStyle().
			Foreground(successColor)

	failStyle = lipgloss.NewStyle().
			Foreground(dangerColor)

	warnStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Padding(1, 2)
)
