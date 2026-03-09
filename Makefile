#go = /usr/lib/go-1.19/bin/go
go = go

all: build

build: test	
	cd src; $(go) build -o ../bin/ruuvi
	cd src; env GOARCH=arm64 $(go) build -o ../bin/ruuvi.arm64
clean:
	rm bin/ruuvi
	rm bin/ruuvi.arm64
env:
	$(go) mod tidy
test:
	cd src; $(go) test

