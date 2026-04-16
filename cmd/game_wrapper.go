package main

import (
	"log"
	"time"

	"tesselbox/pkg/game"

	"github.com/hajimehoshi/ebiten/v2"
)

// GameWrapper wraps the GameManager for Ebiten compatibility
type GameWrapper struct {
	manager      *game.GameManager
	screenWidth  int
	screenHeight int
	lastTime     time.Time
}

// NewGameWrapper creates a new game wrapper
func NewGameWrapper(worldName string, worldSeed int64, creativeMode bool, screenWidth, screenHeight int) *GameWrapper {
	gm := game.NewGameManager(worldName, worldSeed, creativeMode)

	return &GameWrapper{
		manager:      gm,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		lastTime:     time.Now(),
	}
}

// Update implements ebiten.Game
func (gw *GameWrapper) Update() error {
	currentTime := time.Now()
	deltaTime := currentTime.Sub(gw.lastTime).Seconds()
	gw.lastTime = currentTime

	// Update the game manager
	if err := gw.manager.Update(deltaTime); err != nil {
		return err
	}

	return nil
}

// Draw implements ebiten.Game
func (gw *GameWrapper) Draw(screen *ebiten.Image) {
	// TODO: Implement drawing logic using gw.manager fields
	// For now, this is a placeholder
	// The actual drawing logic from the original Game.Draw() should be moved here
	// and adapted to use gw.manager instead of g
}

// Layout implements ebiten.Game
func (gw *GameWrapper) Layout(outsideWidth, outsideHeight int) (int, int) {
	return gw.screenWidth, gw.screenHeight
}

// GetManager returns the underlying GameManager
func (gw *GameWrapper) GetManager() *game.GameManager {
	return gw.manager
}

// Close cleans up resources
func (gw *GameWrapper) Close() error {
	log.Printf("Closing game wrapper...")
	gw.manager.Cleanup()
	return nil
}
