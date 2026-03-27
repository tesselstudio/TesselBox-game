# TesselBox: Hexagonal Sandbox Adventure 

[![Typing SVG](https://readme-typing-svg.demolab.com/?lines=Experience+TesselBox,+a+captivating+open-source+game+engine+built+in+Go.;Dive+into+a+procedural+hexagonal+world,+where+you+can+mine,+build,+and+explore+endless+landscapes.)](https://git.io/typing-svg)
[![GitHub stars](https://img.shields.io/github/stars/tesselstudio/TesselBox-game?style=social)](https://github.com/tesselstudio/TesselBox-game/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/tesselstudio/TesselBox-game?style=social)](https://github.com/tesselstudio/TesselBox-game/network/members)
[![GitHub last commit](https://img.shields.io/github/last-commit/tesselstudio/TesselBox-game)](https://github.com/tesselstudio/TesselBox-game/commits/main)
[![GitHub repo size](https://img.shields.io/github/repo-size/tesselstudio/TesselBox-game)](https://github.com/tesselstudio/TesselBox-game/size)
[![GitHub contributors](https://img.shields.io/github/contributors/tesselstudio/TesselBox-game)](https://github.com/tesselstudio/TesselBox-game/graphs/contributors)


<p align="center">
  <a href="https://github.com/tesselstudio/TesselBox-game/archive/refs/heads/main.zip">
    <img src="download.png" alt="Download TesselBox" width="120" style="border-radius: 6px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
  </a>
</p>

Experience TesselBox, a captivating open-source game engine built in Go. Dive into a procedural hexagonal world, where you can mine, build, and explore endless landscapes☝🏼☝🏼☝🏼☝🏼☝🏼




## 🎮 Gameplay Overview

TesselBox offers a captivating experience in an open-source, hexagonal sandbox environment. Players can engage in mining, building, and exploring vast, procedurally generated landscapes. The game's core features include building unique structures on a hexagonal grid, discovering new terrains through procedural generation, and unearthing resources while exploring the world.

### Features

-   **Hexagonal World-Building**: Craft unique structures in a distinct hexagonal grid, offering natural movement patterns and strategic depth.
-   **Procedural Generation**: Discover new and exciting terrains with every play, ensuring limitless and diverse exploration.
-   **Mining & Exploration**: Unearth resources and uncover the secrets of the world as you explore its vast landscapes.
-   **Open Source**: Join our community and contribute to the evolution of TesselBox!

## 🚀 Quick Start

Ready to play? Get TesselBox running in under a minute!

### Option 1: Play Immediately (Pre-built Binaries)

1. **Download the latest release** from [GitHub Releases](https://github.com/tesselstudio/TesselBox-game/releases)
2. **Choose your platform**:
   - Windows: `tesselbox-windows-amd64.exe`
   - Linux: `tesselbox-linux-amd64`
   - macOS: Build from source (requires native compilation)
3. **Run the binary**:
   ```bash
   # Linux/macOS
   chmod +x tesselbox-*
   ./tesselbox-*
   
   # Windows
   tesselbox-windows-amd64.exe
   ```

### Option 2: Build from Source

1. **Clone the repository**:
    ```bash
    git clone https://github.com/tesselstudio/TesselBox-game.git
    cd TesselBox-game
    ```
2. **Run the game**:
    ```bash
    go run cmd/main.go
    ```

### Option 3: Build Release Binaries

```bash
# Build for current platform
make build

# Build release binaries (Windows + Linux)
make release

# Build for specific platform
make windows
make linux
```

## 🎮 Controls

- **Right-click**: Place blocks
- **Left-click**: Mine/break blocks
- **Number keys 1-9**: Select block types (when implemented)
- **WASD**: Move player
- **Mouse**: Look around

## 🛠️ Development

### Building from Source

**Prerequisites**:
- Go 1.21 or later
- For Linux: `gcc` (for CGO dependencies)
- For macOS: Xcode Command Line Tools
- For Windows: MinGW-w64 (optional, for cross-compilation)

**Development Build**:
```bash
# Quick development build
go run cmd/main.go

# Development binary
make dev

# Full build with optimizations
make build
```

### Project Structure

```
TesselBox-game/
├── cmd/main.go           # Game entry point
├── pkg/                   # Core packages
│   ├── blocks/           # Block system
│   ├── world/            # World generation
│   ├── player/           # Player mechanics
│   └── hexagon/          # Hexagonal math
├── assets/               # Embedded assets
│   └── config/           # Game configuration
├── build/                # Build system
│   ├── build.go         # Cross-platform builder
│   └── generate-icons.sh # Icon generator
└── .github/workflows/    # CI/CD pipelines
```

### Build System Features

✅ **Cross-platform builds** (Windows, Linux, macOS)
✅ **Embedded assets** (single binary distribution)
✅ **Automated releases** (GitHub Actions)
✅ **Icon generation** (platform-specific formats)
✅ **Version management** (semantic versioning)

### Advanced Build Options

```bash
# Custom build with specific options
cd build
go run build.go \
  -os=linux \
  -arch=amd64 \
  -output=tesselbox-custom \
  -release \
  -version=2.0.0

# Generate icons
cd build
./generate-icons.sh

# Create release
cd build
./release.sh 2.0.0
```

### Contributing

We welcome all contributions! Whether you're fixing bugs, adding features, or improving documentation.

**Development Workflow**:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

**Areas needing contribution**:
- 🎮 Block variety system
- ⛏️ Enhanced mining mechanics
- 📦 Inventory management
- 🎨 UI/UX improvements
- 🐛 Bug fixes and optimizations

## 🛠️ Built With

TesselBox is built with modern technologies to deliver a robust and engaging gaming experience.

<p align="center">
    <img src="https://skillicons.dev/icons?i=go,linux,windows,git,github,markdown,vscode,apple" alt="Technology Stack">
</p>

-   **Go**: The primary programming language for the game engine.
-   **Linux, Windows, Apple**: Supported operating systems for running the game.
-   **Git & GitHub**: Version control and collaborative development platform.
-   **Markdown**: Used for documentation and README formatting.
-   **VS Code**: Recommended development environment.

## Star History 

<a href="https://www.star-history.com/">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/image?repos=tesselstudio/TesselBox-game&type=timeline&theme=dark&legend=bottom-right" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/image?repos=tesselstudio/TesselBox-game&type=timeline&legend=bottom-right" />
   <img alt="Star History Chart" src="https://api.star-history.com/image?repos=tesselstudio/TesselBox-game&type=timeline&legend=bottom-right" />
 </picture>
</a>

## 🤝 Contributing

We welcome all contributions! Whether you're a developer, designer, or tester, your input helps us grow. Please see our [contributing.md](contributing.md) for detailed guidelines.

## 👥 Contributors

A big thank you to everyone who has contributed to TesselBox!

[![https://contrib.rocks/image?repo=tesselstudio/TesselBox-game](https://contrib.rocks/image?repo=tesselstudio/TesselBox-game)](https://github.com/tesselstudio/TesselBox-game/graphs/contributors)

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.