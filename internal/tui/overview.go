package tui

import (
	"fmt"
	"strings"

	"github.com/chxmxii/3a/internal/storage"
)

type overviewView struct {
	assessment *storage.Assessment
	resources  []storage.Resource
	findings   []storage.Finding
	costs      []storage.CostEstimate
}

func (v *overviewView) render(width int) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📊 Assessment Overview"))
	b.WriteString("\n\n")

	if v.assessment != nil {
		b.WriteString(headerStyle.Render("Assessment Info"))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("  Profile:  %s\n", v.assessment.Profile))
		b.WriteString(fmt.Sprintf("  Provider: %s\n", v.assessment.Provider))
		b.WriteString(fmt.Sprintf("  Status:   %s\n", v.assessment.Status))
		b.WriteString(fmt.Sprintf("  Started:  %s\n", v.assessment.StartedAt.Format("2006-01-02 15:04:05")))
		if v.assessment.CompletedAt != nil {
			b.WriteString(fmt.Sprintf("  Finished: %s\n", v.assessment.CompletedAt.Format("2006-01-02 15:04:05")))
		}
		b.WriteString(fmt.Sprintf("  Regions:  %s\n", strings.Join(v.assessment.Regions, ", ")))
		b.WriteString("\n")
	}

	// Resource summary.
	b.WriteString(headerStyle.Render("Resources"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  Total: %d\n", len(v.resources)))

	typeCounts := make(map[string]int)
	for _, r := range v.resources {
		typeCounts[r.ResourceType]++
	}
	for t, c := range typeCounts {
		b.WriteString(fmt.Sprintf("    %-20s %d\n", t, c))
	}
	b.WriteString("\n")

	// Findings summary.
	b.WriteString(headerStyle.Render("Findings"))
	b.WriteString("\n")
	severityCounts := map[string]int{"critical": 0, "high": 0, "medium": 0, "low": 0}
	for _, f := range v.findings {
		severityCounts[f.Severity]++
	}
	b.WriteString(fmt.Sprintf("  %s %d  ", severityCriticalStyle.Render("CRITICAL:"), severityCounts["critical"]))
	b.WriteString(fmt.Sprintf("%s %d  ", severityHighStyle.Render("HIGH:"), severityCounts["high"]))
	b.WriteString(fmt.Sprintf("%s %d  ", severityMediumStyle.Render("MEDIUM:"), severityCounts["medium"]))
	b.WriteString(fmt.Sprintf("%s %d\n", severityLowStyle.Render("LOW:"), severityCounts["low"]))
	b.WriteString("\n")

	// Cost summary.
	b.WriteString(headerStyle.Render("Estimated Monthly Cost"))
	b.WriteString("\n")
	totalCost := 0.0
	for _, c := range v.costs {
		if c.MonthlyCost != nil {
			totalCost += *c.MonthlyCost
		}
	}
	b.WriteString(fmt.Sprintf("  $%.2f/month\n", totalCost))

	return b.String()
}
