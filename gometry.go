package gometry

import (
	_ "gometry/driver/cli"
	_ "gometry/driver/file"
	_ "gometry/driver/prometheus"
	"gometry/logging"
	"gometry/registry"
)

// Init initializes and returns a TelemetryInstance for the given identifier.
func Init(identifier ...string) (logging.TelemetryInstance, error) {
	return registry.GetTelemetryInstance(identifier...)
}
