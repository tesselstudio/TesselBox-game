package world

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"tesselbox/pkg/biomes"
	"tesselbox/pkg/blocks"
	"tesselbox/pkg/hexagon"
	"tesselbox/pkg/organisms"
)

const (
	// RenderDistance is the number of chunks to render around the player
	RenderDistance = 4
	// ChunkUnloadDistance is the distance at which chunks are unloaded
	ChunkUnloadDistance = 10
)

const (
	// Spatial hash constants
	SpatialHashCellSize = 100.0 // Size of each spatial hash cell
)

// World represents the game world
type World struct {
	Chunks    map[[2]int]*Chunk
	Seed      int64
	Organisms []*organisms.Organism
	Storage   *WorldStorage
	WorldName string

	// Spatial hash for optimized collision detection
	spatialHash map[[2]int][]*Hexagon // Cell coordinates -> list of hexagons
}

// NewWorld creates a new world
func NewWorld(worldName string) *World {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	world := &World{
		Chunks:      make(map[[2]int]*Chunk),
		Seed:        seed,
		Organisms:   []*organisms.Organism{},
		Storage:     NewWorldStorage(worldName),
		WorldName:   worldName,
		spatialHash: make(map[[2]int][]*Hexagon),
	}

	return world
}

// NewWorldFromStorage creates a world and loads it from storage
func NewWorldFromStorage(worldName string) (*World, error) {
	world := &World{
		Chunks:    make(map[[2]int]*Chunk),
		Organisms: []*organisms.Organism{},
		Storage:   NewWorldStorage(worldName),
		WorldName: worldName,
	}

	// Load metadata to get seed
	metadata, err := world.Storage.GetWorldMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to load world metadata: %w", err)
	}

	world.Seed = metadata.CreatedAt.UnixNano()
	rand.Seed(world.Seed)

	return world, nil
}

// GetChunkCoords returns the chunk coordinates for a given world position
func (w *World) GetChunkCoords(x, y float64) (int, int) {
	chunkX := int(math.Floor(x / GetChunkWidth()))
	chunkY := int(math.Floor(y / GetChunkHeight()))
	return chunkX, chunkY
}

// getSpatialHashKey calculates the spatial hash key for a given position
func (w *World) getSpatialHashKey(x, y float64) [2]int {
	cellX := int(math.Floor(x / SpatialHashCellSize))
	cellY := int(math.Floor(y / SpatialHashCellSize))
	return [2]int{cellX, cellY}
}

// addHexagonToSpatialHash adds a hexagon to the spatial hash
func (w *World) addHexagonToSpatialHash(hex *Hexagon) {
	key := w.getSpatialHashKey(hex.X, hex.Y)
	w.spatialHash[key] = append(w.spatialHash[key], hex)
}

// removeHexagonFromSpatialHash removes a hexagon from the spatial hash
func (w *World) removeHexagonFromSpatialHash(hex *Hexagon) {
	key := w.getSpatialHashKey(hex.X, hex.Y)
	hexagons := w.spatialHash[key]

	// Find and remove the hexagon from the slice
	for i, h := range hexagons {
		if h == hex {
			// Remove element at index i
			w.spatialHash[key] = append(hexagons[:i], hexagons[i+1:]...)
			break
		}
	}

	// Clean up empty cells
	if len(w.spatialHash[key]) == 0 {
		delete(w.spatialHash, key)
	}
}

// rebuildSpatialHash rebuilds the entire spatial hash from all chunks
func (w *World) rebuildSpatialHash() {
	w.spatialHash = make(map[[2]int][]*Hexagon)

	for _, chunk := range w.Chunks {
		for _, hex := range chunk.Hexagons {
			w.addHexagonToSpatialHash(hex)
		}
	}
}

// GetChunk returns the chunk at the given chunk coordinates
func (w *World) GetChunk(chunkX, chunkY int) *Chunk {
	key := [2]int{chunkX, chunkY}
	chunk, exists := w.Chunks[key]
	if !exists {
		// Try to load from storage first
		loadedChunk, err := w.Storage.LoadChunk(chunkX, chunkY)
		if err != nil {
			// If loading fails, generate new chunk
			chunk = NewChunk(chunkX, chunkY)
			w.generateChunk(chunk)
		} else if loadedChunk != nil {
			chunk = loadedChunk
		} else {
			// No saved chunk exists, generate new one
			chunk = NewChunk(chunkX, chunkY)
			w.generateChunk(chunk)
		}

		w.Chunks[key] = chunk

		// Add all hexagons from this chunk to the spatial hash
		for _, hex := range chunk.Hexagons {
			w.addHexagonToSpatialHash(hex)
		}
	}
	chunk.LastAccessed = time.Now()
	return chunk
}

