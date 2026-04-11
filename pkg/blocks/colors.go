package blocks

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// SimpleLRUCache is a thread-safe LRU cache for ebiten images
type SimpleLRUCache struct {
	data    map[string]*ebiten.Image
	keys    []string
	size    int
	maxSize int
}

// NewSimpleLRUCache creates a new LRU cache with max size
func NewSimpleLRUCache(maxSize int) *SimpleLRUCache {
	return &SimpleLRUCache{
		data:    make(map[string]*ebiten.Image),
		keys:    make([]string, 0, maxSize),
		size:    0,
		maxSize: maxSize,
	}
}

// Get retrieves an item from the cache
func (c *SimpleLRUCache) Get(key string) (*ebiten.Image, bool) {
	if val, ok := c.data[key]; ok {
		// Move key to end (most recently used)
		c.moveToEnd(key)
		return val, true
	}
	return nil, false
}

// Set adds or updates an item in the cache
func (c *SimpleLRUCache) Set(key string, value *ebiten.Image) {
	if _, ok := c.data[key]; ok {
		c.data[key] = value
		c.moveToEnd(key)
		return
	}

	if c.size >= c.maxSize {
		// Evict oldest
		oldest := c.keys[0]
		delete(c.data, oldest)
		c.keys = c.keys[1:]
		c.size--
	}

	c.data[key] = value
	c.keys = append(c.keys, key)
	c.size++
}

// moveToEnd moves a key to the end of the keys slice
func (c *SimpleLRUCache) moveToEnd(key string) {
	for i, k := range c.keys {
		if k == key {
			c.keys = append(c.keys[:i], c.keys[i+1:]...)
			c.keys = append(c.keys, key)
			break
		}
	}
}

// Global texture cache with 1000 entry limit
var textureCache = NewSimpleLRUCache(1000)

// ColorVariation defines different ways blocks can vary in color
type ColorVariation int

const (
	VariationNone ColorVariation = iota
	VariationRandom
	VariationGradient
	VariationPattern
	VariationBiome
	VariationAge
	VariationMoisture
	VariationTemperature
)

// BlockColorScheme defines a color scheme for blocks with variations
type BlockColorScheme struct {
	BaseColors        []color.RGBA            `yaml:"baseColors"`
	VariationType     ColorVariation          `yaml:"variationType"`
	VariationRange    float64                 `yaml:"variationRange"`
	PatternColors     []color.RGBA            `yaml:"patternColors"`
	GradientColors    []color.RGBA            `yaml:"gradientColors"`
	BiomeColors       map[string][]color.RGBA `yaml:"biomeColors"`
	AgeColors         []color.RGBA            `yaml:"ageColors"`
	MoistureColors    []color.RGBA            `yaml:"moistureColors"`
	TemperatureColors []color.RGBA            `yaml:"temperatureColors"`
}

// BlockAppearance manages the visual appearance of blocks
type BlockAppearance struct {
	ColorSchemes map[string]*BlockColorScheme
	randSource   *rand.Rand
}

