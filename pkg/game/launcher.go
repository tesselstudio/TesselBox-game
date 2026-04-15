package game

import (
	"log"

	"tesselbox/pkg/config"
)

// LaunchConfig contains configuration for launching the game
type LaunchConfig struct {
	WorldName    string
	WorldSeed    int64
	CreativeMode bool
	ScreenWidth  int
	ScreenHeight int
}

// LaunchGUI starts the game with GUI (Ebiten)
// Note: The actual game logic is currently in cmd/main.go
// This launcher provides a clean entry point for future migration
func LaunchGUI(cfg LaunchConfig) error {
	log.Printf("Launching game GUI: world=%s, creative=%v", cfg.WorldName, cfg.CreativeMode)

	if err := config.EnsureDirectories(); err != nil {
		return err
	}

	// TODO: Migrate to GameManager once main.go is refactored
	// For now, this serves as a placeholder for the launcher architecture
	log.Printf("Game launcher ready - migration to GameManager pending")

	return nil
}
