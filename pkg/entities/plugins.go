package entities

import (
	"fmt"
	"log"
	"path/filepath"
	"plugin"
	"sync"
)

// Plugin defines the interface for all plugins
type Plugin interface {
	GetName() string
	GetVersion() string
	GetDescription() string
	GetAuthor() string
	Initialize(manager *PluginManager) error
	Shutdown() error
	GetComponents() []Component
	GetSystems() []System
	GetTemplates() map[string]*EntityTemplate
	GetDependencies() []string
}

// PluginInfo contains metadata about a plugin
type PluginInfo struct {
	Name         string   `yaml:"name"`
	Version      string   `yaml:"version"`
	Description  string   `yaml:"description"`
	Author       string   `yaml:"author"`
	Website      string   `yaml:"website,omitempty"`
	License      string   `yaml:"license"`
	Dependencies []string `yaml:"dependencies,omitempty"`
	MinVersion   string   `yaml:"minVersion,omitempty"`
	MaxVersion   string   `yaml:"maxVersion,omitempty"`
	Enabled      bool     `yaml:"enabled"`
}

// PluginManager manages loading and unloading of plugins
type PluginManager struct {
	plugins       map[string]Plugin
	pluginInfo    map[string]*PluginInfo
	loadedPlugins map[string]*plugin.Plugin
	entityManager *EntityManager
	systemManager *SystemManager
	eventBus      *EventBus
	mutex         sync.RWMutex
	pluginPath    string
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(entityManager *EntityManager, systemManager *SystemManager, eventBus *EventBus) *PluginManager {
	return &PluginManager{
		plugins:       make(map[string]Plugin),
		pluginInfo:    make(map[string]*PluginInfo),
		loadedPlugins: make(map[string]*plugin.Plugin),
		entityManager: entityManager,
		systemManager: systemManager,
		eventBus:      eventBus,
		pluginPath:    "plugins",
	}
}

// SetPluginPath sets the path to look for plugins
func (pm *PluginManager) SetPluginPath(path string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.pluginPath = path
}

// LoadPlugin loads a plugin by name
func (pm *PluginManager) LoadPlugin(pluginName string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Check if plugin is already loaded
	if _, exists := pm.plugins[pluginName]; exists {
		return fmt.Errorf("plugin %s is already loaded", pluginName)
	}

	// Load plugin from file
	pluginPath := filepath.Join(pm.pluginPath, pluginName+".so")
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %v", pluginName, err)
	}

	// Look for the plugin symbol
	sym, err := plug.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin %s does not export a Plugin symbol: %v", pluginName, err)
	}

	// Convert symbol to Plugin interface
	pluginInstance, ok := sym.(Plugin)
	if !ok {
		return fmt.Errorf("plugin %s does not implement the Plugin interface", pluginName)
	}

	// Check dependencies
	deps := pluginInstance.GetDependencies()
	for _, dep := range deps {
		if _, exists := pm.plugins[dep]; !exists {
			return fmt.Errorf("plugin %s requires dependency %s which is not loaded", pluginName, dep)
		}
	}

	// Initialize plugin
	err = pluginInstance.Initialize(pm)
	if err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %v", pluginName, err)
	}

	// Register plugin components
	components := pluginInstance.GetComponents()
	for _, component := range components {
		RegisterComponent(component.GetType(), component)
	}

	// Register plugin systems
	systems := pluginInstance.GetSystems()
	for _, system := range systems {
		pm.systemManager.RegisterSystem(system)
	}

	// Register plugin templates
	templates := pluginInstance.GetTemplates()
	for templateID := range templates {
		// This would need to be implemented in EntityManager
		log.Printf("Registered template from plugin %s: %s", pluginName, templateID)
	}

	// Store plugin
	pm.plugins[pluginName] = pluginInstance
	pm.loadedPlugins[pluginName] = plug

	log.Printf("Loaded plugin: %s v%s", pluginInstance.GetName(), pluginInstance.GetVersion())
	return nil
}

