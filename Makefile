#Golang 1.17 or later is required
#go = /usr/lib/go-1.19/bin/go
go = go

all: test build

build:
	$(go) build -o ../../bin/ruuvi
	env GOARCH=arm64 $(go) build -o ../../bin/ruuvi.arm64
clean:
	rm ../../bin/ruuvi
	rm ../../bin/ruuvi.arm64
	rm ../../bin/config.yml
env:
	$(go) mod tidy
test:
	$(go) test

