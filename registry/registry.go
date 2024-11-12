package registry

import (
	"fmt"
	"github.com/mdnicolae/gometry/config"
	"github.com/mdnicolae/gometry/logging"
	"sync"
)

type DriverFactory func(config map[string]interface{}) (logging.TelemetryInstance, error)

var defaultDriverIdentifier string

var (
	drivers          = make(map[string]DriverFactory)
	instanceRegistry = make(map[string]logging.TelemetryInstance)
	mu               sync.RWMutex
	initialized      bool
	initOnce         sync.Once
)

// RegisterDriver allows users to register a custom driver with a unique name.
func RegisterDriver(name string, factory DriverFactory) {
	mu.Lock()
	defer mu.Unlock()
	drivers[name] = factory
}

func loadConfig() error {
	configs, err := config.LoadConfig("gometry.json")
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	for _, cfg := range configs {
		factory, exists := drivers[cfg.Driver]
		if !exists {
			return fmt.Errorf("unsupported driver: %s", cfg.Driver)
		}

		instance, err := factory(cfg.Options)
		if err != nil {
			return fmt.Errorf("could not initialize driver %s: %v", cfg.Identifier, err)
		}

		mu.Lock()
		instanceRegistry[cfg.Identifier] = instance
		mu.Unlock()

		if cfg.Default {
			defaultDriverIdentifier = cfg.Identifier
		}
	}
	initialized = true
	return nil
}

// GetTelemetryInstance provides access to a telemetry instance by its identifier.
func GetTelemetryInstance(identifier ...string) (logging.TelemetryInstance, error) {
	initOnce.Do(func() {
		if !initialized {
			if err := loadConfig(); err != nil {
				fmt.Printf("Error loading config: %v\n", err)
			}
		}
	})

	id := defaultDriverIdentifier
	if len(identifier) > 0 {
		id = identifier[0]
	}

	mu.RLock()
	instance, ok := instanceRegistry[id]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("telemetry instance not found: %s", id)
	}
	return instance, nil
}
