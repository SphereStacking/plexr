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
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# OS/Arch detection
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

# Tool versions
GOLANGCI_LINT_VERSION=v1.61.0

# Tool paths
TOOLS_DIR=.tools
GOLANGCI_LINT=$(TOOLS_DIR)/golangci-lint
GOIMPORTS=$(shell which goimports 2>/dev/null || echo "$(GOPATH)/bin/goimports")

# Git hooks
HOOKS_DIR=.githooks

.PHONY: all build clean test coverage deps run help lint fmt fmt-check vet tools dev-setup install hooks

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
dev-setup: tools hooks
	@echo "Setting up Go workspace..."
	@go mod download
	@echo "Development environment ready!"
	@echo ""
	@echo "VSCode users: Install recommended extensions when prompted"
	@echo "Run 'code .' to open in VSCode"

## tools: Install development tools
tools: $(GOLANGCI_LINT) goimports

$(GOLANGCI_LINT):
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@mkdir -p $(TOOLS_DIR)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_DIR) $(GOLANGCI_LINT_VERSION)
	@echo "golangci-lint installed!"

## goimports: Install goimports
goimports:
	@if ! command -v goimports &> /dev/null; then \
		echo "Installing goimports..."; \
		$(GOCMD) install golang.org/x/tools/cmd/goimports@latest; \
	fi

## hooks: Install git hooks
hooks:
	@echo "Installing git hooks..."
	@git config core.hooksPath $(HOOKS_DIR)
	@echo "Git hooks installed!"

## lint: Run linters
lint: $(GOLANGCI_LINT)
	@echo "Running linters..."
	$(GOLANGCI_LINT) run

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...
	@echo "Running goimports..."
	@which goimports > /dev/null 2>&1 || $(GOCMD) install golang.org/x/tools/cmd/goimports@latest
	@export PATH="$$PATH:$$(go env GOPATH)/bin" && goimports -w .

## fmt-check: Check if code is formatted
fmt-check:
	@echo "Checking code formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "The following files are not formatted:"; \
		gofmt -l .; \
		exit 1; \
	fi
	@echo "All files are properly formatted!"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...
