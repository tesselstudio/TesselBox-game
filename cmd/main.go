package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"time"

	"tesselbox/pkg/blocks"
	"tesselbox/pkg/crafting"
	"tesselbox/pkg/gametime"
	"tesselbox/pkg/hexagon"
	"tesselbox/pkg/items"
	"tesselbox/pkg/menu"
	"tesselbox/pkg/player"
	"tesselbox/pkg/save"
	"tesselbox/pkg/weather"
	"tesselbox/pkg/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// stringToBlockType converts a string block type name to blocks.BlockType
func stringToBlockType(blockTypeStr string) blocks.BlockType {
	blockMap := map[string]blocks.BlockType{
		"air":         blocks.AIR,
		"dirt":        blocks.DIRT,
		"grass":       blocks.GRASS,
		"stone":       blocks.STONE,
		"sand":        blocks.SAND,
		"water":       blocks.WATER,
		"log":         blocks.LOG,
		"leaves":      blocks.LEAVES,
		"coal_ore":    blocks.COAL_ORE,
		"iron_ore":    blocks.IRON_ORE,
		"gold_ore":    blocks.GOLD_ORE,
		"diamond_ore": blocks.DIAMOND_ORE,
		"bedrock":     blocks.BEDROCK,
		"glass":       blocks.GLASS,
		"brick":       blocks.BRICK,
		"plank":       blocks.PLANK,
		"cactus":      blocks.CACTUS,
	}
	if bt, ok := blockMap[blockTypeStr]; ok {
		return bt
	}
	return blocks.AIR
}

// minFloat32 returns the minimum of two float32 values
func minFloat32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
	FPS          = 60
)

// Game represents the game state
type Game struct {
	// Game state
	world     *world.World
	player    *player.Player
	inventory *items.Inventory

	// Systems
	craftingSystem *crafting.CraftingSystem
	craftingUI     *crafting.CraftingUI
	menu           *menu.Menu

	// Save system
	saveManager *save.SaveManager
	autoSaver   *save.AutoSaver

	// Day/night cycle
	dayNightCycle *gametime.DayNightCycle

	// Weather system
	weatherSystem *weather.WeatherSystem

	// Camera
	cameraX, cameraY float64

	// Mouse
	mouseX, mouseY int

	// Game state flags
	inMenu     bool
	inGame     bool
	inCrafting bool

	// Timing
	lastTime time.Time

	// For solid color drawing
	whiteImage *ebiten.Image
}

// NewGame creates a new game
func NewGame() *Game {
	// Create a 1x1 white image for solid color drawing
	whiteImage := ebiten.NewImage(1, 1)
	whiteImage.Fill(color.RGBA{255, 255, 255, 255})

	g := &Game{
		world:      world.NewWorld("default"), // Default world name
		player:     player.NewPlayer(0, 0),
		inventory:  items.NewInventory(32),
		cameraX:    0,
		cameraY:    0,
		lastTime:   time.Now(),
		whiteImage: whiteImage,
	}

	// Initialize crafting system
	g.craftingSystem = crafting.NewCraftingSystem()
	if err := g.craftingSystem.LoadRecipes("config/crafting_recipes.yaml"); err != nil {
		log.Printf("Warning: Failed to load crafting recipes: %v", err)
	}
	g.craftingUI = crafting.NewCraftingUI(g.craftingSystem, g.inventory)

	// Load items
	items.LoadItems()

	// Initialize menu
	g.menu = menu.NewMenu()

	// Set default hotbar items
	defaultItems := items.DefaultHotbarItems()
	for i, item := range defaultItems {
		if i < len(g.inventory.Slots) {
			g.inventory.Slots[i] = item
		}
	}

	// Initialize save system
	playerName := "player"
	g.saveManager = save.NewSaveManager("default", playerName)

	// Initialize auto-saver with 5-minute interval
	g.autoSaver = save.NewAutoSaver(g.saveManager, g.createSaveState(), 5*time.Minute)

	// Initialize day/night cycle (10 minute days for gameplay)
	g.dayNightCycle = gametime.NewDayNightCycle(600.0)

	// Initialize weather system
	g.weatherSystem = weather.NewWeatherSystem()

	// Start in menu
	g.inMenu = true
	g.inGame = false

	return g
}

