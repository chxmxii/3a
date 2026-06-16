package config

import "fmt"

// ProfileNotFoundError is returned when a requested profile does not exist in the configuration.
type ProfileNotFoundError struct {
	RequestedName  string
	AvailableNames []string
}

func (e *ProfileNotFoundError) Error() string {
	return fmt.Sprintf("profile %q not found; available profiles: %v", e.RequestedName, e.AvailableNames)
}

// GetProfile returns the profile with the given name from the configuration.
// If no profile matches, it returns a ProfileNotFoundError listing available names.
func GetProfile(cfg *Config, name string) (*AccountProfile, error) {
	available := make([]string, 0, len(cfg.Profiles))
	for i := range cfg.Profiles {
		if cfg.Profiles[i].Name == name {
			return &cfg.Profiles[i], nil
		}
		available = append(available, cfg.Profiles[i].Name)
	}
	return nil, &ProfileNotFoundError{
		RequestedName:  name,
		AvailableNames: available,
	}
}

// ListProfiles returns all configured profiles.
func ListProfiles(cfg *Config) []AccountProfile {
	return cfg.Profiles
}

// AddProfile appends a new profile to the configuration.
func AddProfile(cfg *Config, profile AccountProfile) {
	cfg.Profiles = append(cfg.Profiles, profile)
}
