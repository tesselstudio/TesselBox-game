package save

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"tesselbox/pkg/chest"
	"tesselbox/pkg/enemies"
	"tesselbox/pkg/equipment"
	"tesselbox/pkg/health"
	"tesselbox/pkg/items"
	"tesselbox/pkg/player"
	"tesselbox/pkg/survival"
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
	GameMode   string    `json:"game_mode"` // "creative" or "survival"

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
	HotbarSlots    []InventorySlotData `json:"hotbar_slots"`

	// World state
	CameraX          float64 `json:"camera_x"`
	CameraY          float64 `json:"camera_y"`
	WorldTime        float64 `json:"world_time"`        // Day/night cycle time
	Weather          string  `json:"weather"`           // Current weather
	CurrentDimension string  `json:"current_dimension"` // "overworld" or "randomland"

	// Game state
	InMenu       bool `json:"in_menu"`
	InGame       bool `json:"in_game"`
	InCrafting   bool `json:"in_crafting"`
	CreativeMode bool `json:"creative_mode"`

	// Crafting state
	CraftingStation string   `json:"crafting_station"`
	UnlockedRecipes []string `json:"unlocked_recipes"`

	// Statistics
	BlocksPlaced    int     `json:"blocks_placed"`
	BlocksDestroyed int     `json:"blocks_destroyed"`
	ItemsCrafted    int     `json:"items_crafted"`
	PlayTime        float64 `json:"play_time_seconds"`

	// Survival systems
	SurvivalStats  *SurvivalStatsData   `json:"survival_stats,omitempty"`
	Equipment      []EquipmentSlotData  `json:"equipment,omitempty"`
	BodyPartHealth []BodyPartHealthData `json:"body_part_health,omitempty"`
	Zombies        []ZombieData         `json:"zombies,omitempty"`
	Chests         []ChestData          `json:"chests,omitempty"`
}

// InventorySlotData represents a single inventory slot for serialization
type InventorySlotData struct {
	Type       items.ItemType `json:"type"`
	Quantity   int            `json:"quantity"`
	Durability int            `json:"durability"`
}

// SurvivalStatsData stores survival mode statistics
type SurvivalStatsData struct {
	Hunger         float64   `json:"hunger"`
	MaxHunger      float64   `json:"max_hunger"`
	Thirst         float64   `json:"thirst"`
	MaxThirst      float64   `json:"max_thirst"`
	Stamina        float64   `json:"stamina"`
	MaxStamina     float64   `json:"max_stamina"`
	LastDamageTime time.Time `json:"last_damage_time"`
	IsStarving     bool      `json:"is_starving"`
	IsDehydrated   bool      `json:"is_dehydrated"`
}

// EquipmentSlotData stores a single equipment item
type EquipmentSlotData struct {
	Name          string  `json:"name"`
	Slot          int     `json:"slot"`
	Material      int     `json:"material"`
	ArmorType     int     `json:"armor_type"`
	BaseDefense   float64 `json:"base_defense"`
	Durability    int     `json:"durability"`
	MaxDurability int     `json:"max_durability"`
	GrantsFlight  bool    `json:"grants_flight"`
}

// BodyPartHealthData stores health for each body part
type BodyPartHealthData struct {
	Name      string  `json:"name"`
	Health    float64 `json:"health"`
	MaxHealth float64 `json:"max_health"`
	IsVital   bool    `json:"is_vital"`
}

// ZombieData stores zombie state
type ZombieData struct {
	ID        string  `json:"id"`
	X         float64 `json:"x,omitempty"`
	Y         float64 `json:"y,omitempty"`
	Health    float64 `json:"health"`
	MaxHealth float64 `json:"max_health"`
	Type      int     `json:"type"`
	IsAlive   bool    `json:"is_alive"`
	State     int     `json:"state"`
}

// ChestData stores chest contents and position
type ChestData struct {
	X     float64             `json:"x"`
	Y     float64             `json:"y"`
	Slots []InventorySlotData `json:"slots"`
}

// SaveManager handles unified save game management
type SaveManager struct {
	SaveDir    string
	WorldName  string
	PlayerName string
}

