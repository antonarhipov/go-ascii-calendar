# ASCII Calendar Makefile
# Build automation for the ASCII Calendar terminal application

# Configuration
APP_NAME = ascii-calendar
GO_VERSION = 1.19
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date -u +"%Y-%m-%d %H:%M:%S UTC")
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS = -ldflags "-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.CommitHash=$(COMMIT_HASH)'"

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(APP_NAME) v$(VERSION)..."
	go build $(LDFLAGS) -o $(APP_NAME) .
	@echo "Build complete: $(APP_NAME)"

# Build for all supported platforms
.PHONY: build-all
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p dist
	
	# Linux (64-bit)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(APP_NAME)-linux-amd64 .
	
	# Linux (ARM64)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(APP_NAME)-linux-arm64 .
	
	# macOS (Intel)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(APP_NAME)-darwin-amd64 .
	
	# macOS (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(APP_NAME)-darwin-arm64 .
	
	# Windows (64-bit)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(APP_NAME)-windows-amd64.exe .
	
	# Create checksums
	cd dist && sha256sum * > checksums.txt
	@echo "Cross-platform builds complete in dist/ directory"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run integration tests
.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	go test -v -tags=integration ./...

# Lint the code
.PHONY: lint
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, running basic checks..."; \
		go vet ./...; \
		go fmt ./...; \
	fi

# Format code
.PHONY: format
format:
	@echo "Formatting code..."
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(APP_NAME)
	rm -rf dist/
	rm -f coverage.out coverage.html

# Run the application
.PHONY: run
run: build
	@echo "Starting $(APP_NAME)..."
	./$(APP_NAME)

# Development mode (rebuild and run on changes)
.PHONY: dev
dev:
	@echo "Development mode - rebuild on changes..."
	@if command -v watchexec >/dev/null 2>&1; then \
		watchexec -e go "make run"; \
	else \
		echo "watchexec not installed. Install with: cargo install watchexec-cli"; \
		echo "Falling back to single run..."; \
		make run; \
	fi

# Create a release package
.PHONY: release
release: clean test lint build-all
	@echo "Creating release package v$(VERSION)..."
	@mkdir -p release
	
	# Copy documentation
	cp README.md LICENSE* release/ 2>/dev/null || true
	cp -r docs/ release/ 2>/dev/null || true
	
	# Create platform-specific packages
	@for binary in dist/*; do \
		if [ -f "$$binary" ] && [ "$$(basename $$binary)" != "checksums.txt" ]; then \
			platform=$$(basename $$binary | sed 's/$(APP_NAME)-//'); \
			echo "Creating package for $$platform..."; \
			mkdir -p "release/$(APP_NAME)-$(VERSION)-$$platform"; \
			cp "$$binary" "release/$(APP_NAME)-$(VERSION)-$$platform/"; \
			cp README.md "release/$(APP_NAME)-$(VERSION)-$$platform/" 2>/dev/null || true; \
			cp LICENSE* "release/$(APP_NAME)-$(VERSION)-$$platform/" 2>/dev/null || true; \
			cd release && tar -czf "$(APP_NAME)-$(VERSION)-$$platform.tar.gz" "$(APP_NAME)-$(VERSION)-$$platform/" && cd ..; \
		fi; \
	done
	
	# Copy checksums
	cp dist/checksums.txt release/
	
	@echo "Release packages created in release/ directory"
	@ls -la release/

# Install the application system-wide
.PHONY: install
install: build
	@echo "Installing $(APP_NAME) to /usr/local/bin..."
	@sudo cp $(APP_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(APP_NAME)
	@echo "$(APP_NAME) installed successfully"

# Uninstall the application
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(APP_NAME)..."
	@sudo rm -f /usr/local/bin/$(APP_NAME)
	@echo "$(APP_NAME) uninstalled"

# Check system requirements
.PHONY: check-env
check-env:
	@echo "Checking system requirements..."
	@echo "Go version: $$(go version)"
	@echo "Terminal size: $$(tput cols)x$$(tput lines) (minimum required: 80x24)"
	@echo "Current directory: $$(pwd)"
	@echo "Available space: $$(df -h . | tail -1 | awk '{print $$4}')"
	@if [ $$(tput cols) -lt 80 ] || [ $$(tput lines) -lt 24 ]; then \
		echo "WARNING: Terminal size is smaller than recommended minimum (80x24)"; \
	fi

# Help target
.PHONY: help
help:
	@echo "ASCII Calendar Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  build          - Build the application for current platform"
	@echo "  build-all      - Build for all supported platforms"
	@echo "  deps           - Install Go dependencies"
	@echo "  test           - Run unit tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-integration - Run integration tests"
	@echo "  lint           - Run code linters"
	@echo "  format         - Format code"
	@echo "  clean          - Clean build artifacts"
	@echo "  run            - Build and run the application"
	@echo "  dev            - Development mode (rebuild on changes)"
	@echo "  release        - Create release packages"
	@echo "  install        - Install system-wide"
	@echo "  uninstall      - Uninstall system-wide"
	@echo "  check-env      - Check system requirements"
	@echo "  help           - Show this help message"
	@echo ""
	@echo "Configuration:"
	@echo "  VERSION        - Version string (default: git tag or 'dev')"
	@echo "  APP_NAME       - Application name (default: ascii-calendar)"
	@echo ""
	@echo "Examples:"
	@echo "  make build                    # Build for current platform"
	@echo "  make VERSION=1.0.0 release    # Create v1.0.0 release"
	@echo "  make test-coverage            # Run tests with coverage"