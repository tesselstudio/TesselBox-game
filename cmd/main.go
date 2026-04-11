package main

import (
	"fmt"
	"hash/fnv"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"tesselbox/pkg/audio"
	"tesselbox/pkg/biomes"
	"tesselbox/pkg/blocks"
	"tesselbox/pkg/chest"
	"tesselbox/pkg/combat"
	"tesselbox/pkg/crafting"
	"tesselbox/pkg/enemies"
	"tesselbox/pkg/equipment"
	"tesselbox/pkg/gametime"
	"tesselbox/pkg/health"
	"tesselbox/pkg/hexagon"
	"tesselbox/pkg/input"
	"tesselbox/pkg/items"
	"tesselbox/pkg/player"
	"tesselbox/pkg/plugins"
	"tesselbox/pkg/save"
	"tesselbox/pkg/skin"
	"tesselbox/pkg/survival"
	"tesselbox/pkg/ui"
	"tesselbox/pkg/weather"
	"tesselbox/pkg/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// getTesselboxDir returns the user's home directory .tesselbox folder
func getTesselboxDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home dir can't be found
		return ".tesselbox"
	}
	return filepath.Join(home, ".tesselbox")
}

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
		"workbench":   blocks.WORKBENCH,
		"furnace":     blocks.FURNACE,
		"anvil":       blocks.ANVIL,
		// New blocks
		"gravel":            blocks.GRAVEL,
		"sandstone":         blocks.SANDSTONE,
		"obsidian":          blocks.OBSIDIAN,
		"ice":               blocks.ICE,
		"snow":              blocks.SNOW,
		"torch":             blocks.TORCH,
		"crafting_table":    blocks.CRAFTING_TABLE,
		"chest":             blocks.CHEST,
		"ladder":            blocks.LADDER,
		"fence":             blocks.FENCE,
		"gate":              blocks.GATE,
		"door":              blocks.DOOR,
		"window":            blocks.WINDOW,
		"flower":            blocks.FLOWER,
		"tall_grass":        blocks.TALL_GRASS,
		"mushroom_red":      blocks.MUSHROOM_RED,
		"mushroom_brown":    blocks.MUSHROOM_BROWN,
		"wool":              blocks.WOOL,
		"bookshelf":         blocks.BOOKSHELF,
		"jukebox":           blocks.JUKEBOX,
		"note_block":        blocks.NOTE_BLOCK,
		"pumpkin":           blocks.PUMPKIN,
		"melon":             blocks.MELON,
		"hay_bale":          blocks.HAY_BALE,
		"cobblestone":       blocks.COBBLESTONE,
		"mossy_cobblestone": blocks.MOSSY_COBBLESTONE,
		"stone_bricks":      blocks.STONE_BRICKS,
		"chiseled_stone":    blocks.CHISELED_STONE,
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
	craftingSystem  *crafting.CraftingSystem
	craftingUI      *crafting.CraftingUI
	pluginManager   *plugins.PluginManager
	pluginUI        *plugins.PluginUI
	pluginInstaller *plugins.PluginInstaller
	skinEditor      *skin.SkinEditor
	inputManager    *input.InputManager

	// Save system
	saveManager *save.SaveManager
	autoSaver   *save.AutoSaver

	// Day/night cycle
	dayNightCycle *gametime.DayNightCycle

	// Weather system
	weatherSystem *weather.WeatherSystem

	// Audio system
	audioManager *audio.AudioManager
	soundLibrary *audio.SoundLibrary

	// Footstep audio tracking
	lastFootstepTime time.Time
	footstepCooldown time.Duration

	// Dropped items
	droppedItems []*DroppedItem

	// Object pools for rendering optimization
	vertexPool  [][]ebiten.Vertex
	indicesPool [][]uint16
	poolIndex   int

	// Camera
	cameraX, cameraY float64

	// Mouse
	mouseX, mouseY       int
	hoveredBlockName     string
	leftMouseWasPressed  bool
	rightMouseWasPressed bool

	// Game state flags
	inMenu       bool
	inGame       bool
	inCrafting   bool
	inPluginUI   bool
	inSkinEditor bool
	CreativeMode bool

	// Command system
	commandMode   bool
	commandString string

	// Timing
	lastTime     time.Time
	MiningDamage float64
	MiningSpeed  float64

	// Statistics tracking
	BlocksPlaced    int
	BlocksDestroyed int
	ItemsCrafted    int
	PlayStartTime   time.Time
	TotalPlayTime   time.Duration

	// Crafting station tracking
	CurrentCraftingStation string

	// Unlocked recipes tracking
	UnlockedRecipes map[string]bool

	// For solid color drawing
	whiteImage *ebiten.Image

	// Survival mode systems
	survivalManager *survival.SurvivalManager
	equipmentSet    *equipment.EquipmentSet
	healthSystem    *health.LocationalHealthSystem
	backpackUI      *ui.BackpackUI
	hud             *ui.HUD

	// Chest system
	chestManager *chest.ChestManager
	chestUI      *ui.ChestUI

	// Combat system
	weaponSystem *combat.WeaponSystem

	// Damage indicators
	damageIndicators  *ui.DamageIndicatorManager
	screenFlash       *ui.ScreenFlash
	directionalHitInd *ui.DirectionalHitManager

	// Death screen
	deathScreen *ui.DeathScreen
	isDead      bool

	// Enemy systems
	zombieSpawner *enemies.ZombieSpawner

	// Layer system (surface=0, middle=1, back=2)
	currentLayer int
	totalLayers  int
}

// NewGame creates a new game with default world
func NewGame() *Game {
	return NewGameWithWorld("default", 0)
}

