#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Running Stremax-Lang Test Suite${NC}"
echo "=================================="

# Format check
echo -e "\n${YELLOW}Checking code formatting...${NC}"
if [ "$(gofmt -l . | wc -l)" -gt 0 ]; then
  echo -e "${RED}The following files need formatting:${NC}"
  gofmt -l .
  echo -e "${YELLOW}Running gofmt to fix formatting issues...${NC}"
  gofmt -w .
  echo -e "${GREEN}Formatting fixed.${NC}"
else
  echo -e "${GREEN}All files are properly formatted.${NC}"
fi

# Lint check
echo -e "\n${YELLOW}Running linter...${NC}"
if command -v golint > /dev/null; then
  golint ./...
else
  echo -e "${RED}golint not installed. Run: go install golang.org/x/lint/golint@latest${NC}"
fi

# Vet check
echo -e "\n${YELLOW}Running go vet...${NC}"
go vet ./...

# Unit tests
echo -e "\n${YELLOW}Running unit tests...${NC}"
go test -v ./...

# Tests with race detection
echo -e "\n${YELLOW}Running tests with race detection...${NC}"
go test -race ./...

# Tests with coverage
echo -e "\n${YELLOW}Running tests with coverage...${NC}"
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Benchmarks
echo -e "\n${YELLOW}Running benchmarks...${NC}"
go test -bench=. -benchmem ./...

# Build check
echo -e "\n${YELLOW}Building the project...${NC}"
go build -v ./cmd/stremax

# Run examples
echo -e "\n${YELLOW}Running example programs...${NC}"
if [ -f "./stremax" ]; then
  echo -e "${YELLOW}Running simple.sx${NC}"
  ./stremax run -file ./examples/simple.sx
  
  echo -e "\n${YELLOW}Running arithmetic.sx${NC}"
  ./stremax run -file ./examples/arithmetic.sx
  
  echo -e "\n${YELLOW}Running strings.sx${NC}"
  ./stremax run -file ./examples/strings.sx
  
  echo -e "\n${YELLOW}Running conditionals.sx${NC}"
  ./stremax run -file ./examples/conditionals.sx
  
  echo -e "\n${YELLOW}Running boolean.sx${NC}"
  ./stremax run -file ./examples/boolean.sx
  
  echo -e "\n${YELLOW}Running combined.sx${NC}"
  ./stremax run -file ./examples/combined.sx
else
  echo -e "${RED}stremax binary not found. Build failed?${NC}"
fi

echo -e "\n${GREEN}All tests completed!${NC}" 