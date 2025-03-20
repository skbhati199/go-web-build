#!/bin/bash

# Run unit tests
echo "Running unit tests..."
go test ./internal/... -v

# Run integration tests
echo "Running integration tests..."
go test ./test/integration/... -v

# Run E2E tests
echo "Running E2E tests..."
go test ./test/e2e/... -v

# Run benchmarks
echo "Running benchmarks..."
go test ./test/benchmark/... -bench=. -benchmem