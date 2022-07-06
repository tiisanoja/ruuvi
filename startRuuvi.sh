#!/bin/bash

#Create directory for the logs if it does not exists
mkdir -p /var/log/ruuvi

sleep 15

hciattach /dev/ttyAMA0 bcm43xx 921600

#Set green led on to indicate we are receiving measurements
#You can comment this out if you are not using RaspberryPi
echo 1 > /sys/class/leds/ACT/brightness

./ruuvi 2>&1 |grep -va 'Ruuvi data with'|grep -va ' DATA: \[ '|grep -va 'INFO: Skipping'|grep -va ' MAC: ' >> /var/log/ruuvi/ruuvi.`date +%Y%m%d`.log 

#Turn off green led to indicate that we are not running application anymore
#This will help to see if ruuvi application is not running anymore
#You can comment next row if you are not using RaspberryPi
echo 0 > /sys/class/leds/ACT/brightness
