package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the top-level configuration.
type Config struct {
	DBPath    string           `yaml:"db_path"`
	Steampipe SteampipeConfig  `yaml:"steampipe"`
	Profiles  []AccountProfile `yaml:"profiles"`
}

// SteampipeConfig holds Steampipe connection settings.
type SteampipeConfig struct {
	ConnectionString string `yaml:"connection_string"`
}


// AccountProfile defines a cloud account connection.
type AccountProfile struct {
	Name        string     `yaml:"name"`
	DisplayName string     `yaml:"display_name"`
	Provider    string     `yaml:"provider"` // "aws" or "oci"
	Credentials Credential `yaml:"credentials"`
	Regions     []string   `yaml:"regions"`
}

// Credential specifies how to authenticate.
type Credential struct {
	Type        string `yaml:"type"`         // "profile", "env", "config_file"
	ProfileName string `yaml:"profile_name"` // named profile reference
}

// ConfigError represents an error related to configuration loading or parsing.
type ConfigError struct {
	Path    string
	Message string
	Err     error
}

func (e *ConfigError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("config error (%s): %s: %v", e.Path, e.Message, e.Err)
	}
	return fmt.Sprintf("config error (%s): %s", e.Path, e.Message)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

// DefaultConfigPath returns the default configuration file path: ~/.3a/config.yaml
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("~", ".3a", "config.yaml")
	}
	return filepath.Join(home, ".3a", "config.yaml")
}

// Load reads and parses the YAML configuration file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ConfigError{
			Path:    path,
			Message: "unable to read config file",
			Err:     err,
		}
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, &ConfigError{
			Path:    path,
			Message: "malformed config file",
			Err:     err,
		}
	}

	return &cfg, nil
}

// Save writes the configuration to a YAML file at the given path.
// It creates parent directories if they do not exist.
func Save(path string, cfg *Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return &ConfigError{
			Path:    path,
			Message: "unable to create config directory",
			Err:     err,
		}
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return &ConfigError{
			Path:    path,
			Message: "unable to serialize config",
			Err:     err,
		}
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return &ConfigError{
			Path:    path,
			Message: "unable to write config file",
			Err:     err,
		}
	}

	return nil
}
