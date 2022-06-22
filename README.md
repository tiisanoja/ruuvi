# Ruuvi
Save ruuvi tag measurements to InfluxDB. You can use Graphana to view resukts from database.

## Building

This should be compiled with go 1.15 or later
Set env variable  GO111MODULE so that golang uses correct version of modules
export GO111MODULE=on

You need to have installed golang and make to be able to compile. First you need to add dependencies and only after that build application.
Steps:
1. make env
2. make build

This will create ruuvi binary.


## Executing

It is expected that config.yml is in the same directory with binary.
