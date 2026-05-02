# Makefile for devkit CLI DevOps Toolkit

# Variables
VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)"

# Binary name
BINARY_NAME := devkit

# Build directory
DIST_DIR := dist

# Default target
.DEFAULT_GOAL := build

# Build for current platform
.PHONY: build
build:
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(DIST_DIR)
	go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(DIST_DIR)/$(BINARY_NAME)"

# Run all tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Run tests with verbose output
.PHONY: test-verbose
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(DIST_DIR)/
	@echo "Clean complete"

# Cross-compile for all platforms
.PHONY: release
release:
	@echo "Building release binaries for all platforms..."
	@mkdir -p $(DIST_DIR)
	@echo "Building for linux/amd64..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	@echo "Building for darwin/amd64..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@echo "Building for darwin/arm64..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@echo "Building for windows/amd64..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Release build complete. Binaries in $(DIST_DIR)/"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build binary for current platform (default)"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  clean         - Remove dist/ directory"
	@echo "  release       - Cross-compile for all platforms"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION       - Version string (default: dev)"
	@echo "  COMMIT        - Git commit hash (auto-detected)"
	@echo "  BUILD_DATE    - Build timestamp (auto-generated)"
