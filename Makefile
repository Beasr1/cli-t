# Makefile
BINARY_NAME=cli-t
VERSION?=0.1.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X cli-t/pkg/version.Version=${VERSION} \
                  -X cli-t/pkg/version.Commit=${COMMIT} \
                  -X cli-t/pkg/version.BuildTime=${BUILD_TIME}"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Test parameters
TEST_TIMEOUT=10m
COVERAGE_THRESHOLD=70

.PHONY: all build clean test coverage lint install run fmt dev-deps test-unit test-integration test-bench test-verbose

all: clean lint test build

build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) ${LDFLAGS} -o bin/${BINARY_NAME} cmd/cli-t/main.go
	@echo "Build complete: bin/$(BINARY_NAME)"

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/ coverage.out coverage.html .test/

test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic -timeout $(TEST_TIMEOUT) ./...

test-unit:
	@echo "Running unit tests only..."
	$(GOTEST) -v -short -race -timeout $(TEST_TIMEOUT) ./...

test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -race -run Integration -timeout $(TEST_TIMEOUT) ./...

test-bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem -run=^# -timeout $(TEST_TIMEOUT) ./...

test-verbose:
	@echo "Running tests with verbose output..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic -timeout $(TEST_TIMEOUT) -args -test.v ./...

coverage: test
	@echo "Generating coverage report..."
	@$(GOCMD) tool cover -func=coverage.out | tail -1
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@# Check coverage threshold
	@bash -c 'COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk "{print \$$3}" | sed "s/%//"); \
		COVERAGE_INT=$${COVERAGE%.*}; \
		if [ "$$COVERAGE_INT" -lt "$(COVERAGE_THRESHOLD)" ]; then \
			echo "❌ Coverage $$COVERAGE% is below threshold $(COVERAGE_THRESHOLD)%"; \
			exit 1; \
		else \
			echo "✅ Coverage $$COVERAGE% meets threshold $(COVERAGE_THRESHOLD)%"; \
		fi'

lint:
	@echo "Running linter..."
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Install with: make dev-deps"; \
		echo "Running go vet instead..."; \
		$(GOCMD) vet ./...; \
	fi

fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	$(GOMOD) tidy

install: build
	@echo "Installing $(BINARY_NAME)..."
	@if [ -w /usr/local/bin ]; then \
		cp bin/$(BINARY_NAME) /usr/local/bin/; \
	else \
		sudo cp bin/$(BINARY_NAME) /usr/local/bin/; \
	fi
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"

run: build
	@echo "Running $(BINARY_NAME)..."
	./bin/$(BINARY_NAME)

# Development helpers
dev-deps:
	@echo "Installing development dependencies..."
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOCMD) install github.com/stretchr/testify/...@latest
	$(GOMOD) download
	@echo "Development dependencies installed!"

# CI/CD helpers
ci-test:
	@echo "Running CI tests..."
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic -timeout $(TEST_TIMEOUT) ./...

ci-lint:
	@echo "Running CI linting..."
	golangci-lint run --timeout=5m

# Quick checks
check: fmt lint test-unit
	@echo "✅ All checks passed!"

# Test specific packages
test-pkg:
	@if [ -z "$(PKG)" ]; then \
		echo "Usage: make test-pkg PKG=internal/command"; \
		exit 1; \
	fi
	@echo "Testing package: $(PKG)"
	$(GOTEST) -v -race ./$(PKG)/...

# Watch for changes and run tests
watch:
	@echo "Watching for changes..."
	@if command -v entr &> /dev/null; then \
		find . -name '*.go' | entr -c make test-unit; \
	else \
		echo "entr not installed. Install with: brew install entr (macOS) or apt-get install entr (Linux)"; \
		exit 1; \
	fi

# Generate test mocks (if needed later)
mocks:
	@echo "Generating mocks..."
	@echo "No mocks to generate yet"

# Show test coverage for specific package
cover-pkg:
	@if [ -z "$(PKG)" ]; then \
		echo "Usage: make cover-pkg PKG=internal/command"; \
		exit 1; \
	fi
	@echo "Coverage for package: $(PKG)"
	$(GOTEST) -coverprofile=coverage.out ./$(PKG)/...
	$(GOCMD) tool cover -func=coverage.out
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Create test report
test-report:
	@echo "Generating test report..."
	@mkdir -p .test
	$(GOTEST) -v -race -json ./... > .test/test-output.json
	@echo "Test report generated: .test/test-output.json"

.DEFAULT_GOAL := build

# Help target
help:
	@echo "CLI-T Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all            - Run clean, lint, test, and build"
	@echo "  build          - Build the binary"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run all tests with coverage"
	@echo "  test-unit      - Run unit tests only (fast)"
	@echo "  test-bench     - Run benchmarks"
	@echo "  coverage       - Generate coverage report"
	@echo "  lint           - Run golangci-lint"
	@echo "  fmt            - Format code"
	@echo "  install        - Install binary to /usr/local/bin"
	@echo "  run            - Build and run the binary"
	@echo "  dev-deps       - Install development dependencies"
	@echo "  check          - Run fmt, lint, and unit tests"
	@echo "  watch          - Watch for changes and run tests"
	@echo ""
	@echo "Examples:"
	@echo "  make test-pkg PKG=internal/command"
	@echo "  make cover-pkg PKG=internal/config"