// generateChunk generates terrain for a chunk with biome integration
func (w *World) generateChunk(chunk *Chunk) {
	worldX, worldY := chunk.GetWorldPosition()

	// Create noise generator for this world
	noise := biomes.NewSimplexNoise(float64(w.Seed))

	for row := 0; row < ChunkSize; row++ {
		for col := 0; col < ChunkSize; col++ {
			var x, y float64

			// Calculate hexagon position with interlocking pattern
			if row%2 == 0 {
				x = worldX + float64(col)*HexWidth + HexWidth/2
			} else {
				x = worldX + float64(col)*HexWidth + HexWidth
			}
			y = worldY + float64(row)*HexVSpacing + HexSize

			// Get biome at this position
			biomeType := biomes.GetBiomeAtPosition(x, y, noise)
			biomeProps := biomes.BiomeDefinitions[biomeType]

			// Base terrain height varies by biome
			var baseHeight float64
			switch biomeType {
			case biomes.OCEAN:
				baseHeight = 550.0 // Lower for oceans
			case biomes.DESERT:
				baseHeight = 400.0
			case biomes.MOUNTAINS:
				baseHeight = 350.0 // Higher for mountains
			case biomes.SWAMP:
				baseHeight = 420.0
			default:
				baseHeight = 400.0
			}

			// Large scale terrain features
			terrainNoise := math.Sin(x*0.005) * math.Cos(y*0.005) * 150
			terrainNoise += math.Sin(x*0.01+100) * math.Cos(y*0.01) * 75

			// Medium scale variation
			variationNoise := math.Sin(x*0.02) * math.Cos(y*0.02) * 25

			// Small scale detail
			detailNoise := math.Sin(x*0.05) * math.Cos(y*0.05) * 10

			// Combine all noise layers
			totalNoise := terrainNoise + variationNoise + detailNoise
			surfaceY := baseHeight + totalNoise

			// Determine block type based on depth and biome
			depth := y - surfaceY

			var blockType blocks.BlockType

			if depth < -10 {
				// Above surface - air (unless it's an ocean)
				if biomeType == biomes.OCEAN && depth < 0 && depth > -60 {
					blockType = blocks.WATER
				} else {
					blockType = blocks.AIR
				}
			} else if depth < 5 {
				// Surface layer - determine by biome
				switch biomeType {
				case biomes.DESERT:
					blockType = blocks.SAND
				case biomes.OCEAN:
					blockType = blocks.SAND
				case biomes.SWAMP:
					blockType = blocks.GRASS
				case biomes.MOUNTAINS:
					blockType = blocks.STONE
				default:
					blockType = blocks.GRASS
				}
			} else if depth < 15 {
				// Subsurface layer
				switch biomeType {
				case biomes.DESERT, biomes.OCEAN:
					blockType = blocks.SAND
				case biomes.SWAMP:
					blockType = blocks.DIRT
				default:
					blockType = blocks.DIRT
				}
			} else if depth < 200 {
				// Stone layers with ore generation
				// Use biome ore frequency modifier
				oreFrequency := biomeProps.OreFrequency

				// Add random ore generation
				rand.Seed(w.Seed + int64(x)*1000 + int64(y))
				oreChance := rand.Float64()

				// Adjust ore chance by biome frequency
				if depth > 20 && depth < 50 && oreChance < 0.02*oreFrequency {
					blockType = blocks.COAL_ORE
				} else if depth > 30 && depth < 70 && oreChance < 0.015*oreFrequency {
					blockType = blocks.IRON_ORE
				} else if depth > 40 && depth < 60 && oreChance < 0.008*oreFrequency {
					blockType = blocks.GOLD_ORE
				} else if depth > 50 && depth < 80 && oreChance < 0.004*oreFrequency {
					blockType = blocks.DIAMOND_ORE
				} else {
					blockType = blocks.STONE
				}
			} else {
				// Deep stone or bedrock at very bottom
				if depth > 300 {
					blockType = blocks.BEDROCK
				} else {
					blockType = blocks.STONE
				}
			}

			// Create hexagon (only if not air)
			if blockType != blocks.AIR {
				hexagon := NewHexagon(x, y, HexSize, blockType)
				chunk.AddHexagon(x, y, hexagon)
			}

			// Spawn trees in appropriate biomes
			if depth >= -5 && depth <= 5 && biomeType != biomes.OCEAN && biomeType != biomes.DESERT {
				if biomes.ShouldSpawnTree(x, y, noise) {
					// This is a simplified tree spawn
					// In a full implementation, you'd generate the entire tree structure
					// For now, we just mark this position for potential tree generation
				}
			}
		}
	}

	chunk.Modified = false
}

