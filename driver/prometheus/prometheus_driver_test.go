package prometheus

import (
	"github.com/mdnicolae/gometry/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"testing"
	_ "time"
)

func TestDriver_StartTrace(t *testing.T) {
	driver := &Driver{}
	traceID := "test-trace-id"
	driver.StartTrace(traceID)

	if driver.traceID != traceID {
		t.Errorf("expected traceID %s, got %s", traceID, driver.traceID)
	}
}

func TestDriver_EndTrace(t *testing.T) {
	driver := &Driver{traceID: "test-trace-id"}
	driver.EndTrace()

	if driver.traceID != "" {
		t.Errorf("expected traceID to be empty, got %s", driver.traceID)
	}
}

func TestDriver_Log(t *testing.T) {
	driver := &Driver{
		jobName:  "test-job",
		gateway:  "http://localhost:9091",
		metrics:  make(map[string]prometheus.Gauge),
		registry: prometheus.NewRegistry(),
	}

	attributes := []map[string]interface{}{
		{"execution_time": 1.23},
		{"cpu_load": 0.99},
	}

	err := driver.log(logging.Info, "test message", attributes)
	if err != nil {
		t.Fatalf("log method failed: %v", err)
	}

	metric := driver.metrics["execution_time"]
	if metric == nil {
		t.Fatalf("expected metric 'execution_time' to be registered")
	}

	metricValue := testutil.ToFloat64(metric)
	if metricValue != 1.23 {
		t.Errorf("expected metric value 1.23, got %f", metricValue)
	}
}

func TestDriver_Info(t *testing.T) {
	driver := &Driver{
		jobName:  "test-job",
		gateway:  "http://localhost:9091",
		metrics:  make(map[string]prometheus.Gauge),
		registry: prometheus.NewRegistry(),
	}

	driver.Info("test message", map[string]interface{}{"execution_time": 1.23})

	metric := driver.metrics["execution_time"]
	if metric == nil {
		t.Fatalf("expected metric 'execution_time' to be registered")
	}

	metricValue := testutil.ToFloat64(metric)
	if metricValue != 1.23 {
		t.Errorf("expected metric value 1.23, got %f", metricValue)
	}
}

func TestDriver_Debug(t *testing.T) {
	driver := &Driver{
		jobName:  "test-job",
		gateway:  "http://localhost:9091",
		metrics:  make(map[string]prometheus.Gauge),
		registry: prometheus.NewRegistry(),
	}

	driver.Debug("test message", map[string]interface{}{"cpu_load": 0.99})

	metric := driver.metrics["cpu_load"]
	if metric == nil {
		t.Fatalf("expected metric 'cpu_load' to be registered")
	}

	metricValue := testutil.ToFloat64(metric)
	if metricValue != 0.99 {
		t.Errorf("expected metric value 0.99, got %f", metricValue)
	}
}
