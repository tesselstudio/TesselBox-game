package permissions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// PlayerEntry represents a registered player with their permissions
type PlayerEntry struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	RoleID          string                 `json:"role_id"`
	WorldID         string                 `json:"world_id,omitempty"`
	CustomPerms     map[string]bool        `json:"custom_perms,omitempty"` // Permission overrides
	
	// Economy data
	Balance         float64                `json:"balance"`
	BankBalance     float64                `json:"bank_balance"`
	
	// Stats
	PlayTime        time.Duration          `json:"play_time"`
	SessionCount    int                    `json:"session_count"`
	FirstJoin       time.Time              `json:"first_join"`
	LastSeen        time.Time              `json:"last_seen"`
	
	// Moderation
	Muted           bool                   `json:"muted"`
	MutedUntil      *time.Time             `json:"muted_until,omitempty"`
	Banned          bool                   `json:"banned"`
	BanReason       string                 `json:"ban_reason,omitempty"`
	BannedUntil     *time.Time             `json:"banned_until,omitempty"`
	WarnCount       int                    `json:"warn_count"`
	
	// Social
	Friends         []string               `json:"friends,omitempty"`
	Ignored         []string               `json:"ignored,omitempty"`
	
	// Security
	IPHistory       []string               `json:"ip_history,omitempty"`
	LastIP          string                 `json:"last_ip,omitempty"`
	
	// Metadata
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// NewPlayerEntry creates a new player entry with default role
func NewPlayerEntry(id, name, defaultRole string) *PlayerEntry {
	now := time.Now()
	return &PlayerEntry{
		ID:          id,
		Name:        name,
		RoleID:      defaultRole,
		CustomPerms: make(map[string]bool),
		Balance:     0,
		FirstJoin:   now,
		LastSeen:    now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// SetRole changes the player's role
func (pe *PlayerEntry) SetRole(roleID string) {
	pe.RoleID = roleID
	pe.UpdatedAt = time.Now()
}

// GrantCustomPermission grants a custom permission override
func (pe *PlayerEntry) GrantCustomPermission(node PermissionNode) {
	pe.CustomPerms[string(node)] = true
	pe.UpdatedAt = time.Now()
}

// DenyCustomPermission explicitly denies a permission
func (pe *PlayerEntry) DenyCustomPermission(node PermissionNode) {
	pe.CustomPerms[string(node)] = false
	pe.UpdatedAt = time.Now()
}

// RemoveCustomPermission removes a custom permission
func (pe *PlayerEntry) RemoveCustomPermission(node PermissionNode) {
	delete(pe.CustomPerms, string(node))
	pe.UpdatedAt = time.Now()
}

// AddPlayTime adds to total play time
func (pe *PlayerEntry) AddPlayTime(duration time.Duration) {
	pe.PlayTime += duration
	pe.UpdatedAt = time.Now()
}

// RecordSession records a login session
func (pe *PlayerEntry) RecordSession() {
	pe.SessionCount++
	pe.LastSeen = time.Now()
	pe.UpdatedAt = time.Now()
}

// RecordIP records a login IP
func (pe *PlayerEntry) RecordIP(ip string) {
	if pe.LastIP != ip {
		pe.IPHistory = append(pe.IPHistory, ip)
		if len(pe.IPHistory) > 10 { // Keep last 10 IPs
			pe.IPHistory = pe.IPHistory[len(pe.IPHistory)-10:]
		}
	}
	pe.LastIP = ip
	pe.UpdatedAt = time.Now()
}

// AddFriend adds a friend
func (pe *PlayerEntry) AddFriend(playerID string) bool {
	for _, id := range pe.Friends {
		if id == playerID {
			return false // Already friends
		}
	}
	pe.Friends = append(pe.Friends, playerID)
	pe.UpdatedAt = time.Now()
	return true
}

// RemoveFriend removes a friend
func (pe *PlayerEntry) RemoveFriend(playerID string) bool {
	for i, id := range pe.Friends {
		if id == playerID {
			pe.Friends = append(pe.Friends[:i], pe.Friends[i+1:]...)
			pe.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// IsFriend checks if player is a friend
func (pe *PlayerEntry) IsFriend(playerID string) bool {
	for _, id := range pe.Friends {
		if id == playerID {
			return true
		}
	}
	return false
}

// Ignore ignores a player
func (pe *PlayerEntry) Ignore(playerID string) {
	for _, id := range pe.Ignored {
		if id == playerID {
			return // Already ignored
		}
	}
	pe.Ignored = append(pe.Ignored, playerID)
	pe.UpdatedAt = time.Now()
}

// Unignore removes ignore
func (pe *PlayerEntry) Unignore(playerID string) {
	for i, id := range pe.Ignored {
		if id == playerID {
			pe.Ignored = append(pe.Ignored[:i], pe.Ignored[i+1:]...)
			pe.UpdatedAt = time.Now()
			return
		}
	}
}

// IsIgnored checks if player is ignored
func (pe *PlayerEntry) IsIgnored(playerID string) bool {
	for _, id := range pe.Ignored {
		if id == playerID {
			return true
		}
	}
	return false
}

// Mute mutes the player
func (pe *PlayerEntry) Mute(duration time.Duration) {
	pe.Muted = true
	until := time.Now().Add(duration)
	pe.MutedUntil = &until
	pe.UpdatedAt = time.Now()
}

// Unmute unmutes the player
func (pe *PlayerEntry) Unmute() {
	pe.Muted = false
	pe.MutedUntil = nil
	pe.UpdatedAt = time.Now()
}

// IsMuted checks if player is currently muted
func (pe *PlayerEntry) IsMuted() bool {
	if !pe.Muted {
		return false
	}
	if pe.MutedUntil != nil && time.Now().After(*pe.MutedUntil) {
		pe.Muted = false
		pe.MutedUntil = nil
		return false
	}
	return true
}

// Ban bans the player
func (pe *PlayerEntry) Ban(reason string, duration *time.Duration) {
	pe.Banned = true
	pe.BanReason = reason
	if duration != nil {
		until := time.Now().Add(*duration)
		pe.BannedUntil = &until
	}
	pe.UpdatedAt = time.Now()
}

// Unban unbans the player
func (pe *PlayerEntry) Unban() {
	pe.Banned = false
	pe.BanReason = ""
	pe.BannedUntil = nil
	pe.UpdatedAt = time.Now()
}

// IsBanned checks if player is currently banned
func (pe *PlayerEntry) IsBanned() bool {
	if !pe.Banned {
		return false
	}
	if pe.BannedUntil != nil && time.Now().After(*pe.BannedUntil) {
		pe.Banned = false
		pe.BanReason = ""
		pe.BannedUntil = nil
		return false
	}
	return true
}

// AddWarn increments warning count
func (pe *PlayerEntry) AddWarn() {
	pe.WarnCount++
	pe.UpdatedAt = time.Now()
}

// ClearWarns clears warning count
func (pe *PlayerEntry) ClearWarns() {
	pe.WarnCount = 0
	pe.UpdatedAt = time.Now()
}

// AddBalance adds money to balance
func (pe *PlayerEntry) AddBalance(amount float64) {
	pe.Balance += amount
	pe.UpdatedAt = time.Now()
}

// RemoveBalance removes money from balance (returns false if insufficient)
func (pe *PlayerEntry) RemoveBalance(amount float64) bool {
	if pe.Balance < amount {
		return false
	}
	pe.Balance -= amount
	pe.UpdatedAt = time.Now()
	return true
}

// PlayerRegistry manages all player entries
type PlayerRegistry struct {
	players map[string]*PlayerEntry  // By ID
	byName  map[string]string          // Name -> ID lookup
	mu      sync.RWMutex
	
	storagePath string
}

// NewPlayerRegistry creates a new player registry
func NewPlayerRegistry(storageDir string) *PlayerRegistry {
	return &PlayerRegistry{
		players:     make(map[string]*PlayerEntry),
		byName:      make(map[string]string),
		storagePath: filepath.Join(storageDir, "players.json"),
	}
}

// Register registers a new player
func (pr *PlayerRegistry) Register(id, name, defaultRole string) (*PlayerEntry, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	
	// Check if ID already exists
	if _, exists := pr.players[id]; exists {
		return nil, fmt.Errorf("player with ID '%s' already exists", id)
	}
	
	// Check if name is taken
	if existingID, exists := pr.byName[name]; exists {
		if existingID != id {
			return nil, fmt.Errorf("player name '%s' is already taken", name)
		}
	}
	
	entry := NewPlayerEntry(id, name, defaultRole)
	pr.players[id] = entry
	pr.byName[name] = id
	
	return entry, nil
}

// GetByID gets a player by ID
func (pr *PlayerRegistry) GetByID(id string) (*PlayerEntry, bool) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	
	entry, exists := pr.players[id]
	return entry, exists
}

// GetByName gets a player by name
func (pr *PlayerRegistry) GetByName(name string) (*PlayerEntry, bool) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	
	if id, exists := pr.byName[name]; exists {
		return pr.players[id], true
	}
	return nil, false
}

// Update updates a player entry
func (pr *PlayerRegistry) Update(entry *PlayerEntry) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	
	if _, exists := pr.players[entry.ID]; !exists {
		return fmt.Errorf("player with ID '%s' not found", entry.ID)
	}
	
	// Update name lookup if name changed
	oldEntry := pr.players[entry.ID]
	if oldEntry.Name != entry.Name {
		delete(pr.byName, oldEntry.Name)
		pr.byName[entry.Name] = entry.ID
	}
	
	entry.UpdatedAt = time.Now()
	pr.players[entry.ID] = entry
	
	return nil
}

// Delete removes a player
func (pr *PlayerRegistry) Delete(id string) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	
	entry, exists := pr.players[id]
	if !exists {
		return fmt.Errorf("player with ID '%s' not found", id)
	}
	
	delete(pr.byName, entry.Name)
	delete(pr.players, id)
	
	return nil
}

// GetAll returns all players
func (pr *PlayerRegistry) GetAll() []*PlayerEntry {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	
	result := make([]*PlayerEntry, 0, len(pr.players))
	for _, entry := range pr.players {
		result = append(result, entry)
	}
	
	return result
}

// GetByRole returns all players with a specific role
func (pr *PlayerRegistry) GetByRole(roleID string) []*PlayerEntry {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	
	result := make([]*PlayerEntry, 0)
	for _, entry := range pr.players {
		if entry.RoleID == roleID {
			result = append(result, entry)
		}
	}
	
	return result
}

// Count returns total player count
func (pr *PlayerRegistry) Count() int {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	
	return len(pr.players)
}

// Save saves the registry to disk
func (pr *PlayerRegistry) Save() error {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	
	data := struct {
		Players []*PlayerEntry `json:"players"`
		Count   int            `json:"count"`
	}{
		Players: make([]*PlayerEntry, 0, len(pr.players)),
		Count:   len(pr.players),
	}
	
	for _, entry := range pr.players {
		data.Players = append(data.Players, entry)
	}
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal player data: %w", err)
	}
	
	if err := os.WriteFile(pr.storagePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write player data: %w", err)
	}
	
	return nil
}

