# Makefile

.PHONY: build run clean

build:
	go build -o bart ./cmd/bart

run: build
	./bart

clean:
	go clean
	rm -f bart

test:
	go test ./... -v

fmt:
	go fmt ./...