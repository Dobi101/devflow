.PHONY: run test build fmt clean check

run:
	GOCACHE=/tmp/go-build go run ./cmd/devflow

test:
	GOCACHE=/tmp/go-build go test ./...

build:
	mkdir -p bin
	GOCACHE=/tmp/go-build go build -buildvcs=false -o bin/devflow ./cmd/devflow

fmt:
	gofmt -w cmd internal

clean:
	rm -rf bin

check: fmt test build
