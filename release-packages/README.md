# Package Manager Distribution

This directory contains configuration files for distributing TesselBox through various package managers.

## Homebrew (macOS)

```bash
brew install tesselbox/tesselbox/tesselbox
```

The formula is located in `homebrew/tesselbox.rb`.

## Snap (Linux)

```bash
snap install tesselbox
```

Configuration is in `snapcraft.yaml`.

## Winget (Windows)

```powershell
winget install TesselBox.TesselBox
```

Manifest is in `winget.yaml`.

## Release Process

1. Update version numbers in all package files
2. Generate SHA256 hashes for release binaries
3. Submit packages to respective stores:
   - Homebrew: Create a tap or submit to homebrew-core
   - Snap: Publish to Snap Store
   - Winget: Submit to winget-pkgs repository

## Automated Updates

The GitHub Actions workflow will automatically update SHA256 hashes and version numbers when creating releases.