// NewBlockAppearance creates a new block appearance manager
func NewBlockAppearance() *BlockAppearance {
	return &BlockAppearance{
		ColorSchemes: make(map[string]*BlockColorScheme),
		randSource:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// InitializeColorSchemes sets up default color schemes for blocks
func (ba *BlockAppearance) InitializeColorSchemes() {
	// Grass with biome variations
	ba.ColorSchemes["grass"] = &BlockColorScheme{
		BaseColors: []color.RGBA{
			{124, 169, 84, 255}, // Base green
			{134, 179, 94, 255}, // Lighter green
			{114, 159, 74, 255}, // Darker green
		},
		VariationType:  VariationBiome,
		VariationRange: 0.3,
		BiomeColors: map[string][]color.RGBA{
			"forest": {
				{124, 169, 84, 255}, // Forest green
				{114, 159, 74, 255},
				{134, 179, 94, 255},
			},
			"plains": {
				{144, 189, 104, 255}, // Plains green
				{154, 199, 114, 255},
				{134, 179, 94, 255},
			},
			"desert": {
				{184, 199, 124, 255}, // Dry yellow-green
				{194, 209, 134, 255},
				{174, 189, 114, 255},
			},
			"taiga": {
				{104, 149, 64, 255}, // Dark taiga green
				{94, 139, 54, 255},
				{114, 159, 74, 255},
			},
			"tundra": {
				{164, 179, 134, 255}, // Tundra pale green
				{174, 189, 144, 255},
				{154, 169, 124, 255},
			},
		},
	}

	// Stone with mineral variations
	ba.ColorSchemes["stone"] = &BlockColorScheme{
		BaseColors: []color.RGBA{
			{128, 128, 128, 255}, // Gray stone
			{138, 138, 138, 255}, // Lighter gray
			{118, 118, 118, 255}, // Darker gray
		},
		VariationType:  VariationRandom,
		VariationRange: 0.2,
		PatternColors: []color.RGBA{
			{148, 148, 148, 255}, // Mineral flecks
			{108, 108, 108, 255}, // Dark mineral
			{168, 168, 168, 255}, // Light mineral
		},
	}

	// Sand with desert variations
	ba.ColorSchemes["sand"] = &BlockColorScheme{
		BaseColors: []color.RGBA{
			{238, 203, 173, 255}, // Base sand
			{248, 213, 183, 255}, // Lighter sand
			{228, 193, 163, 255}, // Darker sand
		},
		VariationType:  VariationTemperature,
		VariationRange: 0.25,
		TemperatureColors: []color.RGBA{
			{228, 193, 163, 255}, // Cool sand (more gray)
			{238, 203, 173, 255}, // Normal sand
			{248, 213, 183, 255}, // Warm sand (more yellow)
			{255, 223, 193, 255}, // Hot sand (very yellow)
		},
	}

	// Wood with tree type variations
	ba.ColorSchemes["log"] = &BlockColorScheme{
		BaseColors: []color.RGBA{
			{139, 69, 19, 255}, // Oak wood
			{149, 79, 29, 255}, // Lighter oak
			{129, 59, 9, 255},  // Darker oak
		},
		VariationType:  VariationPattern,
		VariationRange: 0.15,
		PatternColors: []color.RGBA{
			{119, 49, 9, 255},  // Dark grain
			{159, 89, 39, 255}, // Light grain
			{139, 69, 19, 255}, // Base wood
		},
	}

	// Wool with dye colors
	ba.ColorSchemes["wool"] = &BlockColorScheme{
		BaseColors: []color.RGBA{
			{222, 222, 222, 255}, // White wool
		},
		VariationType:  VariationRandom,
		VariationRange: 0.1,
		PatternColors: []color.RGBA{
			{255, 255, 255, 255}, // Pure white
			{242, 242, 242, 255}, // Off-white
			{232, 232, 232, 255}, // Light gray
			{212, 212, 212, 255}, // Medium gray
			{192, 192, 192, 255}, // Dark gray
			{255, 0, 0, 255},     // Red
			{0, 255, 0, 255},     // Green
			{0, 0, 255, 255},     // Blue
			{255, 255, 0, 255},   // Yellow
			{255, 0, 255, 255},   // Magenta
			{0, 255, 255, 255},   // Cyan
			{255, 165, 0, 255},   // Orange
			{128, 0, 128, 255},   // Purple
			{165, 42, 42, 255},   // Brown
			{255, 192, 203, 255}, // Pink
			{0, 0, 0, 255},       // Black
		},
	}

	// Leaves with seasonal variations
	ba.ColorSchemes["leaves"] = &BlockColorScheme{
		BaseColors: []color.RGBA{
			{34, 139, 34, 200}, // Spring green
			{44, 149, 44, 200}, // Lighter spring
			{24, 129, 24, 200}, // Darker spring
		},
		VariationType:  VariationAge,
		VariationRange: 0.4,
		AgeColors: []color.RGBA{
			{34, 139, 34, 200},  // Fresh leaves
			{44, 149, 44, 200},  // Mature leaves
			{54, 159, 54, 200},  // Full grown
			{64, 139, 24, 200},  // Early autumn
			{104, 139, 24, 200}, // Mid autumn
			{144, 89, 24, 200},  // Late autumn
			{134, 69, 19, 200},  // Dry leaves
		},
	}

	// Water with depth variations
	ba.ColorSchemes["water"] = &BlockColorScheme{
		BaseColors: []color.RGBA{
			{64, 164, 223, 180}, // Base water
			{74, 174, 233, 180}, // Lighter water
			{54, 154, 213, 180}, // Darker water
		},
		VariationType:  VariationGradient,
		VariationRange: 0.3,
		GradientColors: []color.RGBA{
			{44, 144, 203, 200}, // Deep water (darker, more opaque)
			{64, 164, 223, 180}, // Medium water
			{84, 184, 243, 160}, // Shallow water (lighter, more transparent)
		},
	}

	// Dirt with moisture variations
	ba.ColorSchemes["dirt"] = &BlockColorScheme{
		BaseColors: []color.RGBA{
			{139, 90, 43, 255},  // Base dirt
			{149, 100, 53, 255}, // Lighter dirt
			{129, 80, 33, 255},  // Darker dirt
		},
		VariationType:  VariationMoisture,
		VariationRange: 0.25,
		MoistureColors: []color.RGBA{
			{109, 60, 13, 255},  // Wet dirt (darker)
			{139, 90, 43, 255},  // Normal dirt
			{159, 110, 63, 255}, // Dry dirt (lighter, more brown)
			{179, 130, 83, 255}, // Very dry dirt (very light brown)
		},
	}

	// Ice with transparency variations
	ba.ColorSchemes["ice"] = &BlockColorScheme{
		BaseColors: []color.RGBA{
			{175, 223, 255, 200}, // Clear ice
			{185, 233, 255, 200}, // Lighter ice
			{165, 213, 255, 200}, // Darker ice
		},
		VariationType:  VariationRandom,
		VariationRange: 0.15,
		PatternColors: []color.RGBA{
			{195, 243, 255, 180}, // Clear patches
			{155, 203, 255, 220}, // Frozen patches
		},
	}

	// Gravel with stone mix variations
	ba.ColorSchemes["gravel"] = &BlockColorScheme{
		BaseColors: []color.RGBA{
			{136, 140, 141, 255}, // Base gravel
			{146, 150, 151, 255}, // Lighter gravel
			{126, 130, 131, 255}, // Darker gravel
		},
		VariationType:  VariationPattern,
		VariationRange: 0.3,
		PatternColors: []color.RGBA{
			{128, 128, 128, 255}, // Stone flecks
			{144, 148, 149, 255}, // Light gravel
			{128, 132, 133, 255}, // Medium gravel
		},
	}

	// Sandstone with layer variations
	ba.ColorSchemes["sandstone"] = &BlockColorScheme{
		BaseColors: []color.RGBA{
			{238, 203, 173, 255}, // Base sandstone
			{248, 213, 183, 255}, // Lighter sandstone
			{228, 193, 163, 255}, // Darker sandstone
		},
		VariationType:  VariationGradient,
		VariationRange: 0.2,
		GradientColors: []color.RGBA{
			{218, 183, 153, 255}, // Bottom layer (darker)
			{238, 203, 173, 255}, // Middle layer
			{255, 223, 193, 255}, // Top layer (lighter)
		},
	}
}

// GetBlockColor returns the color for a block with variations applied
func (ba *BlockAppearance) GetBlockColor(blockType string, x, y int, biome string, depth float64) color.RGBA {
	scheme, exists := ba.ColorSchemes[blockType]
	if !exists {
		// Fallback to basic block color
		if props, ok := BlockDefinitions[blockType]; ok {
			return props.Color
		}
		return color.RGBA{128, 128, 128, 255} // Default gray
	}

	baseColor := ba.selectBaseColor(scheme, x, y, biome, depth)
	variedColor := ba.applyVariation(scheme, baseColor, x, y, biome, depth)

	return variedColor
}

// selectBaseColor selects the appropriate base color based on variation type
func (ba *BlockAppearance) selectBaseColor(scheme *BlockColorScheme, x, y int, biome string, depth float64) color.RGBA {
	switch scheme.VariationType {
	case VariationBiome:
		if biomeColors, exists := scheme.BiomeColors[biome]; exists && len(biomeColors) > 0 {
			index := ba.randSource.Intn(len(biomeColors))
			return biomeColors[index]
		}
		fallthrough
	case VariationTemperature:
		if len(scheme.TemperatureColors) > 0 {
			// Temperature based on y-coordinate (higher = colder)
			tempIndex := int((float64(y) / 1000.0) * float64(len(scheme.TemperatureColors)))
			tempIndex = max(0, min(len(scheme.TemperatureColors)-1, tempIndex))
			return scheme.TemperatureColors[tempIndex]
		}
		fallthrough
	case VariationMoisture:
		if len(scheme.MoistureColors) > 0 {
			// Moisture based on x-coordinate (simplistic)
			moistureIndex := int((float64(x%1000) / 1000.0) * float64(len(scheme.MoistureColors)))
			return scheme.MoistureColors[moistureIndex]
		}
		fallthrough
	case VariationAge:
		if len(scheme.AgeColors) > 0 {
			// Age based on time and position
			ageIndex := int((float64((x+y)%100) / 100.0) * float64(len(scheme.AgeColors)))
			return scheme.AgeColors[ageIndex]
		}
		fallthrough
	case VariationGradient:
		if len(scheme.GradientColors) > 0 {
			// Gradient based on depth
			depthIndex := int(depth * float64(len(scheme.GradientColors)))
			depthIndex = max(0, min(len(scheme.GradientColors)-1, depthIndex))
			return scheme.GradientColors[depthIndex]
		}
		fallthrough
	default:
		if len(scheme.BaseColors) > 0 {
			index := ba.randSource.Intn(len(scheme.BaseColors))
			return scheme.BaseColors[index]
		}
	}

	// Fallback to first base color
	if len(scheme.BaseColors) > 0 {
		return scheme.BaseColors[0]
	}
	return color.RGBA{128, 128, 128, 255}
}

// applyVariation applies color variations to the base color
func (ba *BlockAppearance) applyVariation(scheme *BlockColorScheme, baseColor color.RGBA, x, y int, biome string, depth float64) color.RGBA {
	switch scheme.VariationType {
	case VariationRandom:
		return ba.applyRandomVariation(baseColor, scheme.VariationRange)
	case VariationPattern:
		return ba.applyPatternVariation(baseColor, scheme, x, y)
	case VariationGradient:
		return ba.applyGradientVariation(baseColor, scheme, depth)
	default:
		return baseColor
	}
}

// applyRandomVariation applies random color variation
func (ba *BlockAppearance) applyRandomVariation(baseColor color.RGBA, range_ float64) color.RGBA {
	variation := (ba.randSource.Float64() - 0.5) * 2.0 * range_

	newR := float64(baseColor.R) + variation*255
	newG := float64(baseColor.G) + variation*255
	newB := float64(baseColor.B) + variation*255

	// Clamp to valid range
	newR = math.Max(0, math.Min(255, newR))
	newG = math.Max(0, math.Min(255, newG))
	newB = math.Max(0, math.Min(255, newB))

	return color.RGBA{
		R: uint8(newR),
		G: uint8(newG),
		B: uint8(newB),
		A: baseColor.A,
	}
}

// applyPatternVariation applies pattern-based color variation
func (ba *BlockAppearance) applyPatternVariation(baseColor color.RGBA, scheme *BlockColorScheme, x, y int) color.RGBA {
	if len(scheme.PatternColors) == 0 {
		return baseColor
	}

	// Create a simple pattern based on position
	patternIndex := ((x / 16) + (y / 16)) % len(scheme.PatternColors)
	patternColor := scheme.PatternColors[patternIndex]

	// Blend base color with pattern color
	blendFactor := 0.3

	return color.RGBA{
		R: uint8(float64(baseColor.R)*(1.0-blendFactor) + float64(patternColor.R)*blendFactor),
		G: uint8(float64(baseColor.G)*(1.0-blendFactor) + float64(patternColor.G)*blendFactor),
		B: uint8(float64(baseColor.B)*(1.0-blendFactor) + float64(patternColor.B)*blendFactor),
		A: baseColor.A,
	}
}

// applyGradientVariation applies gradient-based color variation
func (ba *BlockAppearance) applyGradientVariation(baseColor color.RGBA, scheme *BlockColorScheme, depth float64) color.RGBA {
	if len(scheme.GradientColors) < 2 {
		return baseColor
	}

	// Find the two gradient colors to blend between
	depthIndex := int(depth * float64(len(scheme.GradientColors)-1))
	if depthIndex >= len(scheme.GradientColors)-1 {
		return scheme.GradientColors[len(scheme.GradientColors)-1]
	}

	color1 := scheme.GradientColors[depthIndex]
	color2 := scheme.GradientColors[depthIndex+1]

	// Calculate blend factor
	blendFactor := (depth * float64(len(scheme.GradientColors)-1)) - float64(depthIndex)

	return color.RGBA{
		R: uint8(float64(color1.R)*(1.0-blendFactor) + float64(color2.R)*blendFactor),
		G: uint8(float64(color1.G)*(1.0-blendFactor) + float64(color2.G)*blendFactor),
		B: uint8(float64(color1.B)*(1.0-blendFactor) + float64(color2.B)*blendFactor),
		A: uint8(float64(color1.A)*(1.0-blendFactor) + float64(color2.A)*blendFactor),
	}
}

// GenerateBlockTexture generates a texture for a block with color variations
// Uses texture caching for performance
func (ba *BlockAppearance) GenerateBlockTexture(blockType string, x, y int, biome string, depth float64) *ebiten.Image {
	// Create cache key based on block properties
	// Use a simplified key for performance - only cache unique variations
	variation := 0
	if ba.randSource != nil {
		variation = ba.randSource.Intn(4) // 4 variation levels
	}

	cacheKey := fmt.Sprintf("%s_%s_%d", blockType, biome, variation)

	// Try to get from cache first
	if tex, ok := textureCache.Get(cacheKey); ok {
		return tex
	}

	const textureSize = 64
	img := ebiten.NewImage(textureSize, textureSize)

	baseColor := ba.GetBlockColor(blockType, x, y, biome, depth)

	scheme, exists := ba.ColorSchemes[blockType]
	if !exists {
		// Simple solid color texture
		img.Fill(baseColor)
		textureCache.Set(cacheKey, img)
		return img
	}

	// Generate texture based on variation type
	switch scheme.VariationType {
	case VariationPattern:
		ba.generatePatternTexture(img, scheme, baseColor, x, y)
	case VariationGradient:
		ba.generateGradientTexture(img, scheme, baseColor, depth)
	default:
		ba.generateSolidTexture(img, baseColor)
	}

	// Cache the generated texture
	textureCache.Set(cacheKey, img)
	return img
}

// generateSolidTexture generates a solid color texture with slight variation
func (ba *BlockAppearance) generateSolidTexture(img *ebiten.Image, baseColor color.RGBA) {
	img.Fill(baseColor)

	// Add some noise for texture
	const size = 64
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			noise := (ba.randSource.Float64() - 0.5) * 0.1
			newR := float64(baseColor.R) + noise*255
			newG := float64(baseColor.G) + noise*255
			newB := float64(baseColor.B) + noise*255

			// Clamp to valid range
			newR = math.Max(0, math.Min(255, newR))
			newG = math.Max(0, math.Min(255, newG))
			newB = math.Max(0, math.Min(255, newB))

			pixelColor := color.RGBA{
				R: uint8(newR),
				G: uint8(newG),
				B: uint8(newB),
				A: baseColor.A,
			}
			img.Set(x, y, pixelColor)
		}
	}
}

