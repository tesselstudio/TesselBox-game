package main

import (
	"image/color"

	"tesselbox/pkg/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// DrawingExample shows how to adapt the drawing logic to use GameManager
// This is a reference implementation for migrating the Draw() method

// DrawGameScene draws the main game world
func DrawGameScene(screen *ebiten.Image, gm *game.GameManager) {
	// Clear screen
	screen.Fill(color.RGBA{135, 206, 235, 255}) // Sky blue

	// Draw world using camera position
	drawWorld(screen, gm)

	// Draw UI overlay
	drawGameUI(screen, gm)
}

// drawWorld draws the game world with camera offset
func drawWorld(screen *ebiten.Image, gm *game.GameManager) {
	// Use camera position from GameManager
	cameraX := gm.CameraX
	cameraY := gm.CameraY

	// Draw hexagons/blocks
	if gm.World != nil {
		// Ensure chunks are generated around camera position
		gm.World.GetChunksInRange(cameraX, cameraY)

		// Get visible hexagons for rendering
		hexagons := gm.World.GetNearbyHexagons(cameraX, cameraY, 1000)

		for _, hex := range hexagons {
			if hex == nil {
				continue
			}

			// Calculate screen position
			screenX := hex.X - cameraX
			screenY := hex.Y - cameraY

			// Only draw if on screen
			if screenX > -100 && screenX < float64(ScreenWidth)+100 &&
				screenY > -100 && screenY < float64(ScreenHeight)+100 {
				// Draw hexagon (implementation depends on your rendering system)
				drawHexagon(screen, screenX, screenY, hex, gm)
			}
		}
	}

	// Draw player
	if gm.Player != nil {
		playerScreenX := gm.Player.X - cameraX
		playerScreenY := gm.Player.Y - cameraY
		drawPlayer(screen, playerScreenX, playerScreenY, gm.Player)
	}

	// Draw dropped items
	for _, item := range gm.DroppedItems {
		itemScreenX := item.X - cameraX
		itemScreenY := item.Y - cameraY
		drawDroppedItem(screen, itemScreenX, itemScreenY, item)
	}
}

// drawGameUI draws the game UI elements
func drawGameUI(screen *ebiten.Image, gm *game.GameManager) {
	// Draw hotbar
	if gm.Inventory != nil {
		drawHotbar(screen, gm.Inventory)
	}

	// Draw health bar
	if gm.HealthSystem != nil {
		drawHealthBar(screen, gm.HealthSystem)
	}

	// Draw HUD
	if gm.HUD != nil {
		gm.HUD.Draw(screen)
	}

	// Draw debug info if enabled
	drawDebugInfo(screen, gm)
}

// drawHexagon draws a single hexagon
func drawHexagon(screen *ebiten.Image, x, y float64, hex interface{}, gm *game.GameManager) {
	// This is a placeholder - implement based on your hexagon rendering
	// Use gm.WhiteImage for drawing if needed
	ebitenutil.DrawRect(screen, x, y, 30, 30, color.RGBA{100, 100, 100, 255})
}

// drawPlayer draws the player character
func drawPlayer(screen *ebiten.Image, x, y float64, player interface{}) {
	// Draw player sprite or placeholder
	ebitenutil.DrawRect(screen, x-10, y-10, 20, 20, color.RGBA{255, 0, 0, 255})
}

// drawDroppedItem draws a dropped item
func drawDroppedItem(screen *ebiten.Image, x, y float64, item *game.DroppedItem) {
	// Draw item sprite or placeholder
	ebitenutil.DrawRect(screen, x-5, y-5, 10, 10, color.RGBA{255, 255, 0, 255})
}

// drawHotbar draws the inventory hotbar
func drawHotbar(screen *ebiten.Image, inventory interface{}) {
	// Draw hotbar background
	ebitenutil.DrawRect(screen, 10, float64(ScreenHeight)-60, 300, 50, color.RGBA{50, 50, 50, 200})

	// Draw hotbar slots
	for i := 0; i < 9; i++ {
		x := float64(20 + i*30)
		y := float64(ScreenHeight - 50)
		ebitenutil.DrawRect(screen, x, y, 25, 25, color.RGBA{100, 100, 100, 255})
	}
}

// drawHealthBar draws the player health bar
func drawHealthBar(screen *ebiten.Image, healthSystem interface{}) {
	// Draw health bar background
	ebitenutil.DrawRect(screen, 10, 10, 200, 20, color.RGBA{50, 50, 50, 200})

	// Draw health bar fill (placeholder - 75% health)
	ebitenutil.DrawRect(screen, 10, 10, 150, 20, color.RGBA{255, 0, 0, 255})
}

// drawDebugInfo draws debug information
func drawDebugInfo(screen *ebiten.Image, gm *game.GameManager) {
	// Draw FPS
	ebitenutil.DebugPrint(screen, "Debug Info")

	// Draw camera position
	debugText := "Camera: " + string(rune(gm.CameraX)) + ", " + string(rune(gm.CameraY))
	text.Draw(screen, debugText, basicfont.Face7x13, 10, 50, color.White)

	// Draw music status
	musicStatus := "Music: "
	if gm.BackgroundMusicManager.IsPlaying() {
		musicStatus += "Playing (" + gm.BackgroundMusicManager.GetCurrentTrack() + ")"
	} else {
		musicStatus += "Stopped"
	}
	text.Draw(screen, musicStatus, basicfont.Face7x13, 10, 70, color.White)

	// Draw profiler info if enabled
	if gm.Profiler != nil {
		// gm.Profiler.Draw(screen)
	}
}

// Example of how to update the GameWrapper.Draw() method:
/*
func (gw *GameWrapper) Draw(screen *ebiten.Image) {
	state := gw.manager.StateManager.GetState()

	switch state {
	case ui.StateCrafting:
		// Draw game in background
		DrawGameScene(screen, gw.manager)
		// Draw crafting UI overlay
		gw.manager.CraftingUI.Draw(screen)

	case ui.StateBackpack:
		// Draw game in background
		DrawGameScene(screen, gw.manager)
		// Draw backpack UI overlay
		gw.manager.BackpackUI.Draw(screen)

	case ui.StateChest:
		// Draw game in background
		DrawGameScene(screen, gw.manager)
		// Draw chest UI overlay
		if gw.manager.ChestUI != nil {
			gw.manager.ChestUI.Draw(screen)
		}

	case ui.StatePluginUI:
		// Draw game in background
		DrawGameScene(screen, gw.manager)
		// Draw plugin UI overlay
		if gw.manager.PluginUI != nil {
			gw.manager.PluginUI.Draw(screen)
		}

	case ui.StateSkinEditor:
		// Draw skin editor (full screen)
		gw.manager.SkinEditor.Draw(screen)

	case ui.StateGame:
		// Draw main game scene
		DrawGameScene(screen, gw.manager)

		// Draw damage indicators on top
		if gw.manager.DamageIndicators != nil {
			gw.manager.DamageIndicators.Draw(screen, gw.manager.CameraX, gw.manager.CameraY)
		}

		// Draw screen flash
		if gw.manager.ScreenFlash != nil {
			gw.manager.ScreenFlash.Draw(screen, gw.screenWidth, gw.screenHeight)
		}

		// Draw directional hit indicator
		if gw.manager.DirectionalHitInd != nil {
			gw.manager.DirectionalHitInd.Draw(screen, gw.screenWidth, gw.screenHeight)
		}

		// Draw death screen on top
		if gw.manager.DeathScreen != nil {
			gw.manager.DeathScreen.Draw(screen)
		}

		// Draw profiler overlay
		gw.manager.Profiler.Draw(screen)

	case ui.StateMenu:
		// Draw menu
		// TODO: Implement menu drawing
	}
}
*/
