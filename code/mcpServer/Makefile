.PHONY: test build run clean lint

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests for a specific package
test-db:
	go test -v ./pkg/db/...

test-types:
	go test -v ./pkg/types/...

# Build the application
build:
	go build -o bin/mcp-server ./cmd

# Run the application
run: build
	./bin/mcp-server

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f movieDB.db

# Run go fmt on all files
fmt:
	go fmt ./...

# Run go vet
lint:
	go vet ./...
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy

# Install dependencies
deps:
	go get -u ./...
	go mod tidy

# Run quick check (lint + test)
check: lint test

# Display help
help:
	@echo "Available targets:"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-db       - Run database tests only"
	@echo "  test-types    - Run types tests only"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  clean         - Remove build artifacts"
	@echo "  fmt           - Format all Go files"
	@echo "  lint          - Run linters (vet + fmt)"
	@echo "  tidy          - Tidy go.mod"
	@echo "  deps          - Update dependencies"
	@echo "  check         - Run lint + test"
