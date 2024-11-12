# gometry - A telemetry package for Golang lovers

[![GoDoc](https://godoc.org/github.com/mdnicolae/gometry?status.svg)](https://godoc.org/github.com/mdnicolae/gometry)
[![Build Status](https://app.travis-ci.com/mdnicolae/gometry.svg?token=xMRafMALnLewDSGZzBGz)](https://app.travis-ci.com/mdnicolae/gometry)
[![Coverage Status](https://coveralls.io/repos/github/mdnicolae/gometry/badge.svg?branch=main)](https://coveralls.io/github/mdnicolae/gometry?branch=main)


## Installation

```bash
go get github.com/mdnicolae/gometry
```

## Usage

To use this package, you need to have a config file 
where you can specify the drivers you want to use and
the configuration for each driver.

The file needs to be under the root of your project and be named `gometry.json`.

Here is an example of a config file:

```json
[
  {
    "identifier": "cli-driver-identifier",
    "driver": "cli"
  },
  {
    "identifier": "file-driver-identifier",
    "driver": "file",
    "options": {
      "file": "app.log"
    }
  }
]
```

## Drivers

A driver is a type of output where the telemetry data is sent. Identifiers are used to differentiate between the different drivers.
There can be defined as many drivers as you want in the config file, but each driver needs to have a unique identifier. 
You can have multiple drivers of the same type, with different configurations.

By default, the package comes with the following implemented drivers:

### CLI Driver
The CLI driver logs the telemetry data to the console. By default,
it uses colors to differentiate between the different log levels.

If you want, you can disable this with a config like:

```json
{
  "identifier": "cli-no-colors",
  "driver": "cli",
  "options": {
    "colors": false
  }
}
```

For more example usage, please check the `examples > cli ` folder.

### File Driver

The file driver logs the data to a file. The file is created if it does not exist and is appended to if it does.

Config structure:
    
```json
{
    "identifier": "file-driver-identifier",
    "driver": "file",
    "options": {
        "file": "app.log"
    }
}
```

For more example usage, please check the `examples > file ` folder.

### Prometheus Driver

The Prometheus driver is used to send metrics to a Prometheus Gateway. For now, all metrics are sent as Gauge type.

Config structure:

```json
{
  "identifier": "prom",
  "driver": "prometheus",
  "options": {
    "job_name": "gometry_job",
    "gateway": "${PROMETHEUS_GATEWAY}"
  }
}
```
For a simple running example, you can check the `examples > prometheus` folder.
> Note: In order to be able to run this on your local machine, you must have docker-compose installed.
> You can start the Prometheus Gateway by running `docker-compose up` in the `examples > prometheus` folder. 
> To see the metrics, you can access Grafana (http://localhost:3000) in your browser and configure it to use the
> http://prometheus:9090 source, then run the example.


## Custom Drivers

Good news is, if you want to run this package with a custom driver, you can do that, without doing any changes to the package itself.

You just simply need to register a new driver type that implements the necessary methods.

You can see the example in the `examples > influxdb > indluxdb_driver.go` folder where I implemented a custom driver that sends the telemetry data to an InfluxDB instance.

Then, you will be able to use it like it's shown in the `main.go` file in the same folder.

For this driver, the config should look like:
    
```json
  {
    "identifier": "my_custom_driver",
    "driver": "influxdb",
    "options": {
      "url": "${INFLUXDB_URL}",
      "token": "${INFLUXDB_TOKEN}",
      "bucket": "gometry",
      "org": "gometry"
    }
  }
```


## Default Drivers

If you want, you can always select a driver as being your default one.

In this `gometry.json` example we have 3 different drivers defined:

```json
[
  {
    "identifier": "cli-no-color",
    "driver": "cli",
    "options": {
      "colors": false
    }
  },
  {
    "identifier": "cli-driver",
    "driver": "cli",
    "default": true
  },
  {
    "identifier": "file-driver",
    "driver": "file",
    "options": {
      "file": "app-logs.log"
    }
  }
]
```

However, because the second one has the `default` field set to `true`, it will be the one used by default.
You will be able to call the package with:

```go
package main

import (
	"gometry"
)

func main() {
	//If identifier will not be provided, it will use the one marked as default in the gometry.json file
	gmy, _ := gometry.Init()
    gmy.Info("This is an info message", map[string]interface{}{"user": "john_doe"})
}
```

This is helpful when you want to easily change the default driver without being forced to change the actual code.

