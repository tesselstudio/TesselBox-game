package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"tesselbox/pkg/blocks"
	"tesselbox/pkg/crafting"
	"tesselbox/pkg/gametime"
	"tesselbox/pkg/hexagon"
	"tesselbox/pkg/input"
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

// DroppedItem represents an item that has been dropped in the world
type DroppedItem struct {
	Type     items.ItemType
	Quantity int
	X, Y     float64
	VX, VY   float64 // Velocity for physics
	Lifetime time.Time
}

// Game represents the game state
type Game struct {
	// Game state
	world         *world.World
	player        *player.Player
	inventory     *items.Inventory
	selectedBlock string // For creative mode block selection

	// Systems
	craftingSystem *crafting.CraftingSystem
	craftingUI     *crafting.CraftingUI
	menu           *menu.Menu
	pluginManager  interface{} // Will be *entities.EnhancedPluginManager
	inputManager   *input.InputManager

	// Save system
	saveManager *save.SaveManager
	autoSaver   *save.AutoSaver

	// Day/night cycle
	dayNightCycle *gametime.DayNightCycle

	// Weather system
	weatherSystem *weather.WeatherSystem

	// Dropped items
	droppedItems []*DroppedItem

	// Object pools for rendering optimization
	vertexPool [][]ebiten.Vertex
	indicesPool [][]uint16
	poolIndex int

	// Camera
	cameraX, cameraY float64

	// Mouse
	mouseX, mouseY   int
	hoveredBlockName string
	leftMouseWasPressed  bool
	rightMouseWasPressed bool

	// Game state flags
	inMenu       bool
	inGame       bool
	inCrafting   bool
	CreativeMode bool

	// Command system
	commandMode   bool
	commandString string

	// Timing
	lastTime     time.Time
	MiningDamage float64
	MiningSpeed  float64

	// For solid color drawing
	whiteImage *ebiten.Image
}

// NewGame creates a new game
func NewGame() *Game {
	// Create a 1x1 white image for solid color drawing
	whiteImage := ebiten.NewImage(1, 1)
	whiteImage.Fill(color.RGBA{255, 255, 255, 255})

	g := &Game{
		world:         world.NewWorld("default"), // Default world name
		player:        player.NewPlayer(0, 0), // Temporary position, will be updated
		inventory:     items.NewInventory(32),
		selectedBlock: "dirt", // Default to dirt in creative mode
		CreativeMode:  true,   // Enable creative mode by default
		cameraX:       0,
		cameraY:       0,
		lastTime:      time.Now(),
		whiteImage:    whiteImage,
		leftMouseWasPressed:  false,
		rightMouseWasPressed: false,
	}

	// Initialize object pools for rendering optimization
	g.vertexPool = make([][]ebiten.Vertex, 10)
	g.indicesPool = make([][]uint16, 10)
	for i := range g.vertexPool {
		g.vertexPool[i] = make([]ebiten.Vertex, 0, 1000) // Pre-allocate for up to 166 blocks
		g.indicesPool[i] = make([]uint16, 0, 1000)
	}
	g.poolIndex = 0

	// Find a suitable spawn position and set player position
	// Spawn in area where terrain is actually generated (negative coordinates)
	spawnX, spawnY := g.world.FindSpawnPosition(-2000, -2000)
	g.player.SetPosition(spawnX, spawnY)

	// Initialize crafting system
	g.craftingSystem = crafting.NewCraftingSystem()
	if err := g.craftingSystem.LoadRecipesFromAssets(); err != nil {
		log.Printf("Warning: Failed to load crafting recipes: %v", err)
	}
	g.craftingUI = crafting.NewCraftingUI(g.craftingSystem, g.inventory)

	// Initialize input manager
	g.inputManager = input.NewInputManager()

	// Load items
	items.LoadItems()

	// Load blocks
	blocks.LoadBlocks()

	// Add some initial items to inventory for testing
	g.inventory.AddItem(items.DIRT_BLOCK, 64)
	g.inventory.AddItem(items.STONE_BLOCK, 64)
	g.inventory.AddItem(items.GRASS_BLOCK, 64)

	// Initialize menu
	g.menu = menu.NewMenu()
	g.menu.CreativeMode = g.CreativeMode // Set creative mode in menu

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
	// g.inMenu = true
	// g.inGame = false
	g.inMenu = false
	g.inGame = true

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

		// Update dropped items physics
		g.updateDroppedItems(deltaTime)

		// Generate chunks around player first to ensure terrain exists for collision
		centerX, centerY := g.player.GetCenter()
		g.world.GetChunksInRange(centerX, centerY)

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
		centerX, centerY = g.player.GetCenter()
		g.cameraX = centerX - ScreenWidth/2
		g.cameraY = centerY - ScreenHeight/2

		// Unload distant chunks to manage memory
		g.world.UnloadDistantChunks(centerX, centerY)
	}

	return nil
}

