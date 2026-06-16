package tui

import (
	"fmt"
	"strings"

	"github.com/chxmxii/3a/internal/storage"
)

type findingsView struct {
	findings       []storage.Finding
	severityFilter string
	cursor         int
	offset         int
}

func (v *findingsView) render(width, height int) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("🔍 Security Findings"))
	b.WriteString("\n")

	if v.severityFilter != "" {
		b.WriteString(fmt.Sprintf("  Filter: %s (press 'f' to clear)\n", v.severityFilter))
	} else {
		b.WriteString("  Press c/h/m/l to filter by severity, f to clear\n")
	}
	b.WriteString("\n")

	filtered := v.filteredFindings()
	b.WriteString(fmt.Sprintf("  Showing %d of %d findings\n\n", len(filtered), len(v.findings)))

	maxRows := height - 10
	if maxRows < 5 {
		maxRows = 5
	}

	if v.offset > len(filtered)-maxRows {
		v.offset = max(0, len(filtered)-maxRows)
	}

	end := v.offset + maxRows
	if end > len(filtered) {
		end = len(filtered)
	}

	for i := v.offset; i < end; i++ {
		f := filtered[i]

		var sevStyle = normalStyle
		switch f.Severity {
		case "critical":
			sevStyle = severityCriticalStyle
		case "high":
			sevStyle = severityHighStyle
		case "medium":
			sevStyle = severityMediumStyle
		case "low":
			sevStyle = severityLowStyle
		}

		sevLabel := sevStyle.Render(fmt.Sprintf("%-8s", strings.ToUpper(f.Severity)))

		desc := f.Description
		if len(desc) > 70 {
			desc = desc[:67] + "..."
		}

		line := fmt.Sprintf("  %s  %s", sevLabel, desc)
		if i == v.cursor {
			b.WriteString(selectedStyle.Render(line))
		} else {
			b.WriteString(line)
		}
		b.WriteString("\n")

		// Show recommendation for selected item.
		if i == v.cursor && f.Recommendation != "" {
			rec := f.Recommendation
			if len(rec) > 80 {
				rec = rec[:77] + "..."
			}
			b.WriteString(fmt.Sprintf("           %s\n", helpStyle.Render("→ "+rec)))
		}
	}

	return b.String()
}

func (v *findingsView) filteredFindings() []storage.Finding {
	if v.severityFilter == "" {
		return v.findings
	}
	var filtered []storage.Finding
	for _, f := range v.findings {
		if f.Severity == v.severityFilter {
			filtered = append(filtered, f)
		}
	}
	return filtered
}
