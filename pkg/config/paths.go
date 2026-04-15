package config

import (
	"os"
	"path/filepath"
)

// GetTesselboxDir returns the user's .tesselbox directory
func GetTesselboxDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".tesselbox")
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
