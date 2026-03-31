# TesselBox Build System
# Supports cross-platform builds (no icons)

.PHONY: all clean build windows linux darwin release test test-verbose test-coverage test-coverage-html test-integration test-migration test-unit test-race test-bench clean-test

# Default target
all: build

# Build for current platform
build:
	@echo "Building for $(shell go env GOOS)/$(shell go env GOARCH)..."
	@mkdir -p bin
	@go build -ldflags "-X main.Version=v0.3-alpha" -o bin/tesselbox cmd/main.go

# Build for all platforms
release: clean
	@echo "Building release binaries..."
	@mkdir -p bin
	@echo "Building Windows AMD64..."
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.Version=v0.3-alpha -s -w" -trimpath -o bin/tesselbox-windows-amd64.exe cmd/main.go
	@if [ "$(shell go env GOOS)" = "linux" ]; then \
		echo "Building Linux AMD64..."; \
		go build -ldflags "-X main.Version=v0.3-alpha -s -w" -trimpath -o bin/tesselbox-linux-amd64 cmd/main.go; \
	else \
		echo "Warning: Cross-compiling to Linux may fail due to Ebiten CGO dependencies"; \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.Version=v0.3-alpha -s -w" -trimpath -o bin/tesselbox-linux-amd64 cmd/main.go || echo "Linux cross-compilation failed - compile on native hardware"; \
	fi
	@echo "Note: ARM64 builds require native compilation due to Ebiten CGO dependencies"
	@echo "Release binaries built in bin/"
	@echo "Note: macOS builds require native compilation on macOS"

# Platform-specific builds
windows:
	@mkdir -p bin
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.Version=v0.3-alpha" -o bin/tesselbox.exe ./cmd/

linux:
	@mkdir -p bin
	@if [ "$(shell go env GOOS)" = "linux" ]; then \
		go build -ldflags "-X main.Version=v0.3-alpha" -o bin/tesselbox ./cmd/; \
	else \
		echo "Warning: Cross-compiling to Linux may fail due to Ebiten CGO dependencies"; \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.Version=v0.3-alpha" -o bin/tesselbox ./cmd/ || echo "Linux cross-compilation failed - compile on native hardware"; \
	fi

linux-arm64:
	@mkdir -p bin
	@echo "Warning: ARM64 builds may fail due to Ebiten CGO dependencies"
	@echo "For ARM64 builds, compile on native ARM64 hardware"
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-X main.Version=v0.3-alpha" -o bin/tesselbox-arm64 ./cmd/ || echo "ARM64 build failed - compile on native hardware"

darwin:
	@mkdir -p bin
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.Version=v0.3-alpha" -o bin/tesselbox ./cmd/

darwin-arm64:
	@mkdir -p bin
	@echo "Warning: ARM64 builds may fail due to Ebiten CGO dependencies"
	@echo "For ARM64 builds, compile on native ARM64 hardware"
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-X main.Version=v0.3-alpha" -o bin/tesselbox-arm64 ./cmd/ || echo "ARM64 build failed - compile on native hardware"

# Development build (fast)
dev:
	@echo "Building development version..."
	@go build -o tesselbox ./cmd/
	@echo "Development binary: tesselbox"

# Run tests
test:
	@echo "Running all tests..."
	go test ./...

test-verbose:
	@echo "Running all tests with verbose output..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

test-coverage-html:
	@echo "Running tests with HTML coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

test-integration:
	@echo "Running integration tests..."
	go test ./tests/integration/...

test-migration:
	@echo "Running migration tests..."
	go test ./tests/migration/

test-unit:
	@echo "Running unit tests..."
	go test ./tests/unit/...

test-race:
	@echo "Running tests with race detection..."
	go test -race ./...

test-bench:
	@echo "Running benchmarks..."
	go test -bench ./...

clean-test:
	@echo "Cleaning test files..."
	rm -f coverage.out coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f tesselbox tesselbox.exe
	@rm -f rsrc.syso
	@rm -f assets/windows.manifest
	@rm -f tesselbox.desktop
	@rm -f TesselBox.app

# Install build dependencies
deps:
	@echo "Installing build dependencies..."
	@if command -v convert >/dev/null 2>&1; then \
		echo "ImageMagick already installed"; \
	else \
		echo "Install ImageMagick for icon generation:"; \
		echo "  macOS: brew install imagemagick"; \
		echo "  Ubuntu: sudo apt-get install imagemagick"; \
		echo "  Fedora: sudo dnf install ImageMagick"; \
	fi

# Development server (with auto-rebuild)
dev-server:
	@echo "Starting development server..."
	@go install github.com/air-verse/air@latest
	@air -c .air.toml

# Create distribution packages
dist: release
	@echo "Creating distribution packages..."
	@mkdir -p dist
	@cd bin && tar -czf ../dist/tesselbox-linux-amd64.tar.gz tesselbox-linux-amd64
	@cd bin && zip -r ../dist/tesselbox-windows-amd64.zip tesselbox-windows-amd64.exe
	@echo "Distribution packages created in dist/"

# Show help
help:
	@echo "TesselBox Build System"
	@echo ""
	@echo "Targets:"
	@echo "  build      - Build for current platform"
	@echo "  release    - Build for all platforms (release)"
	@echo "  windows    - Build Windows binary"
	@echo "  linux      - Build Linux binary (amd64)"
	@echo "  linux-arm64 - Build Linux binary (arm64)"
	@echo "  darwin     - Build macOS binary (amd64)"
	@echo "  darwin-arm64 - Build macOS binary (arm64)"
	@echo "  dev        - Quick development build"
	@echo "  test-verbose       - Run tests with verbose output"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  test-coverage-html - Run tests with HTML coverage"
	@echo "  test-integration   - Run integration tests only"
	@echo "  test-migration    - Run migration tests only"
	@echo "  test-unit         - Run unit tests only"
	@echo "  test-race         - Run tests with race detection"
	@echo "  test-bench         - Run benchmarks"
	@echo "  clean-test        - Clean test files"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Install build dependencies"
	@echo "  dist       - Create distribution packages"
	@echo "  help       - Show this help"
