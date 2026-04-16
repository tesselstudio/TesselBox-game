package main

import (
	"flag"
	"log"
	"os"
)

// NewMain is the new simplified main function using GameManager
// This demonstrates the refactored architecture
func NewMain() {
	// Parse command line flags
	worldName := flag.String("world", "default", "World name to load or create")
	worldSeed := flag.Int64("seed", 0, "World seed (0 for random)")
	creativeMode := flag.Bool("creative", false, "Enable creative mode")
	flag.Parse()

	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Tesselbox Game Engine")
	log.Printf("World: %s, Seed: %d, Creative: %v", *worldName, *worldSeed, *creativeMode)

	// Launch the game
	if err := LaunchGame(*worldName, *worldSeed, *creativeMode); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

// Note: The original main() function in main.go should be updated to call NewMain()
// or this can be used as a reference for refactoring the existing main()