// NewSaveManager creates a new save manager
func NewSaveManager(worldName, playerName string) *SaveManager {
	// Get user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home dir can't be found
		home = "."
	}
	saveDir := filepath.Join(home, ".tesselbox", "saves", worldName)

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
		Version:    "2.0", // Updated version for enhanced save format
		SaveTime:   time.Now(),
		WorldName:  sm.WorldName,
		PlayerName: sm.PlayerName,
		Seed:       gameState.World.Seed,
		GameMode:   gameState.GameMode,

		// Player state
		PlayerX:         gameState.Player.X,
		PlayerY:         gameState.Player.Y,
		PlayerVX:        gameState.Player.VX,
		PlayerVY:        gameState.Player.VY,
		PlayerHealth:    gameState.PlayerHealth,
		PlayerMaxHealth: gameState.PlayerMaxHealth,
		SelectedSlot:    gameState.Player.SelectedSlot,

		// Camera state
		CameraX: gameState.CameraX,
		CameraY: gameState.CameraY,

		// World state
		WorldTime:        gameState.WorldTime,
		Weather:          gameState.Weather,
		CurrentDimension: gameState.CurrentDimension,

		// Game state
		InMenu:       gameState.InMenu,
		InGame:       gameState.InGame,
		InCrafting:   gameState.InCrafting,
		CreativeMode: gameState.CreativeMode,

		// Crafting state
		CraftingStation: gameState.CraftingStation,
		UnlockedRecipes: gameState.UnlockedRecipes,

		// Statistics
		BlocksPlaced:    gameState.BlocksPlaced,
		BlocksDestroyed: gameState.BlocksDestroyed,
		ItemsCrafted:    gameState.ItemsCrafted,
		PlayTime:        gameState.PlayTime,
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

		// Save hotbar slots (first 10 slots of inventory)
		hotbarSize := 10
		if len(gameState.Inventory.Slots) < hotbarSize {
			hotbarSize = len(gameState.Inventory.Slots)
		}
		saveData.HotbarSlots = make([]InventorySlotData, hotbarSize)
		for i := 0; i < hotbarSize; i++ {
			slot := gameState.Inventory.Slots[i]
			saveData.HotbarSlots[i] = InventorySlotData{
				Type:       slot.Type,
				Quantity:   slot.Quantity,
				Durability: slot.Durability,
			}
		}
	}

	// Save survival systems
	if gameState.SurvivalManager != nil {
		saveData.SurvivalStats = &SurvivalStatsData{
			Hunger:         gameState.SurvivalManager.Hunger,
			MaxHunger:      gameState.SurvivalManager.MaxHunger,
			Thirst:         gameState.SurvivalManager.Thirst,
			MaxThirst:      gameState.SurvivalManager.MaxThirst,
			Stamina:        gameState.SurvivalManager.Stamina,
			MaxStamina:     gameState.SurvivalManager.MaxStamina,
			LastDamageTime: gameState.SurvivalManager.LastDamageTime,
			IsStarving:     gameState.SurvivalManager.IsStarving,
			IsDehydrated:   gameState.SurvivalManager.IsDehydrated,
		}
	}

	// Save equipment
	if gameState.EquipmentSet != nil {
		saveData.Equipment = make([]EquipmentSlotData, 0, len(gameState.EquipmentSet.Slots))
		for slotType, item := range gameState.EquipmentSet.Slots {
			if item != nil {
				saveData.Equipment = append(saveData.Equipment, EquipmentSlotData{
					Name:          item.Name,
					Slot:          int(slotType),
					Material:      int(item.Material),
					ArmorType:     int(item.ArmorType),
					BaseDefense:   item.BaseDefense,
					Durability:    item.Durability,
					MaxDurability: item.MaxDurability,
					GrantsFlight:  item.GrantsFlight,
				})
			}
		}
	}

	// Save body part health
	if gameState.HealthSystem != nil {
		saveData.BodyPartHealth = make([]BodyPartHealthData, len(gameState.HealthSystem.Parts))
		for i, part := range gameState.HealthSystem.Parts {
			saveData.BodyPartHealth[i] = BodyPartHealthData{
				Name:      part.Name,
				Health:    part.Health,
				MaxHealth: part.MaxHealth,
				IsVital:   part.IsVital,
			}
		}
	}

	// Save zombies
	if gameState.ZombieSpawner != nil {
		zombies := gameState.ZombieSpawner.Zombies
		saveData.Zombies = make([]ZombieData, 0, len(zombies))
		for _, zombie := range zombies {
			if zombie.IsAlive {
				saveData.Zombies = append(saveData.Zombies, ZombieData{
					ID:        zombie.ID,
					X:         zombie.X,
					Y:         zombie.Y,
					Health:    zombie.Health,
					MaxHealth: zombie.MaxHealth,
					Type:      int(zombie.Type),
					IsAlive:   zombie.IsAlive,
					State:     int(zombie.State),
				})
			}
		}
	}

	// Save chests
	if gameState.ChestManager != nil {
		chests := gameState.ChestManager.GetAllChests()
		saveData.Chests = make([]ChestData, 0, len(chests))
		for _, chest := range chests {
			slots := make([]InventorySlotData, len(chest.Slots))
			for i, slot := range chest.Slots {
				slots[i] = InventorySlotData{
					Type:       slot.Type,
					Quantity:   slot.Quantity,
					Durability: slot.Durability,
				}
			}
			saveData.Chests = append(saveData.Chests, ChestData{
				X:     chest.X,
				Y:     chest.Y,
				Slots: slots,
			})
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

	// Save player data with atomic write and backup
	if err := sm.saveGameAtomic(saveData); err != nil {
		return err
	}

	return nil
}

// saveGameAtomic performs an atomic save with backup
func (sm *SaveManager) saveGameAtomic(saveData *SaveData) error {
	saveFilename := filepath.Join(sm.SaveDir, fmt.Sprintf("player_%s.json", sm.PlayerName))
	tempFilename := saveFilename + ".tmp"

	// Marshal save data
	saveDataJSON, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %w", err)
	}

	// Create backup of existing save (if it exists)
	if _, err := os.Stat(saveFilename); err == nil {
		backupManager := NewBackupManager(sm.WorldName, sm.PlayerName, sm.SaveDir)
		if err := backupManager.CreateBackup(saveData); err != nil {
			// Log but don't fail - backup is best effort
			fmt.Printf("Warning: Failed to create backup: %v\n", err)
		}
	}

	// Write to temp file
	if err := os.WriteFile(tempFilename, saveDataJSON, 0644); err != nil {
		return fmt.Errorf("failed to write temp save file: %w", err)
	}

	// Validate the temp file
	backupManager := NewBackupManager(sm.WorldName, sm.PlayerName, sm.SaveDir)
	validator := NewSaveValidator(backupManager)
	result := validator.ValidateSaveFile(tempFilename)

	if !result.Valid {
		os.Remove(tempFilename)
		return fmt.Errorf("save validation failed: %v", result.Errors)
	}

	// Atomic rename (replaces existing file atomically)
	if err := os.Rename(tempFilename, saveFilename); err != nil {
		os.Remove(tempFilename)
		return fmt.Errorf("failed to commit save file: %w", err)
	}

	return nil
}

