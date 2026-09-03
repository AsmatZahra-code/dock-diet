// config.go loads the optional .dock-diet.yaml project configuration file.
//
// Supported YAML fields:
//
//	fail_under   int       Minimum acceptable Diet Score (default: 100).
//	                      The scan command exits with code 1 when the score
//	                      falls below this value, enabling CI/CD gating.
//	ignore_rules []string Reserved for future per-rule suppression support.
//
// If the configuration file is absent or contains invalid YAML, LoadConfig
// returns the default configuration silently so that pipelines that have not
// opted in are not accidentally broken.
package scanner

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of .dock-diet.yaml
type Config struct {
	FailUnder int      `yaml:"fail_under"`
	Ignore    []string `yaml:"ignore_rules"` // For future use (e.g., ignoring specific rules)
}

// LoadConfig reads the .dock-diet.yaml file if it exists
func LoadConfig() Config {
	// Default configuration — used when the file is absent, invalid, or
	// when individual fields are omitted (zero-value) in the YAML.
	defaultConfig := Config{
		FailUnder: 100, // Fail the pipeline if any issue is detected.
	}

	content, err := os.ReadFile(".dock-diet.yaml")
	if err != nil {
		// File doesn't exist — return defaults silently.
		return defaultConfig
	}

	var userConfig Config
	err = yaml.Unmarshal(content, &userConfig)
	if err != nil {
		// Invalid YAML — return defaults silently.
		return defaultConfig
	}

	// Apply defaults for any fields the user left unset (zero-value).
	// This ensures an empty file behaves identically to a missing file.
	if userConfig.FailUnder == 0 {
		userConfig.FailUnder = defaultConfig.FailUnder
	}

	return userConfig
}