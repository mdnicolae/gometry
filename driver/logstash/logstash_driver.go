package logstash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gometry/logging"
	"gometry/registry"
	"io"
	"net/http"
	"sync"
	"time"
)

type Driver struct {
	mu      sync.RWMutex
	traceID string
	url     string
	client  *http.Client
}

func init() {
	registry.RegisterDriver("logstash", func(config map[string]interface{}) (logging.TelemetryInstance, error) {
		url, ok := config["url"].(string)
		if !ok {
			return nil, fmt.Errorf("url option is required")
		}
		return &Driver{
			url:    url,
			client: &http.Client{},
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
	d.mu.RLock()
	defer d.mu.RUnlock()
	entry := logging.LogEntry{
		Timestamp:  time.Now(),
		Level:      level,
		Message:    msg,
		Attributes: logging.NormalizeAttributes(attributes),
		TraceID:    d.traceID,
	}
	go d.sendLog(entry)
	return nil
}

func (d *Driver) sendLog(entry logging.LogEntry) {
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Printf("Failed to marshal log entry: %v\n", err)
		return
	}
	req, err := http.NewRequest("POST", d.url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Failed to create HTTP request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := d.client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send log entry: %v\n", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Failed to close response body: %v\n", err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received non-OK response from Logstash: %s\n", resp.Status)
	}
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
