package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"tesselbox/pkg/audio"
	"tesselbox/pkg/blocks"
	"tesselbox/pkg/commands"
	"tesselbox/pkg/crafting"
	"tesselbox/pkg/entities"
	"tesselbox/pkg/gametime"
	"tesselbox/pkg/hexagon"
	"tesselbox/pkg/input"
	"tesselbox/pkg/items"
	"tesselbox/pkg/menu"
	"tesselbox/pkg/player"
	"tesselbox/pkg/save"
	"tesselbox/pkg/settings"
	"tesselbox/pkg/weather"
	"tesselbox/pkg/world"

	"github.com/hajimehoshi/ebiten/v2"
)

// Constants for game configuration
const (
	PoolSize             = 100
	PoolCapacity         = 200
	AutoSaveInterval     = 5 * time.Minute
	DayNightCycleDuration = 20 * time.Minute
	CollisionRange       = 40.0
)

// Game represents the main game state
type Game struct {
	// Core systems
	world         *world.World
	player        *player.Player
	cameraX       float64
	cameraY       float64
	groundLevel   float64
	creativeMode  bool

	// Game state
	inMenu        bool
	inGame        bool
	showInventory bool
	showDebug     bool

	// Systems
	craftingSystem *crafting.CraftingSystem
	craftingUI     *crafting.CraftingUI
	menu           *menu.Menu
	pluginManager  *entities.EnhancedPluginManager
	inputManager   *input.InputManager
	chatHandler    *commands.ChatHandler
	pluginWatcher  *commands.PluginWatcher

	// Save system
	saveManager *save.SaveManager
	autoSaver   *save.AutoSaver

	// Settings system
	settingsManager *settings.Manager

	// Audio system
	audioManager *audio.Manager

	// Block selection
	selectedBlock string

	// Visual properties
	whiteImage    *ebiten.Image
	playerBodyColor color.RGBA
	playerHeadColor color.RGBA

	// Day/Night cycle
	dayNightCycle *gametime.DayNightCycle
	weatherSystem *weather.WeatherSystem

	// Entity system
	entityManager *entities.EntityManager

	// Performance metrics
	lastFPS       float64
	lastTPS       float64

	// Dropped items
	droppedItems []*items.DroppedItem
}

// NewGame creates a new game instance
func NewGame(creativeMode bool) *Game {
	g := &Game{
		creativeMode:  creativeMode,
		inMenu:        true,
		inGame:        false,
		showInventory: false,
		showDebug:     false,
		selectedBlock: "grass",
		groundLevel:   400,
		playerBodyColor: color.RGBA{100, 150, 200, 255},
		playerHeadColor: color.RGBA{255, 220, 180, 255},
	}

	// Create white image for rendering
	g.whiteImage = ebiten.NewImage(1, 1)
	g.whiteImage.Fill(color.White)

	// Initialize world
	g.world = world.NewWorld(0.5)

	// Initialize player
	g.player = player.NewPlayer()

	// Initialize systems
	g.craftingSystem = crafting.NewCraftingSystem()
	g.craftingUI = crafting.NewCraftingUI(g.craftingSystem)

	// Initialize input manager
	g.inputManager = input.NewInputManager()

	// Initialize settings
	g.settingsManager = settings.NewManager("")
	if err := g.settingsManager.Load(); err != nil {
		log.Printf("Warning: Failed to load settings: %v", err)
	}

	// Apply settings
	g.CreativeMode = g.settingsManager.GetBool("creative_mode_default")

	// Initialize audio
	g.audioManager = audio.NewManager(g.settingsManager)
	if err := g.audioManager.Initialize(); err != nil {
		log.Printf("Warning: Failed to initialize audio: %v", err)
	}

	// Initialize menu with settings
	g.menu = menu.NewMenu(g.settingsManager)
	g.menu.CreativeMode = g.CreativeMode
	g.menu.SetMainMenu()

	// Initialize plugin manager
	g.pluginManager = entities.NewEnhancedPluginManager()

	// Initialize chat handler and plugin watcher
	g.chatHandler = commands.NewChatHandler(g.pluginManager)
	g.pluginWatcher = commands.NewPluginWatcher(g.pluginManager, "plugins")
	g.pluginWatcher.Start()

	// Initialize save system
	g.saveManager = save.NewSaveManager("saves")
	g.autoSaver = save.NewAutoSaver(g.saveManager, AutoSaveInterval)

	// Initialize day/night cycle
	g.dayNightCycle = gametime.NewDayNightCycle(DayNightCycleDuration)

	// Initialize weather
	g.weatherSystem = weather.NewWeatherSystem()

	// Initialize entity manager
	g.entityManager = entities.NewEntityManager()

	return g
}

