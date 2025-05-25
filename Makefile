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

.PHONY: all build clean test coverage deps run help lint fmt fmt-check vet tools dev-setup install hooks release release-check release-local

## help: Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Main targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -E "^  (build|test|install|release)" | sort
	@echo ""
	@echo "Development targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -vE "^  (build|test|install|release)" | sort

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

## release-check: Check if ready for release
release-check:
	@echo "🔍 Checking release readiness..."
	@echo ""
	@echo "1️⃣  Checking for uncommitted changes..."
	@if ! git diff --quiet || ! git diff --cached --quiet; then \
		echo "❌ You have uncommitted changes. Please commit or stash them."; \
		exit 1; \
	else \
		echo "✅ Working directory is clean"; \
	fi
	@echo ""
	@echo "2️⃣  Checking if on main branch..."
	@if [ "$$(git rev-parse --abbrev-ref HEAD)" != "main" ]; then \
		echo "❌ You are not on the main branch. Please switch to main."; \
		exit 1; \
	else \
		echo "✅ On main branch"; \
	fi
	@echo ""
	@echo "3️⃣  Checking if main is up to date..."
	@git fetch origin
	@if [ "$$(git rev-parse HEAD)" != "$$(git rev-parse origin/main)" ]; then \
		echo "❌ Your main branch is not up to date. Please pull latest changes."; \
		exit 1; \
	else \
		echo "✅ Main branch is up to date"; \
	fi
	@echo ""
	@echo "4️⃣  Running tests..."
	@if ! $(MAKE) test > /dev/null 2>&1; then \
		echo "❌ Tests failed. Please fix them before releasing."; \
		exit 1; \
	else \
		echo "✅ All tests pass"; \
	fi
	@echo ""
	@echo "5️⃣  Running linter..."
	@if ! $(MAKE) lint > /dev/null 2>&1; then \
		echo "❌ Linting failed. Please fix issues before releasing."; \
		exit 1; \
	else \
		echo "✅ Linting passes"; \
	fi
	@echo ""
	@echo "🎉 Ready for release!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Decide on version number (current: $$(git describe --tags --abbrev=0 2>/dev/null || echo 'no tags yet'))"
	@echo "  2. Run: make release VERSION=v0.1.0"

## release: Create a new release
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION is required. Usage: make release VERSION=v0.1.0"; \
		exit 1; \
	fi
	@if ! echo "$(VERSION)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9\.]+)?$$'; then \
		echo "❌ Invalid version format. Use semantic versioning like v0.1.0 or v1.0.0-rc.1"; \
		exit 1; \
	fi
	@echo "🚀 Creating release $(VERSION)..."
	@echo ""
	@echo "📝 Recent commits since last tag:"
	@echo "================================"
	@git log --oneline $$(git describe --tags --abbrev=0 2>/dev/null || echo "")..HEAD 2>/dev/null || git log --oneline -10
	@echo ""
	@echo "Do you want to create release $(VERSION)? [y/N] " && read ans && [ $${ans:-N} = y ]
	@echo ""
	@echo "📝 Enter release notes (or press Ctrl+D when done):"
	@echo "---------------------------------------------------"
	@NOTES=$$(cat) && \
	git tag -a $(VERSION) -m "Release $(VERSION)" -m "$$NOTES"
	@echo ""
	@echo "✅ Tag $(VERSION) created!"
	@echo ""
	@echo "To push the release, run:"
	@echo "  git push origin $(VERSION)"
	@echo ""
	@echo "Or to push everything including the tag:"
	@echo "  git push origin main --tags"

## release-local: Test release process locally
release-local:
	@echo "🧪 Testing release locally with goreleaser..."
	@if ! command -v goreleaser > /dev/null; then \
		echo "❌ goreleaser not installed. Install from https://goreleaser.com"; \
		exit 1; \
	fi
	goreleaser release --snapshot --clean
	@echo ""
	@echo "✅ Local release artifacts created in ./dist/"
	@echo ""
	@ls -la dist/ | grep -E "(tar\.gz|zip|checksums)" || true
