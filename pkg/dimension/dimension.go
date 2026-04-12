// Package dimension implements the dimension system for TesselBox
// including the Randomland dimension with chaotic terrain generation
package dimension

import (
	"fmt"
	"math/rand"
	"time"

	"tesselbox/pkg/blocks"
	"tesselbox/pkg/enemies"
	"tesselbox/pkg/world"
)

const (
	// Randomland dimensions
	RandomlandWidth  = 500.0
	RandomlandHeight = 500.0
	BedrockThickness = 50.0
	// Generation bounds
	RandomlandMinY = 50.0  // Top bedrock ends here
	RandomlandMaxY = 450.0 // Bottom bedrock starts here
	ReturnPortalX  = 250.0 // Center X
	ReturnPortalY  = 250.0 // Center Y (in the cavern space)
)

// DimensionType represents different dimension types
type DimensionType int

const (
	Overworld DimensionType = iota
	Randomland
)

// Dimension interface for all dimension types
type Dimension interface {
	GetWorld() *world.World
	GetType() DimensionType
	GetName() string
	IsGenerated() bool
}

// RandomlandDimension represents the chaotic Randomland realm
type RandomlandDimension struct {
	World         *world.World
	Type          DimensionType
	Name          string
	Generated     bool
	ReturnPortalX float64
	ReturnPortalY float64
	ZombieSpawner *enemies.ZombieSpawner
	LastVisitTime time.Time
}

// NewRandomlandDimension creates a new Randomland dimension
func NewRandomlandDimension(worldName string) *RandomlandDimension {
	// Use a unique world name to avoid conflicts with player worlds
	dimWorldName := worldName + "__randomland_dim"
	spawner := enemies.NewZombieSpawner(nil) // No day/night cycle in randomland
	// Tune spawn rate for Randomland (2x faster spawning, higher cap)
	spawner.SpawnCooldown = 1500 * time.Millisecond // 1.5s instead of 3s
	spawner.MaxZombies = 25                         // Higher than overworld default of 15
	return &RandomlandDimension{
		World:         world.NewWorld(dimWorldName),
		Type:          Randomland,
		Name:          "Randomland",
		Generated:     false,
		ReturnPortalX: ReturnPortalX,
		ReturnPortalY: ReturnPortalY,
		ZombieSpawner: spawner,
		LastVisitTime: time.Now(),
	}
}

// GetWorld returns the dimension's world
func (r *RandomlandDimension) GetWorld() *world.World {
	return r.World
}

// GetType returns the dimension type
func (r *RandomlandDimension) GetType() DimensionType {
	return r.Type
}

// GetName returns the dimension name
func (r *RandomlandDimension) GetName() string {
	return r.Name
}

// IsGenerated returns whether the dimension has been generated
func (r *RandomlandDimension) IsGenerated() bool {
	return r.Generated
}

// HasReturnPortal checks if the return portal still exists at its position
func (r *RandomlandDimension) HasReturnPortal() bool {
	hex := r.World.GetHexagonAt(r.ReturnPortalX, r.ReturnPortalY)
	if hex == nil {
		return false
	}
	return hex.BlockType == blocks.RANDOMLAND_PORTAL
}

// EnsureReturnPortal recreates the return portal if it was destroyed
func (r *RandomlandDimension) EnsureReturnPortal() {
	if r.HasReturnPortal() {
		return
	}

	fmt.Println("WARNING: Return portal was destroyed! Regenerating...")
	r.createReturnPortal()
}

