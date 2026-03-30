package main

import (
	"log"
	"math/rand"

	"tesselbox/pkg/blocks"
	"tesselbox/pkg/entities"
	"tesselbox/pkg/items"
	"tesselbox/pkg/menu"
	"tesselbox/pkg/player"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Constants for item drop physics
const (
	DropOffsetX      = 30.0
	DropOffsetY      = 30.0
	DropVelocityBase = 3.0
	DropVelocityVar  = 2.0
	DropLifeSeconds  = 300
	RandomDivisor    = 100.0
)

func (g *Game) handleBlockPlacement() {
	mouseX, mouseY := ebiten.CursorPosition()
	worldX := float64(mouseX) + g.cameraX
	worldY := float64(mouseY) + g.cameraY

	gridX, gridY := g.world.PixelToGrid(worldX, worldY)

	if g.world.IsValidPosition(gridX, gridY) {
		if g.world.GetBlock(gridX, gridY) == nil {
			// Check distance to player for placement range
			playerX, playerY := g.player.GetPosition()
			dx := worldX - playerX
			dy := worldY - playerY
			dist := dx*dx + dy*dy
			placeRange := g.player.GetPlaceRange() * g.player.GetPlaceRange()

			if dist <= placeRange {
				// Try to get block from inventory
				inv := g.player.GetInventory()
				hotbarSlot := inv.GetSelectedSlot()
				itemStack := inv.GetHotbarSlot(hotbarSlot)

				if itemStack != nil && itemStack.Count > 0 {
					// Get the block type from the item
					blockType, exists := blocks.GetBlockByID(itemStack.Item.ID)
					if exists {
						// Place the block
						if err := g.world.PlaceBlock(gridX, gridY, blockType); err == nil {
							// Decrement item count if not in creative mode
							if !g.CreativeMode {
								inv.RemoveFromHotbar(hotbarSlot, 1)
							}
							g.audioManager.PlayPlaceSound()
						}
					}
				}
			}
		}
	}
}

func (g *Game) handleMining() {
	mouseX, mouseY := ebiten.CursorPosition()
	worldX := float64(mouseX) + g.cameraX
	worldY := float64(mouseY) + g.cameraY

	gridX, gridY := g.world.PixelToGrid(worldX, worldY)

	// Check distance to player
	playerX, playerY := g.player.GetPosition()
	dx := worldX - playerX
	dy := worldY - playerY
	dist := dx*dx + dy*dy

	if dist > g.player.GetReach()*g.player.GetReach() {
		return
	}

	if !g.world.IsValidPosition(gridX, gridY) {
		return
	}

	block := g.world.GetBlock(gridX, gridY)
	if block == nil {
		return
	}

	// Calculate mining damage
	damage := g.player.GetMiningSpeed()
	if g.CreativeMode {
		damage = block.Hardness * 100 // Instant break
	}

	// Apply damage to block
	if g.world.DamageBlock(gridX, gridY, damage) {
		// Block was destroyed
		g.audioManager.PlayMineSound()

		// Create dropped item
		if !g.CreativeMode {
			blockDef, exists := blocks.GetBlock(block.ID)
			if exists && blockDef.DropItem != "" {
				item, itemExists := items.GetItem(blockDef.DropItem)
				if itemExists {
					// Create dropped item with some random velocity
					velX := (rand.Float64()*DropVelocityVar - DropVelocityVar/2) + DropVelocityBase*float64(rand.Intn(2)*2-1)
					velY := (rand.Float64()*DropVelocityVar - DropVelocityVar/2) - DropVelocityBase

					droppedItem := items.NewDroppedItem(
						item,
						worldX+DropOffsetX,
						worldY+DropOffsetY,
						velX,
						velY,
						DropLifeSeconds,
					)
					g.droppedItems = append(g.droppedItems, droppedItem)
				}
			}
		}
	}
}

func (g *Game) handleMenuAction(action menu.MenuAction) {
	switch action {
	case menu.ActionStartGame:
		g.inMenu = false
		g.inGame = true

		spawnX, spawnY := g.world.FindSpawnPosition(-2000, -2000)
		g.player.SetPosition(spawnX, spawnY)
		g.player.SetVelocity(0, 0)

	case menu.ActionBack:
		g.menu.SetMainMenu()

	case menu.ActionOpenBlockLibrary:
		g.menu.SetBlockLibraryMenu()

	case menu.ActionSelectBlock:
		g.selectedBlock = g.menu.SelectedBlock
		g.menu.SetMainMenu()

	case menu.ActionOpenSettings:
		g.menu.SetSettingsMenu()

	case menu.ActionToggleSound:
		if g.settingsManager != nil {
			muted := g.settingsManager.GetBool("mute")
			if err := g.settingsManager.SetBool("mute", !muted); err != nil {
				log.Printf("Failed to toggle sound: %v", err)
			} else {
				g.audioManager.UpdateSettings()
				if muted {
					log.Println("Sound enabled")
				} else {
					log.Println("Sound muted")
				}
			}
			g.menu.SetSettingsMenu()
		}

	case menu.ActionToggleFullscreen:
		if g.settingsManager != nil {
			fullscreen := g.settingsManager.GetBool("fullscreen")
			if err := g.settingsManager.SetBool("fullscreen", !fullscreen); err != nil {
				log.Printf("Failed to toggle fullscreen: %v", err)
			} else {
				ebiten.SetFullscreen(!fullscreen)
			}
			g.menu.SetSettingsMenu()
		}

	case menu.ActionToggleDebug:
		if g.settingsManager != nil {
			debug := g.settingsManager.GetBool("show_debug")
			if err := g.settingsManager.SetBool("show_debug", !debug); err != nil {
				log.Printf("Failed to toggle debug: %v", err)
			} else {
				g.showDebug = !debug
			}
			g.menu.SetSettingsMenu()
		}

	case menu.ActionExit:
		log.Println("Exiting game...")
	}
}

func (g *Game) updateDroppedItems() {
	for i := len(g.droppedItems) - 1; i >= 0; i-- {
		item := g.droppedItems[i]
		item.Update()

		// Check for player pickup
		playerX, playerY := g.player.GetPosition()
		if item.CanPickup(playerX, playerY) {
			// Try to add to inventory
			if g.player.GetInventory().AddItem(item.Item, 1) {
				g.audioManager.PlayInventorySound()
				// Remove from dropped items
				g.droppedItems = append(g.droppedItems[:i], g.droppedItems[i+1:]...)
			}
		}

		// Remove expired items
		if item.IsExpired() {
			g.droppedItems = append(g.droppedItems[:i], g.droppedItems[i+1:]...)
		}
	}
}