// generatePatternTexture generates a patterned texture
func (ba *BlockAppearance) generatePatternTexture(img *ebiten.Image, scheme *BlockColorScheme, baseColor color.RGBA, blockX, blockY int) {
	const size = 64

	// Fill with base color
	img.Fill(baseColor)

	// Apply pattern
	if len(scheme.PatternColors) > 0 {
		for x := 0; x < size; x += 8 {
			for y := 0; y < size; y += 8 {
				// Create pattern based on position
				patternIndex := ((blockX + x/8) + (blockY + y/8)) % len(scheme.PatternColors)
				patternColor := scheme.PatternColors[patternIndex]

				// Draw pattern block
				for px := x; px < min(x+8, size); px++ {
					for py := y; py < min(y+8, size); py++ {
						if ba.randSource.Float64() < 0.3 { // 30% pattern coverage
							img.Set(px, py, patternColor)
						}
					}
				}
			}
		}
	}
}

// generateGradientTexture generates a gradient texture
func (ba *BlockAppearance) generateGradientTexture(img *ebiten.Image, scheme *BlockColorScheme, baseColor color.RGBA, depth float64) {
	const size = 64

	if len(scheme.GradientColors) < 2 {
		img.Fill(baseColor)
		return
	}

	// Create vertical gradient
	for y := 0; y < size; y++ {
		gradientDepth := float64(y) / float64(size)
		color := ba.applyGradientVariation(baseColor, scheme, gradientDepth)

		for x := 0; x < size; x++ {
			img.Set(x, y, color)
		}
	}
}

// Utility functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Global block appearance instance
var GlobalBlockAppearance *BlockAppearance

func init() {
	GlobalBlockAppearance = NewBlockAppearance()
	GlobalBlockAppearance.InitializeColorSchemes()
}
