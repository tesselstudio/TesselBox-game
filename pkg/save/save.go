package save

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"tesselbox/pkg/items"
	"tesselbox/pkg/player"
	"tesselbox/pkg/world"
)

// SaveData represents the complete game state that can be saved/loaded
type SaveData struct {
	// Metadata
	Version    string    `json:"version"`
	SaveTime   time.Time `json:"save_time"`
	WorldName  string    `json:"world_name"`
	PlayerName string    `json:"player_name"`
	Seed       int64     `json:"seed"`

	// Player state
	PlayerX         float64 `json:"player_x"`
	PlayerY         float64 `json:"player_y"`
	PlayerVX        float64 `json:"player_vx"`
	PlayerVY        float64 `json:"player_vy"`
	PlayerHealth    float64 `json:"player_health"`
	PlayerMaxHealth float64 `json:"player_max_health"`
	SelectedSlot    int     `json:"selected_slot"`

	// Inventory state
	InventorySlots []InventorySlotData `json:"inventory_slots"`

	// World state (chunks are saved separately by world storage)
	CameraX float64 `json:"camera_x"`
	CameraY float64 `json:"camera_y"`

	// Game state
	InMenu     bool `json:"in_menu"`
	InGame     bool `json:"in_game"`
	InCrafting bool `json:"in_crafting"`
}

// InventorySlotData represents a single inventory slot for serialization
type InventorySlotData struct {
	Type       items.ItemType `json:"type"`
	Quantity   int            `json:"quantity"`
	Durability int            `json:"durability"`
}

// SaveManager handles unified save game management
type SaveManager struct {
	SaveDir    string
	WorldName  string
	PlayerName string
}

// NewSaveManager creates a new save manager
func NewSaveManager(worldName, playerName string) *SaveManager {
	saveDir := filepath.Join("saves", worldName)

	// Create save directory if it doesn't exist
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		os.MkdirAll(saveDir, 0755)
	}

	return &SaveManager{
		SaveDir:    saveDir,
		WorldName:  worldName,
		PlayerName: playerName,
	}
}

// SaveGame saves the complete game state
func (sm *SaveManager) SaveGame(gameState *GameState) error {
	// Create save data
	saveData := &SaveData{
		Version:    "1.0",
		SaveTime:   time.Now(),
		WorldName:  sm.WorldName,
		PlayerName: sm.PlayerName,
		Seed:       gameState.World.Seed,

		// Player state
		PlayerX:         gameState.Player.X,
		PlayerY:         gameState.Player.Y,
		PlayerVX:        gameState.Player.VX,
		PlayerVY:        gameState.Player.VY,
		PlayerHealth:    gameState.Player.Health,
		PlayerMaxHealth: gameState.Player.MaxHealth,
		SelectedSlot:    gameState.Player.SelectedSlot,

		// Camera state
		CameraX: gameState.CameraX,
		CameraY: gameState.CameraY,

		// Game state
		InMenu:     gameState.InMenu,
		InGame:     gameState.InGame,
		InCrafting: gameState.InCrafting,
	}

	// Convert inventory to serializable format
	if gameState.Inventory != nil {
		saveData.InventorySlots = make([]InventorySlotData, len(gameState.Inventory.Slots))
		for i, slot := range gameState.Inventory.Slots {
			saveData.InventorySlots[i] = InventorySlotData{
				Type:       slot.Type,
				Quantity:   slot.Quantity,
				Durability: slot.Durability,
			}
		}
	}

	// Save world chunks first
	worldStorage := world.NewWorldStorage(sm.WorldName)
	if err := worldStorage.SaveWorld(gameState.World); err != nil {
		return fmt.Errorf("failed to save world: %w", err)
	}

	// Save world metadata
	chunkCount := len(gameState.World.Chunks)
	if err := worldStorage.SaveWorldMetadata(chunkCount); err != nil {
		return fmt.Errorf("failed to save world metadata: %w", err)
	}

	// Save player data
	saveFilename := filepath.Join(sm.SaveDir, fmt.Sprintf("player_%s.json", sm.PlayerName))
	saveDataJSON, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %w", err)
	}

	if err := os.WriteFile(saveFilename, saveDataJSON, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

// LoadGame loads the complete game state
func (sm *SaveManager) LoadGame() (*SaveData, error) {
	saveFilename := filepath.Join(sm.SaveDir, fmt.Sprintf("player_%s.json", sm.PlayerName))

	data, err := os.ReadFile(saveFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("save file not found for player %s in world %s", sm.PlayerName, sm.WorldName)
		}
		return nil, fmt.Errorf("failed to read save file: %w", err)
	}

	var saveData SaveData
	err = json.Unmarshal(data, &saveData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal save data: %w", err)
	}

	return &saveData, nil
}

