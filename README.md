# TesselBox: Hexagonal Sandbox Adventure 

[![Typing SVG](https://readme-typing-svg.demolab.com/?lines=Experience+TesselBox,+a+captivating+open-source+game+engine+built+in+Go.;Dive+into+a+procedural+hexagonal+world,+where+you+can+mine,+build,+and+explore+endless+landscapes+with+advanced+plugin+support+and+modding+capabilities.)](https://git.io/typing-svg)

[![GitHub stars](https://img.shields.io/github/stars/tesselstudio/TesselBox-game?style=social)](https://github.com/tesselstudio/TesselBox-game/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/tesselstudio/TesselBox-game?style=social)](https://github.com/tesselstudio/TesselBox-game/network/members)
[![GitHub last commit](https://img.shields.io/github/last-commit/tesselstudio/TesselBox-game)](https://github.com/tesselstudio/TesselBox-game/commits/main)
[![GitHub repo size](https://img.shields.io/github/repo-size/tesselstudio/TesselBox-game)](https://github.com/tesselstudio/TesselBox-game/size)
[![GitHub contributors](https://img.shields.io/github/contributors/tesselstudio/TesselBox-game)](https://github.com/tesselstudio/TesselBox-game/graphs/contributors)

## Table of Contents
- [Quick Start](#-quick-start)
- [Gameplay Overview](#-gameplay-overview)
- [Controls](#-controls)
- [Blocks & Biomes](#-blocks--biomes)
- [Custom Block Rendering](#-custom-block-rendering)
- [Skin Editor](#-skin-editor)
- [Plugin Development](#-plugin-development)
- [Game Systems](#-game-systems)
- [Development](#-development)
- [Contributing](#-contributing)

##  Quick Start

### Installation

```bash
# Clone repository
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# Build and run
go run cmd/main.go
```

### System Requirements

#### **Minimum Requirements**
- **OS**: Windows 10, macOS 10.14, or Linux (Ubuntu 18.04+)
- **CPU**: Dual-core processor at 2.0 GHz
- **RAM**: 4 GB RAM
- **Graphics**: OpenGL 3.3 compatible graphics card
- **Storage**: 500 MB available space
- **Audio**: Standard audio output device

#### **Recommended Requirements**
- **OS**: Windows 11, macOS 12, or Linux (Ubuntu 20.04+)
- **CPU**: Quad-core processor at 3.0 GHz
- **RAM**: 8 GB RAM
- **Graphics**: Dedicated graphics card with 2GB VRAM
- **Storage**: 2 GB available space (for saves and plugins)
- **Audio**: Stereo speakers or headphones

Experience TesselBox, a captivating open-source game engine built in Go. Dive into a procedural hexagonal world, where you can mine, build, and explore endless landscapes with advanced plugin support, custom skins, dynamic weather systems, and comprehensive game mechanics.

##  Gameplay Overview

TesselBox offers a captivating experience in an open-source, hexagonal sandbox environment with advanced systems. Players can engage in mining, building, and exploring vast, procedurally generated landscapes with dynamic weather, day/night cycles, and rich audio systems. The game's core features include building unique structures on a hexagonal grid, discovering new terrains through procedural generation, unearthing resources while exploring the world, creating custom content through a powerful plugin system, designing custom player skins, and experiencing immersive environmental systems.

###  World Selection System
TesselBox features a complete Minecraft-style world management system that allows players to:
- **Browse Existing Worlds**: View all saved worlds with metadata (name, seed, game mode, last saved time)
- **Create New Worlds**: Generate fresh worlds with custom names, seeds, and game modes
- **Delete Unwanted Worlds**: Remove worlds from the save list
- **Game Mode Selection**: Choose between Creative and Survival modes
- **Seed Sharing**: Copy and share world seeds for identical terrain generation

###  Enhanced Seed Generation
The game includes a comprehensive seed system for consistent world generation:
- **Deterministic Generation**: Same seed always produces identical worlds
- **Seed Customization**: Set specific seeds or generate random ones
- **Known Seeds Library**: Quick access to popular seeds (12345, 67890, 42, 1337, etc.)
- **Seed Sharing**: Copy seeds to share with friends for identical world experiences
- **World Consistency**: Proper seed usage throughout terrain generation, ore placement, and organism spawning

### Features
-   **Hexagonal World-Building**: Craft unique structures in a distinct hexagonal grid, offering natural movement patterns and strategic depth.
-   **Procedural Generation**: Discover new and exciting terrains with every play, ensuring limitless and diverse exploration.
-   **Mining & Exploration**: Unearth resources and uncover the secrets of the world as you explore its vast landscapes.
-   **Advanced Plugin System**: Create custom content, new game mechanics, and unique experiences with a comprehensive plugin architecture.
-   **Skin Editor System**: Design custom player skins with pixel-by-pixel editing, real-time preview, and persistent storage.
-   **Multi-Color Block System**: Varied block appearances with biome-specific colors, patterns, and dynamic rendering.
-   **Dynamic Weather System**: Experience realistic weather with rain, storms, snow, and clear skies with particle effects.
-   **Day/Night Cycle System**: Immersive lighting system with dawn, morning, noon, afternoon, dusk, and night phases.
-   **Rich Audio System**: Comprehensive audio with music, sound effects, and ambient sounds with 3D positioning.
-   **Advanced Crafting System**: Complex crafting with workstations, recipes, and quality systems.
-   **Save System**: Automatic and manual saving with world state persistence and backup systems.
-   **Input Management**: Configurable controls with action mapping and input device support.
-   **Entity-Component Architecture**: Modern ECS system for efficient entity management and game logic.
-   **Open Source**: Join our community and contribute to the evolution of TesselBox!

###  Survival Mode
- **Basic Implementation**: Survival mode toggle is available (`/survival` command)
- **Creative/Survival Toggle**: Switch between creative building and survival challenges
- **Mode Persistence**: Game mode selection is saved with each world
- **Current Status**: Basic framework exists, but core survival mechanics (health, hunger, stamina) are planned for future implementation
- **Development Focus**: Currently prioritizing world generation and creative building features

##  Controls

### Basic Movement
- **W/A/S/D** - Move player (up/left/down/right)
- **Space** - Jump
- **Shift** - Run (hold for faster movement)
- **F** - Toggle flying mode (creative mode)

### Block Interaction
- **Left Mouse Button** - Break/Mine blocks
- **Right Mouse Button** - Place blocks
- **Mouse Movement** - Look around and aim

### Creative Mode Features
- **B Key** - Open Block Library menu
- **P Key** - Open Plugin Manager
- **C Key** - Open Crafting menu
- **R Key** - Interact with crafting stations
- **Escape** - Close menus and return to game
- **Arrow Keys** - Navigate menu options
- **Enter/Space** - Select menu item
- **Mouse Click** - Select menu items directly

### Block Selection
1. Press **B** to open the block library
2. Navigate with **arrow keys** or **mouse**
3. Select a block with **Enter**, **Space**, or **mouse click**
4. The selected block appears in the bottom-left UI
5. Use **right mouse button** to place the selected block

### UI Information
- **Bottom-left corner** shows currently selected block
- **Instructions** display when in creative mode
- **Hotbar** shows inventory items (survival mode)
- **Crosshair** indicates block placement target
- **Weather indicator** shows current weather conditions
- **Time display** shows current game time

### Game Management
- **Escape** - Pause/Menu
- **Tab** - Toggle debug info (development)
- **F11** - Toggle fullscreen (if supported)
- **F5** - Quick save
- **F9** - Quick load

###  World Selection Controls
- **Arrow Keys** - Navigate world list and menu options
- **Enter/Space** - Select world or confirm action
- **Escape** - Go back to previous menu
- **Mouse Click** - Direct selection of worlds and menu items
- **B Key** - Create new world from world selection menu

###  Seed System Controls
- **Up/Down Arrows** - Cycle through preset world names or known seeds
- **Enter** - Confirm world creation with current settings
- **Random Seed** - Generate new random seed for unique world
- **Copy/Paste** - Share seeds with other players

##  Blocks & Biomes

### Complete Block List (48 Types)

#### Basic Blocks (8 types)
1. **AIR** - Empty space (no variation)
2. **DIRT** - Dirt block with humidity-based color variations
3. **GRASS** - Grass-covered dirt with biome-specific colors
4. **STONE** - Stone block with mineral patterns
5. **SAND** - Sand block with temperature variations
6. **WATER** - Water (liquid) with depth-based transparency
7. **LOG** - Wood log (tropical, temperate, pine variants)
8. **LEAVES** - Tree leaves (hollow blocks with transparency)

#### Ore Blocks (4 types)
9. **COAL_ORE** - Coal ore with brightness variations
10. **IRON_ORE** - Iron ore with rust effects
11. **GOLD_ORE** - Gold ore with shine variations
12. **DIAMOND_ORE** - Diamond ore with sparkle effects

#### Building Blocks (14 types)
13. **BEDROCK** - Unbreakable bottom layer
14. **GLASS** - Transparent glass with fog effects
15. **BRICK** - Brick block with mortar variations
16. **PLANK** - Wood planks with grain patterns
17. **CACTUS** - Cactus block with flowering states
18. **GRAVEL** - Gravel block with stone mix variations
19. **SANDSTONE** - Sandstone with layered gradients
20. **OBSIDIAN** - Obsidian with glassy effects
21. **ICE** - Ice block with crack patterns
22. **SNOW** - Snow block with texture variations
23. **COBBLESTONE** - Cobblestone with weathering
24. **MOSSY_COBBLESTONE** - Mossy cobblestone
25. **STONE_BRICKS** - Stone bricks with mortar
26. **CHISELED_STONE** - Chiseled stone with carvings

#### Functional Blocks (11 types)
27. **WORKBENCH** - Crafting table
28. **FURNACE** - Smelting furnace
29. **ANVIL** - Repair station
30. **CRAFTING_TABLE** - Advanced crafting
31. **CHEST** - Storage chest
32. **LADDER** - Climbing ladder
33. **FENCE** - Barrier fence
34. **GATE** - Fence gate
35. **DOOR** - Door block
36. **WINDOW** - Window block
37. **TORCH** - Light source

#### Decorative/Natural Blocks (11 types)
38. **FLOWER** - Flower decoration
39. **TALL_GRASS** - Tall grass
40. **MUSHROOM_RED** - Red mushroom
41. **MUSHROOM_BROWN** - Brown mushroom
42. **WOOL** - Wool block
43. **BOOKSHELF** - Bookshelf
44. **JUKEBOX** - Music player
45. **NOTE_BLOCK** - Musical block
46. **PUMPKIN** - Pumpkin block
47. **MELON** - Melon block
48. **HAY_BALE** - Hay bale

### Humidity System
All blocks (except bedrock) respond to humidity levels:
- **Dry (0.1-0.3)**: Lighter colors, clean surfaces, sparse details
- **Normal (0.3-0.7)**: Standard colors, moderate textures, normal details
- **Wet (0.7-1.0)**: Darker colors, wet effects, lush details

### Complete Biome List (14 Types)

#### Basic Biomes (4 types)
1. **PLAINS** - Moderate temperature and humidity, grass surface
2. **FOREST** - Cool and humid, high tree density
3. **DESERT** - Hot and dry, sand surface, no trees
4. **BADLANDS** - Hot and dry, red sand, hilly terrain
5. **MOUNTAINS** - Cold and dry, stone surface, rich ore deposits

#### Water Biomes (4 types)
6. **OCEAN** - Moderate temperature, saturated humidity
7. **SWAMP** - Moderate temperature, very humid
8. **CORAL_REEF** - Warm, saturated humidity
9. **MANGROVE** - Warm, very humid, high tree density

#### Cold Biomes (3 types)
10. **TAIGA** - Very cold, moderate humidity
11. **TUNDRA** - Extremely cold, dry, snow surface
12. **ICE_FIELDS** - Freezing, dry, ice surface

#### Hot/Humid Biomes (2 types)
13. **JUNGLE** - Hot, very humid, highest tree density
14. **SAVANNA** - Hot, dry, sparse tree coverage

#### Special Biomes (1 type)
15. **VOLCANIC** - Extremely hot, dry, obsidian surface, very rich ore deposits

### Biome Properties
- **Tree Density**: 0.0 (no trees) to 1.0 (dense forest)
- **Ore Frequency**: Multiplier for ore generation (0.5x = half, 2.0x = double)
- **Temperature**: 0.0 (freezing) to 1.0 (extremely hot)
- **Humidity**: 0.0 (arid) to 1.0 (saturated)

##  Custom Block Rendering

This guide shows you different methods to create custom appearances for individual blocks in TesselBox.

### Method 1: Built-in Color System (Recommended)

The easiest way is to use the existing color variation system:

```go
// Get block color with biome, humidity, and position variations
blockKey := "stone"  // or get from block type
biome := "mountains"
depth := 0.5  // 0.0 = surface, 1.0 = deep underground

blockColor := blocks.GlobalBlockAppearance.GetBlockColor(
    blockKey, 
    int(blockX), 
    int(blockY), 
    biome, 
    depth
)

// Draw with this color
g.drawHexagon(screen, blockX, blockY, blockColor)
```

**Features:**
- Automatic biome-specific colors
- Humidity variations (Dry/Normal/Wet)
- Position-based noise for natural variation
- Depth-based coloring for underground blocks

### Method 2: Custom Texture Generation

Use the custom renderer for special effects:

```go
// Create custom parameters
customParams := map[string]interface{}{
    "roughness": 0.5,  // For stone
    "wave_speed": 2.0, // For water
    "color": color.RGBA{255, 100, 255, 255}, // For crystals
}

// Generate custom texture
texture := blocks.GlobalCustomRenderer.GenerateCustomTexture(
    "stone", 
    int(blockX), 
    int(blockY), 
    customParams
)

// Apply texture to block
g.drawTexturedHexagon(screen, blockX, blockY, texture)
```

**Available Custom Textures:**
- `stone` - Rough stone with adjustable roughness
- `grass` - Grass with blade patterns
- `water` - Animated water with waves
- `custom_crystal` - Crystalline structures with custom colors

### Method 3: Direct Pixel Manipulation

For complete control, draw directly:

```go
func (g *Game) drawSpecialDiamondOre(screen *ebiten.Image, x, y float64) {
    // Base stone color
    stoneColor := color.RGBA{54, 54, 54, 255}
    g.drawHexagon(screen, x, y, stoneColor)
    
    // Add sparkling diamond facets
    diamondColor := color.RGBA{185, 242, 255, 255}
    
    // Draw sparkles at random positions
    for i := 0; i < 5; i++ {
        sparkleX := x + (rand.Float64() - 0.5) * 20
        sparkleY := y + (rand.Float64() - 0.5) * 20
        g.drawSparkle(screen, sparkleX, sparkleY, diamondColor)
    }
}
```

### Integration in Main Game

To use custom rendering, modify the `drawBlock` function in `cmd/main.go`:

```go
func (g *Game) drawBlock(screen *ebiten.Image, block *world.Hexagon) {
    if block.BlockType == blocks.AIR {
        return
    }

    blockKey := getBlockKeyFromType(block.BlockType)
    
    // Use custom rendering for special blocks
    switch blockKey {
    case "diamond_ore":
        g.drawSpecialDiamondOre(screen, block.X, block.Y)
    case "gold_ore":
        g.drawSpecialGoldOre(screen, block.X, block.Y)
    case "water":
        g.drawAnimatedWater(screen, block.X, block.Y)
    default:
        // Use standard rendering with variations
        g.drawBlockWithVariations(screen, block)
    }
}
```

### Advanced Techniques

#### Biome-Specific Wood Types

```go
func (g *Game) getWoodTypeForBiome(biome string) string {
    switch biome {
    case "jungle":
        return "tropical_log"
    case "forest":
        return "temperate_log"  
    case "taiga":
        return "pine_log"
    default:
        return "log"
    }
}
```

#### Humidity-Based Color Modification

```go
func (g *Game) applyHumidityEffect(baseColor color.RGBA, humidity float64) color.RGBA {
    if humidity < 0.3 {
        // Dry: lighter colors
        return color.RGBA{
            R: uint8(float64(baseColor.R) * 1.2),
            G: uint8(float64(baseColor.G) * 1.2), 
            B: uint8(float64(baseColor.B) * 1.2),
            A: baseColor.A,
        }
    } else if humidity > 0.7 {
        // Wet: darker colors
        return color.RGBA{
            R: uint8(float64(baseColor.R) * 0.7),
            G: uint8(float64(baseColor.G) * 0.7),
            B: uint8(float64(baseColor.B) * 0.7), 
            A: baseColor.A,
        }
    }
    return baseColor // Normal humidity
}
```

#### Animated Blocks

```go
func (g *Game) drawAnimatedLava(screen *ebiten.Image, x, y float64, time float64) {
    // Pulsing lava effect
    pulse := math.Sin(time * 0.003) * 0.3 + 0.7
    
    lavaColor := color.RGBA{
        R: uint8(255 * pulse),
        G: uint8(100 * pulse),
        B: uint8(0),
        A: 255,
    }
    
    g.drawHexagon(screen, x, y, lavaColor)
    
    // Add glowing particles
    g.drawLavaParticles(screen, x, y, time)
}
```

### Performance Tips

1. **Cache Textures**: Generate custom textures once and reuse them
2. **Batch Similar Blocks**: Group blocks by appearance for efficient rendering
3. **Limit Complex Effects**: Use custom rendering only for special blocks
4. **LOD System**: Use simpler rendering for distant blocks

##  Skin Editor

TesselBox includes a comprehensive pixel-by-pixel skin editor that allows you to create custom player skins with real-time preview and persistent storage.

###  Drawing Tools
- **Pencil** (B key): Draw pixels with adjustable brush sizes (1-3)
- **Eraser** (E key): Remove pixels
- **Fill** (F key): Flood fill areas with selected color
- **Eyedropper** (I key): Pick colors from the canvas
- **Line** (L key): Draw straight lines
- **Rectangle** (R key): Draw rectangles
- **Circle** (C key): Draw circles

###  Canvas Features
- **64x64 pixel canvas** with zoom support (mouse wheel)
- **Grid overlay** when zoomed in for precise editing
- **Real-time coordinate display**
- **Transparent background support**
- **Multiple brush sizes** (1-3 keys)

###  Preview System
- **Live preview** showing skin on player model
- **Animated preview** with rotation
- **Full-screen preview mode**
- **Real-time updates** as you draw

###  File Management
- **Automatic saving** on exit
- **Skin persistence** across game sessions
- **Multiple skin support**
- **JSON-based storage** in `skins/` directory

###  History & Editing
- **Undo/Redo** (Ctrl+Z/Ctrl+Y)
- **50-step history** buffer
- **Non-destructive editing**

### Skin Editor Controls

#### Mouse Controls
- **Left Click**: Draw/Select tools and colors
- **Mouse Wheel**: Zoom in/out
- **Right Click**: Reset zoom (with Ctrl)

#### Keyboard Shortcuts
- **B**: Pencil tool
- **E**: Eraser tool  
- **F**: Fill tool
- **I**: Eyedropper tool
- **L**: Line tool
- **R**: Rectangle tool
- **C**: Circle tool
- **1-3**: Brush size selection
- **Ctrl+Z**: Undo
- **Ctrl+Y**: Redo
- **ESC**: Save and return to menu
- **Space**: Toggle animation (preview mode)

### How to Use the Skin Editor

1. **Open Skin Editor**: Press **ESC** to open main menu, then select "SKIN EDITOR" (only available in creative mode)
2. **Select Tool**: Use keyboard shortcuts or click tool panel
3. **Choose Color**: Click color palette or use eyedropper tool
4. **Draw**: Click and drag on canvas to draw pixels
5. **Adjust Brush**: Use 1-3 keys for different brush sizes
6. **Zoom**: Use mouse wheel to zoom in/out
7. **Preview**: See real-time preview on player model
8. **Save**: Press ESC to automatically save and return to menu

### Skin Storage Structure
```
skins/
|-- config.json          # Skin configuration
|-- Default.json          # Default skin file
|-- CustomSkin.json       # Custom skin files
`-- temp/                 # Temporary files
```

##  Plugin Development

TesselBox features a comprehensive plugin system that allows you to create custom content, new game mechanics, and unique experiences. The plugin system includes:

- **Entity-Component Integration**: Full access to ECS architecture
- **Event System**: Subscribe to and publish game events
- **World Interaction**: Modify blocks and terrain
- **Configuration Management**: YAML-based plugin configuration
- **Hot Reloading**: Develop plugins without restarting the game
- **Security Framework**: Safe plugin execution with permissions
- **Plugin Marketplace**: Browse and install community plugins
- **UI Plugin Manager**: Manage plugins directly from the game interface

### In-Game Plugin Management

#### Plugin Manager Access
- Press **P key** in creative mode to open the Plugin Manager
- Or select "PLUGIN MANAGER" from the main menu

#### Plugin Manager Features
- **Marketplace**: Browse available plugins by category
- **Installation**: One-click install with progress tracking
- **Management**: Enable/disable installed plugins
- **Search**: Find specific plugins quickly
- **Details**: View plugin information and dependencies
- **Updates**: Automatic plugin updating system

#### Plugin Commands
```bash
# List all plugins
/plugin list

# Load a plugin
/plugin load myplugin

# Get plugin information
/plugin info myplugin

# Reload a plugin
/plugin reload myplugin

# Unload a plugin
/plugin unload myplugin
```

### Getting Started with Plugins

```bash
# Create a new plugin
mkdir plugins/myplugin
cd plugins/myplugin

# Follow the plugin development guide below
```

### Sample Plugin Categories
- **Blocks**: New block types with custom properties
- **Items**: New items with unique behaviors
- **Creatures**: Custom mobs and animals
- **Gameplay**: New game mechanics and systems
- **Visual**: UI enhancements and visual effects
- **Tools**: Utility plugins for developers
- **Audio**: Custom sound effects and music

##  Game Systems

### Multi-Color Block System

TesselBox features an advanced multi-color block system that provides visual variety and dynamic rendering based on environmental factors.

####  Color Variation Types
- **Random Variation**: Subtle color differences for natural appearance
- **Gradient Effects**: Smooth color transitions based on depth or position
- **Pattern Variations**: Textured patterns like stripes, checkerboards
- **Biome-Specific Colors**: Different colors based on world biomes
- **Age-Based Colors**: Blocks that change color over time
- **Moisture-Based**: Color variations based on water proximity
- **Temperature-Based**: Color changes based on climate/temperature

####  Biome System
- **Forest**: Rich greens and browns
- **Plains**: Light greens and yellows
- **Desert**: Dry yellows and oranges
- **Taiga**: Dark greens and grays
- **Tundra**: Pale greens and whites

####  Block Types with Variations
- **Grass**: Multiple shades with biome variations
- **Stone**: Mineral flecks and patterns
- **Sand**: Desert temperature variations
- **Wood**: Different tree types and grain patterns
- **Wool**: Full dye color support
- **Leaves**: Seasonal color changes
- **Water**: Depth-based transparency and color
- **Dirt**: Moisture-based color variations
- **Ice**: Transparency and crack patterns
- **Gravel**: Stone mix variations
- **Sandstone**: Layer color gradients

####  Technical Implementation
- **Procedural Generation**: Real-time color calculation
- **Performance Optimized**: Efficient batch rendering
- **Configurable**: Easy to add new block types and variations
- **Extensible**: Plugin system for custom color schemes

###  Weather System

The weather system provides dynamic environmental conditions that affect gameplay and atmosphere.

####  Weather Types
- **Clear**: Normal visibility with bright skies
- **Rain**: Reduced visibility, wet terrain, particle effects
- **Storm**: Heavy rain, lightning, strong winds, dark skies
- **Snow**: Cold conditions, snow accumulation, winter ambience

####  Weather Effects
- **Particle Systems**: Realistic rain, snow, and wind effects
- **Lighting Changes**: Dynamic sky color and ambient lighting
- **Sound Effects**: Weather-appropriate audio with spatial audio
- **Terrain Interaction**: Wet surfaces, snow accumulation
- **Player Impact**: Movement speed changes, visibility reduction

####  Weather Cycling
- **Automatic Changes**: Weather transitions every 5-15 minutes
- **Gradual Transitions**: Smooth changes between weather types
- **Seasonal Patterns**: Different weather probabilities based on game season
- **Biome Influence**: Weather varies by geographic region

###  Day/Night Cycle System

The day/night cycle provides immersive lighting and time-based gameplay mechanics.

####  Time Phases
- **Dawn**: Soft morning light, 5:00-7:00
- **Morning**: Bright daylight, 7:00-11:00  
- **Noon**: Peak brightness, 11:00-13:00
- **Afternoon**: Warm afternoon light, 13:00-17:00
- **Dusk**: Golden hour sunset, 17:00-19:00
- **Night**: Dark with moonlight, 19:00-23:00
- **Midnight**: Darkest period, 23:00-5:00

####  Lighting Effects
- **Dynamic Sky Colors**: Realistic sunrise/sunset gradients
- **Ambient Lighting**: Brightness changes throughout the day
- **Shadow Rendering**: Dynamic shadows based on sun position
- **Moon Phases**: Realistic moon cycle with different illumination
- **Star Rendering**: Night sky with visible stars and constellations

####  Time Configuration
- **Adjustable Day Length**: Configurable minutes per game day (default: 10 minutes)
- **Pause Function**: Ability to pause time progression
- **Time Display**: Current game time shown in UI
- **Speed Control**: Variable time speed for testing/debugging

###  Audio System

The comprehensive audio system provides immersive soundscapes with music, sound effects, and environmental audio.

####  Audio Types
- **Music**: Background music with different tracks for various situations
- **Sound Effects (SFX)**: Interactive sounds for actions and events
- **Ambient Sounds**: Environmental audio for atmosphere and immersion
- **3D Audio**: Spatial positioning for realistic sound propagation

####  Audio Features
- **Dynamic Mixing**: Automatic volume adjustment based on game state
- **Audio Streaming**: Efficient streaming for large audio files
- **Sound Categorization**: Organized audio libraries for easy management
- **Volume Controls**: Separate controls for music, SFX, and ambient audio
- **Audio Caching**: Intelligent caching for optimal performance

####  Audio Content
- **Music Tracks**: Peaceful exploration, intense combat, ambient themes
- **Sound Effects**: Block placement/breaking, player actions, UI interactions
- **Ambient Sounds**: Weather effects, biome-specific audio, creature sounds
- **Dynamic Audio**: Context-sensitive audio that responds to gameplay

###  Game Systems Architecture

#### Core Systems Overview

##### Entity-Component System (ECS)
```
Entities --> Components --> Systems
```
- **Entities**: Game objects with unique IDs
- **Components**: Data containers (render, physics, inventory)
- **Systems**: Logic processors (rendering, physics, behavior)

##### Plugin System Architecture
```
Enhanced Plugin Manager
|-- Plugin Discovery
|-- Configuration Manager
|-- Hot Reload System
|-- Security Manager
|-- Marketplace Integration
`-- Dependency Resolution
```

##### Rendering Pipeline
```
Game World --> Block Processing --> Entity Rendering --> UI Overlay --> Final Output
```

##### Audio Pipeline
```
Audio Sources --> Mixer --> Spatial Processing --> Output Device
```

##### Save System Architecture
```
Game State --> Serialization --> File Storage --> Compression --> Backup System
```

#### System Interactions

##### World Generation Pipeline
1. **Terrain Generation**: Procedural world creation
2. **Biome Placement**: Environmental zone distribution
3. **Feature Generation**: Structures, caves, resources
4. **Entity Spawning**: Creatures and items placement
5. **Lighting Calculation**: Initial light map generation

##### Game Loop Pipeline
1. **Input Processing**: Handle user input and actions
2. **Entity Updates**: Process entity components and systems
3. **World Updates**: Update blocks, weather, time
4. **Physics Simulation**: Handle collisions and movement
5. **Audio Updates**: Process positional audio
6. **Rendering**: Draw all visual elements
7. **UI Updates**: Update interface elements

##### Save/Load Pipeline
1. **State Capture**: Gather all game state data
2. **Serialization**: Convert to storage format
3. **Compression**: Optimize file size
4. **File Writing**: Store to disk with backup
5. **Verification**: Ensure data integrity

##  Development

### Build System

TesselBox uses a simple Make-based build system with cross-platform support.

#### Local Development
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

#### Features
   Cross-platform builds (Windows, Linux, macOS Intel/ARM)
   Optimized builds (stripped binaries for release)
   Simple Make-based system (no external build scripts)
   Development builds (fast iteration)
   Package creation (tar.gz and zip distributions)

#### Build Output
Builds create binaries in the `bin/` directory:
- `tesselbox-windows-amd64.exe`
- `tesselbox-linux-amd64`
- `tesselbox-darwin-amd64`
- `tesselbox-darwin-arm64`

#### Build Targets

**Core Targets**
- `build` - Build for current platform
- `release` - Build optimized binaries for all platforms
- `clean` - Remove build artifacts
- `dev` - Quick development build

**Platform-specific**
- `windows` - Windows AMD64 binary
- `linux` - Linux AMD64 binary
- `darwin` - macOS Intel binary
- `linux-arm64` - Linux ARM64 binary (native compilation required)
- `darwin-arm64` - macOS ARM64 binary (native compilation required)

**Utilities**
- `test` - Run all tests
- `test-verbose` - Run tests with verbose output
- `test-coverage` - Run tests with coverage
- `deps` - Install build dependencies
- `dist` - Create distribution packages
- `help` - Show available targets

#### Build Notes
- ARM64 builds require native compilation due to Ebiten CGO dependencies
- macOS builds should be performed on macOS hardware
- Icons are generated separately if needed
- Version is set to `v0.3-alpha` by default
- Release builds include optimization flags (`-s -w`, `-trimpath`, version injection)

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

#### Game Configuration
Configuration files are stored in:
- **Windows**: `%APPDATA%/TesselBox/`
- **macOS**: `~/Library/Application Support/TesselBox/`
- **Linux**: `~/.config/tesselbox/`

#### Key Configuration Files
- `config.json`: Main game settings
- `controls.json`: Input bindings
- `audio.json`: Audio settings
- `graphics.json`: Display and rendering options
- `plugins/`: Installed plugins directory
- `skins/`: Custom skins directory
- `saves/`: World save files

##  Contributing

We welcome all contributions! Whether you're a developer, designer, or tester, your input helps us grow. Please see our [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for community guidelines.

**Areas needing contribution**:
-   Block variety system
-   Enhanced mining mechanics
-   Inventory management
-   UI/UX improvements
-   Skin editor enhancements (advanced tools, import/export)
-   Plugin development (create amazing plugins!)
-   Block color system (new biomes, patterns)
-   Bug fixes and optimizations

##  Built With

TesselBox is built with modern technologies to deliver a robust and engaging gaming experience.

<p align="center">
    <img src="https://skillicons.dev/icons?i=go,linux,windows,git,github,markdown,vscode,apple" alt="Technology Stack">
</p>

-   **Go**: The primary programming language for the game engine.
-   **Linux, Windows, Apple**: Supported operating systems for running the game.
-   **Git & GitHub**: Version control and collaborative development platform.
-   **Markdown**: Used for documentation and README formatting.
-   **VS Code**: Recommended development environment.

##  Star History

<a href="https://www.star-history.com/">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/image?repos=tesselstudio/TesselBox-game&type=timeline&theme=dark&legend=bottom-right" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/image?repos=tesselstudio/TesselBox-game&type=timeline&legend=bottom-right" />
   <img alt="Star History Chart" src="https://api.star-history.com/image?repos=tesselstudio/TesselBox-game&type=timeline&legend=bottom-right" />
 </picture>
</a>

##  Contributors

A big thank you to everyone who has contributed to TesselBox!

[![https://contrib.rocks/image?repo=tesselstudio/TesselBox-game](https://contrib.rocks/image?repo=tesselstudio/TesselBox-game)](https://github.com/tesselstudio/TesselBox-game/graphs/contributors)

##  License

This project is licensed under the MIT License. See [LICENSE](LICENSE) file for more details.

---

** Happy Gaming, Plugin Development, and Skin Design! **
