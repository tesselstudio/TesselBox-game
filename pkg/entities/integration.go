package entities

import (
	"fmt"
	"log"
)

// GameWorld represents the integrated game world using the new entity system
type GameWorld struct {
	entityManager     *EntityManager
	systemManager     *SystemManager
	pluginManager     *PluginManager
	enhancedPluginMgr *EnhancedPluginManager
	configManager     *PluginConfigManager
	dataLoader        *DataLoader
	eventBus          *EventBus
	initialized       bool
}

// NewGameWorld creates a new integrated game world
func NewGameWorld() *GameWorld {
	// Create core systems
	eventBus := NewEventBus()
	entityManager := NewEntityManager()
	systemManager := NewSystemManager()
	pluginManager := NewPluginManager(entityManager, systemManager, eventBus)
	enhancedPluginMgr := NewEnhancedPluginManager(entityManager, systemManager, eventBus)
	configManager := NewPluginConfigManager("config")
	dataLoader := NewDataLoader(entityManager)

	// Create the integrated world
	world := &GameWorld{
		entityManager:     entityManager,
		systemManager:     systemManager,
		pluginManager:     pluginManager,
		enhancedPluginMgr: enhancedPluginMgr,
		configManager:     configManager,
		dataLoader:        dataLoader,
		eventBus:          eventBus,
		initialized:       false,
	}

	return world
}

// Initialize initializes the game world
func (gw *GameWorld) Initialize() error {
	if gw.initialized {
		return fmt.Errorf("game world is already initialized")
	}

	log.Println("Initializing game world...")

	// Load global plugin configuration
	if err := gw.configManager.LoadGlobalConfig(); err != nil {
		log.Printf("Warning: Failed to load global plugin config: %v", err)
	}

	// Register core systems
	gw.registerCoreSystems()

	// Load all entity data
	if err := gw.dataLoader.LoadAll(); err != nil {
		return fmt.Errorf("failed to load entity data: %v", err)
	}

	// Initialize enhanced plugin system
	if err := gw.initializePluginSystem(); err != nil {
		return fmt.Errorf("failed to initialize plugin system: %v", err)
	}

	// Set up event listeners
	gw.setupEventListeners()

	gw.initialized = true
	log.Println("Game world initialized successfully")

	return nil
}

// registerCoreSystems registers all core systems
func (gw *GameWorld) registerCoreSystems() {
	// Register core systems
	gw.systemManager.RegisterSystem(NewRenderSystem())
	gw.systemManager.RegisterSystem(NewPhysicsSystem(9.8))
	gw.systemManager.RegisterSystem(NewBehaviorSystem())
	gw.systemManager.RegisterSystem(NewInventorySystem())
	gw.systemManager.RegisterSystem(NewCraftingSystem())
	gw.systemManager.RegisterSystem(NewToolSystem())
	gw.systemManager.RegisterSystem(NewCombatSystem())

	log.Println("Registered core systems")
}

// initializePluginSystem initializes the enhanced plugin system
func (gw *GameWorld) initializePluginSystem() error {
	log.Println("Initializing enhanced plugin system...")

	// Get global configuration
	globalConfig := gw.configManager.GetGlobalConfig()

	// Set plugin directories
	pluginDirs := []string{"plugins"}
	if globalConfig.PluginDirectory != "" {
		pluginDirs = append(pluginDirs, globalConfig.PluginDirectory)
	}

	for _, dir := range pluginDirs {
		if err := gw.enhancedPluginMgr.AddPluginDirectory(dir); err != nil {
			log.Printf("Warning: Failed to add plugin directory %s: %v", dir, err)
		}
	}

	// Enable hot reload if configured
	gw.enhancedPluginMgr.EnableHotReload(globalConfig.HotReload)

	// Discover and load plugins
	if err := gw.enhancedPluginMgr.DiscoverAndLoad(); err != nil {
		log.Printf("Warning: Failed to discover and load plugins: %v", err)
	}

	// Load built-in plugins for compatibility
	if err := gw.loadBuiltinPlugins(); err != nil {
		log.Printf("Warning: Failed to load built-in plugins: %v", err)
	}

	log.Println("Enhanced plugin system initialized successfully")
	return nil
}

