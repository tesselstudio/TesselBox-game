# TesselBox Plugin System

## Overview

The TesselBox game now supports a flexible plugin system that allows content to be modular and removable. This increases flexibility for developers and modders.

## Architecture

### Core Components

The plugin system isolates the following game content into removable plugins:

- **Blocks** - All block types, properties, and behaviors
- **Creatures** - Enemy types, AI, and spawning logic  
- **Organisms** - Plants, trees, and environmental objects
- **Audio** - Sound effects, music, and ambient audio

### Plugin Interface

All plugins implement the `GamePlugin` interface:

```go
type GamePlugin interface {
    // Plugin metadata
    ID() string
    Name() string
    Version() string
    Description() string
    Author() string
    
    // Plugin lifecycle
    Initialize() error
    Shutdown() error
    
    // Content providers
    GetBlockTypes() []blocks.BlockType
    GetBlockDefinition(blockType blocks.BlockType) (*BlockDefinition, bool)
    GetCreatureTypes() []creatures.CreatureType
    GetCreatureDefinition(creatureType creatures.CreatureType) (*CreatureDefinition, bool)
    GetOrganismTypes() []organisms.OrganismType
    GetOrganismDefinition(organismType organisms.OrganismType) (*OrganismDefinition, bool)
    GetAudioTypes() []audio.AudioType
    GetAudioDefinition(audioType audio.AudioType) (*AudioDefinition, bool)
    
    // Game hooks
    GenerateChunk(world *world.World, chunkX, chunkZ int) error
    SpawnOrganisms(world *world.World) error
    SpawnCreatures(world *world.World) error
    OnBlockPlaced(x, y, z int, blockType blocks.BlockType) error
    OnBlockBroken(x, y, z int, blockType blocks.BlockType) error
    OnCreatureSpawn(creature *creatures.Creature) error
    OnCreatureDeath(creature *creatures.Creature) error
    OnTick(world *world.World, deltaTime float64) error
}
```

## Default Plugin

The `DefaultPlugin` contains all original game content:

- **ID**: `default`
- **Name**: `TesselBox Default Content`
- **Version**: `1.0.0`
- **Author**: `TesselBox Team`

### Features

- **96 Block Types**: Air, Dirt, Grass, Stone, Water, Ores, etc.
- **3 Creature Types**: Slime, Spider, Zombie
- **5 Organism Types**: Tree, Bush, Flower, Mushroom, Venus Flytrap
- **3 Audio Types**: Music, SFX, Ambient

### Removal

The default plugin can be **completely removed**:

```go
// Remove default content
pluginManager.DisablePlugin("default")
pluginManager.UnregisterPlugin("default")
```

## Plugin Manager

The `PluginManager` handles plugin lifecycle:

```go
// Create plugin manager
manager := plugins.NewPluginManager()

// Register plugins
manager.RegisterPlugin(plugins.NewDefaultPlugin())
manager.RegisterPlugin(plugins.NewExamplePlugin())

// Enable/disable plugins
manager.EnablePlugin("default")
manager.DisablePlugin("default")

// Get active plugins
activePlugins := manager.GetActivePlugins()
```

## Creating Custom Plugins

### 1. Create Plugin File

```go
package plugins

import (
    "tesselbox/pkg/blocks"
    "tesselbox/pkg/creatures"
    "tesselbox/pkg/organisms"
    "tesselbox/pkg/audio"
    "tesselbox/pkg/world"
)

// MyCustomPlugin implements GamePlugin
type MyCustomPlugin struct {
    initialized bool
}

// Implement all GamePlugin interface methods...
func (mcp *MyCustomPlugin) ID() string {
    return "my_custom"
}

// ... implement other required methods
```

### 2. Plugin Registration

```go
// In main game initialization
func initGamePlugins(manager *plugins.PluginManager) {
    // Register custom plugins
    manager.RegisterPlugin(plugins.NewMyCustomPlugin())
    
    // Enable default plugin (optional)
    manager.EnablePlugin("default")
    
    // Enable custom plugins
    manager.EnablePlugin("my_custom")
}
```

## Benefits

### For Developers

- **Modular Design**: Content is isolated and independent
- **Easy Testing**: Plugins can be enabled/disabled individually
- **Custom Content**: Add new blocks, creatures, etc. without modifying core
- **Version Control**: Different content versions can coexist
- **Hot Reloading**: Plugins can be reloaded during development

### For Players

- **Mod Support**: Players can install custom content
- **Content Variety**: More blocks, creatures, and features
- **Performance**: Unused plugins can be disabled
- **Compatibility**: Multiple content packs can work together

## Migration

### From Old System

The previous monolithic system had:
- Hardcoded content in core packages
- No way to remove default content
- Difficult to add custom content
- Tight coupling between systems

### To New Plugin System

1. **Content Isolated**: All game content moved to plugins
2. **Interface-Based**: Standard plugin interface for all content
3. **Manager Controlled**: Central plugin management system
4. **Backward Compatible**: Default plugin provides all original content

## File Structure

```
pkg/plugins/
├── interface.go          # Plugin interface and manager
├── default.go           # Default game content plugin
├── example_plugin.go    # Example custom plugin
└── README.md            # This documentation
```

## Usage Examples

### Basic Plugin Management

```go
// Initialize plugin system
manager := plugins.NewPluginManager()

// Load default content
defaultPlugin := plugins.NewDefaultPlugin()
manager.RegisterPlugin(defaultPlugin)
manager.EnablePlugin("default")

// Game loop - call plugin hooks
for _, plugin := range manager.GetActivePlugins() {
    plugin.OnTick(world, deltaTime)
    plugin.GenerateChunk(world, chunkX, chunkZ)
}
```

### Custom Content Example

```go
// Create custom blocks plugin
type MagicBlocksPlugin struct{}

func (mbp *MagicBlocksPlugin) GetBlockTypes() []blocks.BlockType {
    return []blocks.BlockType{
        blocks.NewCustomBlock("magic_stone", "#FF00FF"),
        blocks.NewCustomBlock("magic_wood", "#00FF00"),
    }
}

// Register and enable
manager.RegisterPlugin(plugins.NewMagicBlocksPlugin())
manager.EnablePlugin("magic_blocks")
```

## Future Enhancements

- **Plugin Dependencies**: Define dependencies between plugins
- **Configuration Files**: JSON/YAML plugin configs
- **Hot Reloading**: Reload plugins without restart
- **Plugin Marketplace**: Browse and install plugins in-game
- **Version Management**: Automatic plugin updates
- **Sandboxing**: Security isolation for plugins

This plugin system provides maximum flexibility while maintaining compatibility with existing game code.
