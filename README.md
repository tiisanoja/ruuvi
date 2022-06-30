# Ruuvi
Project provides application to store RuuviTag measurements to InfluxDB. You can then use for example Graphana to view results from database. This application has been tested on RaspberryPi. So it should work at least there.

Following values are stored to database. Stored values are mainly taken from sensor. There are few values which are calculated based on measuremants. Calculated values are marked with **`Calculated`** -tag.

### Weather
* Temperature (°C)
* Pressure (hPa)
* Humidity (%)
* Absolutely humidity (g/m2) **`Calculated`**
* Dew point (°C) **`Calculated`**

### Movement
* Acceleration (x,y,z) (mG)

### Hardware
* Battery voltage (mV)
* Transmit power (dBm)

Application listen only [RAWv2](https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-5-rawv2) format. Really old versions of RuuviTag might have still Data Format 3 which is not supported. There is also Data Format 8 which is encrypted version of data format. That is not supported right now. 

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

Before executing binary setup config.yml. Provide URL for InfluxDB. You need also to give name for sensors and provide MAC address. Only those sensors are stored to db which MAC address is provided in config.yml. Each RuuviTag sensor is stored in 15s interval by defaul. True interval is some where 15s-17s because Ruuvi Tag is sending every 2s measurements. In config.yml you can specify interval if something else is needed.

Config.yml needs to be in the same directory with binary.

## DB

Data is stored to InfluxDB. Currently supported version is 1.x. Application stores data to database named *weather*. Used presission to store measurements is a second.

