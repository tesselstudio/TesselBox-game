package main

import (
	"log"
	"tesselbox/pkg/game"

	"github.com/hajimehoshi/ebiten/v2"
)

// GameLifecycle manages the game startup and shutdown
type GameLifecycle struct {
	wrapper *GameWrapper
}

// NewGameLifecycle creates a new game lifecycle manager
func NewGameLifecycle(worldName string, worldSeed int64, creativeMode bool) (*GameLifecycle, error) {
	// Initialize storage
	if err := initTesselboxStorage(); err != nil {
		return nil, err
	}

	// Create game wrapper
	wrapper := NewGameWrapper(worldName, worldSeed, creativeMode, ScreenWidth, ScreenHeight)

	// Start background music
	if err := wrapper.manager.StartBackgroundMusic(); err != nil {
		log.Printf("Warning: Failed to start background music: %v", err)
	}

	return &GameLifecycle{
		wrapper: wrapper,
	}, nil
}

// Run starts the game loop
func (gl *GameLifecycle) Run() error {
	log.Printf("Starting game loop...")

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Tesselbox")

	if err := ebiten.RunGame(gl.wrapper); err != nil {
		return err
	}

	return nil
}

// Shutdown cleanly shuts down the game
func (gl *GameLifecycle) Shutdown() error {
	log.Printf("Shutting down game...")

	if gl.wrapper != nil {
		gl.wrapper.manager.StopBackgroundMusic()
		gl.wrapper.manager.Cleanup()
	}

	return nil
}

// GetManager returns the underlying GameManager
func (gl *GameLifecycle) GetManager() *game.GameManager {
	if gl.wrapper != nil {
		return gl.wrapper.manager
	}
	return nil
}

// LaunchGame is the main entry point for starting the game
func LaunchGame(worldName string, worldSeed int64, creativeMode bool) error {
	log.Printf("Launching Tesselbox game...")
	log.Printf("World: %s, Seed: %d, Creative: %v", worldName, worldSeed, creativeMode)

	// Create lifecycle manager
	lifecycle, err := NewGameLifecycle(worldName, worldSeed, creativeMode)
	if err != nil {
		return err
	}

	// Run the game
	if err := lifecycle.Run(); err != nil {
		return err
	}

	// Shutdown
	return lifecycle.Shutdown()
}
