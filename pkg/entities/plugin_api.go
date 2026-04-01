package entities

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
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
	pluginName     string
	allowedActions map[string]bool
}

// NewPluginAPI creates a new plugin API instance for a specific plugin
func NewPluginAPI(manager *PluginManager, pluginName string) *PluginAPI {
	api := &PluginAPI{
		manager:        manager,
		entityManager:  manager.entityManager,
		systemManager:  manager.systemManager,
		eventBus:       manager.eventBus,
		pluginName:     pluginName,
		allowedActions: make(map[string]bool),
	}

	// Initialize with basic safe permissions
	api.initializeDefaultPermissions()

	return api
}

// initializeDefaultPermissions sets up safe default permissions for plugins
func (api *PluginAPI) initializeDefaultPermissions() {
	// Basic permissions that are generally safe
	safePermissions := []string{
		"entity.get",
		"entity.find",
		"event.subscribe",
		"template.get",
		"system.get", // Read-only system access
	}

	for _, perm := range safePermissions {
		api.allowedActions[perm] = true
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

	// Check resource limits (this would need context parameter in real implementation)
	// For now, just log the operation
	log.Printf("Plugin %s creating entity from template %s", api.pluginName, templateID)

	// Create entity using template
	entityID := fmt.Sprintf("plugin_%s_%s", api.pluginName, templateID)
	entity, err := api.entityManager.CreateEntityFromTemplate(templateID, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to create entity from template %s: %v", templateID, err)
	}

	log.Printf("Plugin %s created entity %s from template %s", api.pluginName, entity.ID, templateID)
	return entity, nil
}

// CreateCustomEntity creates a new entity with custom components
func (api *PluginAPI) CreateCustomEntity(entityType string, components map[string]interface{}) (*Entity, error) {
	if !api.hasPermission("entity.create_custom") {
		return nil, fmt.Errorf("plugin %s does not have permission to create custom entities", api.pluginName)
	}

	entityID := fmt.Sprintf("plugin_%s_custom_%s", api.pluginName, entityType)
	entity := &Entity{
		ID:         entityID,
		Type:       entityType,
		Components: make(map[string]Component),
		Tags:       []string{},
	}

	// Add components
	for compType, compData := range components {
		component, err := api.entityManager.createComponentFromData(compType, compData)
		if err != nil {
			return nil, fmt.Errorf("failed to create component %s: %v", compType, err)
		}
		entity.AddComponent(component)
	}

	// Register entity
	api.entityManager.RegisterEntity(entity)
	log.Printf("Plugin %s created custom entity %s of type %s", api.pluginName, entity.ID, entityType)
	return entity, nil
}

// RemoveEntity removes an entity from the world
func (api *PluginAPI) RemoveEntity(entityID uint64) error {
	if !api.hasPermission("entity.remove") {
		return fmt.Errorf("plugin %s does not have permission to remove entities", api.pluginName)
	}

	idStr := fmt.Sprintf("%d", entityID)
	api.entityManager.RemoveEntity(idStr)
	log.Printf("Plugin %s removed entity %s", api.pluginName, idStr)
	return nil
}

// GetEntity gets an entity by ID
func (api *PluginAPI) GetEntity(entityID uint64) (*Entity, error) {
	if !api.hasPermission("entity.get") {
		return nil, fmt.Errorf("plugin %s does not have permission to get entities", api.pluginName)
	}

	idStr := fmt.Sprintf("%d", entityID)
	entity, exists := api.entityManager.GetEntity(idStr)
	if !exists {
		return nil, fmt.Errorf("entity %s not found", idStr)
	}
	return entity, nil
}

// FindEntities finds entities matching criteria
func (api *PluginAPI) FindEntities(criteria EntityCriteria) ([]*Entity, error) {
	if !api.hasPermission("entity.find") {
		return nil, fmt.Errorf("plugin %s does not have permission to find entities", api.pluginName)
	}

	// Simple implementation - return all entities for now
	// Use FindEntitiesByType with empty type to get all
	results := api.entityManager.FindEntitiesByType("")
	return results, nil
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
	component, err := api.entityManager.createComponentFromData(componentType, data)
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

	idStr := fmt.Sprintf("%d", entityID)
	entity, exists := api.entityManager.GetEntity(idStr)
	if !exists {
		return fmt.Errorf("entity %s not found", idStr)
	}

	entity.AddComponent(component)
	log.Printf("Plugin %s added component %s to entity %s", api.pluginName, component.GetType(), idStr)
	return nil
}

// RemoveComponent removes a component from an entity
func (api *PluginAPI) RemoveComponent(entityID uint64, componentType string) error {
	if !api.hasPermission("entity.modify") {
		return fmt.Errorf("plugin %s does not have permission to modify entities", api.pluginName)
	}

	idStr := fmt.Sprintf("%d", entityID)
	entity, exists := api.entityManager.GetEntity(idStr)
	if !exists {
		return fmt.Errorf("entity %s not found", idStr)
	}

	entity.RemoveComponent(componentType)
	log.Printf("Plugin %s removed component %s from entity %s", api.pluginName, componentType, idStr)
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
	event.Source = api.pluginName
	event.Timestamp = time.Now()

	api.eventBus.Publish(event.Type, event)
	log.Printf("Plugin %s published event %s", api.pluginName, event.Type)
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

		// Add plugin context to event data
		if event.Data == nil {
			event.Data = make(map[string]interface{})
		}
		if dataMap, ok := event.Data.(map[string]interface{}); ok {
			dataMap["handling_plugin"] = api.pluginName
		}

		handler(event)
	}

	// Convert string to EventType
	eventTypeEnum := EventType(eventType)
	api.eventBus.Subscribe(eventTypeEnum, wrappedHandler)
	log.Printf("Plugin %s subscribed to event %s", api.pluginName, eventType)
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

	template, exists := api.entityManager.GetTemplate(templateID)
	if !exists {
		return nil, fmt.Errorf("template %s not found", templateID)
	}
	return template, nil
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
		// Default to deny - strict security model
		// Plugins must explicitly request permissions
		return false
	}
	return allowed
}

// ============================================================================
// Plugin Context
// ============================================================================

// PluginContext provides context for plugin operations
type PluginContext struct {
	PluginName     string
	API            *PluginAPI
	Cancel         context.CancelFunc
	Done           <-chan struct{}
	Timeout        time.Duration
	MaxOperations  int
	OperationCount int64
}

// CreateContext creates a new plugin context
func (api *PluginAPI) CreateContext() *PluginContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &PluginContext{
		PluginName:     api.pluginName,
		API:            api,
		Cancel:         cancel,
		Done:           ctx.Done(),
		Timeout:        30 * time.Second, // Default 30 second timeout
		MaxOperations:  1000,             // Default max operations
		OperationCount: 0,
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
