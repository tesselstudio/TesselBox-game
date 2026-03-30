package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	DefaultScreenWidth  = 1024
	DefaultScreenHeight = 768
	GameTitle           = "TesselBox v0.3-alpha"
)

func main() {
	// Parse command line flags
	creativeMode := flag.Bool("creative", false, "Enable creative mode")
	debugMode := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// Create game instance
	game := NewGame(*creativeMode)
	if *debugMode {
		game.showDebug = true
	}

	// Configure ebiten
	ebiten.SetWindowSize(DefaultScreenWidth, DefaultScreenHeight)
	ebiten.SetWindowTitle(GameTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}
