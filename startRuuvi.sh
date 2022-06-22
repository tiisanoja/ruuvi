#!/bin/bash

sleep 15

hciattach /dev/ttyAMA0 bcm43xx 921600

/usr/local/bin/LED-act-on.sh
./ruuvi >/dev/null 2>&1
/usr/local/bin/LED-act-off.sh
