package entities

import (
	"fmt"
	"image/color"
	"log"
	"sync"
	"time"

	"tesselbox/pkg/blocks"
	"tesselbox/pkg/items"
	"tesselbox/pkg/organisms"
	"tesselbox/pkg/player"
	"tesselbox/pkg/world"
)

// GameIntegration integrates the entity system with the main game
type GameIntegration struct {
	migrationManager *MigrationManager
	bridge          *Bridge
	
	// Game references
	gameWorld       *world.World
	gamePlayer      *player.Player
	
	// Entity tracking
	entityEntities  map[string]*Entity // Maps game objects to entities
	
	// Performance tracking
	lastUpdateTime  time.Time
	updateCount     int
}

// NewGameIntegration creates a new game integration
func NewGameIntegration() *GameIntegration {
	return &GameIntegration{
		migrationManager: GetGlobalMigrationManager(),
		bridge:          GetGlobalMigrationManager().GetBridge(),
		entityEntities:  make(map[string]*Entity),
		lastUpdateTime:  time.Now(),
	}
}

// Initialize initializes the game integration
func (gi *GameIntegration) Initialize(gameWorld *world.World, gamePlayer *player.Player) error {
	gi.gameWorld = gameWorld
	gi.gamePlayer = gamePlayer
	
	log.Printf("Game integration initialized")
	return nil
}

// Update updates the entity system
func (gi *GameIntegration) Update(deltaTime float64) {
	// Update migration manager
	gi.migrationManager.Update(deltaTime)
	
	// Update entity system
	gi.bridge.Update(deltaTime)
	
	// Track performance
	gi.updateCount++
	if time.Since(gi.lastUpdateTime) >= time.Second {
		log.Printf("Entity System: %d updates/sec", gi.updateCount)
		gi.updateCount = 0
		gi.lastUpdateTime = time.Now()
	}
}

// ============================================================================
// Block System Integration
// ============================================================================

// GetBlockColorIntegrated gets block color with migration support
func (gi *GameIntegration) GetBlockColorIntegrated(blockType string) color.RGBA {
	if gi.migrationManager.ShouldUseNewSystem("blocks") {
		gi.migrationManager.RecordNewSystemUsage("block_color")
		return gi.bridge.GetBlockColor(blockType)
	}
	
	gi.migrationManager.RecordOldSystemUsage("block_color")
	return blocks.ColorByType(blockType)
}

// GetBlockHardnessIntegrated gets block hardness with migration support
func (gi *GameIntegration) GetBlockHardnessIntegrated(blockType string) float64 {
	if gi.migrationManager.ShouldUseNewSystem("blocks") {
		gi.migrationManager.RecordNewSystemUsage("block_hardness")
		return gi.bridge.GetBlockHardness(blockType)
	}
	
	gi.migrationManager.RecordOldSystemUsage("block_hardness")
	return blocks.HardnessByType(blockType)
}

// CreateBlockIntegrated creates a block using the new system if enabled
func (gi *GameIntegration) CreateBlockIntegrated(blockType string, x, y, z float64) error {
	if gi.migrationManager.ShouldUseNewSystem("blocks") {
		gi.migrationManager.RecordNewSystemUsage("create_block")
		return gi.bridge.CreateBlock(blockType, x, y, z)
	}
	
	gi.migrationManager.RecordOldSystemUsage("create_block")
	return nil // Use old system
}

// ============================================================================
// Item System Integration
// ============================================================================

// GetItemPropertiesIntegrated gets item properties with migration support
func (gi *GameIntegration) GetItemPropertiesIntegrated(itemType items.ItemType) map[string]interface{} {
	if gi.migrationManager.ShouldUseNewSystem("items") {
		gi.migrationManager.RecordNewSystemUsage("item_properties")
		return gi.bridge.GetItemProperties(itemType)
	}
	
	gi.migrationManager.RecordOldSystemUsage("item_properties")
	oldProps := items.GetItemProperties(itemType)
	return map[string]interface{}{
		"name":        oldProps.Name,
		"stackSize":   oldProps.StackSize,
		"durability":  oldProps.Durability,
		"isTool":      oldProps.IsTool,
		"toolPower":   oldProps.ToolPower,
	}
}

