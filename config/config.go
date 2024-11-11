package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// TelemetryConfig defines the structure for each logging configuration.
type TelemetryConfig struct {
	Identifier string                 `json:"identifier"`
	Driver     string                 `json:"driver"`
	Default    bool                   `json:"default"`
	Options    map[string]interface{} `json:"options"` // Options can contain environment variable placeholders
}

// LoadConfig loads and parses the configuration from a specified JSON file.
func LoadConfig(path string) ([]TelemetryConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %v", err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			fmt.Printf("could not close config file: %v\n", err)
		}
	}(file)

	var configs []TelemetryConfig
	if err := json.NewDecoder(file).Decode(&configs); err != nil {
		return nil, fmt.Errorf("could not decode config file: %v", err)
	}

	// Expand environment variables in the Options map
	configs = ExpandConfigEnvVariables(configs)

	for _, cfg := range configs {
		// Assert that the required fields are present
		if cfg.Identifier == "" || cfg.Driver == "" {
			return nil, fmt.Errorf("missing required fields in config: %+v", cfg)
		}
		ValidateRequiredOptionsForDriver(cfg)
	}

	return configs, nil
}

func ValidateRequiredOptionsForDriver(cfg TelemetryConfig) {
	switch cfg.Driver {
	case "file":
		if cfg.Options["file"] == "" {
			panic(fmt.Sprintf("missing required option 'file' for file driver: %+v", cfg))
		}
	}
}
