# Variables
BINARY_NAME=gobuild
BUILD_DIR=build
GO_FILES=$(shell find . -name '*.go')
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# Default target
.PHONY: all build test clean run deps fmt lint deploy help

all: deps fmt test build

# Build the project
build:
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gobuild

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean up build artifacts
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)

# Run the application
run: build
	@echo "Running ${BINARY_NAME}..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Lint the code
lint:
	@echo "Linting code..."
	golangci-lint run

# Deploy the application
deploy:
	@echo "Deploying ${BINARY_NAME}..."
	@if [ -z "$(DEPLOY_TOKEN)" ]; then \
		echo "Error: DEPLOY_TOKEN is not set"; \
		exit 1; \
	fi
	@echo "Deploying version ${VERSION}..."
	./scripts/deploy.sh

# Help
help:
	@echo "Available commands:"
	@echo "  make          - Build the project (default)"
	@echo "  make all      - Run deps, fmt, test, and build"
	@echo "  make build    - Build the project"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build directory"
	@echo "  make run      - Build and run the application"
	@echo "  make deps     - Install dependencies"
	@echo "  make fmt      - Format the code"
	@echo "  make lint     - Run linter"
	@echo "  make deploy   - Deploy the application"
	@echo "  make help     - Show this help message"