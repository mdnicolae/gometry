package main

import (
	"fmt"
	"gometry"
	"gometry/logging"
	"os"
)

func main() {
	os.Setenv("INFLUXDB_URL", "http://localhost:8086")
	os.Setenv("INFLUXDB_TOKEN", "my-token")
	gmy, err := gometry.Init("my_custom_driver")
	if err != nil {
		fmt.Printf("Failed to initialize gmy: %v\n", err)
		return
	}

	defer func(logger logging.TelemetryInstance) {
		err := logger.Close()
		if err != nil {
			fmt.Printf("Failed to close gmy: %v\n", err)
		}
	}(gmy)

	//Now you can use your custom driver to log messages
	gmy.Info("This is an info message", map[string]interface{}{"user": "john_doe"})
	gmy.Debug("Debugging details", map[string]interface{}{"module": "auth"})
}
