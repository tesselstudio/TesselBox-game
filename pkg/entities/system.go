package entities

import (
	"log"
	"time"

)

// System represents a game system that processes entities
type System interface {
	GetName() string
	GetRequiredComponents() []string
	Process(deltaTime float64, entities []*Entity)
	Matches(entity *Entity) bool
}

// SystemManager manages all game systems
type SystemManager struct {
	systems map[string]System
	entities map[string]*Entity
	eventBus *EventBus
}

// NewSystemManager creates a new system manager
func NewSystemManager() *SystemManager {
	return &SystemManager{
		systems: make(map[string]System),
		entities: make(map[string]*Entity),
		eventBus: NewEventBus(),
	}
}

// RegisterSystem registers a system with the manager
func (sm *SystemManager) RegisterSystem(system System) {
	sm.systems[system.GetName()] = system
	log.Printf("Registered system: %s", system.GetName())
}

// UnregisterSystem unregisters a system
func (sm *SystemManager) UnregisterSystem(systemName string) {
	delete(sm.systems, systemName)
	log.Printf("Unregistered system: %s", systemName)
}

// AddEntity adds an entity to be processed by systems
func (sm *SystemManager) AddEntity(entity *Entity) {
	sm.entities[entity.ID] = entity
	sm.eventBus.Publish(EventEntityAdded, EntityEvent{Entity: entity})
}

// RemoveEntity removes an entity from processing
func (sm *SystemManager) RemoveEntity(entityID string) {
	if entity, exists := sm.entities[entityID]; exists {
		delete(sm.entities, entityID)
		sm.eventBus.Publish(EventEntityRemoved, EntityEvent{Entity: entity})
	}
}

// Update processes all systems
func (sm *SystemManager) Update(deltaTime float64) {
	// Convert entities to slice for processing
	entityList := make([]*Entity, 0, len(sm.entities))
	for _, entity := range sm.entities {
		entityList = append(entityList, entity)
	}

	// Process each system
	for _, system := range sm.systems {
		start := time.Now()
		system.Process(deltaTime, entityList)
		duration := time.Since(start)
		
		if duration > time.Millisecond*10 {
			log.Printf("System %s took %v to process", system.GetName(), duration)
		}
	}
}

// GetEntities returns all entities
func (sm *SystemManager) GetEntities() map[string]*Entity {
	return sm.entities
}

// GetEntity returns an entity by ID
func (sm *SystemManager) GetEntity(id string) (*Entity, bool) {
	entity, exists := sm.entities[id]
	return entity, exists
}

// FindEntitiesByComponent finds entities that have specific components
func (sm *SystemManager) FindEntitiesByComponent(componentTypes ...string) []*Entity {
	var results []*Entity
	for _, entity := range sm.entities {
		matches := true
		for _, compType := range componentTypes {
			if !entity.HasComponent(compType) {
				matches = false
				break
			}
		}
		if matches {
			results = append(results, entity)
		}
	}
	return results
}

// FindEntitiesByType finds all entities of a specific type
func (sm *SystemManager) FindEntitiesByType(entityType string) []*Entity {
	var results []*Entity
	for _, entity := range sm.entities {
		if entity.Type == entityType {
			results = append(results, entity)
		}
	}
	return results
}

// FindEntitiesByTag finds entities with specific tags
func (sm *SystemManager) FindEntitiesByTag(tags ...string) []*Entity {
	var results []*Entity
	for _, entity := range sm.entities {
		matches := true
		for _, tag := range tags {
			if !entity.HasTag(tag) {
				matches = false
				break
			}
		}
		if matches {
			results = append(results, entity)
		}
	}
	return results
}

// GetEventBus returns the event bus
func (sm *SystemManager) GetEventBus() *EventBus {
	return sm.eventBus
}

// ============================================================================
// Core Systems
// ============================================================================

// RenderSystem handles rendering of entities
type RenderSystem struct {
	name string
	requiredComponents []string
}

func NewRenderSystem() *RenderSystem {
	return &RenderSystem{
		name: "render",
		requiredComponents: []string{"render"},
	}
}

func (rs *RenderSystem) GetName() string { return rs.name }
func (rs *RenderSystem) GetRequiredComponents() []string { return rs.requiredComponents }
func (rs *RenderSystem) Matches(entity *Entity) bool {
	return entity.HasComponent("render")
}

func (rs *RenderSystem) Process(deltaTime float64, entities []*Entity) {
	for _, entity := range entities {
		if !rs.Matches(entity) {
			continue
		}

		renderComp, _ := entity.GetComponent("render")
		if render, ok := renderComp.(*RenderComponent); ok {
			// Update animation if needed
			if render.Animated {
				// Animation logic would go here
			}
		}
	}
}

// PhysicsSystem handles physics simulation
type PhysicsSystem struct {
	name string
	requiredComponents []string
	gravity float64
}

