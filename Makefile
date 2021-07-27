.PHONY: all, build, clean, test, coverage
go_command := docker run --rm -v "$(shell pwd)":/usr/src/myapp -w /usr/src/myapp golang:1.16 go
coverage_file := coverage.out
coverage_html := coverage.html

all: build

build:
	$(go_command) build -v

test:
	$(go_command) test ./...

coverage:
	$(go_command) test -coverprofile="$(coverage_file)" ./...
	$(go_command) tool cover -html="$(coverage_file)" -o "$(coverage_html)"

clean:
	git clean -fdX .

lint:
	docker run --rm -it -v "$(shell pwd)":/app -w /app golangci/golangci-lint golangci-lint run -v
