package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	err := w.Close()
	if err != nil {
		return ""
	}
	os.Stdout = old
	_, err = buf.ReadFrom(r)
	if err != nil {
		return ""
	}
	return buf.String()
}

func TestDriver_Info(t *testing.T) {
	driver := &Driver{colors: false}
	output := captureOutput(func() {
		driver.Info("Test info message")
	})

	expectedKeywords := "INFO: Test info message [traceID: ]"
	if !strings.Contains(output, expectedKeywords) {
		t.Errorf("expected %q, got %q", expectedKeywords, output)
	}
}

func TestDriver_Debug(t *testing.T) {
	driver := &Driver{colors: false}
	output := captureOutput(func() {
		driver.Debug("Test debug message")
	})

	expectedKeywords := "DEBUG: Test debug message [traceID: ]"
	if !strings.Contains(output, expectedKeywords) {
		t.Errorf("expected %q, got %q", expectedKeywords, output)
	}
}

func TestDriver_Warning(t *testing.T) {
	driver := &Driver{colors: false}
	output := captureOutput(func() {
		driver.Warning("Test warning message")
	})

	expectedKeywords := "WARNING: Test warning message [traceID: ]"
	if !strings.Contains(output, expectedKeywords) {
		t.Errorf("expected %q, got %q", expectedKeywords, output)
	}
}

func TestDriver_Error(t *testing.T) {
	driver := &Driver{colors: false}
	output := captureOutput(func() {
		driver.Error("Test error message")
	})

	expectedKeywords := "ERROR: Test error message [traceID: ]"
	if !strings.Contains(output, expectedKeywords) {
		t.Errorf("expected %q, got %q", expectedKeywords, output)
	}
}

func TestDriver_Critical(t *testing.T) {
	driver := &Driver{colors: false}
	output := captureOutput(func() {
		driver.Critical("Test critical message")
	})

	expectedKeywords := "CRITICAL: Test critical message [traceID: ]"
	if !strings.Contains(output, expectedKeywords) {
		t.Errorf("expected %q, got %q", expectedKeywords, output)
	}
}

func TestDriver_InfoWithColor(t *testing.T) {
	driver := &Driver{colors: true}
	output := captureOutput(func() {
		driver.Info("Test info message")
	})

	expectedKeywords := "INFO: Test info message [traceID: ]"
	if !strings.HasPrefix(output, "\033[37m") {
		t.Errorf("expected color white, got %q", output)
	}
	if !strings.Contains(output, expectedKeywords) {
		t.Errorf("expected %q, got %q", expectedKeywords, output)
	}
}

func TestDriver_DebugWithColor(t *testing.T) {
	driver := &Driver{colors: true}
	output := captureOutput(func() {
		driver.Debug("Test debug message")
	})

	expectedKeywords := "DEBUG: Test debug message [traceID: ]"
	if !strings.HasPrefix(output, "\033[34m") {
		t.Errorf("expected color blue, got %q", output)
	}
	if !strings.Contains(output, expectedKeywords) {
		t.Errorf("expected %q, got %q", expectedKeywords, output)
	}
}

func TestDriver_WarningWithColor(t *testing.T) {
	driver := &Driver{colors: true}
	output := captureOutput(func() {
		driver.Warning("Test warning message")
	})

	expectedKeywords := "WARNING: Test warning message [traceID: ]"
	if !strings.HasPrefix(output, "\033[33m") {
		t.Errorf("expected color yellow, got %q", output)
	}
	if !strings.Contains(output, expectedKeywords) {
		t.Errorf("expected %q, got %q", expectedKeywords, output)
	}
}

func TestDriver_ErrorWithColor(t *testing.T) {
	driver := &Driver{colors: true}
	output := captureOutput(func() {
		driver.Error("Test error message")
	})

	expectedKeywords := "ERROR: Test error message [traceID: ]"
	if !strings.HasPrefix(output, "\033[35m") {
		t.Errorf("expected color magenta, got %q", output)
	}
	if !strings.Contains(output, expectedKeywords) {
		t.Errorf("expected %q, got %q", expectedKeywords, output)
	}
}

func TestDriver_CriticalWithColor(t *testing.T) {
	driver := &Driver{colors: true}
	output := captureOutput(func() {
		driver.Critical("Test critical message")
	})

	expectedKeywords := "CRITICAL: Test critical message [traceID: ]"
	if !strings.HasPrefix(output, "\033[31m") {
		t.Errorf("expected color red, got %q", output)
	}
	if !strings.Contains(output, expectedKeywords) {
		t.Errorf("expected %q, got %q", expectedKeywords, output)
	}
}

func TestDriver_Close(t *testing.T) {
	driver := &Driver{}
	err := driver.Close()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
