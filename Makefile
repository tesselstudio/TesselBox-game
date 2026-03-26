# TesselBox Build System
# Supports cross-platform builds with icon embedding

.PHONY: all clean build windows linux darwin release test icons

# Default target
all: build

# Build for current platform
build:
	@echo "Building for $(shell go env GOOS)/$(shell go env GOARCH)..."
	@go run build/build.go

# Build for all platforms
release: clean icons
	@echo "Building release binaries..."
	@go run build/build.go -os=windows -arch=amd64 -output=bin/tesselbox-windows-amd64.exe -release
	@go run build/build.go -os=linux -arch=amd64 -output=bin/tesselbox-linux-amd64 -release
	@go run build/build.go -os=darwin -arch=amd64 -output=bin/tesselbox-darwin-amd64 -release
	@go run build/build.go -os=darwin -arch=arm64 -output=bin/tesselbox-darwin-arm64 -release
	@echo "Release binaries built in bin/"

# Platform-specific builds
windows:
	@mkdir -p bin
	@go run build/build.go -os=windows -arch=amd64 -output=bin/tesselbox.exe

linux:
	@mkdir -p bin
	@go run build/build.go -os=linux -arch=amd64 -output=bin/tesselbox

darwin:
	@mkdir -p bin
	@go run build/build.go -os=darwin -arch=amd64 -output=bin/tesselbox

# Generate placeholder icons
icons:
	@echo "Generating placeholder icons..."
	@if command -v convert >/dev/null 2>&1; then \
		./assets/icons/create_icons.sh; \
	else \
		echo "ImageMagick not found. Install with: brew install imagemagick or apt-get install imagemagick"; \
	fi

# Development build (fast)
dev:
	@echo "Building development version..."
	@go build -o tesselbox cmd/main.go
	@echo "Development binary: tesselbox"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

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
	@go install github.com/akavel/rsrc@latest
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
	@cd bin && tar -czf ../dist/tesselbox-darwin-amd64.tar.gz TesselBox.app
	@cd bin && tar -czf ../dist/tesselbox-darwin-arm64.tar.gz TesselBox.app
	@echo "Distribution packages created in dist/"

# Show help
help:
	@echo "TesselBox Build System"
	@echo ""
	@echo "Targets:"
	@echo "  build      - Build for current platform"
	@echo "  release    - Build for all platforms (release)"
	@echo "  windows    - Build Windows binary"
	@echo "  linux      - Build Linux binary"
	@echo "  darwin     - Build macOS binary"
	@echo "  icons      - Generate placeholder icons"
	@echo "  dev        - Quick development build"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Install build dependencies"
	@echo "  dist       - Create distribution packages"
	@echo "  help       - Show this help"
