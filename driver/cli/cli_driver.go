package cli

import (
	"fmt"
	"github.com/mdnicolae/gometry/logging"
	"github.com/mdnicolae/gometry/registry"
	"sync"
	"time"
)

type Driver struct {
	mu      sync.RWMutex
	traceID string
	colors  bool
}

func init() {
	registry.RegisterDriver("cli", func(config map[string]interface{}) (logging.TelemetryInstance, error) {
		colors, ok := config["colors"].(bool)
		if !ok {
			colors = true // Default to true if not specified
		}
		return &Driver{colors: colors}, nil
	})
}

func (c *Driver) StartTrace(traceID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.traceID = traceID
}

func (c *Driver) EndTrace() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.traceID = ""
}

func (c *Driver) log(level logging.LogLevel, msg string, attributes []map[string]interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry := logging.LogEntry{
		Timestamp:  time.Now(),
		Level:      level,
		Message:    msg,
		Attributes: logging.NormalizeAttributes(attributes),
		TraceID:    c.traceID,
	}
	color := c.getColor(level)
	if len(entry.Attributes) == 0 {
		fmt.Printf("%s[%s] %s: %s [traceID: %s]%s\n", color, entry.Timestamp.Format(time.RFC3339Nano), level.String(), entry.Message, entry.TraceID, resetColor)
		return nil
	}
	fmt.Printf("%s[%s] %s: %s - %v [traceID: %s]%s\n", color, entry.Timestamp.Format(time.RFC3339Nano), level.String(), entry.Message, entry.Attributes, entry.TraceID, resetColor)
	return nil
}

func (c *Driver) getColor(level logging.LogLevel) string {
	if c.colors == false {
		return ""
	}
	switch level {
	case logging.Info:
		return "\033[37m" // White
	case logging.Debug:
		return "\033[34m" // Blue
	case logging.Warning:
		return "\033[33m" // Yellow
	case logging.Error:
		return "\033[35m" // Magenta
	case logging.Critical:
		return "\033[31m" // Red
	default:
		return "\033[37m" // Default to white
	}
}

const resetColor = "\033[0m"

func (c *Driver) Info(msg string, attributes ...map[string]interface{}) {
	_ = c.log(logging.Info, msg, attributes)
}

func (c *Driver) Debug(msg string, attributes ...map[string]interface{}) {
	_ = c.log(logging.Debug, msg, attributes)
}

func (c *Driver) Warning(msg string, attributes ...map[string]interface{}) {
	_ = c.log(logging.Warning, msg, attributes)
}

func (c *Driver) Error(msg string, attributes ...map[string]interface{}) {
	_ = c.log(logging.Error, msg, attributes)
}

func (c *Driver) Critical(msg string, attributes ...map[string]interface{}) {
	_ = c.log(logging.Critical, msg, attributes)
}

func (c *Driver) Close() error {
	return nil
}
