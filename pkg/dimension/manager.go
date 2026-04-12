// Package dimension provides dimension management for teleportation
// and world switching in TesselBox
package dimension

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"tesselbox/pkg/blocks"
	"tesselbox/pkg/enemies"
	"tesselbox/pkg/player"
	"tesselbox/pkg/world"
)

// Manager handles dimension switching and state management
type Manager struct {
	CurrentDimension     DimensionType
	OverworldWorld       *world.World
	RandomlandDim        *RandomlandDimension
	PlayerLastOverworldX float64
	PlayerLastOverworldY float64
	StoragePath          string

	// Teleport cooldown
	lastTeleportTime time.Time
	teleportCooldown time.Duration
}

// NewManager creates a new dimension manager
func NewManager(overworld *world.World, storageDir string) *Manager {
	return &Manager{
		CurrentDimension:     Overworld,
		OverworldWorld:       overworld,
		RandomlandDim:        nil, // Created on first use
		PlayerLastOverworldX: 0,
		PlayerLastOverworldY: 0,
		StoragePath:          filepath.Join(storageDir, "dimensions"),
		teleportCooldown:     2 * time.Second, // 2 second cooldown between teleports
	}
}

// CanTeleport checks if teleport is off cooldown
func (m *Manager) CanTeleport() bool {
	return time.Since(m.lastTeleportTime) >= m.teleportCooldown
}

// GetTeleportCooldownRemaining returns remaining cooldown time
func (m *Manager) GetTeleportCooldownRemaining() time.Duration {
	elapsed := time.Since(m.lastTeleportTime)
	if elapsed >= m.teleportCooldown {
		return 0
	}
	return m.teleportCooldown - elapsed
}

// GetCurrentWorld returns the currently active world
func (m *Manager) GetCurrentWorld() *world.World {
	switch m.CurrentDimension {
	case Randomland:
		if m.RandomlandDim != nil {
			return m.RandomlandDim.GetWorld()
		}
		return nil
	default:
		return m.OverworldWorld
	}
}

// GetCurrentDimensionName returns the name of current dimension
func (m *Manager) GetCurrentDimensionName() string {
	switch m.CurrentDimension {
	case Randomland:
		return "Randomland"
	default:
		return "Overworld"
	}
}

// IsInRandomland returns true if currently in Randomland
func (m *Manager) IsInRandomland() bool {
	return m.CurrentDimension == Randomland
}

// TeleportToRandomland teleports player to Randomland
// progressCallback is optional and receives generation progress updates
func (m *Manager) TeleportToRandomland(player *player.Player, progressCallback func(float64, string)) error {
	// Check teleport cooldown
	if !m.CanTeleport() {
		return fmt.Errorf("portal on cooldown (%.1f seconds remaining)", m.GetTeleportCooldownRemaining().Seconds())
	}

	// Save current position in overworld
	m.PlayerLastOverworldX, m.PlayerLastOverworldY = player.GetCenter()

	// Initialize Randomland if needed
	if m.RandomlandDim == nil {
		m.RandomlandDim = NewRandomlandDimension(m.OverworldWorld.WorldName)
		if err := m.RandomlandDim.Generate(progressCallback); err != nil {
			return fmt.Errorf("failed to generate Randomland: %w", err)
		}
	}

	// Record teleport time for cooldown
	m.lastTeleportTime = time.Now()

	// Switch dimension
	m.CurrentDimension = Randomland

	// Position player at return portal and reset velocity
	spawnX, spawnY := m.RandomlandDim.GetSpawnPosition()
	player.SetPosition(spawnX, spawnY)
	player.SetVelocity(0, 0) // Reset velocity to prevent launch into void

	fmt.Printf("Teleported to Randomland! (Return to overworld at %.1f, %.1f)\n",
		m.PlayerLastOverworldX, m.PlayerLastOverworldY)

	return nil
}

// TeleportToOverworld teleports player back to overworld
func (m *Manager) TeleportToOverworld(player *player.Player) error {
	if m.CurrentDimension != Randomland {
		return fmt.Errorf("not currently in Randomland")
	}

	// Check teleport cooldown
	if !m.CanTeleport() {
		return fmt.Errorf("portal on cooldown (%.1f seconds remaining)", m.GetTeleportCooldownRemaining().Seconds())
	}

	// Save and unload Randomland chunks to free memory
	if m.RandomlandDim != nil && m.RandomlandDim.World != nil {
		m.RandomlandDim.World.UnloadAllChunks()
	}

	// Record teleport time for cooldown
	m.lastTeleportTime = time.Now()

	// Switch dimension
	m.CurrentDimension = Overworld

	// Return player to saved position (or spawn if none saved)
	player.SetPosition(m.PlayerLastOverworldX, m.PlayerLastOverworldY)

	fmt.Println("Returned to overworld")

	return nil
}

// CheckReturnPortalProximity checks if player is near return portal in Randomland
func (m *Manager) CheckReturnPortalProximity(player *player.Player) bool {
	if m.CurrentDimension != Randomland || m.RandomlandDim == nil {
		return false
	}

	px, py := player.GetCenter()
	return m.RandomlandDim.IsNearReturnPortal(px, py, 60.0) // 60 pixel tolerance
}

// Update updates the current dimension
func (m *Manager) Update(player *player.Player, deltaTime float64) {
	if m.CurrentDimension == Randomland && m.RandomlandDim != nil {
		px, py := player.GetCenter()
		m.RandomlandDim.Update(px, py, deltaTime)

		// Safety check: ensure return portal exists (player could have destroyed it)
		// Check every 5 seconds to avoid performance impact
		if int(deltaTime*60)%300 == 0 {
			m.RandomlandDim.EnsureReturnPortal()
		}
	}
}

