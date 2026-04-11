package world

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"tesselbox/pkg/biomes"
	"tesselbox/pkg/blocks"
	"tesselbox/pkg/creatures"
	"tesselbox/pkg/gametime"
	"tesselbox/pkg/hexagon"
	"tesselbox/pkg/organisms"
)

const (
	// RenderDistance is the number of chunks to render around the player
	RenderDistance = 3
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
	Creatures []*creatures.Creature
	Storage   *WorldStorage
	WorldName string

	// Spatial hash for optimized collision detection
	spatialHash map[[2]int][]*Hexagon // Cell coordinates -> list of hexagons

	// Noise generator for terrain generation (cached for performance)
	noiseGenerator *biomes.SimplexNoise

	// Loading state to prevent deadlocks
	loadingChunks map[[2]int]bool // Fixed: Track chunks being loaded
	loadingMutex  sync.Mutex      // Fixed: Protect loading state
}

// NewWorld creates a new world
func NewWorld(worldName string) *World {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	world := &World{
		Chunks:        make(map[[2]int]*Chunk),
		Seed:          seed,
		Organisms:     []*organisms.Organism{},
		Creatures:     []*creatures.Creature{},
		Storage:       NewWorldStorage(worldName),
		WorldName:     worldName,
		spatialHash:   make(map[[2]int][]*Hexagon),
		loadingChunks: make(map[[2]int]bool),
	}

	// Initialize noise generator for terrain generation
	world.noiseGenerator = biomes.NewSimplexNoise(float64(seed))

	return world
}

// SetSeed sets the world seed and regenerates the noise generator
func (w *World) SetSeed(seed int64) {
	w.Seed = seed
	rand.Seed(seed)
	w.noiseGenerator = biomes.NewSimplexNoise(float64(seed))

	// Clear existing chunks to force regeneration with new seed
	w.Chunks = make(map[[2]int]*Chunk)
	w.spatialHash = make(map[[2]int][]*Hexagon)
}

// GetSeed returns the current world seed
func (w *World) GetSeed() int64 {
	return w.Seed
}

// GetNoiseGenerator returns the world's noise generator
func (w *World) GetNoiseGenerator() *biomes.SimplexNoise {
	return w.noiseGenerator
}

