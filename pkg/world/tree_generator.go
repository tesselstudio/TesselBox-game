// Package world implements tree generation for natural world structures
package world

import (
	"math"
	"math/rand"

	"tesselbox/pkg/biomes"
	"tesselbox/pkg/blocks"
)

// TreeType represents different tree variants
type TreeType int

const (
	Oak TreeType = iota
	Birch
	Spruce
	Jungle
	Acacia
)

// TreeStructure defines a tree's properties
type TreeStructure struct {
	Type       TreeType
	Height     int
	TrunkWidth int
	LeafRadius int
	LeafHeight int
}

// TreeGenerator handles tree generation
type TreeGenerator struct {
	world       *World
	rng         *rand.Rand
	TreeDensity float64 // Trees per chunk (0-1)
}

// NewTreeGenerator creates a new tree generator
func NewTreeGenerator(world *World) *TreeGenerator {
	return &TreeGenerator{
		world:       world,
		rng:         rand.New(rand.NewSource(world.Seed)),
		TreeDensity: 0.15, // 15% of chunks have trees
	}
}

// GenerateTreesForChunk generates trees in a chunk based on biome
func (tg *TreeGenerator) GenerateTreesForChunk(chunk *Chunk) {
	// Determine biome
	biome := tg.getBiomeAtChunk(chunk.ChunkX, chunk.ChunkY)

	// Check if trees should generate in this chunk
	if tg.rng.Float64() > tg.TreeDensity {
		return
	}

	// Number of trees in this chunk
	numTrees := tg.rng.Intn(3) + 1 // 1-3 trees per chunk

	for i := 0; i < numTrees; i++ {
		tg.generateTree(chunk, biome)
	}
}

// getBiomeAtChunk determines the biome at a chunk position
func (tg *TreeGenerator) getBiomeAtChunk(chunkX, chunkY int) biomes.BiomeType {
	// Simplified biome determination based on world coordinates
	// In a real implementation, this would use noise functions
	_ = float64(chunkX) * float64(ChunkSize) * HexWidth
	worldY := float64(chunkY) * float64(ChunkSize) * HexVSpacing

	if worldY < -500 {
		return biomes.DESERT
	} else if worldY > 500 {
		return biomes.TUNDRA
	} else if worldY > 300 {
		return biomes.TAIGA
	} else {
		return biomes.FOREST
	}
}

// generateTree creates a single tree in a chunk
func (tg *TreeGenerator) generateTree(chunk *Chunk, biome biomes.BiomeType) {
	// Get tree type based on biome
	treeType := tg.selectTreeType(biome)
	structure := tg.getTreeStructure(treeType)

	// Random position in chunk
	localX := tg.rng.Intn(ChunkSize-4) + 2
	localY := tg.rng.Intn(ChunkSize-4) + 2

	// Find ground level
	groundY := tg.findGroundLevel(chunk, localX, localY)
	if groundY < 0 {
		return // No valid ground found
	}

	// Place trunk
	tg.placeTrunk(chunk, localX, groundY, structure)

	// Place leaves
	tg.placeLeaves(chunk, localX, groundY, structure)
}

// selectTreeType chooses a tree type based on biome
func (tg *TreeGenerator) selectTreeType(biome biomes.BiomeType) TreeType {
	switch biome {
	case biomes.FOREST:
		if tg.rng.Float64() < 0.7 {
			return Oak
		}
		return Birch
	case biomes.TAIGA:
		return Spruce
	case biomes.JUNGLE:
		return Jungle
	case biomes.SAVANNA:
		return Acacia
	case biomes.DESERT:
		return Oak // Small oak in desert oases
	default:
		return Oak
	}
}