// ZombieData represents a zombie for serialization
type ZombieData struct {
	ID        string  `json:"id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Health    float64 `json:"health"`
	MaxHealth float64 `json:"max_health"`
	Type      int     `json:"type"`
	IsAlive   bool    `json:"is_alive"`
	State     int     `json:"state"`
}

// DimensionState represents save data for dimensions
type DimensionState struct {
	RandomlandGenerated bool         `json:"randomland_generated"`
	ReturnPortalX       float64      `json:"return_portal_x"`
	ReturnPortalY       float64      `json:"return_portal_y"`
	LastOverworldX      float64      `json:"last_overworld_x"`
	LastOverworldY      float64      `json:"last_overworld_y"`
	RandomlandZombies   []ZombieData `json:"randomland_zombies,omitempty"`
}

// Save saves dimension state
func (m *Manager) Save() error {
	// Ensure storage directory exists
	if err := os.MkdirAll(m.StoragePath, 0755); err != nil {
		return fmt.Errorf("failed to create dimension storage: %w", err)
	}

	// Save Randomland world if it exists
	if m.RandomlandDim != nil && m.RandomlandDim.IsGenerated() {
		if err := m.RandomlandDim.World.SaveWorld(); err != nil {
			return fmt.Errorf("failed to save randomland world: %w", err)
		}
	}

	state := DimensionState{
		RandomlandGenerated: m.RandomlandDim != nil && m.RandomlandDim.IsGenerated(),
		LastOverworldX:      m.PlayerLastOverworldX,
		LastOverworldY:      m.PlayerLastOverworldY,
	}

	if m.RandomlandDim != nil {
		state.ReturnPortalX = m.RandomlandDim.ReturnPortalX
		state.ReturnPortalY = m.RandomlandDim.ReturnPortalY

		// Save Randomland zombies
		if m.RandomlandDim.ZombieSpawner != nil {
			state.RandomlandZombies = make([]ZombieData, 0, len(m.RandomlandDim.ZombieSpawner.Zombies))
			for _, zombie := range m.RandomlandDim.ZombieSpawner.Zombies {
				if zombie.IsAlive {
					state.RandomlandZombies = append(state.RandomlandZombies, ZombieData{
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
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal dimension state: %w", err)
	}

	filename := filepath.Join(m.StoragePath, "dimension_state.json")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write dimension state: %w", err)
	}

	return nil
}

// Load loads dimension state
func (m *Manager) Load() error {
	filename := filepath.Join(m.StoragePath, "dimension_state.json")

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// No saved state yet, that's fine
			return nil
		}
		return fmt.Errorf("failed to read dimension state: %w", err)
	}

	var state DimensionState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to unmarshal dimension state: %w", err)
	}

	// Restore state
	m.PlayerLastOverworldX = state.LastOverworldX
	m.PlayerLastOverworldY = state.LastOverworldY

	// If Randomland was generated, recreate it and load the world
	if state.RandomlandGenerated {
		m.RandomlandDim = NewRandomlandDimension(m.OverworldWorld.WorldName)
		m.RandomlandDim.ReturnPortalX = state.ReturnPortalX
		m.RandomlandDim.ReturnPortalY = state.ReturnPortalY

		// Load the saved Randomland world
		if err := m.RandomlandDim.World.LoadWorldArea(250, 250, 10); err != nil {
			// If loading fails, we'll regenerate
			fmt.Printf("Warning: Failed to load Randomland world, will regenerate: %v\n", err)
			m.RandomlandDim.Generated = false
		} else {
			// Mark as generated after successful load
			m.RandomlandDim.Generated = true
			fmt.Println("Randomland world loaded successfully")
		}

		// Restore Randomland zombies
		if len(state.RandomlandZombies) > 0 && m.RandomlandDim.ZombieSpawner != nil {
			m.RandomlandDim.ZombieSpawner.Zombies = make([]*enemies.Zombie, 0, len(state.RandomlandZombies))
			for _, zombieData := range state.RandomlandZombies {
				zombie := enemies.NewZombie(
					len(m.RandomlandDim.ZombieSpawner.Zombies),
					enemies.ZombieType(zombieData.Type),
					zombieData.X,
					zombieData.Y,
				)
				zombie.ID = zombieData.ID
				zombie.Health = zombieData.Health
				zombie.MaxHealth = zombieData.MaxHealth
				zombie.IsAlive = zombieData.IsAlive
				zombie.State = enemies.ZombieState(zombieData.State)
				m.RandomlandDim.ZombieSpawner.Zombies = append(m.RandomlandDim.ZombieSpawner.Zombies, zombie)
			}
			// Update NextID to avoid ID conflicts
			m.RandomlandDim.ZombieSpawner.NextID = len(m.RandomlandDim.ZombieSpawner.Zombies) + 1
			fmt.Printf("Restored %d Randomland zombies\n", len(m.RandomlandDim.ZombieSpawner.Zombies))
		}
	}

	return nil
}

// CanTeleportToRandomland checks if player is on a portal block in overworld
func (m *Manager) CanTeleportToRandomland(player *player.Player) bool {
	if m.CurrentDimension != Overworld {
		return false
	}

	// Check if standing on a portal block
	px, py := player.GetCenter()
	world := m.OverworldWorld

	hex := world.GetHexagonAt(px, py)
	if hex == nil {
		return false
	}

	// Check if it's a portal block
	return hex.BlockType == blocks.RANDOMLAND_PORTAL
}
