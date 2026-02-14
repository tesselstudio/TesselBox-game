package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SaveData represents what we want to persist for each player
type SaveData struct {
	UserID      string    `json:"user_id"`
	PlayerX     float64   `json:"player_x"`
	PlayerY     float64   `json:"player_y"`
	WorldName   string    `json:"world_name"`
	MinedBlocks []string  `json:"mined_blocks"` // Format: "x,y"
	Inventory   []string  `json:"inventory"`    // Item IDs
	Health      float64   `json:"health"`
	LastSaved   time.Time `json:"last_saved"`
	Version     string    `json:"version"`
}

// SavePlayerProgress saves the game state to a local JSON file named after the user
func (s *AuthService) SavePlayerProgress(data SaveData) error {
	if s.currentUser == nil {
		return fmt.Errorf("cannot save: no user logged in")
	}

	// Create a saves directory if it doesn't exist
	saveDir := "saves"
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		os.Mkdir(saveDir, 0755)
	}

	// Use a sanitized version of the User ID as the filename
	filename := filepath.Join(saveDir, fmt.Sprintf("user_%s.json", s.currentUser.ID))

	// Set metadata
	data.UserID = s.currentUser.ID
	data.LastSaved = time.Now()
	data.Version = "1.0"

	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, file, 0644)
}

// LoadPlayerProgress retrieves the saved state for the current user
func (s *AuthService) LoadPlayerProgress() (*SaveData, error) {
	if s.currentUser == nil {
		return nil, fmt.Errorf("cannot load: no user logged in")
	}

	filename := filepath.Join("saves", fmt.Sprintf("user_%s.json", s.currentUser.ID))

	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data SaveData
	err = json.Unmarshal(file, &data)
	return &data, err
}

// GetPlayerSaveData creates SaveData from current player state
func (s *AuthService) GetPlayerSaveData(playerX, playerY float64, worldName string, minedBlocks []string, inventory []string, health float64) SaveData {
	return SaveData{
		UserID:      s.currentUser.ID,
		PlayerX:     playerX,
		PlayerY:     playerY,
		WorldName:   worldName,
		MinedBlocks: minedBlocks,
		Inventory:   inventory,
		Health:      health,
		LastSaved:   time.Now(),
		Version:     "1.0",
	}
}

// SavePlayerWorld saves player progress with world integration
func (s *AuthService) SavePlayerWorld(playerX, playerY float64, worldName string, minedBlocks []string, inventory []string, health float64) error {
	data := s.GetPlayerSaveData(playerX, playerY, worldName, minedBlocks, inventory, health)
	return s.SavePlayerProgress(data)
}

// LoadPlayerWorld loads player progress and returns world data
func (s *AuthService) LoadPlayerWorld() (*SaveData, error) {
	return s.LoadPlayerProgress()
}

// DeletePlayerSave removes a player's save file
func (s *AuthService) DeletePlayerSave() error {
	if s.currentUser == nil {
		return fmt.Errorf("cannot delete: no user logged in")
	}

	filename := filepath.Join("saves", fmt.Sprintf("user_%s.json", s.currentUser.ID))

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	return os.Remove(filename)
}

// ListPlayerSaves returns a list of all player save files
func ListPlayerSaves() ([]string, error) {
	saveDir := "saves"

	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(saveDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read saves directory: %w", err)
	}

	var saves []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			saves = append(saves, entry.Name())
		}
	}

	return saves, nil
}
