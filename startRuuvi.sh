#!/bin/bash

#Create directory for the logs if it does not exists
mkdir -p /var/log/ruuvi

sleep 15

hciattach /dev/ttyAMA0 bcm43xx 921600

/usr/local/bin/LED-act-on.sh

#./ruuvi >/dev/null 2>&1
./ruuvi 2>&1 |grep -va 'Ruuvi data with'|grep -va ' DATA: \[ '|grep -va 'INFO: Skipping'|grep -va ' MAC: ' >> /var/log/ruuvi/ruuvi.`date +%Y%m%d`.log 

/usr/local/bin/LED-act-off.sh
