# Ruuvi
Project provides application to store RuuviTag measurements to InfluxDB. You can then use for example Graphana to view results from database. This application has been tested on RaspberryPi. So it should work at least there.

Following values are stored to database. Stored values are mainly taken message sent from sensor. There are few values which are calculated based on measuremants. Calculated values are marked with **`Calculated`** -tag.

### Weather
* Temperature (°C)
* Pressure (hPa)
* Humidity (%)
* Absolutely humidity (g/m2) **`Calculated`** See Note 1.
* Dew point (°C) **`Calculated`** See Note 2.

Note 1: http://www.finwx.net/forum/index.php?topic=390.0

Note 2: Dew point is calculated using formula found in https://en.wikipedia.org/wiki/Dew_point. A well-known approximation formula is used to calculate the dew point. Formula can be found below "Calculating the dew point".
b and c values used in the formula are:  b = 17.62, c = 243.12°C
Result is approximation. Amount of error is unclear because it is at least mixture of usage well-known approximation formula and these b and c values has some error depending on the measured temperature. More details can be found from: https://en.wikipedia.org/wiki/Dew_point


### Movement
* Acceleration (x,y,z) (mG)

### Hardware
* Battery voltage (mV)
* Transmit power (dBm)

Application listen only [RAWv2](https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-5-rawv2) format. Really old versions of RuuviTag might have still Data Format 3 which is not supported. There is also Data Format 8 which is encrypted version of data format. That is not supported right now.

## Building

Building requires module support from golang. Building requires at least golang version 1.17. Building requires that you have installed make and golang. You might need to change Makefile to specify where go binary can be found.

Building:
1. make

This will first run unit tests and then create ruuvi binary to directory ../../bin. If unit tests are not passed, binary is not created.


## Running

Before running the binary, setup config.yml. Provide database settings for InfluxDB. In config.yml there are examples for InfluxDB 1.8 and 2.x. Please note that those two versions are configured totally different way. You need also to give name for sensors and provide MAC addresses of RuuviTags. Only those sensors are stored to the database which MAC address is provided in config.yml. Each RuuviTag sensor is stored in 15s interval by default to save some disk space. True interval is bit more because RuuviTag is sending measurements between 1-3s depending on configuration and firmware installed to RuuviTag. In config.yml you can specify storing interval if something else than 15s is needed.

There is startRuuvi.sh which can be used at least on RaspberryPi to start application. It will trun green led on when application is runnig. You can comment that part from the script if you do not want that functionality.

Config.yml needs to be in the same directory with binary.

### Example
1. mkdir -p /opt/ruuvi
2. Copy config.yml, ruuvi and startRuuvi.sh to /opt/ruuvi directory
3. cd /opt/ruuvi
4. ./startRuuvi.sh

Error log is generated to /var/log/ruuvi directory. It will use starting day as part of the log file name (ruuvi.<date in form of YYYYMMDD>.log). Logs do not rollover execpt if you start application daily.

## DB

Data is stored to InfluxDB. Supported version are 1.8 and 2.x. Application stores data to bucket which is configured in config.yml. Default bucket is *weather*. Used presission to store measurements is a second.

