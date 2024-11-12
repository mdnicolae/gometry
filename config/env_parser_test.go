package config

import (
	"os"
	"testing"
)

func TestExpandEnvVariables(t *testing.T) {
	err := os.Setenv("EXISTING_VAR", "value")
	if err != nil {
		t.Fatalf("could not set environment variable: %v", err)
	}
	defer func() {
		err := os.Unsetenv("EXISTING_VAR")
		if err != nil {
			t.Errorf("could not unset environment variable: %v", err)
		}
	}()

	tests := []struct {
		input    string
		expected string
	}{
		{"${EXISTING_VAR}", "value"},
		{"${NON_EXISTING_VAR}", ""},
		{"prefix_${EXISTING_VAR}_suffix", "prefix_value_suffix"},
		{"no_env_var", "no_env_var"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := expandEnvVariables(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestExpandConfigEnvVariables(t *testing.T) {
	// Set up environment variables for testing
	err := os.Setenv("EXISTING_VAR", "value")
	if err != nil {
		t.Fatalf("could not set environment variable: %v", err)
	}
	defer func() {
		err := os.Unsetenv("EXISTING_VAR")
		if err != nil {
			t.Errorf("could not unset environment variable: %v", err)
		}
	}()

	configs := []TelemetryConfig{
		{
			Options: map[string]interface{}{
				"option1": "${EXISTING_VAR}",
				"option2": "static_value",
				"option3": "${NON_EXISTING_VAR}",
			},
		},
	}

	expectedConfigs := []TelemetryConfig{
		{
			Options: map[string]interface{}{
				"option1": "value",
				"option2": "static_value",
				"option3": "",
			},
		},
	}

	result := ExpandConfigEnvVariables(configs)
	for i, cfg := range result {
		for key, value := range cfg.Options {
			if value != expectedConfigs[i].Options[key] {
				t.Errorf("expected %q for key %q, got %q", expectedConfigs[i].Options[key], key, value)
			}
		}
	}
}
