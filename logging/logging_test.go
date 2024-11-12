package logging

import (
	"reflect"
	"testing"
	"time"
)

type MockTelemetryInstance struct {
	entries []LogEntry
	traceID string
}

func (m *MockTelemetryInstance) log(level LogLevel, msg string, attributes ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp:  time.Now(),
		Level:      level,
		Message:    msg,
		Attributes: NormalizeAttributes(attributes),
		TraceID:    m.traceID,
	}
	m.entries = append(m.entries, entry)
}

func (m *MockTelemetryInstance) Info(msg string, attributes ...map[string]interface{}) {
	m.log(Info, msg, attributes...)
}

func (m *MockTelemetryInstance) Debug(msg string, attributes ...map[string]interface{}) {
	m.log(Debug, msg, attributes...)
}

func (m *MockTelemetryInstance) Warning(msg string, attributes ...map[string]interface{}) {
	m.log(Warning, msg, attributes...)
}

func (m *MockTelemetryInstance) Error(msg string, attributes ...map[string]interface{}) {
	m.log(Error, msg, attributes...)
}

func (m *MockTelemetryInstance) Critical(msg string, attributes ...map[string]interface{}) {
	m.log(Critical, msg, attributes...)
}

func (m *MockTelemetryInstance) StartTrace(traceID string) {
	m.traceID = traceID
}

func (m *MockTelemetryInstance) EndTrace() {
	m.traceID = ""
}

func (m *MockTelemetryInstance) Close() error {
	return nil
}

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{Critical, "CRITICAL"},
		{Error, "ERROR"},
		{Warning, "WARNING"},
		{Info, "INFO"},
		{Debug, "DEBUG"},
		{LogLevel(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestTelemetryInstance(t *testing.T) {
	mock := &MockTelemetryInstance{}

	mock.StartTrace("test-trace-id")
	mock.Info("info message", map[string]interface{}{"key": "value"})
	mock.Debug("debug message")
	mock.Warning("warning message")
	mock.Error("error message")
	mock.Critical("critical message")
	mock.EndTrace()

	if len(mock.entries) != 5 {
		t.Fatalf("expected 5 log entries, got %d", len(mock.entries))
	}

	tests := []struct {
		index     int
		level     LogLevel
		message   string
		attribute map[string]interface{}
	}{
		{0, Info, "info message", map[string]interface{}{"key": "value"}},
		{1, Debug, "debug message", map[string]interface{}{}},
		{2, Warning, "warning message", map[string]interface{}{}},
		{3, Error, "error message", map[string]interface{}{}},
		{4, Critical, "critical message", map[string]interface{}{}},
	}

	for _, tt := range tests {
		entry := mock.entries[tt.index]
		if entry.Level != tt.level {
			t.Errorf("expected level %s, got %s", tt.level.String(), entry.Level.String())
		}
		if entry.Message != tt.message {
			t.Errorf("expected message %s, got %s", tt.message, entry.Message)
		}
		if !reflect.DeepEqual(entry.Attributes, tt.attribute) {
			t.Errorf("expected attributes %v, got %v", tt.attribute, entry.Attributes)
		}
		if entry.TraceID != "test-trace-id" {
			t.Errorf("expected trace ID 'test-trace-id', got %s", entry.TraceID)
		}
	}
}
