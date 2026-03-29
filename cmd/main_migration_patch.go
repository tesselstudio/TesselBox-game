package main

import (
	"image/color"
	"log"
	"time"

	"tesselbox/pkg/entities"
)

// Integration with the new entity system
var entityIntegration *entities.GameIntegration

// initializeEntitySystem initializes the entity system integration
func (g *Game) initializeEntitySystem() {
	// Initialize migration system
	entities.InitializeMigration()
	
	// Create game integration
	entityIntegration = entities.GetGlobalGameIntegration()
	
	// Initialize with game world and player
	if err := entityIntegration.Initialize(g.world, g.player); err != nil {
		log.Printf("Warning: Failed to initialize entity integration: %v", err)
		return
	}
	
	log.Printf("Entity system integration initialized")
}

// updateEntitySystem updates the entity system
func (g *Game) updateEntitySystem(deltaTime float64) {
	if entityIntegration != nil {
		entityIntegration.Update(deltaTime)
	}
}

// ============================================================================
// Patched Functions with Entity System Integration
// ============================================================================

// patchedNewGame patches the NewGame function to include entity system
func patchedNewGame() *Game {
	// Create the original game
	game := NewGame()
	
	// Initialize entity system
	game.initializeEntitySystem()
	
	return game
}

// patchedUpdate patches the Update function to include entity system
func (g *Game) patchedUpdate() error {
	// Calculate delta time
	currentTime := time.Now()
	deltaTime := currentTime.Sub(g.lastTime).Seconds()
	g.lastTime = currentTime
	
	// Update entity system
	g.updateEntitySystem(deltaTime)
	
	// Call original update
	return g.Update()
}

// patchedGetBlockColor patches block color queries to use entity system
func patchedGetBlockColor(blockType string) color.RGBA {
	if entityIntegration != nil {
		return entities.GetBlockColor(blockType)
	}
	// Fallback to original
	return blocks.ColorByType(blockType)
}

// patchedGetBlockHardness patches block hardness queries to use entity system
func patchedGetBlockHardness(blockType string) float64 {
	if entityIntegration != nil {
		return entities.GetBlockHardness(blockType)
	}
	// Fallback to original
	return blocks.HardnessByType(blockType)
}

// patchedCompleteMining patches mining completion to publish events
func (g *Game) patchedCompleteMining(targetHex *world.Hexagon) {
	// Call original completeMining
	g.completeMining(targetHex)
	
	// Publish block broken event
	if entityIntegration != nil {
		blockType := getBlockKeyFromType(targetHex.BlockType)
		entities.PublishBlockBroken(blockType, targetHex.X, targetHex.Y, 0, "player")
	}
}

// patchedHandleBlockPlacement patches block placement to publish events
func (g *Game) patchedHandleBlockPlacement() {
	// Call original handleBlockPlacement
	g.handleBlockPlacement()
	
	// Publish block placed event (if successful)
	if entityIntegration != nil {
		// This would need to be enhanced to track successful placement
		var blockType string
		if g.CreativeMode && g.selectedBlock != "" {
			blockType = g.selectedBlock
		} else {
			selectedItem := g.inventory.GetSelectedItem()
			if selectedItem != nil && selectedItem.Type != items.NONE {
				props := items.GetItemProperties(selectedItem.Type)
				if props != nil && props.IsPlaceable {
					blockType = props.BlockType
				}
			}
		}
		
		if blockType != "" {
			mouseWorldX := float64(g.mouseX) + g.cameraX
			mouseWorldY := float64(g.mouseY) + g.cameraY
			entities.PublishBlockPlaced(blockType, mouseWorldX, mouseWorldY, 0, "player")
		}
	}
}

// patchedDropItem patches item dropping to publish events
func (g *Game) patchedDropItem() {
	selectedItem := g.inventory.GetSelectedItem()
	if selectedItem == nil || selectedItem.Type == items.NONE {
		return
	}
	
	// Publish item used event before dropping
	if entityIntegration != nil {
		itemName := items.ItemNameByID(selectedItem.Type)
		entities.PublishItemUsed(itemName, 1, "player", "", true)
	}
	
	// Call original dropItem
	g.dropItem()
}

// ============================================================================
// Migration Commands
// ============================================================================

// executeEntityCommand executes entity system commands
func (g *Game) executeEntityCommand(command string) {
	// Parse command
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return
	}
	
	cmd := strings.ToLower(parts[0])
	args := parts[1:]
	
	switch cmd {
	case "migration":
		if len(args) == 0 {
			g.showMigrationStatus()
			return
		}
		
		subCmd := strings.ToLower(args[0])
		switch subCmd {
		case "step":
			if len(args) < 2 {
				log.Printf("Usage: /migration step <step>")
				return
			}
			step, err := strconv.Atoi(args[1])
			if err != nil {
				log.Printf("Invalid step: %s", args[1])
				return
			}
			entities.SetMigrationStepCommand(step)
			
		case "status":
			g.showMigrationStatus()
			
		case "toggle":
			entities.ToggleMigrationCommand()
			
		default:
			log.Printf("Unknown migration command: %s", subCmd)
		}
		
	case "entities":
		g.showEntityStatus()
		
	case "debug":
		if entityIntegration != nil {
			entityIntegration.DebugPrint()
		}
		
	default:
		log.Printf("Unknown entity command: %s", cmd)
	}
}

// showMigrationStatus shows migration status
func (g *Game) showMigrationStatus() {
	status := entities.GetMigrationStatusCommand()
	log.Printf("=== Migration Status ===")
	log.Printf("Step: %v", status["migration_step"])
	log.Printf("Enabled: %v", status["enabled"])
	log.Printf("Ratio: %.2f%%", status["migration_ratio"].(float64)*100)
	
	if oldUsage, ok := status["old_system_usage"].(map[string]int); ok {
		log.Printf("Old System Usage: %v", oldUsage)
	}
	
	if newUsage, ok := status["new_system_usage"].(map[string]int); ok {
		log.Printf("New System Usage: %v", newUsage)
	}
	log.Printf("========================")
}

// showEntityStatus shows entity system status
func (g *Game) showEntityStatus() {
	if entityIntegration == nil {
		log.Printf("Entity system not initialized")
		return
	}
	
	stats := entityIntegration.GetStatistics()
	log.Printf("=== Entity System Status ===")
	
	if migration, ok := stats["migration"].(map[string]interface{}); ok {
		log.Printf("Migration Step: %v", migration["migration_step"])
		log.Printf("Migration Ratio: %.2f%%", migration["migration_ratio"].(float64)*100)
	}
	
	if entityWorld, ok := stats["entity_world"].(map[string]interface{}); ok {
		log.Printf("Total Entities: %v", entityWorld["total_entities"])
	}
	
	if performance, ok := stats["performance"].(map[string]interface{}); ok {
		log.Printf("Updates/sec: %v", performance["updates_per_second"])
	}
	
	log.Printf("==========================")
}

// ============================================================================
// Patch Application
// ============================================================================

// applyEntitySystemPatches applies all entity system patches
func applyEntitySystemPatches() {
	// This function would be called to replace the original functions
	// In a real implementation, this would be done more elegantly
	log.Printf("Entity system patches applied")
}

// ============================================================================
// Cleanup
// ============================================================================

// shutdownEntitySystem shuts down the entity system
func shutdownEntitySystem() {
	if entityIntegration != nil {
		if err := entityIntegration.Shutdown(); err != nil {
			log.Printf("Error shutting down entity integration: %v", err)
		}
	}
	
	entities.ShutdownMigration()
	log.Printf("Entity system shut down")
}
