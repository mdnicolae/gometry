package gometry

import (
	_ "github.com/mdnicolae/gometry/driver/cli"
	_ "github.com/mdnicolae/gometry/driver/file"
	_ "github.com/mdnicolae/gometry/driver/prometheus"
	"github.com/mdnicolae/gometry/logging"
	"github.com/mdnicolae/gometry/registry"
)

// Init initializes and returns a TelemetryInstance for the given identifier.
func Init(identifier ...string) (logging.TelemetryInstance, error) {
	return registry.GetTelemetryInstance(identifier...)
}