// Generate generates the randomland terrain if not already generated
// Returns error if generation fails (with panic recovery)
// progressCallback is optional and receives 0.0-1.0 progress updates
func (r *RandomlandDimension) Generate(progressCallback func(float64, string)) (err error) {
	if r.Generated {
		return nil
	}

	// Panic recovery for graceful error handling
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic during Randomland generation: %v", rec)
			fmt.Printf("ERROR: %v\n", err)
			r.Generated = false
		}
	}()

	fmt.Println("Generating Randomland dimension...")

	// Set a fixed seed for Randomland (different from overworld)
	r.World.SetSeed(123456789)

	// Generate bedrock ceiling (top layer) - 20%
	if progressCallback != nil {
		progressCallback(0.1, "Generating bedrock ceiling...")
	}
	r.generateBedrockLayer(0, BedrockThickness)

	// Generate bedrock floor (bottom layer) - 20%
	if progressCallback != nil {
		progressCallback(0.2, "Generating bedrock floor...")
	}
	r.generateBedrockLayer(RandomlandHeight-BedrockThickness, RandomlandHeight)

	// Generate random blocks in the middle - 50%
	if progressCallback != nil {
		progressCallback(0.3, "Generating randomized terrain...")
	}
	r.generateRandomBlocks()

	// Create return portal at center - 80%
	if progressCallback != nil {
		progressCallback(0.8, "Placing return portal...")
	}
	r.createReturnPortal()

	// Spawn zombies in open spaces - 100%
	if progressCallback != nil {
		progressCallback(0.9, "Spawning inhabitants...")
	}
	r.spawnZombies()

	if progressCallback != nil {
		progressCallback(1.0, "Generation complete!")
	}

	r.Generated = true
	fmt.Println("Randomland generation complete!")
	return nil
}

// generateBedrockLayer generates a solid layer of bedrock
func (r *RandomlandDimension) generateBedrockLayer(startY, endY float64) {
	hexSize := world.HexSize
	hexWidth := hexSize * 2
	hexHeight := hexSize * 1.732 // sqrt(3)

	for x := 0.0; x < RandomlandWidth; x += hexWidth * 0.75 {
		for y := startY; y < endY; y += hexHeight {
			// Offset every other row
			actualX := x
			row := int(y / hexHeight)
			if row%2 == 1 {
				actualX += hexWidth * 0.375
			}

			// Only place if within bounds
			if actualX < RandomlandWidth {
				r.World.AddHexagonAt(actualX, y, blocks.BEDROCK)
			}
		}
	}
}

// generateRandomBlocks fills the middle area with completely random blocks
func (r *RandomlandDimension) generateRandomBlocks() {
	hexSize := world.HexSize
	hexWidth := hexSize * 2
	hexHeight := hexSize * 1.732

	// All placeable block types (excluding AIR, BEDROCK, and special blocks)
	placeableBlocks := []blocks.BlockType{
		blocks.DIRT,
		blocks.GRASS,
		blocks.STONE,
		blocks.SAND,
		blocks.WATER,
		blocks.LOG,
		blocks.LEAVES,
		blocks.TROPICAL_LOG,
		blocks.TEMPERATE_LOG,
		blocks.PINE_LOG,
		blocks.TROPICAL_LEAVES,
		blocks.TEMPERATE_LEAVES,
		blocks.PINE_LEAVES,
		blocks.COAL_ORE,
		blocks.IRON_ORE,
		blocks.GOLD_ORE,
		blocks.DIAMOND_ORE,
		blocks.GLASS,
		blocks.BRICK,
		blocks.PLANK,
		blocks.CACTUS,
		blocks.WORKBENCH,
		blocks.FURNACE,
		blocks.ANVIL,
		blocks.GRAVEL,
		blocks.SANDSTONE,
		blocks.OBSIDIAN,
		blocks.ICE,
		blocks.SNOW,
		blocks.TORCH,
		blocks.CRAFTING_TABLE,
		blocks.CHEST,
		blocks.LADDER,
		blocks.FENCE,
		blocks.GATE,
		blocks.DOOR,
		blocks.WINDOW,
		blocks.FLOWER,
		blocks.TALL_GRASS,
		blocks.MUSHROOM_RED,
		blocks.MUSHROOM_BROWN,
		blocks.WOOL,
		blocks.BOOKSHELF,
		blocks.JUKEBOX,
		blocks.NOTE_BLOCK,
		blocks.PUMPKIN,
		blocks.MELON,
		blocks.HAY_BALE,
		blocks.COBBLESTONE,
		blocks.MOSSY_COBBLESTONE,
		blocks.STONE_BRICKS,
		blocks.CHISELED_STONE,
	}

	for x := 0.0; x < RandomlandWidth; x += hexWidth * 0.75 {
		for y := RandomlandMinY; y < RandomlandMaxY; y += hexHeight {
			// Offset every other row
			actualX := x
			row := int((y - RandomlandMinY) / hexHeight)
			if row%2 == 1 {
				actualX += hexWidth * 0.375
			}

			if actualX < RandomlandWidth {
				// 70% chance to place a block (leave some air pockets)
				if rand.Float64() < 0.7 {
					// Pick random block type
					blockType := placeableBlocks[rand.Intn(len(placeableBlocks))]
					r.World.AddHexagonAt(actualX, y, blockType)
				}
			}
		}
	}
}

