package warps

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// WarpType represents the type of warp
type WarpType int

const (
	WarpPublic WarpType = iota
	WarpPrivate
	WarpShop
	WarpEvent
	WarpSpawn
)

// String returns warp type name
func (w WarpType) String() string {
	switch w {
	case WarpPublic:
		return "public"
	case WarpPrivate:
		return "private"
	case WarpShop:
		return "shop"
	case WarpEvent:
		return "event"
	case WarpSpawn:
		return "spawn"
	}
	return "unknown"
}

// Warp represents a teleport location
type Warp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Location
	WorldID string  `json:"world_id"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`

	// Type & Access
	Type           WarpType `json:"type"`
	OwnerID        string   `json:"owner_id"`
	IsPublic       bool     `json:"is_public"`
	AllowedPlayers []string `json:"allowed_players,omitempty"`
	RequiredRank   string   `json:"required_rank,omitempty"`

	// Cost
	Cost     float64       `json:"cost"`
	Cooldown time.Duration `json:"cooldown"`

	// Stats
	UseCount  int       `json:"use_count"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
}

// NewWarp creates a new warp
func NewWarp(id, name, worldID, ownerID string, warpType WarpType, x, y float64) *Warp {
	return &Warp{
		ID:             id,
		Name:           name,
		WorldID:        worldID,
		OwnerID:        ownerID,
		Type:           warpType,
		X:              x,
		Y:              y,
		IsPublic:       warpType == WarpPublic || warpType == WarpSpawn,
		AllowedPlayers: make([]string, 0),
		Cost:           0,
		Cooldown:       0,
		UseCount:       0,
		CreatedAt:      time.Now(),
	}
}

// SetCost sets the warp cost
func (w *Warp) SetCost(cost float64) {
	w.Cost = cost
}

// SetCooldown sets the cooldown
func (w *Warp) SetCooldown(cooldown time.Duration) {
	w.Cooldown = cooldown
}

// AllowPlayer allows a player to use this warp
func (w *Warp) AllowPlayer(playerID string) {
	for _, id := range w.AllowedPlayers {
		if id == playerID {
			return // Already allowed
		}
	}
	w.AllowedPlayers = append(w.AllowedPlayers, playerID)
}

// DisallowPlayer removes player access
func (w *Warp) DisallowPlayer(playerID string) {
	for i, id := range w.AllowedPlayers {
		if id == playerID {
			w.AllowedPlayers = append(w.AllowedPlayers[:i], w.AllowedPlayers[i+1:]...)
			return
		}
	}
}

// CanUse checks if a player can use this warp
func (w *Warp) CanUse(playerID string, playerRank string) bool {
	// Owner can always use
	if w.OwnerID == playerID {
		return true
	}

	// Public warps are accessible to all
	if w.IsPublic {
		return true
	}

	// Check allowed list
	for _, id := range w.AllowedPlayers {
		if id == playerID {
			return true
		}
	}

	// Check rank requirement
	if w.RequiredRank != "" && playerRank == w.RequiredRank {
		return true
	}

	return false
}

// RecordUse records a use
func (w *Warp) RecordUse() {
	w.UseCount++
	w.LastUsed = time.Now()
}

// Home represents a player's home location
type Home struct {
	ID       string  `json:"id"`
	PlayerID string  `json:"player_id"`
	Name     string  `json:"name"`
	WorldID  string  `json:"world_id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`

	// Settings
	IsDefault bool `json:"is_default"`
	IsPublic  bool `json:"is_public"`

	// Stats
	CreatedAt time.Time `json:"created_at"`
	UseCount  int       `json:"use_count"`
}

// WarpManager manages warps and homes
type WarpManager struct {
	warps   map[string]*Warp
	homes   map[string]*Home // key: "playerID_homeName"
	byWorld map[string][]string

	// Default costs
	defaultCost     float64
	defaultCooldown time.Duration

	storagePath string
}

// NewWarpManager creates a new warp manager
func NewWarpManager(storageDir string) *WarpManager {
	return &WarpManager{
		warps:           make(map[string]*Warp),
		homes:           make(map[string]*Home),
		byWorld:         make(map[string][]string),
		defaultCost:     5.0,
		defaultCooldown: 5 * time.Second,
		storagePath:     filepath.Join(storageDir, "warps.json"),
	}
}

// CreateWarp creates a new warp
func (wm *WarpManager) CreateWarp(id, name, worldID, ownerID string, warpType WarpType, x, y float64) (*Warp, error) {
	if _, exists := wm.warps[id]; exists {
		return nil, fmt.Errorf("warp with ID '%s' already exists", id)
	}

	warp := NewWarp(id, name, worldID, ownerID, warpType, x, y)
	warp.Cost = wm.defaultCost
	warp.Cooldown = wm.defaultCooldown

	wm.warps[id] = warp
	wm.byWorld[worldID] = append(wm.byWorld[worldID], id)

	return warp, nil
}

// GetWarp gets a warp by ID
func (wm *WarpManager) GetWarp(warpID string) (*Warp, bool) {
	warp, exists := wm.warps[warpID]
	return warp, exists
}

// GetWarpsByWorld gets warps in a world
func (wm *WarpManager) GetWarpsByWorld(worldID string) []*Warp {
	warpIDs := wm.byWorld[worldID]
	warps := make([]*Warp, 0, len(warpIDs))

	for _, id := range warpIDs {
		if warp, exists := wm.warps[id]; exists {
			warps = append(warps, warp)
		}
	}

	return warps
}

