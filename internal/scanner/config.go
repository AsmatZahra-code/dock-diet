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
	// Default configuration
	defaultConfig := Config{
		FailUnder: 100, // By default, fail if score is below 100 (any issue)
	}

	content, err := os.ReadFile(".dock-diet.yaml")
	if err != nil {
		// If file doesn't exist, return default config silently
		return defaultConfig
	}

	var userConfig Config
	err = yaml.Unmarshal(content, &userConfig)
	if err != nil {
		// If YAML is invalid, print warning and return default
		return defaultConfig
	}

	return userConfig
}