// NewGameWithWorld creates a new game with a specific world name and seed
func NewGameWithWorld(worldName string, worldSeed int64) *Game {
	// Create a 1x1 white image for solid color drawing
	whiteImage := ebiten.NewImage(1, 1)
	whiteImage.Fill(color.RGBA{255, 255, 255, 255})

	g := &Game{
		world:                  world.NewWorld(worldName), // Create world with name
		player:                 player.NewPlayer(400, 300),
		inventory:              items.NewInventory(32),
		selectedBlock:          "dirt",
		CreativeMode:           true,
		cameraX:                0,
		cameraY:                0,
		lastTime:               time.Now(),
		whiteImage:             whiteImage,
		leftMouseWasPressed:    false,
		rightMouseWasPressed:   false,
		PlayStartTime:          time.Now(),
		BlocksPlaced:           0,
		BlocksDestroyed:        0,
		ItemsCrafted:           0,
		CurrentCraftingStation: "",
		UnlockedRecipes:        make(map[string]bool),
	}

	// Set the world seed if provided (non-zero)
	if worldSeed != 0 {
		g.world.SetSeed(worldSeed)
		log.Printf("World '%s' created with seed: %d", worldName, worldSeed)
	} else {
		log.Printf("World '%s' created with random seed: %d", worldName, g.world.GetSeed())
	}

	// Initialize object pools for rendering optimization
	g.vertexPool = make([][]ebiten.Vertex, 10)
	g.indicesPool = make([][]uint16, 10)
	for i := range g.vertexPool {
		g.vertexPool[i] = make([]ebiten.Vertex, 0, 1000)
		g.indicesPool[i] = make([]uint16, 0, 1000)
	}
	g.poolIndex = 0

	// Find spawn position
	spawnX, spawnY := g.world.FindSpawnPosition(0, 0)
	groundY := spawnY
	for checkY := spawnY; checkY < spawnY+1000; checkY += 30 {
		hex := g.world.GetHexagonAt(spawnX, checkY)
		if hex != nil && hex.BlockType != blocks.AIR {
			blockKey := getBlockKeyFromType(hex.BlockType)
			def := blocks.BlockDefinitions[blockKey]
			if def != nil && def.Solid {
				groundY = checkY
				break
			}
		}
	}
	if groundY > spawnY {
		spawnY = groundY - 200
	}
	g.player.SetPosition(spawnX, spawnY)
	g.world.GetChunksInRange(spawnX, spawnY)
	log.Printf("Player spawned at: (%.1f, %.1f) in world '%s'", spawnX, spawnY, worldName)

	// Initialize crafting system
	g.craftingSystem = crafting.NewCraftingSystem()
	if err := g.craftingSystem.LoadRecipesFromAssets(); err != nil {
		log.Printf("Warning: Failed to load crafting recipes: %v", err)
	}
	g.craftingSystem.OnItemCrafted = func(recipeID string) {
		g.ItemsCrafted++
		g.UnlockedRecipes[recipeID] = true
	}
	g.craftingUI = crafting.NewCraftingUI(g.craftingSystem, g.inventory)

	// Initialize input manager
	g.inputManager = input.NewInputManager()

	// Load game assets
	log.Printf("Loading game assets for world '%s'...", worldName)
	items.LoadItems()
	blocks.LoadBlocks()

	// Add initial items
	g.inventory.AddItem(items.DIRT_BLOCK, 64)
	g.inventory.AddItem(items.STONE_BLOCK, 64)
	g.inventory.AddItem(items.GRASS_BLOCK, 64)
	g.inventory.AddItem(items.WORKBENCH, 1)
	g.inventory.AddItem(items.WOODEN_PICKAXE, 1)
	g.inventory.AddItem(items.COAL, 10)

	// Initialize plugin system
	g.pluginManager = plugins.NewPluginManager()
	defaultPlugin := plugins.NewDefaultPlugin()
	g.pluginManager.RegisterPlugin(defaultPlugin)
	g.pluginManager.EnablePlugin("default")

	// Initialize skin editor
	g.skinEditor = skin.NewSkinEditor()

	// Initialize save system with world name
	g.saveManager = save.NewSaveManager(worldName, "player")

	// Initialize day/night cycle
	g.dayNightCycle = gametime.NewDayNightCycle(600.0)

	// Initialize weather system
	g.weatherSystem = weather.NewWeatherSystem()

	// Initialize audio system
	g.audioManager = audio.NewAudioManager()
	g.soundLibrary = audio.NewSoundLibrary(g.audioManager)

	// Load audio files with validation
	log.Printf("Loading audio system...")
	loader := audio.NewAudioLoader(g.audioManager)
	if err := loader.LoadAllAudio(); err != nil {
		log.Printf("Warning: Failed to load audio files: %v", err)
		// Clean up any partially loaded resources before loading placeholders
		g.audioManager.Cleanup()
		// Load placeholder sounds if real audio files are missing
		if err := loader.LoadPlaceholderSounds(); err != nil {
			log.Printf("Warning: Failed to load placeholder sounds: %v", err)
		} else {
			log.Printf("Loaded placeholder audio sounds")
		}
	} else {
		log.Printf("Audio system loaded successfully")
	}

	// Initialize sound library
	g.soundLibrary.InitializeDefaultSounds()

	// Validate audio system
	loadedSounds := g.audioManager.GetLoadedSounds()
	if len(loadedSounds) == 0 {
		log.Printf("WARNING: No audio sounds available")
	} else {
		log.Printf("Audio validation passed: %d sounds loaded", len(loadedSounds))
	}

	// Initialize footstep tracking
	g.lastFootstepTime = time.Now()
	g.footstepCooldown = 400 * time.Millisecond // Footstep every 400ms when walking

	// Initialize survival mode systems
	log.Printf("Initializing survival mode systems...")

	// Create survival manager (start in survival mode)
	gameMode := survival.ModeSurvival
	if g.CreativeMode {
		gameMode = survival.ModeCreative
	}
	g.survivalManager = survival.NewSurvivalManager(gameMode, g.player, g.inventory)

	// Create equipment set
	g.equipmentSet = equipment.NewEquipmentSet()

	// Create locational health system
	g.healthSystem = health.NewLocationalHealthSystem()

	// Create backpack UI
	g.backpackUI = ui.NewBackpackUI(ScreenWidth, ScreenHeight, g.inventory, g.equipmentSet, g.healthSystem)
	g.backpackUI.SetSelectedSlot(g.player.GetSelectedSlot())

	// Create zombie spawner
	g.zombieSpawner = enemies.NewZombieSpawner(g.dayNightCycle)

	// Set up damage callback for zombie attacks
	g.zombieSpawner.OnPlayerDamage = func(damage float64, zombieX, zombieY float64) {
		// Trigger screen flash
		if g.screenFlash != nil {
			g.screenFlash.Trigger(color.RGBA{255, 0, 0, 100}, 500*time.Millisecond)
		}

		// Determine damage tier based on damage amount
		var tier ui.DamageTier
		isCritical := false
		switch {
		case damage >= 15:
			tier = ui.TierPurple // Fatal
			isCritical = true
		case damage >= 10:
			tier = ui.TierRed // Severe
			isCritical = true
		case damage >= 5:
			tier = ui.TierYellow // Moderate
		default:
			tier = ui.TierGreen // Low
		}

		// Spawn damage indicator
		if g.damageIndicators != nil {
			g.damageIndicators.SpawnDamageIndicator(g.player.X, g.player.Y, damage, tier, isCritical)
		}

		// Trigger directional hit indicator
		if g.directionalHitInd != nil {
			// Determine direction based on zombie position relative to player
			if zombieX < g.player.X {
				g.directionalHitInd.TriggerHit(ui.DirLeft)
			} else {
				g.directionalHitInd.TriggerHit(ui.DirRight)
			}
		}
	}

	// Add some starter equipment for testing
	g.equipmentSet.EquipItem(equipment.CreateArmor("Leather Cap", equipment.SlotHelmet, equipment.MaterialLeather, equipment.ArmorLight), equipment.SlotHelmet)
	g.equipmentSet.EquipItem(equipment.CreateArmor("Leather Tunic", equipment.SlotChestplate, equipment.MaterialLeather, equipment.ArmorLight), equipment.SlotChestplate)
	g.equipmentSet.EquipItem(equipment.CreateArmor("Leather Pants", equipment.SlotLeggings, equipment.MaterialLeather, equipment.ArmorLight), equipment.SlotLeggings)
	g.equipmentSet.EquipItem(equipment.CreateArmor("Leather Boots", equipment.SlotBoots, equipment.MaterialLeather, equipment.ArmorLight), equipment.SlotBoots)

	// Create wings and equip them
	wings := equipment.CreateWings("Angel", equipment.MaterialCloth)
	g.equipmentSet.EquipItem(wings, equipment.SlotWings)

	// Create HUD
	g.hud = ui.NewHUD(ScreenWidth, ScreenHeight, g.survivalManager, g.equipmentSet, g.healthSystem, g.dayNightCycle)

	// Create chest system
	g.chestManager = chest.NewChestManager(worldName)

	// Create combat system
	g.weaponSystem = combat.NewWeaponSystem()
	// Load existing chests
	if err := g.chestManager.LoadChests(); err != nil {
		log.Printf("Failed to load chests: %v", err)
	}
	g.chestUI = ui.NewChestUI(ScreenWidth, ScreenHeight, g.chestManager, g.inventory)

	// Create damage indicators
	g.damageIndicators = ui.NewDamageIndicatorManager(ScreenWidth, ScreenHeight)
	g.screenFlash = ui.NewScreenFlash()
	g.directionalHitInd = ui.NewDirectionalHitManager()

	// Create death screen
	g.deathScreen = ui.NewDeathScreen(ScreenWidth, ScreenHeight)
	g.deathScreen.OnRespawn = func() {
		g.respawnPlayer()
	}
	g.deathScreen.OnMainMenu = func() {
		g.returnToMainMenu()
	}

	// Initialize layer system (surface=0, middle=1, back=2)
	g.currentLayer = 0
	g.totalLayers = 3
	log.Printf("Layer system initialized: surface=0, middle=1, back=2")

	log.Printf("Survival systems initialized: Equipment slots filled, wings equipped, HUD ready")

	// Start in game mode (no menu)
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

	// Update play time statistics
	g.TotalPlayTime += time.Duration(deltaTime * float64(time.Second))

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

	// Handle backpack UI
	if g.backpackUI.IsOpen() {
		if err := g.backpackUI.Update(); err != nil {
			return err
		}

		// Handle escape to close backpack
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.backpackUI.Toggle()
		}
		return nil
	}

	// Handle chest UI
	if g.chestUI != nil && g.chestUI.IsOpen() {
		if err := g.chestUI.Update(); err != nil {
			return err
		}

		// Handle escape to close chest
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.chestUI.Close()
		}
		return nil
	}

	// Handle plugin UI
	if g.inPluginUI {
		if err := g.pluginUI.Update(); err != nil {
			return err
		}

		// Handle escape to close plugin UI
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.inPluginUI = false
		}
		return nil
	}

	// Handle skin editor
	if g.inSkinEditor {
		if err := g.skinEditor.Update(); err != nil {
			log.Printf("Skin editor update error: %v", err)
		}

		// Handle escape to close skin editor
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			// Auto-save skin before exiting
			if err := g.skinEditor.SaveSkin(); err != nil {
				log.Printf("Failed to save skin on exit: %v", err)
			}
			g.inSkinEditor = false
			g.inGame = true
			g.playUISound("close")
		}
		return nil
	}

	// Handle death screen
	if g.deathScreen != nil && g.deathScreen.IsActive() {
		if err := g.deathScreen.Update(); err != nil {
			return err
		}
		// Skip other updates when dead
		return nil
	}

	// Handle game state
	if g.inGame {
		g.handleGameInput()

		// Update player with delta time (framerate-independent)
		g.player.Update(deltaTime)

		// Update mining progress
		g.updateMining(deltaTime)

		// Check if player died
		g.checkPlayerDeath()

		// Update day/night cycle
		g.dayNightCycle.Update()

		// Update survival systems
		g.survivalManager.Update(deltaTime)
		g.healthSystem.Update(deltaTime)

		// Update damage indicators
		if g.damageIndicators != nil {
			g.damageIndicators.Update(deltaTime)
		}
		if g.screenFlash != nil {
			g.screenFlash.Update()
		}
		if g.directionalHitInd != nil {
			g.directionalHitInd.Update()
		}

		// Update zombies
		ambientLight := g.dayNightCycle.AmbientLight
		g.zombieSpawner.Update(deltaTime, g.player, ambientLight)

		// Update weather system
		g.weatherSystem.Update(deltaTime, ScreenWidth, ScreenHeight)

		// Update audio system (clean up finished sounds)
		g.audioManager.Update()

		// Handle footstep sounds
		g.handleFootstepAudio()

		// Update dropped items physics
		g.updateDroppedItems(deltaTime)

		// Generate chunks around player first to ensure terrain exists for collision
		centerX, centerY := g.player.GetCenter()
		g.world.GetChunksInRange(centerX, centerY)

		// Apply collision-aware position update
		nearbyHexagons := g.world.GetNearbyHexagons(g.player.X, g.player.Y, 300)

		g.player.UpdateWithCollision(deltaTime, func(minX, minY, maxX, maxY float64) bool {
			for _, hex := range nearbyHexagons {
				if hex == nil {
					continue // Fixed: Add nil check for hexagon
				}

				blockKey := getBlockKeyFromType(hex.BlockType)
				def := blocks.BlockDefinitions[blockKey]
				if def == nil || !def.Solid {
					continue
				}

				hexMinX := hex.X - hex.Size
				hexMinY := hex.Y - hex.Size
				hexMaxX := hex.X + hex.Size
				hexMaxY := hex.Y + hex.Size

				collision := !(maxX < hexMinX || minX > hexMaxX || maxY < hexMinY || minY > hexMaxY)
				if collision {
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

	// Toggle flying with F key (requires wings in survival mode)
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		canFly := g.CreativeMode || g.equipmentSet.CanFly()
		if canFly {
			g.player.SetFlying(!g.player.GetIsFlying())
			if g.player.GetIsFlying() {
				log.Printf("Flying enabled - wings equipped")
			}
		} else if !g.CreativeMode {
			log.Printf("Cannot fly - no wings equipped")
		}
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

	// Layer switching: I = go deeper (toward back), K = go toward surface
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		if g.currentLayer < g.totalLayers-1 {
			g.currentLayer++
			layerNames := []string{"surface", "middle", "back"}
			log.Printf("Moved deeper to layer: %s (%d)", layerNames[g.currentLayer], g.currentLayer)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		if g.currentLayer > 0 {
			g.currentLayer--
			layerNames := []string{"surface", "middle", "back"}
			log.Printf("Moved toward surface to layer: %s (%d)", layerNames[g.currentLayer], g.currentLayer)
		}
	}

	// Open crafting menu
	if g.inputManager.IsActionJustPressed("crafting") {
		g.playUISound("open")
		g.inCrafting = true
		g.craftingUI.Toggle()
	}

	// Open backpack
	if g.inputManager.IsActionJustPressed("inventory") {
		g.playUISound("open")
		g.backpackUI.Toggle()
	}

	// Open plugin manager (P key)
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.inPluginUI = true
		g.inGame = false
		g.playUISound("open")
	}

	// Open skin editor (S key) - only in creative mode and not during gameplay
	if inpututil.IsKeyJustPressed(ebiten.KeyS) && g.CreativeMode && !g.inGame {
		g.inSkinEditor = true
		g.playUISound("open")
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

	// Quick save (F5)
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		if err := g.SaveGame(); err != nil {
			log.Printf("Failed to save game: %v", err)
		} else {
			log.Println("Game saved successfully")
		}
	}

	// Backup save (F6)
	if inpututil.IsKeyJustPressed(ebiten.KeyF6) {
		if g.saveManager != nil {
			if err := g.saveManager.BackupSave(); err != nil {
				log.Printf("Failed to backup save: %v", err)
			} else {
				log.Println("Save backup created successfully")
			}
		}
	}

	// Save info (F7)
	if inpututil.IsKeyJustPressed(ebiten.KeyF7) {
		if g.saveManager != nil {
			if info, err := g.saveManager.GetSaveInfo(); err != nil {
				log.Printf("Failed to get save info: %v", err)
			} else {
				log.Printf("Save Info - World: %s, Player: %s, Mode: %s, Play Time: %.1f min",
					info.WorldName, info.PlayerName, info.GameMode, info.PlayTime/60)
				log.Printf("Stats - Blocks Placed: %d, Destroyed: %d, Items Crafted: %d",
					info.BlocksPlaced, info.BlocksDestroyed, info.ItemsCrafted)
			}
		}
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
		// Left mouse just pressed - start mining and weapon attack
		g.startMining()
		g.performWeaponAttack()
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

	// Hotbar selection using input manager with bounds checking
	if g.inputManager.IsActionJustPressed("hotbar_1") {
		if g.inventory.SelectSlot(0) {
			g.playItemSound("hotbar_select")
		}
	} else if g.inputManager.IsActionJustPressed("hotbar_2") {
		if g.inventory.SelectSlot(1) {
			g.playItemSound("hotbar_select")
		}
	} else if g.inputManager.IsActionJustPressed("hotbar_3") {
		if g.inventory.SelectSlot(2) {
			g.playItemSound("hotbar_select")
		}
	} else if g.inputManager.IsActionJustPressed("hotbar_4") {
		if g.inventory.SelectSlot(3) {
			g.playItemSound("hotbar_select")
		}
	} else if g.inputManager.IsActionJustPressed("hotbar_5") {
		if g.inventory.SelectSlot(4) {
			g.playItemSound("hotbar_select")
		}
	} else if g.inputManager.IsActionJustPressed("hotbar_6") {
		if g.inventory.SelectSlot(5) {
			g.playItemSound("hotbar_select")
		}
	} else if g.inputManager.IsActionJustPressed("hotbar_7") {
		if g.inventory.SelectSlot(6) {
			g.playItemSound("hotbar_select")
		}
	} else if g.inputManager.IsActionJustPressed("hotbar_8") {
		if g.inventory.SelectSlot(7) {
			g.playItemSound("hotbar_select")
		}
	} else if g.inputManager.IsActionJustPressed("hotbar_9") {
		if g.inventory.SelectSlot(8) {
			g.playItemSound("hotbar_select")
		}
	}

	// Scroll wheel for hotbar
	_, scrollY := ebiten.Wheel()
	if scrollY > 0 {
		g.inventory.PrevSlot()
		g.playItemSound("hotbar_select")
	} else if scrollY < 0 {
		g.inventory.NextSlot()
		g.playItemSound("hotbar_select")
	}

	// Inventory management shortcuts
	if inpututil.IsKeyJustPressed(ebiten.KeyS) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		// Ctrl+S: Sort inventory
		g.inventory.SortInventory()
		g.inventory.ConsolidateItems()
		log.Printf("Inventory sorted and consolidated")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		// I: Show inventory stats
		stats := g.inventory.GetInventoryStats()
		log.Printf("Inventory: %d/%d slots used, %d items, %d types",
			stats["used_slots"], stats["total_slots"],
			stats["total_items"], stats["unique_types"])
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		// Ctrl+C: Consolidate items only
		g.inventory.ConsolidateItems()
		log.Printf("Items consolidated")
	}

	// Command system
	if g.inputManager.IsActionJustPressed("command") {
		g.commandMode = true
		g.commandString = ""
	}

	if g.commandMode {
		// Handle typing letters with input validation
		for key := ebiten.KeyA; key <= ebiten.KeyZ; key++ {
			if inpututil.IsKeyJustPressed(key) {
				// Fixed: Limit command length to prevent buffer overflow
				if len(g.commandString) < 100 {
					g.commandString += string(rune('a' + int(key-ebiten.KeyA)))
				}
			}
		}
		// Space with validation
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			if len(g.commandString) < 100 && len(g.commandString) > 0 {
				// Don't allow leading or multiple consecutive spaces
				lastChar := g.commandString[len(g.commandString)-1]
				if lastChar != ' ' {
					g.commandString += " "
				}
			}
		}
		// Numbers with validation
		for key := ebiten.Key0; key <= ebiten.Key9; key++ {
			if inpututil.IsKeyJustPressed(key) {
				if len(g.commandString) < 100 {
					g.commandString += string(rune('0' + int(key-ebiten.Key0)))
				}
			}
		}
		// Backspace
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			if len(g.commandString) > 0 {
				g.commandString = g.commandString[:len(g.commandString)-1]
			}
		}
		// Enter to execute
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if len(g.commandString) > 0 {
				g.executeCommand(g.commandString)
			}
			g.commandMode = false
			g.commandString = ""
		}
		// Escape to cancel
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.commandMode = false
			g.commandString = ""
		}
	}
}