// CreateItemIntegrated creates an item using the new system if enabled
func (gi *GameIntegration) CreateItemIntegrated(itemType string, quantity int, playerID string) error {
	if gi.migrationManager.ShouldUseNewSystem("items") {
		gi.migrationManager.RecordNewSystemUsage("create_item")
		return gi.bridge.CreateItem(itemType, quantity, playerID)
	}
	
	gi.migrationManager.RecordOldSystemUsage("create_item")
	return nil // Use old system
}

// ============================================================================
// Organism System Integration
// ============================================================================

// CreateAdvancedEntity creates an entity with custom components
func (gi *GameIntegration) CreateAdvancedEntity(entityType, entityID string, x, y, z float64, quantity int, playerID, componentType string) error {
	if gi.migrationManager.ShouldUseNewSystem("blocks") || 
	   gi.migrationManager.ShouldUseNewSystem("items") || 
	   gi.migrationManager.ShouldUseNewSystem("organisms") {
		
		gi.migrationManager.RecordNewSystemUsage("create_advanced_entity")
		
		// Create entity with custom component
		entity := NewEntity(entityID, entityType)
		entity.Metadata["position"] = [3]float64{x, y, z}
		
		// Add custom component based on type
		switch componentType {
		case "render":
			renderComp := &RenderComponent{
				Type:       "render",
				Color:      color.RGBA{255, 255, 255, 255},
				Visible:    true,
				LightLevel: 15,
				Scale:      1.0,
			}
			entity.AddComponent(renderComp)
			
		case "physics":
			physicsComp := &PhysicsComponent{
				Type:       "physics",
				Hardness:   1.0,
				Density:    1.0,
				Solid:      true,
				Gravity:    true,
				Mass:       1.0,
			}
			entity.AddComponent(physicsComp)
			
		case "behavior":
			behaviorComp := &BehaviorComponent{
				Type:         "behavior",
				AIType:       "basic",
				Passive:      true,
				CurrentState: "idle",
				SightRange:   10.0,
				Speed:        1.0,
			}
			entity.AddComponent(behaviorComp)
			
		case "inventory":
			inventoryComp := &InventoryComponent{
				Type:         "inventory",
				StackSize:    64,
				MaxDurability: 100,
				Container:    true,
				Slots:        32,
				Contents:     make(map[string]int),
				Weight:       0.0,
			}
			entity.AddComponent(inventoryComp)
			
		case "crafting":
			craftingComp := &CraftingComponent{
				Type:         "crafting",
				Craftable:    true,
				Recipe:       make(map[string]int),
				Results:      make(map[string]int),
				CraftingTime: 3000 * time.Millisecond,
				Category:     "custom",
				Tier:         1,
			}
			entity.AddComponent(craftingComp)
			
		default:
			return fmt.Errorf("unknown component type: %s", componentType)
		}
		
		// Add basic metadata
		entity.Metadata["name"] = entityID
		entity.Metadata["description"] = fmt.Sprintf("Advanced %s with %s component", entityType, componentType)
		entity.Metadata["created_at"] = time.Now()
		entity.Tags = []string{"advanced", componentType, entityType}
		
		// Add entity to the world
		world := gi.bridge.GetWorld()
		world.GetSystemManager().AddEntity(entity)
		
		return nil
	}
	
	gi.migrationManager.RecordOldSystemUsage("create_advanced_entity")
	return fmt.Errorf("advanced entity creation not available in old system")
}

