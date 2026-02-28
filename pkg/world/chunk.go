package world

import (
	"fmt"
	"time"

	"tesselbox/pkg/blocks"
)

const (
	// ChunkSize is the number of hexagons per chunk dimension
	ChunkSize = 32
)

// GetChunkWidth returns the width of a chunk in world coordinates
func GetChunkWidth() float64 {
	return float64(ChunkSize) * HexWidth
}

// GetChunkHeight returns the height of a chunk in world coordinates
func GetChunkHeight() float64 {
	return float64(ChunkSize) * HexVSpacing
}

// Chunk represents a chunk of the world containing multiple hexagons
type Chunk struct {
	ChunkX       int
	ChunkY       int
	Hexagons     map[[2]int]*Hexagon
	Modified     bool
	LastAccessed time.Time
	LastSaved    time.Time
}

// NewChunk creates a new chunk
func NewChunk(chunkX, chunkY int) *Chunk {
	return &Chunk{
		ChunkX:       chunkX,
		ChunkY:       chunkY,
		Hexagons:     make(map[[2]int]*Hexagon),
		Modified:     false,
		LastAccessed: time.Now(),
		LastSaved:    time.Time{},
	}
}

// GetWorldPosition returns the world position of the chunk
func (c *Chunk) GetWorldPosition() (float64, float64) {
	worldX := float64(c.ChunkX) * GetChunkWidth()
	worldY := float64(c.ChunkY) * GetChunkHeight()
	return worldX, worldY
}

// GetHexagon returns the hexagon at the given world coordinates
func (c *Chunk) GetHexagon(x, y float64) *Hexagon {
	worldX, worldY := c.GetWorldPosition()

	// Calculate local row
	localRow := int((y - worldY) / HexVSpacing)

	// Calculate local column with proper offset for interlocking
	localCol := int((x - worldX - HexWidth/2) / HexWidth)
	if localRow%2 == 0 {
		localCol = int((x - worldX - HexWidth/2) / HexWidth)
	} else {
		localCol = int((x - worldX) / HexWidth)
	}

	key := [2]int{localCol, localRow}
	return c.Hexagons[key]
}

// AddHexagon adds a hexagon to the chunk
func (c *Chunk) AddHexagon(x, y float64, hexagon *Hexagon) {
	worldX, worldY := c.GetWorldPosition()

	// Calculate local row
	localRow := int((y - worldY) / HexVSpacing)

	// Calculate local column with proper offset for interlocking
	localCol := int((x - worldX - HexWidth/2) / HexWidth)
	if localRow%2 == 0 {
		localCol = int((x - worldX - HexWidth/2) / HexWidth)
	} else {
		localCol = int((x - worldX) / HexWidth)
	}

	hexagon.ChunkX = c.ChunkX
	hexagon.ChunkY = c.ChunkY
	key := [2]int{localCol, localRow}
	c.Hexagons[key] = hexagon
	c.Modified = true
}

// RemoveHexagon removes a hexagon from the chunk
func (c *Chunk) RemoveHexagon(x, y float64) bool {
	worldX, worldY := c.GetWorldPosition()

	// Calculate local row
	localRow := int((y - worldY) / HexVSpacing)

	// Calculate local column with proper offset for interlocking
	localCol := int((x - worldX - HexWidth/2) / HexWidth)
	if localRow%2 == 0 {
		localCol = int((x - worldX - HexWidth/2) / HexWidth)
	} else {
		localCol = int((x - worldX) / HexWidth)
	}

	key := [2]int{localCol, localRow}
	if _, ok := c.Hexagons[key]; ok {
		delete(c.Hexagons, key)
		c.Modified = true
		return true
	}
	return false
}

// ChunkData represents the serializable data for a chunk
type ChunkData struct {
	ChunkX   int                           `json:"chunk_x"`
	ChunkY   int                           `json:"chunk_y"`
	Hexagons map[string]*SerializedHexagon `json:"hexagons"`
}

// SerializedHexagon represents a hexagon that can be serialized to JSON
type SerializedHexagon struct {
	X         float64          `json:"x"`
	Y         float64          `json:"y"`
	Size      float64          `json:"size"`
	BlockType blocks.BlockType `json:"block_type"`
	Health    float64          `json:"health"`
}

// ToChunkData converts a chunk to serializable format
func (c *Chunk) ToChunkData() *ChunkData {
	hexagons := make(map[string]*SerializedHexagon)
	for key, hex := range c.Hexagons {
		keyStr := fmt.Sprintf("%d,%d", key[0], key[1])
		hexagons[keyStr] = &SerializedHexagon{
			X:         hex.X,
			Y:         hex.Y,
			Size:      hex.Size,
			BlockType: hex.BlockType,
			Health:    hex.Health,
		}
	}

	return &ChunkData{
		ChunkX:   c.ChunkX,
		ChunkY:   c.ChunkY,
		Hexagons: hexagons,
	}
}

// FromChunkData loads chunk data from serializable format
func (c *Chunk) FromChunkData(data *ChunkData) {
	c.ChunkX = data.ChunkX
	c.ChunkY = data.ChunkY
	c.Hexagons = make(map[[2]int]*Hexagon)

	for keyStr, serHex := range data.Hexagons {
		var key [2]int
		fmt.Sscanf(keyStr, "%d,%d", &key[0], &key[1])

		hex := &Hexagon{
			X:         serHex.X,
			Y:         serHex.Y,
			Size:      serHex.Size,
			BlockType: serHex.BlockType,
			Health:    serHex.Health,
			ChunkX:    c.ChunkX,
			ChunkY:    c.ChunkY,
		}
		c.Hexagons[key] = hex
	}

	c.Modified = false
	c.LastSaved = time.Now()
}