// getTreeStructure returns the dimensions for a tree type
func (tg *TreeGenerator) getTreeStructure(t TreeType) TreeStructure {
	switch t {
	case Oak:
		return TreeStructure{
			Type:       Oak,
			Height:     tg.rng.Intn(4) + 4, // 4-7 blocks
			TrunkWidth: 1,
			LeafRadius: 2,
			LeafHeight: 3,
		}
	case Birch:
		return TreeStructure{
			Type:       Birch,
			Height:     tg.rng.Intn(3) + 5, // 5-7 blocks
			TrunkWidth: 1,
			LeafRadius: 2,
			LeafHeight: 2,
		}
	case Spruce:
		return TreeStructure{
			Type:       Spruce,
			Height:     tg.rng.Intn(5) + 7, // 7-11 blocks
			TrunkWidth: 1,
			LeafRadius: 1,
			LeafHeight: 4,
		}
	case Jungle:
		return TreeStructure{
			Type:       Jungle,
			Height:     tg.rng.Intn(6) + 8, // 8-13 blocks
			TrunkWidth: 2,
			LeafRadius: 3,
			LeafHeight: 4,
		}
	case Acacia:
		return TreeStructure{
			Type:       Acacia,
			Height:     tg.rng.Intn(3) + 4, // 4-6 blocks
			TrunkWidth: 1,
			LeafRadius: 2,
			LeafHeight: 2,
		}
	default:
		return TreeStructure{Type: Oak, Height: 5, TrunkWidth: 1, LeafRadius: 2, LeafHeight: 3}
	}
}

// findGroundLevel finds the y-coordinate of the ground at a given x position
func (tg *TreeGenerator) findGroundLevel(chunk *Chunk, x, startY int) int {
	// Scan downward to find solid ground
	for y := startY; y < ChunkSize; y++ {
		if hex := chunk.GetHexagon(float64(x), float64(y)); hex != nil {
			if hex.BlockType == blocks.GRASS || hex.BlockType == blocks.DIRT {
				return y - 1 // Return position above ground
			}
		}
	}
	return -1
}

// placeTrunk places the tree trunk blocks
func (tg *TreeGenerator) placeTrunk(chunk *Chunk, x, y int, structure TreeStructure) {
	for i := 0; i < structure.Height; i++ {
		targetY := y - i
		if targetY < 0 {
			continue
		}

		hex := chunk.GetHexagon(float64(x), float64(targetY))
		if hex != nil && hex.BlockType == blocks.AIR {
			hex.BlockType = blocks.LOG
		}

		// Wider trunk for jungle trees
		if structure.TrunkWidth > 1 {
			for dx := 0; dx < structure.TrunkWidth; dx++ {
				for dz := 0; dz < structure.TrunkWidth; dz++ {
					if dx == 0 && dz == 0 {
						continue
					}
					hex2 := chunk.GetHexagon(float64(x+dx), float64(targetY+dz))
					if hex2 != nil && hex2.BlockType == blocks.AIR {
						hex2.BlockType = blocks.LOG
					}
				}
			}
		}
	}
}

// placeLeaves places the leaf blocks around the tree
func (tg *TreeGenerator) placeLeaves(chunk *Chunk, x, y int, structure TreeStructure) {
	// Start placing leaves from top of trunk downward
	topY := y - structure.Height + 1

	for ly := 0; ly < structure.LeafHeight; ly++ {
		currentY := topY + ly
		if currentY < 0 {
			continue
		}

		// Calculate radius at this height (tapering for spruce)
		radius := structure.LeafRadius
		if structure.Type == Spruce {
			// Spruce leaves taper more
			radius = int(math.Max(1, float64(structure.LeafRadius-ly/2)))
		}

		// Place leaves in a circle around the trunk
		for dx := -radius; dx <= radius; dx++ {
			for dz := -radius; dz <= radius; dz++ {
				// Skip trunk position
				if dx == 0 && dz == 0 && ly < structure.LeafHeight-1 {
					continue
				}

				// Circular shape
				dist := math.Sqrt(float64(dx*dx + dz*dz))
				if dist > float64(radius)+0.5 {
					continue
				}

				targetX := x + dx
				targetZ := currentY + dz

				hex := chunk.GetHexagon(float64(targetX), float64(targetZ))
				if hex != nil && hex.BlockType == blocks.AIR {
					hex.BlockType = blocks.LEAVES
				}
			}
		}
	}
}

// SetTreeDensity sets the tree generation density (0-1)
func (tg *TreeGenerator) SetTreeDensity(density float64) {
	if density < 0 {
		density = 0
	}
	if density > 1 {
		density = 1
	}
	tg.TreeDensity = density
}
