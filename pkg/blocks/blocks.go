package blocks

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"gopkg.in/yaml.v3"
	"tesselbox/assets"
)


// LiquidType represents the type of liquid
type LiquidType int

const (
	WATER_LIQUID LiquidType = iota
	LAVA_LIQUID
)

// LiquidProperties defines the properties of a liquid type
type LiquidProperties struct {
	Type         LiquidType
	Name         string
	Color        color.RGBA
	Density      float64  // How heavy the liquid is
	Viscosity    float64  // How thick/runny the liquid is (0-1)
	FlowSpeed    float64  // How fast the liquid flows
	LightLevel   int      // Light emission level
	Transparent  bool
	Gravity      bool
	SpreadRate   int      // How far liquid spreads horizontally
}


// BlockType represents the type of a block
type BlockType int

const (
	AIR BlockType = iota
	DIRT
	GRASS
	STONE
	SAND
	LOG
	LEAVES
	COAL_ORE
	IRON_ORE
	GOLD_ORE
	DIAMOND_ORE
	BEDROCK
	GLASS
	BRICK
	PLANK
	CACTUS
	WORKBENCH
	FURNACE
	ANVIL
	// New decorative and functional blocks
	GRAVEL
	SANDSTONE
	OBSIDIAN
	ICE
	SNOW
	TORCH
	CRAFTING_TABLE
	CHEST
	LADDER
	FENCE
	GATE
	DOOR
	WINDOW
	FLOWER
	TALL_GRASS
	MUSHROOM_RED
	MUSHROOM_BROWN
	WOOL
	BOOKSHELF
	JUKEBOX
	NOTE_BLOCK
	PUMPKIN
	MELON
	HAY_BALE
	COBBLESTONE
	MOSSY_COBBLESTONE
	STONE_BRICKS
	CHISELED_STONE
)

// BlockProperties defines the properties of a block type
type BlockProperties struct {
	ID          BlockType
	Name        string
	Color       color.RGBA
	TopColor    color.RGBA   // For blocks with different top face colors
	SideColor   color.RGBA   // For blocks with different side face colors
	Colors      []color.RGBA // Palette of colors for procedural textures
	Hardness    float64
	Transparent bool
	Solid       bool
	Collectible bool
	Flammable   bool
	LightLevel  int
	Gravity     bool
	Viscosity   float64       // For liquids
	Pattern     string        // "solid", "striped", "checkerboard", etc.
	Texture     *ebiten.Image // Optional texture for pixel-by-pixel appearance
}

// BlockJSON represents the YAML structure for block definitions
type BlockJSON struct {
	ID          string                 `yaml:"id"`
	Name        string                 `yaml:"name"`
	Color       []uint8                `yaml:"color"`
	TopColor    []uint8                `yaml:"topColor,omitempty"`
	SideColor   []uint8                `yaml:"sideColor,omitempty"`
	Colors      [][]uint8              `yaml:"colors,omitempty"`
	Hardness    float64                `yaml:"hardness"`
	Transparent bool                   `yaml:"transparent"`
	Solid       bool                   `yaml:"solid"`
	Collectible bool                   `yaml:"collectible"`
	Flammable   bool                   `yaml:"flammable"`
	LightLevel  int                    `yaml:"lightLevel"`
	Gravity     bool                   `yaml:"gravity"`
	Viscosity   float64                `yaml:"viscosity"`
	Pattern     string                 `yaml:"pattern"`
	UI          map[string]interface{} `yaml:"ui"`
	Function    map[string]interface{} `yaml:"function"`
	Texture     string                 `yaml:"texture,omitempty"`
}

// BlockDefinitions holds all block type definitions
var BlockDefinitions = make(map[string]*BlockProperties)

// LiquidDefinitions holds all liquid type definitions
var LiquidDefinitions = make(map[LiquidType]*LiquidProperties)

// Initialize liquid definitions
func init() {
	LiquidDefinitions[WATER_LIQUID] = &LiquidProperties{
		Type:        WATER_LIQUID,
		Name:        "Water",
		Color:       color.RGBA{64, 164, 223, 180},
		Density:     1.0,
		Viscosity:   0.3,  // Quite runny
		FlowSpeed:   2.0,  // Flows moderately fast
		LightLevel:  0,
		Transparent: true,
		Gravity:     true,
		SpreadRate:  7,    // Spreads up to 7 blocks horizontally
	}
	
	LiquidDefinitions[LAVA_LIQUID] = &LiquidProperties{
		Type:        LAVA_LIQUID,
		Name:        "Lava",
		Color:       color.RGBA{255, 100, 0, 200},
		Density:     3.0,  // Much heavier than water
		Viscosity:   0.8,  // Thick and slow
		FlowSpeed:   0.5,  // Flows slowly
		LightLevel:  12,   // Emits light
		Transparent: false,
		Gravity:     true,
		SpreadRate:  3,    // Doesn't spread as far
	}
}