// executeCommand parses and executes the given command string
func (g *Game) executeCommand(command string) {
	// Trim leading slash if present
	command = strings.TrimPrefix(command, "/")

	// Fixed: Sanitize input to prevent command injection
	command = strings.TrimSpace(command)

	// Limit command length
	if len(command) > 100 {
		log.Printf("Command too long")
		return
	}

	// Split command into parts
	parts := strings.Fields(command)
	if len(parts) == 0 {
		log.Printf("Empty command")
		return
	}

	// Fixed: Validate command characters
	for _, part := range parts {
		for _, r := range part {
			if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-') {
				log.Printf("Invalid characters in command")
				return
			}
		}
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
		vy := -2.0                                     // Upward velocity for arc motion

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
			g.CurrentCraftingStation = "workbench"
		case "furnace":
			station = crafting.STATION_FURNACE
			g.CurrentCraftingStation = "furnace"
		case "anvil":
			station = crafting.STATION_ANVIL
			g.CurrentCraftingStation = "anvil"
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
	g.CurrentCraftingStation = "" // Reset to no station
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

	// Clamp mining progress to prevent overflow
	if g.player.MiningProgress > 100.0 {
		g.player.MiningProgress = 100.0
	}

	// Check if mining is complete
	if g.player.MiningProgress >= 100.0 {
		// Mining complete - destroy the block
		g.completeMining(targetHex)
		g.player.StopMining()
	}
}

