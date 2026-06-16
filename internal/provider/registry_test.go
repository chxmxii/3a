package provider

import (
	"context"
	"errors"
	"testing"
)

// mockProvider is a minimal Provider implementation for testing.
type mockProvider struct {
	name string
}

func (m *mockProvider) Name() string                          { return m.name }
func (m *mockProvider) Authenticate(_ context.Context) error  { return nil }
func (m *mockProvider) Discoverer() Discoverer                { return nil }
func (m *mockProvider) MetricsClient() MetricsClient          { return nil }
func (m *mockProvider) PricingClient() PricingClient          { return nil }

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry returned nil")
	}
	if r.providers == nil {
		t.Fatal("providers map is nil")
	}
}

func TestRegisterAndGet(t *testing.T) {
	r := NewRegistry()

	factory := func(cfg ProviderConfig) (Provider, error) {
		return &mockProvider{name: cfg.ProviderType}, nil
	}

	r.Register("aws", factory)

	cfg := ProviderConfig{
		ProviderType: "aws",
		Credentials: CredentialSource{
			Type:        "profile",
			ProfileName: "default",
		},
		Regions: []string{"us-east-1"},
	}

	p, err := r.Get(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "aws" {
		t.Errorf("expected provider name %q, got %q", "aws", p.Name())
	}
}

func TestGetUnregisteredProvider(t *testing.T) {
	r := NewRegistry()

	cfg := ProviderConfig{
		ProviderType: "gcp",
		Credentials: CredentialSource{
			Type: "env",
		},
		Regions: []string{"us-central1"},
	}

	_, err := r.Get(cfg)
	if err == nil {
		t.Fatal("expected error for unregistered provider, got nil")
	}

	expected := `no provider registered for type "gcp"`
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestMultipleProviders(t *testing.T) {
	r := NewRegistry()

	awsFactory := func(cfg ProviderConfig) (Provider, error) {
		return &mockProvider{name: "aws"}, nil
	}
	ociFactory := func(cfg ProviderConfig) (Provider, error) {
		return &mockProvider{name: "oci"}, nil
	}

	r.Register("aws", awsFactory)
	r.Register("oci", ociFactory)

	// Get AWS provider
	awsCfg := ProviderConfig{
		ProviderType: "aws",
		Credentials: CredentialSource{
			Type:        "profile",
			ProfileName: "prod",
		},
		Regions: []string{"us-east-1", "eu-west-1"},
	}
	p, err := r.Get(awsCfg)
	if err != nil {
		t.Fatalf("unexpected error getting aws: %v", err)
	}
	if p.Name() != "aws" {
		t.Errorf("expected aws, got %q", p.Name())
	}

	// Get OCI provider
	ociCfg := ProviderConfig{
		ProviderType: "oci",
		Credentials: CredentialSource{
			Type:        "config_file",
			ProfileName: "DEFAULT",
		},
		Regions: []string{"us-ashburn-1"},
	}
	p, err = r.Get(ociCfg)
	if err != nil {
		t.Fatalf("unexpected error getting oci: %v", err)
	}
	if p.Name() != "oci" {
		t.Errorf("expected oci, got %q", p.Name())
	}
}

func TestFactoryError(t *testing.T) {
	r := NewRegistry()

	errFactory := func(cfg ProviderConfig) (Provider, error) {
		return nil, errors.New("authentication failed")
	}

	r.Register("aws", errFactory)

	cfg := ProviderConfig{
		ProviderType: "aws",
		Credentials: CredentialSource{
			Type:        "profile",
			ProfileName: "invalid",
		},
		Regions: []string{"us-east-1"},
	}

	_, err := r.Get(cfg)
	if err == nil {
		t.Fatal("expected error from factory, got nil")
	}
	if err.Error() != "authentication failed" {
		t.Errorf("expected 'authentication failed', got %q", err.Error())
	}
}

func TestRegisterOverwrite(t *testing.T) {
	r := NewRegistry()

	firstFactory := func(cfg ProviderConfig) (Provider, error) {
		return &mockProvider{name: "first"}, nil
	}
	secondFactory := func(cfg ProviderConfig) (Provider, error) {
		return &mockProvider{name: "second"}, nil
	}

	r.Register("aws", firstFactory)
	r.Register("aws", secondFactory)

	cfg := ProviderConfig{
		ProviderType: "aws",
		Credentials:  CredentialSource{Type: "env"},
		Regions:      []string{"us-east-1"},
	}

	p, err := r.Get(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "second" {
		t.Errorf("expected overwritten factory to produce 'second', got %q", p.Name())
	}
}
