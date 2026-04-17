package config

import (
	"os"
	"path/filepath"
)

// GetTesselboxDir returns the user's tesselbox storage directory
func GetTesselboxDir() string {
	// Use user home directory for application data
	if homeDir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(homeDir, ".tesselbox")
	}
	// Fallback to current directory if home dir can't be determined
	return ".tesselbox"
}

// GetSavesDir returns the saves directory
func GetSavesDir() string {
	return filepath.Join(GetTesselboxDir(), "saves")
}

// GetWorldSaveDir returns the save directory for a specific world
func GetWorldSaveDir(worldName string) string {
	return filepath.Join(GetSavesDir(), worldName)
}

// GetWorldsDir returns the worlds directory
func GetWorldsDir() string {
	return filepath.Join(GetTesselboxDir(), "worlds")
}

// GetSkinsDir returns the skins directory
func GetSkinsDir() string {
	return filepath.Join(GetTesselboxDir(), "skins")
}

// GetChestFile returns the path to the chests file for a world
func GetChestFile(worldName string) string {
	return filepath.Join(GetWorldSaveDir(worldName), "chests.json")
}

// EnsureDirectories creates all necessary directories
func EnsureDirectories() error {
	dirs := []string{
		GetTesselboxDir(),
		GetSavesDir(),
		GetWorldsDir(),
		GetSkinsDir(),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}