func NewPhysicsSystem(gravity float64) *PhysicsSystem {
	return &PhysicsSystem{
		name: "physics",
		requiredComponents: []string{"physics"},
		gravity: gravity,
	}
}

func (ps *PhysicsSystem) GetName() string { return ps.name }
func (ps *PhysicsSystem) GetRequiredComponents() []string { return ps.requiredComponents }
func (ps *PhysicsSystem) Matches(entity *Entity) bool {
	return entity.HasComponent("physics")
}

func (ps *PhysicsSystem) Process(deltaTime float64, entities []*Entity) {
	for _, entity := range entities {
		if !ps.Matches(entity) {
			continue
		}

		physicsComp, _ := entity.GetComponent("physics")
		if physics, ok := physicsComp.(*PhysicsComponent); ok {
			// Apply gravity if enabled
			if physics.Gravity {
				// Gravity logic would go here
			}

			// Handle collisions if enabled
			if physics.Collision {
				// Collision detection would go here
			}
		}
	}
}

// BehaviorSystem handles AI and entity behaviors
type BehaviorSystem struct {
	name string
	requiredComponents []string
}

func NewBehaviorSystem() *BehaviorSystem {
	return &BehaviorSystem{
		name: "behavior",
		requiredComponents: []string{"behavior"},
	}
}

func (bs *BehaviorSystem) GetName() string { return bs.name }
func (bs *BehaviorSystem) GetRequiredComponents() []string { return bs.requiredComponents }
func (bs *BehaviorSystem) Matches(entity *Entity) bool {
	return entity.HasComponent("behavior")
}

func (bs *BehaviorSystem) Process(deltaTime float64, entities []*Entity) {
	for _, entity := range entities {
		if !bs.Matches(entity) {
			continue
		}

		behaviorComp, _ := entity.GetComponent("behavior")
		if behavior, ok := behaviorComp.(*BehaviorComponent); ok {
			// Process AI behavior
			switch behavior.AIType {
			case "passive":
				bs.processPassiveBehavior(entity, behavior, deltaTime)
			case "hostile":
				bs.processHostileBehavior(entity, behavior, deltaTime)
			case "neutral":
				bs.processNeutralBehavior(entity, behavior, deltaTime)
			}
		}
	}
}

func (bs *BehaviorSystem) processPassiveBehavior(entity *Entity, behavior *BehaviorComponent, deltaTime float64) {
	// Passive entities don't attack unless provoked
}

func (bs *BehaviorSystem) processHostileBehavior(entity *Entity, behavior *BehaviorComponent, deltaTime float64) {
	// Hostile entities actively seek and attack targets
	// Target acquisition logic would go here
}

func (bs *BehaviorSystem) processNeutralBehavior(entity *Entity, behavior *BehaviorComponent, deltaTime float64) {
	// Neutral entities defend themselves but don't seek targets
	// Self-defense logic would go here
}

// InventorySystem handles inventory management
type InventorySystem struct {
	name string
	requiredComponents []string
}

func NewInventorySystem() *InventorySystem {
	return &InventorySystem{
		name: "inventory",
		requiredComponents: []string{"inventory"},
	}
}

func (is *InventorySystem) GetName() string { return is.name }
func (is *InventorySystem) GetRequiredComponents() []string { return is.requiredComponents }
func (is *InventorySystem) Matches(entity *Entity) bool {
	return entity.HasComponent("inventory")
}

func (is *InventorySystem) Process(deltaTime float64, entities []*Entity) {
	for _, entity := range entities {
		if !is.Matches(entity) {
			continue
		}

		inventoryComp, _ := entity.GetComponent("inventory")
		if inventory, ok := inventoryComp.(*InventoryComponent); ok {
			// Update durability if applicable
			if inventory.CurrentDurability > 0 && inventory.MaxDurability > 0 {
				// Durability decay logic would go here
			}

			// Process container contents
			if inventory.Container {
				// Container logic would go here
			}
		}
	}
}

// CraftingSystem handles crafting operations
type CraftingSystem struct {
	name string
	requiredComponents []string
	recipes map[string]*CraftingRecipe
}

type CraftingRecipe struct {
	ID       string
	Name     string
	Inputs   map[string]int
	Outputs  map[string]int
	RequiredTool string
	CraftingTime time.Duration
	Category string
	Tier     int
}

func NewCraftingSystem() *CraftingSystem {
	return &CraftingSystem{
		name: "crafting",
		requiredComponents: []string{"crafting"},
		recipes: make(map[string]*CraftingRecipe),
	}
}

func (cs *CraftingSystem) GetName() string { return cs.name }
func (cs *CraftingSystem) GetRequiredComponents() []string { return cs.requiredComponents }
func (cs *CraftingSystem) Matches(entity *Entity) bool {
	return entity.HasComponent("crafting")
}

func (cs *CraftingSystem) Process(deltaTime float64, entities []*Entity) {
	for _, entity := range entities {
		if !cs.Matches(entity) {
			continue
		}

		craftingComp, _ := entity.GetComponent("crafting")
		if crafting, ok := craftingComp.(*CraftingComponent); ok {
			if crafting.Craftable {
				// Crafting logic would go here
			}
		}
	}
}

