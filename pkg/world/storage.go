package world

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// WorldStorage handles persistent storage of world data
type WorldStorage struct {
	WorldDir string
}

// NewWorldStorage creates a new world storage instance
func NewWorldStorage(worldName string) *WorldStorage {
	worldDir := filepath.Join("worlds", worldName)

	// Create world directory if it doesn't exist
	if _, err := os.Stat(worldDir); os.IsNotExist(err) {
		os.MkdirAll(worldDir, 0755)
	}

	return &WorldStorage{
		WorldDir: worldDir,
	}
}

// SaveChunk saves a single chunk to disk
func (ws *WorldStorage) SaveChunk(chunk *Chunk) error {
	if !chunk.Modified {
		return nil // Skip saving unchanged chunks
	}

	chunkData := chunk.ToChunkData()

	data, err := json.MarshalIndent(chunkData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal chunk data: %w", err)
	}

	filename := filepath.Join(ws.WorldDir, fmt.Sprintf("chunk_%d_%d.json", chunk.ChunkX, chunk.ChunkY))

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write chunk file: %w", err)
	}

	chunk.Modified = false
	chunk.LastSaved = time.Now()

	return nil
}

// LoadChunk loads a chunk from disk
func (ws *WorldStorage) LoadChunk(chunkX, chunkY int) (*Chunk, error) {
	filename := filepath.Join(ws.WorldDir, fmt.Sprintf("chunk_%d_%d.json", chunkX, chunkY))

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Chunk doesn't exist, return nil
		}
		return nil, fmt.Errorf("failed to read chunk file: %w", err)
	}

	var chunkData ChunkData
	err = json.Unmarshal(data, &chunkData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chunk data: %w", err)
	}

	chunk := NewChunk(chunkX, chunkY)
	chunk.FromChunkData(&chunkData)

	return chunk, nil
}

// SaveWorld saves all modified chunks in the world
func (ws *WorldStorage) SaveWorld(world *World) error {
	var saveErrors []error

	for key, chunk := range world.Chunks {
		if chunk.Modified {
			err := ws.SaveChunk(chunk)
			if err != nil {
				saveErrors = append(saveErrors, fmt.Errorf("chunk %d,%d: %w", key[0], key[1], err))
			}
		}
	}

	if len(saveErrors) > 0 {
		return fmt.Errorf("failed to save some chunks: %v", saveErrors)
	}

	return nil
}

// LoadWorld loads chunks around a specific position
func (ws *WorldStorage) LoadWorld(world *World, centerX, centerY float64, radius int) error {
	centerChunkX, centerChunkY := world.GetChunkCoords(centerX, centerY)

	for dx := -radius; dx <= radius; dx++ {
		for dy := -radius; dy <= radius; dy++ {
			chunkX := centerChunkX + dx
			chunkY := centerChunkY + dy

			key := [2]int{chunkX, chunkY}
			if _, exists := world.Chunks[key]; !exists {
				chunk, err := ws.LoadChunk(chunkX, chunkY)
				if err != nil {
					return fmt.Errorf("failed to load chunk %d,%d: %w", chunkX, chunkY, err)
				}

				if chunk != nil {
					world.Chunks[key] = chunk
				}
			}
		}
	}

	return nil
}

// GetWorldMetadata returns metadata about the saved world
func (ws *WorldStorage) GetWorldMetadata() (*WorldMetadata, error) {
	metadataFile := filepath.Join(ws.WorldDir, "metadata.json")

	data, err := os.ReadFile(metadataFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default metadata if file doesn't exist
			return &WorldMetadata{
				CreatedAt:  time.Now(),
				LastSaved:  time.Now(),
				ChunkCount: 0,
				Version:    "1.0",
			}, nil
		}
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var metadata WorldMetadata
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &metadata, nil
}

// SaveWorldMetadata saves world metadata
func (ws *WorldStorage) SaveWorldMetadata(chunkCount int) error {
	metadata := WorldMetadata{
		CreatedAt:  time.Now(),
		LastSaved:  time.Now(),
		ChunkCount: chunkCount,
		Version:    "1.0",
	}

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	metadataFile := filepath.Join(ws.WorldDir, "metadata.json")
	err = os.WriteFile(metadataFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// WorldMetadata contains metadata about a saved world
type WorldMetadata struct {
	CreatedAt  time.Time `json:"created_at"`
	LastSaved  time.Time `json:"last_saved"`
	ChunkCount int       `json:"chunk_count"`
	Version    string    `json:"version"`
}

// ListSavedWorlds returns a list of all saved world names
func ListSavedWorlds() ([]string, error) {
	worldsDir := "worlds"

	if _, err := os.Stat(worldsDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(worldsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read worlds directory: %w", err)
	}

	var worlds []string
	for _, entry := range entries {
		if entry.IsDir() {
			worlds = append(worlds, entry.Name())
		}
	}

	return worlds, nil
}

// DeleteWorld removes a saved world from disk
func DeleteWorld(worldName string) error {
	worldDir := filepath.Join("worlds", worldName)

	if _, err := os.Stat(worldDir); os.IsNotExist(err) {
		return nil // World doesn't exist, nothing to delete
	}

	err := os.RemoveAll(worldDir)
	if err != nil {
		return fmt.Errorf("failed to delete world directory: %w", err)
	}

	return nil
}