// UnloadPlugin unloads a plugin by name
func (pm *PluginManager) UnloadPlugin(pluginName string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	plugin, exists := pm.plugins[pluginName]
	if !exists {
		return fmt.Errorf("plugin %s is not loaded", pluginName)
	}

	// Check if other plugins depend on this one
	for name, p := range pm.plugins {
		deps := p.GetDependencies()
		for _, dep := range deps {
			if dep == pluginName {
				return fmt.Errorf("cannot unload plugin %s because %s depends on it", pluginName, name)
			}
		}
	}

	// Shutdown plugin
	err := plugin.Shutdown()
	if err != nil {
		log.Printf("Warning: plugin %s shutdown failed: %v", pluginName, err)
	}

	// Remove plugin systems
	systems := plugin.GetSystems()
	for _, system := range systems {
		pm.systemManager.UnregisterSystem(system.GetName())
	}

	// Remove plugin
	delete(pm.plugins, pluginName)
	delete(pm.loadedPlugins, pluginName)

	log.Printf("Unloaded plugin: %s", pluginName)
	return nil
}

// ReloadPlugin reloads a plugin
func (pm *PluginManager) ReloadPlugin(pluginName string) error {
	if err := pm.UnloadPlugin(pluginName); err != nil {
		return err
	}
	return pm.LoadPlugin(pluginName)
}

// GetPlugin returns a plugin by name
func (pm *PluginManager) GetPlugin(pluginName string) (Plugin, bool) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	plugin, exists := pm.plugins[pluginName]
	return plugin, exists
}

// ListPlugins returns a list of loaded plugin names
func (pm *PluginManager) ListPlugins() []string {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	plugins := make([]string, 0, len(pm.plugins))
	for name := range pm.plugins {
		plugins = append(plugins, name)
	}
	return plugins
}

// GetPluginInfo returns information about a plugin
func (pm *PluginManager) GetPluginInfo(pluginName string) (*PluginInfo, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	info, exists := pm.pluginInfo[pluginName]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginName)
	}
	return info, nil
}

// IsLoaded checks if a plugin is loaded
func (pm *PluginManager) IsLoaded(pluginName string) bool {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	_, exists := pm.plugins[pluginName]
	return exists
}

// LoadAllPlugins loads all plugins from the plugin directory
func (pm *PluginManager) LoadAllPlugins() error {
	// This would scan the plugin directory and load all .so files
	// For now, we'll just log that this would be implemented
	log.Printf("Plugin loading from directory %s would be implemented here", pm.pluginPath)
	return nil
}

// UnloadAllPlugins unloads all plugins
func (pm *PluginManager) UnloadAllPlugins() error {
	plugins := pm.ListPlugins()
	for _, pluginName := range plugins {
		if err := pm.UnloadPlugin(pluginName); err != nil {
			log.Printf("Failed to unload plugin %s: %v", pluginName, err)
		}
	}
	return nil
}

// ============================================================================
// Plugin Base Implementation
// ============================================================================

// BasePlugin provides a base implementation for plugins
type BasePlugin struct {
	name         string
	version      string
	description  string
	author       string
	initialized  bool
	components   []Component
	systems      []System
	templates    map[string]*EntityTemplate
	dependencies []string
}

// NewBasePlugin creates a new base plugin
func NewBasePlugin(name, version, description, author string) *BasePlugin {
	return &BasePlugin{
		name:         name,
		version:      version,
		description:  description,
		author:       author,
		initialized:  false,
		components:   make([]Component, 0),
		systems:      make([]System, 0),
		templates:    make(map[string]*EntityTemplate),
		dependencies: make([]string, 0),
	}
}

func (bp *BasePlugin) GetName() string           { return bp.name }
func (bp *BasePlugin) GetVersion() string        { return bp.version }
func (bp *BasePlugin) GetDescription() string    { return bp.description }
func (bp *BasePlugin) GetAuthor() string         { return bp.author }
func (bp *BasePlugin) GetDependencies() []string { return bp.dependencies }

func (bp *BasePlugin) Initialize(manager *PluginManager) error {
	if bp.initialized {
		return fmt.Errorf("plugin %s is already initialized", bp.name)
	}

	bp.initialized = true
	log.Printf("Initialized plugin: %s", bp.name)
	return nil
}

func (bp *BasePlugin) Shutdown() error {
	if !bp.initialized {
		return fmt.Errorf("plugin %s is not initialized", bp.name)
	}

	bp.initialized = false
	log.Printf("Shutdown plugin: %s", bp.name)
	return nil
}