// Update updates the game state
func (g *Game) Update() error {
	// Calculate delta time for framerate-independent movement
	currentTime := time.Now()
	deltaTime := currentTime.Sub(g.lastTime).Seconds()
	g.lastTime = currentTime

	// Handle menu state
	if g.inMenu {
		action := g.menu.Update()
		g.handleMenuAction(action)
		return nil
	}

	// Handle crafting UI
	if g.inCrafting {
		if err := g.craftingUI.Update(); err != nil {
			return err
		}

		// Handle escape to close crafting
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.inCrafting = false
		}
		return nil
	}

	// Handle game state
	if g.inGame {
		g.handleGameInput()

		// Update player with delta time (framerate-independent)
		g.player.Update(deltaTime)

		// Update mining progress
		g.updateMining(deltaTime)

		// Update day/night cycle
		g.dayNightCycle.Update()

		// Update weather system
		g.weatherSystem.Update(deltaTime, ScreenWidth, ScreenHeight)

		// Apply collision-aware position update
		nearbyHexagons := g.world.GetNearbyHexagons(g.player.X, g.player.Y, 300)
		g.player.UpdateWithCollision(deltaTime, func(minX, minY, maxX, maxY float64) bool {
			for _, hex := range nearbyHexagons {
				def := blocks.BlockDefinitions[getBlockKeyFromType(hex.BlockType)]
				if def == nil || !def.Solid {
					continue
				}

				hexMinX := hex.X - hex.Size
				hexMinY := hex.Y - hex.Size
				hexMaxX := hex.X + hex.Size
				hexMaxY := hex.Y + hex.Size

				if !(maxX < hexMinX || minX > hexMaxX || maxY < hexMinY || minY > hexMaxY) {
					return true
				}
			}
			return false
		})

		// Update camera to follow player
		centerX, centerY := g.player.GetCenter()
		g.cameraX = centerX - ScreenWidth/2
		g.cameraY = centerY - ScreenHeight/2

		// Generate chunks around player
		g.world.GetChunksInRange(centerX, centerY)

		// Unload distant chunks to manage memory
		g.world.UnloadDistantChunks(centerX, centerY)
	}

	return nil
}

// handleMenuAction handles menu actions
func (g *Game) handleMenuAction(action menu.MenuAction) {
	switch action {
	case menu.ActionStartGame:
		g.inMenu = false
		g.inGame = true
		g.player.SetPosition(0, 0)
		g.player.SetVelocity(0, 0)

	case menu.ActionOpenSettings:
		g.menu.SetSettingsMenu()

	case menu.ActionBack:
		if g.menu.CurrentMenu == menu.MenuSettings {
			g.menu.SetMainMenu()
		} else {
			g.menu.SetMainMenu()
		}

	case menu.ActionExit:
		os.Exit(0)
		return
	}
}

// handleGameInput handles game input
func (g *Game) handleGameInput() {
	// Handle keyboard input
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.MovingLeft = true
		g.player.MovingRight = false
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.MovingRight = true
		g.player.MovingLeft = false
	} else {
		g.player.MovingLeft = false
		g.player.MovingRight = false
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Jump()
	}

	// Open crafting menu
	if inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.inCrafting = true
		g.craftingUI.Toggle()
	}

	// Interact with crafting stations
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.interactWithStation()
	}

	// Drop item (Q key)
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.dropItem()
	}

	// Load game (F9)
	if inpututil.IsKeyJustPressed(ebiten.KeyF9) {
		if err := g.LoadGame(); err != nil {
			log.Printf("Failed to load game: %v", err)
		} else {
			log.Println("Game loaded successfully")
		}
	}

	// Return to menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.inGame = false
		g.inMenu = true
		g.menu.SetMainMenu()
	}

	// Handle mouse input
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	// Start mining when left mouse button is pressed
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.startMining()
	}

	// Attack with space bar (when not jumping)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.player.IsOnGround() {
		g.player.Attack()
	}

	// Block placement with right click (only when not attacking)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) && !g.player.IsAttacking() {
		g.handleBlockPlacement()
	}

	// Hotbar selection
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.inventory.SelectSlot(0)
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.inventory.SelectSlot(1)
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.inventory.SelectSlot(2)
	} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.inventory.SelectSlot(3)
	} else if inpututil.IsKeyJustPressed(ebiten.Key5) {
		g.inventory.SelectSlot(4)
	} else if inpututil.IsKeyJustPressed(ebiten.Key6) {
		g.inventory.SelectSlot(5)
	} else if inpututil.IsKeyJustPressed(ebiten.Key7) {
		g.inventory.SelectSlot(6)
	} else if inpututil.IsKeyJustPressed(ebiten.Key8) {
		g.inventory.SelectSlot(7)
	} else if inpututil.IsKeyJustPressed(ebiten.Key9) {
		g.inventory.SelectSlot(8)
	}

	// Scroll wheel for hotbar
	_, scrollY := ebiten.Wheel()
	if scrollY > 0 {
		g.inventory.PrevSlot()
	} else if scrollY < 0 {
		g.inventory.NextSlot()
	}
}

// dropItem drops the currently selected item from inventory
func (g *Game) dropItem() {
	selectedItem := g.inventory.GetSelectedItem()
	if selectedItem == nil || selectedItem.Type == items.NONE {
		return // Nothing to drop
	}

	// Remove one item from the selected slot
	if g.inventory.RemoveItem(1) {
		log.Printf("Dropped %v", selectedItem.Type)
		// TODO: Implement actual item dropping in world (create item entity)
	}
}

