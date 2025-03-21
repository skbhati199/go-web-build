#!/bin/bash

# Test script for template generation and build tools

echo "Testing React JavaScript template..."
mkdir -p /tmp/test-react-js
cd /tmp/test-react-js
go run /Users/sonukumar/go-web-build/cmd/gobuild/main.go create test-app --framework react --template javascript

if [ $? -eq 0 ]; then
  echo "✅ React JavaScript template created successfully"
  cd test-app
  go run /Users/sonukumar/go-web-build/cmd/gobuild/main.go build --mode development
  if [ $? -eq 0 ]; then
    echo "✅ Development build successful"
  else
    echo "❌ Development build failed"
  fi
else
  echo "❌ React JavaScript template creation failed"
fi

echo "Testing React TypeScript template..."
mkdir -p /tmp/test-react-ts
cd /tmp/test-react-ts
go run /Users/sonukumar/go-web-build/cmd/gobuild/main.go create test-app --framework react --template typescript

if [ $? -eq 0 ]; then
  echo "✅ React TypeScript template created successfully"
  cd test-app
  go run /Users/sonukumar/go-web-build/cmd/gobuild/main.go build --mode development
  if [ $? -eq 0 ]; then
    echo "✅ Development build successful"
  else
    echo "❌ Development build failed"
  fi
else
  echo "❌ React TypeScript template creation failed"
fi

echo "Test completed"