func (bp *BasePlugin) GetComponents() []Component               { return bp.components }
func (bp *BasePlugin) GetSystems() []System                     { return bp.systems }
func (bp *BasePlugin) GetTemplates() map[string]*EntityTemplate { return bp.templates }

// AddComponent adds a component to the plugin
func (bp *BasePlugin) AddComponent(component Component) {
	bp.components = append(bp.components, component)
}

// AddSystem adds a system to the plugin
func (bp *BasePlugin) AddSystem(system System) {
	bp.systems = append(bp.systems, system)
}

// AddTemplate adds a template to the plugin
func (bp *BasePlugin) AddTemplate(templateID string, template *EntityTemplate) {
	bp.templates[templateID] = template
}

// AddDependency adds a dependency to the plugin
func (bp *BasePlugin) AddDependency(dependency string) {
	bp.dependencies = append(bp.dependencies, dependency)
}

// ============================================================================
// Example Plugin Implementations
// ============================================================================

// MagicPlugin is an example plugin that adds magic components
type MagicPlugin struct {
	*BasePlugin
}

func NewMagicPlugin() *MagicPlugin {
	base := NewBasePlugin(
		"magic",
		"1.0.0",
		"Adds magic components and systems to the game",
		"TesselBox Team",
	)

	plugin := &MagicPlugin{BasePlugin: base}

	// Add magic components
	plugin.AddComponent(&MagicComponent{})

	// Add magic system
	plugin.AddSystem(NewMagicSystem())

	// Add magic templates
	plugin.AddTemplate("magic_wand", &EntityTemplate{
		ID:   "magic_wand",
		Type: "item",
		Name: "Magic Wand",
		Tags: []string{"item", "tool", "magic"},
		Components: map[string]interface{}{
			"render": map[string]interface{}{
				"type":    "render",
				"color":   []uint8{255, 0, 255, 255},
				"pattern": "solid",
				"scale":   0.8,
			},
			"tool": map[string]interface{}{
				"type":       "tool",
				"toolType":   "magic",
				"power":      8.0,
				"efficiency": 4.0,
				"durability": 100,
				"effective":  []string{"all"},
			},
			"magic": map[string]interface{}{
				"type":       "magic",
				"manaCost":   10,
				"spellPower": 15,
				"spells":     []string{"fireball", "teleport"},
			},
		},
	})

	return plugin
}

// MagicComponent represents magic properties
type MagicComponent struct {
	Type       string   `yaml:"type"`
	ManaCost   int      `yaml:"manaCost"`
	SpellPower int      `yaml:"spellPower"`
	Spells     []string `yaml:"spells"`
}

func (c *MagicComponent) GetType() string { return "magic" }
func (c *MagicComponent) Clone() Component {
	clone := *c
	if len(c.Spells) > 0 {
		clone.Spells = make([]string, len(c.Spells))
		copy(clone.Spells, c.Spells)
	}
	return &clone
}
func (c *MagicComponent) Merge(other Component) {
	if mc, ok := other.(*MagicComponent); ok {
		if mc.ManaCost != 0 {
			c.ManaCost = mc.ManaCost
		}
		if mc.SpellPower != 0 {
			c.SpellPower = mc.SpellPower
		}
		if len(mc.Spells) > 0 {
			c.Spells = mc.Spells
		}
	}
}
func (c *MagicComponent) Validate() error {
	if c.ManaCost < 0 {
		return fmt.Errorf("mana cost cannot be negative")
	}
	if c.SpellPower < 0 {
		return fmt.Errorf("spell power cannot be negative")
	}
	return nil
}

// MagicSystem handles magic-related logic
type MagicSystem struct {
	name               string
	requiredComponents []string
}

func NewMagicSystem() *MagicSystem {
	return &MagicSystem{
		name:               "magic",
		requiredComponents: []string{"magic"},
	}
}

func (ms *MagicSystem) GetName() string                 { return ms.name }
func (ms *MagicSystem) GetRequiredComponents() []string { return ms.requiredComponents }
func (ms *MagicSystem) Matches(entity *Entity) bool {
	return entity.HasComponent("magic")
}

func (ms *MagicSystem) Process(deltaTime float64, entities []*Entity) {
	for _, entity := range entities {
		if !ms.Matches(entity) {
			continue
		}

		magicComp, _ := entity.GetComponent("magic")
		if magic, ok := magicComp.(*MagicComponent); ok {
			// Process magic logic
			_ = magic // Magic processing would go here
		}
	}
}

