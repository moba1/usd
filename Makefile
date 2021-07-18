.PHONY: all, build, clean

all: build

build:
	docker run --rm -v "$(shell pwd)":/usr/src/myapp -w /usr/src/myapp golang:1.16 go build -v

clean:
	rm -f usd
