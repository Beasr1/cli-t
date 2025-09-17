# Makefile
BINARY_NAME=cli-t
VERSION?=0.1.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X github.com/yourusername/cli-t/pkg/version.Version=${VERSION} \
                  -X github.com/yourusername/cli-t/pkg/version.Commit=${COMMIT} \
                  -X github.com/yourusername/cli-t/pkg/version.BuildTime=${BUILD_TIME}"

.PHONY: all build clean test coverage lint install

all: clean lint test build

build:
	go build ${LDFLAGS} -o bin/${BINARY_NAME} cmd/cli-t/main.go

clean:
	go clean
	rm -rf bin/ coverage.out

test:
	go test -v -race ./...

coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint:
	golangci-lint run

install: build
	sudo mv bin/${BINARY_NAME} /usr/local/bin/

run: build
	./bin/${BINARY_NAME}

# Development helpers
dev-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go mod download

fmt:
	go fmt ./...

.DEFAULT_GOAL := build