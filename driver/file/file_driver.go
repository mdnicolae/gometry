package file

import (
	"fmt"
	"github.com/mdnicolae/gometry/logging"
	"github.com/mdnicolae/gometry/registry"
	"os"
	"sync"
	"time"
)

type Driver struct {
	mu      sync.RWMutex
	traceID string
	file    *os.File
}

func init() {
	registry.RegisterDriver("file", func(config map[string]interface{}) (logging.TelemetryInstance, error) {
		fileName, ok := config["file"].(string)
		if !ok {
			return nil, fmt.Errorf("file option is required")
		}
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("could not open file: %v", err)
		}
		return &Driver{file: file}, nil
	})
}

func (f *Driver) StartTrace(traceID string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.traceID = traceID
}

func (f *Driver) EndTrace() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.traceID = ""
}

func (f *Driver) log(level logging.LogLevel, msg string, attributes []map[string]interface{}) error {
	f.mu.RLock()
	defer f.mu.RUnlock()
	entry := logging.LogEntry{
		Timestamp:  time.Now(),
		Level:      level,
		Message:    msg,
		Attributes: logging.NormalizeAttributes(attributes),
		TraceID:    f.traceID,
	}
	logLine := fmt.Sprintf("[%s] %s: %s - %v [traceID: %s]\n", entry.Timestamp.Format(time.RFC3339Nano), level.String(), entry.Message, entry.Attributes, entry.TraceID)
	if _, err := f.file.WriteString(logLine); err != nil {
		return fmt.Errorf("could not write to file: %v", err)
	}
	return nil
}

func (f *Driver) Info(msg string, attributes ...map[string]interface{}) {
	_ = f.log(logging.Info, msg, attributes)
}

func (f *Driver) Debug(msg string, attributes ...map[string]interface{}) {
	_ = f.log(logging.Debug, msg, attributes)
}

func (f *Driver) Warning(msg string, attributes ...map[string]interface{}) {
	_ = f.log(logging.Warning, msg, attributes)
}

func (f *Driver) Error(msg string, attributes ...map[string]interface{}) {
	_ = f.log(logging.Error, msg, attributes)
}

func (f *Driver) Critical(msg string, attributes ...map[string]interface{}) {
	_ = f.log(logging.Critical, msg, attributes)
}

func (f *Driver) Close() error {
	return f.file.Close()
}