// loadBuiltinPlugins loads built-in plugins
func (gw *GameWorld) loadBuiltinPlugins() error {
	factory := NewPluginFactory()
	availablePlugins := factory.ListPlugins()

	for _, pluginName := range availablePlugins {
		plugin, err := factory.CreatePlugin(pluginName)
		if err != nil {
			log.Printf("Failed to create plugin %s: %v", pluginName, err)
			continue
		}

		// Initialize plugin directly (bypassing plugin manager for built-ins)
		err = plugin.Initialize(gw.pluginManager)
		if err != nil {
			log.Printf("Failed to initialize plugin %s: %v", pluginName, err)
			continue
		}

		// Register plugin components
		components := plugin.GetComponents()
		for _, component := range components {
			RegisterComponent(component.GetType(), component)
		}

		// Register plugin systems
		systems := plugin.GetSystems()
		for _, system := range systems {
			gw.systemManager.RegisterSystem(system)
		}

		// Register plugin templates
		templates := plugin.GetTemplates()
		for templateID := range templates {
			// This would need to be implemented in EntityManager
			log.Printf("Registered template from built-in plugin: %s", templateID)
		}

		log.Printf("Loaded built-in plugin: %s v%s", plugin.GetName(), plugin.GetVersion())
	}

	return nil
}

// setupEventListeners sets up event listeners
func (gw *GameWorld) setupEventListeners() {
	// Entity events
	gw.eventBus.Subscribe(EventEntityAdded, func(event Event) {
		if entityEvent, ok := event.Data.(EntityEvent); ok {
			log.Printf("Entity added: %s (%s)", entityEvent.Entity.ID, entityEvent.Entity.Type)
		}
	})

	gw.eventBus.Subscribe(EventEntityRemoved, func(event Event) {
		if entityEvent, ok := event.Data.(EntityEvent); ok {
			log.Printf("Entity removed: %s (%s)", entityEvent.Entity.ID, entityEvent.Entity.Type)
		}
	})

	// Block events
	gw.eventBus.Subscribe(EventBlockPlaced, func(event Event) {
		if blockEvent, ok := event.Data.(BlockEvent); ok {
			log.Printf("Block placed: %s at (%.1f, %.1f, %.1f)",
				blockEvent.BlockType, blockEvent.Position.X, blockEvent.Position.Y, blockEvent.Position.Z)
		}
	})

	gw.eventBus.Subscribe(EventBlockBroken, func(event Event) {
		if blockEvent, ok := event.Data.(BlockEvent); ok {
			log.Printf("Block broken: %s at (%.1f, %.1f, %.1f)",
				blockEvent.BlockType, blockEvent.Position.X, blockEvent.Position.Y, blockEvent.Position.Z)
		}
	})

	// Item events
	gw.eventBus.Subscribe(EventItemUsed, func(event Event) {
		if itemEvent, ok := event.Data.(ItemEvent); ok {
			log.Printf("Item used: %s x%d by %s", itemEvent.ItemType, itemEvent.Quantity, itemEvent.PlayerID)
		}
	})

	gw.eventBus.Subscribe(EventItemCrafted, func(event Event) {
		if itemEvent, ok := event.Data.(ItemEvent); ok {
			log.Printf("Item crafted: %s x%d by %s", itemEvent.ItemType, itemEvent.Quantity, itemEvent.PlayerID)
		}
	})

	// Combat events
	gw.eventBus.Subscribe(EventAttack, func(event Event) {
		if combatEvent, ok := event.Data.(CombatEvent); ok {
			log.Printf("Attack: %s -> %s (%.1f damage)",
				combatEvent.AttackerID, combatEvent.TargetID, combatEvent.Damage)
		}
	})

	gw.eventBus.Subscribe(EventDeath, func(event Event) {
		if combatEvent, ok := event.Data.(CombatEvent); ok {
			log.Printf("Death: %s killed by %s", combatEvent.TargetID, combatEvent.AttackerID)
		}
	})

	log.Println("Event listeners set up")
}

// Update updates the game world
func (gw *GameWorld) Update(deltaTime float64) {
	if !gw.initialized {
		return
	}

	// Process events
	gw.eventBus.ProcessBatch()

	// Update all systems
	gw.systemManager.Update(deltaTime)
}

// Shutdown shuts down the game world
func (gw *GameWorld) Shutdown() error {
	if !gw.initialized {
		return nil
	}

	log.Println("Shutting down game world...")

	// Unload all plugins
	err := gw.pluginManager.UnloadAllPlugins()
	if err != nil {
		log.Printf("Error unloading plugins: %v", err)
	}

	// Clear event bus
	gw.eventBus.Clear()

	// Clear entities
	gw.entityManager = NewEntityManager()

	gw.initialized = false
	log.Println("Game world shut down")

	return nil
}

// ============================================================================
// Entity Creation Helpers
// ============================================================================