func (g *Game) handleMenuAction(action menu.MenuAction) {
	switch action {
	case menu.ActionStartGame:
		g.inMenu = false
		g.inGame = true
		
		// Find a suitable spawn position and set player position
		// Spawn in area where terrain is actually generated (negative coordinates)
		spawnX, spawnY := g.world.FindSpawnPosition(-2000, -2000)
		g.player.SetPosition(spawnX, spawnY)
		g.player.SetVelocity(0, 0)

	case menu.ActionBack:
		if g.menu.CurrentMenu == menu.MenuBlockLibrary {
			g.menu.SetMainMenu()
		} else {
			g.menu.SetMainMenu()
		}

	case menu.ActionOpenBlockLibrary:
		g.menu.SetBlockLibraryMenu()

	case menu.ActionSelectBlock:
		g.selectedBlock = g.menu.SelectedBlock
		g.menu.SetMainMenu()

	case menu.ActionExit:
		os.Exit(0)
		return
	}
}

func (g *Game) handleGameInput() {
	// Handle keyboard input using input manager
	if g.inputManager.IsActionPressed("move_left") {
		g.player.MovingLeft = true
		g.player.MovingRight = false
	} else if g.inputManager.IsActionPressed("move_right") {
		g.player.MovingRight = true
		g.player.MovingLeft = false
	} else {
		g.player.MovingLeft = false
		g.player.MovingRight = false
	}

	if g.inputManager.IsActionPressed("jump") {
		g.player.Jump()
	}

	// Toggle flying with F key
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.player.SetFlying(!g.player.GetIsFlying())
	}

	// Handle vertical movement when flying
	if g.player.GetIsFlying() {
		if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
			g.player.SetMovingUp(true)
			g.player.SetMovingDown(false)
		} else if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
			g.player.SetMovingDown(true)
			g.player.SetMovingUp(false)
		} else {
			g.player.SetMovingUp(false)
			g.player.SetMovingDown(false)
		}
	}

	// Open crafting menu
	if g.inputManager.IsActionJustPressed("crafting") || g.inputManager.IsActionJustPressed("inventory") {
		g.inCrafting = true
		g.craftingUI.Toggle()
	}

	// Open block library (B key) - only in creative mode
	if inpututil.IsKeyJustPressed(ebiten.KeyB) && g.CreativeMode {
		g.inMenu = true
		g.inGame = false
		g.menu.SetBlockLibraryMenu()
	}

	// Interact with crafting stations
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.interactWithStation()
	}

	// Drop item
	if g.inputManager.IsActionJustPressed("drop") {
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
	if g.inputManager.IsActionJustPressed("menu") {
		g.inGame = false
		g.inMenu = true
		g.menu.SetMainMenu()
	}

	// Handle mouse input
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	// Detect hovered block for tooltip
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY
	hoveredHex := g.world.GetHexagonAt(mouseWorldX, mouseWorldY)
	if hoveredHex != nil && hoveredHex.BlockType != blocks.AIR {
		g.hoveredBlockName = getBlockKeyFromType(hoveredHex.BlockType)
	} else {
		g.hoveredBlockName = ""
	}

	// Mouse input for mining and placement
	staticLeftPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	staticRightPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	
	// Track previous state to detect "just pressed"
	if !g.leftMouseWasPressed && staticLeftPressed {
		// Left mouse just pressed - start mining
		g.startMining()
	}
	if g.leftMouseWasPressed && !staticLeftPressed {
		// Left mouse just released - stop mining
		g.player.StopMining()
	}
	if !g.rightMouseWasPressed && staticRightPressed {
		// Right mouse just pressed - place block
		g.handleBlockPlacement()
	}
	
	// Update state
	g.leftMouseWasPressed = staticLeftPressed
	g.rightMouseWasPressed = staticRightPressed

	// Hotbar selection using input manager
	if g.inputManager.IsActionJustPressed("hotbar_1") {
		g.inventory.SelectSlot(0)
	} else if g.inputManager.IsActionJustPressed("hotbar_2") {
		g.inventory.SelectSlot(1)
	} else if g.inputManager.IsActionJustPressed("hotbar_3") {
		g.inventory.SelectSlot(2)
	} else if g.inputManager.IsActionJustPressed("hotbar_4") {
		g.inventory.SelectSlot(3)
	} else if g.inputManager.IsActionJustPressed("hotbar_5") {
		g.inventory.SelectSlot(4)
	} else if g.inputManager.IsActionJustPressed("hotbar_6") {
		g.inventory.SelectSlot(5)
	} else if g.inputManager.IsActionJustPressed("hotbar_7") {
		g.inventory.SelectSlot(6)
	} else if g.inputManager.IsActionJustPressed("hotbar_8") {
		g.inventory.SelectSlot(7)
	} else if g.inputManager.IsActionJustPressed("hotbar_9") {
		g.inventory.SelectSlot(8)
	}

	// Scroll wheel for hotbar
	_, scrollY := ebiten.Wheel()
	if scrollY > 0 {
		g.inventory.PrevSlot()
	} else if scrollY < 0 {
		g.inventory.NextSlot()
	}

	// Command system
	if g.inputManager.IsActionJustPressed("command") {
		g.commandMode = true
		g.commandString = ""
	}

	if g.commandMode {
		// Handle typing letters
		for key := ebiten.KeyA; key <= ebiten.KeyZ; key++ {
			if inpututil.IsKeyJustPressed(key) {
				g.commandString += string(rune('a' + int(key-ebiten.KeyA)))
			}
		}
		// Space
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.commandString += " "
		}
		// Backspace
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			if len(g.commandString) > 0 {
				g.commandString = g.commandString[:len(g.commandString)-1]
			}
		}
		// Enter to execute
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.executeCommand(g.commandString)
			g.commandMode = false
		}
		// Escape to cancel
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.commandMode = false
		}
	}
}

