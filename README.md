# Ruuvi
Application stores RuuviTag measurements to InfluxDB. You can then use for example Graphana to view results from database. This application has been tested on RaspberryPi. So it should work at least there.

Following values are stored to database. Stored values are mainly taken message sent from sensor. There are few values which are calculated based on measuremants. Calculated values are marked with **`Calculated`** -tag.

### Weather
* Temperature (°C)
* Pressure (hPa)
* Humidity (%)
* Absolutely humidity (g/m³) **`Calculated`** (**\***)
* Dew point (°C) **`Calculated`** (**\*\***)
* Wet Bulb Temperature (°C) **`Calculated`** (Development on going)

(**\***) Absolutely humidity approximation is calculated based on temperature and humidity%.
Absolutely humidity is calculated using Bolton formula for steam saturated pressure. Formula can be found [here](https://carnotcycle.wordpress.com/2012/08/04/how-to-convert-relative-humidity-to-absolute-humidity/comment-page-1/).
***Note!*** Returned value is approximation. See links for detail for error. Also measurements has error which are effecting to result of approximation of absolutely humidity. This calculated value is without any warranty!
 
(**\*\***) Dew point is calculated using a well-known approximation formula found in [Wikipedia](https://en.wikipedia.org/wiki/Dew_point). Formula can be found below "Calculating the dew point".
b and c values used in the formula are:  b = 17.62, c = 243.12°C
***Note!*** Result is approximation. See links for detail for error. Also measurements has error which are effecting to result of approximation of dew point. More details can be found from [here](https://en.wikipedia.org/wiki/Dew_point) This calculated value is without any warranty!

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

This will first run unit tests and then create ruuvi binary to directory ../../bin. If unit tests are not passed, binary is not created. Check Makefile and change it so that go -compiler can be found.


## Running

Before running the binary, setup config.yml. Provide database settings for InfluxDB. In config.yml there are examples for InfluxDB 1.8 and 2.x. Please note that those two versions are configured totally different way. You need also to give name for sensors and provide MAC addresses of RuuviTags. Only those sensors are stored to the database which MAC address is provided in config.yml. Each RuuviTag sensor is stored in 15s interval by default to save some disk space. True interval is bit more because RuuviTag is sending measurements between 1-3s depending on configuration and firmware installed to RuuviTag. In config.yml you can specify storing interval if something else than 15s is needed.

There is startRuuvi.sh which can be used at least on RaspberryPi to start application. It will trun green led on when application is runnig. You can comment that part from the script if you do not want that functionality.

Config.yml needs to be in the same directory with binary.

### Example
1. mkdir -p /opt/ruuvi
2. Copy config.yml, ruuvi and startRuuvi.sh to /opt/ruuvi directory
3. cd /opt/ruuvi
4. ./startRuuvi.sh

Error log is generated to /var/log/ruuvi directory. It will use starting day as part of the log file name (ruuvi.<date in form of YYYYMMDD>.log). Logs do not rollover except if you start application daily.

## Database

Data is stored to InfluxDB. Supported version by used client are 1.8 and 2.x. Application stores data to bucket, which is configured in config.yml. Default bucket is *weather*. Used presission in a databse to store measurements is a second. Application has been tested against InfluxDB 1.8 but now on only InfluxDB 2.X will be verified. InfluxDB 1.8 should work as long as used client supports 1.8.

## Grafana
 
Grafana can be used to present measurements from database. It has good support for InfluxDB. See more from [here](https://grafana.com/oss/grafana/).
