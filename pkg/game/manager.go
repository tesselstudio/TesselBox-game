package game

import (
	"image/color"
	"log"
	"time"

	"tesselbox/pkg/audio"
	"tesselbox/pkg/chest"
	"tesselbox/pkg/combat"
	"tesselbox/pkg/crafting"
	"tesselbox/pkg/debug"
	"tesselbox/pkg/enemies"
	"tesselbox/pkg/equipment"
	"tesselbox/pkg/gametime"
	"tesselbox/pkg/health"
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
)

// GameManager coordinates all game subsystems
type GameManager struct {
	// Core systems
	World        *world.World
	Player       *player.Player
	Inventory    *items.Inventory
	StateManager *ui.StateManager

	// Crafting
	CraftingSystem *crafting.CraftingSystem
	CraftingUI     *crafting.CraftingUI

	// Plugins
	PluginManager   *plugins.PluginManager
	PluginUI        *plugins.PluginUI
	PluginInstaller *plugins.PluginInstaller

	// Skin
	SkinEditor *skin.SkinEditor

	// Input
	InputManager *input.InputManager

	// Save system
	SaveManager *save.SaveManager
	AutoSaver   *save.AutoSaver

	// Time systems
	DayNightCycle *gametime.DayNightCycle
	WeatherSystem *weather.WeatherSystem

	// Audio
	AudioManager           *audio.AudioManager
	SoundLibrary           *audio.SoundLibrary
	BackgroundMusicManager *audio.BackgroundMusicManager

	// Debug
	RecoveryHandler *debug.RecoveryHandler
	Profiler        *debug.PerformanceProfiler

	// Survival systems
	SurvivalManager *survival.SurvivalManager
	EquipmentSet    *equipment.EquipmentSet
	HealthSystem    *health.LocationalHealthSystem
	BackpackUI      *ui.BackpackUI
	HUD             *ui.HUD

	// Chest system
	ChestManager *chest.ChestManager
	ChestUI      *ui.ChestUI

	// Combat
	WeaponSystem *combat.WeaponSystem

	// UI effects
	DamageIndicators  *ui.DamageIndicatorManager
	ScreenFlash       *ui.ScreenFlash
	DirectionalHitInd *ui.DirectionalHitManager
	DeathScreen       *ui.DeathScreen

	// Enemy systems
	ZombieSpawner *enemies.ZombieSpawner

	// Game state
	SelectedBlock string
	CreativeMode  bool
	CommandMode   bool
	CommandString string

	// Statistics
	BlocksPlaced    int
	BlocksDestroyed int
	ItemsCrafted    int
	PlayStartTime   time.Time
	TotalPlayTime   time.Duration

	// Crafting station tracking
	CurrentCraftingStation string
	UnlockedRecipes        map[string]bool

	// Rendering
	CameraX, CameraY float64
	CurrentLayer     int
	WhiteImage       *ebiten.Image

	// Timing
	LastTime         time.Time
	LastFootstepTime time.Time
	FootstepCooldown time.Duration
	MiningDamage     float64
	MiningSpeed      float64

	// Dropped items
	DroppedItems []*DroppedItem

	// Object pools for rendering
	VertexPool  [][]ebiten.Vertex
	IndicesPool [][]uint16
	PoolIndex   int

	// Mouse state
	MouseX, MouseY       int
	HoveredBlockName     string
	LeftMouseWasPressed  bool
	RightMouseWasPressed bool

	// Mining state
	IsMining        bool
	MiningProgress  float64
	MiningStartTime time.Time
}

// DroppedItem represents an item dropped in the world
type DroppedItem struct {
	X, Y     float64
	Item     *items.Item
	PickupAt time.Time
}

