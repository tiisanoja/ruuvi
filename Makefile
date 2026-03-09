#go = /usr/lib/go-1.19/bin/go
go = go

all: build

build: test
	$(go) build -o bin/ruuvi
	env GOARCH=arm64 $(go) build -o bin/ruuvi.arm64
clean:
	rm bin/ruuvi
	rm bin/ruuvi.arm64
env:
	$(go) mod tidy
test:
	$(go) test

