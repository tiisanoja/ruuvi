# Ruuvi
Saves ruuvi tag measurements to InfluxDB. You can use Graphana to view results from database. This code has been run on RaspberryPi. So it should work at least there.

Application stores Temperature (°C), Pressure (hPa), Humidity (%). It also calculates absolutely humidity (g/m2) and dew point (°C) and stores those to InfluxDB. Application listen only [RAWv2](https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-5-rawv2) format. Really old versions of RuuviTag might have still Data Format 3 which is not supported. There is also Data Format 8 which is encrypted version of data format. That is not supported right now. 

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

Before executing setup config.yml. You need to give name for sensors and provide MAC address of those sensors. Only those sensors are stored to db which MAC address is provided in config.yml. Each device is stored in 15s interval. True interval is some where 15s-17s because Ruuvi Tag is sending every 2s measurements.

It is expected that config.yml is in the same directory with binary.

## DB

Database name is weather. Presission is second.

