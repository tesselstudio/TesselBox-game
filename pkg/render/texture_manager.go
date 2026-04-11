package render

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// TextureManager provides centralized texture caching and resource management
type TextureManager struct {
	cache       map[string]*ebiten.Image
	whiteImage  *ebiten.Image
	atlas       *ebiten.Image
	atlasSlots  map[string]image.Rectangle
	atlasSize   int
	atlasCursor image.Point
	mu          sync.RWMutex
	maxCacheSize int
}

var (
	instance *TextureManager
	once     sync.Once
)

// GetTextureManager returns the singleton texture manager instance
func GetTextureManager() *TextureManager {
	once.Do(func() {
		instance = &TextureManager{
			cache:        make(map[string]*ebiten.Image),
			atlasSlots:   make(map[string]image.Rectangle),
			atlasSize:    4096,
			maxCacheSize: 1000,
		}
		instance.whiteImage = ebiten.NewImage(1, 1)
		instance.whiteImage.Fill(color.RGBA{255, 255, 255, 255})
		instance.atlas = ebiten.NewImage(instance.atlasSize, instance.atlasSize)
	})
	return instance
}

// GetWhiteImage returns the shared 1x1 white image for solid color drawing
func (tm *TextureManager) GetWhiteImage() *ebiten.Image {
	return tm.whiteImage
}

// GetCachedTexture retrieves a cached texture by key
func (tm *TextureManager) GetCachedTexture(key string) (*ebiten.Image, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	img, ok := tm.cache[key]
	return img, ok
}

// CacheTexture stores a texture in the cache
func (tm *TextureManager) CacheTexture(key string, img *ebiten.Image) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	if len(tm.cache) >= tm.maxCacheSize {
		tm.evictOldest()
	}
	
	tm.cache[key] = img
}

// GenerateBlockTextureKey creates a cache key for block textures
func (tm *TextureManager) GenerateBlockTextureKey(blockType string, variation int) string {
	return fmt.Sprintf("block:%s:%d", blockType, variation)
}

// GetOrCreateTexture gets a cached texture or creates it using the provided function
func (tm *TextureManager) GetOrCreateTexture(key string, create func() *ebiten.Image) *ebiten.Image {
	if img, ok := tm.GetCachedTexture(key); ok {
		return img
	}
	
	img := create()
	tm.CacheTexture(key, img)
	return img
}

// GetAtlasSlot returns the atlas position for a texture, adding it if needed
func (tm *TextureManager) GetAtlasSlot(key string, img *ebiten.Image) (image.Rectangle, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	if slot, ok := tm.atlasSlots[key]; ok {
		return slot, true
	}
	
	// Get image bounds
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	
	// Check if we can fit in atlas
	if tm.atlasCursor.X+w > tm.atlasSize {
		tm.atlasCursor.X = 0
		tm.atlasCursor.Y += 64 // Row height
	}
	
	if tm.atlasCursor.Y+h > tm.atlasSize {
		return image.Rectangle{}, false // Atlas full
	}
	
	// Create slot
	slot := image.Rect(tm.atlasCursor.X, tm.atlasCursor.Y, tm.atlasCursor.X+w, tm.atlasCursor.Y+h)
	tm.atlasSlots[key] = slot
	
	// Copy image to atlas
	tm.atlas.DrawImage(img, &ebiten.DrawImageOptions{})
	
	// Advance cursor
	tm.atlasCursor.X += w
	
	return slot, true
}

// evictOldest removes oldest cached entries when cache is full
func (tm *TextureManager) evictOldest() {
	// Simple eviction: remove first 10% of entries
	toRemove := tm.maxCacheSize / 10
	count := 0
	for key := range tm.cache {
		if count >= toRemove {
			break
		}
		delete(tm.cache, key)
		count++
	}
}

// ClearCache clears all cached textures (call when switching worlds)
func (tm *TextureManager) ClearCache() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	for _, img := range tm.cache {
		img.Dispose()
	}
	
	tm.cache = make(map[string]*ebiten.Image)
	tm.atlasSlots = make(map[string]image.Rectangle)
	tm.atlasCursor = image.Point{}
	tm.atlas.Clear()
}

// GetStats returns cache statistics
func (tm *TextureManager) GetStats() (cacheSize, atlasSlots int) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return len(tm.cache), len(tm.atlasSlots)
}
