#!/bin/bash

# Run all tests and verify functionality
echo "Running comprehensive tests..."

# Unit Tests
echo "🧪 Running unit tests..."
go test ./internal/pkg/... -v

# Integration Tests
echo "🔄 Running integration tests..."
go test ./internal/cmd/... -v -tags=integration

# E2E Tests
echo "🎯 Running E2E tests..."
go test ./test/e2e/... -v

# Framework Adapter Tests
echo "🏗️  Testing framework adapters..."
go test ./internal/pkg/framework/... -v

# Template Tests
echo "📝 Testing template system..."
go test ./internal/pkg/template/... -v

# Build System Tests
echo "🔨 Testing build system..."
go test ./internal/pkg/builder/... -v