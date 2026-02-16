# TesselBox - English README
## Hexagonal Voxel Game

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

A 2D sandbox adventure game inspired by *Terraria*, but built on a **hexagonal grid**.

Explore worlds, mine resources, build structures, craft items, fight enemies, and survive — all in beautiful hex tiles.

## Game Features

### ✅ **Complete Features**
- **Hexagonal World Generation** - Procedurally generated worlds with biomes
- **Mining & Crafting** - Tool-based mining with different material speeds
- **Block Placement** - Right-click to place blocks with ghost preview
- **Inventory System** - 32-slot inventory with hotbar (9 slots)
- **Combat System** - Health/damage system with attack animations
- **Day/Night Cycle** - Dynamic lighting and time progression
- **Weather Effects** - Rain, snow, and storm systems
- **Save/Load System** - Persistent world state with auto-save

### 🎮 **Controls**
- **WASD / Arrow Keys**: Movement
- **Space**: Jump / Attack
- **Left Click**: Mine blocks
- **Right Click**: Place blocks
- **E**: Open crafting menu
- **Q**: Drop selected item
- **Mouse Wheel**: Hotbar selection
- **1-9**: Direct hotbar selection
- **F5**: Manual save
- **F9**: Manual load
- **ESC**: Menu / Close menus

## Installation & Setup

### Prerequisites
- **Go 1.19+** - Core engine
- **Git** - Version control

### Quick Start
```bash
# Clone the repository
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Build the game
go build ./cmd/client

# Run the game
./client
```

### Development Setup
```bash
# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build for development
go build -tags debug ./cmd/client
```

## System Requirements

### Minimum
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: Dual-core processor
- **RAM**: 4GB
- **GPU**: OpenGL 3.3+ compatible
- **Storage**: 500MB free space

### Recommended
- **CPU**: Quad-core processor
- **RAM**: 8GB+
- **GPU**: Dedicated graphics card
- **Storage**: 1GB+ free space

## Architecture

### Core Technologies
- **Language**: Go (Golang)
- **Graphics**: Ebiten (2D game library)
- **Build System**: Go modules

### Project Structure
```
TesselBox/
├── cmd/client/          # Main game executable
├── pkg/                 # Core packages
│   ├── world/          # World generation & management
│   ├── player/         # Player mechanics & physics
│   ├── blocks/         # Block types & properties
│   ├── items/          # Item system & crafting
│   ├── crafting/       # Crafting recipes & UI
│   ├── weather/        # Weather simulation
│   ├── gametime/       # Day/night cycle
│   ├── save/           # Save/load functionality
│   └── render/         # Rendering & UI systems
├── config/             # Configuration files
└── assets/             # Game assets (if any)
```

## Contributing

### For Developers
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go coding standards
- Add tests for new features
- Update documentation
- Ensure cross-platform compatibility

## License

**MIT License** - See [LICENSE](LICENSE) file for details.

## Credits

- **Inspired by**: Terraria game mechanics
- **Built with**: Ebiten game engine
- **Contributors**: Open source community

## Support

- **Issues**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **Discussions**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **Wiki**: [Project Wiki](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*Enjoy exploring the hexagonal world of TesselBox!*