// GetOrganismPropertiesIntegrated gets organism properties with migration support
func (gi *GameIntegration) GetOrganismPropertiesIntegrated(orgType organisms.OrganismType) map[string]interface{} {
	if gi.migrationManager.ShouldUseNewSystem("organisms") {
		gi.migrationManager.RecordNewSystemUsage("organism_properties")
		return gi.bridge.GetOrganismProperties(orgType)
	}
	
	gi.migrationManager.RecordOldSystemUsage("organism_properties")
	// Convert int to string for lookup
	orgName := ""
	for name, oType := range organisms.OrganismTypeMap {
		if oType == orgType {
			orgName = name
			break
		}
	}
	return map[string]interface{}{
		"type": orgName,
	}
}

// CreateOrganismIntegrated creates an organism using the new system if enabled
func (gi *GameIntegration) CreateOrganismIntegrated(orgType string, x, y, z float64) error {
	if gi.migrationManager.ShouldUseNewSystem("organisms") {
		gi.migrationManager.RecordNewSystemUsage("create_organism")
		return gi.bridge.CreateOrganism(orgType, x, y, z)
	}
	
	gi.migrationManager.RecordOldSystemUsage("create_organism")
	return nil // Use old system
}

// ============================================================================
// Event Integration
// ============================================================================

// PublishBlockPlacedEvent publishes a block placed event
func (gi *GameIntegration) PublishBlockPlacedEvent(blockType string, x, y, z float64, playerID string) {
	if gi.bridge != nil && gi.bridge.world != nil {
		eventBus := gi.bridge.world.GetEventBus()
		eventBus.PublishWithSource(EventBlockPlaced, "game", 
			CreateBlockEvent(blockType, x, y, z, playerID, ""))
	}
}

// PublishBlockBrokenEvent publishes a block broken event
func (gi *GameIntegration) PublishBlockBrokenEvent(blockType string, x, y, z float64, playerID, toolUsed string) {
	if gi.bridge != nil && gi.bridge.world != nil {
		eventBus := gi.bridge.world.GetEventBus()
		eventBus.PublishWithSource(EventBlockBroken, "game", 
			CreateBlockEvent(blockType, x, y, z, playerID, toolUsed))
	}
}

// PublishItemUsedEvent publishes an item used event
func (gi *GameIntegration) PublishItemUsedEvent(itemType string, quantity int, playerID, targetID string, success bool) {
	if gi.bridge != nil && gi.bridge.world != nil {
		eventBus := gi.bridge.world.GetEventBus()
		eventBus.PublishWithSource(EventItemUsed, "game", 
			CreateItemEvent(itemType, quantity, playerID, targetID, success))
	}
}

// PublishItemCraftedEvent publishes an item crafted event
func (gi *GameIntegration) PublishItemCraftedEvent(itemType string, quantity int, playerID string) {
	if gi.bridge != nil && gi.bridge.world != nil {
		eventBus := gi.bridge.world.GetEventBus()
		eventBus.PublishWithSource(EventItemCrafted, "game", 
			CreateItemEvent(itemType, quantity, playerID, "", true))
	}
}

// PublishAttackEvent publishes an attack event
func (gi *GameIntegration) PublishAttackEvent(attackerID, targetID string, damage float64, weaponType string, critical bool) {
	if gi.bridge != nil && gi.bridge.world != nil {
		eventBus := gi.bridge.world.GetEventBus()
		eventBus.PublishWithSource(EventAttack, "game", 
			CreateCombatEvent(attackerID, targetID, damage, weaponType, critical, false))
	}
}

// ============================================================================
// Performance and Debugging
// ============================================================================

// GetStatistics returns integration statistics
func (gi *GameIntegration) GetStatistics() map[string]interface{} {
	stats := make(map[string]interface{})
	
	// Migration statistics
	stats["migration"] = gi.migrationManager.GetMigrationProgress()
	
	// Entity world statistics
	if gi.bridge != nil {
		stats["entity_world"] = gi.bridge.GetWorld().GetStatistics()
	}
	
	// Performance statistics
	stats["performance"] = map[string]interface{}{
		"updates_per_second": gi.updateCount,
		"last_update":        gi.lastUpdateTime,
		"entity_count":       len(gi.entityEntities),
	}
	
	return stats
}

