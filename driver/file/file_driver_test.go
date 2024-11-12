package file

import (
	"fmt"
	"github.com/mdnicolae/gometry/logging"
	"os"
	"strings"
	"testing"
)

func TestDriver_StartTrace(t *testing.T) {
	driver := &Driver{}
	traceID := "test-trace-id"
	driver.StartTrace(traceID)

	if driver.traceID != traceID {
		t.Errorf("expected traceID %q, got %q", traceID, driver.traceID)
	}
}

func TestDriver_EndTrace(t *testing.T) {
	driver := &Driver{traceID: "test-trace-id"}
	driver.EndTrace()

	if driver.traceID != "" {
		t.Errorf("expected traceID to be empty, got %q", driver.traceID)
	}
}

func TestDriver_LogMethods(t *testing.T) {
	tests := []struct {
		name       string
		logMethod  func(driver *Driver, msg string, attributes ...map[string]interface{})
		level      string
		msg        string
		attributes []map[string]interface{}
	}{
		{"Info", (*Driver).Info, "INFO", "info message", []map[string]interface{}{{"key": "value"}}},
		{"Debug", (*Driver).Debug, "DEBUG", "debug message", []map[string]interface{}{{"key": "value"}}},
		{"Warning", (*Driver).Warning, "WARNING", "warning message", []map[string]interface{}{{"key": "value"}}},
		{"Error", (*Driver).Error, "ERROR", "error message", []map[string]interface{}{{"key": "value"}}},
		{"Critical", (*Driver).Critical, "CRITICAL", "critical message", []map[string]interface{}{{"key": "value"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "logfile")
			if err != nil {
				t.Fatalf("could not create temp file: %v", err)
			}
			defer os.Remove(file.Name())

			driver := &Driver{file: file}
			tt.logMethod(driver, tt.msg, tt.attributes...)

			file.Seek(0, 0)
			content := make([]byte, 1024)
			n, _ := file.Read(content)
			logLine := string(content[:n])

			// Split the log line to ignore the timestamp
			logParts := strings.SplitN(logLine, " ", 2)
			if len(logParts) < 2 {
				t.Fatalf("log line format is incorrect: %q", logLine)
			}

			expectedLogLine := fmt.Sprintf("%s: %s - %v [traceID: ]\n", tt.level, tt.msg, logging.NormalizeAttributes(tt.attributes))
			if logParts[1] != expectedLogLine {
				t.Errorf("expected log line %q, got %q", expectedLogLine, logParts[1])
			}
		})
	}
}
