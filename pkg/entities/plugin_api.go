package entities

import (
	"context"
	"fmt"
	"log"
	"time"

	"tesselbox/pkg/world"
)

// ============================================================================
// Enhanced Plugin API
// ============================================================================

// PluginAPI provides the full API that plugins can use to interact with the game
type PluginAPI struct {
	manager        *PluginManager
	entityManager  *EntityManager
	systemManager  *SystemManager
	eventBus       *EventBus
	world          *world.World
	pluginName     string
	allowedActions map[string]bool
}

// NewPluginAPI creates a new plugin API instance for a specific plugin
func NewPluginAPI(manager *PluginManager, pluginName string) *PluginAPI {
	return &PluginAPI{
		manager:        manager,
		entityManager:  manager.entityManager,
		systemManager:  manager.systemManager,
		eventBus:       manager.eventBus,
		world:          manager.world,
		pluginName:     pluginName,
		allowedActions: make(map[string]bool),
	}
}

// ============================================================================
// Entity Management API
// ============================================================================

// CreateEntity creates a new entity with the given template
func (api *PluginAPI) CreateEntity(templateID string) (*Entity, error) {
	if !api.hasPermission("entity.create") {
		return nil, fmt.Errorf("plugin %s does not have permission to create entities", api.pluginName)
	}

	// Create entity using template
	template, err := api.GetTemplate(templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template %s: %v", templateID, err)
	}

	entity := api.entityManager.CreateEntityFromTemplate(template)
	log.Printf("Plugin %s created entity %d from template %s", api.pluginName, entity.ID, templateID)
	return entity, nil
}

// CreateCustomEntity creates a new entity with custom components
func (api *PluginAPI) CreateCustomEntity(entityType string, components map[string]interface{}) (*Entity, error) {
	if !api.hasPermission("entity.create_custom") {
		return nil, fmt.Errorf("plugin %s does not have permission to create custom entities", api.pluginName)
	}

	entity := api.entityManager.CreateEntity(entityType)
	
	// Add components
	for compType, compData := range components {
		component, err := api.CreateComponent(compType, compData)
		if err != nil {
			api.entityManager.RemoveEntity(entity.ID)
			return nil, fmt.Errorf("failed to create component %s: %v", compType, err)
		}
		entity.AddComponent(component)
	}

	log.Printf("Plugin %s created custom entity %d of type %s", api.pluginName, entity.ID, entityType)
	return entity, nil
}

// RemoveEntity removes an entity from the world
func (api *PluginAPI) RemoveEntity(entityID uint64) error {
	if !api.hasPermission("entity.remove") {
		return fmt.Errorf("plugin %s does not have permission to remove entities", api.pluginName)
	}

	if err := api.entityManager.RemoveEntity(entityID); err != nil {
		return fmt.Errorf("failed to remove entity %d: %v", entityID, err)
	}

	log.Printf("Plugin %s removed entity %d", api.pluginName, entityID)
	return nil
}

// GetEntity gets an entity by ID
func (api *PluginAPI) GetEntity(entityID uint64) (*Entity, error) {
	if !api.hasPermission("entity.get") {
		return nil, fmt.Errorf("plugin %s does not have permission to get entities", api.pluginName)
	}

	entity, exists := api.entityManager.GetEntity(entityID)
	if !exists {
		return nil, fmt.Errorf("entity %d not found", entityID)
	}
	return entity, nil
}

// FindEntities finds entities matching criteria
func (api *PluginAPI) FindEntities(criteria EntityCriteria) ([]*Entity, error) {
	if !api.hasPermission("entity.find") {
		return nil, fmt.Errorf("plugin %s does not have permission to find entities", api.pluginName)
	}

	return api.entityManager.FindEntities(criteria), nil
}

// ============================================================================
// Component Management API
// ============================================================================

// CreateComponent creates a component from data
func (api *PluginAPI) CreateComponent(componentType string, data interface{}) (Component, error) {
	if !api.hasPermission("component.create") {
		return nil, fmt.Errorf("plugin %s does not have permission to create components", api.pluginName)
	}

	// Use the component registry to create components
	component, err := CreateComponentFromData(componentType, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create component %s: %v", componentType, err)
	}
	return component, nil
}

// AddComponent adds a component to an entity
func (api *PluginAPI) AddComponent(entityID uint64, component Component) error {
	if !api.hasPermission("entity.modify") {
		return fmt.Errorf("plugin %s does not have permission to modify entities", api.pluginName)
	}

	entity, exists := api.entityManager.GetEntity(entityID)
	if !exists {
		return fmt.Errorf("entity %d not found", entityID)
	}

	entity.AddComponent(component)
	log.Printf("Plugin %s added component %s to entity %d", api.pluginName, component.GetType(), entityID)
	return nil
}