// interactWithStation checks for nearby crafting stations and opens the crafting UI
func (g *Game) interactWithStation() {
	playerX, playerY := g.player.GetCenter()
	interactionRange := 100.0 // pixels

	// Check nearby hexagons for crafting stations
	nearbyHexagons := g.world.GetNearbyHexagons(playerX, playerY, interactionRange)

	for _, hex := range nearbyHexagons {
		if hex.BlockType == blocks.AIR {
			continue
		}

		blockKey := getBlockKeyFromType(hex.BlockType)

		var station crafting.CraftingStation
		switch blockKey {
		case "workbench":
			station = crafting.STATION_WORKBENCH
		case "furnace":
			station = crafting.STATION_FURNACE
		case "anvil":
			station = crafting.STATION_ANVIL
		default:
			continue
		}

		// Found a station, set it and open UI
		g.craftingUI.SetStation(station)
		g.inCrafting = true
		g.craftingUI.Toggle()
		return // Only interact with the first found station
	}

	// No station found, open general crafting
	g.craftingUI.SetStation(crafting.STATION_NONE)
	g.inCrafting = true
	g.craftingUI.Toggle()
}

// startMining starts mining the block under the mouse cursor
func (g *Game) startMining() {
	// Convert mouse position to world coordinates
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY

	// Find the hexagon at mouse position
	targetHex := g.world.GetHexagonAt(mouseWorldX, mouseWorldY)
	if targetHex == nil {
		return
	}

	// Check if player can reach the block
	if !g.player.CanReach(targetHex.X, targetHex.Y) {
		return
	}

	// Start mining the target block
	g.player.StartMining(targetHex)
}

// updateMining updates mining progress and handles block destruction
func (g *Game) updateMining(deltaTime float64) {
	if !g.player.IsMining() {
		return
	}

	targetHex := g.player.GetMiningTarget()
	if targetHex == nil {
		g.player.StopMining()
		return
	}

	// Calculate mining speed based on tool and block
	miningSpeed := g.calculateMiningSpeed(targetHex.BlockType)

	// Update mining progress
	progressIncrease := miningSpeed * deltaTime
	g.player.MiningProgress += progressIncrease

	// Check if mining is complete
	if g.player.MiningProgress >= 100.0 {
		// Mining complete - destroy the block
		g.completeMining(targetHex)
		g.player.StopMining()
	}
}

// completeMining handles the completion of mining (block destruction and item pickup)
func (g *Game) completeMining(targetHex *world.Hexagon) {
	// Get the block type before removing
	blockType := targetHex.BlockType

	// Get the exact world position before removing
	x, y := targetHex.X, targetHex.Y
	g.world.RemoveHexagonAt(x, y)

	// Use item durability
	g.inventory.UseItem()

	// Add mined item to inventory
	minedItemType := g.getItemFromBlockType(blockType)
	if minedItemType != items.NONE {
		if !g.inventory.AddItem(minedItemType, 1) {
			// Inventory full - could implement dropping item here
			log.Printf("Inventory full, cannot pick up %v", minedItemType)
		}
	}
}

// handleMining handles block mining
func (g *Game) handleMining() {
	// Convert mouse position to world coordinates
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY

	// Find the hexagon at mouse position directly using world coordinates
	targetHex := g.world.GetHexagonAt(mouseWorldX, mouseWorldY)
	if targetHex == nil {
		return
	}

	// Check if player can reach the block
	if !g.player.CanReach(targetHex.X, targetHex.Y) {
		return
	}

	// Calculate mining damage based on tool and block hardness
	damage := g.calculateMiningDamage(targetHex.BlockType)

	// Damage the block
	targetHex.TakeDamage(damage)

	// Check if block is destroyed
	if targetHex.Health <= 0 {
		// Get the block type before removing
		blockType := targetHex.BlockType

		// Get the exact world position before removing
		x, y := targetHex.X, targetHex.Y
		g.world.RemoveHexagonAt(x, y)

		// Use item durability
		g.inventory.UseItem()

		// Add mined item to inventory
		minedItemType := g.getItemFromBlockType(blockType)
		if minedItemType != items.NONE {
			if !g.inventory.AddItem(minedItemType, 1) {
				// Inventory full - could implement dropping item here
				log.Printf("Inventory full, cannot pick up %v", minedItemType)
			}
		}
	}
}

// handleBlockPlacement handles block placement
func (g *Game) handleBlockPlacement() {
	// Get selected item
	selectedItem := g.inventory.GetSelectedItem()
	if selectedItem == nil || selectedItem.Type == items.NONE {
		return
	}

	// Get item properties
	props := items.GetItemProperties(selectedItem.Type)
	if props == nil || !props.IsPlaceable {
		return
	}

	// Convert mouse position to world coordinates
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY

	// Use the correct world position for placement
	centerX, centerY, _, _ := world.PixelToHexCenter(mouseWorldX, mouseWorldY)

	// Placement validation: check if position is valid
	if !g.canPlaceBlockAt(centerX, centerY) {
		return // Cannot place block here
	}

	// Place block at the calculated position
	blockType := stringToBlockType(props.BlockType)
	g.world.AddHexagonAt(centerX, centerY, blockType)

	// Remove item from inventory
	g.inventory.RemoveItem(1)
}

