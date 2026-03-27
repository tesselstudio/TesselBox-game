# TesselBox Build System

## Overview
Complete cross-platform build system for TesselBox with embedded assets and automated releases.

## Files Created

### Core Build System
- `build/build.go` - Main build script with cross-platform support
- `build/go.mod` - Build module dependencies
- `build/version.go` - Version info generator

### Icon Generation
- `build/generate-icons.sh` - Creates platform-specific icons
- `assets/icons/` - Generated icon files (.ico, .icns, .png)

### Automation
- `.github/workflows/release.yml` - GitHub Actions CI/CD
- `build/release.sh` - Automated release script

## Usage

### Local Development
```bash
# Build for current platform
make build

# Build for all platforms
make release

# Build for specific platform
make windows
make linux
make darwin
```

### Direct Build Script
```bash
cd build
go run build.go -os=linux -arch=amd64 -output=tesselbox -release
```

### Generate Icons
```bash
cd build
./generate-icons.sh
```

### Release Process
```bash
cd build
./release.sh 2.0.0
```

## Features

✅ **Cross-platform builds** (Windows, Linux, macOS Intel/ARM)
✅ **Embedded assets** (single binary distribution)
✅ **Icon generation** (platform-specific formats)
✅ **Automated releases** (GitHub Actions)
✅ **Version management** (semantic versioning)
✅ **Optimized builds** (stripped binaries)
✅ **CI/CD pipeline** (automated testing and building)

## Output

Builds create optimized binaries in the `bin/` directory:
- `tesselbox-windows-amd64.exe`
- `tesselbox-linux-amd64`
- `tesselbox-darwin-amd64`
- `tesselbox-darwin-arm64`

## Release Automation

1. Push tag to trigger GitHub Actions
2. Automatic cross-platform builds
3. Generated release with all binaries
4. Automatic release notes
5. Artifact upload

## Next Steps

To complete the distribution system:

1. **Package Managers**: Homebrew, Snap, Winget formulas
2. **Game Platforms**: Itch.io, Steam integration
3. **Installer Creation**: MSI, DMG, AppImage packages
4. **Auto-updater**: Built-in update mechanism

The build system is now complete and ready for production use! 🚀
