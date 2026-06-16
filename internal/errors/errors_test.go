package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestAuthError_Error(t *testing.T) {
	cause := fmt.Errorf("invalid credentials")
	err := &AuthError{
		Provider:       "aws",
		CredentialType: "profile",
		Cause:          cause,
	}

	msg := err.Error()
	if msg != "authentication failed for aws provider using profile: invalid credentials" {
		t.Errorf("unexpected error message: %s", msg)
	}
}

func TestAuthError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("token expired")
	err := &AuthError{
		Provider:       "oci",
		CredentialType: "config_file",
		Cause:          cause,
	}

	if err.Unwrap() != cause {
		t.Error("Unwrap() did not return the wrapped cause")
	}
}

func TestAuthError_ErrorsIs(t *testing.T) {
	sentinel := fmt.Errorf("base error")
	err := &AuthError{
		Provider:       "aws",
		CredentialType: "env",
		Cause:          sentinel,
	}

	if !errors.Is(err, sentinel) {
		t.Error("errors.Is() should find the wrapped sentinel error")
	}
}

func TestAuthError_ErrorsAs(t *testing.T) {
	cause := fmt.Errorf("access denied")
	err := &AuthError{
		Provider:       "aws",
		CredentialType: "profile",
		Cause:          cause,
	}

	var target *AuthError
	if !errors.As(err, &target) {
		t.Fatal("errors.As() should match *AuthError")
	}
	if target.Provider != "aws" {
		t.Errorf("expected provider 'aws', got %q", target.Provider)
	}
	if target.CredentialType != "profile" {
		t.Errorf("expected credential type 'profile', got %q", target.CredentialType)
	}
}

func TestDiscoveryError_Error(t *testing.T) {
	cause := fmt.Errorf("timeout")
	err := &DiscoveryError{
		Service: "ec2",
		Region:  "us-east-1",
		Retries: 3,
		Cause:   cause,
	}

	msg := err.Error()
	if msg != "discovery failed for ec2 in us-east-1 after 3 retries: timeout" {
		t.Errorf("unexpected error message: %s", msg)
	}
}

func TestDiscoveryError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("connection refused")
	err := &DiscoveryError{
		Service: "s3",
		Region:  "eu-west-1",
		Retries: 2,
		Cause:   cause,
	}

	if err.Unwrap() != cause {
		t.Error("Unwrap() did not return the wrapped cause")
	}
}

func TestDiscoveryError_ErrorsIs(t *testing.T) {
	sentinel := fmt.Errorf("network error")
	err := &DiscoveryError{
		Service: "rds",
		Region:  "ap-southeast-1",
		Retries: 3,
		Cause:   sentinel,
	}

	if !errors.Is(err, sentinel) {
		t.Error("errors.Is() should find the wrapped sentinel error")
	}
}

func TestDiscoveryError_ErrorsAs(t *testing.T) {
	cause := fmt.Errorf("rate limited")
	err := &DiscoveryError{
		Service: "lambda",
		Region:  "us-west-2",
		Retries: 3,
		Cause:   cause,
	}

	var target *DiscoveryError
	if !errors.As(err, &target) {
		t.Fatal("errors.As() should match *DiscoveryError")
	}
	if target.Service != "lambda" {
		t.Errorf("expected service 'lambda', got %q", target.Service)
	}
	if target.Region != "us-west-2" {
		t.Errorf("expected region 'us-west-2', got %q", target.Region)
	}
	if target.Retries != 3 {
		t.Errorf("expected retries 3, got %d", target.Retries)
	}
}

func TestRuleEvaluationError_Error(t *testing.T) {
	cause := fmt.Errorf("metadata missing")
	err := &RuleEvaluationError{
		RuleID:     "CIS-2.1.1",
		ResourceID: "arn:aws:s3:::my-bucket",
		Cause:      cause,
	}

	msg := err.Error()
	if msg != "rule CIS-2.1.1 failed for resource arn:aws:s3:::my-bucket: metadata missing" {
		t.Errorf("unexpected error message: %s", msg)
	}
}

func TestRuleEvaluationError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("nil pointer")
	err := &RuleEvaluationError{
		RuleID:     "WA-3.2",
		ResourceID: "i-12345",
		Cause:      cause,
	}

	if err.Unwrap() != cause {
		t.Error("Unwrap() did not return the wrapped cause")
	}
}

func TestRuleEvaluationError_ErrorsIs(t *testing.T) {
	sentinel := fmt.Errorf("evaluation failure")
	err := &RuleEvaluationError{
		RuleID:     "SEC-1.0",
		ResourceID: "vpc-abc",
		Cause:      sentinel,
	}

	if !errors.Is(err, sentinel) {
		t.Error("errors.Is() should find the wrapped sentinel error")
	}
}

func TestRuleEvaluationError_ErrorsAs(t *testing.T) {
	cause := fmt.Errorf("parse error")
	err := &RuleEvaluationError{
		RuleID:     "PERF-4.1",
		ResourceID: "db-instance-1",
		Cause:      cause,
	}

	var target *RuleEvaluationError
	if !errors.As(err, &target) {
		t.Fatal("errors.As() should match *RuleEvaluationError")
	}
	if target.RuleID != "PERF-4.1" {
		t.Errorf("expected rule ID 'PERF-4.1', got %q", target.RuleID)
	}
	if target.ResourceID != "db-instance-1" {
		t.Errorf("expected resource ID 'db-instance-1', got %q", target.ResourceID)
	}
}

func TestNestedErrorChain(t *testing.T) {
	// Test a chain: RuleEvaluationError wraps DiscoveryError wraps base error
	base := fmt.Errorf("connection reset")
	discovery := &DiscoveryError{
		Service: "ec2",
		Region:  "us-east-1",
		Retries: 3,
		Cause:   base,
	}
	rule := &RuleEvaluationError{
		RuleID:     "SEC-1.0",
		ResourceID: "i-abc",
		Cause:      discovery,
	}

	// errors.Is should traverse the chain
	if !errors.Is(rule, base) {
		t.Error("errors.Is() should find the base error through the chain")
	}

	// errors.As should find intermediate error
	var disc *DiscoveryError
	if !errors.As(rule, &disc) {
		t.Fatal("errors.As() should find *DiscoveryError in chain")
	}
	if disc.Service != "ec2" {
		t.Errorf("expected service 'ec2', got %q", disc.Service)
	}
}