// NewGameManager creates a new game manager with all subsystems initialized
func NewGameManager(worldName string, worldSeed int64, creativeMode bool) *GameManager {
	log.Printf("Initializing GameManager for world '%s'...", worldName)

	gm := &GameManager{
		World:            world.NewWorld(worldName),
		Player:           player.NewPlayer(400, 300),
		Inventory:        items.NewInventory(32),
		StateManager:     ui.NewStateManager(),
		SelectedBlock:    "dirt",
		CreativeMode:     creativeMode,
		LastTime:         time.Now(),
		LastFootstepTime: time.Now(),
		FootstepCooldown: 400 * time.Millisecond,
		PlayStartTime:    time.Now(),
		UnlockedRecipes:  make(map[string]bool),
		CurrentLayer:     0,
		DroppedItems:     make([]*DroppedItem, 0),
		VertexPool:       make([][]ebiten.Vertex, 10),
		IndicesPool:      make([][]uint16, 10),
	}

	// Initialize white image for rendering
	gm.WhiteImage = ebiten.NewImage(1, 1)
	gm.WhiteImage.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})

	// Initialize object pools
	for i := range gm.VertexPool {
		gm.VertexPool[i] = make([]ebiten.Vertex, 0, 1000)
		gm.IndicesPool[i] = make([]uint16, 0, 1000)
	}

	// Set world seed
	if worldSeed != 0 {
		gm.World.SetSeed(worldSeed)
		log.Printf("World '%s' created with seed: %d", worldName, worldSeed)
	} else {
		log.Printf("World '%s' created with random seed: %d", worldName, gm.World.GetSeed())
	}

	// Find spawn position
	spawnX, spawnY := gm.World.FindSpawnPosition(0, 0)
	gm.Player.SetPosition(spawnX, spawnY)
	gm.World.GetChunksInRange(spawnX, spawnY)
	log.Printf("Player spawned at: (%.1f, %.1f)", spawnX, spawnY)

	// Initialize crafting system
	gm.CraftingSystem = crafting.NewCraftingSystem()
	if err := gm.CraftingSystem.LoadRecipesFromAssets(); err != nil {
		log.Printf("Warning: Failed to load crafting recipes: %v", err)
	}
	gm.CraftingUI = crafting.NewCraftingUI(gm.CraftingSystem, gm.Inventory)

	// Initialize input manager
	gm.InputManager = input.NewInputManager()

	// Load game assets
	items.LoadItems()

	// Add initial items
	gm.Inventory.AddItem(items.DIRT_BLOCK, 64)
	gm.Inventory.AddItem(items.STONE_BLOCK, 64)
	gm.Inventory.AddItem(items.GRASS_BLOCK, 64)
	gm.Inventory.AddItem(items.WORKBENCH, 1)
	gm.Inventory.AddItem(items.WOODEN_PICKAXE, 1)
	gm.Inventory.AddItem(items.COAL, 10)

	// Initialize plugin system
	gm.PluginManager = plugins.NewPluginManager()
	defaultPlugin := plugins.NewDefaultPlugin()
	gm.PluginManager.RegisterPlugin(defaultPlugin)
	gm.PluginManager.EnablePlugin("default")

	// Initialize skin editor
	gm.SkinEditor = skin.NewSkinEditor()

	// Initialize save system
	gm.SaveManager = save.NewSaveManager(worldName, "player")

	// Initialize time systems
	gm.DayNightCycle = gametime.NewDayNightCycle(600.0)
	gm.WeatherSystem = weather.NewWeatherSystem()

	// Initialize audio system
	gm.AudioManager = audio.NewAudioManager()
	gm.SoundLibrary = audio.NewSoundLibrary(gm.AudioManager)
	gm.BackgroundMusicManager = audio.NewBackgroundMusicManager(gm.AudioManager)

	// Initialize debug systems
	gm.RecoveryHandler = debug.NewRecoveryHandler("", func(info debug.PanicInfo) {
		log.Printf("Recovered from panic: %s", info.Message)
	})
	gm.Profiler = debug.NewPerformanceProfiler()

	// Initialize survival systems
	gameMode := survival.ModeSurvival
	if creativeMode {
		gameMode = survival.ModeCreative
	}
	gm.SurvivalManager = survival.NewSurvivalManager(gameMode, gm.Player, gm.Inventory)
	gm.EquipmentSet = equipment.NewEquipmentSet()
	gm.HealthSystem = health.NewLocationalHealthSystem()
	gm.BackpackUI = ui.NewBackpackUI(800, 600, gm.Inventory, gm.EquipmentSet, gm.HealthSystem)

	// Initialize chest system
	gm.ChestManager = chest.NewChestManager(worldName)

	// Initialize combat
	gm.WeaponSystem = combat.NewWeaponSystem()

	// Initialize UI effects
	gm.DamageIndicators = ui.NewDamageIndicatorManager(800, 600)
	gm.ScreenFlash = ui.NewScreenFlash()
	gm.DirectionalHitInd = ui.NewDirectionalHitManager()
	gm.DeathScreen = ui.NewDeathScreen(800, 600)

	// Initialize enemy systems
	gm.ZombieSpawner = enemies.NewZombieSpawner(gm.DayNightCycle)

	log.Printf("GameManager initialized successfully")
	return gm
}

// StartBackgroundMusic starts the background music loop
func (gm *GameManager) StartBackgroundMusic() error {
	log.Printf("Starting background music...")
	return gm.BackgroundMusicManager.PlayBackgroundMusic("gameplay_music")
}

// StopBackgroundMusic stops the background music
func (gm *GameManager) StopBackgroundMusic() {
	gm.BackgroundMusicManager.StopBackgroundMusic()
}

// Update updates all game systems
func (gm *GameManager) Update(deltaTime float64) error {
	defer gm.RecoveryHandler.Recover()

	gm.TotalPlayTime += time.Duration(deltaTime * float64(time.Second))

	state := gm.StateManager.GetState()

	switch state {
	case ui.StateCrafting:
		return gm.CraftingUI.Update()
	case ui.StateBackpack:
		return gm.BackpackUI.Update()
	case ui.StateChest:
		if gm.ChestUI != nil {
			return gm.ChestUI.Update()
		}
	case ui.StatePluginUI:
		if gm.PluginUI != nil {
			return gm.PluginUI.Update()
		}
	case ui.StateSkinEditor:
		return gm.SkinEditor.Update()
	case ui.StateDeathScreen:
		return gm.DeathScreen.Update()
	case ui.StateGame:
		// Update game systems
		gm.Player.Update(deltaTime)
		gm.DayNightCycle.Update()
		gm.SurvivalManager.Update(deltaTime)
		gm.HealthSystem.Update(deltaTime)
		if gm.DamageIndicators != nil {
			gm.DamageIndicators.Update(deltaTime)
		}
		if gm.ScreenFlash != nil {
			gm.ScreenFlash.Update()
		}
		if gm.DirectionalHitInd != nil {
			gm.DirectionalHitInd.Update()
		}
	}

	// Update audio system
	gm.AudioManager.Update()
	gm.BackgroundMusicManager.Update()

	return nil
}

// Cleanup cleans up all resources
func (gm *GameManager) Cleanup() {
	log.Printf("Cleaning up GameManager...")
	gm.StopBackgroundMusic()
	gm.AudioManager.Cleanup()
}