// canPlaceBlockAt checks if a block can be placed at the given position
func (g *Game) canPlaceBlockAt(x, y float64) bool {
	// Check if there's already a block at this position
	existingHex := g.world.GetHexagonAt(x, y)
	if existingHex != nil {
		return false // Cannot place on existing block
	}

	// Check if player is too close (prevent placing blocks inside player)
	playerCenterX, playerCenterY := g.player.GetCenter()
	distance := math.Sqrt((x-playerCenterX)*(x-playerCenterX) + (y-playerCenterY)*(y-playerCenterY))
	minDistance := g.player.Width + world.HexSize
	if distance < minDistance {
		return false // Too close to player
	}

	return true
}

// drawBlockPlacementPreview draws a preview of where a block will be placed
func (g *Game) drawBlockPlacementPreview(screen *ebiten.Image) {
	// Only show preview if player is holding a placeable item
	selectedItem := g.inventory.GetSelectedItem()
	if selectedItem == nil || selectedItem.Type == items.NONE {
		return
	}

	props := items.GetItemProperties(selectedItem.Type)
	if props == nil || !props.IsPlaceable {
		return
	}

	// Get mouse world position
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY

	// Calculate placement position
	centerX, centerY, _, _ := world.PixelToHexCenter(mouseWorldX, mouseWorldY)

	// Check if placement is valid
	if !g.canPlaceBlockAt(centerX, centerY) {
		return // Don't show preview for invalid positions
	}

	// Calculate screen position
	screenX := centerX - g.cameraX
	screenY := centerY - g.cameraY

	// Check if preview is on screen
	if screenX < -world.HexSize || screenX > ScreenWidth+world.HexSize ||
		screenY < -world.HexSize || screenY > ScreenHeight+world.HexSize {
		return
	}

	// Get block color for preview
	blockType := stringToBlockType(props.BlockType)
	blockKey := getBlockKeyFromType(blockType)
	blockDef := blocks.BlockDefinitions[blockKey]
	if blockDef == nil {
		return
	}

	// Draw translucent preview hexagon
	corners := hexagon.GetHexCorners(screenX, screenY, world.HexSize)
	vertices := make([]ebiten.Vertex, len(corners))
	for i, corner := range corners {
		vertices[i] = ebiten.Vertex{
			DstX:   float32(corner[0]),
			DstY:   float32(corner[1]),
			ColorR: float32(blockDef.Color.R) / 255.0 * 0.7, // More translucent
			ColorG: float32(blockDef.Color.G) / 255.0 * 0.7,
			ColorB: float32(blockDef.Color.B) / 255.0 * 0.7,
			ColorA: 0.5, // Semi-transparent
		}
	}

	indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5}
	screen.DrawTriangles(vertices, indices, g.whiteImage, nil)
}

// drawMiningProgress draws the mining progress bar above the block being mined
func (g *Game) drawMiningProgress(screen *ebiten.Image) {
	if !g.player.IsMining() {
		return
	}

	targetHex := g.player.GetMiningTarget()
	if targetHex == nil {
		return
	}

	// Calculate screen position of the block
	screenX := targetHex.X - g.cameraX
	screenY := targetHex.Y - g.cameraY

	// Check if block is on screen
	if screenX < -world.HexSize || screenX > ScreenWidth+world.HexSize ||
		screenY < -world.HexSize || screenY > ScreenHeight+world.HexSize {
		return
	}

	// Progress bar dimensions and position
	barWidth := 60.0
	barHeight := 8.0
	barX := screenX - barWidth/2
	barY := screenY - world.HexSize - 15 // Above the hexagon

	// Background (gray)
	ebitenutil.DrawRect(screen, barX, barY, barWidth, barHeight, g.colorToRGB(100, 100, 100))

	// Progress fill (green)
	progressRatio := g.player.MiningProgress / 100.0
	if progressRatio > 1.0 {
		progressRatio = 1.0
	}
	fillWidth := barWidth * progressRatio
	ebitenutil.DrawRect(screen, barX, barY, fillWidth, barHeight, g.colorToRGB(100, 255, 100))

	// Border
	ebitenutil.DrawRect(screen, barX, barY, barWidth, 1, g.colorToRGB(0, 0, 0))
	ebitenutil.DrawRect(screen, barX, barY+barHeight-1, barWidth, 1, g.colorToRGB(0, 0, 0))
	ebitenutil.DrawRect(screen, barX, barY, 1, barHeight, g.colorToRGB(0, 0, 0))
	ebitenutil.DrawRect(screen, barX+barWidth-1, barY, 1, barHeight, g.colorToRGB(0, 0, 0))
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.inMenu {
		g.menu.Draw(screen)
		return
	}

	if g.inCrafting {
		// Draw game in background
		g.drawGameScene(screen)
		// Draw crafting UI overlay
		g.craftingUI.Draw(screen)
		return
	}

	if g.inGame {
		g.drawGameScene(screen)
		g.drawUI(screen)
		g.drawDebugInfo(screen)
	}
}