// DebugPrint prints debug information
func (gi *GameIntegration) DebugPrint() {
	stats := gi.GetStatistics()
	log.Printf("=== Game Integration Debug ===")
	log.Printf("Migration Step: %v", stats["migration"].(map[string]interface{})["migration_step"])
	log.Printf("Migration Ratio: %.2f%%", stats["migration"].(map[string]interface{})["migration_ratio"].(float64)*100)
	
	if entityStats, ok := stats["entity_world"].(map[string]interface{}); ok {
		log.Printf("Entity Count: %v", entityStats["total_entities"])
	}
	
	log.Printf("Performance: %v updates/sec", stats["performance"].(map[string]interface{})["updates_per_second"])
	log.Printf("===============================")
}

// Shutdown shuts down the game integration
func (gi *GameIntegration) Shutdown() error {
	log.Printf("Shutting down game integration...")
	
	// Print final statistics
	gi.DebugPrint()
	
	// Shutdown migration manager
	return gi.migrationManager.Shutdown()
}

// ============================================================================
// Global Integration Instance
// ============================================================================

var globalGameIntegration *GameIntegration
var gameIntegrationOnce sync.Once

// GetGlobalGameIntegration returns the global game integration
func GetGlobalGameIntegration() *GameIntegration {
	gameIntegrationOnce.Do(func() {
		globalGameIntegration = NewGameIntegration()
	})
	return globalGameIntegration
}

// ============================================================================
// Convenience Functions for Direct Integration
// ============================================================================

// InitializeGameIntegration initializes the global game integration
func InitializeGameIntegration(gameWorld *world.World, gamePlayer *player.Player) error {
	return GetGlobalGameIntegration().Initialize(gameWorld, gamePlayer)
}

// UpdateGameIntegration updates the global game integration
func UpdateGameIntegration(deltaTime float64) {
	GetGlobalGameIntegration().Update(deltaTime)
}

// GetBlockColor gets block color with migration support
func GetBlockColor(blockType string) color.RGBA {
	return GetGlobalGameIntegration().GetBlockColorIntegrated(blockType)
}

// GetBlockHardness gets block hardness with migration support
func GetBlockHardness(blockType string) float64 {
	return GetGlobalGameIntegration().GetBlockHardnessIntegrated(blockType)
}

// GetItemProperties gets item properties with migration support
func GetItemProperties(itemType items.ItemType) map[string]interface{} {
	return GetGlobalGameIntegration().GetItemPropertiesIntegrated(itemType)
}

// GetOrganismProperties gets organism properties with migration support
func GetOrganismProperties(orgType organisms.OrganismType) map[string]interface{} {
	return GetGlobalGameIntegration().GetOrganismPropertiesIntegrated(orgType)
}

// PublishGameEvent publishes various game events
func PublishBlockPlaced(blockType string, x, y, z float64, playerID string) {
	GetGlobalGameIntegration().PublishBlockPlacedEvent(blockType, x, y, z, playerID)
}

func PublishBlockBroken(blockType string, x, y, z float64, playerID, toolUsed string) {
	GetGlobalGameIntegration().PublishBlockBrokenEvent(blockType, x, y, z, playerID, toolUsed)
}

func PublishItemUsed(itemType string, quantity int, playerID, targetID string, success bool) {
	GetGlobalGameIntegration().PublishItemUsedEvent(itemType, quantity, playerID, targetID, success)
}

func PublishItemCrafted(itemType string, quantity int, playerID string) {
	GetGlobalGameIntegration().PublishItemCraftedEvent(itemType, quantity, playerID)
}

func PublishAttack(attackerID, targetID string, damage float64, weaponType string, critical bool) {
	GetGlobalGameIntegration().PublishAttackEvent(attackerID, targetID, damage, weaponType, critical)
}
