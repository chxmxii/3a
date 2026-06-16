package checklist

import (
	"fmt"

	"github.com/chxmxii/3a/internal/storage"
)

// Engine generates a compliance checklist from assessment findings.
type Engine struct {
	store *storage.Store
}

// NewEngine creates a new checklist engine.
func NewEngine(store *storage.Store) *Engine {
	return &Engine{store: store}
}

// Generate creates a checklist summary from findings in the given assessment.
func (e *Engine) Generate(assessmentID string) (*ChecklistSummary, error) {
	findings, err := e.store.GetFindingsByAssessment(assessmentID)
	if err != nil {
		return nil, fmt.Errorf("loading findings: %w", err)
	}

	// Group findings by control_id.
	byControl := make(map[string][]storage.Finding)
	for _, f := range findings {
		key := f.ControlID + ": " + f.StandardName
		byControl[key] = append(byControl[key], f)
	}

	summary := &ChecklistSummary{
		ByCategory: make(map[string][]CheckItem),
	}

	// Define all checks.
	checks := allChecks()

	for _, check := range checks {
		item := CheckItem{
			Name:        check.Name,
			Description: check.Description,
			Category:    check.Category,
			Status:      StatusPass, // default to pass
		}

		// Look for findings matching this check.
		for key, controlFindings := range byControl {
			if matchesCheck(key, check) {
				// Determine status based on severity.
				for _, f := range controlFindings {
					item.ResourceIDs = append(item.ResourceIDs, f.ResourceID)
					switch f.Severity {
					case "critical", "high":
						item.Status = StatusFail
					case "medium":
						if item.Status != StatusFail {
							item.Status = StatusWarn
						}
					case "low", "informational":
						if item.Status == StatusPass {
							item.Status = StatusWarn
						}
					}
				}
				item.Details = fmt.Sprintf("%d resource(s) affected", len(controlFindings))
			}
		}

		summary.Items = append(summary.Items, item)
		summary.ByCategory[item.Category] = append(summary.ByCategory[item.Category], item)

		switch item.Status {
		case StatusPass:
			summary.PassCount++
		case StatusFail:
			summary.FailCount++
		case StatusWarn:
			summary.WarnCount++
		}
	}

	return summary, nil
}

type checkDef struct {
	Name        string
	Description string
	Category    string
	ControlIDs  []string
}

func matchesCheck(key string, check checkDef) bool {
	for _, id := range check.ControlIDs {
		if len(key) >= len(id) && key[:len(id)] == id {
			return true
		}
	}
	return false
}

func allChecks() []checkDef {
	return []checkDef{
		{Name: "S3 Public Access", Description: "No S3 buckets allow public access", Category: "Security", ControlIDs: []string{"SEC-001"}},
		{Name: "Security Group Restrictions", Description: "Security groups restrict access on dangerous ports", Category: "Security", ControlIDs: []string{"SEC-002"}},
		{Name: "EBS Encryption", Description: "All EBS volumes are encrypted", Category: "Security", ControlIDs: []string{"SEC-003"}},
		{Name: "RDS Private Access", Description: "RDS instances are not publicly accessible", Category: "Security", ControlIDs: []string{"SEC-004"}},
		{Name: "IAM MFA Enabled", Description: "All IAM users have MFA enabled", Category: "Security", ControlIDs: []string{"SEC-005"}},
		{Name: "EKS Private Endpoint", Description: "EKS clusters use private endpoints", Category: "Security", ControlIDs: []string{"SEC-006"}},
		{Name: "S3 Encryption", Description: "All S3 buckets have default encryption", Category: "Security", ControlIDs: []string{"SEC-007"}},
		{Name: "OCI Bucket Privacy", Description: "Object Storage buckets are not public", Category: "Security", ControlIDs: []string{"SEC-008"}},
		{Name: "OCI NSG Rules", Description: "NSGs restrict ingress appropriately", Category: "Security", ControlIDs: []string{"SEC-009"}},
		{Name: "OCI Volume Encryption", Description: "Block volumes use customer-managed keys", Category: "Security", ControlIDs: []string{"SEC-010"}},
		{Name: "OCI DB Protection", Description: "Database systems have NSG protection", Category: "Security", ControlIDs: []string{"SEC-011"}},
	}
}