// drawGameScene draws the game world and player
func (g *Game) drawGameScene(screen *ebiten.Image) {
	// Draw background with day/night cycle sky color
	skyR, skyG, skyB := g.dayNightCycle.GetSkyColor()
	skyColor := g.colorToRGB(int(skyR*255), int(skyG*255), int(skyB*255))
	screen.Fill(skyColor)

	// Get visible blocks
	px, py := g.player.GetCenter()
	visibleBlocks := g.world.GetVisibleBlocks(px, py)

	// Draw blocks with batching by color for performance
	g.drawBlocksBatched(screen, visibleBlocks)

	// Draw player
	g.drawPlayer(screen)

	// Draw block placement preview
	g.drawBlockPlacementPreview(screen)

	// Draw mining progress bar
	g.drawMiningProgress(screen)

	// Draw weather particles
	g.weatherSystem.Draw(screen, g.cameraX, g.cameraY)
}

// drawBlocksBatched draws blocks in batches grouped by color for performance optimization
func (g *Game) drawBlocksBatched(screen *ebiten.Image, blockList []*world.Hexagon) {
	// Group blocks by color and damage state for batching
	colorGroups := make(map[string][]*world.Hexagon)

	for _, block := range blockList {
		if block.BlockType == blocks.AIR {
			continue
		}

		props := blocks.BlockDefinitions[getBlockKeyFromType(block.BlockType)]
		if props == nil {
			continue
		}

		// Create a key based on color and damage state for grouping
		damageState := "normal"
		if block.Health < block.MaxHealth {
			damageRatio := block.Health / block.MaxHealth
			if damageRatio > 0.5 {
				damageState = "minor_damage"
			} else {
				damageState = "major_damage"
			}
		}

		colorKey := fmt.Sprintf("%d_%d_%d_%s", props.Color.R, props.Color.G, props.Color.B, damageState)
		colorGroups[colorKey] = append(colorGroups[colorKey], block)
	}

	// Draw each color group as a batch
	for _, groupBlocks := range colorGroups {
		if len(groupBlocks) == 0 {
			continue
		}

		// Get properties from first block in group
		firstBlock := groupBlocks[0]
		props := blocks.BlockDefinitions[getBlockKeyFromType(firstBlock.BlockType)]
		if props == nil {
			continue
		}

		// Prepare vertices for all blocks in this color group
		totalVertices := len(groupBlocks) * 6 // 6 vertices per hexagon
		vertices := make([]ebiten.Vertex, 0, totalVertices)
		indices := make([]uint16, 0, len(groupBlocks)*6) // 6 indices per hexagon

		baseIndex := uint16(0)

		for _, block := range groupBlocks {
			// Calculate screen position
			screenX := block.X - g.cameraX
			screenY := block.Y - g.cameraY

			// Check if block is on screen
			if screenX < -100 || screenX > ScreenWidth+100 ||
				screenY < -100 || screenY > ScreenHeight+100 {
				continue
			}

			// Get hexagon corners
			corners := hexagon.GetHexCorners(screenX, screenY, world.HexSize)

			// Calculate color with damage darkening
			r := float32(props.Color.R) / 255.0
			gc := float32(props.Color.G) / 255.0
			b := float32(props.Color.B) / 255.0
			a := float32(props.Color.A) / 255.0

			// Apply damage darkening
			if block.Health < block.MaxHealth {
				damageRatio := block.Health / block.MaxHealth
				r *= float32(math.Min(1.0, float64(damageRatio)))
				gc *= float32(math.Min(1.0, float64(damageRatio)))
				b *= float32(math.Min(1.0, float64(damageRatio)))
			}

			// Add vertices for this hexagon
			for _, corner := range corners {
				vertices = append(vertices, ebiten.Vertex{
					DstX:   float32(corner[0]),
					DstY:   float32(corner[1]),
					ColorR: r,
					ColorG: gc,
					ColorB: b,
					ColorA: a,
				})
			}

			// Add indices for this hexagon
			hexIndices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5}
			for _, idx := range hexIndices {
				indices = append(indices, baseIndex+idx)
			}

			baseIndex += 6
		}

		// Draw the entire batch
		if len(vertices) > 0 {
			screen.DrawTriangles(vertices, indices, g.whiteImage, nil)
		}
	}
}

