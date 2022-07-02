go = /usr/lib/go-1.18/bin/go
#go = go

build:
	$(go) build -o ../../bin/ruuvi
clean:
	rm ../../bin/ruuvi
	rm ../../bin/config.yml
env:
	$(go) mod tidy