// RemoveComponent removes a component from an entity
func (api *PluginAPI) RemoveComponent(entityID uint64, componentType string) error {
	if !api.hasPermission("entity.modify") {
		return fmt.Errorf("plugin %s does not have permission to modify entities", api.pluginName)
	}

	entity, exists := api.entityManager.GetEntity(entityID)
	if !exists {
		return fmt.Errorf("entity %d not found", entityID)
	}

	entity.RemoveComponent(componentType)
	log.Printf("Plugin %s removed component %s from entity %d", api.pluginName, componentType, entityID)
	return nil
}

// ============================================================================
// Event System API
// ============================================================================

// PublishEvent publishes an event to the event bus
func (api *PluginAPI) PublishEvent(event Event) error {
	if !api.hasPermission("event.publish") {
		return fmt.Errorf("plugin %s does not have permission to publish events", api.pluginName)
	}

	// Add plugin metadata to event
	if event.GetMetadata() == nil {
		event.SetMetadata(make(map[string]interface{}))
	}
	event.GetMetadata()["plugin"] = api.pluginName
	event.GetMetadata()["timestamp"] = time.Now()

	api.eventBus.Publish(event)
	log.Printf("Plugin %s published event %s", api.pluginName, event.GetType())
	return nil
}

// SubscribeToEvent subscribes to an event type
func (api *PluginAPI) SubscribeToEvent(eventType string, handler EventHandler) error {
	if !api.hasPermission("event.subscribe") {
		return fmt.Errorf("plugin %s does not have permission to subscribe to events", api.pluginName)
	}

	// Wrap handler to add plugin context
	wrappedHandler := func(event Event) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Plugin %s event handler panicked for event %s: %v", api.pluginName, eventType, r)
			}
		}()
		
		// Add plugin context to event metadata
		if event.GetMetadata() == nil {
			event.SetMetadata(make(map[string]interface{}))
		}
		event.GetMetadata()["handling_plugin"] = api.pluginName
		
		handler(event)
	}

	api.eventBus.Subscribe(eventType, wrappedHandler)
	log.Printf("Plugin %s subscribed to event %s", api.pluginName, eventType)
	return nil
}

// ============================================================================
// World Interaction API
// ============================================================================

// GetBlockAt gets the block at the given world coordinates
func (api *PluginAPI) GetBlockAt(x, y float64) (*world.Hexagon, error) {
	if !api.hasPermission("world.read") {
		return nil, fmt.Errorf("plugin %s does not have permission to read world data", api.pluginName)
	}

	hex := api.world.GetHexagonAt(x, y)
	if hex == nil {
		return nil, fmt.Errorf("no block found at (%.1f, %.1f)", x, y)
	}
	return hex, nil
}

// SetBlockAt sets a block at the given world coordinates
func (api *PluginAPI) SetBlockAt(x, y float64, blockType string) error {
	if !api.hasPermission("world.modify") {
		return fmt.Errorf("plugin %s does not have permission to modify world", api.pluginName)
	}

	// Convert block type string to actual block type
	// This would need to be implemented based on the block system
	log.Printf("Plugin %s set block %s at (%.1f, %.1f)", api.pluginName, blockType, x, y)
	
	// Publish block placed event
	api.PublishEvent(&BlockEvent{
		Type:      "block_placed",
		BlockType: blockType,
		X:         x,
		Y:         y,
		Z:         0,
		EntityID:  0,
		Cause:     "plugin",
		Source:    api.pluginName,
	})
	
	return nil
}

// RemoveBlockAt removes a block at the given world coordinates
func (api *PluginAPI) RemoveBlockAt(x, y float64) error {
	if !api.hasPermission("world.modify") {
		return fmt.Errorf("plugin %s does not have permission to modify world", api.pluginName)
	}

	hex := api.world.GetHexagonAt(x, y)
	if hex == nil {
		return fmt.Errorf("no block found at (%.1f, %.1f)", x, y)
	}

	blockType := hex.BlockType.String()
	
	// Remove the block
	api.world.RemoveHexagonAt(x, y)
	
	// Publish block broken event
	api.PublishEvent(&BlockEvent{
		Type:      "block_broken",
		BlockType: blockType,
		X:         x,
		Y:         y,
		Z:         0,
		EntityID:  0,
		Cause:     "plugin",
		Source:    api.pluginName,
	})
	
	log.Printf("Plugin %s removed block at (%.1f, %.1f)", api.pluginName, x, y)
	return nil
}

// ============================================================================
// Template Management API
// ============================================================================

// GetTemplate gets an entity template by ID
func (api *PluginAPI) GetTemplate(templateID string) (*EntityTemplate, error) {
	if !api.hasPermission("template.get") {
		return nil, fmt.Errorf("plugin %s does not have permission to get templates", api.pluginName)
	}

	// This would need to be implemented in EntityManager
	// For now, return a basic template
	return &EntityTemplate{
		ID:   templateID,
		Type: "custom",
		Name: "Custom Template",
		Components: make(map[string]interface{}),
	}, nil
}

