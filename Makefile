.PHONY: build test clean install release docker lint security coverage help

BINARY_NAME=nexa
VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")

build:
	@echo "Building $(BINARY_NAME)..."
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/$(BINARY_NAME) ./cmd/nexa

test:
	@echo "Running tests..."
	go test -v -race ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/ dist/

install: build
	@echo "Installing..."
	sudo cp bin/$(BINARY_NAME) /usr/local/bin/

release: test
	@echo "Building release binaries..."
	./scripts/build_linux.sh $(VERSION)
	./scripts/build_windows.ps1 $(VERSION)

docker:
	@echo "Building Docker image..."
	docker build -t nexa:$(VERSION) -f scripts/docker/Dockerfile .

lint:
	@echo "Running linter..."
	golangci-lint run

security:
	@echo "Running security scan..."
	gosec ./...

coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build     - Build the binary"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build artifacts"
	@echo "  install   - Install to system"
	@echo "  release   - Build release binaries"
	@echo "  docker    - Build Docker image"
	@echo "  lint      - Run linter"
	@echo "  security  - Run security scan"
	@echo "  coverage  - Generate coverage report"