// Load loads the registry from disk
func (pr *PlayerRegistry) Load() error {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	
	data, err := os.ReadFile(pr.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No existing data
		}
		return fmt.Errorf("failed to read player data: %w", err)
	}
	
	var loaded struct {
		Players []*PlayerEntry `json:"players"`
	}
	
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal player data: %w", err)
	}
	
	pr.players = make(map[string]*PlayerEntry)
	pr.byName = make(map[string]string)
	
	for _, entry := range loaded.Players {
		if entry.CustomPerms == nil {
			entry.CustomPerms = make(map[string]bool)
		}
		pr.players[entry.ID] = entry
		pr.byName[entry.Name] = entry.ID
	}
	
	return nil
}

// SearchByName searches for players by name (partial match)
func (pr *PlayerRegistry) SearchByName(query string) []*PlayerEntry {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	
	result := make([]*PlayerEntry, 0)
	query = fmt.Sprintf("%%%s%%", query) // Simple contains check
	
	for name, id := range pr.byName {
		// Simple case-insensitive contains
		if len(query) > 2 && containsIgnoreCase(name, query[1:len(query)-1]) {
			result = append(result, pr.players[id])
		}
	}
	
	return result
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		len(s) > 0 && len(substr) > 0 && containsIgnoreCaseHelper(s, substr))
}

func containsIgnoreCaseHelper(s, substr string) bool {
	// Simple implementation - for production use strings.Contains(strings.ToLower(s), strings.ToLower(substr))
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if toLower(s[i+j]) != toLower(substr[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}
