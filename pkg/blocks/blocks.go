package blocks

import (
	"image/color"
	"os"

	"gopkg.in/yaml.v3"
)

// BlockType represents the type of a block
type BlockType int

const (
	AIR BlockType = iota
	DIRT
	GRASS
	STONE
	SAND
	WATER
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
)

// BlockProperties defines the properties of a block type
type BlockProperties struct {
	ID          BlockType
	Name        string
	Color       color.RGBA
	TopColor    color.RGBA // For blocks with different top face colors
	SideColor   color.RGBA // For blocks with different side face colors
	Hardness    float64
	Transparent bool
	Solid       bool
	Collectible bool
	Flammable   bool
	LightLevel  int
	Gravity     bool
	Viscosity   float64 // For liquids
}

// BlockJSON represents the YAML structure for block definitions
type BlockJSON struct {
	ID          string                 `yaml:"id"`
	Name        string                 `yaml:"name"`
	Color       []uint8                `yaml:"color"`
	TopColor    []uint8                `yaml:"topColor,omitempty"`
	SideColor   []uint8                `yaml:"sideColor,omitempty"`
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
}

// BlockDefinitions holds all block type definitions
var BlockDefinitions = make(map[string]*BlockProperties)

var BlockTypeMap = map[string]BlockType{
	"air":         AIR,
	"dirt":        DIRT,
	"grass":       GRASS,
	"stone":       STONE,
	"sand":        SAND,
	"water":       WATER,
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
}

// LoadBlocks loads block definitions from YAML files
func LoadBlocks() {
	loadBlocksFromFile("config/blocks.yaml")
	loadMods()
}

// loadBlocksFromFile loads blocks from a specific YAML file
func loadBlocksFromFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		// If file doesn't exist, skip
		return
	}
	var blocks map[string]*BlockJSON
	err = yaml.Unmarshal(data, &blocks)
	if err != nil {
		// Log error or panic
		panic(err)
	}
	for id, b := range blocks {
		props := &BlockProperties{
			ID:          BlockTypeMap[id],
			Name:        b.Name,
			Color:       color.RGBA{b.Color[0], b.Color[1], b.Color[2], b.Color[3]},
			Hardness:    b.Hardness,
			Transparent: b.Transparent,
			Solid:       b.Solid,
			Collectible: b.Collectible,
			Flammable:   b.Flammable,
			LightLevel:  b.LightLevel,
			Gravity:     b.Gravity,
			Viscosity:   b.Viscosity,
		}
		if b.TopColor != nil && len(b.TopColor) == 4 {
			props.TopColor = color.RGBA{b.TopColor[0], b.TopColor[1], b.TopColor[2], b.TopColor[3]}
		}
		if b.SideColor != nil && len(b.SideColor) == 4 {
			props.SideColor = color.RGBA{b.SideColor[0], b.SideColor[1], b.SideColor[2], b.SideColor[3]}
		}
		BlockDefinitions[id] = props
	}
}

// loadMods loads mod block definitions
func loadMods() {
	loadBlocksFromFile("mods/blocks.yaml")
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