// drawUI draws the user interface
func (g *Game) drawUI(screen *ebiten.Image) {
	// Draw hotbar
	hotbarWidth := 400
	hotbarHeight := 50
	hotbarX := (ScreenWidth - hotbarWidth) / 2
	hotbarY := ScreenHeight - hotbarHeight - 20

	slotWidth := hotbarWidth / 8

	for i := 0; i < 8; i++ {
		slotX := hotbarX + i*slotWidth
		slotY := hotbarY

		// Draw slot background
		bgColor := g.colorToRGB(100, 100, 100)
		if i == g.inventory.Selected {
			bgColor = g.colorToRGB(150, 150, 150)
		}

		ebitenutil.DrawRect(screen, float64(slotX), float64(slotY), float64(slotWidth-2), float64(hotbarHeight), bgColor)

		// Draw slot border
		ebitenutil.DrawRect(screen, float64(slotX), float64(slotY), float64(slotWidth-2), 2, g.colorToRGB(0, 0, 0))
		ebitenutil.DrawRect(screen, float64(slotX), float64(slotY+hotbarHeight-2), float64(slotWidth-2), 2, g.colorToRGB(0, 0, 0))
		ebitenutil.DrawRect(screen, float64(slotX), float64(slotY), 2, float64(hotbarHeight), g.colorToRGB(0, 0, 0))
		ebitenutil.DrawRect(screen, float64(slotX+slotWidth-4), float64(slotY), 2, float64(hotbarHeight), g.colorToRGB(0, 0, 0))

		// Draw item if present
		if i < len(g.inventory.Slots) {
			item := g.inventory.Slots[i]
			if item.Type != items.NONE {
				itemColor := items.ItemColorByID(item.Type)
				ebitenutil.DrawRect(screen, float64(slotX+10), float64(slotY+10), float64(slotWidth-20), float64(hotbarHeight-20),
					g.colorToRGB(int(itemColor.R), int(itemColor.G), int(itemColor.B)))

				// Draw quantity indicator
				if item.Quantity > 1 {
					ebitenutil.DrawRect(screen, float64(slotX+slotWidth-20), float64(slotY+hotbarHeight-20), 8, 8, g.colorToRGB(255, 255, 255))
				}
			}
		}
	}

	// Draw health bar
	healthBarWidth := 200
	healthBarHeight := 20
	healthBarX := 20
	healthBarY := 20

	healthRatio := g.player.Health / g.player.MaxHealth

	// Background
	ebitenutil.DrawRect(screen, float64(healthBarX), float64(healthBarY), float64(healthBarWidth), float64(healthBarHeight), g.colorToRGB(50, 50, 50))

	// Health
	ebitenutil.DrawRect(screen, float64(healthBarX), float64(healthBarY), float64(int(float64(healthBarWidth)*healthRatio)), float64(healthBarHeight), g.colorToRGB(200, 50, 50))

	// Border
	ebitenutil.DrawRect(screen, float64(healthBarX), float64(healthBarY), float64(healthBarWidth), 2, g.colorToRGB(0, 0, 0))
	ebitenutil.DrawRect(screen, float64(healthBarX), float64(healthBarY+healthBarHeight-2), float64(healthBarWidth), 2, g.colorToRGB(0, 0, 0))
	ebitenutil.DrawRect(screen, float64(healthBarX), float64(healthBarY), 2, float64(healthBarHeight), g.colorToRGB(0, 0, 0))
	ebitenutil.DrawRect(screen, float64(healthBarX+healthBarWidth-2), float64(healthBarY), 2, float64(healthBarHeight), g.colorToRGB(0, 0, 0))
}

// drawDebugInfo draws debug information
func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	px, py := g.player.GetCenter()
	vx, vy := g.player.GetVelocity()

	timeInfo := g.dayNightCycle.GetTimeString()
	lightInfo := fmt.Sprintf("Ambient: %.2f, Sky: %.2f, Block: %.2f",
		g.dayNightCycle.AmbientLight, g.dayNightCycle.SkyLight, g.dayNightCycle.BlockLight)

	_, weatherIntensity, weatherName := g.weatherSystem.GetWeatherInfo()
	weatherInfo := fmt.Sprintf("Weather: %s (%.1f)", weatherName, weatherIntensity)

	info := fmt.Sprintf("Pos: (%.1f, %.1f)\nVel: (%.1f, %.1f)\nFPS: %.1f\nOnGround: %v\nDelta: %.4f\n%s\n%s\n%s",
		px, py, vx, vy, ebiten.ActualFPS(), g.player.IsOnGround(), time.Since(g.lastTime).Seconds(),
		timeInfo, lightInfo, weatherInfo)

	ebitenutil.DebugPrint(screen, info)
}