// CreateBlock creates a block entity
func (gw *GameWorld) CreateBlock(blockType string, x, y, z float64) (*Entity, error) {
	if !gw.initialized {
		return nil, fmt.Errorf("game world not initialized")
	}

	entityID := fmt.Sprintf("block_%s_%.0f_%.0f_%.0f", blockType, x, y, z)
	entity, err := gw.entityManager.CreateEntityFromTemplate(blockType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to create block %s: %v", blockType, err)
	}

	// Add position metadata
	entity.Metadata["position"] = map[string]float64{
		"x": x, "y": y, "z": z,
	}

	// Register entity
	gw.systemManager.AddEntity(entity)

	// Publish event
	gw.eventBus.PublishWithSource(EventBlockPlaced, "world",
		CreateBlockEvent(blockType, x, y, z, "", ""))

	return entity, nil
}

// CreateItem creates an item entity
func (gw *GameWorld) CreateItem(itemType string, quantity int, playerID string) (*Entity, error) {
	if !gw.initialized {
		return nil, fmt.Errorf("game world not initialized")
	}

	entityID := fmt.Sprintf("item_%s_%s", itemType, playerID)
	entity, err := gw.entityManager.CreateEntityFromTemplate(itemType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to create item %s: %v", itemType, err)
	}

	// Set quantity in inventory component
	if inventoryComp, has := entity.GetComponent("inventory"); has {
		if inventory, ok := inventoryComp.(*InventoryComponent); ok {
			if inventory.Contents == nil {
				inventory.Contents = make(map[string]int)
			}
			inventory.Contents[itemType] = quantity
		}
	}

	// Add owner metadata
	entity.Metadata["owner"] = playerID

	// Register entity
	gw.systemManager.AddEntity(entity)

	// Publish event
	gw.eventBus.PublishWithSource(EventItemCrafted, "world",
		CreateItemEvent(itemType, quantity, playerID, "", true))

	return entity, nil
}

// CreateOrganism creates an organism entity
func (gw *GameWorld) CreateOrganism(organismType string, x, y, z float64) (*Entity, error) {
	if !gw.initialized {
		return nil, fmt.Errorf("game world not initialized")
	}

	entityID := fmt.Sprintf("organism_%s_%.0f_%.0f_%.0f", organismType, x, y, z)
	entity, err := gw.entityManager.CreateEntityFromTemplate(organismType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to create organism %s: %v", organismType, err)
	}

	// Add position metadata
	entity.Metadata["position"] = map[string]float64{
		"x": x, "y": y, "z": z,
	}

	// Register entity
	gw.systemManager.AddEntity(entity)

	// Publish event
	gw.eventBus.PublishWithSource(EventEntityAdded, "world",
		CreateEntityEvent(entity, "", nil, nil))

	return entity, nil
}

// ============================================================================
// Query Methods
// ============================================================================

// GetEntitiesByType returns all entities of a specific type
func (gw *GameWorld) GetEntitiesByType(entityType string) []*Entity {
	if !gw.initialized {
		return nil
	}
	return gw.systemManager.FindEntitiesByType(entityType)
}

// GetEntitiesByTag returns all entities with a specific tag
func (gw *GameWorld) GetEntitiesByTag(tag string) []*Entity {
	if !gw.initialized {
		return nil
	}
	return gw.systemManager.FindEntitiesByTag(tag)
}

// GetEntitiesByComponent returns all entities that have a specific component
func (gw *GameWorld) GetEntitiesByComponent(componentType string) []*Entity {
	if !gw.initialized {
		return nil
	}
	return gw.systemManager.FindEntitiesByComponent(componentType)
}

// GetEntity returns an entity by ID
func (gw *GameWorld) GetEntity(entityID string) (*Entity, bool) {
	if !gw.initialized {
		return nil, false
	}
	return gw.systemManager.GetEntity(entityID)
}

// ============================================================================
// Legacy Integration Methods
// ============================================================================

// These methods provide compatibility with the existing game code

// GetBlockColor returns the color for a block type (legacy compatibility)
func (gw *GameWorld) GetBlockColor(blockType string) (interface{}, error) {
	if !gw.initialized {
		return nil, fmt.Errorf("game world not initialized")
	}

	entities := gw.GetEntitiesByType(blockType)
	if len(entities) == 0 {
		// Try to create a temporary entity to get the color
		entity, err := gw.entityManager.CreateEntityFromTemplate(blockType, "temp")
		if err != nil {
			return nil, fmt.Errorf("block type %s not found", blockType)
		}

		renderComp, has := entity.GetComponent("render")
		if !has {
			return nil, fmt.Errorf("block %s has no render component", blockType)
		}

		render := renderComp.(*RenderComponent)
		return render.Color, nil
	}

	// Get color from first entity
	entity := entities[0]
	renderComp, has := entity.GetComponent("render")
	if !has {
		return nil, fmt.Errorf("block %s has no render component", blockType)
	}

	render := renderComp.(*RenderComponent)
	return render.Color, nil
}