// ApplySaveData applies loaded save data to the game state
func (sm *SaveManager) ApplySaveData(saveData *SaveData, gameState *GameState) error {
	// Apply player state
	gameState.Player.X = saveData.PlayerX
	gameState.Player.Y = saveData.PlayerY
	gameState.Player.VX = saveData.PlayerVX
	gameState.Player.VY = saveData.PlayerVY
	gameState.Player.Health = saveData.PlayerHealth
	gameState.Player.MaxHealth = saveData.PlayerMaxHealth
	gameState.Player.SelectedSlot = saveData.SelectedSlot

	// Apply camera state
	gameState.CameraX = saveData.CameraX
	gameState.CameraY = saveData.CameraY

	// Apply game state
	gameState.InMenu = saveData.InMenu
	gameState.InGame = saveData.InGame
	gameState.InCrafting = saveData.InCrafting

	// Apply inventory state
	if gameState.Inventory != nil && len(saveData.InventorySlots) > 0 {
		// Ensure inventory has enough slots
		if len(saveData.InventorySlots) > len(gameState.Inventory.Slots) {
			// Expand inventory if needed
			newSlots := make([]items.Item, len(saveData.InventorySlots))
			copy(newSlots, gameState.Inventory.Slots)
			for i := len(gameState.Inventory.Slots); i < len(newSlots); i++ {
				newSlots[i] = items.Item{Type: items.NONE, Quantity: 0, Durability: -1}
			}
			gameState.Inventory.Slots = newSlots
		}

		// Restore inventory slots
		for i, slotData := range saveData.InventorySlots {
			if i < len(gameState.Inventory.Slots) {
				gameState.Inventory.Slots[i] = items.Item{
					Type:       slotData.Type,
					Quantity:   slotData.Quantity,
					Durability: slotData.Durability,
				}
			}
		}
	}

	// Load world chunks around player position
	worldStorage := world.NewWorldStorage(sm.WorldName)
	if err := worldStorage.LoadWorld(gameState.World, saveData.PlayerX, saveData.PlayerY, 5); err != nil {
		return fmt.Errorf("failed to load world: %w", err)
	}

	return nil
}

// DeleteSave deletes a save file
func (sm *SaveManager) DeleteSave() error {
	saveFilename := filepath.Join(sm.SaveDir, fmt.Sprintf("player_%s.json", sm.PlayerName))

	if _, err := os.Stat(saveFilename); os.IsNotExist(err) {
		return nil // Save doesn't exist, nothing to delete
	}

	return os.Remove(saveFilename)
}

// ListSaves returns a list of all save files for this world
func (sm *SaveManager) ListSaves() ([]string, error) {
	entries, err := os.ReadDir(sm.SaveDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read save directory: %w", err)
	}

	var saves []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			saves = append(saves, entry.Name())
		}
	}

	return saves, nil
}

// GetSaveInfo returns metadata about a save file
func (sm *SaveManager) GetSaveInfo() (*SaveInfo, error) {
	saveData, err := sm.LoadGame()
	if err != nil {
		return nil, err
	}

	worldStorage := world.NewWorldStorage(sm.WorldName)
	worldMetadata, err := worldStorage.GetWorldMetadata()
	if err != nil {
		return nil, err
	}

	return &SaveInfo{
		PlayerName:     saveData.PlayerName,
		WorldName:      saveData.WorldName,
		SaveTime:       saveData.SaveTime,
		Version:        saveData.Version,
		PlayerHealth:   saveData.PlayerHealth,
		PlayerX:        saveData.PlayerX,
		PlayerY:        saveData.PlayerY,
		ChunkCount:     worldMetadata.ChunkCount,
		WorldCreatedAt: worldMetadata.CreatedAt,
	}, nil
}

// SaveInfo contains metadata about a save file
type SaveInfo struct {
	PlayerName     string    `json:"player_name"`
	WorldName      string    `json:"world_name"`
	SaveTime       time.Time `json:"save_time"`
	Version        string    `json:"version"`
	PlayerHealth   float64   `json:"player_health"`
	PlayerX        float64   `json:"player_x"`
	PlayerY        float64   `json:"player_y"`
	ChunkCount     int       `json:"chunk_count"`
	WorldCreatedAt time.Time `json:"world_created_at"`
}

// GameState represents the current game state (used for save/load operations)
type GameState struct {
	World      *world.World
	Player     *player.Player
	Inventory  *items.Inventory
	CameraX    float64
	CameraY    float64
	InMenu     bool
	InGame     bool
	InCrafting bool
}

// AutoSaver handles automatic saving at intervals
type AutoSaver struct {
	saveManager *SaveManager
	gameState   *GameState
	interval    time.Duration
	lastSave    time.Time
	enabled     bool
	stopChan    chan bool
}

// NewAutoSaver creates a new auto-saver
func NewAutoSaver(saveManager *SaveManager, gameState *GameState, interval time.Duration) *AutoSaver {
	return &AutoSaver{
		saveManager: saveManager,
		gameState:   gameState,
		interval:    interval,
		enabled:     false,
		stopChan:    make(chan bool),
	}
}

// Start starts the auto-saver
func (as *AutoSaver) Start() {
	if as.enabled {
		return
	}
	as.enabled = true
	go as.autoSaveLoop()
}

// Stop stops the auto-saver
func (as *AutoSaver) Stop() {
	if !as.enabled {
		return
	}
	as.enabled = false
	as.stopChan <- true
}

// SetInterval changes the auto-save interval
func (as *AutoSaver) SetInterval(interval time.Duration) {
	as.interval = interval
}

// ForceSave forces an immediate save
func (as *AutoSaver) ForceSave() error {
	as.lastSave = time.Now()
	return as.saveManager.SaveGame(as.gameState)
}

// autoSaveLoop runs the auto-save loop
func (as *AutoSaver) autoSaveLoop() {
	ticker := time.NewTicker(as.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if as.enabled {
				if err := as.ForceSave(); err != nil {
					// Log error but continue running
					fmt.Printf("Auto-save failed: %v\n", err)
				}
			}
		case <-as.stopChan:
			return
		}
	}
}

// ListAllWorlds returns a list of all saved worlds
func ListAllWorlds() ([]string, error) {
	return world.ListSavedWorlds()
}

// DeleteWorld deletes an entire world and all its saves
func DeleteWorld(worldName string) error {
	return world.DeleteWorld(worldName)
}
