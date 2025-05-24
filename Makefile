# Plexr Makefile

# Variables
BINARY_NAME=plexr
MAIN_PATH=./cmd/plexr
BUILD_DIR=./build
DIST_DIR=./dist

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build variables
VERSION?=dev
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# OS/Arch detection
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

.PHONY: all build clean test coverage deps run help

## help: Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## all: Build for current platform
all: clean deps test build

## build: Build the binary for current platform
build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

## build-all: Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	rm -f coverage.txt coverage.html

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

## coverage: Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## run: Run the application
run: build
	$(BUILD_DIR)/$(BINARY_NAME)

## install: Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

## dev-setup: Setup development environment
dev-setup:
	@echo "Setting up development environment..."
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOGET) -u github.com/goreleaser/goreleaser
	@echo "Development environment ready!"

## lint: Run linters
lint:
	@echo "Running linters..."
	golangci-lint run

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...
