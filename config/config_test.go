package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configContent := `[
		{
			"identifier": "test",
			"driver": "file",
			"default": true,
			"options": {
				"file": "test.log"
			}
		}
	]`
	tmpFile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("could not remove temp file: %v", err)
		}
	}(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatalf("could not write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("could not close temp file: %v", err)
	}

	configs, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(configs) != 1 {
		t.Fatalf("expected 1 config, got %d", len(configs))
	}

	cfg := configs[0]
	if cfg.Identifier != "test" {
		t.Errorf("expected identifier 'test', got '%s'", cfg.Identifier)
	}
	if cfg.Driver != "file" {
		t.Errorf("expected driver 'file', got '%s'", cfg.Driver)
	}
	if !cfg.Default {
		t.Errorf("expected default to be true, got false")
	}
	if cfg.Options["file"] != "test.log" {
		t.Errorf("expected file option 'test.log', got '%s'", cfg.Options["file"])
	}
}
