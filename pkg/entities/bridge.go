package entities

import (
	"image/color"
	"log"

	"tesselbox/pkg/blocks"
	"tesselbox/pkg/items"
	"tesselbox/pkg/organisms"
)

// Bridge provides compatibility between old and new systems
type Bridge struct {
	world *GameWorld
}

// NewBridge creates a new bridge instance
func NewBridge() *Bridge {
	world := NewGameWorld()
	if err := world.Initialize(); err != nil {
		log.Printf("Warning: Failed to initialize entity world: %v", err)
	}

	return &Bridge{world: world}
}

// GetWorld returns the entity world
func (b *Bridge) GetWorld() *GameWorld {
	return b.world
}

// Legacy compatibility methods for blocks
func (b *Bridge) GetBlockColor(blockType string) color.RGBA {
	if blockColor, err := b.world.GetBlockColor(blockType); err == nil {
		// Try to convert from RGBA slice
		if rgbaSlice, ok := blockColor.([]uint8); ok && len(rgbaSlice) == 4 {
			return color.RGBA{R: rgbaSlice[0], G: rgbaSlice[1], B: rgbaSlice[2], A: rgbaSlice[3]}
		}
	}
	// Fallback to old system
	return blocks.ColorByType(blockType)
}

func (b *Bridge) GetBlockHardness(blockType string) float64 {
	// Try new system first
	entities := b.world.GetEntitiesByType("block")
	for _, entity := range entities {
		if entity.Type == blockType {
			if physics, has := entity.GetComponent("physics"); has {
				if phys, ok := physics.(*PhysicsComponent); ok {
					return phys.Hardness
				}
			}
		}
	}
	// Fallback to old system
	return blocks.HardnessByType(blockType)
}

// Legacy compatibility methods for items
func (b *Bridge) GetItemProperties(itemType items.ItemType) map[string]interface{} {
	itemName := items.ItemNameByID(itemType)
	if props, err := b.world.GetItemProperties(itemName); err == nil {
		if propMap, ok := props.(map[string]interface{}); ok {
			return propMap
		}
	}
	// Fallback to old system
	oldProps := items.GetItemProperties(itemType)
	return map[string]interface{}{
		"name":       oldProps.Name,
		"stackSize":  oldProps.StackSize,
		"durability": oldProps.Durability,
		"isTool":     oldProps.IsTool,
		"toolPower":  oldProps.ToolPower,
	}
}

// Legacy compatibility methods for organisms
func (b *Bridge) GetOrganismProperties(orgType organisms.OrganismType) map[string]interface{} {
	// Convert int to string for lookup
	orgName := ""
	for name, oType := range organisms.OrganismTypeMap {
		if oType == orgType {
			orgName = name
			break
		}
	}
	if orgName == "" {
		return map[string]interface{}{"type": "unknown"}
	}

	if props, err := b.world.GetOrganismProperties(orgName); err == nil {
		if propMap, ok := props.(map[string]interface{}); ok {
			return propMap
		}
	}
	// Fallback to old system
	return map[string]interface{}{
		"type": orgName,
	}
}

// Create methods using new system
func (b *Bridge) CreateBlock(blockType string, x, y, z float64) error {
	_, err := b.world.CreateBlock(blockType, x, y, z)
	return err
}

func (b *Bridge) CreateItem(itemType string, quantity int, playerID string) error {
	_, err := b.world.CreateItem(itemType, quantity, playerID)
	return err
}

func (b *Bridge) CreateOrganism(orgType string, x, y, z float64) error {
	_, err := b.world.CreateOrganism(orgType, x, y, z)
	return err
}

// Update the entity world
func (b *Bridge) Update(deltaTime float64) {
	if b.world != nil && b.world.IsInitialized() {
		b.world.Update(deltaTime)
	}
}

// Shutdown the bridge
func (b *Bridge) Shutdown() error {
	if b.world != nil {
		return b.world.Shutdown()
	}
	return nil
}