// TechPlugin is an example plugin that adds technology components
type TechPlugin struct {
	*BasePlugin
}

func NewTechPlugin() *TechPlugin {
	base := NewBasePlugin(
		"tech",
		"1.0.0",
		"Adds technology components and systems to the game",
		"TesselBox Team",
	)

	plugin := &TechPlugin{BasePlugin: base}

	// Add tech components
	plugin.AddComponent(&TechComponent{})

	// Add tech system
	plugin.AddSystem(NewTechSystem())

	return plugin
}

// TechComponent represents technology properties
type TechComponent struct {
	Type       string   `yaml:"type"`
	PowerLevel int      `yaml:"powerLevel"`
	Efficiency float64  `yaml:"efficiency"`
	TechType   string   `yaml:"techType"`
	Upgrades   []string `yaml:"upgrades"`
}

func (c *TechComponent) GetType() string { return "tech" }
func (c *TechComponent) Clone() Component {
	clone := *c
	if len(c.Upgrades) > 0 {
		clone.Upgrades = make([]string, len(c.Upgrades))
		copy(clone.Upgrades, c.Upgrades)
	}
	return &clone
}
func (c *TechComponent) Merge(other Component) {
	if tc, ok := other.(*TechComponent); ok {
		if tc.PowerLevel != 0 {
			c.PowerLevel = tc.PowerLevel
		}
		if tc.Efficiency != 0 {
			c.Efficiency = tc.Efficiency
		}
		if tc.TechType != "" {
			c.TechType = tc.TechType
		}
		if len(tc.Upgrades) > 0 {
			c.Upgrades = tc.Upgrades
		}
	}
}
func (c *TechComponent) Validate() error {
	if c.PowerLevel < 0 {
		return fmt.Errorf("power level cannot be negative")
	}
	if c.Efficiency < 0 || c.Efficiency > 1 {
		return fmt.Errorf("efficiency must be between 0 and 1")
	}
	return nil
}

// TechSystem handles technology-related logic
type TechSystem struct {
	name               string
	requiredComponents []string
}

func NewTechSystem() *TechSystem {
	return &TechSystem{
		name:               "tech",
		requiredComponents: []string{"tech"},
	}
}

func (ts *TechSystem) GetName() string                 { return ts.name }
func (ts *TechSystem) GetRequiredComponents() []string { return ts.requiredComponents }
func (ts *TechSystem) Matches(entity *Entity) bool {
	return entity.HasComponent("tech")
}

func (ts *TechSystem) Process(deltaTime float64, entities []*Entity) {
	for _, entity := range entities {
		if !ts.Matches(entity) {
			continue
		}

		techComp, _ := entity.GetComponent("tech")
		if tech, ok := techComp.(*TechComponent); ok {
			// Process technology logic
			_ = tech // Tech processing would go here
		}
	}
}

// ============================================================================
// Plugin Factory
// ============================================================================

// PluginFactory creates plugin instances
type PluginFactory struct {
	creators map[string]func() Plugin
}

// NewPluginFactory creates a new plugin factory
func NewPluginFactory() *PluginFactory {
	factory := &PluginFactory{
		creators: make(map[string]func() Plugin),
	}

	// Register built-in plugins
	factory.RegisterCreator("magic", func() Plugin { return NewMagicPlugin() })
	factory.RegisterCreator("tech", func() Plugin { return NewTechPlugin() })

	return factory
}

// RegisterCreator registers a plugin creator function
func (pf *PluginFactory) RegisterCreator(name string, creator func() Plugin) {
	pf.creators[name] = creator
}

// CreatePlugin creates a plugin by name
func (pf *PluginFactory) CreatePlugin(name string) (Plugin, error) {
	creator, exists := pf.creators[name]
	if !exists {
		return nil, fmt.Errorf("unknown plugin type: %s", name)
	}
	return creator(), nil
}

// ListPlugins returns a list of available plugin names
func (pf *PluginFactory) ListPlugins() []string {
	plugins := make([]string, 0, len(pf.creators))
	for name := range pf.creators {
		plugins = append(plugins, name)
	}
	return plugins
}
