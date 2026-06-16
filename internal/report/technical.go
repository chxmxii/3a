package report

// TechnicalReport provides detailed technical findings.
type TechnicalReport struct {
	Executive    ExecutiveSummary
	ResourcesByType map[string]int
	ResourcesByRegion map[string]int
	FindingsByCategory map[string][]FindingDetail
	CostBreakdown map[string]float64
	Relationships int
}

// FindingDetail holds a finding with additional context.
type FindingDetail struct {
	Severity       string
	ResourceID     string
	Description    string
	Recommendation string
	ControlID      string
}

// BuildTechnicalReport creates a full technical report from the data.
func BuildTechnicalReport(data *ReportData) TechnicalReport {
	report := TechnicalReport{
		Executive:          BuildExecutiveSummary(data),
		ResourcesByType:    make(map[string]int),
		ResourcesByRegion:  make(map[string]int),
		FindingsByCategory: make(map[string][]FindingDetail),
		CostBreakdown:      make(map[string]float64),
		Relationships:      len(data.Relationships),
	}

	for _, r := range data.Resources {
		report.ResourcesByType[r.ResourceType]++
		report.ResourcesByRegion[r.Region]++
	}

	for _, f := range data.Findings {
		detail := FindingDetail{
			Severity:       f.Severity,
			ResourceID:     f.ResourceID,
			Description:    f.Description,
			Recommendation: f.Recommendation,
			ControlID:      f.ControlID,
		}
		report.FindingsByCategory[f.Category] = append(report.FindingsByCategory[f.Category], detail)
	}

	for _, c := range data.Costs {
		if c.MonthlyCost != nil {
			report.CostBreakdown[c.Category] += *c.MonthlyCost
		}
	}

	return report
}