// GetPublicWarps gets all public warps
func (wm *WarpManager) GetPublicWarps() []*Warp {
	result := make([]*Warp, 0)
	for _, warp := range wm.warps {
		if warp.IsPublic {
			result = append(result, warp)
		}
	}
	return result
}

// GetAvailableWarps gets warps a player can use
func (wm *WarpManager) GetAvailableWarps(playerID string, playerRank string, worldID string) []*Warp {
	result := make([]*Warp, 0)

	for _, warp := range wm.warps {
		if warp.WorldID == worldID && warp.CanUse(playerID, playerRank) {
			result = append(result, warp)
		}
	}

	return result
}

// DeleteWarp deletes a warp
func (wm *WarpManager) DeleteWarp(warpID string) error {
	warp, exists := wm.warps[warpID]
	if !exists {
		return fmt.Errorf("warp not found")
	}

	// Remove from byWorld
	worldWarps := wm.byWorld[warp.WorldID]
	for i, id := range worldWarps {
		if id == warpID {
			wm.byWorld[warp.WorldID] = append(worldWarps[:i], worldWarps[i+1:]...)
			break
		}
	}

	delete(wm.warps, warpID)
	return nil
}

// SetHome sets a home for a player
func (wm *WarpManager) SetHome(playerID, name, worldID string, x, y float64) (*Home, error) {
	key := generateHomeKey(playerID, name)

	// Check if updating existing
	if home, exists := wm.homes[key]; exists {
		home.WorldID = worldID
		home.X = x
		home.Y = y
		return home, nil
	}

	home := &Home{
		ID:        fmt.Sprintf("home_%s_%d", playerID, time.Now().Unix()),
		PlayerID:  playerID,
		Name:      name,
		WorldID:   worldID,
		X:         x,
		Y:         y,
		CreatedAt: time.Now(),
	}

	// If first home, make it default
	playerHomes := wm.GetPlayerHomes(playerID)
	if len(playerHomes) == 0 {
		home.IsDefault = true
	}

	wm.homes[key] = home
	return home, nil
}

// GetHome gets a player's home
func (wm *WarpManager) GetHome(playerID, name string) (*Home, bool) {
	key := generateHomeKey(playerID, name)
	home, exists := wm.homes[key]
	return home, exists
}

// GetDefaultHome gets a player's default home
func (wm *WarpManager) GetDefaultHome(playerID string) (*Home, bool) {
	homes := wm.GetPlayerHomes(playerID)
	for _, home := range homes {
		if home.IsDefault {
			return home, true
		}
	}
	// Return first home if no default
	if len(homes) > 0 {
		return homes[0], true
	}
	return nil, false
}

// GetPlayerHomes gets all homes for a player
func (wm *WarpManager) GetPlayerHomes(playerID string) []*Home {
	result := make([]*Home, 0)
	prefix := playerID + "_"

	for key, home := range wm.homes {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			result = append(result, home)
		}
	}

	return result
}

// DeleteHome deletes a home
func (wm *WarpManager) DeleteHome(playerID, name string) error {
	key := generateHomeKey(playerID, name)
	if _, exists := wm.homes[key]; !exists {
		return fmt.Errorf("home not found")
	}

	delete(wm.homes, key)
	return nil
}

// SetDefaultHome sets the default home
func (wm *WarpManager) SetDefaultHome(playerID, name string) error {
	homes := wm.GetPlayerHomes(playerID)

	for _, home := range homes {
		home.IsDefault = (home.Name == name)
	}

	return nil
}

// GetHomeCount returns number of homes for a player
func (wm *WarpManager) GetHomeCount(playerID string) int {
	return len(wm.GetPlayerHomes(playerID))
}

// CanSetHome checks if player can set another home
func (wm *WarpManager) CanSetHome(playerID string, maxHomes int) bool {
	return wm.GetHomeCount(playerID) < maxHomes
}

// generateHomeKey generates a key for home lookup
func generateHomeKey(playerID, name string) string {
	return fmt.Sprintf("%s_%s", playerID, name)
}

// Save saves warps and homes
func (wm *WarpManager) Save() error {
	data := struct {
		Warps map[string]*Warp `json:"warps"`
		Homes map[string]*Home `json:"homes"`
	}{
		Warps: wm.warps,
		Homes: wm.homes,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	if err := os.WriteFile(wm.storagePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}

// Load loads warps and homes
func (wm *WarpManager) Load() error {
	data, err := os.ReadFile(wm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read: %w", err)
	}

	var loaded struct {
		Warps map[string]*Warp `json:"warps"`
		Homes map[string]*Home `json:"homes"`
	}

	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	wm.warps = loaded.Warps
	if wm.warps == nil {
		wm.warps = make(map[string]*Warp)
	}

	wm.homes = loaded.Homes
	if wm.homes == nil {
		wm.homes = make(map[string]*Home)
	}

	// Rebuild byWorld index
	wm.byWorld = make(map[string][]string)
	for _, warp := range wm.warps {
		wm.byWorld[warp.WorldID] = append(wm.byWorld[warp.WorldID], warp.ID)
	}

	return nil
}

// GetTopWarps returns most used warps
func (wm *WarpManager) GetTopWarps(count int) []*Warp {
	warps := make([]*Warp, 0, len(wm.warps))
	for _, warp := range wm.warps {
		warps = append(warps, warp)
	}

	// Sort by use count
	for i := 0; i < len(warps); i++ {
		for j := i + 1; j < len(warps); j++ {
			if warps[i].UseCount < warps[j].UseCount {
				warps[i], warps[j] = warps[j], warps[i]
			}
		}
	}

	if count > len(warps) {
		count = len(warps)
	}

	return warps[:count]
}
