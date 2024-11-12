package config

import (
	"fmt"
	"os"
	"regexp"
)

var envVarRegex = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)}`)

// expandEnvVariables replaces placeholders like ${VAR_NAME} with the environment variable values.
func expandEnvVariables(value string) string {
	return envVarRegex.ReplaceAllStringFunc(value, func(envVar string) string {
		key := envVar[2 : len(envVar)-1] // Extract key from ${KEY}
		if val, exists := os.LookupEnv(key); exists {
			return val
		}
		fmt.Printf("Warning: environment variable %s not set\n", key)
		return ""
	})
}

func ExpandConfigEnvVariables(configs []TelemetryConfig) []TelemetryConfig {
	for i, cfg := range configs {
		for key, value := range cfg.Options {
			// Only process if value is a string
			if strVal, ok := value.(string); ok {
				// Expand environment variables within each option value
				cfg.Options[key] = expandEnvVariables(strVal)
			}
		}
		configs[i] = cfg
	}
	return configs
}
