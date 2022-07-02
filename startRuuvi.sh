#!/bin/bash

sleep 15

hciattach /dev/ttyAMA0 bcm43xx 921600

#Set green led on to indicate that we are receiving measurements
#You can comment this out if you are not using RaspberryPi
echo 1 > /sys/class/leds/ACT/brightness

./ruuvi >/dev/null 2>&1

#Trun off green led to indicate that we are not running application anymore
#This will help to see if ruuvi application is not running anymore
#You can comment nexy row  if you are not using RaspberryPi
echo 0 > /sys/class/leds/ACT/brightness
