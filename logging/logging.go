package logging

import (
	"time"
)

type LogLevel int

const (
	Critical LogLevel = iota
	Error
	Warning
	Info
	Debug
)

type LogEntry struct {
	Timestamp  time.Time
	Level      LogLevel
	Message    string
	Attributes map[string]interface{}
	TraceID    string
}

type TelemetryInstance interface {
	Info(msg string, attributes ...map[string]interface{})
	Debug(msg string, attributes ...map[string]interface{})
	Warning(msg string, attributes ...map[string]interface{})
	Error(msg string, attributes ...map[string]interface{})
	Critical(msg string, attributes ...map[string]interface{})
	StartTrace(traceID string)
	EndTrace()
	Close() error
}

func (l LogLevel) String() string {
	switch l {
	case Critical:
		return "CRITICAL"
	case Error:
		return "ERROR"
	case Warning:
		return "WARNING"
	case Info:
		return "INFO"
	case Debug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}