// completeMining handles the completion of mining (block destruction and item drop)
func (g *Game) completeMining(targetHex *world.Hexagon) {
	// Get the block type before removing
	blockType := targetHex.BlockType

	// Play mining complete sound
	g.playBlockSound("break", blockType)

	// Use the exact hexagon coordinates for removal
	x, y := targetHex.X, targetHex.Y
	g.world.RemoveHexagonAt(x, y)

	// Use item durability
	g.inventory.UseItem()

	// Drop mined item as floating item (like Minecraft) instead of adding directly to inventory
	minedItemType := g.getItemFromBlockType(blockType)
	if minedItemType != items.NONE {
		// Spawn floating item at the mined block position with slight random velocity
		vx := float64(rand.Intn(60)-30) / 10.0   // Random horizontal velocity: -3.0 to 3.0
		vy := -3.0 - float64(rand.Intn(20))/10.0 // Upward velocity with variation: -3.0 to -5.0

		droppedItem := &DroppedItem{
			Type:     minedItemType,
			Quantity: 1,
			X:        x,
			Y:        y - 10, // Slightly above the block center
			VX:       vx,
			VY:       vy,
			Lifetime: time.Now().Add(5 * time.Minute), // Items disappear after 5 minutes
		}

		g.droppedItems = append(g.droppedItems, droppedItem)
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

		// Track statistics
		g.BlocksDestroyed++

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

	// Check if clicking on an existing chest block
	if g.handleChestInteraction(mouseWorldX, mouseWorldY) {
		return // Opened chest, don't place block
	}

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

	// Track statistics
	g.BlocksPlaced++

	// Play block placement sound
	g.playBlockSound("place", blockType)

	// Remove item from inventory only if not in creative mode
	if !g.CreativeMode {
		g.inventory.RemoveItem(1)
	}
}

// handleChestInteraction checks if player clicked on a chest and opens it
// Returns true if a chest was interacted with
func (g *Game) handleChestInteraction(mouseWorldX, mouseWorldY float64) bool {
	// Check reach distance
	playerX, playerY := g.player.GetCenter()
	dist := math.Sqrt(math.Pow(mouseWorldX-playerX, 2) + math.Pow(mouseWorldY-playerY, 2))
	if dist > 150 { // Chest interaction range
		return false
	}

	// Get the chunk at mouse position
	chunkX, chunkY := g.world.GetChunkCoords(mouseWorldX, mouseWorldY)
	chunk := g.world.GetChunk(chunkX, chunkY)
	if chunk == nil {
		return false
	}

	// Check if there's a chest block at the clicked position
	worldX, worldY := chunk.GetWorldPosition()
	localRow := int((mouseWorldY - worldY) / world.HexVSpacing)
	var localCol int
	if localRow%2 == 0 {
		localCol = int((mouseWorldX - worldX - world.HexWidth/2) / world.HexWidth)
	} else {
		localCol = int((mouseWorldX - worldX - world.HexWidth) / world.HexWidth)
	}

	// Calculate block position
	var blockX, blockY float64
	if localRow%2 == 0 {
		blockX = worldX + float64(localCol)*world.HexWidth + world.HexWidth/2
	} else {
		blockX = worldX + float64(localCol)*world.HexWidth + world.HexWidth
	}
	blockY = worldY + float64(localRow)*world.HexVSpacing + world.HexSize

	// Check if there's a chest at this position
	// We need to check the hexagon at this position
	hex := chunk.GetHexagon(float64(localCol), float64(localRow))
	if hex != nil && hex.BlockType == blocks.CHEST {
		// Open the chest UI
		if g.chestUI != nil && g.chestManager != nil {
			g.chestUI.OpenChest(blockX, blockY)
			return true
		}
	}

	return false
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
	gravity := 9.8 * 100.0 // Scale gravity for pixel space
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
		item.VX *= friction            // Apply friction
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
				// Play pickup sound
				g.playItemSound("pickup")
				// Remove picked up item
				g.droppedItems = append(g.droppedItems[:i], g.droppedItems[i+1:]...)
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.inCrafting {
		// Draw game in background
		g.drawGameScene(screen)
		// Draw crafting UI overlay
		g.craftingUI.Draw(screen)
		return
	}

	if g.backpackUI.IsOpen() {
		// Draw game in background
		g.drawGameScene(screen)
		// Draw backpack UI overlay
		g.backpackUI.Draw(screen)
		return
	}

	if g.chestUI != nil && g.chestUI.IsOpen() {
		// Draw game in background
		g.drawGameScene(screen)
		// Draw chest UI overlay
		g.chestUI.Draw(screen)
		return
	}

	if g.inPluginUI {
		// Draw game in background
		g.drawGameScene(screen)
		// Draw plugin UI overlay
		g.pluginUI.Draw(screen)
		return
	}

	if g.inSkinEditor {
		// Draw skin editor (full screen)
		g.skinEditor.Draw(screen)
		return
	}

	if g.inGame {
		g.drawGameScene(screen)
		g.drawUI(screen)
		g.drawDebugInfo(screen)

		// Draw damage indicators (on top of everything)
		if g.damageIndicators != nil {
			g.damageIndicators.Draw(screen, g.cameraX, g.cameraY)
		}
		if g.screenFlash != nil {
			g.screenFlash.Draw(screen, ScreenWidth, ScreenHeight)
		}
		if g.directionalHitInd != nil {
			g.directionalHitInd.Draw(screen, ScreenWidth, ScreenHeight)
		}

		// Draw death screen on top of everything
		if g.deathScreen != nil {
			g.deathScreen.Draw(screen)
		}
	}
}

// drawGameScene draws the game world and player with layer support
func (g *Game) drawGameScene(screen *ebiten.Image) {
	// Draw background with day/night cycle sky color
	skyR, skyG, skyB := g.dayNightCycle.GetSkyColor()
	skyColor := g.colorToRGB(int(skyR*255), int(skyG*255), int(skyB*255))
	screen.Fill(skyColor)

	// Get player position
	px, py := g.player.GetCenter()

	// Draw back layers first (with blur effect based on distance)
	for layer := g.totalLayers - 1; layer > g.currentLayer; layer-- {
		blurAmount := float64(layer-g.currentLayer) * 0.3 // More blur for further layers
		g.drawLayer(screen, px, py, layer, blurAmount)
	}

	// Draw current layer (no blur)
	g.drawLayer(screen, px, py, g.currentLayer, 0)

	// Draw player (only on current layer)
	g.drawPlayer(screen)

	// Draw weapon swing effect (only on current layer)
	g.drawWeaponSwing(screen)

	// Draw zombies (only on current layer)
	g.drawZombies(screen)

	// Draw dropped items (only on current layer)
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

// drawLayer draws a specific layer of the world with optional blur effect
func (g *Game) drawLayer(screen *ebiten.Image, px, py float64, layer int, blurAmount float64) {
	// Get visible blocks for this layer
	visibleBlocks := g.world.GetVisibleBlocksForLayer(px, py, layer)

	if len(visibleBlocks) == 0 {
		return
	}

	// Apply blur effect if needed
	if blurAmount > 0 {
		// Draw blocks with reduced opacity and darkened colors for blur effect
		g.drawBlocksBatchedWithBlur(screen, visibleBlocks, blurAmount)
	} else {
		// Draw normally
		g.drawBlocksBatched(screen, visibleBlocks)
	}
}

// drawBlocksBatchedWithBlur draws blocks with a blur/darkening effect
func (g *Game) drawBlocksBatchedWithBlur(screen *ebiten.Image, blockList []*world.Hexagon, blurAmount float64) {
	// Calculate opacity and darkening based on blur amount
	opacity := uint8(255 * (1.0 - blurAmount*0.5))
	darkenFactor := 1.0 - blurAmount*0.4

	for _, block := range blockList {
		if block.BlockType == blocks.AIR {
			continue
		}

		// Calculate screen position
		screenX := block.X - g.cameraX
		screenY := block.Y - g.cameraY

		// Skip if off screen
		if screenX < -world.HexSize*2 || screenX > ScreenWidth+world.HexSize*2 ||
			screenY < -world.HexSize*2 || screenY > ScreenHeight+world.HexSize*2 {
			continue
		}

		// Get base color with darkening
		r, gr, b, a := block.ActiveColor.RGBA()
		darkenedR := uint8(float64(r/257) * darkenFactor)
		darkenedG := uint8(float64(gr/257) * darkenFactor)
		darkenedB := uint8(float64(b/257) * darkenFactor)

		// Apply opacity
		finalA := uint8(float64(a/257) * float64(opacity) / 255.0)

		// Draw hexagon using triangle fan
		if len(block.Corners) == 6 {
			centerX := screenX
			centerY := screenY

			// Create vertices for the hexagon
			vertices := make([]ebiten.Vertex, 7)

			// Center vertex
			vertices[0] = ebiten.Vertex{
				DstX:   float32(centerX),
				DstY:   float32(centerY),
				SrcX:   0,
				SrcY:   0,
				ColorR: float32(darkenedR) / 255.0,
				ColorG: float32(darkenedG) / 255.0,
				ColorB: float32(darkenedB) / 255.0,
				ColorA: float32(finalA) / 255.0,
			}

			// Corner vertices
			for i, corner := range block.Corners {
				vertices[i+1] = ebiten.Vertex{
					DstX:   float32(corner[0] - g.cameraX),
					DstY:   float32(corner[1] - g.cameraY),
					SrcX:   0,
					SrcY:   0,
					ColorR: float32(darkenedR) / 255.0,
					ColorG: float32(darkenedG) / 255.0,
					ColorB: float32(darkenedB) / 255.0,
					ColorA: float32(finalA) / 255.0,
				}
			}

			// Create indices for triangle fan
			indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5, 0, 5, 6, 0, 6, 1}

			// Draw the hexagon
			screen.DrawTriangles(vertices, indices, g.whiteImage, nil)
		}
	}
}

// drawUI draws the user interface
func (g *Game) drawUI(screen *ebiten.Image) {
	// Draw HUD (survival bars)
	if g.hud != nil {
		g.hud.Draw(screen)
	}

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

	// Draw selected block info (Creative Mode)
	if g.CreativeMode {
		selectedBlockText := fmt.Sprintf("Selected: %s", strings.Title(g.selectedBlock))
		ebitenutil.DebugPrintAt(screen, selectedBlockText, 10, ScreenHeight-100)

		// Draw instructions
		instructions := "Press B to open block library"
		ebitenutil.DebugPrintAt(screen, instructions, 10, ScreenHeight-80)

		placementInstructions := "Right Click to place, Left Click to break"
		ebitenutil.DebugPrintAt(screen, placementInstructions, 10, ScreenHeight-60)
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

	// Layer info
	layerNames := []string{"surface", "middle", "back"}
	layerInfo := fmt.Sprintf("Layer: %s (%d/%d)", layerNames[g.currentLayer], g.currentLayer, g.totalLayers-1)

	info := fmt.Sprintf("Pos: (%.1f, %.1f)\nVel: (%.1f, %.1f)\nFPS: %.1f\nOnGround: %v\nDelta: %.4f\n%s\n%s\n%s\n%s",
		px, py, vx, vy, ebiten.ActualFPS(), g.player.IsOnGround(), time.Since(g.lastTime).Seconds(),
		timeInfo, lightInfo, weatherInfo, layerInfo)

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

	// Get block color with multi-color variations
	biome := "forest"                 // This would come from world biome system
	depth := float64(blockR) / 1000.0 // Simple depth calculation
	blockColor := blocks.GlobalBlockAppearance.GetBlockColor(blockKey, int(block.X), int(block.Y), biome, depth)

	// Apply damage darkening
	blockColor = color.RGBA{
		R: uint8(float64(blockColor.R) * damageRatio),
		G: uint8(float64(blockColor.G) * damageRatio),
		B: uint8(float64(blockColor.B) * damageRatio),
		A: blockColor.A,
	}

	// Apply hover effect
	if isHovered {
		blockColor.R = uint8(minFloat32(255, float32(blockColor.R)+50))
		blockColor.G = uint8(minFloat32(255, float32(blockColor.G)+50))
		blockColor.B = uint8(minFloat32(255, float32(blockColor.B)+50))
	}

	// Use the multi-color system for rendering
	if props.Pattern == "striped" || props.Pattern == "" { // Default to solid if empty
		// Draw as solid or striped with multi-color support
		color1 := blockColor
		color2 := blockColor

		// Create darker shade for pattern
		if props.TopColor.A > 0 {
			color2 = props.TopColor
		} else {
			color2 = color.RGBA{
				R: uint8(float64(blockColor.R) * 0.7),
				G: uint8(float64(blockColor.G) * 0.7),
				B: uint8(float64(blockColor.B) * 0.7),
				A: blockColor.A,
			}
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

	// Draw floating player name above player
	playerName := "Player" // Could be customizable
	nameX := int(screenX + float64(g.player.Width)/2 - float64(len(playerName)*4))
	nameY := int(screenY - 35)

	// Draw name background for visibility
	ebitenutil.DrawRect(screen, float64(nameX-5), float64(nameY-2), float64(len(playerName)*8+10), 14, color.RGBA{0, 0, 0, 150})
	ebitenutil.DebugPrintAt(screen, playerName, nameX, nameY)

	// Draw player using custom skin renderer
	if g.skinEditor != nil && g.skinEditor.GetSkinData() != nil {
		renderer := skin.NewPlayerRenderer(g.skinEditor.GetSkinData())
		renderer.SetScale(1.0)
		renderer.Draw(screen, screenX+float64(g.player.Width)/2, screenY+float64(g.player.Height)/2)
	} else {
		// Fallback to blocky square player (Minecraft-style)
		bodyColor := g.colorToRGB(139, 69, 19)   // Brown body
		headColor := g.colorToRGB(255, 200, 150) // Skin tone head

		// Draw bigger square body (50x50)
		ebitenutil.DrawRect(screen, screenX-5, screenY+10, float64(g.player.Width)+10, float64(g.player.Height)-10, bodyColor)

		// Draw square head (25x25) centered on top
		headSize := 25.0
		headX := screenX + (float64(g.player.Width)-headSize)/2
		headY := screenY - 5
		ebitenutil.DrawRect(screen, headX, headY, headSize, headSize, headColor)

		// Draw simple blocky arms
		armWidth := 10.0
		armHeight := 30.0
		armColor := g.colorToRGB(139, 69, 19) // Same as body

		// Left arm
		ebitenutil.DrawRect(screen, screenX-armWidth-5, screenY+20, armWidth, armHeight, armColor)
		// Right arm
		ebitenutil.DrawRect(screen, screenX+float64(g.player.Width)+5, screenY+20, armWidth, armHeight, armColor)

		// Draw blocky legs
		legWidth := 12.0
		legHeight := 20.0
		legColor := g.colorToRGB(0, 0, 139) // Blue pants

		// Left leg
		ebitenutil.DrawRect(screen, screenX+8, screenY+float64(g.player.Height)-legHeight, legWidth, legHeight, legColor)
		// Right leg
		ebitenutil.DrawRect(screen, screenX+float64(g.player.Width)-legWidth-8, screenY+float64(g.player.Height)-legHeight, legWidth, legHeight, legColor)
	}
}

// drawZombies draws all active zombies
func (g *Game) drawZombies(screen *ebiten.Image) {
	for _, zombie := range g.zombieSpawner.Zombies {
		if !zombie.IsAlive {
			continue
		}

		screenX := zombie.X - g.cameraX
		screenY := zombie.Y - g.cameraY

		// Don't draw if off screen
		if screenX < -100 || screenX > ScreenWidth+100 || screenY < -100 || screenY > ScreenHeight+100 {
			continue
		}

		// Determine zombie color based on type and state
		var bodyColor color.RGBA
		if zombie.IsBurning {
			// Burning zombies are orange/red
			bodyColor = color.RGBA{255, 100, 50, 255}
		} else {
			switch zombie.Type {
			case enemies.ZombieNormal:
				bodyColor = color.RGBA{50, 100, 50, 255} // Green
			case enemies.ZombieFast:
				bodyColor = color.RGBA{100, 150, 50, 255} // Yellow-green
			case enemies.ZombieStrong:
				bodyColor = color.RGBA{80, 80, 80, 255} // Dark gray
			case enemies.ZombieTank:
				bodyColor = color.RGBA{30, 60, 30, 255} // Dark green
			default:
				bodyColor = color.RGBA{50, 100, 50, 255}
			}
		}

		// Draw zombie body
		ebitenutil.DrawRect(screen, screenX, screenY, zombie.Width, zombie.Height, bodyColor)

		// Draw zombie head (slightly smaller)
		headSize := zombie.Width * 0.7
		headX := screenX + (zombie.Width-headSize)/2
		headY := screenY - headSize/2
		ebitenutil.DrawRect(screen, headX, headY, headSize, headSize, color.RGBA{100, 150, 100, 255})

		// Draw eyes (white with red pupils when hostile)
		eyeSize := headSize / 5
		ebitenutil.DrawRect(screen, headX+headSize*0.2, headY+headSize*0.3, eyeSize, eyeSize, color.RGBA{255, 255, 255, 255})
		ebitenutil.DrawRect(screen, headX+headSize*0.6, headY+headSize*0.3, eyeSize, eyeSize, color.RGBA{255, 255, 255, 255})

		// Draw red pupils
		pupilSize := eyeSize / 2
		ebitenutil.DrawRect(screen, headX+headSize*0.2+eyeSize/4, headY+headSize*0.3+eyeSize/4, pupilSize, pupilSize, color.RGBA{255, 0, 0, 255})
		ebitenutil.DrawRect(screen, headX+headSize*0.6+eyeSize/4, headY+headSize*0.3+eyeSize/4, pupilSize, pupilSize, color.RGBA{255, 0, 0, 255})

		// Draw health bar above zombie
		healthBarWidth := zombie.Width
		healthBarHeight := 4.0
		healthPct := zombie.Health / zombie.MaxHealth
		healthWidth := healthBarWidth * healthPct

		// Background (gray)
		ebitenutil.DrawRect(screen, screenX, screenY-10, healthBarWidth, healthBarHeight, color.RGBA{50, 50, 50, 255})

		// Health (green to red based on health)
		var healthColor color.RGBA
		if healthPct > 0.5 {
			healthColor = color.RGBA{0, 255, 0, 255}
		} else if healthPct > 0.25 {
			healthColor = color.RGBA{255, 255, 0, 255}
		} else {
			healthColor = color.RGBA{255, 0, 0, 255}
		}
		ebitenutil.DrawRect(screen, screenX, screenY-10, healthWidth, healthBarHeight, healthColor)
	}
}

// drawWeaponSwing draws the weapon swing effect
func (g *Game) drawWeaponSwing(screen *ebiten.Image) {
	if g.weaponSystem == nil || !g.weaponSystem.IsSwinging() {
		return
	}

	// Get player position
	playerX, playerY := g.player.GetCenter()
	screenX := playerX - g.cameraX
	screenY := playerY - g.cameraY

	// Get swing arc
	startAngle, endAngle, radius := g.weaponSystem.GetSwingArc()
	progress := g.weaponSystem.GetSwingProgress()

	// Draw swing arc visualization
	arcColor := color.RGBA{255, 255, 255, 150}
	if progress < 0.5 {
		// Swing is active - brighter color
		arcColor = color.RGBA{255, 255, 200, 200}
	}

	// Draw a simple arc representation (triangular sweep)
	midAngle := (startAngle + endAngle) / 2
	angleRad := midAngle * math.Pi / 180

	// Draw line from player to arc edge
	endX := screenX + math.Cos(angleRad)*radius
	endY := screenY + math.Sin(angleRad)*radius

	// Draw the swing line
	ebitenutil.DrawLine(screen, screenX, screenY, endX, endY, arcColor)

	// Draw arc edges
	startRad := startAngle * math.Pi / 180
	endRad := endAngle * math.Pi / 180
	ebitenutil.DrawLine(screen, screenX, screenY, screenX+math.Cos(startRad)*radius*0.5, screenY+math.Sin(startRad)*radius*0.5, arcColor)
	ebitenutil.DrawLine(screen, screenX, screenY, screenX+math.Cos(endRad)*radius*0.5, screenY+math.Sin(endRad)*radius*0.5, arcColor)
}

// performWeaponAttack executes a weapon swing attack
func (g *Game) performWeaponAttack() {
	if g.weaponSystem == nil || g.zombieSpawner == nil {
		return
	}

	// Get player center position
	playerX, playerY := g.player.GetCenter()

	// Get mouse position in world coordinates
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY

	// Calculate weapon damage based on equipped weapon
	damage := 5.0 // Base unarmed damage
	selectedItem := g.inventory.GetSelectedItem()
	if selectedItem != nil && selectedItem.Type != items.NONE {
		props := items.ItemDefinitions[selectedItem.Type]
		if props != nil && props.IsWeapon {
			damage = props.WeaponDamage
		}
	}

	// Perform attack
	results := g.weaponSystem.PerformAttack(playerX, playerY, mouseWorldX, mouseWorldY, damage, g.zombieSpawner.Zombies)

	// Apply damage to hit zombies and show indicators
	for _, result := range results {
		if result.Hit {
			// Find the zombie that was hit and apply damage
			for _, zombie := range g.zombieSpawner.Zombies {
				if zombie.IsAlive &&
					math.Abs(zombie.X+zombie.Width/2-result.HitX) < 10 &&
					math.Abs(zombie.Y+zombie.Height/2-result.HitY) < 10 {
					// Apply damage
					zombie.TakeDamage(result.Damage)

					// Show damage indicator with appropriate tier color
					if g.damageIndicators != nil {
						var tier ui.DamageTier
						switch result.Tier {
						case combat.CritTierPurple:
							tier = ui.TierPurple // Fatal - instant death
						case combat.CritTierRed:
							tier = ui.TierRed // Severe damage
						case combat.CritTierYellow:
							tier = ui.TierYellow // Moderate damage
						default:
							tier = ui.TierGreen // Low damage
						}
						g.damageIndicators.SpawnDamageIndicator(zombie.X, zombie.Y, result.Damage, tier, result.IsCritical)
					}

					break
				}
			}
		}
	}
}

// respawnPlayer respawns the player at a safe location
func (g *Game) respawnPlayer() {
	// Reset player position (spawn at world origin or safe location)
	g.player.X = 0
	g.player.Y = 0
	g.player.VX = 0
	g.player.VY = 0

	// Restore health
	g.player.Health = g.player.MaxHealth
	if g.healthSystem != nil {
		// Heal all body parts
		for i := 0; i < 6; i++ {
			g.healthSystem.HealBodyPart(health.BodyPart(i), 9999) // Full heal
		}
	}

	// Restore survival stats
	if g.survivalManager != nil {
		g.survivalManager.Hunger = g.survivalManager.MaxHunger * 0.5 // 50% hunger on respawn
		g.survivalManager.Thirst = g.survivalManager.MaxThirst * 0.5
		g.survivalManager.Stamina = g.survivalManager.MaxStamina
		g.survivalManager.IsStarving = false
		g.survivalManager.IsDehydrated = false
	}

	// Clear death state
	g.isDead = false

	// In hardcore mode, we might want to drop inventory or reset progress
	// For now, just respawn with current inventory

	log.Printf("Player respawned at (%.0f, %.0f)", g.player.X, g.player.Y)
}

// returnToMainMenu returns the player to the main menu
func (g *Game) returnToMainMenu() {
	// Save game before returning to menu
	if err := g.SaveGame(); err != nil {
		log.Printf("Failed to save game before returning to menu: %v", err)
	}

	// Reset game state
	g.inGame = false
	g.inMenu = true
	g.isDead = false

	// Stop auto-save
	g.StopAutoSave()

	log.Printf("Returning to main menu")
}

// checkPlayerDeath checks if player has died and triggers death screen
func (g *Game) checkPlayerDeath() {
	if g.isDead {
		return
	}

	// Check if player health is 0
	if g.player.Health <= 0 || (g.healthSystem != nil && g.healthSystem.OverallHealth <= 0) {
		g.isDead = true

		// Determine cause of death
		cause := "Unknown"
		if g.survivalManager != nil && g.survivalManager.IsStarving {
			cause = "Starvation"
		} else if g.survivalManager != nil && g.survivalManager.IsDehydrated {
			cause = "Dehydration"
		} else {
			cause = "Killed by Zombie"
		}

		// Trigger death screen
		if g.deathScreen != nil {
			g.deathScreen.Trigger(cause)
		}

		log.Printf("Player died: %s", cause)
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

// createSaveState creates a save state from the current game state
func (g *Game) createSaveState() *save.GameState {
	return &save.GameState{
		World:     g.world,
		Player:    g.player,
		Inventory: g.inventory,
		CameraX:   g.cameraX,
		CameraY:   g.cameraY,
		// Crafting state
		CraftingStation: g.CurrentCraftingStation,
		UnlockedRecipes: g.getUnlockedRecipesAsSlice(),

		// Statistics tracking
		BlocksPlaced:    g.BlocksPlaced,
		BlocksDestroyed: g.BlocksDestroyed,
		ItemsCrafted:    g.ItemsCrafted,
		PlayTime:        g.TotalPlayTime.Seconds(),

		// Survival systems
		SurvivalManager: g.survivalManager,
		EquipmentSet:    g.equipmentSet,
		HealthSystem:    g.healthSystem,

		// Enemy systems
		ZombieSpawner: g.zombieSpawner,

		// Storage systems
		ChestManager: g.chestManager,
	}
}

// getUnlockedRecipesAsSlice converts the unlocked recipes map to a slice
func (g *Game) getUnlockedRecipesAsSlice() []string {
	recipes := make([]string, 0, len(g.UnlockedRecipes))
	for recipeID := range g.UnlockedRecipes {
		recipes = append(recipes, recipeID)
	}
	return recipes
}

// SaveGame saves the current game state
func (g *Game) SaveGame() error {
	if g.saveManager == nil {
		return fmt.Errorf("save manager not initialized")
	}

	// Save main game state
	if err := g.saveManager.SaveGame(g.createSaveState()); err != nil {
		return err
	}

	// Save chests
	if g.chestManager != nil {
		if err := g.chestManager.SaveChests(); err != nil {
			log.Printf("Failed to save chests: %v", err)
			return err
		}
	}

	return nil
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

// handleFootstepAudio handles footstep sounds based on player movement
func (g *Game) handleFootstepAudio() {
	// Only play footsteps if player is moving and on ground
	if (g.player.MovingLeft || g.player.MovingRight) && g.player.OnGround && !g.player.IsFlying {
		// Check if enough time has passed since last footstep
		if time.Since(g.lastFootstepTime) >= g.footstepCooldown {
			// Determine surface type based on player position
			surfaceType := g.getSurfaceTypeAtPlayer()
			g.soundLibrary.PlayFootstep(surfaceType)
			g.lastFootstepTime = time.Now()
		}
	}
}

// getSurfaceTypeAtPlayer determines the surface type under the player
func (g *Game) getSurfaceTypeAtPlayer() string {
	// Check the block directly under the player
	playerX, playerY := g.player.GetCenter()
	groundY := playerY + g.player.Height/2 + 5 // Slightly below player center

	hex := g.world.GetHexagonAt(playerX, groundY)
	if hex == nil || hex.BlockType == blocks.AIR {
		return "air" // In air, no footstep sound
	}

	// Convert block type to surface type
	switch hex.BlockType {
	case blocks.GRASS:
		return "grass"
	case blocks.DIRT:
		return "dirt"
	case blocks.STONE, blocks.BEDROCK, blocks.BRICK:
		return "stone"
	case blocks.SAND:
		return "sand"
	case blocks.LOG:
		return "wood"
	default:
		return "stone" // Default to stone
	}
}

// playBlockSound plays a sound for block interactions
func (g *Game) playBlockSound(action string, blockType blocks.BlockType) {
	blockName := getBlockKeyFromType(blockType)
	g.soundLibrary.PlayBlockSound(action, blockName)
}

// playItemSound plays a sound for item interactions
func (g *Game) playItemSound(action string) {
	g.soundLibrary.PlayItemSound(action)
}

// playUISound plays a sound for UI interactions
func (g *Game) playUISound(action string) {
	g.soundLibrary.PlayUISound(action)
}

// playCraftingSound plays a sound for crafting interactions
func (g *Game) playCraftingSound(action string) {
	g.soundLibrary.PlayCraftingSound(action)
}

// playMusic plays background music based on context
func (g *Game) playMusic(context string) {
	g.soundLibrary.PlayMusic(context)
}

// updateAudioContext updates the audio context based on game state
func (g *Game) updateAudioContext() {
	// Determine current biome (simplified)
	biome := biomes.PLAINS // Default biome
	if g.player.Y < -500 {
		biome = biomes.PLAINS // Keep as plains for underground for now
	}

	// Determine if underground
	isUnderground := g.player.Y < 0

	// Determine time of day
	timeOfDay := "day"
	if g.dayNightCycle != nil {
		// This would need to be implemented in day/night cycle
		// timeOfDay = g.dayNightCycle.GetTimeOfDay()
	}

	// Get weather from weather system
	currentWeather := "clear" // Default weather
	if g.weatherSystem != nil {
		currentWeather = g.weatherSystem.GetCurrentWeather()
	}

	g.soundLibrary.UpdateContext(biome, currentWeather, isUnderground, timeOfDay)
}

// TUI model for Bubble Tea
type model struct {
	choices       []string
	cursor        int
	selected      int
	width         int
	height        int
	currentScreen string // "main", "worlds", "settings", "skin_editor", "plugins"
	shouldExit    bool

	// Settings state
	soundEnabled    bool
	graphicsQuality string // "low", "medium", "high"
	difficulty      string // "easy", "normal", "hard"

	// Multi-world state
	worlds        []string
	worldSeeds    map[string]int64 // Map world names to their seeds
	selectedWorld int
	newWorldName  string
	newWorldSeed  string // Seed as string input (empty = random)

	// Skin editor state
	skinColors    []string
	selectedColor int
	skinName      string

	// Plugin manager state
	plugins        []PluginInfo
	selectedPlugin int
}

// PluginInfo holds plugin information for the plugin manager
type PluginInfo struct {
	Name        string
	Version     string
	Enabled     bool
	Description string
}

// Init implements tea.Model interface
func (m model) Init() tea.Cmd {
	return nil
}

// Initial model for TUI
func initialModel() model {
	// Initialize with sample worlds and their seeds (like Minecraft)
	worldSeeds := map[string]int64{
		"World 1": 12345, // Known seed for sharing
		"World 2": 67890, // Another known seed
	}

	return model{
		choices:         []string{"Play (Select World)", "Create New World", "Skin Editor", "Plugin Manager", "Settings", "Exit"},
		cursor:          0,
		selected:        0,
		width:           80,
		height:          20,
		currentScreen:   "main",
		shouldExit:      false,
		soundEnabled:    true,
		graphicsQuality: "medium",
		difficulty:      "normal",
		worlds:          []string{"World 1", "World 2"}, // Default worlds
		worldSeeds:      worldSeeds,
		selectedWorld:   0,
		newWorldSeed:    "", // Empty = random seed
		skinColors:      []string{"Default", "Red", "Blue", "Green", "Purple"},
		selectedColor:   0,
		skinName:        "Player",
		plugins: []PluginInfo{
			{Name: "Minimap", Version: "1.0", Enabled: true, Description: "Shows a minimap"},
			{Name: "Auto-Save", Version: "2.1", Enabled: true, Description: "Auto-saves every 5 minutes"},
			{Name: "Debug Tools", Version: "0.5", Enabled: false, Description: "Developer debugging tools"},
		},
		selectedPlugin: 0,
	}
}

// TUI update function for Bubble Tea
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.currentScreen == "main" {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter", " ":
				m.selected = m.cursor
				choice := m.choices[m.cursor]
				if choice == "Settings" {
					m.currentScreen = "settings"
					m.choices = []string{"Sound: ON", "Graphics: Medium", "Difficulty: Normal", "Back"}
					m.cursor = 0
				} else if choice == "Exit" {
					fmt.Println("Goodbye!")
					os.Exit(0)
				} else if choice == "Play (Select World)" {
					m.currentScreen = "worlds"
					m.cursor = 0
				} else if choice == "Create New World" {
					m.currentScreen = "create_world"
					m.newWorldName = ""
					m.cursor = 0
				} else if choice == "Skin Editor" {
					m.currentScreen = "skin_editor"
					m.cursor = 0
				} else if choice == "Plugin Manager" {
					m.currentScreen = "plugins"
					m.cursor = 0
				} else if choice == "Back" {
					// This shouldn't happen in main screen
				} else {
					// Start Game
					return m, tea.Quit
				}
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		} else if m.currentScreen == "settings" {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter", " ":
				if m.choices[m.cursor] == "Sound: ON" || m.choices[m.cursor] == "Sound: OFF" {
					m.soundEnabled = !m.soundEnabled
					if m.soundEnabled {
						m.choices[0] = "Sound: ON"
					} else {
						m.choices[0] = "Sound: OFF"
					}
				} else if m.choices[m.cursor] == "Graphics: Low" || m.choices[m.cursor] == "Graphics: Medium" || m.choices[m.cursor] == "Graphics: High" {
					// Cycle through graphics quality
					if m.graphicsQuality == "low" {
						m.graphicsQuality = "medium"
						m.choices[1] = "Graphics: Medium"
					} else if m.graphicsQuality == "medium" {
						m.graphicsQuality = "high"
						m.choices[1] = "Graphics: High"
					} else {
						m.graphicsQuality = "low"
						m.choices[1] = "Graphics: Low"
					}
				} else if m.choices[m.cursor] == "Difficulty: Easy" || m.choices[m.cursor] == "Difficulty: Normal" || m.choices[m.cursor] == "Difficulty: Hard" {
					// Cycle through difficulty
					if m.difficulty == "easy" {
						m.difficulty = "normal"
						m.choices[2] = "Difficulty: Normal"
					} else if m.difficulty == "normal" {
						m.difficulty = "hard"
						m.choices[2] = "Difficulty: Hard"
					} else {
						m.difficulty = "easy"
						m.choices[2] = "Difficulty: Easy"
					}
				} else if m.choices[m.cursor] == "Back" {
					m.currentScreen = "main"
					m.choices = []string{"Play (Select World)", "Create New World", "Skin Editor", "Plugin Manager", "Settings", "Exit"}
					m.cursor = 4 // Return to Settings option
				}
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		} else if m.currentScreen == "worlds" {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.worlds) {
					m.cursor++
				}
			case "enter", " ":
				if m.cursor < len(m.worlds) {
					m.selectedWorld = m.cursor
					return m, tea.Quit // Start game with selected world
				} else if m.cursor == len(m.worlds) { // Back option
					m.currentScreen = "main"
					m.cursor = 0
				}
			case "ctrl+c", "q":
				m.currentScreen = "main"
				m.cursor = 0
			}
		} else if m.currentScreen == "create_world" {
			switch msg.String() {
			case "tab":
				// Toggle between name and seed fields
				if m.cursor == 0 {
					m.cursor = 1
				} else {
					m.cursor = 0
				}
			case "enter", " ":
				if m.newWorldName != "" {
					// Generate or parse seed
					var seed int64
					if m.newWorldSeed == "" {
						// Random seed
						seed = time.Now().UnixNano()
					} else {
						// Try to parse seed as integer
						parsedSeed, err := strconv.ParseInt(m.newWorldSeed, 10, 64)
						if err != nil {
							// Use string hash as seed
							h := fnv.New64a()
							h.Write([]byte(m.newWorldSeed))
							seed = int64(h.Sum64())
						} else {
							seed = parsedSeed
						}
					}

					// Add world with seed
					m.worlds = append(m.worlds, m.newWorldName)
					if m.worldSeeds == nil {
						m.worldSeeds = make(map[string]int64)
					}
					m.worldSeeds[m.newWorldName] = seed

					m.currentScreen = "worlds"
					m.cursor = len(m.worlds) - 1
				}
			case "ctrl+c", "q":
				m.currentScreen = "main"
				m.cursor = 1
			}
		} else if m.currentScreen == "skin_editor" {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.skinColors)+1 {
					m.cursor++
				}
			case "enter", " ":
				if m.cursor < len(m.skinColors) {
					m.selectedColor = m.cursor
				} else if m.cursor == len(m.skinColors) { // Save & Back
					m.currentScreen = "main"
					m.cursor = 2
				}
			case "ctrl+c", "q":
				m.currentScreen = "main"
				m.cursor = 2
			}
		} else if m.currentScreen == "plugins" {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.plugins) {
					m.cursor++
				}
			case "enter", " ":
				if m.cursor < len(m.plugins) {
					// Toggle plugin enabled state
					m.plugins[m.cursor].Enabled = !m.plugins[m.cursor].Enabled
				} else if m.cursor == len(m.plugins) { // Back
					m.currentScreen = "main"
					m.cursor = 3
				}
			case "ctrl+c", "q":
				m.currentScreen = "main"
				m.cursor = 3
			}
		}
	}

	return m, nil
}

