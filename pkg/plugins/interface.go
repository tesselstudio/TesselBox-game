package plugins

import (
	"fmt"
	"log"
	"tesselbox/pkg/audio"
	"tesselbox/pkg/blocks"
	"tesselbox/pkg/creatures"
	"tesselbox/pkg/organisms"
	"tesselbox/pkg/world"
)

// GamePlugin defines the interface for all game content plugins
type GamePlugin interface {
	// Plugin metadata
	ID() string
	Name() string
	Version() string
	Description() string
	Author() string

	// Plugin lifecycle methods
	Initialize() error
	Shutdown() error

	// Content providers
	GetBlockTypes() []blocks.BlockType
	GetBlockDefinition(blockType blocks.BlockType) (*BlockDefinition, bool)
	GetBlockProperties(blockType blocks.BlockType) (map[string]interface{}, bool)
	GetCreatureTypes() []creatures.CreatureType
	GetCreatureDefinition(creatureType creatures.CreatureType) (*CreatureDefinition, bool)
	GetOrganismTypes() []organisms.OrganismType
	GetOrganismDefinition(organismType organisms.OrganismType) (*OrganismDefinition, bool)
	GetAudioTypes() []audio.AudioType
	GetAudioDefinition(audioType audio.AudioType) (*AudioDefinition, bool)

	// World generation hooks
	GenerateChunk(world *world.World, chunkX, chunkZ int) error
	SpawnOrganisms(world *world.World) error
	SpawnCreatures(world *world.World) error

	// Game hooks
	OnBlockPlaced(x, y, z int, blockType blocks.BlockType) error
	OnBlockBroken(x, y, z int, blockType blocks.BlockType) error
	OnCreatureSpawn(creature *creatures.Creature) error
	OnCreatureDeath(creature *creatures.Creature) error
	OnTick(world *world.World, deltaTime float64) error
}

// PluginManager manages multiple game content plugins
type PluginManager struct {
	plugins map[string]GamePlugin
	active  map[string]bool
}

// NewPluginManager creates a new plugin manager
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]GamePlugin),
		active:  make(map[string]bool),
	}
}

// RegisterPlugin registers a new plugin
func (pm *PluginManager) RegisterPlugin(plugin GamePlugin) error {
	id := plugin.ID()
	if _, exists := pm.plugins[id]; exists {
		return fmt.Errorf("plugin with ID '%s' is already registered", id)
	}

	pm.plugins[id] = plugin
	pm.active[id] = false

	log.Printf("Registered plugin: %s v%s by %s", plugin.Name(), plugin.Version(), plugin.Author())
	return nil
}

// UnregisterPlugin removes a plugin
func (pm *PluginManager) UnregisterPlugin(pluginID string) error {
	if _, exists := pm.plugins[pluginID]; !exists {
		return fmt.Errorf("plugin with ID '%s' is not registered", pluginID)
	}

	// Shutdown plugin if active
	if pm.active[pluginID] {
		pm.plugins[pluginID].Shutdown()
	}

	delete(pm.plugins, pluginID)
	delete(pm.active, pluginID)

	log.Printf("Unregistered plugin: %s", pluginID)
	return nil
}

// EnablePlugin enables a registered plugin
func (pm *PluginManager) EnablePlugin(pluginID string) error {
	plugin, exists := pm.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin with ID '%s' is not registered", pluginID)
	}

	if pm.active[pluginID] {
		return fmt.Errorf("plugin '%s' is already active", pluginID)
	}

	if err := plugin.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize plugin '%s': %w", pluginID, err)
	}

	pm.active[pluginID] = true
	log.Printf("Enabled plugin: %s", pluginID)
	return nil
}

// DisablePlugin disables an active plugin
func (pm *PluginManager) DisablePlugin(pluginID string) error {
	plugin, exists := pm.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin with ID '%s' is not registered", pluginID)
	}

	if !pm.active[pluginID] {
		return fmt.Errorf("plugin '%s' is not active", pluginID)
	}

	if err := plugin.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown plugin '%s': %w", pluginID, err)
	}

	pm.active[pluginID] = false
	log.Printf("Disabled plugin: %s", pluginID)
	return nil
}

// GetActivePlugins returns all active plugins
func (pm *PluginManager) GetActivePlugins() []GamePlugin {
	var activePlugins []GamePlugin
	for id, active := range pm.active {
		if active {
			activePlugins = append(activePlugins, pm.plugins[id])
		}
	}
	return activePlugins
}

// GetPlugin returns a specific plugin by ID
func (pm *PluginManager) GetPlugin(pluginID string) (GamePlugin, bool) {
	plugin, exists := pm.plugins[pluginID]
	return plugin, exists
}

// IsPluginActive checks if a plugin is active
func (pm *PluginManager) IsPluginActive(pluginID string) bool {
	return pm.active[pluginID]
}

// Definition types for plugins
type BlockDefinition struct {
	Type        blocks.BlockType
	Name        string
	Hardness    float64
	Color       string
	Transparent bool
	Solid       bool
}

type CreatureDefinition struct {
	Type   creatures.CreatureType
	Name   string
	Health float64
	Damage float64
	Speed  float64
	Color  string
}

type OrganismDefinition struct {
	Type   organisms.OrganismType
	Name   string
	Height float64
	Color  string
}

type AudioDefinition struct {
	Type   audio.AudioType
	Name   string
	Volume float64
}
