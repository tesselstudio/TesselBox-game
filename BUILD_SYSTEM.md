# TesselBox Build System

## Overview
Simple Make-based build system for TesselBox with cross-platform support.

## Usage

### Local Development
```bash
# Build for current platform
make build

# Build for all platforms (release optimized)
make release

# Build for specific platform
make windows
make linux
make darwin
make linux-arm64
make darwin-arm64

# Quick development build
make dev

# Clean build artifacts
make clean

# Install build dependencies
make deps

# Create distribution packages
make dist
```

## Features

✅ **Cross-platform builds** (Windows, Linux, macOS Intel/ARM)
✅ **Optimized builds** (stripped binaries for release)
✅ **Simple Make-based system** (no external build scripts)
✅ **Development builds** (fast iteration)
✅ **Package creation** (tar.gz and zip distributions)

## Output

Builds create binaries in the `bin/` directory:
- `tesselbox-windows-amd64.exe`
- `tesselbox-linux-amd64`
- `tesselbox-darwin-amd64`
- `tesselbox-darwin-arm64`

Release builds include optimization flags:
- `-s -w` for stripped binaries
- `-trimpath` for cleaner stack traces
- Version injection via ldflags

## Build Targets

### Core Targets
- `build` - Build for current platform
- `release` - Build optimized binaries for all platforms
- `clean` - Remove build artifacts
- `dev` - Quick development build

### Platform-specific
- `windows` - Windows AMD64 binary
- `linux` - Linux AMD64 binary
- `darwin` - macOS Intel binary
- `linux-arm64` - Linux ARM64 binary (native compilation required)
- `darwin-arm64` - macOS ARM64 binary (native compilation required)

### Utilities
- `test` - Run all tests
- `test-verbose` - Run tests with verbose output
- `test-coverage` - Run tests with coverage
- `deps` - Install build dependencies
- `dist` - Create distribution packages
- `help` - Show available targets

## Notes

- ARM64 builds require native compilation due to Ebiten CGO dependencies
- macOS builds should be performed on macOS hardware
- Icons are generated separately if needed
- Version is set to `v0.3-alpha` by default

The build system is now simplified and uses native Make commands! 🚀