// TUI view function for Bubble Tea
func (m model) View() string {
	s := strings.Builder{}

	// Determine title based on screen
	var title string
	switch m.currentScreen {
	case "main":
		title = "\033[38;5;196m╔════════════════════════════════════════════════════════════════════╗\033[0m\n"
		title += "\033[38;5;196m║\033[38;5;208m  ╔══╗  \033[38;5;226mT\033[38;5;46mE\033[38;5;51mS\033[38;5;99mS\033[38;5;129mE\033[38;5;165mL\033[38;5;196mB\033[38;5;208mO\033[38;5;226mX  ╔══╗  \033[38;5;196m║\033[0m\n"
		title += "\033[38;5;196m║\033[38;5;208m  ║  ║  \033[38;5;46m║\033[38;5;51m║\033[38;5;99m║\033[38;5;129m║\033[38;5;165m║\033[38;5;196m║\033[38;5;208m║\033[38;5;226m║  ║  ║  \033[38;5;196m║\033[0m\n"
		title += "\033[38;5;196m║\033[38;5;208m  ╚══╝  ╚══╝  \033[38;5;46m╚══╝\033[38;5;51m╚══╝\033[38;5;99m╚══╝\033[38;5;129m╚══╝\033[38;5;165m╚══╝\033[38;5;196m╚══╝\033[38;5;208m╚══╝\033[38;5;226m╚══╝  ║  ║  \033[38;5;196m║\033[0m\n"
		title += "\033[38;5;196m╚══════════════════════════════════════════════════════════════════════════╝\033[0m\n"
	case "settings":
		title = "\033[38;5;208m╔══════════════════════════════════════════════════════════╗\033[0m\n"
		title += "\033[38;5;208m║\033[38;5;226m  ⚙ \033[38;5;46mS\033[38;5;51mE\033[38;5;129mT\033[38;5;165mT\033[38;5;196mI\033[38;5;208mN\033[38;5;226mG\033[38;5;208mS  ⚙  \033[38;5;208m║\033[0m\n"
		title += "\033[38;5;208m║\033[38;5;226m                                             ║  \033[38;5;208m║\033[0m\n"
		title += "\033[38;5;208m╚══════════════════════════════════════════════════════════════════════════════╝\033[0m\n"
	case "worlds":
		title = "\033[38;5;51m╔══════════════════════════════════════════════════════════╗\033[0m\n"
		title += "\033[38;5;51m║\033[38;5;46m  🌍 \033[38;5;99mS\033[38;5;129mE\033[38;5;165mL\033[38;5;196mE\033[38;5;208mC\033[38;5;226mT\033[38;5;46m \033[38;5;51mW\033[38;5;99mO\033[38;5;129mR\033[38;5;165mL\033[38;5;196mD\033[38;5;208m  🌍  \033[38;5;51m║\033[0m\n"
		title += "\033[38;5;51m╚══════════════════════════════════════════════════════════════════════════════╝\033[0m\n"
	case "create_world":
		title = "\033[38;5;46m╔══════════════════════════════════════════════════════════╗\033[0m\n"
		title += "\033[38;5;46m║\033[38;5;51m  ✨ \033[38;5;99mC\033[38;5;129mR\033[38;5;165mE\033[38;5;196mA\033[38;5;208mT\033[38;5;226mE\033[38;5;46m \033[38;5;51mW\033[38;5;99mO\033[38;5;129mR\033[38;5;165mL\033[38;5;196mD\033[38;5;208m  ✨  \033[38;5;46m║\033[0m\n"
		title += "\033[38;5;46m╚══════════════════════════════════════════════════════════════════════════════╝\033[0m\n"
	case "skin_editor":
		title = "\033[38;5;165m╔══════════════════════════════════════════════════════════╗\033[0m\n"
		title += "\033[38;5;165m║\033[38;5;196m  🎨 \033[38;5;208mS\033[38;5;226mK\033[38;5;46mI\033[38;5;51mN\033[38;5;99m \033[38;5;129mE\033[38;5;165mD\033[38;5;196mI\033[38;5;208mT\033[38;5;226mO\033[38;5;46mR\033[38;5;165m  🎨  \033[38;5;165m║\033[0m\n"
		title += "\033[38;5;165m╚══════════════════════════════════════════════════════════════════════════════╝\033[0m\n"
	case "plugins":
		title = "\033[38;5;226m╔══════════════════════════════════════════════════════════╗\033[0m\n"
		title += "\033[38;5;226m║\033[38;5;46m  🔌 \033[38;5;51mP\033[38;5;99mL\033[38;5;129mU\033[38;5;165mG\033[38;5;196mI\033[38;5;208mN\033[38;5;226m \033[38;5;46mM\033[38;5;51mA\033[38;5;99mN\033[38;5;129mA\033[38;5;165mG\033[38;5;196mE\033[38;5;208mR\033[38;5;226m  🔌  \033[38;5;226m║\033[0m\n"
		title += "\033[38;5;226m╚══════════════════════════════════════════════════════════════════════════════╝\033[0m\n"
	default:
		title = "\033[38;5;196m╔══════════════════════════════════════════════════════════╗\033[0m\n"
		title += "\033[38;5;196m║\033[38;5;226m  TESSELBOX  \033[38;5;196m║\033[0m\n"
		title += "\033[38;5;196m╚══════════════════════════════════════════════════════════════════════════════╝\033[0m\n"
	}

	titlePadding := (m.width - 60) / 2
	if titlePadding < 0 {
		titlePadding = 0
	}
	s.WriteString(strings.Repeat(" ", titlePadding))
	s.WriteString(title)
	s.WriteString(strings.Repeat(" ", titlePadding))
	s.WriteString("\n\n")

	// Render content based on screen
	switch m.currentScreen {
	case "main", "settings":
		// Menu items with color highlighting
		for i, choice := range m.choices {
			padding := (m.width - len(choice) - 4) / 2
			s.WriteString(strings.Repeat(" ", padding))

			if m.cursor == i {
				if m.currentScreen == "main" {
					s.WriteString("\033[38;5;46;1m") // Bright green bold
				} else {
					s.WriteString("\033[38;5;208;1m") // Bright orange bold
				}
				s.WriteString("👉 ")
				s.WriteString(choice)
				s.WriteString("\033[0m")
			} else {
				if strings.Contains(choice, "ON") {
					s.WriteString("\033[38;5;46m")
				} else if strings.Contains(choice, "OFF") {
					s.WriteString("\033[38;5;196m")
				} else if strings.Contains(choice, "High") {
					s.WriteString("\033[38;5;51m")
				} else if strings.Contains(choice, "Low") {
					s.WriteString("\033[38;5;208m")
				} else if strings.Contains(choice, "Hard") {
					s.WriteString("\033[38;5;196m")
				} else if strings.Contains(choice, "Easy") {
					s.WriteString("\033[38;5;46m")
				} else {
					s.WriteString("\033[38;5;250m")
				}
				s.WriteString("  ")
				s.WriteString(choice)
				s.WriteString("\033[0m")
			}
			s.WriteString("\n")
		}

	case "worlds":
		// World selection list with seeds
		s.WriteString("  Available Worlds:\n\n")
		for i, worldName := range m.worlds {
			// Get seed for this world
			seed := int64(0)
			if m.worldSeeds != nil {
				if s, ok := m.worldSeeds[worldName]; ok {
					seed = s
				}
			}

			// Format: "World Name (seed: 12345)"
			seedStr := fmt.Sprintf(" (seed: %d)", seed)
			fullText := worldName + seedStr

			padding := (m.width - len(fullText) - 4) / 2
			s.WriteString(strings.Repeat(" ", padding))

			if m.cursor == i {
				s.WriteString("\033[38;5;46;1m👉 ")
				s.WriteString(worldName)
				s.WriteString("\033[38;5;240m" + seedStr + "\033[0m")
			} else {
				s.WriteString("   \033[38;5;250m")
				s.WriteString(worldName)
				s.WriteString("\033[38;5;240m" + seedStr + "\033[0m")
			}
			s.WriteString("\n")
		}
		// Back option
		padding := (m.width - 6) / 2
		s.WriteString(strings.Repeat(" ", padding))
		if m.cursor == len(m.worlds) {
			s.WriteString("\033[38;5;208;1m👉 Back\033[0m")
		} else {
			s.WriteString("\033[38;5;250m   Back\033[0m")
		}
		s.WriteString("\n")

	case "create_world":
		// Create world screen with seed input
		s.WriteString("\n")

		// World name field
		nameLabel := "World Name:"
		namePadding := (m.width - len(nameLabel) - 25) / 2
		s.WriteString(strings.Repeat(" ", namePadding))
		if m.cursor == 0 {
			s.WriteString("\033[38;5;46;1m👉 " + nameLabel + "\033[0m ")
		} else {
			s.WriteString("   " + nameLabel + " ")
		}
		s.WriteString(m.newWorldName)
		if m.cursor == 0 {
			s.WriteString("_") // Cursor indicator
		}
		s.WriteString("\n\n")

		// Seed field
		seedLabel := "Seed (optional):"
		seedDisplay := m.newWorldSeed
		if seedDisplay == "" {
			seedDisplay = "(random)"
		}
		seedPadding := (m.width - len(seedLabel) - 25) / 2
		s.WriteString(strings.Repeat(" ", seedPadding))
		if m.cursor == 1 {
			s.WriteString("\033[38;5;208;1m👉 " + seedLabel + "\033[0m ")
		} else {
			s.WriteString("   " + seedLabel + " ")
		}
		if m.newWorldSeed == "" {
			s.WriteString("\033[38;5;240m" + seedDisplay + "\033[0m")
		} else {
			s.WriteString(seedDisplay)
		}
		if m.cursor == 1 {
			s.WriteString("_") // Cursor indicator
		}
		s.WriteString("\n\n")

		// Show some example seeds
		s.WriteString("\033[38;5;240m")
		seedHint := "Known seeds: 12345, 67890, 42, 1337, 69420"
		hintPadding := (m.width - len(seedHint)) / 2
		s.WriteString(strings.Repeat(" ", hintPadding))
		s.WriteString(seedHint)
		s.WriteString("\033[0m\n\n")

		// Instructions
		instr := "Tab: Switch fields  •  Enter: Create  •  q: Cancel"
		padding2 := (m.width - len(instr)) / 2
		s.WriteString(strings.Repeat(" ", padding2))
		s.WriteString("\033[38;5;240m" + instr + "\033[0m")
		s.WriteString("\n")

	case "skin_editor":
		// Skin color selection
		s.WriteString("  Select Skin Color:\n\n")
		for i, color := range m.skinColors {
			padding := (m.width - len(color) - 4) / 2
			s.WriteString(strings.Repeat(" ", padding))

			var colorCode string
			switch color {
			case "Red":
				colorCode = "\033[38;5;196m"
			case "Blue":
				colorCode = "\033[38;5;51m"
			case "Green":
				colorCode = "\033[38;5;46m"
			case "Purple":
				colorCode = "\033[38;5;165m"
			default:
				colorCode = "\033[38;5;226m"
			}

			if m.cursor == i {
				s.WriteString("\033[38;5;46;1m👉 " + colorCode + color + "\033[0m")
			} else if m.selectedColor == i {
				s.WriteString("   " + colorCode + color + "\033[0m")
			} else {
				s.WriteString("   \033[38;5;250m" + color + "\033[0m")
			}
			s.WriteString("\n")
		}
		// Save & Back option
		padding := (m.width - 14) / 2
		s.WriteString(strings.Repeat(" ", padding))
		if m.cursor == len(m.skinColors) {
			s.WriteString("\033[38;5;46;1m👉 Save & Back\033[0m")
		} else {
			s.WriteString("   \033[38;5;250mSave & Back\033[0m")
		}
		s.WriteString("\n")

	case "plugins":
		// Plugin list
		s.WriteString("  Installed Plugins:\n\n")
		for i, plugin := range m.plugins {
			padding := (m.width - len(plugin.Name) - 15) / 2
			s.WriteString(strings.Repeat(" ", padding))

			if m.cursor == i {
				s.WriteString("\033[38;5;46;1m👉 ")
			} else {
				s.WriteString("   ")
			}

			if plugin.Enabled {
				s.WriteString("\033[38;5;46m[ON]\033[0m ")
			} else {
				s.WriteString("\033[38;5;196m[OFF]\033[0m ")
			}

			s.WriteString(plugin.Name)
			s.WriteString(" \033[38;5;240mv" + plugin.Version + "\033[0m")
			if m.cursor == i {
				s.WriteString("\033[0m")
			}
			s.WriteString("\n")
		}
		// Back option
		padding := (m.width - 6) / 2
		s.WriteString(strings.Repeat(" ", padding))
		if m.cursor == len(m.plugins) {
			s.WriteString("\033[38;5;208;1m👉 Back\033[0m")
		} else {
			s.WriteString("\033[38;5;250m   Back\033[0m")
		}
		s.WriteString("\n")
	}

	// Footer with instructions
	s.WriteString("\n")
	var footer string
	switch m.currentScreen {
	case "main":
		footer = "↑/k: Move    ↓/j: Move    Enter: Select    q/Ctrl+C: Quit"
	case "worlds", "plugins":
		footer = "↑/k: Move    ↓/j: Move    Enter: Select    q/Ctrl+C: Back"
	case "create_world":
		footer = "Type: Name    Enter: Create    q: Cancel"
	default:
		footer = "↑/k: Change    ↓/j: Change    Enter: Toggle    q/Ctrl+C: Back"
	}
	footerPadding := (m.width - len(footer)) / 2
	s.WriteString(strings.Repeat(" ", footerPadding))
	s.WriteString("\033[38;5;240m") // Dim gray
	s.WriteString(footer)
	s.WriteString("\033[0m") // Reset color

	return s.String()
}

