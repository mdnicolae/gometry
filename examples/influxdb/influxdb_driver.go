package main

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/mdnicolae/gometry/logging"
	"github.com/mdnicolae/gometry/registry"
	"sync"
	"time"
)

type InfluxDriver struct {
	mu       sync.RWMutex
	traceID  string
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
	bucket   string
	org      string
}

func init() {
	registry.RegisterDriver("influxdb", func(config map[string]interface{}) (logging.TelemetryInstance, error) {
		url, ok := config["url"].(string)
		if !ok {
			return nil, fmt.Errorf("url option is required")
		}
		token, ok := config["token"].(string)
		if !ok {
			return nil, fmt.Errorf("token option is required")
		}
		bucket, ok := config["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("bucket option is required")
		}
		org, ok := config["org"].(string)
		if !ok {
			return nil, fmt.Errorf("org option is required")
		}
		client := influxdb2.NewClient(url, token)
		writeAPI := client.WriteAPIBlocking(org, bucket)
		return &InfluxDriver{
			client:   client,
			writeAPI: writeAPI,
			bucket:   bucket,
			org:      org,
		}, nil
	})
}

func (d *InfluxDriver) StartTrace(traceID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.traceID = traceID
}

func (d *InfluxDriver) EndTrace() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.traceID = ""
}

func (d *InfluxDriver) log(level logging.LogLevel, msg string, attributes []map[string]interface{}) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	tags := map[string]string{
		"traceID": d.traceID,
		"level":   level.String(),
	}
	fields := map[string]interface{}{
		"message": msg,
	}
	for _, attr := range attributes {
		for key, value := range attr {
			fields[key] = value
		}
	}

	point := influxdb2.NewPoint("log", tags, fields, time.Now())
	return d.writeAPI.WritePoint(context.Background(), point)
}

func (d *InfluxDriver) Info(msg string, attributes ...map[string]interface{}) {
	_ = d.log(logging.Info, msg, attributes)
}

func (d *InfluxDriver) Debug(msg string, attributes ...map[string]interface{}) {
	_ = d.log(logging.Debug, msg, attributes)
}

func (d *InfluxDriver) Warning(msg string, attributes ...map[string]interface{}) {
	_ = d.log(logging.Warning, msg, attributes)
}

func (d *InfluxDriver) Error(msg string, attributes ...map[string]interface{}) {
	_ = d.log(logging.Error, msg, attributes)
}

func (d *InfluxDriver) Critical(msg string, attributes ...map[string]interface{}) {
	_ = d.log(logging.Critical, msg, attributes)
}

func (d *InfluxDriver) Close() error {
	d.client.Close()
	return nil
}
