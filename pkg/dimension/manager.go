// Package dimension provides dimension management for teleportation
// and world switching in TesselBox
package dimension

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"tesselbox/pkg/player"
	"tesselbox/pkg/world"
)

// Manager handles dimension switching and state management
type Manager struct {
	CurrentDimension   DimensionType
	OverworldWorld     *world.World
	RandomlandDim      *RandomlandDimension
	PlayerLastOverworldX float64
	PlayerLastOverworldY float64
	StoragePath        string
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
	}
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
func (m *Manager) TeleportToRandomland(player *player.Player) error {
	// Save current position in overworld
	m.PlayerLastOverworldX, m.PlayerLastOverworldY = player.GetCenter()

	// Initialize Randomland if needed
	if m.RandomlandDim == nil {
		m.RandomlandDim = NewRandomlandDimension()
		m.RandomlandDim.Generate()
	}

	// Switch dimension
	m.CurrentDimension = Randomland

	// Position player at return portal
	spawnX, spawnY := m.RandomlandDim.GetSpawnPosition()
	player.SetPosition(spawnX, spawnY)

	fmt.Printf("Teleported to Randomland! (Return to overworld at %.1f, %.1f)\n",
		m.PlayerLastOverworldX, m.PlayerLastOverworldY)

	return nil
}

// TeleportToOverworld teleports player back to overworld
func (m *Manager) TeleportToOverworld(player *player.Player) error {
	if m.CurrentDimension != Randomland {
		return fmt.Errorf("not currently in Randomland")
	}

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
	}
}

// DimensionState represents save data for dimensions
type DimensionState struct {
	RandomlandGenerated bool    `json:"randomland_generated"`
	ReturnPortalX       float64 `json:"return_portal_x"`
	ReturnPortalY       float64 `json:"return_portal_y"`
	LastOverworldX      float64 `json:"last_overworld_x"`
	LastOverworldY      float64 `json:"last_overworld_y"`
}

// Save saves dimension state
func (m *Manager) Save() error {
	// Ensure storage directory exists
	if err := os.MkdirAll(m.StoragePath, 0755); err != nil {
		return fmt.Errorf("failed to create dimension storage: %w", err)
	}

	state := DimensionState{
		RandomlandGenerated: m.RandomlandDim != nil && m.RandomlandDim.IsGenerated(),
		LastOverworldX:      m.PlayerLastOverworldX,
		LastOverworldY:      m.PlayerLastOverworldY,
	}

	if m.RandomlandDim != nil {
		state.ReturnPortalX = m.RandomlandDim.ReturnPortalX
		state.ReturnPortalY = m.RandomlandDim.ReturnPortalY
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

	// If Randomland was generated, recreate it (but don't regenerate terrain)
	if state.RandomlandGenerated {
		m.RandomlandDim = NewRandomlandDimension()
		m.RandomlandDim.ReturnPortalX = state.ReturnPortalX
		m.RandomlandDim.ReturnPortalY = state.ReturnPortalY
		// Mark as generated so we don't regenerate
		m.RandomlandDim.Generated = true
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

	// Check if it's a portal block (we'll add RANDOMLAND_PORTAL type)
	// For now, use OBSIDIAN as a placeholder
	return hex.BlockType == 98 // RANDOMLAND_PORTAL type number
}