// drawBlock draws a hexagonal block
func (g *Game) drawBlock(screen *ebiten.Image, block *world.Hexagon) {
	if block.BlockType == blocks.AIR {
		return
	}

	props := blocks.BlockDefinitions[getBlockKeyFromType(block.BlockType)]
	if props == nil {
		return
	}

	// Calculate screen position
	screenX := block.X - g.cameraX
	screenY := block.Y - g.cameraY

	// Check if block is on screen
	if screenX < -100 || screenX > ScreenWidth+100 ||
		screenY < -100 || screenY > ScreenHeight+100 {
		return
	}

	// Get hexagon corners
	corners := hexagon.GetHexCorners(screenX, screenY, world.HexSize)

	// Create polygon vertices
	vertices := make([]ebiten.Vertex, len(corners))
	for i, corner := range corners {
		vertices[i] = ebiten.Vertex{
			DstX:   float32(corner[0]),
			DstY:   float32(corner[1]),
			ColorR: float32(props.Color.R) / 255.0,
			ColorG: float32(props.Color.G) / 255.0,
			ColorB: float32(props.Color.B) / 255.0,
			ColorA: float32(props.Color.A) / 255.0,
		}
	}

	// Draw filled hexagon
	indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5}

	// Check if mouse is hovering over this block
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY

	q, r := hexagon.PixelToHex(mouseWorldX, mouseWorldY, world.HexSize)
	hoverHex := hexagon.HexRound(q, r)

	// Get block's hex coordinates
	blockQ, blockR := hexagon.PixelToHex(block.X, block.Y, world.HexSize)
	blockHex := hexagon.HexRound(blockQ, blockR)

	if hoverHex.Q == blockHex.Q && hoverHex.R == blockHex.R {
		// Highlight hovered block
		for i := range vertices {
			vertices[i].ColorR = minFloat32(1.0, vertices[i].ColorR+0.2)
			vertices[i].ColorG = minFloat32(1.0, vertices[i].ColorG+0.2)
			vertices[i].ColorB = minFloat32(1.0, vertices[i].ColorB+0.2)
		}
	}

	// Apply damage darkening
	if block.Health < block.MaxHealth {
		damageRatio := block.Health / block.MaxHealth
		for i := range vertices {
			vertices[i].ColorR *= float32(math.Min(1.0, float64(damageRatio)))
			vertices[i].ColorG *= float32(math.Min(1.0, float64(damageRatio)))
			vertices[i].ColorB *= float32(math.Min(1.0, float64(damageRatio)))
		}
	}

	screen.DrawTriangles(vertices, indices, g.whiteImage, nil)
}

// drawPlayer draws the player
func (g *Game) drawPlayer(screen *ebiten.Image) {
	screenX := g.player.X - g.cameraX
	screenY := g.player.Y - g.cameraY

	// Draw player body
	bodyColor := g.colorToRGB(255, 100, 100)
	if g.player.IsAttacking() {
		// Flash red during attack
		progress := g.player.GetAttackProgress()
		if progress < 0.5 {
			// Bright red during wind-up
			bodyColor = g.colorToRGB(255, 150, 150)
		} else {
			// Dark red during attack
			bodyColor = g.colorToRGB(200, 50, 50)
		}
	}
	ebitenutil.DrawRect(screen, screenX, screenY, float64(g.player.Width), float64(g.player.Height), bodyColor)

	// Draw player head (simple representation)
	headColor := g.colorToRGB(255, 200, 150)
	if g.player.IsAttacking() {
		// Angry expression during attack
		headColor = g.colorToRGB(255, 150, 100)
	}
	ebitenutil.DrawRect(screen, screenX+10, screenY+5, 20, 20, headColor)

	// Draw attack effect
	if g.player.IsAttacking() {
		g.drawAttackEffect(screen, screenX, screenY)
	}
}

// drawAttackEffect draws visual effects during attacks
func (g *Game) drawAttackEffect(screen *ebiten.Image, playerX, playerY float64) {
	progress := g.player.GetAttackProgress()

	// Draw attack arc/sweep effect
	if progress > 0.3 && progress < 0.8 {
		attackRange := g.player.GetAttackRange()

		// Simple attack arc visualization
		for angle := -0.5; angle <= 0.5; angle += 0.1 {
			endX := playerX + g.player.Width/2 + math.Cos(angle)*attackRange
			endY := playerY + g.player.Height/2 + math.Sin(angle)*attackRange
			ebitenutil.DrawLine(screen, playerX+g.player.Width/2, playerY+g.player.Height/2, endX, endY, g.colorToRGB(255, 255, 100))
		}
	}
}

