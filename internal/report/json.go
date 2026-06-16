package report

import (
	"encoding/json"
)

// JSONReport is the JSON-serializable report structure.
type JSONReport struct {
	Executive   ExecutiveSummary       `json:"executive_summary"`
	Technical   TechnicalReport        `json:"technical_report"`
	Assessment  assessmentInfo         `json:"assessment"`
	Resources   resourceSummary        `json:"resources"`
	Findings    findingsSummary        `json:"findings"`
	Cost        costSummary            `json:"cost"`
}

type assessmentInfo struct {
	ID       string   `json:"id"`
	Profile  string   `json:"profile"`
	Provider string   `json:"provider"`
	Status   string   `json:"status"`
	Regions  []string `json:"regions"`
}

type resourceSummary struct {
	Total    int            `json:"total"`
	ByType   map[string]int `json:"by_type"`
	ByRegion map[string]int `json:"by_region"`
}

type findingsSummary struct {
	Total    int            `json:"total"`
	BySeverity map[string]int `json:"by_severity"`
}

type costSummary struct {
	TotalMonthly float64            `json:"total_monthly"`
	ByCategory   map[string]float64 `json:"by_category"`
}

// renderJSON generates a JSON report.
func renderJSON(data *ReportData) (string, error) {
	tech := BuildTechnicalReport(data)
	exec := tech.Executive

	sevCounts := map[string]int{
		"critical": exec.CriticalCount,
		"high":     exec.HighCount,
		"medium":   exec.MediumCount,
		"low":      exec.LowCount,
	}

	report := JSONReport{
		Executive: exec,
		Technical: tech,
		Assessment: assessmentInfo{
			ID:       data.Assessment.ID,
			Profile:  data.Assessment.Profile,
			Provider: data.Assessment.Provider,
			Status:   data.Assessment.Status,
			Regions:  data.Assessment.Regions,
		},
		Resources: resourceSummary{
			Total:    len(data.Resources),
			ByType:   tech.ResourcesByType,
			ByRegion: tech.ResourcesByRegion,
		},
		Findings: findingsSummary{
			Total:      len(data.Findings),
			BySeverity: sevCounts,
		},
		Cost: costSummary{
			TotalMonthly: exec.MonthlyCost,
			ByCategory:   tech.CostBreakdown,
		},
	}

	out, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}
