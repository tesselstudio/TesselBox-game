// Package chest implements the chest storage system for TesselBox
package chest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"tesselbox/pkg/config"
	"tesselbox/pkg/items"
)

const (
	ChestSlots = 27 // 3x9 chest inventory
)

// ChestInventory represents the contents of a single chest
type ChestInventory struct {
	X     float64      `json:"x"`
	Y     float64      `json:"y"`
	Slots []items.Item `json:"slots"`
}

// ChestManager manages all chests in the world
type ChestManager struct {
	chests   map[string]*ChestInventory // Key: "x,y" format
	filePath string
	mutex    sync.RWMutex
}

// NewChestManager creates a new chest manager
func NewChestManager(worldName string) *ChestManager {
	// Get chest file path from config
	filePath := config.GetChestFile(worldName)

	return &ChestManager{
		chests:   make(map[string]*ChestInventory),
		filePath: filePath,
	}
}

// GetChestKey generates a unique key for a chest position
func GetChestKey(x, y float64) string {
	return fmt.Sprintf("%.0f,%.0f", x, y)
}

// GetChest gets or creates a chest at the specified position
func (cm *ChestManager) GetChest(x, y float64) *ChestInventory {
	key := GetChestKey(x, y)

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if chest, exists := cm.chests[key]; exists {
		return chest
	}

	// Create new chest
	chest := &ChestInventory{
		X:     x,
		Y:     y,
		Slots: make([]items.Item, ChestSlots),
	}

	// Initialize empty slots
	for i := range chest.Slots {
		chest.Slots[i] = items.Item{Type: items.NONE, Quantity: 0, Durability: -1}
	}

	cm.chests[key] = chest
	return chest
}

// RemoveChest removes a chest at the specified position
func (cm *ChestManager) RemoveChest(x, y float64) {
	key := GetChestKey(x, y)

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	delete(cm.chests, key)
}

// ChestExists checks if a chest exists at the position
func (cm *ChestManager) ChestExists(x, y float64) bool {
	key := GetChestKey(x, y)

	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	_, exists := cm.chests[key]
	return exists
}

// AddItemToChest adds an item to the first available slot in a chest
func (cm *ChestManager) AddItemToChest(x, y float64, itemType items.ItemType, quantity int) bool {
	chest := cm.GetChest(x, y)

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Try to stack with existing items first
	itemProps := items.ItemDefinitions[itemType]
	if itemProps != nil && itemProps.StackSize > 1 {
		for i := range chest.Slots {
			if chest.Slots[i].Type == itemType && chest.Slots[i].Quantity < itemProps.StackSize {
				canAdd := itemProps.StackSize - chest.Slots[i].Quantity
				if canAdd >= quantity {
					chest.Slots[i].Quantity += quantity
					return true
				}
				chest.Slots[i].Quantity = itemProps.StackSize
				quantity -= canAdd
			}
		}
	}

	// Find empty slot for remaining quantity
	for i := range chest.Slots {
		if chest.Slots[i].Type == items.NONE {
			chest.Slots[i] = items.Item{
				Type:       itemType,
				Quantity:   quantity,
				Durability: itemProps.Durability,
			}
			return true
		}
	}

	return false // Chest is full
}

// RemoveItemFromChest removes items from a specific slot
func (cm *ChestManager) RemoveItemFromChest(x, y float64, slotIndex, quantity int) bool {
	chest := cm.GetChest(x, y)

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if slotIndex < 0 || slotIndex >= len(chest.Slots) {
		return false
	}

	if chest.Slots[slotIndex].Type == items.NONE || chest.Slots[slotIndex].Quantity < quantity {
		return false
	}

	chest.Slots[slotIndex].Quantity -= quantity
	if chest.Slots[slotIndex].Quantity <= 0 {
		chest.Slots[slotIndex] = items.Item{Type: items.NONE, Quantity: 0, Durability: -1}
	}

	return true
}

// GetChestContents returns the contents of a chest
func (cm *ChestManager) GetChestContents(x, y float64) []items.Item {
	chest := cm.GetChest(x, y)

	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	// Return a copy
	contents := make([]items.Item, len(chest.Slots))
	copy(contents, chest.Slots)
	return contents
}

// SetChestContents sets the contents of a chest
func (cm *ChestManager) SetChestContents(x, y float64, contents []items.Item) {
	chest := cm.GetChest(x, y)

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if len(contents) <= len(chest.Slots) {
		copy(chest.Slots, contents)
	}
}

// SaveChests saves all chest data to disk
func (cm *ChestManager) SaveChests() error {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	// Ensure directory exists
	dir := filepath.Dir(cm.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create chests directory: %w", err)
	}

	// Convert to serializable format
	chestList := make([]*ChestInventory, 0, len(cm.chests))
	for _, chest := range cm.chests {
		chestList = append(chestList, chest)
	}

	// Marshal and save
	data, err := json.MarshalIndent(chestList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal chests: %w", err)
	}

	if err := os.WriteFile(cm.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write chests file: %w", err)
	}

	return nil
}

// LoadChests loads chest data from disk
func (cm *ChestManager) LoadChests() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Check if file exists
	if _, err := os.Stat(cm.filePath); os.IsNotExist(err) {
		// No chests to load (first time playing)
		return nil
	}

	// Read file
	data, err := os.ReadFile(cm.filePath)
	if err != nil {
		return fmt.Errorf("failed to read chests file: %w", err)
	}

	// Unmarshal
	var chestList []*ChestInventory
	if err := json.Unmarshal(data, &chestList); err != nil {
		return fmt.Errorf("failed to unmarshal chests: %w", err)
	}

	// Populate map
	for _, chest := range chestList {
		key := GetChestKey(chest.X, chest.Y)
		// Ensure slots are initialized
		if len(chest.Slots) == 0 {
			chest.Slots = make([]items.Item, ChestSlots)
			for i := range chest.Slots {
				chest.Slots[i] = items.Item{Type: items.NONE, Quantity: 0, Durability: -1}
			}
		}
		cm.chests[key] = chest
	}

	return nil
}

// GetAllChests returns all chest positions and contents
func (cm *ChestManager) GetAllChests() []*ChestInventory {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	chestList := make([]*ChestInventory, 0, len(cm.chests))
	for _, chest := range cm.chests {
		chestList = append(chestList, chest)
	}
	return chestList
}

// GetChestCount returns the number of chests in the world
func (cm *ChestManager) GetChestCount() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	return len(cm.chests)
}
