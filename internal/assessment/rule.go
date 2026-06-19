package assessment

import (
	"context"

	"github.com/chxmxii/a3/internal/provider"
	"github.com/chxmxii/a3/internal/storage"
)

// Severity represents the impact level of a finding.
type Severity string

const (
	SeverityCritical      Severity = "critical"
	SeverityHigh          Severity = "high"
	SeverityMedium        Severity = "medium"
	SeverityLow           Severity = "low"
	SeverityInformational Severity = "informational"
)

// FindingCategory represents the assessment domain of a finding.
type FindingCategory string

const (
	CategorySecurity              FindingCategory = "Security"
	CategoryReliability           FindingCategory = "Reliability"
	CategoryPerformance           FindingCategory = "Performance"
	CategoryCostOptimization      FindingCategory = "Cost Optimization"
	CategoryOperationalExcellence FindingCategory = "Operational Excellence"
)

// Rule is a single compliance check evaluated against resources.
type Rule interface {
	// ID returns the unique rule identifier.
	ID() string

	// Standard returns the standard this rule belongs to (e.g., "CIS AWS 1.5").
	Standard() string

	// ControlID returns the specific control (e.g., "2.1.1").
	ControlID() string

	// Category returns the finding category.
	Category() FindingCategory

	// AppliesTo returns which resource types this rule evaluates.
	AppliesTo() []provider.ResourceType

	// Evaluate checks a resource and returns findings.
	Evaluate(ctx context.Context, resource storage.Resource) ([]Finding, error)
}

// Finding represents a single compliance violation or observation.
type Finding struct {
	Severity       Severity
	ResourceID     string
	Description    string // max 500 chars
	Recommendation string // max 1000 chars
	StandardName   string
	ControlID      string
	Category       FindingCategory
}