// LoadGame loads the complete game state with validation and recovery
func (sm *SaveManager) LoadGame() (*SaveData, error) {
	saveFilename := filepath.Join(sm.SaveDir, fmt.Sprintf("player_%s.json", sm.PlayerName))

	// First, validate the save file
	backupManager := NewBackupManager(sm.WorldName, sm.PlayerName, sm.SaveDir)
	validator := NewSaveValidator(backupManager)
	result := validator.ValidateSaveFile(saveFilename)

	if !result.Valid {
		// Save is corrupted - attempt recovery
		if result.CanRecover && result.BackupPath != "" {
			fmt.Printf("Save file is corrupted. Attempting recovery from backup: %s\n", result.BackupPath)
			recoveredData, err := validator.TryRecover()
			if err != nil {
				return nil, fmt.Errorf("save file corrupted and recovery failed: %v (errors: %v)", err, result.Errors)
			}
			fmt.Printf("Successfully recovered from backup. Warnings: %v\n", result.Warnings)
			return recoveredData, nil
		}
		return nil, fmt.Errorf("save file corrupted and no backup available: %v", result.Errors)
	}

	// Validate passed, load the data
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

	// Log any warnings
	if len(result.Warnings) > 0 {
		fmt.Printf("Save loaded with warnings: %v\n", result.Warnings)
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
	gameState.Player.SelectedSlot = saveData.SelectedSlot

	// Apply enhanced player state
	gameState.PlayerHealth = saveData.PlayerHealth
	gameState.PlayerMaxHealth = saveData.PlayerMaxHealth

	// Apply camera state
	gameState.CameraX = saveData.CameraX
	gameState.CameraY = saveData.CameraY

	// Apply world state
	gameState.WorldTime = saveData.WorldTime
	gameState.Weather = saveData.Weather
	gameState.CurrentDimension = saveData.CurrentDimension

	// Validate dimension consistency: if saved in randomland but randomland not generated,
	// force back to overworld to prevent void fall
	if gameState.CurrentDimension == "randomland" {
		// Check if dimension state file exists
		dimStatePath := filepath.Join(sm.SaveDir, "dimensions", "dimension_state.json")
		if _, err := os.Stat(dimStatePath); os.IsNotExist(err) {
			// Dimension state doesn't exist, reset to overworld
			fmt.Println("Warning: Saved in Randomland but no dimension data found. Resetting to overworld.")
			gameState.CurrentDimension = "overworld"
		}
	}

	// Apply game state
	gameState.InMenu = saveData.InMenu
	gameState.InGame = saveData.InGame
	gameState.InCrafting = saveData.InCrafting
	gameState.CreativeMode = saveData.CreativeMode
	gameState.GameMode = saveData.GameMode

	// Apply crafting state
	gameState.CraftingStation = saveData.CraftingStation
	gameState.UnlockedRecipes = saveData.UnlockedRecipes

	// Apply statistics
	gameState.BlocksPlaced = saveData.BlocksPlaced
	gameState.BlocksDestroyed = saveData.BlocksDestroyed
	gameState.ItemsCrafted = saveData.ItemsCrafted
	gameState.PlayTime = saveData.PlayTime

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

		// Restore hotbar slots if available
		if len(saveData.HotbarSlots) > 0 {
			for i, slotData := range saveData.HotbarSlots {
				if i < len(gameState.Inventory.Slots) {
					gameState.Inventory.Slots[i] = items.Item{
						Type:       slotData.Type,
						Quantity:   slotData.Quantity,
						Durability: slotData.Durability,
					}
				}
			}
		}
	}

	// Apply survival stats
	if saveData.SurvivalStats != nil && gameState.SurvivalManager != nil {
		gameState.SurvivalManager.Hunger = saveData.SurvivalStats.Hunger
		gameState.SurvivalManager.MaxHunger = saveData.SurvivalStats.MaxHunger
		gameState.SurvivalManager.Thirst = saveData.SurvivalStats.Thirst
		gameState.SurvivalManager.MaxThirst = saveData.SurvivalStats.MaxThirst
		gameState.SurvivalManager.Stamina = saveData.SurvivalStats.Stamina
		gameState.SurvivalManager.MaxStamina = saveData.SurvivalStats.MaxStamina
		gameState.SurvivalManager.LastDamageTime = saveData.SurvivalStats.LastDamageTime
		gameState.SurvivalManager.IsStarving = saveData.SurvivalStats.IsStarving
		gameState.SurvivalManager.IsDehydrated = saveData.SurvivalStats.IsDehydrated
	}

	// Apply equipment
	if len(saveData.Equipment) > 0 && gameState.EquipmentSet != nil {
		for _, eqData := range saveData.Equipment {
			slot := equipment.EquipmentSlot(eqData.Slot)
			item := &equipment.EquipmentItem{
				Name:          eqData.Name,
				Slot:          slot,
				Material:      equipment.ArmorMaterial(eqData.Material),
				ArmorType:     equipment.ArmorType(eqData.ArmorType),
				BaseDefense:   eqData.BaseDefense,
				Durability:    eqData.Durability,
				MaxDurability: eqData.MaxDurability,
				GrantsFlight:  eqData.GrantsFlight,
			}
			gameState.EquipmentSet.EquipItem(item, slot)
		}
	}

	// Apply body part health
	if len(saveData.BodyPartHealth) > 0 && gameState.HealthSystem != nil {
		totalHealth := 0.0
		totalMax := 0.0
		for i, partData := range saveData.BodyPartHealth {
			if i < len(gameState.HealthSystem.Parts) {
				gameState.HealthSystem.Parts[i].Health = partData.Health
				gameState.HealthSystem.Parts[i].MaxHealth = partData.MaxHealth
				gameState.HealthSystem.Parts[i].IsVital = partData.IsVital
				// Weight vital parts more heavily
				if partData.IsVital {
					totalHealth += partData.Health * 2
					totalMax += partData.MaxHealth * 2
				} else {
					totalHealth += partData.Health
					totalMax += partData.MaxHealth
				}
			}
		}
		// Recalculate overall health
		gameState.HealthSystem.OverallHealth = totalHealth
		gameState.HealthSystem.MaxOverallHealth = totalMax
	}

	// Load zombies
	if len(saveData.Zombies) > 0 && gameState.ZombieSpawner != nil {
		// Clear existing zombies and recreate from save
		gameState.ZombieSpawner.Zombies = make([]*enemies.Zombie, 0, len(saveData.Zombies))
		for _, zombieData := range saveData.Zombies {
			zombie := enemies.NewZombie(
				len(gameState.ZombieSpawner.Zombies),
				enemies.ZombieType(zombieData.Type),
				zombieData.X,
				zombieData.Y,
			)
			zombie.Health = zombieData.Health
			zombie.MaxHealth = zombieData.MaxHealth
			zombie.State = enemies.ZombieState(zombieData.State)
			gameState.ZombieSpawner.Zombies = append(gameState.ZombieSpawner.Zombies, zombie)
		}
		// Update NextID to avoid ID conflicts
		gameState.ZombieSpawner.NextID = len(gameState.ZombieSpawner.Zombies) + 1
	}

	// Load chests
	if len(saveData.Chests) > 0 && gameState.ChestManager != nil {
		for _, chestData := range saveData.Chests {
			slots := make([]items.Item, len(chestData.Slots))
			for i, slot := range chestData.Slots {
				slots[i] = items.Item{
					Type:       slot.Type,
					Quantity:   slot.Quantity,
					Durability: slot.Durability,
				}
			}
			gameState.ChestManager.SetChestContents(chestData.X, chestData.Y, slots)
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
		PlayerX:        saveData.PlayerX,
		PlayerY:        saveData.PlayerY,
		ChunkCount:     worldMetadata.ChunkCount,
		WorldCreatedAt: worldMetadata.CreatedAt,
		// Enhanced fields
		GameMode:        saveData.GameMode,
		PlayTime:        saveData.PlayTime,
		BlocksPlaced:    saveData.BlocksPlaced,
		BlocksDestroyed: saveData.BlocksDestroyed,
		ItemsCrafted:    saveData.ItemsCrafted,
		Seed:            saveData.Seed,
	}, nil
}

// SaveInfo contains metadata about a save file
type SaveInfo struct {
	PlayerName     string    `json:"player_name"`
	WorldName      string    `json:"world_name"`
	SaveTime       time.Time `json:"save_time"`
	Version        string    `json:"version"`
	PlayerX        float64   `json:"player_x"`
	PlayerY        float64   `json:"player_y"`
	ChunkCount     int       `json:"chunk_count"`
	WorldCreatedAt time.Time `json:"world_created_at"`
	// Enhanced fields
	GameMode        string  `json:"game_mode"`
	PlayTime        float64 `json:"play_time"`
	BlocksPlaced    int     `json:"blocks_placed"`
	BlocksDestroyed int     `json:"blocks_destroyed"`
	ItemsCrafted    int     `json:"items_crafted"`
	Seed            int64   `json:"seed"`
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

	// Enhanced game state
	CreativeMode     bool
	GameMode         string  // "creative" or "survival"
	WorldTime        float64 // Day/night cycle time
	Weather          string  // Current weather
	CurrentDimension string  // "overworld" or "randomland"
	PlayerHealth     float64
	PlayerMaxHealth  float64

	// Crafting state
	CraftingStation string
	UnlockedRecipes []string

	// Statistics
	BlocksPlaced    int
	BlocksDestroyed int
	ItemsCrafted    int
	PlayTime        float64

	// Survival systems
	SurvivalManager *survival.SurvivalManager
	EquipmentSet    *equipment.EquipmentSet
	HealthSystem    *health.LocationalHealthSystem

	// Enemy systems
	ZombieSpawner *enemies.ZombieSpawner

	// Storage systems
	ChestManager *chest.ChestManager
}

// AutoSaver handles automatic saving at intervals
type AutoSaver struct {
	saveManager *SaveManager
	gameState   *GameState
	interval    time.Duration
	lastSave    time.Time
	enabled     bool
	stopChan    chan bool
	mutex       sync.RWMutex // Fixed: Add mutex for thread safety
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
	as.mutex.Lock()
	defer as.mutex.Unlock()

	if as.enabled {
		return
	}
	as.enabled = true
	go as.autoSaveLoop()
}

// Stop stops the auto-saver
func (as *AutoSaver) Stop() {
	as.mutex.Lock()
	defer as.mutex.Unlock()

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
	as.mutex.RLock()
	defer as.mutex.RUnlock()

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
			as.mutex.RLock()
			enabled := as.enabled
			as.mutex.RUnlock()

			if enabled {
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

// BackupSave creates a backup of the current save
func (sm *SaveManager) BackupSave() error {
	saveData, err := sm.LoadGame()
	if err != nil {
		return err
	}

	// Create backup filename with timestamp
	timestamp := saveData.SaveTime.Format("20060102_150405")
	backupFilename := filepath.Join(sm.SaveDir, fmt.Sprintf("backup_%s_%s.json", sm.PlayerName, timestamp))

	saveDataJSON, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup data: %w", err)
	}

	return os.WriteFile(backupFilename, saveDataJSON, 0644)
}

// ListBackups returns a list of all backup files for this world
func (sm *SaveManager) ListBackups() ([]string, error) {
	entries, err := os.ReadDir(sm.SaveDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var backups []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "backup_") {
			backups = append(backups, entry.Name())
		}
	}

	return backups, nil
}
