package assessment

import (
	"context"
	"log"

	"github.com/chxmxii/3a/internal/provider"
	"github.com/chxmxii/3a/internal/storage"
)

// Engine evaluates resources against compliance rules.
type Engine struct {
	store *storage.Store
	rules []Rule
}

// NewEngine creates a new assessment engine with the given store and rules.
func NewEngine(store *storage.Store, rules []Rule) *Engine {
	return &Engine{
		store: store,
		rules: rules,
	}
}

// Run executes the assessment for the given assessment ID.
// It loads all resources, evaluates applicable rules, and persists findings.
// Rule evaluation errors are logged and skipped; the engine continues processing.
func (e *Engine) Run(ctx context.Context, assessmentID string) error {
	resources, err := e.store.GetResourcesByAssessment(assessmentID)
	if err != nil {
		return err
	}

	var allFindings []Finding

	for _, resource := range resources {
		applicable := e.findApplicableRules(provider.ResourceType(resource.ResourceType))

		for _, rule := range applicable {
			findings, err := rule.Evaluate(ctx, resource)
			if err != nil {
				log.Printf("assessment: rule %s failed for resource %s: %v", rule.ID(), resource.ResourceID, err)
				continue
			}
			allFindings = append(allFindings, findings...)
		}
	}

	// Persist all findings to storage.
	for i := range allFindings {
		f := &storage.Finding{
			AssessmentID:   assessmentID,
			Severity:       string(allFindings[i].Severity),
			ResourceID:     allFindings[i].ResourceID,
			Description:    allFindings[i].Description,
			Recommendation: allFindings[i].Recommendation,
			StandardName:   allFindings[i].StandardName,
			ControlID:      allFindings[i].ControlID,
			Category:       string(allFindings[i].Category),
		}
		if err := e.store.InsertFinding(f); err != nil {
			log.Printf("assessment: failed to persist finding for resource %s: %v", allFindings[i].ResourceID, err)
		}
	}

	return nil
}

// findApplicableRules returns rules whose AppliesTo includes the given resource type.
func (e *Engine) findApplicableRules(resourceType provider.ResourceType) []Rule {
	var applicable []Rule
	for _, rule := range e.rules {
		for _, rt := range rule.AppliesTo() {
			if rt == resourceType {
				applicable = append(applicable, rule)
				break
			}
		}
	}
	return applicable
}
