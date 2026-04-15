package ui

import "sync"

// GameState represents the current game state
type GameState int

const (
	StateMenu GameState = iota
	StateGame
	StateCrafting
	StateBackpack
	StateChest
	StatePluginUI
	StateSkinEditor
	StateDeathScreen
)

// StateManager manages game state transitions
type StateManager struct {
	current GameState
	mu      sync.RWMutex
}

// NewStateManager creates a new state manager
func NewStateManager() *StateManager {
	return &StateManager{
		current: StateMenu,
	}
}

// GetState returns the current game state
func (sm *StateManager) GetState() GameState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.current
}

// SetState sets the game state
func (sm *StateManager) SetState(state GameState) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.current = state
}

// IsInGame returns true if the game is actively running
func (sm *StateManager) IsInGame() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.current == StateGame
}

// IsInMenu returns true if in menu state
func (sm *StateManager) IsInMenu() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.current == StateMenu
}

// IsModalOpen returns true if any modal is open
func (sm *StateManager) IsModalOpen() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.current != StateGame && sm.current != StateMenu
}

// String returns the string representation of the state
func (s GameState) String() string {
	switch s {
	case StateMenu:
		return "Menu"
	case StateGame:
		return "Game"
	case StateCrafting:
		return "Crafting"
	case StateBackpack:
		return "Backpack"
	case StateChest:
		return "Chest"
	case StatePluginUI:
		return "PluginUI"
	case StateSkinEditor:
		return "SkinEditor"
	case StateDeathScreen:
		return "DeathScreen"
	default:
		return "Unknown"
	}
}
