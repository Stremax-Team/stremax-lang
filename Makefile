.PHONY: build test clean run

# Build the stremax binary
build:
	go build -o stremax ./cmd/stremax

# Run all tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -f stremax
	rm -f coverage.out

# Run a stremax program
run:
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make run FILE=<file>"; \
		echo "Example: make run FILE=./examples/simple.sx"; \
		exit 1; \
	fi
	./stremax run -file $(FILE)

# Install stremax globally
install:
	go install ./cmd/stremax

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	@if command -v golint > /dev/null; then \
		golint ./...; \
	else \
		echo "golint not installed. Run: go install golang.org/x/lint/golint@latest"; \
	fi

# Vet code
vet:
	go vet ./...

# Run all code quality checks
check: fmt lint vet

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the stremax binary"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  clean          - Clean build artifacts"
	@echo "  run FILE=<file> - Run a stremax program"
	@echo "  install        - Install stremax globally"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  vet            - Vet code"
	@echo "  check          - Run all code quality checks"
	@echo "  help           - Show this help message" 