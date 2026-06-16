package provider

import "fmt"

// ProviderConfig holds authentication and connection details.
type ProviderConfig struct {
	ProviderType string
	Credentials  CredentialSource
	Regions      []string
}

// CredentialSource describes how to authenticate.
type CredentialSource struct {
	Type        string // "profile", "env", "config_file"
	ProfileName string // AWS named profile or OCI config profile
}

// ProviderFactory creates a Provider from configuration.
type ProviderFactory func(cfg ProviderConfig) (Provider, error)

// Registry manages available provider implementations.
type Registry struct {
	providers map[string]ProviderFactory
}

// NewRegistry creates a new empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]ProviderFactory),
	}
}

// Register adds a provider factory under the given name.
func (r *Registry) Register(name string, factory ProviderFactory) {
	r.providers[name] = factory
}

// Get creates a Provider instance for the given config.
// Returns an error if no factory is registered for the provider type.
func (r *Registry) Get(cfg ProviderConfig) (Provider, error) {
	factory, ok := r.providers[cfg.ProviderType]
	if !ok {
		return nil, fmt.Errorf("no provider registered for type %q", cfg.ProviderType)
	}
	return factory(cfg)
}
