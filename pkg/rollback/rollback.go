package rollback

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ChangeType represents type of change
type ChangeType int

const (
	ChangeBlockPlace ChangeType = iota
	ChangeBlockBreak
	ChangeInventoryAdd
	ChangeInventoryRemove
)

// BlockChange represents a block modification
type BlockChange struct {
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	WorldID string  `json:"world_id"`
	OldType string  `json:"old_type"`
	NewType string  `json:"new_type"`
	OldData string  `json:"old_data,omitempty"`
	NewData string  `json:"new_data,omitempty"`
}

// InventoryChange represents inventory modification
type InventoryChange struct {
	PlayerID string `json:"player_id"`
	Slot     int    `json:"slot"`
	OldItem  string `json:"old_item"`
	NewItem  string `json:"new_item"`
	Quantity int    `json:"quantity"`
}

// RollbackEntry represents a single change
type RollbackEntry struct {
	ID        string     `json:"id"`
	Timestamp time.Time  `json:"timestamp"`
	PlayerID  string     `json:"player_id"`
	Type      ChangeType `json:"type"`

	// Either block or inventory change
	BlockChange     *BlockChange     `json:"block_change,omitempty"`
	InventoryChange *InventoryChange `json:"inventory_change,omitempty"`
}

// RollbackLog tracks changes for rollback
type RollbackLog struct {
	WorldID    string          `json:"world_id"`
	Entries    []RollbackEntry `json:"entries"`
	MaxEntries int             `json:"max_entries"`
}

// NewRollbackLog creates a new log
func NewRollbackLog(worldID string) *RollbackLog {
	return &RollbackLog{
		WorldID:    worldID,
		Entries:    make([]RollbackEntry, 0),
		MaxEntries: 10000,
	}
}

// AddEntry adds a change to the log
func (rl *RollbackLog) AddEntry(entry RollbackEntry) {
	rl.Entries = append(rl.Entries, entry)

	// Trim old entries
	if len(rl.Entries) > rl.MaxEntries {
		rl.Entries = rl.Entries[len(rl.Entries)-rl.MaxEntries:]
	}
}

// GetChangesSince returns changes after a timestamp
func (rl *RollbackLog) GetChangesSince(since time.Time, playerID string) []RollbackEntry {
	result := make([]RollbackEntry, 0)

	for _, entry := range rl.Entries {
		if entry.Timestamp.After(since) {
			if playerID == "" || entry.PlayerID == playerID {
				result = append(result, entry)
			}
		}
	}

	return result
}

// GetChangesInRegion returns changes in an area
func (rl *RollbackLog) GetChangesInRegion(minX, minY, maxX, maxY float64, since time.Time) []RollbackEntry {
	result := make([]RollbackEntry, 0)

	for _, entry := range rl.Entries {
		if entry.Timestamp.After(since) && entry.BlockChange != nil {
			bc := entry.BlockChange
			if bc.X >= minX && bc.X <= maxX && bc.Y >= minY && bc.Y <= maxY {
				result = append(result, entry)
			}
		}
	}

	return result
}

// RollbackManager manages rollback data
type RollbackManager struct {
	logs map[string]*RollbackLog // worldID -> log

	storagePath string
}

// NewRollbackManager creates new manager
func NewRollbackManager(storageDir string) *RollbackManager {
	return &RollbackManager{
		logs:        make(map[string]*RollbackLog),
		storagePath: filepath.Join(storageDir, "rollback"),
	}
}

// GetOrCreateLog gets or creates log for world
func (rm *RollbackManager) GetOrCreateLog(worldID string) *RollbackLog {
	if log, exists := rm.logs[worldID]; exists {
		return log
	}

	log := NewRollbackLog(worldID)
	rm.logs[worldID] = log
	return log
}