// AddRecipe adds a crafting recipe to the system
func (cs *CraftingSystem) AddRecipe(recipe *CraftingRecipe) {
	cs.recipes[recipe.ID] = recipe
}

// GetRecipe returns a recipe by ID
func (cs *CraftingSystem) GetRecipe(id string) (*CraftingRecipe, bool) {
	recipe, exists := cs.recipes[id]
	return recipe, exists
}

// CanCraft checks if an entity can craft a specific recipe
func (cs *CraftingSystem) CanCraft(entity *Entity, recipeID string) bool {
	recipe, exists := cs.GetRecipe(recipeID)
	if !exists {
		return false
	}

	inventoryComp, hasInventory := entity.GetComponent("inventory")
	if !hasInventory {
		return false
	}

	inventory := inventoryComp.(*InventoryComponent)

	// Check if entity has required items
	for itemID, quantity := range recipe.Inputs {
		if currentQty, has := inventory.Contents[itemID]; !has || currentQty < quantity {
			return false
		}
	}

	// Check if entity has required tool
	if recipe.RequiredTool != "" {
		// Tool checking logic would go here
	}

	return true
}

// ToolSystem handles tool operations and durability
type ToolSystem struct {
	name string
	requiredComponents []string
}

func NewToolSystem() *ToolSystem {
	return &ToolSystem{
		name: "tool",
		requiredComponents: []string{"tool"},
	}
}

func (ts *ToolSystem) GetName() string { return ts.name }
func (ts *ToolSystem) GetRequiredComponents() []string { return ts.requiredComponents }
func (ts *ToolSystem) Matches(entity *Entity) bool {
	return entity.HasComponent("tool")
}

func (ts *ToolSystem) Process(deltaTime float64, entities []*Entity) {
	for _, entity := range entities {
		if !ts.Matches(entity) {
			continue
		}

		toolComp, _ := entity.GetComponent("tool")
		if tool, ok := toolComp.(*ToolComponent); ok {
			// Update tool durability
			if tool.Durability > 0 && tool.MaxDurability > 0 {
				// Durability logic would go here
			}
		}
	}
}

// UseTool uses a tool and reduces durability
func (ts *ToolSystem) UseTool(entity *Entity, targetBlockType string) bool {
	toolComp, hasTool := entity.GetComponent("tool")
	if !hasTool {
		return false
	}

	tool := toolComp.(*ToolComponent)

	// Check if tool is effective against target
	isEffective := false
	for _, effective := range tool.Effective {
		if effective == "all" {
			isEffective = true
			break
		}
	}

	if !isEffective {
		return false
	}

	// Reduce durability
	if tool.Durability > 0 {
		tool.Durability--
		if tool.Durability <= 0 {
			// Tool is broken
			return false
		}
	}

	return true
}

// CombatSystem handles combat mechanics
type CombatSystem struct {
	name string
	requiredComponents []string
}

func NewCombatSystem() *CombatSystem {
	return &CombatSystem{
		name: "combat",
		requiredComponents: []string{"combat"},
	}
}

func (cs *CombatSystem) GetName() string { return cs.name }
func (cs *CombatSystem) GetRequiredComponents() []string { return cs.requiredComponents }
func (cs *CombatSystem) Matches(entity *Entity) bool {
	return entity.HasComponent("combat")
}

func (cs *CombatSystem) Process(deltaTime float64, entities []*Entity) {
	for _, entity := range entities {
		if !cs.Matches(entity) {
			continue
		}

		combatComp, _ := entity.GetComponent("combat")
		if combat, ok := combatComp.(*CombatComponent); ok {
			// Health regeneration
			if combat.Regeneration > 0 && combat.Health < combat.MaxHealth {
				combat.Health += combat.Regeneration * deltaTime
				if combat.Health > combat.MaxHealth {
					combat.Health = combat.MaxHealth
				}
			}

			// Process combat effects
			for range combat.Effects {
				// Effect processing would go here
			}
		}
	}
}

// Attack performs an attack from one entity to another
func (cs *CombatSystem) Attack(attacker, target *Entity) bool {
	attackerCombat, hasAttackerCombat := attacker.GetComponent("combat")
	targetCombat, hasTargetCombat := target.GetComponent("combat")

	if !hasAttackerCombat || !hasTargetCombat {
		return false
	}

	attackerStats := attackerCombat.(*CombatComponent)
	targetStats := targetCombat.(*CombatComponent)

	// Calculate damage
	damage := attackerStats.Damage

	// Apply armor reduction
	damage -= targetStats.Armor
	if damage < 0 {
		damage = 0
	}

	// Apply damage to target
	targetStats.Health -= damage

	// Check if target is defeated
	if targetStats.Health <= 0 {
		targetStats.Health = 0
		// Target defeated logic would go here
		return true
	}

	return false
}