// Update implements ebiten.Game interface
func (g *Game) Update() error {
	if g.inMenu {
		return g.updateMenu()
	}

	if !g.inGame {
		return nil
	}

	// Update day/night cycle
	g.dayNightCycle.Update()

	// Update weather
	g.weatherSystem.Update()

	// Update player
	if err := g.updatePlayer(); err != nil {
		return err
	}

	// Update camera to follow player
	g.updateCamera()

	// Handle input
	g.handleInput()

	// Update entities
	g.entityManager.Update(g.world, g.player)

	// Update dropped items
	g.updateDroppedItems()

	// Auto-save
	g.autoSaver.Update(g)

	return nil
}

// updateMenu handles menu input
func (g *Game) updateMenu() error {
	if action := g.menu.Update(); action != menu.ActionNone {
		g.handleMenuAction(action)
	}
	return nil
}

// updatePlayer updates player physics and state
func (g *Game) updatePlayer() error {
	g.player.Update()

	// Apply gravity
	g.player.ApplyGravity()

	// Check collisions with ground
	if g.player.GetY() > g.groundLevel {
		g.player.SetY(g.groundLevel)
		g.player.SetVelocityY(0)
	}

	return nil
}

// updateCamera updates camera position to follow player
func (g *Game) updateCamera() {
	screenWidth, screenHeight := ebiten.WindowSize()
	g.cameraX = g.player.GetX() - float64(screenWidth)/2
	g.cameraY = g.player.GetY() - float64(screenHeight)/2
}

// handleInput processes game input
func (g *Game) handleInput() {
	g.inputManager.Update()

	// Toggle inventory
	if g.inputManager.IsInventoryToggled() {
		g.showInventory = !g.showInventory
		g.audioManager.PlayInventorySound()
	}

	// Toggle debug
	if g.inputManager.IsDebugToggled() {
		g.showDebug = !g.showDebug
	}

	// Handle block placement/mining
	if g.inputManager.IsPrimaryActionPressed() {
		g.handleMining()
	} else if g.inputManager.IsSecondaryActionPressed() {
		g.handleBlockPlacement()
	}

	// Handle hotbar selection
	g.handleHotbarSelection()

	// Handle jumping
	if g.inputManager.IsJumpPressed() && g.player.IsOnGround() {
		g.player.Jump()
		g.audioManager.PlayJumpSound()
	}
}

// handleHotbarSelection processes hotbar key input
func (g *Game) handleHotbarSelection() {
	keys := []ebiten.Key{
		ebiten.Key1, ebiten.Key2, ebiten.Key3,
		ebiten.Key4, ebiten.Key5, ebiten.Key6,
		ebiten.Key7, ebiten.Key8, ebiten.Key9,
	}
	for i, key := range keys {
		if ebiten.IsKeyPressed(key) {
			g.player.GetInventory().SelectSlot(i)
			g.audioManager.PlayInventorySound()
			break
		}
	}
}

// Draw implements ebiten.Game interface
func (g *Game) Draw(screen *ebiten.Image) {
	if g.inMenu {
		g.menu.Draw(screen)
		return
	}

	// Clear screen with sky color based on time of day
	skyColor := g.dayNightCycle.GetSkyColor()
	screen.Fill(skyColor)

	// Render world
	g.renderWorld(screen)

	// Render entities
	g.renderEntities(screen)

	// Render player
	g.renderPlayer(screen)

	// Render dropped items
	g.renderDroppedItems(screen)

	// Render UI
	g.renderUI(screen)

	// Render debug info
	if g.showDebug {
		g.renderDebugInfo(screen)
	}
}

// Layout implements ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