// GetNearbyHexagons returns hexagons within a radius of the given position (optimized with spatial hash)
func (w *World) GetNearbyHexagons(x, y, radius float64) []*Hexagon {
	hexagons := []*Hexagon{}

	// Calculate the range of spatial hash cells to check
	minX := x - radius
	maxX := x + radius
	minY := y - radius
	maxY := y + radius

	minCellX := int(math.Floor(minX / SpatialHashCellSize))
	maxCellX := int(math.Floor(maxX / SpatialHashCellSize))
	minCellY := int(math.Floor(minY / SpatialHashCellSize))
	maxCellY := int(math.Floor(maxY / SpatialHashCellSize))

	radiusSq := radius * radius

	// Check all relevant spatial hash cells
	for cellX := minCellX; cellX <= maxCellX; cellX++ {
		for cellY := minCellY; cellY <= maxCellY; cellY++ {
			key := [2]int{cellX, cellY}
			cellHexagons := w.spatialHash[key]

			// Check each hexagon in this cell
			for _, hex := range cellHexagons {
				hx := hex.X - x
				hy := hex.Y - y
				if hx*hx+hy*hy <= radiusSq {
					hexagons = append(hexagons, hex)
				}
			}
		}
	}

	return hexagons
}

// GetHexagonAt returns the hexagon at the given world position
func (w *World) GetHexagonAt(x, y float64) *Hexagon {
	chunkX, chunkY := w.GetChunkCoords(x, y)
	chunk := w.GetChunk(chunkX, chunkY)
	return chunk.GetHexagon(x, y)
}

// AddHexagonAt adds a hexagon at the given world position
func (w *World) AddHexagonAt(x, y float64, blockType blocks.BlockType) {
	centerX, centerY, _, _ := PixelToHexCenter(x, y)
	hexagon := NewHexagon(centerX, centerY, HexSize, blockType)

	chunkX, chunkY := w.GetChunkCoords(centerX, centerY)
	chunk := w.GetChunk(chunkX, chunkY)
	chunk.AddHexagon(centerX, centerY, hexagon)

	// Add to spatial hash
	w.addHexagonToSpatialHash(hexagon)
}

// RemoveHexagonAt removes the hexagon at the given world position
func (w *World) RemoveHexagonAt(x, y float64) bool {
	chunkX, chunkY := w.GetChunkCoords(x, y)
	chunk := w.GetChunk(chunkX, chunkY)

	// Get the hexagon before removing it
	hexagon := chunk.GetHexagon(x, y)
	if hexagon == nil {
		return false
	}

	// Remove from spatial hash first
	w.removeHexagonFromSpatialHash(hexagon)

	// Then remove from chunk
	return chunk.RemoveHexagon(x, y)
}

// UnloadDistantChunks unloads chunks that are far from the player
func (w *World) UnloadDistantChunks(playerX, playerY float64) {
	playerChunkX, playerChunkY := w.GetChunkCoords(playerX, playerY)
	toDelete := [][2]int{}

	for key, chunk := range w.Chunks {
		dx := chunk.ChunkX - playerChunkX
		dy := chunk.ChunkY - playerChunkY
		distance := math.Sqrt(float64(dx*dx + dy*dy))

		if distance > ChunkUnloadDistance {
			toDelete = append(toDelete, key)
		}
	}

	// Save and unload distant chunks
	for _, key := range toDelete {
		chunk := w.Chunks[key]

		// Save modified chunks before unloading
		if chunk.Modified {
			err := w.Storage.SaveChunk(chunk)
			if err != nil {
				// Log error but continue - don't prevent unloading due to save failure
				fmt.Printf("Warning: Failed to save chunk %d,%d before unloading: %v\n", key[0], key[1], err)
			}
		}

		// Remove all hexagons from this chunk from the spatial hash
		for _, hex := range chunk.Hexagons {
			w.removeHexagonFromSpatialHash(hex)
		}

		// Remove chunk from memory
		delete(w.Chunks, key)
	}
}

