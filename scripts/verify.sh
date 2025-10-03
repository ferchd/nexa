#!/bin/bash
set -e

echo "Verifying Nexa build..."

# Check Go version
echo "Go version: $(go version)"

# Download dependencies
echo "Downloading dependencies..."
go mod download
go mod verify

# Run tests
echo "Running tests..."
go test -v ./...

# Build binary
echo "Building binary..."
go build -o nexa ./cmd/nexa

# Test basic functionality
echo "Testing basic functionality..."
./nexa --version
./nexa --help || true

echo "All checks passed!"