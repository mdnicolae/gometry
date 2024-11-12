package prometheus

import (
	"fmt"
	"github.com/mdnicolae/gometry/logging"
	"github.com/mdnicolae/gometry/registry"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"sync"
)

type Driver struct {
	mu       sync.RWMutex
	traceID  string
	jobName  string
	gateway  string
	metrics  map[string]prometheus.Gauge
	registry *prometheus.Registry
}

func init() {
	registry.RegisterDriver("prometheus", func(config map[string]interface{}) (logging.TelemetryInstance, error) {
		jobName, ok := config["job_name"].(string)
		if !ok {
			return nil, fmt.Errorf("job_name option is required")
		}
		gateway, ok := config["gateway"].(string)
		if !ok {
			return nil, fmt.Errorf("gateway option is required")
		}
		return &Driver{
			jobName:  jobName,
			gateway:  gateway,
			metrics:  make(map[string]prometheus.Gauge),
			registry: prometheus.NewRegistry(),
		}, nil
	})
}

func (d *Driver) StartTrace(traceID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.traceID = traceID
}

func (d *Driver) EndTrace() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.traceID = ""
}

func (d *Driver) log(level logging.LogLevel, msg string, attributes []map[string]interface{}) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, attr := range attributes {
		for key, value := range attr {
			metricName := key
			if val, ok := value.(float64); ok {
				if _, exists := d.metrics[metricName]; !exists {
					d.metrics[metricName] = prometheus.NewGauge(prometheus.GaugeOpts{
						Name: metricName,
						Help: "Custom metric",
					})
					d.registry.MustRegister(d.metrics[metricName])
				}
				d.metrics[metricName].Set(val)

				go func(metric prometheus.Gauge, traceID, level string) {
					if err := push.New(d.gateway, d.jobName).
						Collector(metric).
						Grouping("traceID", traceID).
						Grouping("level", level).
						Push(); err != nil {
						fmt.Printf("could not push to gateway: %v\n", err)
					}
				}(d.metrics[metricName], d.traceID, level.String())
			}
		}
	}

	return nil
}

func (d *Driver) Info(msg string, attributes ...map[string]interface{}) {
	_ = d.log(logging.Info, msg, attributes)
}

func (d *Driver) Debug(msg string, attributes ...map[string]interface{}) {
	_ = d.log(logging.Debug, msg, attributes)
}

func (d *Driver) Warning(msg string, attributes ...map[string]interface{}) {
	_ = d.log(logging.Warning, msg, attributes)
}

func (d *Driver) Error(msg string, attributes ...map[string]interface{}) {
	_ = d.log(logging.Error, msg, attributes)
}

func (d *Driver) Critical(msg string, attributes ...map[string]interface{}) {
	_ = d.log(logging.Critical, msg, attributes)
}

func (d *Driver) Close() error {
	return nil
}