// LogBlockChange logs a block change
func (rm *RollbackManager) LogBlockChange(worldID, playerID string, x, y float64, oldType, newType string) {
	log := rm.GetOrCreateLog(worldID)

	entry := RollbackEntry{
		ID:        fmt.Sprintf("rb_%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		PlayerID:  playerID,
		Type:      ChangeBlockPlace,
		BlockChange: &BlockChange{
			X:       x,
			Y:       y,
			WorldID: worldID,
			OldType: oldType,
			NewType: newType,
		},
	}

	if newType == "" {
		entry.Type = ChangeBlockBreak
	}

	log.AddEntry(entry)
}

// LogInventoryChange logs inventory change
func (rm *RollbackManager) LogInventoryChange(worldID, playerID string, slot int, oldItem, newItem string, quantity int) {
	log := rm.GetOrCreateLog(worldID)

	entryType := ChangeInventoryAdd
	if quantity < 0 {
		entryType = ChangeInventoryRemove
	}

	entry := RollbackEntry{
		ID:        fmt.Sprintf("rb_%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		PlayerID:  playerID,
		Type:      entryType,
		InventoryChange: &InventoryChange{
			PlayerID: playerID,
			Slot:     slot,
			OldItem:  oldItem,
			NewItem:  newItem,
			Quantity: quantity,
		},
	}

	log.AddEntry(entry)
}

// RollbackPlayer rolls back a player's changes
func (rm *RollbackManager) RollbackPlayer(worldID, playerID string, since time.Time) []RollbackEntry {
	log := rm.GetOrCreateLog(worldID)
	changes := log.GetChangesSince(since, playerID)

	// In real implementation, actually reverse the changes
	// For now, just return the list

	return changes
}

// RollbackRegion rolls back changes in a region
func (rm *RollbackManager) RollbackRegion(worldID string, minX, minY, maxX, maxY float64, since time.Time) []RollbackEntry {
	log := rm.GetOrCreateLog(worldID)
	changes := log.GetChangesInRegion(minX, minY, maxX, maxY, since)

	// In real implementation, reverse the changes

	return changes
}

// UndoLastChange undoes the most recent change by a player
func (rm *RollbackManager) UndoLastChange(worldID, playerID string) *RollbackEntry {
	log := rm.GetOrCreateLog(worldID)

	// Find last change by player
	for i := len(log.Entries) - 1; i >= 0; i-- {
		if log.Entries[i].PlayerID == playerID {
			// In real implementation, undo this change
			return &log.Entries[i]
		}
	}

	return nil
}

// GetRecentChanges gets recent changes for inspection
func (rm *RollbackManager) GetRecentChanges(worldID string, count int) []RollbackEntry {
	log := rm.GetOrCreateLog(worldID)

	if count > len(log.Entries) {
		count = len(log.Entries)
	}

	start := len(log.Entries) - count
	if start < 0 {
		start = 0
	}

	return log.Entries[start:]
}

// Save saves rollback logs
func (rm *RollbackManager) Save() error {
	for worldID, log := range rm.logs {
		path := filepath.Join(rm.storagePath, fmt.Sprintf("%s.json", worldID))

		if err := os.MkdirAll(rm.storagePath, 0755); err != nil {
			return err
		}

		data, err := json.MarshalIndent(log, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal: %w", err)
		}

		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}
	}

	return nil
}

// Load loads rollback logs
func (rm *RollbackManager) Load() error {
	files, err := os.ReadDir(rm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		worldID := file.Name()[:len(file.Name())-5] // Remove .json
		path := filepath.Join(rm.storagePath, file.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var log RollbackLog
		if err := json.Unmarshal(data, &log); err != nil {
			continue
		}

		rm.logs[worldID] = &log
	}

	return nil
}

// ClearOldEntries removes entries older than duration
func (rm *RollbackManager) ClearOldEntries(maxAge time.Duration) int {
	removed := 0
	cutoff := time.Now().Add(-maxAge)

	for _, log := range rm.logs {
		newEntries := make([]RollbackEntry, 0)
		for _, entry := range log.Entries {
			if entry.Timestamp.After(cutoff) {
				newEntries = append(newEntries, entry)
			} else {
				removed++
			}
		}
		log.Entries = newEntries
	}

	return removed
}
