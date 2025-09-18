#!/bin/bash
# test/run-tests.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test coverage threshold
COVERAGE_THRESHOLD=70

echo -e "${BLUE}=== CLI-T Test Suite ===${NC}\n"

# Run go mod tidy first
echo -e "${YELLOW}Running go mod tidy...${NC}"
go mod tidy

# Run linter
echo -e "\n${YELLOW}Running linter...${NC}"
if command -v golangci-lint &> /dev/null; then
    golangci-lint run
else
    echo -e "${YELLOW}golangci-lint not installed, skipping...${NC}"
fi

# Run tests with coverage
echo -e "\n${YELLOW}Running tests with coverage...${NC}"
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Check test results
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
else
    echo -e "${RED}✗ Tests failed!${NC}"
    exit 1
fi

# Display coverage report
echo -e "\n${YELLOW}Coverage Report:${NC}"
go tool cover -func=coverage.out | tail -1

# Check coverage threshold
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
COVERAGE_INT=${COVERAGE%.*}

if [ "$COVERAGE_INT" -lt "$COVERAGE_THRESHOLD" ]; then
    echo -e "${RED}✗ Coverage ${COVERAGE}% is below threshold ${COVERAGE_THRESHOLD}%${NC}"
    exit 1
else
    echo -e "${GREEN}✓ Coverage ${COVERAGE}% meets threshold ${COVERAGE_THRESHOLD}%${NC}"
fi

# Run benchmarks (optional)
if [ "$1" == "--bench" ]; then
    echo -e "\n${YELLOW}Running benchmarks...${NC}"
    go test -bench=. -benchmem ./...
fi

# Generate HTML coverage report (optional)
if [ "$1" == "--html" ]; then
    echo -e "\n${YELLOW}Generating HTML coverage report...${NC}"
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}Coverage report generated: coverage.html${NC}"
fi

echo -e "\n${GREEN}=== Test Suite Complete ===${NC}"