// runTUI provides beautiful terminal interface using Bubble Tea
func runTUI() {
	p := tea.NewProgram(initialModel())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running TUI: %v", err)
		os.Exit(1)
	}

	m := finalModel.(model)

	// Handle game start from different screens
	worldName := "default"
	worldSeed := int64(0)

	if m.currentScreen == "worlds" && m.selectedWorld >= 0 && m.selectedWorld < len(m.worlds) {
		// User selected a world from the list
		worldName = m.worlds[m.selectedWorld]
		if m.worldSeeds != nil {
			if seed, ok := m.worldSeeds[worldName]; ok {
				worldSeed = seed
			}
		}
		startGameWithGUI(worldName, worldSeed)
	} else if m.currentScreen == "main" && m.selected == 0 {
		// Legacy: Start Game from main menu (use first world or default)
		if len(m.worlds) > 0 {
			worldName = m.worlds[0]
			if m.worldSeeds != nil {
				if seed, ok := m.worldSeeds[worldName]; ok {
					worldSeed = seed
				}
			}
		}
		startGameWithGUI(worldName, worldSeed)
	}
}

// startGameWithGUI starts the game with Ebiten GUI engine
func startGameWithGUI(worldName string, worldSeed int64) {
	fmt.Printf("🚀 Starting game with world '%s' (seed: %d)...\n", worldName, worldSeed)

	// Load block definitions before any world generation
	blocks.LoadBlocks()

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Tesselbox v2.0 - Hexagon Sandbox")
	ebiten.SetTPS(FPS)

	// Enable input
	ebiten.SetCursorMode(ebiten.CursorModeVisible)

	// Create game with specified world and seed
	game := NewGameWithWorld(worldName, worldSeed)

	// Start auto-saver
	game.StartAutoSave()

	// Run the game with Ebiten engine
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Failed to run game: %v", err)
	}
}

