# TesselBox: Hexagonal Sandbox Adventure 

[![Typing SVG](https://readme-typing-svg.demolab.com/?lines=Experience+TesselBox,+a+captivating+open-source+game+engine+built+in+Go.;Dive+into+a+procedural+hexagonal+world,+where+you+can+mine,+build,+and+explore+endless+landscapes+with+advanced+plugin+support+and+modding+capabilities.)](https://git.io/typing-svg)


[![GitHub stars](https://img.shields.io/github/stars/tesselstudio/TesselBox-game?style=social)](https://github.com/tesselstudio/TesselBox-game/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/tesselstudio/TesselBox-game?style=social)](https://github.com/tesselstudio/TesselBox-game/network/members)
[![GitHub last commit](https://img.shields.io/github/last-commit/tesselstudio/TesselBox-game)](https://github.com/tesselstudio/TesselBox-game/commits/main)
[![GitHub repo size](https://img.shields.io/github/repo-size/tesselstudio/TesselBox-game)](https://github.com/tesselstudio/TesselBox-game/size)
[![GitHub contributors](https://img.shields.io/github/contributors/tesselstudio/TesselBox-game)](https://github.com/tesselstudio/TesselBox-game/graphs/contributors)

<p align="center">
  <a href="https://github.com/tesselstudio/TesselBox-game/archive/refs/heads/main.zip">
    <img src="https://raw.githubusercontent.com/tesselstudio/TesselBox-game/main/download.png" alt="Download TesselBox" width="120" style="border-radius: 6px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
  </a>
</p>

TesselBox is an open-source sandbox game engine built in Go. Dive into a procedural hexagonal world where you can mine, build, and explore endless landscapes with advanced plugin support and modding capabilities.

## 🎮 Gameplay Overview

TesselBox offers an immersive experience in an open-source, hexagonal sandbox environment with advanced plugin support. Players engage in mining, building, and exploring vast, procedurally generated landscapes. The game's core features include building unique structures on a hexagonal grid, discovering new terrains through procedural generation, unearthing resources while exploring the world, and creating custom content through a powerful plugin system.

### Features

-   **Hexagonal World-Building**: Craft unique structures in a distinct hexagonal grid, offering natural movement patterns and strategic depth.
-   **Procedural Generation**: Discover new and exciting terrains with every play, ensuring limitless and diverse exploration.
-   **Mining & Exploration**: Unearth resources and uncover the secrets of the world as you explore its vast landscapes.
-   **Advanced Plugin System**: Create custom content, new game mechanics, and unique experiences with a comprehensive plugin architecture.
-   **Entity-Component Architecture**: Modern ECS system for efficient entity management and game logic.
-   **Open Source**: Join our community and contribute to the evolution of TesselBox!

## 🚀 Quick Start

### Installation

```bash
# Clone repository
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Build and run
go run cmd/main.go
```

### Development

```bash
# Development build
make dev

# Full build
make build

# Run tests
make test
```

## 🔌 Plugin Development

TesselBox features a comprehensive plugin system that allows you to create custom content, new game mechanics, and unique experiences. The plugin system includes:

- **Entity-Component Integration**: Full access to ECS architecture
- **Event System**: Subscribe to and publish game events
- **World Interaction**: Modify blocks and terrain
- **Configuration Management**: YAML-based plugin configuration
- **Hot Reloading**: Develop plugins without restarting the game
- **Security Framework**: Safe plugin execution with permissions

### Getting Started with Plugins

```bash
# Create a new plugin
mkdir plugins/myplugin
cd plugins/myplugin

# Follow the plugin development guide below
```

### In-Game Plugin Management

```bash
# List all plugins
/plugin list

# Load a plugin
/plugin load myplugin

# Get plugin information
/plugin info myplugin

# Reload a plugin
/plugin reload myplugin
```

## 🏗️ Architecture

### Entity-Component System (ECS)

```
Entities → Components → Systems
```

- **Entities**: Game objects with unique IDs
- **Components**: Data containers (render, physics, inventory)
- **Systems**: Logic processors (rendering, physics, behavior)

### Plugin System Architecture

```
Enhanced Plugin Manager
├── Plugin Discovery
├── Configuration Manager
├── Hot Reload System
└── Security Manager
```

## 🛠️ Development

### Build System

```bash
# Build for current platform
make build

# Build for all platforms
make release

# Cross-platform builds
make windows
make linux
make darwin
```

### Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific test types
make test-unit
make test-integration
make test-migration
```

## 🤝 Contributing

We welcome all contributions! Whether you're a developer, designer, or tester, your input helps us grow. Please see our [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for community guidelines.

**Areas needing contribution**:
- 🎮 Block variety system
- ⛏️ Enhanced mining mechanics
- 📦 Inventory management
- 🎨 UI/UX improvements
- 🔌 Plugin development (create amazing plugins!)
- 🐛 Bug fixes and optimizations

## 🛠️ Built With

TesselBox is built with modern technologies to deliver a robust gaming experience across platforms.

<p align="center">
    <img src="https://skillicons.dev/icons?i=go,linux,windows,apple,git,github,markdown,vscode" alt="Technology Stack">
</p>

-   **Go**: High-performance language for the game engine
-   **Cross-Platform**: Linux, Windows, and macOS support
-   **Git & GitHub**: Version control and collaboration
-   **Markdown**: Documentation format
-   **VS Code**: Recommended IDE

## 📊 Star History

<a href="https://www.star-history.com/">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/image?repos=tesselstudio/TesselBox-game&type=timeline&theme=dark&legend=bottom-right" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/image?repos=tesselstudio/TesselBox-game&type=timeline&legend=bottom-right" />
   <img alt="Star History Chart" src="https://api.star-history.com/image?repos=tesselstudio/TesselBox-game&type=timeline&legend=bottom-right" />
 </picture>
</a>

## 👥 Contributors

A big thank you to everyone who has contributed to TesselBox!

[![https://contrib.rocks/image?repo=tesselstudio/TesselBox-game](https://contrib.rocks/image?repo=tesselstudio/TesselBox-game)](https://github.com/tesselstudio/TesselBox-game/graphs/contributors)

## 📜 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) file for more details.

---



**🎮 Happy Gaming and Plugin Development! 🚀**
