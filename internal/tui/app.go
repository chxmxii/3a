package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chxmxii/3a/internal/storage"
)

// View represents the currently active TUI view.
type View int

const (
	ViewOverview View = iota
	ViewInventory
	ViewArchitecture
	ViewFindings
	ViewCost
)

// Model is the root Bubble Tea model for the 3A TUI.
type Model struct {
	store        *storage.Store
	assessmentID string
	activeView   View
	width        int
	height       int

	// View data.
	overview     overviewView
	inventory    inventoryView
	architecture architectureView
	findings     findingsView
	cost         costView

	loaded bool
	err    error
}

// dataLoadedMsg is sent when data has been loaded from the store.
type dataLoadedMsg struct{}

// errMsg wraps an error for the TUI.
type errMsg struct{ err error }

// NewModel creates a new TUI model for the given assessment.
func NewModel(store *storage.Store, assessmentID string) Model {
	return Model{
		store:        store,
		assessmentID: assessmentID,
		activeView:   ViewOverview,
	}
}

// Init starts the TUI.
func (m Model) Init() tea.Cmd {
	return m.loadData
}

func (m Model) loadData() tea.Msg {
	assessment, err := m.store.GetAssessment(m.assessmentID)
	if err != nil {
		return errMsg{err}
	}

	resources, err := m.store.GetResourcesByAssessment(m.assessmentID)
	if err != nil {
		return errMsg{err}
	}

	findings, err := m.store.GetFindingsByAssessment(m.assessmentID)
	if err != nil {
		return errMsg{err}
	}

	relationships, err := m.store.GetRelationshipsByAssessment(m.assessmentID)
	if err != nil {
		return errMsg{err}
	}

	costs, err := m.store.GetCostsByAssessment(m.assessmentID)
	if err != nil {
		return errMsg{err}
	}

	m.overview = overviewView{
		assessment: assessment,
		resources:  resources,
		findings:   findings,
		costs:      costs,
	}
	m.inventory = inventoryView{resources: resources}
	m.architecture = architectureView{resources: resources, relationships: relationships}
	m.findings = findingsView{findings: findings}
	m.cost = costView{costs: costs, resources: resources}

	return dataLoadedMsg{}
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.activeView = ViewOverview
		case "2":
			m.activeView = ViewInventory
		case "3":
			m.activeView = ViewArchitecture
		case "4":
			m.activeView = ViewFindings
		case "5":
			m.activeView = ViewCost

		// View-specific keys.
		case "up", "k":
			switch m.activeView {
			case ViewInventory:
				if m.inventory.cursor > 0 {
					m.inventory.cursor--
				}
			case ViewFindings:
				if m.findings.cursor > 0 {
					m.findings.cursor--
				}
			}
		case "down", "j":
			switch m.activeView {
			case ViewInventory:
				filtered := m.inventory.filteredResources()
				if m.inventory.cursor < len(filtered)-1 {
					m.inventory.cursor++
				}
			case ViewFindings:
				filtered := m.findings.filteredFindings()
				if m.findings.cursor < len(filtered)-1 {
					m.findings.cursor++
				}
			}
		case "c":
			if m.activeView == ViewFindings {
				m.findings.severityFilter = "critical"
				m.findings.cursor = 0
			}
		case "h":
			if m.activeView == ViewFindings {
				m.findings.severityFilter = "high"
				m.findings.cursor = 0
			}
		case "m":
			if m.activeView == ViewFindings {
				m.findings.severityFilter = "medium"
				m.findings.cursor = 0
			}
		case "l":
			if m.activeView == ViewFindings {
				m.findings.severityFilter = "low"
				m.findings.cursor = 0
			}
		case "f":
			if m.activeView == ViewFindings {
				m.findings.severityFilter = ""
				m.findings.cursor = 0
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case dataLoadedMsg:
		m.loaded = true

	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

// View renders the TUI.
func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit.", m.err)
	}

	if !m.loaded {
		return "Loading assessment data...\n"
	}

	var content string
	switch m.activeView {
	case ViewOverview:
		content = m.overview.render(m.width)
	case ViewInventory:
		content = m.inventory.render(m.width, m.height)
	case ViewArchitecture:
		content = m.architecture.render(m.width)
	case ViewFindings:
		content = m.findings.render(m.width, m.height)
	case ViewCost:
		content = m.cost.render(m.width)
	}

	// Navigation bar.
	nav := m.renderNav()

	// Help bar.
	help := helpStyle.Render("  q: quit • ↑/↓: navigate • 1-5: switch views")

	return nav + "\n" + content + "\n" + help
}

func (m Model) renderNav() string {
	tabs := []struct {
		key  string
		name string
		view View
	}{
		{"1", "Overview", ViewOverview},
		{"2", "Inventory", ViewInventory},
		{"3", "Architecture", ViewArchitecture},
		{"4", "Findings", ViewFindings},
		{"5", "Cost", ViewCost},
	}

	var parts []string
	for _, tab := range tabs {
		label := fmt.Sprintf(" %s:%s ", tab.key, tab.name)
		if tab.view == m.activeView {
			parts = append(parts, selectedStyle.Render(label))
		} else {
			parts = append(parts, normalStyle.Render(label))
		}
	}

	return " " + fmt.Sprintf("%s", joinStrings(parts, "│"))
}

func joinStrings(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}