// executeCommand parses and executes the given command string
func (g *Game) executeCommand(command string) {
	// Trim leading slash if present
	command = strings.TrimPrefix(command, "/")

	// Split command into parts
	parts := strings.Fields(command)
	if len(parts) == 0 {
		log.Printf("Empty command")
		return
	}

	cmd := strings.ToLower(parts[0])
	args := parts[1:]

	switch cmd {
	case "help":
		log.Printf("Available commands: help, give, creative, survival, tp, plugin list, plugin load, plugin unload, plugin reload")
	case "give":
		if len(args) < 2 {
			log.Printf("Usage: /give <item_type> <quantity>")
			return
		}
		itemTypeStr := args[0]
		quantityStr := args[1]

		// Parse quantity
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil || quantity <= 0 {
			log.Printf("Invalid quantity: %s", quantityStr)
			return
		}

		// Find item type by name
		var itemType items.ItemType
		found := false
		for it, props := range items.ItemDefinitions {
			if strings.EqualFold(props.Name, itemTypeStr) {
				itemType = it
				found = true
				break
			}
		}
		if !found {
			log.Printf("Unknown item: %s", itemTypeStr)
			return
		}

		// Add to inventory
		if g.inventory.AddItem(itemType, quantity) {
			log.Printf("Gave %d %s", quantity, itemTypeStr)
		} else {
			log.Printf("Inventory full, could not give items")
		}
	case "creative":
		g.CreativeMode = true
		log.Printf("Switched to creative mode")
	case "survival":
		g.CreativeMode = false
		log.Printf("Switched to survival mode")
	case "tp":
		if len(args) < 2 {
			log.Printf("Usage: /tp <x> <y>")
			return
		}
		xStr, yStr := args[0], args[1]
		x, errX := strconv.ParseFloat(xStr, 64)
		y, errY := strconv.ParseFloat(yStr, 64)
		if errX != nil || errY != nil {
			log.Printf("Invalid coordinates: %s %s", xStr, yStr)
			return
		}
		g.player.SetPosition(x, y)
		g.player.SetVelocity(0, 0)
		log.Printf("Teleported to (%.1f, %.1f)", x, y)
	case "plugin":
		if len(args) < 1 {
			log.Printf("Usage: /plugin <list|load|unload|reload> [plugin_name]")
			return
		}
		g.handlePluginCommand(args[0], args[1:])
	default:
		log.Printf("Unknown command: %s", cmd)
	}
}

