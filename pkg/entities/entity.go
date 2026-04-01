package entities

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Entity represents any game object using a component-based architecture
type Entity struct {
	ID         string                 `yaml:"id" json:"id"`
	Type       string                 `yaml:"type" json:"type"`
	Name       string                 `yaml:"name" json:"name"`
	Components map[string]Component   `yaml:"-" json:"-"`
	Tags       []string               `yaml:"tags" json:"tags"`
	Metadata   map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// Component defines the interface for all entity components
type Component interface {
	GetType() string
	Clone() Component
	Merge(other Component)
	Validate() error
}

// ComponentRegistry holds all registered component types
var ComponentRegistry = make(map[string]reflect.Type)

// RegisterComponent registers a component type for dynamic loading
func RegisterComponent(name string, component Component) {
	ComponentRegistry[name] = reflect.TypeOf(component).Elem()
}

// CreateComponent creates a component instance by name
func CreateComponent(componentType string) (Component, error) {
	if compType, ok := ComponentRegistry[componentType]; ok {
		component := reflect.New(compType).Interface().(Component)
		return component, nil
	}
	return nil, fmt.Errorf("unknown component type: %s", componentType)
}

// NewEntity creates a new entity with the specified ID and type
func NewEntity(id, entityType string) *Entity {
	return &Entity{
		ID:         id,
		Type:       entityType,
		Components: make(map[string]Component),
		Tags:       make([]string, 0),
		Metadata:   make(map[string]interface{}),
	}
}

// AddComponent adds a component to the entity
func (e *Entity) AddComponent(component Component) {
	e.Components[component.GetType()] = component
}

// GetComponent gets a component by type
func (e *Entity) GetComponent(componentType string) (Component, bool) {
	component, exists := e.Components[componentType]
	return component, exists
}

// RemoveComponent removes a component from the entity
func (e *Entity) RemoveComponent(componentType string) {
	delete(e.Components, componentType)
}

// HasComponent checks if the entity has a specific component
func (e *Entity) HasComponent(componentType string) bool {
	_, exists := e.Components[componentType]
	return exists
}

// HasTag checks if the entity has a specific tag
func (e *Entity) HasTag(tag string) bool {
	for _, t := range e.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// AddTag adds a tag to the entity
func (e *Entity) AddTag(tag string) {
	if !e.HasTag(tag) {
		e.Tags = append(e.Tags, tag)
	}
}

// RemoveTag removes a tag from the entity
func (e *Entity) RemoveTag(tag string) {
	for i, t := range e.Tags {
		if t == tag {
			e.Tags = append(e.Tags[:i], e.Tags[i+1:]...)
			break
		}
	}
}

// Clone creates a deep copy of the entity
func (e *Entity) Clone() *Entity {
	clone := &Entity{
		ID:         e.ID + "_clone",
		Type:       e.Type,
		Name:       e.Name,
		Components: make(map[string]Component),
		Tags:       make([]string, len(e.Tags)),
		Metadata:   make(map[string]interface{}),
	}

	// Copy tags
	copy(clone.Tags, e.Tags)

	// Copy components
	for compType, component := range e.Components {
		clone.Components[compType] = component.Clone()
	}

	// Copy metadata
	for k, v := range e.Metadata {
		clone.Metadata[k] = v
	}

	return clone
}

// ToJSON converts the entity to JSON representation
func (e *Entity) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON creates an entity from JSON representation
func FromJSON(data []byte) (*Entity, error) {
	var entity Entity
	err := json.Unmarshal(data, &entity)
	if err != nil {
		return nil, err
	}

	// Reconstruct components from JSON
	entity.Components = make(map[string]Component)
	// Component reconstruction would need to be implemented based on saved data

	return &entity, nil
}

// EntityTemplate represents a template for creating entities
type EntityTemplate struct {
	ID         string                 `yaml:"id"`
	Type       string                 `yaml:"type"`
	Name       string                 `yaml:"name"`
	Components map[string]interface{} `yaml:"components"`
	Tags       []string               `yaml:"tags"`
	Metadata   map[string]interface{} `yaml:"metadata,omitempty"`
	Inherits   []string               `yaml:"inherits,omitempty"`
}

// EntityManager manages entity creation and templates
type EntityManager struct {
	templates map[string]*EntityTemplate
	entities  map[string]*Entity
}

// NewEntityManager creates a new entity manager
func NewEntityManager() *EntityManager {
	return &EntityManager{
		templates: make(map[string]*EntityTemplate),
		entities:  make(map[string]*Entity),
	}
}

// LoadTemplates loads entity templates from YAML data
func (em *EntityManager) LoadTemplates(data []byte) error {
	var templates map[string]*EntityTemplate
	err := yaml.Unmarshal(data, &templates)
	if err != nil {
		return err
	}

	for id, template := range templates {
		template.ID = id
		em.templates[id] = template
	}

	return nil
}

// CreateEntityFromTemplate creates an entity from a template
func (em *EntityManager) CreateEntityFromTemplate(templateID string, entityID string) (*Entity, error) {
	template, exists := em.templates[templateID]
	if !exists {
		return nil, fmt.Errorf("template not found: %s", templateID)
	}

	entity := NewEntity(entityID, template.Type)
	entity.Name = template.Name
	entity.Tags = make([]string, len(template.Tags))
	copy(entity.Tags, template.Tags)

	// Copy metadata
	for k, v := range template.Metadata {
		entity.Metadata[k] = v
	}

	// Process inheritance
	for _, parentID := range template.Inherits {
		parent, err := em.CreateEntityFromTemplate(parentID, entityID+"_parent")
		if err != nil {
			return nil, fmt.Errorf("failed to process inheritance from %s: %v", parentID, err)
		}

		// Merge parent components
		for compType, component := range parent.Components {
			if !entity.HasComponent(compType) {
				entity.AddComponent(component.Clone())
			}
		}

		// Merge parent tags
		for _, tag := range parent.Tags {
			entity.AddTag(tag)
		}
	}

	// Create components from template
	for compType, compData := range template.Components {
		component, err := em.createComponentFromData(compType, compData)
		if err != nil {
			return nil, fmt.Errorf("failed to create component %s: %v", compType, err)
		}

		// Merge with existing component if present
		if existing, has := entity.GetComponent(compType); has {
			existing.Merge(component)
		} else {
			entity.AddComponent(component)
		}
	}

	return entity, nil
}

// createComponentFromData creates a component from raw data
func (em *EntityManager) createComponentFromData(compType string, data interface{}) (Component, error) {
	component, err := CreateComponent(compType)
	if err != nil {
		return nil, err
	}

	// Convert data to YAML and back to unmarshal into component
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlData, component)
	if err != nil {
		return nil, err
	}

	err = component.Validate()
	if err != nil {
		return nil, err
	}

	return component, nil
}

// RegisterEntity registers an entity with the manager
func (em *EntityManager) RegisterEntity(entity *Entity) {
	em.entities[entity.ID] = entity
}

// GetEntity retrieves an entity by ID
func (em *EntityManager) GetEntity(id string) (*Entity, bool) {
	entity, exists := em.entities[id]
	return entity, exists
}

// RemoveEntity removes an entity from the manager
func (em *EntityManager) RemoveEntity(id string) {
	delete(em.entities, id)
}

// FindEntitiesByTag finds all entities with a specific tag
func (em *EntityManager) FindEntitiesByTag(tag string) []*Entity {
	var results []*Entity
	for _, entity := range em.entities {
		if entity.HasTag(tag) {
			results = append(results, entity)
		}
	}
	return results
}

// FindEntitiesByType finds all entities of a specific type
func (em *EntityManager) FindEntitiesByType(entityType string) []*Entity {
	var results []*Entity
	for _, entity := range em.entities {
		if entity.Type == entityType {
			results = append(results, entity)
		}
	}
	return results
}

// FindEntitiesByComponent finds all entities that have a specific component
func (em *EntityManager) FindEntitiesByComponent(componentType string) []*Entity {
	var results []*Entity
	for _, entity := range em.entities {
		if entity.HasComponent(componentType) {
			results = append(results, entity)
		}
	}
	return results
}

// GetTemplate returns a template by ID
func (em *EntityManager) GetTemplate(id string) (*EntityTemplate, bool) {
	template, exists := em.templates[id]
	return template, exists
}

// ListTemplates returns all template IDs
func (em *EntityManager) ListTemplates() []string {
	templates := make([]string, 0, len(em.templates))
	for id := range em.templates {
		templates = append(templates, id)
	}
	return templates
}

// GenerateUniqueID generates a unique entity ID
func GenerateUniqueID(prefix string) string {
	// Simple ID generation - in production use UUID or proper ID generator
	return prefix + "_" + strconv.FormatInt(int64(len(prefix)*1000), 10)
}

// ValidateEntityID checks if an entity ID is valid
func ValidateEntityID(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("entity ID cannot be empty")
	}
	if strings.Contains(id, " ") {
		return fmt.Errorf("entity ID cannot contain spaces")
	}
	return nil
}