// createReturnPortal creates the return portal at the center
func (r *RandomlandDimension) createReturnPortal() {
	// Clear area around portal
	portalX := r.ReturnPortalX
	portalY := r.ReturnPortalY

	// Create a small safe zone around return portal (3x3 area)
	for x := portalX - 60; x <= portalX+60; x += 30 {
		for y := portalY - 50; y <= portalY+50; y += 26 {
			// Remove any existing blocks
			r.World.RemoveHexagonAt(x, y)
		}
	}

	// Place the return portal block
	r.World.AddHexagonAt(portalX, portalY, blocks.RANDOMLAND_PORTAL)

	fmt.Printf("Return portal created at (%.1f, %.1f)\n", portalX, portalY)
}

// spawnZombies spawns zombies in open spaces
func (r *RandomlandDimension) spawnZombies() {
	maxZombies := 20
	spawnAttempts := 50
	spawned := 0

	for i := 0; i < spawnAttempts && spawned < maxZombies; i++ {
		// Random position within the cavern
		x := rand.Float64() * RandomlandWidth
		y := RandomlandMinY + rand.Float64()*(RandomlandMaxY-RandomlandMinY)

		// Check if position is valid (not in solid blocks, has ground below)
		hex := r.World.GetHexagonAt(x, y)
		groundHex := r.World.GetHexagonAt(x, y+30)

		// Valid spawn: air at position, solid ground below
		if hex == nil && groundHex != nil && groundHex.BlockType != blocks.AIR {
			// Check distance from return portal (don't spawn too close)
			dx := x - r.ReturnPortalX
			dy := y - r.ReturnPortalY
			distance := dx*dx + dy*dy

			if distance > 10000 { // At least 100 pixels away
				// Create zombie
				zombie := enemies.NewZombie(spawned, enemies.ZombieNormal, x, y)
				r.ZombieSpawner.Zombies = append(r.ZombieSpawner.Zombies, zombie)
				spawned++
			}
		}
	}

	fmt.Printf("Spawned %d zombies in Randomland\n", spawned)
}

// GetSpawnPosition returns a safe spawn position in Randomland
func (r *RandomlandDimension) GetSpawnPosition() (float64, float64) {
	// Return the return portal location (safe zone)
	return r.ReturnPortalX, r.ReturnPortalY - 50 // Spawn slightly above portal
}

// Update updates the dimension (zombies, etc.)
func (r *RandomlandDimension) Update(playerX, playerY float64, deltaTime float64) {
	// Create collision function for zombies
	collisionFunc := func(minX, minY, maxX, maxY float64) bool {
		nearbyHexagons := r.World.GetNearbyHexagons((minX+maxX)/2, (minY+maxY)/2, 100)
		for _, hex := range nearbyHexagons {
			if hex == nil {
				continue
			}
			// Simple AABB collision check
			if hex.X > minX && hex.X < maxX && hex.Y > minY && hex.Y < maxY {
				if hex.BlockType != blocks.AIR {
					return true
				}
			}
		}
		return false
	}

	// Find spawn function
	spawnFunc := func(x, y float64) (float64, float64) {
		return r.World.FindSpawnPosition(x, y)
	}

	// Update zombies (ambient light = 0.5 for randomland - always dim)
	r.ZombieSpawner.Update(deltaTime, nil, 0.5, collisionFunc, spawnFunc)
}

// IsNearReturnPortal checks if a position is near the return portal
func (r *RandomlandDimension) IsNearReturnPortal(x, y float64, tolerance float64) bool {
	dx := x - r.ReturnPortalX
	dy := y - r.ReturnPortalY
	return dx*dx+dy*dy <= tolerance*tolerance
}
