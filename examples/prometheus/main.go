package main

import (
	"fmt"
	"gometry"
	"gometry/logging"
	"math/rand"
	"os"
	"time"
)

func main() {
	//Set os environment variable GOMETRY_CONFIG to the path of the configuration file
	os.Setenv("PROMETHEUS_GATEWAY", "http://localhost:9091")

	prom, err := gometry.Init("prom")
	if err != nil {
		fmt.Printf("Failed to initialize prom: %v\n", err)
		return
	}

	defer func(prom logging.TelemetryInstance) {
		err := prom.Close()
		if err != nil {
			fmt.Printf("Failed to close prom: %v\n", err)
		}
	}(prom)

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)

		prom.StartTrace(fmt.Sprintf("x-trace-id-%d", time.Second))

		prom.Info("execution_time", map[string]interface{}{"execution_time": rand.Float32()})
		prom.Debug("cpu_load", map[string]interface{}{"cpu_load": 0.99})
		prom.Warning("disk_space", map[string]interface{}{"disk_space": 0.55})
		prom.Error("error_count", map[string]interface{}{"error_count": 100})
		prom.Critical("critical_error", map[string]interface{}{"critical_errors": 10})

		prom.EndTrace()
	}
}
