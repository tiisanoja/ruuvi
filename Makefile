go = /usr/lib/go-1.15/bin/go
#go = go

build:
	$(go) build 
clean:
	rm ../bin/ruuvi
	rm ../bin/config.yml
env:
	$(go) get github.com/influxdata/influxdb1-client/v2
	$(go) get github.com/paypal/gatt
	$(go) get github.com/spf13/viper

