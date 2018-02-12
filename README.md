# naturallight-strategy
This repo is for the new strategy server that will collect and process data from the car on the fly

## Requirements
- [golang](https://golang.org/) 1.9

## How to install and collect data
1. Run `make install` to install the `telemetry` command in `$GOPATH/bin`
2. Run `telemetry collect <run_name>` to collect data
3. Data will appear in `.blog` and `.csv` files in your current directory

## How to use
All relevant scripts should be defined in `Makefile` and can be run like `make test`, etc.
