package main

import (
	"fmt"
	"github.com/mdnicolae/gometry"
	"github.com/mdnicolae/gometry/logging"
	"math/rand"
)

func main() {
	gmy, err := gometry.Init("file-driver")
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

	// Start a trace with random id
	gmy.StartTrace(fmt.Sprintf("x-trace-id-%d", rand.Intn(1000000)))

	// Log messages at different levels
	gmy.Info("This is an info message", map[string]interface{}{"user": "john_doe"})
	gmy.Debug("Debugging details", map[string]interface{}{"module": "auth"})
	gmy.Warning("This is a warning", map[string]interface{}{"disk_space": "low"})
	gmy.Error("An error occurred", map[string]interface{}{"error_code": 500})
	gmy.Critical("Critical error occurred", map[string]interface{}{"error_code": 500})

	gmy.EndTrace()

	// Log without trace
	gmy.Info("This is an info message")
}
