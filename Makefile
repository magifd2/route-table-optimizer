# Makefile for route-table-optimizer

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary name
BINARY_NAME=route-table-optimizer
BINARY_UNIX=$(BINARY_NAME)

# Build directory
BUILD_DIR=build

# Default target
all: test build

# Build the binary for the current platform
build:
	$(GOBUILD) -o $(BINARY_NAME) .

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -rf $(BUILD_DIR)

# Cross-compilation targets for a release
release: release-linux release-windows release-mac

release-linux:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(BUILD_DIR)/linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/linux/$(BINARY_NAME) .

release-windows:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(BUILD_DIR)/windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/windows/$(BINARY_NAME).exe .

release-mac:
	@echo "Building for macOS (amd64)..."
	@mkdir -p $(BUILD_DIR)/mac
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/mac/$(BINARY_UNIX) .

.PHONY: all build test clean release release-linux release-windows release-mac