// Layout defines the game's layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// colorToRGB converts RGB values to a color
func (g *Game) colorToRGB(rVal, gVal, bVal int) color.RGBA {
	return color.RGBA{
		R: uint8(rVal),
		G: uint8(gVal),
		B: uint8(bVal),
		A: 255,
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// getBlockKeyFromType converts a BlockType to its string key
func getBlockKeyFromType(blockType blocks.BlockType) string {
	switch blockType {
	case blocks.AIR:
		return "air"
	case blocks.DIRT:
		return "dirt"
	case blocks.GRASS:
		return "grass"
	case blocks.STONE:
		return "stone"
	case blocks.SAND:
		return "sand"
	case blocks.WATER:
		return "water"
	case blocks.LOG:
		return "log"
	case blocks.LEAVES:
		return "leaves"
	case blocks.COAL_ORE:
		return "coal_ore"
	case blocks.IRON_ORE:
		return "iron_ore"
	case blocks.GOLD_ORE:
		return "gold_ore"
	case blocks.DIAMOND_ORE:
		return "diamond_ore"
	case blocks.BEDROCK:
		return "bedrock"
	case blocks.GLASS:
		return "glass"
	case blocks.BRICK:
		return "brick"
	case blocks.PLANK:
		return "plank"
	case blocks.CACTUS:
		return "cactus"
	default:
		return "dirt"
	}
}

// calculateMiningSpeed calculates how fast mining progresses (progress per second)
func (g *Game) calculateMiningSpeed(blockType blocks.BlockType) float64 {
	// Get block hardness
	blockKey := getBlockKeyFromType(blockType)
	blockDef := blocks.BlockDefinitions[blockKey]
	if blockDef == nil {
		return 10.0 // Default speed
	}

	// Base mining speed (progress per second for hand mining)
	baseSpeed := 5.0

	// Get tool properties
	selectedItem := g.inventory.GetSelectedItem()
	if selectedItem != nil && selectedItem.Type != items.NONE {
		itemProps := items.GetItemProperties(selectedItem.Type)
		if itemProps != nil && itemProps.IsTool {
			// Tool speed = base speed * tool power / sqrt(block hardness)
			// This makes tools much more effective against harder blocks
			baseSpeed = baseSpeed * itemProps.ToolPower / math.Sqrt(math.Max(0.1, blockDef.Hardness))
		}
	}

	// Ensure minimum mining speed
	if baseSpeed < 1.0 {
		baseSpeed = 1.0
	}

	// Cap maximum mining speed to prevent instant mining
	if baseSpeed > 100.0 {
		baseSpeed = 100.0
	}

	return baseSpeed
}

// calculateMiningDamage calculates the damage dealt to a block per mining tick
func (g *Game) calculateMiningDamage(blockType blocks.BlockType) float64 {
	// Get block hardness
	blockKey := getBlockKeyFromType(blockType)
	blockDef := blocks.BlockDefinitions[blockKey]
	if blockDef == nil {
		return 1.0 // Default damage
	}

	// Base damage (damage per tick for hand mining)
	baseDamage := 1.0

	// Get tool properties
	selectedItem := g.inventory.GetSelectedItem()
	if selectedItem != nil && selectedItem.Type != items.NONE {
		itemProps := items.GetItemProperties(selectedItem.Type)
		if itemProps != nil && itemProps.IsTool {
			// Tool damage = base damage * tool power / block hardness
			// This makes tools much more effective against harder blocks
			baseDamage = baseDamage * itemProps.ToolPower / math.Max(0.1, blockDef.Hardness)
		}
	}

	// Ensure minimum damage
	if baseDamage < 0.1 {
		baseDamage = 0.1
	}

	// Cap maximum damage to prevent instant destruction
	if baseDamage > 10.0 {
		baseDamage = 10.0
	}

	return baseDamage
}

// getItemFromBlockType converts a block type to the corresponding item type
func (g *Game) getItemFromBlockType(blockType blocks.BlockType) items.ItemType {
	switch blockType {
	case blocks.DIRT:
		return items.DIRT_BLOCK
	case blocks.GRASS:
		return items.GRASS_BLOCK
	case blocks.STONE:
		return items.STONE_BLOCK
	case blocks.SAND:
		return items.SAND_BLOCK
	case blocks.LOG:
		return items.LOG_BLOCK
	case blocks.COAL_ORE:
		return items.COAL
	case blocks.IRON_ORE:
		return items.IRON_INGOT // Could be iron ore, but using ingot for now
	case blocks.GOLD_ORE:
		return items.GOLD_INGOT
	case blocks.DIAMOND_ORE:
		return items.DIAMOND
	default:
		return items.NONE
	}
}

// openURL opens a URL in the default browser
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// createSaveState creates a save state from the current game state
func (g *Game) createSaveState() *save.GameState {
	return &save.GameState{
		World:      g.world,
		Player:     g.player,
		Inventory:  g.inventory,
		CameraX:    g.cameraX,
		CameraY:    g.cameraY,
		InMenu:     g.inMenu,
		InGame:     g.inGame,
		InCrafting: g.inCrafting,
	}
}

// SaveGame saves the current game state
func (g *Game) SaveGame() error {
	if g.saveManager == nil {
		return fmt.Errorf("save manager not initialized")
	}
	return g.saveManager.SaveGame(g.createSaveState())
}

// LoadGame loads a game state
func (g *Game) LoadGame() error {
	if g.saveManager == nil {
		return fmt.Errorf("save manager not initialized")
	}

	saveData, err := g.saveManager.LoadGame()
	if err != nil {
		return err
	}

	return g.saveManager.ApplySaveData(saveData, g.createSaveState())
}

// StartAutoSave starts the auto-saver
func (g *Game) StartAutoSave() {
	if g.autoSaver != nil {
		g.autoSaver.Start()
	}
}

// StopAutoSave stops the auto-saver
func (g *Game) StopAutoSave() {
	if g.autoSaver != nil {
		g.autoSaver.Stop()
	}
}

// Main function
func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Tesselbox v2.0 - Hexagon Sandbox")
	ebiten.SetTPS(FPS)

	game := NewGame()

	// Start auto-saver
	game.StartAutoSave()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
