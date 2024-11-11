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



