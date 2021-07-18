.PHONY: all

all:
	docker run --rm -v "$(shell pwd)":/usr/src/myapp -w /usr/src/myapp golang:1.16 go build -v
