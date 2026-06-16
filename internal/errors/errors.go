package errors

import "fmt"

// AuthError indicates authentication failure.
type AuthError struct {
	Provider       string
	CredentialType string
	Cause          error
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("authentication failed for %s provider using %s: %v",
		e.Provider, e.CredentialType, e.Cause)
}

func (e *AuthError) Unwrap() error {
	return e.Cause
}

// DiscoveryError indicates a service/region scan failure after retries.
type DiscoveryError struct {
	Service string
	Region  string
	Retries int
	Cause   error
}

func (e *DiscoveryError) Error() string {
	return fmt.Sprintf("discovery failed for %s in %s after %d retries: %v",
		e.Service, e.Region, e.Retries, e.Cause)
}

func (e *DiscoveryError) Unwrap() error {
	return e.Cause
}

// RuleEvaluationError indicates a rule failed for a resource.
type RuleEvaluationError struct {
	RuleID     string
	ResourceID string
	Cause      error
}

func (e *RuleEvaluationError) Error() string {
	return fmt.Sprintf("rule %s failed for resource %s: %v",
		e.RuleID, e.ResourceID, e.Cause)
}

func (e *RuleEvaluationError) Unwrap() error {
	return e.Cause
}