var BlockTypeMap = map[string]BlockType{
	"air":         AIR,
	"dirt":        DIRT,
	"grass":       GRASS,
	"stone":       STONE,
	"sand":        SAND,
	"log":         LOG,
	"leaves":      LEAVES,
	"coal_ore":    COAL_ORE,
	"iron_ore":    IRON_ORE,
	"gold_ore":    GOLD_ORE,
	"diamond_ore": DIAMOND_ORE,
	"bedrock":     BEDROCK,
	"glass":       GLASS,
	"brick":       BRICK,
	"plank":       PLANK,
	"cactus":      CACTUS,
	"workbench":   WORKBENCH,
	"furnace":     FURNACE,
	"anvil":       ANVIL,
	// New blocks
	"gravel":           GRAVEL,
	"sandstone":       SANDSTONE,
	"obsidian":        OBSIDIAN,
	"ice":             ICE,
	"snow":            SNOW,
	"torch":           TORCH,
	"crafting_table":  CRAFTING_TABLE,
	"chest":           CHEST,
	"ladder":          LADDER,
	"fence":           FENCE,
	"gate":            GATE,
	"door":            DOOR,
	"window":          WINDOW,
	"flower":          FLOWER,
	"tall_grass":      TALL_GRASS,
	"mushroom_red":    MUSHROOM_RED,
	"mushroom_brown":  MUSHROOM_BROWN,
	"wool":            WOOL,
	"bookshelf":       BOOKSHELF,
	"jukebox":         JUKEBOX,
	"note_block":      NOTE_BLOCK,
	"pumpkin":         PUMPKIN,
	"melon":           MELON,
	"hay_bale":        HAY_BALE,
	"cobblestone":     COBBLESTONE,
	"mossy_cobblestone": MOSSY_COBBLESTONE,
	"stone_bricks":    STONE_BRICKS,
	"chiseled_stone":  CHISELED_STONE,
}

// LoadBlocks loads block definitions from YAML files
func LoadBlocks() {
	LoadBlocksFromAssets()
	loadMods()
}

// LoadBlocksFromAssets loads block definitions from embedded assets
func LoadBlocksFromAssets() {
	loadBlocksFromEmbedded()
}

// loadBlocksFromEmbedded loads blocks from embedded YAML data with improved validation
func loadBlocksFromEmbedded() {
	data, err := assets.GetConfigFile("blocks.yaml")
	if err != nil {
		log.Printf("Warning: Failed to load blocks.yaml from embedded assets: %v", err)
		log.Printf("Loading default block configurations...")
		loadDefaultBlocks()
		return
	}
	
	var blocks map[string]*BlockJSON
	if err := yaml.Unmarshal(data, &blocks); err != nil {
		log.Printf("Warning: Failed to parse blocks.yaml: %v", err)
		log.Printf("Loading default block configurations...")
		loadDefaultBlocks()
		return
	}
	
	// Validate loaded blocks
	if len(blocks) == 0 {
		log.Printf("Warning: No valid blocks found in blocks.yaml")
		log.Printf("Loading default block configurations...")
		loadDefaultBlocks()
		return
	}
	
	loadedCount := 0
	for id, b := range blocks {
		// Validate block data
		if b.Name == "" || len(b.Color) < 4 {
			log.Printf("Warning: Invalid block data for %s, skipping", id)
			continue
		}
		
		props := &BlockProperties{
			ID:          BlockTypeMap[id],
			Name:        b.Name,
			Color:       color.RGBA{b.Color[0], b.Color[1], b.Color[2], b.Color[3]},
			Hardness:    validateHardness(b.Hardness),
			Transparent: validateBool(b.Transparent, false),
			Solid:       validateBool(b.Solid, true),
			Collectible: validateBool(b.Collectible, true),
			Flammable:   validateBool(b.Flammable, false),
			LightLevel:  validateLightLevel(b.LightLevel),
			Gravity:     validateBool(b.Gravity, false),
			Viscosity:   validateViscosity(b.Viscosity),
			Pattern:     b.Pattern,
		}
		
		if len(b.TopColor) == 4 {
			props.TopColor = color.RGBA{b.TopColor[0], b.TopColor[1], b.TopColor[2], b.TopColor[3]}
		}
		if len(b.SideColor) == 4 {
			props.SideColor = color.RGBA{b.SideColor[0], b.SideColor[1], b.SideColor[2], b.SideColor[3]}
		}

		// Parse color palette
		if len(b.Colors) > 0 {
			props.Colors = make([]color.RGBA, len(b.Colors))
			for i, c := range b.Colors {
				if len(c) == 4 {
					props.Colors[i] = color.RGBA{c[0], c[1], c[2], c[3]}
				}
			}
		}

		BlockDefinitions[id] = props
		loadedCount++
	}
	
	if loadedCount == 0 {
		log.Printf("Warning: No valid blocks could be loaded from blocks.yaml")
		log.Printf("Loading default block configurations...")
		loadDefaultBlocks()
	} else {
		log.Printf("Successfully loaded %d block configurations", loadedCount)
	}
}