// GetItemProperties returns properties for an item type (legacy compatibility)
func (gw *GameWorld) GetItemProperties(itemType string) (interface{}, error) {
	if !gw.initialized {
		return nil, fmt.Errorf("game world not initialized")
	}

	entities := gw.GetEntitiesByType(itemType)
	if len(entities) == 0 {
		// Try to create a temporary entity to get properties
		entity, err := gw.entityManager.CreateEntityFromTemplate(itemType, "temp")
		if err != nil {
			return nil, fmt.Errorf("item type %s not found", itemType)
		}

		// Return all components as properties
		properties := make(map[string]interface{})
		for compType, component := range entity.Components {
			properties[compType] = component
		}
		return properties, nil
	}

	// Get properties from first entity
	entity := entities[0]
	properties := make(map[string]interface{})
	for compType, component := range entity.Components {
		properties[compType] = component
	}
	return properties, nil
}

// GetOrganismProperties returns properties for an organism type (legacy compatibility)
func (gw *GameWorld) GetOrganismProperties(organismType string) (interface{}, error) {
	if !gw.initialized {
		return nil, fmt.Errorf("game world not initialized")
	}

	entities := gw.GetEntitiesByType(organismType)
	if len(entities) == 0 {
		// Try to create a temporary entity to get properties
		entity, err := gw.entityManager.CreateEntityFromTemplate(organismType, "temp")
		if err != nil {
			return nil, fmt.Errorf("organism type %s not found", organismType)
		}

		// Return all components as properties
		properties := make(map[string]interface{})
		for compType, component := range entity.Components {
			properties[compType] = component
		}
		return properties, nil
	}

	// Get properties from first entity
	entity := entities[0]
	properties := make(map[string]interface{})
	for compType, component := range entity.Components {
		properties[compType] = component
	}
	return properties, nil
}

// ============================================================================
// Utility Methods
// ============================================================================

// GetEntityManager returns the entity manager
func (gw *GameWorld) GetEntityManager() *EntityManager {
	return gw.entityManager
}

// GetSystemManager returns the system manager
func (gw *GameWorld) GetSystemManager() *SystemManager {
	return gw.systemManager
}

// GetPluginManager returns the plugin manager
func (gw *GameWorld) GetPluginManager() *PluginManager {
	return gw.pluginManager
}

// GetEventBus returns the event bus
func (gw *GameWorld) GetEventBus() *EventBus {
	return gw.eventBus
}

// IsInitialized returns whether the world is initialized
func (gw *GameWorld) IsInitialized() bool {
	return gw.initialized
}

// GetStatistics returns world statistics
func (gw *GameWorld) GetStatistics() map[string]interface{} {
	if !gw.initialized {
		return map[string]interface{}{
			"initialized": false,
		}
	}

	entities := gw.systemManager.GetEntities()
	typeCounts := make(map[string]int)
	componentCounts := make(map[string]int)
	tagCounts := make(map[string]int)

	for _, entity := range entities {
		// Count by type
		typeCounts[entity.Type]++

		// Count by components
		for compType := range entity.Components {
			componentCounts[compType]++
		}

		// Count by tags
		for _, tag := range entity.Tags {
			tagCounts[tag]++
		}
	}

	return map[string]interface{}{
		"initialized":      true,
		"total_entities":   len(entities),
		"entity_types":     typeCounts,
		"component_counts": componentCounts,
		"tag_counts":       tagCounts,
		"loaded_plugins":   gw.pluginManager.ListPlugins(),
		"loaded_templates": gw.entityManager.ListTemplates(),
	}
}

// DebugPrint prints debug information about the world
func (gw *GameWorld) DebugPrint() {
	if !gw.initialized {
		log.Println("Game world not initialized")
		return
	}

	stats := gw.GetStatistics()
	log.Printf("=== Game World Statistics ===")
	log.Printf("Initialized: %v", stats["initialized"])
	log.Printf("Total Entities: %d", stats["total_entities"])
	log.Printf("Entity Types: %v", stats["entity_types"])
	log.Printf("Component Counts: %v", stats["component_counts"])
	log.Printf("Tag Counts: %v", stats["tag_counts"])
	log.Printf("Loaded Plugins: %v", stats["loaded_plugins"])
	log.Printf("Loaded Templates: %v", stats["loaded_templates"])
	log.Printf("============================")
}