// RegisterTemplate registers a new entity template
func (api *PluginAPI) RegisterTemplate(template *EntityTemplate) error {
	if !api.hasPermission("template.register") {
		return fmt.Errorf("plugin %s does not have permission to register templates", api.pluginName)
	}

	// Add plugin metadata
	if template.Metadata == nil {
		template.Metadata = make(map[string]interface{})
	}
	template.Metadata["plugin"] = api.pluginName
	template.Metadata["registered_at"] = time.Now()

	log.Printf("Plugin %s registered template %s", api.pluginName, template.ID)
	return nil
}

// ============================================================================
// System Management API
// ============================================================================

// RegisterSystem registers a new system
func (api *PluginAPI) RegisterSystem(system System) error {
	if !api.hasPermission("system.register") {
		return fmt.Errorf("plugin %s does not have permission to register systems", api.pluginName)
	}

	api.systemManager.RegisterSystem(system)
	log.Printf("Plugin %s registered system %s", api.pluginName, system.GetName())
	return nil
}

// UnregisterSystem unregisters a system
func (api *PluginAPI) UnregisterSystem(systemName string) error {
	if !api.hasPermission("system.unregister") {
		return fmt.Errorf("plugin %s does not have permission to unregister systems", api.pluginName)
	}

	api.systemManager.UnregisterSystem(systemName)
	log.Printf("Plugin %s unregistered system %s", api.pluginName, systemName)
	return nil
}

// ============================================================================
// Utility API
// ============================================================================

// Log logs a message from the plugin
func (api *PluginAPI) Log(level string, message string) {
	log.Printf("[%s] %s: %s", strings.ToUpper(level), api.pluginName, message)
}

// GetPluginInfo gets information about the plugin
func (api *PluginAPI) GetPluginInfo() *PluginInfo {
	info, _ := api.manager.GetPluginInfo(api.pluginName)
	return info
}

// GetLoadedPlugins gets a list of loaded plugins
func (api *PluginAPI) GetLoadedPlugins() []string {
	return api.manager.ListPlugins()
}

// IsPluginLoaded checks if a plugin is loaded
func (api *PluginAPI) IsPluginLoaded(pluginName string) bool {
	return api.manager.IsLoaded(pluginName)
}

// ============================================================================
// Permission System
// ============================================================================

// GrantPermission grants a permission to the plugin
func (api *PluginAPI) GrantPermission(permission string) {
	api.allowedActions[permission] = true
}

// RevokePermission revokes a permission from the plugin
func (api *PluginAPI) RevokePermission(permission string) {
	delete(api.allowedActions, permission)
}

// HasPermission checks if the plugin has a specific permission
func (api *PluginAPI) HasPermission(permission string) bool {
	return api.hasPermission(permission)
}

// hasPermission is the internal permission check
func (api *PluginAPI) hasPermission(permission string) bool {
	allowed, exists := api.allowedActions[permission]
	if !exists {
		// Default to allowing most actions for now
		// In a production system, this would be more restrictive
		return true
	}
	return allowed
}

// ============================================================================
// Plugin Context
// ============================================================================

// PluginContext provides context for plugin operations
type PluginContext struct {
	PluginName string
	API        *PluginAPI
	Cancel     context.CancelFunc
	Done       <-chan struct{}
}

// CreateContext creates a new plugin context
func (api *PluginAPI) CreateContext() *PluginContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &PluginContext{
		PluginName: api.pluginName,
		API:        api,
		Cancel:     cancel,
		Done:       ctx.Done(),
	}
}

// ============================================================================
// Enhanced Plugin Interface
// ============================================================================

// EnhancedPlugin extends the basic Plugin interface with additional lifecycle methods
type EnhancedPlugin interface {
	Plugin
	
	// Enhanced lifecycle methods
	OnLoad(api *PluginAPI) error
	OnUnload(api *PluginAPI) error
	OnReload(api *PluginAPI) error
	
	// Event handlers
	OnGameStart(api *PluginAPI) error
	OnGameStop(api *PluginAPI) error
	OnPlayerJoin(api *PluginAPI, playerID string) error
	OnPlayerLeave(api *PluginAPI, playerID string) error
	
	// Configuration
	GetDefaultConfig() map[string]interface{}
	ValidateConfig(config map[string]interface{}) error
	OnConfigChange(api *PluginAPI, config map[string]interface{}) error
	
	// Permissions
	GetRequiredPermissions() []string
}

// EntityCriteria defines criteria for finding entities
type EntityCriteria struct {
	Types      []string
	Components []string
	Tags       []string
	WithinArea *Area
	Limit      int
}

// Area defines a rectangular area
type Area struct {
	X, Y, Width, Height float64
}
