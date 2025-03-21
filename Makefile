# Variables
BINARY_NAME=gobuild
BUILD_DIR=build

# Default target
all: build test clean

# Build the project
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gobuild

# Run tests
test:
	go test ./...

# Clean up build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Run the application
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Install dependencies
deps:
	go mod tidy

# Format the code
fmt:
	go fmt ./...

# Lint the code
lint:
	golangci-lint run

dploy:
	./$(BUILD_DIR)/$(BINARY_NAME)

# Help
help:
	@echo "Makefile commands:"
	@echo "  make build   - Build the project"
	@echo "  make test    - Run tests"
	@echo "  make clean   - Clean up build artifacts"
	@echo "  make run     - Run the application"
	@echo "  make deps    - Install dependencies"
	@echo "  make fmt     - Format the code"
	@echo "  make lint    - Lint the code"