// handlePluginCommand handles plugin-related commands
func (g *Game) handlePluginCommand(action string, args []string) {
	if g.pluginManager == nil {
		log.Printf("Plugin manager not initialized")
		return
	}

	switch action {
	case "list":
		log.Printf("Plugin management not yet implemented in main game")
		log.Printf("Plugins would be listed here")
	case "load":
		if len(args) < 1 {
			log.Printf("Usage: /plugin load <plugin_name>")
			return
		}
		log.Printf("Plugin loading not yet implemented in main game")
		log.Printf("Would load plugin: %s", args[0])
	case "unload":
		if len(args) < 1 {
			log.Printf("Usage: /plugin unload <plugin_name>")
			return
		}
		log.Printf("Plugin unloading not yet implemented in main game")
		log.Printf("Would unload plugin: %s", args[0])
	case "reload":
		if len(args) < 1 {
			log.Printf("Usage: /plugin reload <plugin_name>")
			return
		}
		log.Printf("Plugin reloading not yet implemented in main game")
		log.Printf("Would reload plugin: %s", args[0])
	default:
		log.Printf("Unknown plugin action: %s", action)
		log.Printf("Available actions: list, load, unload, reload")
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
		// Get player position to drop item in front of player
		playerX, playerY := g.player.GetCenter()
		
		// Drop item slightly in front of player based on facing direction
		dropX := playerX + 30.0 // Drop in front of player
		dropY := playerY - 10.0 // Slightly above player center
		
		// Add some random velocity for natural dropping motion
		vx := float64(rand.Intn(100)-50) / 100.0 * 2.0 // Random horizontal velocity
		vy := -2.0 // Upward velocity for arc motion
		
		// Create dropped item entity
		droppedItem := &DroppedItem{
			Type:     selectedItem.Type,
			Quantity: 1,
			X:        dropX,
			Y:        dropY,
			VX:       vx,
			VY:       vy,
			Lifetime: time.Now().Add(5 * time.Minute), // Items disappear after 5 minutes
		}
		
		// Add to dropped items list
		g.droppedItems = append(g.droppedItems, droppedItem)
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
	if targetHex == nil || targetHex.BlockType == blocks.AIR {
		return
	}

	// Check if block is unbreakable
	blockKey := getBlockKeyFromType(targetHex.BlockType)
	blockDef := blocks.BlockDefinitions[blockKey]
	if blockDef != nil && blockDef.Hardness <= 0 {
		return // Cannot mine unbreakable blocks
	}

	// Check if player can reach the block
	if !g.player.CanReach(targetHex.X, targetHex.Y) {
		return
	}

	// In creative mode, destroy blocks instantly
	if g.CreativeMode {
		g.completeMining(targetHex)
		return
	}

	// Start mining the block (survival mode)
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

	// Check if block is unbreakable
	blockKey := getBlockKeyFromType(targetHex.BlockType)
	blockDef := blocks.BlockDefinitions[blockKey]
	if blockDef != nil && blockDef.Hardness <= 0 {
		g.player.StopMining()
		return // Cannot mine unbreakable blocks
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

	// Use the exact hexagon coordinates for removal
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

	// Trigger gravity fall
	q, r := hexagon.PixelToHex(x, y, world.HexSize)
	hex := hexagon.HexRound(q, r)
	g.fallBlocks(hex.Q, hex.R)
}

// fallBlocks makes blocks above the destroyed position fall if they have gravity
func (g *Game) fallBlocks(q, r int) {
	currentR := r + 1
	for {
		// Get pixel position for this hex
		h := hexagon.Hexagon{Q: q, R: currentR}
		x, y := hexagon.HexToPixel(h, world.HexSize)

		// Check if there's a block here
		block := g.world.GetHexagonAt(x, y)
		if block == nil || block.BlockType == blocks.AIR {
			break
		}

		// Check if it has gravity
		props := blocks.BlockDefinitions[getBlockKeyFromType(block.BlockType)]
		if props == nil || !props.Gravity {
			break
		}

		// Move it down
		newH := hexagon.Hexagon{Q: q, R: currentR - 1}
		newX, newY := hexagon.HexToPixel(newH, world.HexSize)

		g.world.RemoveHexagonAt(x, y)
		g.world.AddHexagonAt(newX, newY, block.BlockType)

		currentR++
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

		// Trigger gravity fall
		q, r := hexagon.PixelToHex(x, y, world.HexSize)
		hex := hexagon.HexRound(q, r)
		g.fallBlocks(hex.Q, hex.R)
	}
}

// handleBlockPlacement handles block placement
func (g *Game) handleBlockPlacement() {
	var blockTypeToPlace string

	if g.CreativeMode && g.selectedBlock != "" {
		// Use selected block from library in creative mode
		blockTypeToPlace = g.selectedBlock
	} else {
		// Normal inventory-based placement
		selectedItem := g.inventory.GetSelectedItem()
		if selectedItem == nil || selectedItem.Type == items.NONE {
			return
		}

		// Get item properties
		props := items.GetItemProperties(selectedItem.Type)
		if props == nil || !props.IsPlaceable {
			return
		}

		blockTypeToPlace = props.BlockType
	}

	// Convert mouse position to world coordinates
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY

	// Find which hexagon the mouse is over using the same system as world generation
	// Convert world coordinates to local chunk coordinates
	chunkX, chunkY := g.world.GetChunkCoords(mouseWorldX, mouseWorldY)
	chunk := g.world.GetChunk(chunkX, chunkY)
	if chunk == nil {
		return
	}
	
	// Get chunk world position
	worldX, worldY := chunk.GetWorldPosition()
	
	// Calculate local row/col using the same formula as world generation
	localRow := int((mouseWorldY - worldY) / world.HexVSpacing)
	var localCol int
	if localRow%2 == 0 {
		localCol = int((mouseWorldX - worldX - world.HexWidth/2) / world.HexWidth)
	} else {
		localCol = int((mouseWorldX - worldX - world.HexWidth) / world.HexWidth)
	}
	
	// Convert back to world coordinates using the same formula as world generation
	var placeX, placeY float64
	if localRow%2 == 0 {
		placeX = worldX + float64(localCol)*world.HexWidth + world.HexWidth/2
	} else {
		placeX = worldX + float64(localCol)*world.HexWidth + world.HexWidth
	}
	placeY = worldY + float64(localRow)*world.HexVSpacing + world.HexSize

	// Placement validation: check if position is valid
	if !g.canPlaceBlockAt(placeX, placeY) {
		return // Cannot place block here
	}

	// Check if player can reach the block placement position
	if !g.player.CanReach(placeX, placeY) {
		return // Too far from player
	}

	// Place block at the calculated position
	blockType := stringToBlockType(blockTypeToPlace)
	g.world.AddHexagonAt(placeX, placeY, blockType)

	// Remove item from inventory only if not in creative mode
	if !g.CreativeMode {
		g.inventory.RemoveItem(1)
	}
}

// canPlaceBlockAt checks if a block can be placed at the given position
func (g *Game) canPlaceBlockAt(x, y float64) bool {
	// Use the coordinates directly since we're now using snapped hexagon centers
	centerX, centerY := x, y
	
	// Get the chunk directly and check the exact hexagon position
	chunkX, chunkY := g.world.GetChunkCoords(centerX, centerY)
	chunk := g.world.GetChunk(chunkX, chunkY)
	if chunk == nil {
		return false
	}
	
	// Use the same calculation as chunk.GetHexagon for exact lookup
	worldX, worldY := chunk.GetWorldPosition()
	
	localRow := int((centerY - worldY) / world.HexVSpacing)
	localCol := int((centerX - worldX - world.HexWidth/2) / world.HexWidth)
	if localRow%2 == 0 {
		localCol = int((centerX - worldX - world.HexWidth/2) / world.HexWidth)
	} else {
		localCol = int((centerX - worldX) / world.HexWidth)
	}
	
	key := [2]int{localCol, localRow}
	existingHex := chunk.Hexagons[key]
	
	if existingHex != nil {
		return false // Cannot place on existing block
	}

	// Check if player is too close (prevent placing blocks inside player)
	playerCenterX, playerCenterY := g.player.GetCenter()
	distance := math.Sqrt((x-playerCenterX)*(x-playerCenterX) + (y-playerCenterY)*(y-playerCenterY))
	minDistance := g.player.Width/2 + 2 // Very small buffer - just prevent direct overlap
	if distance < minDistance {
		return false // Too close to player
	}

	return true
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

// updateDroppedItems updates physics for all dropped items
func (g *Game) updateDroppedItems(deltaTime float64) {
	gravity := 9.8 * 50.0 // Scale gravity for pixel space
	friction := 0.98       // Air friction
	
	// Update items in reverse order to safely remove expired items
	for i := len(g.droppedItems) - 1; i >= 0; i-- {
		item := g.droppedItems[i]
		
		// Check if item has expired
		if time.Now().After(item.Lifetime) {
			g.droppedItems = append(g.droppedItems[:i], g.droppedItems[i+1:]...)
			continue
		}
		
		// Apply physics
		item.VY += gravity * deltaTime // Apply gravity
		item.VX *= friction           // Apply friction
		item.VY *= friction
		
		// Update position
		item.X += item.VX * deltaTime
		item.Y += item.VY * deltaTime
		
		// Simple ground collision - check if item hit ground
		groundY := float64(1000) // Simple ground level for now
		if item.Y > groundY {
			item.Y = groundY
			item.VY = 0
			item.VX *= 0.8 // Ground friction
		}
		
		// Check for player pickup (proximity check)
		playerX, playerY := g.player.GetCenter()
		distance := math.Sqrt((item.X-playerX)*(item.X-playerX) + (item.Y-playerY)*(item.Y-playerY))
		if distance < 30.0 { // Pickup range
			// Try to add to inventory
			if g.inventory.AddItem(item.Type, item.Quantity) {
				// Remove picked up item
				g.droppedItems = append(g.droppedItems[:i], g.droppedItems[i+1:]...)
			}
		}
	}
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

	// Draw dropped items
	g.drawDroppedItems(screen)

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

		// Check if block is on screen before grouping
		screenX := block.X - g.cameraX
		screenY := block.Y - g.cameraY
		if screenX < -100 || screenX > ScreenWidth+100 ||
			screenY < -100 || screenY > ScreenHeight+100 {
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

		// Prepare vertices for all blocks in this color group using object pools
		totalVertices := len(groupBlocks) * 6 // 6 vertices per hexagon
		
		// Get pooled slices and reset them
		poolIdx := g.poolIndex % len(g.vertexPool)
		vertices := g.vertexPool[poolIdx][:0] // Reset length but keep capacity
		indices := g.indicesPool[poolIdx][:0] // Reset length but keep capacity
		
		// Ensure capacity is sufficient
		if cap(vertices) < totalVertices {
			vertices = make([]ebiten.Vertex, 0, totalVertices)
		}
		if cap(indices) < totalVertices {
			indices = make([]uint16, 0, totalVertices)
		}
		
		g.poolIndex++ // Rotate through pools

		baseIndex := uint16(0)

		for _, block := range groupBlocks {
			// Calculate screen position
			screenX := block.X - g.cameraX
			screenY := block.Y - g.cameraY

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

		// Draw the batch
		if props.Texture != nil {
			// Draw textured hexagons using triangles with texture mapping
			for _, block := range groupBlocks {
				screenX := block.X - g.cameraX
				screenY := block.Y - g.cameraY

				// Get hexagon corners
				corners := hexagon.GetHexCorners(screenX, screenY, world.HexSize)

				// Prepare vertices with texture coordinates
				vertices := make([]ebiten.Vertex, len(corners))
				for i, corner := range corners {
					relativeX := corner[0] - screenX
					relativeY := corner[1] - screenY
					srcX := 32 + (relativeX/30)*32
					srcY := 32 + (relativeY/30)*32

					r := float32(1.0)
					gc := float32(1.0)
					b := float32(1.0)
					a := float32(1.0)

					// Apply damage darkening
					if block.Health < block.MaxHealth {
						damageRatio := block.Health / block.MaxHealth
						r *= float32(damageRatio)
						gc *= float32(damageRatio)
						b *= float32(damageRatio)
					}

					vertices[i] = ebiten.Vertex{
						DstX:   float32(corner[0]),
						DstY:   float32(corner[1]),
						SrcX:   float32(srcX),
						SrcY:   float32(srcY),
						ColorR: r,
						ColorG: gc,
						ColorB: b,
						ColorA: a,
					}
				}

				// Indices for hexagon triangles
				indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5}

				screen.DrawTriangles(vertices, indices, props.Texture, nil)
			}
		} else {
			// Draw solid colors using triangles
			// Prepare vertices for all blocks in this color group
			totalVertices := len(groupBlocks) * 6 // 6 vertices per hexagon
			vertices := make([]ebiten.Vertex, 0, totalVertices)
			indices := make([]uint16, 0, len(groupBlocks)*6) // 6 indices per hexagon

			baseIndex := uint16(0)

			for _, block := range groupBlocks {
				// Calculate screen position
				screenX := block.X - g.cameraX
				screenY := block.Y - g.cameraY

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

			// Draw the triangle batch
			if len(vertices) > 0 {
				screen.DrawTriangles(vertices, indices, g.whiteImage, nil)
			}
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

				// Draw item name
				itemName := items.ItemNameByID(item.Type)
				// Truncate if too long
				if len(itemName) > 8 {
					itemName = itemName[:8]
				}
				ebitenutil.DebugPrintAt(screen, itemName, slotX+5, slotY+hotbarHeight-5)

				// Draw quantity indicator
				if item.Quantity > 1 {
					quantityStr := fmt.Sprintf("%d", item.Quantity)
					ebitenutil.DebugPrintAt(screen, quantityStr, slotX+slotWidth-15, slotY+hotbarHeight-5)
				}
			}
		}
	}

	// Draw hovered block tooltip
	if g.hoveredBlockName != "" {
		ebitenutil.DebugPrintAt(screen, strings.Title(g.hoveredBlockName), g.mouseX+10, g.mouseY-20)
	}
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

// drawDroppedItems renders all dropped items in the world
func (g *Game) drawDroppedItems(screen *ebiten.Image) {
	for _, item := range g.droppedItems {
		// Calculate screen position
		screenX := item.X - g.cameraX
		screenY := item.Y - g.cameraY
		
		// Skip if off screen
		if screenX < -50 || screenX > ScreenWidth+50 || screenY < -50 || screenY > ScreenHeight+50 {
			continue
		}
		
		// Get item properties for color
		itemProps := items.GetItemProperties(item.Type)
		if itemProps == nil {
			continue
		}
		
		// Draw item as a small square with item color
		itemSize := 16.0
		halfSize := itemSize / 2.0
		
		// Create simple rectangle representation
		color := g.colorToRGB(int(itemProps.IconColor.R), int(itemProps.IconColor.G), int(itemProps.IconColor.B))
		ebitenutil.DrawRect(screen, screenX-halfSize, screenY-halfSize, itemSize, itemSize, color)
		
		// Draw border
		borderColor := g.colorToRGB(0, 0, 0)
		ebitenutil.DrawRect(screen, screenX-halfSize, screenY-halfSize, itemSize, 1, borderColor)
		ebitenutil.DrawRect(screen, screenX-halfSize, screenY-halfSize+itemSize-1, itemSize, 1, borderColor)
		ebitenutil.DrawRect(screen, screenX-halfSize, screenY-halfSize, 1, itemSize, borderColor)
		ebitenutil.DrawRect(screen, screenX-halfSize+itemSize-1, screenY-halfSize, 1, itemSize, borderColor)
		
		// Draw quantity if > 1
		if item.Quantity > 1 {
			quantityStr := fmt.Sprintf("%d", item.Quantity)
			ebitenutil.DebugPrintAt(screen, quantityStr, int(screenX+halfSize-5), int(screenY-halfSize-5))
		}
	}
}

// drawBlock draws a hexagonal block
func (g *Game) drawBlock(screen *ebiten.Image, block *world.Hexagon) {
	if block.BlockType == blocks.AIR {
		return
	}

	blockKey := getBlockKeyFromType(block.BlockType)
	props := blocks.BlockDefinitions[blockKey]
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

	// Center of hexagon
	centerX := screenX
	centerY := screenY

	// Check if mouse is hovering over this block
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY

	q, r := hexagon.PixelToHex(mouseWorldX, mouseWorldY, world.HexSize)
	hoverHex := hexagon.HexRound(q, r)

	// Get block's hex coordinates
	blockQ, blockR := hexagon.PixelToHex(block.X, block.Y, world.HexSize)
	blockHex := hexagon.HexRound(blockQ, blockR)

	isHovered := hoverHex.Q == blockHex.Q && hoverHex.R == blockHex.R

	// Damage ratio
	damageRatio := 1.0
	if block.Health < block.MaxHealth {
		damageRatio = block.Health / block.MaxHealth
	}

	if props.Pattern == "striped" || props.Pattern == "" { // Default to solid if empty
		// Draw as solid or striped
		color1 := props.Color
		color2 := props.Color
		if props.TopColor.A > 0 {
			color2 = props.TopColor
		} else {
			// Darker shade
			color2 = color.RGBA{
				R: uint8(float64(props.Color.R) * 0.7),
				G: uint8(float64(props.Color.G) * 0.7),
				B: uint8(float64(props.Color.B) * 0.7),
				A: props.Color.A,
			}
		}

		// Apply damage darkening
		color1 = color.RGBA{
			R: uint8(float64(color1.R) * damageRatio),
			G: uint8(float64(color1.G) * damageRatio),
			B: uint8(float64(color1.B) * damageRatio),
			A: color1.A,
		}
		color2 = color.RGBA{
			R: uint8(float64(color2.R) * damageRatio),
			G: uint8(float64(color2.G) * damageRatio),
			B: uint8(float64(color2.B) * damageRatio),
			A: color2.A,
		}

		if isHovered {
			color1.R = uint8(minFloat32(255, float32(color1.R)+50))
			color1.G = uint8(minFloat32(255, float32(color1.G)+50))
			color1.B = uint8(minFloat32(255, float32(color1.B)+50))
			color2.R = uint8(minFloat32(255, float32(color2.R)+50))
			color2.G = uint8(minFloat32(255, float32(color2.G)+50))
			color2.B = uint8(minFloat32(255, float32(color2.B)+50))
		}

		switch props.Pattern {
		case "solid", "":
			// Single color
			vertices := make([]ebiten.Vertex, len(corners))
			for i, corner := range corners {
				r := float32(color1.R) / 255.0
				gc := float32(color1.G) / 255.0
				b := float32(color1.B) / 255.0
				a := float32(color1.A) / 255.0
				vertices[i] = ebiten.Vertex{
					DstX:   float32(corner[0]),
					DstY:   float32(corner[1]),
					ColorR: r,
					ColorG: gc,
					ColorB: b,
					ColorA: a,
				}
			}

			indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5}
			screen.DrawTriangles(vertices, indices, g.whiteImage, nil)
		case "striped":
			// Draw each triangle with alternating colors
			for i := 0; i < 6; i++ {
				triangleIndices := []uint16{0, uint16(i + 1), uint16((i+1)%6 + 1)}
				vertices := []ebiten.Vertex{
					{
						DstX: float32(centerX),
						DstY: float32(centerY),
					},
					{
						DstX: float32(corners[i][0]),
						DstY: float32(corners[i][1]),
					},
					{
						DstX: float32(corners[(i+1)%6][0]),
						DstY: float32(corners[(i+1)%6][1]),
					},
				}

				var triColor color.RGBA
				if i%2 == 0 {
					triColor = color1
				} else {
					triColor = color2
				}

				for j := range vertices {
					vertices[j].ColorR = float32(triColor.R) / 255.0
					vertices[j].ColorG = float32(triColor.G) / 255.0
					vertices[j].ColorB = float32(triColor.B) / 255.0
					vertices[j].ColorA = float32(triColor.A) / 255.0
				}

				screen.DrawTriangles(vertices, triangleIndices, g.whiteImage, nil)
			}
		}
	} else {
		// Fallback to solid
		vertices := make([]ebiten.Vertex, len(corners))
		for i, corner := range corners {
			r := float32(props.Color.R) / 255.0
			gc := float32(props.Color.G) / 255.0
			b := float32(props.Color.B) / 255.0
			a := float32(props.Color.A) / 255.0

			if isHovered {
				r = minFloat32(1.0, r+0.2)
				gc = minFloat32(1.0, gc+0.2)
				b = minFloat32(1.0, b+0.2)
			}

			if block.Health < block.MaxHealth {
				damageRatio := block.Health / block.MaxHealth
				r *= float32(math.Min(1.0, float64(damageRatio)))
				gc *= float32(math.Min(1.0, float64(damageRatio)))
				b *= float32(math.Min(1.0, float64(damageRatio)))
			}

			vertices[i] = ebiten.Vertex{
				DstX:   float32(corner[0]),
				DstY:   float32(corner[1]),
				ColorR: r,
				ColorG: gc,
				ColorB: b,
				ColorA: a,
			}
		}

		indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5}
		screen.DrawTriangles(vertices, indices, g.whiteImage, nil)
	}
}

// drawPlayer draws the player
func (g *Game) drawPlayer(screen *ebiten.Image) {
	screenX := g.player.X - g.cameraX
	screenY := g.player.Y - g.cameraY

	// Draw player body
	bodyColor := g.colorToRGB(255, 100, 100)
	ebitenutil.DrawRect(screen, screenX, screenY, float64(g.player.Width), float64(g.player.Height), bodyColor)

	// Draw player head (simple representation)
	headColor := g.colorToRGB(255, 200, 150)
	ebitenutil.DrawRect(screen, screenX+10, screenY+5, 20, 20, headColor)
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
	return g.calculateMiningDamage(blockType) * 60
}

// calculateMiningDamage calculates the damage dealt to a block per mining tick
func (g *Game) calculateMiningDamage(blockType blocks.BlockType) float64 {
	// Get block hardness
	blockKey := getBlockKeyFromType(blockType)
	blockDef := blocks.BlockDefinitions[blockKey]
	if blockDef == nil {
		return 1.0 // Default damage
	}

	// Unbreakable blocks (hardness <= 0)
	if blockDef.Hardness <= 0 {
		return 0
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
	// Load block definitions before any world generation
	blocks.LoadBlocks()

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Tesselbox v2.0 - Hexagon Sandbox")
	ebiten.SetTPS(FPS)
	
	// Enable input
	ebiten.SetCursorMode(ebiten.CursorModeVisible)

	game := NewGame()

	// Start auto-saver
	game.StartAutoSave()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
