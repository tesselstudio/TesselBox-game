package main

import (
	"fmt"
	"image/color"

	"tesselbox/pkg/hexagon"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Constants for rendering
const (
	CullingPadding    = 100.0
	HotbarSlotSize    = 50
	HotbarSlotCount   = 9
	HotbarYOffset     = 10
	HotbarBackgroundColor = 0x333333
	HotbarSelectedColor   = 0x666666
	ItemNameMaxLength = 20
	DroppedItemCulling = 2000
	ItemSize          = 32
)

// renderWorld draws the game world
func (g *Game) renderWorld(screen *ebiten.Image) {
	// Calculate visible region with padding
	screenWidth, screenHeight := ebiten.WindowSize()
	minX := g.cameraX - CullingPadding
	minY := g.cameraY - CullingPadding
	maxX := g.cameraX + float64(screenWidth) + CullingPadding
	maxY := g.cameraY + float64(screenHeight) + CullingPadding

	// Convert to grid coordinates
	startGridX, startGridY := g.world.PixelToGrid(minX, minY)
	endGridX, endGridY := g.world.PixelToGrid(maxX, maxY)

	// Render blocks
	for y := startGridY; y <= endGridY; y++ {
		for x := startGridX; x <= endGridX; x++ {
			block := g.world.GetBlock(x, y)
			if block != nil {
				px, py := g.world.GridToPixel(x, y)
				screenX := px - g.cameraX
				screenY := py - g.cameraY

				// Draw hexagon block
				hexagon.DrawHexagon(screen, screenX, screenY, block.Color)

				// Show mining progress if being mined
				progress := g.world.GetMiningProgress(x, y)
				if progress > 0 {
					g.renderMiningBar(screen, px, py, progress)
				}
			}
		}
	}
}

// renderEntities draws all game entities
func (g *Game) renderEntities(screen *ebiten.Image) {
	g.entityManager.Render(screen, g.cameraX, g.cameraY)
}

// renderPlayer draws the player character
func (g *Game) renderPlayer(screen *ebiten.Image) {
	playerX, playerY := g.player.GetPosition()
	screenX := playerX - g.cameraX
	screenY := playerY - g.cameraY

	// Draw player body
	bodySize := 30.0
	bodyImg := ebiten.NewImage(int(bodySize), int(bodySize*1.5))
	bodyImg.Fill(g.playerBodyColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenX-bodySize/2, screenY-bodySize*1.5/2)
	screen.DrawImage(bodyImg, op)

	// Draw player head
	headSize := 20.0
	headImg := ebiten.NewImage(int(headSize), int(headSize))
	headImg.Fill(g.playerHeadColor)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenX-headSize/2, screenY-bodySize*1.5/2-headSize)
	screen.DrawImage(headImg, op)
}

// renderDroppedItems draws dropped items in the world
func (g *Game) renderDroppedItems(screen *ebiten.Image) {
	for _, item := range g.droppedItems {
		// Culling - don't render if too far from player
		playerX, playerY := g.player.GetPosition()
		dx := item.X - playerX
		dy := item.Y - playerY
		if dx*dx+dy*dy > DroppedItemCulling*DroppedItemCulling {
			continue
		}

		screenX := item.X - g.cameraX
		screenY := item.Y - g.cameraY

		// Draw item representation
		itemImg := ebiten.NewImage(ItemSize, ItemSize)
		itemImg.Fill(color.RGBA{200, 200, 100, 255})

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(screenX-ItemSize/2, screenY-ItemSize/2)
		screen.DrawImage(itemImg, op)
	}
}

// renderUI draws the game UI
func (g *Game) renderUI(screen *ebiten.Image) {
	// Draw hotbar at bottom
	g.renderHotbar(screen)

	// Draw inventory if open
	if g.showInventory {
		g.renderInventory(screen)
	}

	// Draw crafting UI if open
	if g.craftingUI.IsOpen() {
		g.craftingUI.Draw(screen)
	}
}

// renderHotbar draws the hotbar at the bottom of the screen
func (g *Game) renderHotbar(screen *ebiten.Image) {
	screenWidth, screenHeight := ebiten.WindowSize()

	// Calculate hotbar position (centered at bottom)
	hotbarWidth := HotbarSlotSize*HotbarSlotCount + (HotbarSlotCount-1)*5
	startX := (screenWidth - hotbarWidth) / 2
	y := screenHeight - HotbarSlotSize - HotbarYOffset

	// Draw hotbar background
	hotbarBg := ebiten.NewImage(hotbarWidth+20, HotbarSlotSize+10)
	hotbarBg.Fill(color.RGBA{uint8(HotbarBackgroundColor >> 16), uint8(HotbarBackgroundColor >> 8), uint8(HotbarBackgroundColor), 200})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(startX-10), float64(y-5))
	screen.DrawImage(hotbarBg, op)

	// Draw slots
	inv := g.player.GetInventory()
	for i := 0; i < HotbarSlotCount && i < inv.GetHotbarSize(); i++ {
		x := startX + i*(HotbarSlotSize+5)

		// Draw slot background
		slotColor := color.RGBA{100, 100, 100, 255}
		if i == inv.GetSelectedSlot() {
			slotColor = color.RGBA{uint8(HotbarSelectedColor >> 16), uint8(HotbarSelectedColor >> 8), uint8(HotbarSelectedColor), 255}
		}

		slotImg := ebiten.NewImage(HotbarSlotSize, HotbarSlotSize)
		slotImg.Fill(slotColor)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(slotImg, op)

		// Draw item if present
		itemStack := inv.GetHotbarSlot(i)
		if itemStack != nil && itemStack.Count > 0 {
			// Draw item representation
			itemImg := ebiten.NewImage(HotbarSlotSize-10, HotbarSlotSize-10)
			itemImg.Fill(color.RGBA{200, 150, 100, 255})
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x+5), float64(y+5))
			screen.DrawImage(itemImg, op)

			// Draw count
			if itemStack.Count > 1 {
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", itemStack.Count), x+HotbarSlotSize-15, y+HotbarSlotSize-15)
			}
		}

		// Draw slot number
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i+1), x+2, y+2)
	}
}

// renderInventory draws the full inventory
func (g *Game) renderInventory(screen *ebiten.Image) {
	screenWidth, screenHeight := ebiten.WindowSize()

	// Draw inventory background
	invWidth := 400
	invHeight := 300
	x := (screenWidth - invWidth) / 2
	y := (screenHeight - invHeight) / 2

	bg := ebiten.NewImage(invWidth, invHeight)
	bg.Fill(color.RGBA{30, 30, 40, 230})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(bg, op)

	// Draw title
	ebitenutil.DebugPrintAt(screen, "INVENTORY", x+10, y+10)
}

// renderDebugInfo draws debug information
func (g *Game) renderDebugInfo(screen *ebiten.Image) {
	playerX, playerY := g.player.GetPosition()
	gridX, gridY := g.world.PixelToGrid(playerX, playerY)

	debugInfo := fmt.Sprintf(
		"FPS: %.0f | TPS: %.0f\n"+
		"Player: (%.1f, %.1f) [%d, %d]\n"+
		"Time: %s\n"+
		"Weather: %s\n"+
		"Creative: %v",
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		playerX, playerY,
		gridX, gridY,
		g.dayNightCycle.GetTimeString(),
		g.weatherSystem.GetCurrentWeather(),
		g.CreativeMode,
	)

	ebitenutil.DebugPrint(screen, debugInfo)
}