// startGameCLI starts the game in CLI mode (for testing only)
func startGameCLI() {
	worldName := "default"
	if len(os.Args) > 2 {
		worldName = os.Args[2]
	}

	fmt.Printf("Starting game in world: %s\n", worldName)
	fmt.Println("Game running... Press Ctrl+C to exit")

	// Create game instance without GUI
	game := NewGame()
	game.inMenu = false
	game.inGame = true

	// Simple game loop
	for {
		// Update game state
		if err := game.Update(); err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}

		// Simple status display
		px, py := game.player.GetCenter()
		fmt.Printf("\rPosition: (%.1f, %.1f)  Chunks: %d", px, py, len(game.world.Chunks))

		time.Sleep(50 * time.Millisecond)
	}
}

// listWorldsCLI lists available worlds
func listWorldsCLI() {
	fmt.Println("Available worlds:")

	// Check tesselbox directory
	tesselboxDir := getTesselboxDir()
	worldDir := filepath.Join(tesselboxDir, "worlds")
	if _, err := os.Stat(worldDir); os.IsNotExist(err) {
		fmt.Println("No worlds found")
		return
	}

	files, err := os.ReadDir(worldDir)
	if err != nil {
		fmt.Printf("Error reading worlds: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No worlds found")
		return
	}

	for i, file := range files {
		if file.IsDir() {
			fmt.Printf("%d. %s\n", i+1, file.Name())
		}
	}
}