// Validation functions for block properties
func validateHardness(hardness float64) float64 {
	if hardness < 0 {
		return 1.0 // Default hardness
	}
	if hardness > 100 {
		return 100.0 // Cap maximum hardness
	}
	return hardness
}

func validateBool(value bool, defaultValue bool) bool {
	// In YAML parsing, this is mostly a pass-through, but useful for consistency
	return value
}

func validateLightLevel(level int) int {
	if level < 0 {
		return 0
	}
	if level > 15 {
		return 15
	}
	return level
}

func validateViscosity(viscosity float64) float64 {
	if viscosity < 0 {
		return 0.0
	}
	if viscosity > 1.0 {
		return 1.0
	}
	return viscosity
}

// loadDefaultBlocks loads essential default block configurations
func loadDefaultBlocks() {
	defaultBlocks := map[string]*BlockProperties{
		"air": {
			ID:          AIR,
			Name:        "Air",
			Color:       color.RGBA{0, 0, 0, 0},
			Hardness:    0,
			Transparent: true,
			Solid:       false,
			Collectible: false,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     false,
		},
		"dirt": {
			ID:          DIRT,
			Name:        "Dirt",
			Color:       color.RGBA{139, 90, 43, 255},
			Hardness:    0.5,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		"grass": {
			ID:          GRASS,
			Name:        "Grass",
			Color:       color.RGBA{124, 169, 84, 255},
			TopColor:    color.RGBA{124, 169, 84, 255},
			SideColor:   color.RGBA{139, 90, 43, 255},
			Hardness:    0.6,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		"stone": {
			ID:          STONE,
			Name:        "Stone",
			Color:       color.RGBA{128, 128, 128, 255},
			Hardness:    1.5,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		"sand": {
			ID:          SAND,
			Name:        "Sand",
			Color:       color.RGBA{238, 203, 173, 255},
			Hardness:    0.5,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		"log": {
			ID:          LOG,
			Name:        "Log",
			Color:       color.RGBA{139, 69, 19, 255},
			Hardness:    1.0,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   true,
			LightLevel:  0,
			Gravity:     true,
		},
		"leaves": {
			ID:          LEAVES,
			Name:        "Leaves",
			Color:       color.RGBA{34, 139, 34, 200},
			Hardness:    0.2,
			Transparent: true,
			Solid:       true,
			Collectible: true,
			Flammable:   true,
			LightLevel:  0,
			Gravity:     false,
		},
		"coal_ore": {
			ID:          COAL_ORE,
			Name:        "Coal Ore",
			Color:       color.RGBA{54, 54, 54, 255},
			Hardness:    1.5,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		"iron_ore": {
			ID:          IRON_ORE,
			Name:        "Iron Ore",
			Color:       color.RGBA{183, 183, 183, 255},
			Hardness:    2.0,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		"gold_ore": {
			ID:          GOLD_ORE,
			Name:        "Gold Ore",
			Color:       color.RGBA{255, 215, 0, 255},
			Hardness:    2.0,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		"diamond_ore": {
			ID:          DIAMOND_ORE,
			Name:        "Diamond Ore",
			Color:       color.RGBA{185, 242, 255, 255},
			Hardness:    3.0,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		"bedrock": {
			ID:          BEDROCK,
			Name:        "Bedrock",
			Color:       color.RGBA{64, 64, 64, 255},
			Hardness:    -1, // Unbreakable
			Transparent: false,
			Solid:       true,
			Collectible: false,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     false,
		},
		"workbench": {
			ID:          WORKBENCH,
			Name:        "Workbench",
			Color:       color.RGBA{139, 69, 19, 255},
			Hardness:    1.0,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   true,
			LightLevel:  0,
			Gravity:     true,
		},
		"furnace": {
			ID:          FURNACE,
			Name:        "Furnace",
			Color:       color.RGBA{169, 169, 169, 255},
			Hardness:    1.5,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		"anvil": {
			ID:          ANVIL,
			Name:        "Anvil",
			Color:       color.RGBA{192, 192, 192, 255},
			Hardness:    2.0,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		// New decorative and functional blocks
		"gravel": {
			ID:          GRAVEL,
			Name:        "Gravel",
			Color:       color.RGBA{136, 140, 141, 255},
			Hardness:    0.6,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		"sandstone": {
			ID:          SANDSTONE,
			Name:        "Sandstone",
			Color:       color.RGBA{238, 203, 173, 255},
			Hardness:    0.8,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		"obsidian": {
			ID:          OBSIDIAN,
			Name:        "Obsidian",
			Color:       color.RGBA{27, 23, 23, 255},
			Hardness:    5.0,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		"ice": {
			ID:          ICE,
			Name:        "Ice",
			Color:       color.RGBA{175, 223, 255, 200},
			Hardness:    0.5,
			Transparent: true,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     false,
		},
		"snow": {
			ID:          SNOW,
			Name:        "Snow",
			Color:       color.RGBA{255, 255, 255, 255},
			Hardness:    0.2,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     false,
		},
		"torch": {
			ID:          TORCH,
			Name:        "Torch",
			Color:       color.RGBA{255, 200, 100, 255},
			Hardness:    0.1,
			Transparent: true,
			Solid:       false,
			Collectible: true,
			Flammable:   false,
			LightLevel:  14,
			Gravity:     false,
		},
		"crafting_table": {
			ID:          CRAFTING_TABLE,
			Name:        "Crafting Table",
			Color:       color.RGBA{139, 90, 43, 255},
			Hardness:    1.0,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   true,
			LightLevel:  0,
			Gravity:     true,
		},
		"chest": {
			ID:          CHEST,
			Name:        "Chest",
			Color:       color.RGBA{139, 90, 19, 255},
			Hardness:    1.0,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   true,
			LightLevel:  0,
			Gravity:     true,
		},
		"wool": {
			ID:          WOOL,
			Name:        "Wool",
			Color:       color.RGBA{222, 222, 222, 255},
			Hardness:    0.3,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   true,
			LightLevel:  0,
			Gravity:     true,
		},
		"flower": {
			ID:          FLOWER,
			Name:        "Flower",
			Color:       color.RGBA{255, 100, 100, 255},
			Hardness:    0.1,
			Transparent: true,
			Solid:       false,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     false,
		},
		"pumpkin": {
			ID:          PUMPKIN,
			Name:        "Pumpkin",
			Color:       color.RGBA{255, 140, 0, 255},
			Hardness:    0.5,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
		"cobblestone": {
			ID:          COBBLESTONE,
			Name:        "Cobblestone",
			Color:       color.RGBA{128, 128, 128, 255},
			Hardness:    1.5,
			Transparent: false,
			Solid:       true,
			Collectible: true,
			Flammable:   false,
			LightLevel:  0,
			Gravity:     true,
		},
	}

	// Load default blocks
	loadedCount := 0
	for id, props := range defaultBlocks {
		BlockDefinitions[id] = props
		loadedCount++
	}
	
	log.Printf("Successfully loaded %d default block configurations", loadedCount)
}

// loadMods loads mod block definitions
func loadMods() {
	// Mods are currently disabled in embedded build
}

// generateProceduralTexture creates a texture using a color palette with random pixels
func generateProceduralTexture(colors []color.RGBA, id BlockType) *ebiten.Image {
	if len(colors) == 0 {
		return nil
	}
	const size = 64
	img := ebiten.NewImage(size, size)
	rand.Seed(int64(id) * 1000) // Deterministic seed per block type
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			idx := rand.Intn(len(colors))
			img.Set(x, y, colors[idx])
		}
	}
	return img
}

// ColorByType returns the color for a block type string
func ColorByType(blockType string) color.RGBA {
	if props, ok := BlockDefinitions[blockType]; ok {
		return props.Color
	}
	return BlockDefinitions["dirt"].Color
}

// HardnessByType returns the hardness for a block type string
func HardnessByType(blockType string) float64 {
	if props, ok := BlockDefinitions[blockType]; ok {
		return props.Hardness
	}
	return 1.0
}

// TransparentByType returns if a block type is transparent
func TransparentByType(blockType string) bool {
	if props, ok := BlockDefinitions[blockType]; ok {
		return props.Transparent
	}
	return false
}

// SolidByType returns if a block type is solid
func SolidByType(blockType string) bool {
	if props, ok := BlockDefinitions[blockType]; ok {
		return props.Solid
	}
	return true
}

// CollectibleByType returns if a block type is collectible
func CollectibleByType(blockType string) bool {
	if props, ok := BlockDefinitions[blockType]; ok {
		return props.Collectible
	}
	return true
}

// ValidBlockType checks if a block type string is valid
func ValidBlockType(blockType string) bool {
	_, ok := BlockDefinitions[blockType]
	return ok
}

// GetLiquidProperties returns the properties of a liquid type
func GetLiquidProperties(liquidType LiquidType) *LiquidProperties {
	return LiquidDefinitions[liquidType]
}

// IsValidLiquid checks if a liquid type is valid
func IsValidLiquid(liquidType LiquidType) bool {
	_, ok := LiquidDefinitions[liquidType]
	return ok
}

// LiquidPhysics represents the physics state of a liquid cell
type LiquidPhysics struct {
	LiquidType   LiquidType
	Level        float64  // 0.0 to 1.0 (empty to full)
	Flowing      bool     // Whether this liquid is actively flowing
	Source       bool     // Whether this is a source block
	UpdateTime   float64  // Last time this liquid was updated
	Pressure     float64  // Liquid pressure for physics calculations
}

// UpdateLiquidPhysics updates liquid flow based on gravity and terrain
func UpdateLiquidPhysics(liquid *LiquidPhysics, deltaTime float64, neighbors []*LiquidPhysics) {
	props := GetLiquidProperties(liquid.LiquidType)
	if props == nil {
		return
	}
	
	// Apply gravity
	if props.Gravity && liquid.Level > 0.0 {
		// Check if liquid can flow down
		canFlowDown := true
		for _, neighbor := range neighbors {
			if neighbor != nil && neighbor.Level < 1.0 {
				canFlowDown = true
				break
			}
		}
		
		if canFlowDown {
			// Calculate flow rate based on viscosity
			flowRate := props.FlowSpeed * deltaTime * (1.0 - props.Viscosity)
			liquid.Level = max(0.0, liquid.Level-flowRate)
			liquid.Flowing = true
		}
	}
	
	// Horizontal spreading
	if liquid.Level > 0.0 && !liquid.Source {
		// Spread to adjacent cells based on pressure and viscosity
		spreadAmount := (liquid.Level * props.FlowSpeed * deltaTime) / float64(props.SpreadRate)
		spreadAmount *= (1.0 - props.Viscosity) // More viscous liquids spread less
		
		liquid.Level = max(0.0, liquid.Level-spreadAmount)
		if spreadAmount > 0.01 {
			liquid.Flowing = true
		}
	}
	
	// Update time
	liquid.UpdateTime += deltaTime
	
	// Stop flowing if level is too low
	if liquid.Level < 0.01 {
		liquid.Level = 0.0
		liquid.Flowing = false
	}
}

// CalculateLiquidShape calculates the visual shape of liquid based on level and terrain
func CalculateLiquidShape(liquid *LiquidPhysics, terrainHeight float64) []float64 {
	props := GetLiquidProperties(liquid.LiquidType)
	if props == nil || liquid.Level <= 0.0 {
		return []float64{}
	}
	
	// Basic liquid shape - can be enhanced for more realistic rendering
	shape := make([]float64, 8) // 4 points for a quad
	
	// Calculate liquid surface height based on level and viscosity
	surfaceHeight := terrainHeight + (liquid.Level * 0.8) // Liquid doesn't fill completely
	
	// Adjust for viscosity (thicker liquids have flatter surfaces)
	if props.Viscosity > 0.5 {
		surfaceHeight = terrainHeight + (liquid.Level * 0.9)
	}
	
	// Simple quad shape - can be replaced with more complex liquid mesh
	shape[0] = 0.0 // x1
	shape[1] = surfaceHeight // y1
	shape[2] = 1.0 // x2
	shape[3] = surfaceHeight // y2
	shape[4] = 1.0 // x3
	shape[5] = terrainHeight // y3
	shape[6] = 0.0 // x4
	shape[7] = terrainHeight // y4
	
	return shape
}
