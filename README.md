# TesselBox

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey)](https://github.com/tesselstudio/TesselBox-game)

## Quick Start

```bash
# Clone and run
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game
go run cmd/main.go
```

**Requirements:** Go 1.21+, OpenGL 3.3, 4GB RAM

## Features

- **Hexagonal Grid** - Unique 6-sided world-building with natural movement patterns
- **Procedural Worlds** - Deterministic generation with customizable seeds
- **Dual Game Modes** - Creative (unlimited building) and Survival (resource management, health, combat)
- **Dynamic Systems** - Day/night cycle, weather, biomes, temperature, humidity
- **Combat & Survival** - Zombie enemies, locational health, equipment/armor, weapon swinging
- **Crafting** - Workbenches, furnaces, anvils with recipe unlocks
- **Storage** - 27-slot chests with full inventory management
- **Plugins** - Extend with custom blocks, items, creatures, and mechanics
- **Skin Editor** - 64x64 pixel editor with real-time preview
- **Save System** - Persistent worlds with player data, chest contents, and enemy states

## Controls

| Key | Action |
|-----|--------|
| W/A/S/D | Move |
| Space | Jump |
| Shift | Sprint |
| F | Toggle fly (if wings equipped) |
| LMB | Mine / Attack |
| RMB | Place block |
| B | Block library (creative) |
| C | Crafting menu |
| E | Inventory |
| Q | Drop item |
| R | Interact with crafting station |
| I | Go to deeper layer (middle → back) |
| K | Go to surface layer (back → middle → surface) |
| M | Toggle music |
| P | Open plugin manager |
| H | Toggle controls help display |
| ESC | Menu / Close |
| F5 | Quick save |
| F6 | Backup save |
| F7 | Save info (shows stats in console) |
| F9 | Quick load |
| F11 | Toggle profiler (debug) |
| F12 | Toggle debug info display |

## Development

```bash
# Build for current platform
make build

# Build for all platforms (release)
make release

# Quick dev build
make dev

# Run tests
make test
```

**Build Targets:** `windows`, `linux`, `darwin`, `linux-arm64`, `darwin-arm64`

## Contributing

We welcome contributions! See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for guidelines.

**Areas needing help:**
- Block variety & biomes
- Enhanced mining mechanics
- UI/UX improvements
- Plugin development
- Bug fixes

## Star History

<a href="https://www.star-history.com/">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/image?repos=tesselstudio/TesselBox-game&type=timeline&theme=dark&legend=bottom-right" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/image?repos=tesselstudio/TesselBox-game&type=timeline&legend=bottom-right" />
   <img alt="Star History Chart" src="https://api.star-history.com/image?repos=tesselstudio/TesselBox-game&type=timeline&legend=bottom-right" />
 </picture>
</a>

## Contributors

A big thank you to everyone who has contributed to TesselBox!

[![https://contrib.rocks/image?repo=tesselstudio/TesselBox-game](https://contrib.rocks/image?repo=tesselstudio/TesselBox-game)](https://github.com/tesselstudio/TesselBox-game/graphs/contributors)

##  License

This project is licensed under the MIT License. See [LICENSE](LICENSE) file for more details.

---

**Happy Gaming, Plugin Development, and Skin Design!**