// createWorldCLI creates a new world
func createWorldCLI() {
	worldName := "NewWorld"
	if len(os.Args) > 3 {
		worldName = os.Args[3]
	}

	fmt.Printf("Creating world: %s\n", worldName)

	// Create world directory
	tesselboxDir := getTesselboxDir()
	worldDir := filepath.Join(tesselboxDir, "worlds", worldName)
	if err := os.MkdirAll(worldDir, 0755); err != nil {
		fmt.Printf("Error creating world: %v\n", err)
		return
	}

	fmt.Printf("World '%s' created successfully!\n", worldName)
}

// deleteWorldCLI deletes a world
func deleteWorldCLI() {
	worldName := "NewWorld"
	if len(os.Args) > 3 {
		worldName = os.Args[3]
	}

	fmt.Printf("Deleting world: %s\n", worldName)

	// Delete world directory
	tesselboxDir := getTesselboxDir()
	worldDir := filepath.Join(tesselboxDir, "worlds", worldName)
	if err := os.RemoveAll(worldDir); err != nil {
		fmt.Printf("Error deleting world: %v\n", err)
		return
	}

	fmt.Printf("World '%s' deleted successfully!\n", worldName)
}

// cleanupAudio cleans up audio resources when shutting down
func (g *Game) cleanupAudio() {
	if g.audioManager != nil {
		g.audioManager.Cleanup()
	}
}

// Main function
func main() {
	// Default to TUI mode
	fmt.Println("🎮 TESSELBOX - Beautiful Terminal Interface 🎮")
	runTUI()
}