// ValidateSeed checks if a seed is valid
func (w *World) ValidateSeed(seed int64) bool {
	// Seeds can be any int64, but we could add validation here
	return seed != 0 // Avoid zero seed for better randomness
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
	if hex == nil {
		return // Fixed: Add nil check
	}

	key := w.getSpatialHashKey(hex.X, hex.Y)
	hexagons := w.spatialHash[key]

	// Find and remove the hexagon from the slice using coordinates instead of pointer comparison
	for i, h := range hexagons {
		if h != nil && h.X == hex.X && h.Y == hex.Y {
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

	// Check if chunk already exists
	chunk, exists := w.Chunks[key]
	if exists {
		chunk.LastAccessed = time.Now()
		return chunk
	}

	// Check if chunk is currently being loaded to prevent deadlock
	w.loadingMutex.Lock()
	if w.loadingChunks[key] {
		w.loadingMutex.Unlock()
		// Chunk is being loaded, wait and return existing or nil
		if chunk, exists := w.Chunks[key]; exists {
			chunk.LastAccessed = time.Now()
			return chunk
		}
		return nil
	}
	w.loadingChunks[key] = true
	w.loadingMutex.Unlock()

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

	// Mark loading as complete
	w.loadingMutex.Lock()
	delete(w.loadingChunks, key)
	w.loadingMutex.Unlock()

	chunk.LastAccessed = time.Now()
	return chunk
}

// generateChunk generates terrain for a chunk with biome integration
func (w *World) generateChunk(chunk *Chunk) {
	worldX, worldY := chunk.GetWorldPosition()

	// Use cached noise generator for better performance
	noise := w.noiseGenerator

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
			case biomes.CORAL_REEF:
				baseHeight = 560.0 // Slightly lower for coral reefs
			case biomes.MANGROVE:
				baseHeight = 545.0 // Slightly above ocean
			case biomes.DESERT:
				baseHeight = 380.0 // Lower for deserts
			case biomes.SAVANNA:
				baseHeight = 390.0 // Slightly higher than plains
			case biomes.MOUNTAINS:
				baseHeight = 320.0 // Higher for mountains
			case biomes.VOLCANIC:
				baseHeight = 300.0 // Highest for volcanic regions
			case biomes.SWAMP:
				baseHeight = 420.0
			case biomes.TAIGA:
				baseHeight = 380.0 // Cooler climate
			case biomes.TUNDRA:
				baseHeight = 360.0 // Cold, flat terrain
			case biomes.JUNGLE:
				baseHeight = 400.0 // Tropical elevation
			case biomes.ICE_FIELDS:
				baseHeight = 340.0 // Ice cap elevation
			default:
				baseHeight = 400.0
			}

			// Enhanced multi-layer terrain noise for more realistic terrain
			// Continental scale features (mountains, valleys)
			continentalNoise := noise.Noise2D(x*0.001, y*0.001) * 200

			// Regional scale features (hills, ridges)
			regionalNoise := noise.Noise2D(x*0.003, y*0.003) * 100

			// Local scale features (small hills, dunes)
			localNoise := noise.Noise2D(x*0.01, y*0.01) * 40

			// Detail scale features (small variations)
			detailNoise := noise.Noise2D(x*0.05, y*0.05) * 8

			// River and valley cutting
			riverNoise := noise.Noise2D(x*0.015, y*0.015) * 30
			if riverNoise < -0.3 {
				riverNoise *= 2.0 // Deepen valleys
			}

			// Combine all noise layers with biome-specific weighting
			terrainNoise := continentalNoise + regionalNoise + localNoise + detailNoise + riverNoise

			// Combine all noise layers
			surfaceY := baseHeight + terrainNoise

			// Determine block type based on depth and biome
			depth := y - surfaceY

			var blockType blocks.BlockType

			if depth < -10 {
				// Above surface - air (unless it's an ocean)
				if biomeType == biomes.OCEAN && depth < 0 && depth > -60 {
					blockType = blocks.AIR // Ocean areas are air for now
				} else {
					blockType = blocks.AIR
				}
			} else if depth < 5 {
				// Surface layer - determine by biome
				switch biomeType {
				case biomes.DESERT:
					blockType = blocks.SAND
				case biomes.OCEAN, biomes.CORAL_REEF:
					blockType = blocks.SAND
				case biomes.MANGROVE:
					blockType = blocks.GRASS
				case biomes.SWAMP:
					blockType = blocks.GRASS
				case biomes.MOUNTAINS, biomes.VOLCANIC:
					blockType = blocks.STONE
				case biomes.TAIGA:
					blockType = blocks.GRASS
				case biomes.TUNDRA:
					blockType = blocks.SNOW
				case biomes.JUNGLE:
					blockType = blocks.GRASS
				case biomes.SAVANNA:
					blockType = blocks.GRASS
				case biomes.ICE_FIELDS:
					blockType = blocks.ICE
				default:
					blockType = blocks.GRASS
				}
			} else if depth < 15 {
				// Subsurface layer
				switch biomeType {
				case biomes.DESERT, biomes.OCEAN, biomes.CORAL_REEF:
					blockType = blocks.SAND
				case biomes.SWAMP, biomes.MANGROVE, biomes.JUNGLE, biomes.SAVANNA:
					blockType = blocks.DIRT
				case biomes.TAIGA:
					blockType = blocks.DIRT
				case biomes.TUNDRA:
					blockType = blocks.SNOW
				case biomes.ICE_FIELDS:
					blockType = blocks.ICE
				case biomes.MOUNTAINS, biomes.VOLCANIC:
					blockType = blocks.STONE
				default:
					blockType = blocks.DIRT
				}
			} else if depth < 200 {
				// Stone layers with ore generation
				// Use biome ore frequency modifier
				oreFrequency := biomeProps.OreFrequency

				// Enhanced ore generation with more variety
				// Use position-based seed for consistent ore generation
				posSeed := w.Seed + int64(x)*1000 + int64(y)
				rand.Seed(posSeed)
				oreChance := rand.Float64()

				// Add vein detection for more realistic ore deposits
				veinNoise := noise.Noise2D(x*0.1, y*0.1)
				isVein := veinNoise > 0.7

				// Adjust ore chance by biome frequency and vein detection
				if isVein {
					// Higher chance in veins
					oreChance *= 3.0
				}

				// Multiple ore types with realistic depth distribution
				if depth > 15 && depth < 60 && oreChance < 0.03*oreFrequency {
					blockType = blocks.COAL_ORE
				} else if depth > 25 && depth < 80 && oreChance < 0.025*oreFrequency {
					blockType = blocks.IRON_ORE
				} else if depth > 35 && depth < 65 && oreChance < 0.015*oreFrequency {
					blockType = blocks.GOLD_ORE
				} else if depth > 45 && depth < 90 && oreChance < 0.008*oreFrequency {
					blockType = blocks.DIAMOND_ORE
				} else if depth > 20 && depth < 55 && oreChance < 0.02*oreFrequency && biomeType == biomes.VOLCANIC {
					blockType = blocks.OBSIDIAN // Volcanic regions have obsidian
				} else if depth > 10 && depth < 40 && oreChance < 0.012*oreFrequency {
					blockType = blocks.GRAVEL // Common stone deposit
				} else if depth > 5 && depth < 30 && oreChance < 0.01*oreFrequency {
					blockType = blocks.SANDSTONE // Sedimentary deposits
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

			// Enhanced organism spawning for all biomes
			if (blockType == blocks.GRASS || blockType == blocks.SAND || blockType == blocks.SNOW || blockType == blocks.ICE) && depth >= -2 && depth <= 2 {
				// Use position-based seed for consistent organism spawning
				posSeed := w.Seed + int64(x)*10000 + int64(y)*10000
				rand.Seed(posSeed)
				spawnChance := rand.Float64()

				var orgType string
				var spawnProbability float64

				switch biomeType {
				case biomes.FOREST:
					if spawnChance < 0.15 {
						orgType = "tree"
						spawnProbability = 0.15
					} else if spawnChance < 0.25 {
						orgType = "bush"
						spawnProbability = 0.10
					} else if spawnChance < 0.35 {
						orgType = "flower"
						spawnProbability = 0.10
					}
				case biomes.JUNGLE:
					if spawnChance < 0.25 {
						orgType = "tree"
						spawnProbability = 0.25
					} else if spawnChance < 0.35 {
						orgType = "bush"
						spawnProbability = 0.10
					}
				case biomes.TAIGA:
					if spawnChance < 0.12 {
						orgType = "tree"
						spawnProbability = 0.12
					} else if spawnChance < 0.18 {
						orgType = "bush"
						spawnProbability = 0.06
					}
				case biomes.DESERT:
					if spawnChance < 0.02 {
						orgType = "cactus"
						spawnProbability = 0.02
					} else if spawnChance < 0.05 {
						orgType = "dead_bush"
						spawnProbability = 0.03
					}
				case biomes.SAVANNA:
					if spawnChance < 0.08 {
						orgType = "tree"
						spawnProbability = 0.08
					} else if spawnChance < 0.12 {
						orgType = "bush"
						spawnProbability = 0.04
					}
				case biomes.TUNDRA:
					if spawnChance < 0.01 {
						orgType = "ice_shrub"
						spawnProbability = 0.01
					}
				case biomes.ICE_FIELDS:
					if spawnChance < 0.005 {
						orgType = "ice_spike"
						spawnProbability = 0.005
					}
				case biomes.SWAMP:
					if spawnChance < 0.08 {
						orgType = "bush"
						spawnProbability = 0.08
					} else if spawnChance < 0.12 {
						orgType = "flower"
						spawnProbability = 0.04
					}
				case biomes.MANGROVE:
					if spawnChance < 0.10 {
						orgType = "mangrove_tree"
						spawnProbability = 0.10
					}
				case biomes.CORAL_REEF:
					if spawnChance < 0.05 {
						orgType = "coral"
						spawnProbability = 0.05
					}
				case biomes.MOUNTAINS:
					if spawnChance < 0.03 {
						orgType = "bush"
						spawnProbability = 0.03
					}
				case biomes.VOLCANIC:
					if spawnChance < 0.02 {
						orgType = "lava_rock"
						spawnProbability = 0.02
					}
				default: // PLAINS
					if spawnChance < 0.05 {
						orgType = "tree"
						spawnProbability = 0.05
					} else if spawnChance < 0.10 {
						orgType = "bush"
						spawnProbability = 0.05
					} else if spawnChance < 0.15 {
						orgType = "flower"
						spawnProbability = 0.05
					}
				}

				if spawnChance < spawnProbability {
					// Convert pixel coordinates to hex coordinates
					q, r := hexagon.PixelToHex(x, y, HexSize)
					hexagonCoords := hexagon.HexRound(q, r)
					orgHex, _ := hexagon.AxialToHex(hexagonCoords.Q, hexagonCoords.R)

					organism := organisms.CreateOrganism(orgType, x, y, orgHex)
					if organism != nil {
						w.Organisms = append(w.Organisms, organism)
					}
				}
			}

			// Spawn trees in appropriate biomes (legacy comment, now handled above)
			if depth >= -5 && depth <= 5 && biomeType != biomes.OCEAN && biomeType != biomes.DESERT {
				if biomes.ShouldSpawnTree(x, y, noise) {
					// This is now handled by the organism spawning above
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

// GetHexagonAt returns the hexagon at the given world position with tolerance
func (w *World) GetHexagonAt(x, y float64) *Hexagon {
	chunkX, chunkY := w.GetChunkCoords(x, y)
	chunk := w.GetChunk(chunkX, chunkY)

	// First try exact position
	hex := chunk.GetHexagon(x, y)
	if hex != nil {
		return hex
	}

	// If no exact match, search nearby hexagons within tolerance
	tolerance := 30.0 // pixels - roughly half hexagon size
	nearbyHexagons := w.GetNearbyHexagons(x, y, tolerance)

	var closestHex *Hexagon
	minDistance := math.MaxFloat64

	for _, nearbyHex := range nearbyHexagons {
		dx := nearbyHex.X - x
		dy := nearbyHex.Y - y
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance < minDistance {
			minDistance = distance
			closestHex = nearbyHex
		}
	}

	return closestHex
}

// AddHexagonAt adds a hexagon at the given world position
func (w *World) AddHexagonAt(x, y float64, blockType blocks.BlockType) {
	// Use the coordinates directly - don't convert to center
	hexagon := NewHexagon(x, y, HexSize, blockType)

	chunkX, chunkY := w.GetChunkCoords(x, y)
	chunk := w.GetChunk(chunkX, chunkY)
	chunk.AddHexagon(x, y, hexagon)

	// Add to spatial hash
	w.addHexagonToSpatialHash(hexagon)
}

// RemoveHexagonAt removes the hexagon at the given world position
func (w *World) RemoveHexagonAt(x, y float64) bool {
	chunkX, chunkY := w.GetChunkCoords(x, y)
	chunk := w.GetChunk(chunkX, chunkY)

	// Get the hexagon using direct coordinates (same as AddHexagon)
	hexagon := chunk.GetHexagonDirect(x, y)
	if hexagon == nil {
		return false
	}

	// Remove from spatial hash first
	w.removeHexagonFromSpatialHash(hexagon)

	// Then remove from chunk using direct coordinates
	return chunk.RemoveHexagonDirect(x, y)
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

// GetVisibleBlocks returns all visible blocks based on camera position with frustum culling
func (w *World) GetVisibleBlocks(cameraX, cameraY float64) []*Hexagon {
	// Calculate visible area with some padding for smoother edge transitions
	visibleWidth := 1280.0 + 200.0 // ScreenWidth + padding
	visibleHeight := 720.0 + 200.0 // ScreenHeight + padding

	// Convert to world coordinates
	halfWidth := visibleWidth / 2.0
	halfHeight := visibleHeight / 2.0

	left := cameraX - halfWidth
	right := cameraX + halfWidth
	top := cameraY - halfHeight
	bottom := cameraY + halfHeight

	// Get blocks in the visible area
	allBlocks := w.GetNearbyHexagons(cameraX, cameraY, float64(RenderDistance)*GetChunkWidth())

	// Filter blocks that are actually visible (frustum culling)
	var visibleBlocks []*Hexagon
	for _, hex := range allBlocks {
		if hex.X >= left && hex.X <= right && hex.Y >= top && hex.Y <= bottom {
			visibleBlocks = append(visibleBlocks, hex)
		}
	}

	return visibleBlocks
}

// GetVisibleBlocksForLayer returns visible blocks for a specific layer
func (w *World) GetVisibleBlocksForLayer(cameraX, cameraY float64, layer int) []*Hexagon {
	// Calculate visible area with some padding for smoother edge transitions
	visibleWidth := 1280.0 + 200.0 // ScreenWidth + padding
	visibleHeight := 720.0 + 200.0 // ScreenHeight + padding

	// Convert to world coordinates
	halfWidth := visibleWidth / 2.0
	halfHeight := visibleHeight / 2.0

	left := cameraX - halfWidth
	right := cameraX + halfWidth
	top := cameraY - halfHeight
	bottom := cameraY + halfHeight

	// Get blocks in the visible area
	allBlocks := w.GetNearbyHexagons(cameraX, cameraY, float64(RenderDistance)*GetChunkWidth())

	// Filter blocks by layer and visibility
	var visibleBlocks []*Hexagon
	for _, hex := range allBlocks {
		if hex.X >= left && hex.X <= right && hex.Y >= top && hex.Y <= bottom && hex.Z == layer {
			visibleBlocks = append(visibleBlocks, hex)
		}
	}

	return visibleBlocks
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

// SpawnCreatures spawns creatures in the world based on time of day and location
func (w *World) SpawnCreatures(dayNightCycle *gametime.DayNightCycle, playerX, playerY float64) {
	// Only spawn at night or dusk for hostile creatures
	timeOfDay := dayNightCycle.GetCurrentTimeOfDay()
	isNightTime := timeOfDay == gametime.Dusk || timeOfDay == gametime.Night || timeOfDay == gametime.Midnight

	// Limit total creatures to prevent overcrowding
	maxCreatures := 20
	if len(w.Creatures) >= maxCreatures {
		return
	}

	// Spawn area around player
	spawnRadius := 200.0
	spawnAttempts := 5

	for i := 0; i < spawnAttempts && len(w.Creatures) < maxCreatures; i++ {
		// Random position around player
		angle := rand.Float64() * 2 * math.Pi
		distance := rand.Float64() * spawnRadius
		spawnX := playerX + math.Cos(angle)*distance
		spawnY := playerY + math.Sin(angle)*distance

		// Check if position is valid (not in solid blocks)
		hex := w.GetHexagonAt(spawnX, spawnY)
		if hex == nil || hex.BlockType != blocks.AIR {
			continue // Can't spawn in solid blocks
		}

		// Determine creature type based on biome and time
		var creatureType creatures.CreatureType
		if isNightTime {
			// Night spawns
			switch rand.Intn(3) {
			case 0:
				creatureType = creatures.SLIME
			case 1:
				creatureType = creatures.SPIDER
			case 2:
				creatureType = creatures.ZOMBIE
			}
		} else {
			// Day spawns - only passive or less aggressive creatures
			// For now, only spawn at night
			continue
		}

		// Create and add creature
		// Convert pixel coordinates to hex coordinates
		q, r := hexagon.PixelToHex(hex.X, hex.Y, HexSize)
		hexagonCoords := hexagon.HexRound(q, r)
		creatureHex, _ := hexagon.AxialToHex(hexagonCoords.Q, hexagonCoords.R)
		creature := creatures.NewCreature(creatureType, spawnX, spawnY, creatureHex)
		w.Creatures = append(w.Creatures, creature)
	}
}

// UpdateCreatures updates all creatures in the world
func (w *World) UpdateCreatures(playerX, playerY float64, deltaTime float64) {
	// Create isBlocked function for creatures
	isBlocked := func(x, y float64) bool {
		hex := w.GetHexagonAt(x, y)
		return hex != nil && hex.BlockType != blocks.AIR
	}

	for _, creature := range w.Creatures {
		creature.Update(playerX, playerY, deltaTime, isBlocked)
	}
}

// RemoveDeadCreatures removes creatures that have died
func (w *World) RemoveDeadCreatures() {
	var aliveCreatures []*creatures.Creature
	for _, creature := range w.Creatures {
		if creature.IsAlive() {
			aliveCreatures = append(aliveCreatures, creature)
		}
	}
	w.Creatures = aliveCreatures
}

// FindSpawnPosition finds a suitable spawn position near the given coordinates
func (w *World) FindSpawnPosition(centerX, centerY float64) (float64, float64) {
	// Search radius for finding suitable spawn
	searchRadius := 500.0
	stepSize := 50.0

	// Check in a spiral pattern around the center
	for radius := 0.0; radius <= searchRadius; radius += stepSize {
		if radius == 0 {
			// Check center position first
			if spawnX, spawnY := w.checkPositionForSpawn(centerX, centerY); spawnX != 0 || spawnY != 0 {
				return spawnX, spawnY
			}
		} else {
			// Check positions in a circle at this radius
			steps := int(2 * math.Pi * radius / stepSize)
			for i := 0; i < steps; i++ {
				angle := float64(i) / float64(steps) * 2 * math.Pi
				checkX := centerX + math.Cos(angle)*radius
				checkY := centerY + math.Sin(angle)*radius

				if spawnX, spawnY := w.checkPositionForSpawn(checkX, checkY); spawnX != 0 || spawnY != 0 {
					return spawnX, spawnY
				}
			}
		}
	}

	// Fallback: generate terrain at center position and return guaranteed safe spawn
	w.GetChunksInRange(centerX, centerY)

	// Find ground level at center position
	var groundY float64 = centerY + 400 // Start searching from typical terrain height
	for checkY := centerY + 300; checkY <= centerY+600; checkY += 10.0 {
		hex := w.GetHexagonAt(centerX, checkY)
		if hex != nil && hex.BlockType != blocks.AIR {
			groundY = checkY
			break
		}
	}

	return centerX, groundY - 50 // Spawn 50 pixels above ground
}

// checkPositionForSpawn checks if a position is suitable for spawning
func (w *World) checkPositionForSpawn(x, y float64) (float64, float64) {
	// Ensure chunks are loaded around this position
	w.GetChunksInRange(x, y)

	// Check for solid ground below this position
	// Look for a solid block within a reasonable distance below (based on world generation heights)
	maxGroundCheck := 600.0 // Up to ocean depth
	minGroundCheck := 300.0 // Start around typical terrain height

	for checkY := y + minGroundCheck; checkY <= y+maxGroundCheck; checkY += 10.0 {
		hex := w.GetHexagonAt(x, checkY)
		if hex != nil && hex.BlockType != blocks.AIR {
			// Found solid ground, spawn 50 pixels above it
			return x, checkY - 50
		}
	}

	// No suitable ground found
	return 0, 0
}

// GetCreaturesInArea returns creatures within a certain radius of a point
func (w *World) GetCreaturesInArea(centerX, centerY, radius float64) []*creatures.Creature {
	var nearbyCreatures []*creatures.Creature
	for _, creature := range w.Creatures {
		dx := creature.X - centerX
		dy := creature.Y - centerY
		distance := math.Sqrt(dx*dx + dy*dy)
		if distance <= radius {
			nearbyCreatures = append(nearbyCreatures, creature)
		}
	}
	return nearbyCreatures
}