// GetNearbyOrganisms returns organisms within a radius of the given position
func (w *World) GetNearbyOrganisms(x, y, radius float64) []*organisms.Organism {
	nearby := []*organisms.Organism{}
	radiusSq := radius * radius

	for _, org := range w.Organisms {
		dx := org.X - x
		dy := org.Y - y
		if dx*dx+dy*dy <= radiusSq {
			nearby = append(nearby, org)
		}
	}

	return nearby
}

// GetVisibleBlocks returns all visible blocks based on camera position
func (w *World) GetVisibleBlocks(cameraX, cameraY float64) []*Hexagon {
	return w.GetNearbyHexagons(cameraX, cameraY, float64(RenderDistance)*GetChunkWidth())
}

// GetChunksInRange ensures chunks are generated around the given position
func (w *World) GetChunksInRange(x, y float64) {
	chunkX, chunkY := w.GetChunkCoords(x, y)
	chunkRadius := RenderDistance

	for dx := -chunkRadius; dx <= chunkRadius; dx++ {
		for dy := -chunkRadius; dy <= chunkRadius; dy++ {
			w.GetChunk(chunkX+dx, chunkY+dy)
		}
	}
}

// SetBlock sets a block at the given hex coordinates
func (w *World) SetBlock(hex hexagon.Hexagon, chunkZ int, blockType blocks.BlockType) {
	// Convert hex coordinates to world position
	x, y := hexagon.HexToPixel(hex, HexSize)

	if blockType == blocks.AIR {
		w.RemoveHexagonAt(x, y)
	} else {
		w.AddHexagonAt(x, y, blockType)
	}
}

// DamageBlock applies damage to a block and returns true if destroyed
func (w *World) DamageBlock(hex hexagon.Hexagon, chunkZ int, damage float64) bool {
	// Convert hex coordinates to world position
	x, y := hexagon.HexToPixel(hex, HexSize)

	h := w.GetHexagonAt(x, y)
	if h == nil {
		return false
	}

	h.Health -= damage
	if h.Health <= 0 {
		w.RemoveHexagonAt(x, y)
		return true
	}
	return false
}
func (w *World) GetOrganismAt(x, y, tolerance float64) *organisms.Organism {
	toleranceSq := tolerance * tolerance

	for _, org := range w.Organisms {
		dx := org.X - x
		dy := org.Y - y
		if dx*dx+dy*dy <= toleranceSq {
			return org
		}
	}

	return nil
}

// SaveWorld saves the current world state to storage
func (w *World) SaveWorld() error {
	if w.Storage == nil {
		return fmt.Errorf("world storage not initialized")
	}

	err := w.Storage.SaveWorld(w)
	if err != nil {
		return err
	}

	// Save metadata
	return w.Storage.SaveWorldMetadata(len(w.Chunks))
}

// LoadWorldArea loads chunks around a specific position from storage
func (w *World) LoadWorldArea(centerX, centerY float64, radius int) error {
	if w.Storage == nil {
		return fmt.Errorf("world storage not initialized")
	}

	return w.Storage.LoadWorld(w, centerX, centerY, radius)
}

// AutoSave periodically saves the world
func (w *World) AutoSave(interval time.Duration, stopChan <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := w.SaveWorld(); err != nil {
				fmt.Printf("Auto-save failed: %v\n", err)
			} else {
				fmt.Printf("World auto-saved successfully\n")
			}
		case <-stopChan:
			return
		}
	}
}

// RemoveOrganism removes an organism from the world
func (w *World) RemoveOrganism(x, y float64) {
	for i, org := range w.Organisms {
		if org.X == x && org.Y == y {
			w.Organisms = append(w.Organisms[:i], w.Organisms[i+1:]...)
			return
		}
	}
}
