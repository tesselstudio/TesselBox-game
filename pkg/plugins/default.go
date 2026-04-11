package plugins

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"tesselbox/pkg/audio"
	"tesselbox/pkg/blocks"
	"tesselbox/pkg/creatures"
	"tesselbox/pkg/organisms"
	"tesselbox/pkg/world"
)

// DefaultPlugin contains all the original game content
type DefaultPlugin struct {
	initialized bool
}

// NewDefaultPlugin creates the default game content plugin
func NewDefaultPlugin() *DefaultPlugin {
	return &DefaultPlugin{}
}

// ID returns the plugin identifier
func (dp *DefaultPlugin) ID() string {
	return "default"
}

// Name returns the plugin name
func (dp *DefaultPlugin) Name() string {
	return "TesselBox Default Content"
}

// Version returns the plugin version
func (dp *DefaultPlugin) Version() string {
	return "1.0.0"
}

// Description returns the plugin description
func (dp *DefaultPlugin) Description() string {
	return "Default game content including blocks, creatures, organisms, and audio"
}

// Author returns the plugin author
func (dp *DefaultPlugin) Author() string {
	return "TesselBox Team"
}

// Initialize sets up the default content
func (dp *DefaultPlugin) Initialize() error {
	if dp.initialized {
		return nil
	}

	log.Println("Initializing default content plugin...")
	dp.initialized = true
	log.Println("Default content plugin initialized successfully")
	return nil
}

// Shutdown cleans up default content
func (dp *DefaultPlugin) Shutdown() error {
	if !dp.initialized {
		return nil
	}

	log.Println("Shutting down default content plugin...")
	dp.initialized = false
	log.Println("Default content plugin shut down successfully")
	return nil
}

// GetBlockTypes returns all available block types
func (dp *DefaultPlugin) GetBlockTypes() []blocks.BlockType {
	// Return core block types - these would be loaded from config
	return []blocks.BlockType{
		blocks.AIR, blocks.DIRT, blocks.GRASS, blocks.STONE,
		blocks.SAND, blocks.WATER, blocks.LOG, blocks.LEAVES,
		blocks.COAL_ORE, blocks.IRON_ORE, blocks.GOLD_ORE, blocks.DIAMOND_ORE,
		blocks.BEDROCK, blocks.GLASS, blocks.BRICK, blocks.PLANK,
		blocks.CACTUS, blocks.WORKBENCH, blocks.FURNACE, blocks.ANVIL,
		blocks.GRAVEL, blocks.SANDSTONE, blocks.OBSIDIAN, blocks.ICE,
		blocks.SNOW, blocks.TORCH, blocks.CRAFTING_TABLE, blocks.CHEST,
		blocks.LADDER, blocks.FENCE, blocks.GATE, blocks.DOOR,
		blocks.CHISELED_STONE,
	}
}

// GetBlockDefinition returns a specific block definition
func (dp *DefaultPlugin) GetBlockDefinition(blockType blocks.BlockType) (*BlockDefinition, bool) {
	// Create a basic block definition
	def := &BlockDefinition{
		Type:        blockType,
		Name:        stringToBlockName(blockType),
		Hardness:    1.0,
		Color:       getDefaultBlockColor(blockType),
		Transparent: isBlockTransparent(blockType),
		Solid:       isBlockSolid(blockType),
	}
	return def, true
}

// GetBlockProperties returns additional block properties
func (dp *DefaultPlugin) GetBlockProperties(blockType blocks.BlockType) (map[string]interface{}, bool) {
	properties := map[string]interface{}{
		"hardness":    getDefaultBlockHardness(blockType),
		"light_level": getDefaultBlockLightLevel(blockType),
		"flammable":   isBlockFlammable(blockType),
		"breakable":   isBlockBreakable(blockType),
	}
	return properties, true
}

// GetCreatureTypes returns all available creature types
func (dp *DefaultPlugin) GetCreatureTypes() []creatures.CreatureType {
	return []creatures.CreatureType{
		creatures.SLIME,
		creatures.SPIDER,
		creatures.ZOMBIE,
	}
}

// GetCreatureDefinition returns a specific creature definition
func (dp *DefaultPlugin) GetCreatureDefinition(creatureType creatures.CreatureType) (*CreatureDefinition, bool) {
	// Create a basic creature definition
	def := &CreatureDefinition{
		Type:   creatureType,
		Name:   stringToCreatureName(creatureType),
		Health: getDefaultCreatureHealth(creatureType),
		Damage: getDefaultCreatureDamage(creatureType),
		Speed:  getDefaultCreatureSpeed(creatureType),
		Color:  getDefaultCreatureColor(creatureType),
	}
	return def, true
}

// GetOrganismTypes returns all available organism types
func (dp *DefaultPlugin) GetOrganismTypes() []organisms.OrganismType {
	return []organisms.OrganismType{
		organisms.TREE,
		organisms.BUSH,
		organisms.FLOWER,
		organisms.MUSHROOM,
		organisms.VENUS_FLYTRAP,
	}
}

// GetOrganismDefinition returns a specific organism definition
func (dp *DefaultPlugin) GetOrganismDefinition(organismType organisms.OrganismType) (*OrganismDefinition, bool) {
	// Create a basic organism definition
	def := &OrganismDefinition{
		Type:   organismType,
		Name:   stringToOrganismName(organismType),
		Height: getDefaultOrganismHeight(organismType),
		Color:  getDefaultOrganismColor(organismType),
	}
	return def, true
}

// GetAudioTypes returns all available audio types
func (dp *DefaultPlugin) GetAudioTypes() []audio.AudioType {
	return []audio.AudioType{
		audio.AudioTypeMusic,
		audio.AudioTypeSFX,
		audio.AudioTypeAmbient,
	}
}

// GetAudioDefinition returns a specific audio definition
func (dp *DefaultPlugin) GetAudioDefinition(audioType audio.AudioType) (*AudioDefinition, bool) {
	// Create a basic audio definition
	def := &AudioDefinition{
		Type:   audioType,
		Name:   stringToAudioName(audioType),
		Volume: getDefaultAudioVolume(audioType),
	}
	return def, true
}

// GenerateChunk handles world generation for default content
func (dp *DefaultPlugin) GenerateChunk(world *world.World, chunkX, chunkZ int) error {
	log.Printf("Generating chunk at %d,%d", chunkX, chunkZ)
	return nil
}

// SpawnOrganisms handles organism spawning for default content
func (dp *DefaultPlugin) SpawnOrganisms(world *world.World) error {
	// Spawn basic organisms
	for i := 0; i < 50; i++ {
		dp.spawnRandomOrganism(world)
	}
	return nil
}

// SpawnCreatures handles creature spawning for default content
func (dp *DefaultPlugin) SpawnCreatures(world *world.World) error {
	// Spawn basic creatures
	for i := 0; i < 10; i++ {
		dp.spawnRandomCreature(world)
	}
	return nil
}

// OnBlockPlaced handles block placement events
func (dp *DefaultPlugin) OnBlockPlaced(x, y, z int, blockType blocks.BlockType) error {
	log.Printf("Block placed: %s at (%d,%d,%d)", stringToBlockName(blockType), x, y, z)
	return nil
}

// OnBlockBroken handles block breaking events
func (dp *DefaultPlugin) OnBlockBroken(x, y, z int, blockType blocks.BlockType) error {
	log.Printf("Block broken: %s at (%d,%d,%d)", stringToBlockName(blockType), x, y, z)
	return nil
}

// OnCreatureSpawn handles creature spawn events
func (dp *DefaultPlugin) OnCreatureSpawn(creature *creatures.Creature) error {
	log.Printf("Creature spawned: %s", stringToCreatureName(creature.Type))
	return nil
}

// OnCreatureDeath handles creature death events
func (dp *DefaultPlugin) OnCreatureDeath(creature *creatures.Creature) error {
	log.Printf("Creature died: %s", stringToCreatureName(creature.Type))
	return nil
}

// OnTick handles per-tick updates
func (dp *DefaultPlugin) OnTick(world *world.World, deltaTime float64) error {
	return nil
}

// Private helper methods

func (dp *DefaultPlugin) spawnRandomOrganism(world *world.World) {
	organismTypes := []organisms.OrganismType{
		organisms.TREE,
		organisms.BUSH,
		organisms.FLOWER,
		organisms.MUSHROOM,
		organisms.VENUS_FLYTRAP,
	}

	organismType := organismTypes[rand.Intn(len(organismTypes))]

	// Create a random position
	x := rand.Intn(100)
	y := 50

	organism := &organisms.Organism{
		Type: organismType,
		X:    float64(x),
		Y:    float64(y),
	}

	world.Organisms = append(world.Organisms, organism)
}

func (dp *DefaultPlugin) spawnRandomCreature(world *world.World) {
	creatureTypes := []creatures.CreatureType{
		creatures.SLIME,
		creatures.SPIDER,
		creatures.ZOMBIE,
	}

	creatureType := creatureTypes[rand.Intn(len(creatureTypes))]

	// Create a random position
	x := rand.Intn(100)
	y := 50

	creature := &creatures.Creature{
		ID:        fmt.Sprintf("creature_%d", rand.Intn(10000)),
		Type:      creatureType,
		X:         float64(x),
		Y:         float64(y),
		Health:    getDefaultCreatureHealth(creatureType),
		MaxHealth: getDefaultCreatureHealth(creatureType),
		Damage:    getDefaultCreatureDamage(creatureType),
		Speed:     getDefaultCreatureSpeed(creatureType),
		AIState:   creatures.IDLE,
		IsHostile: isCreatureHostile(creatureType),
		SpawnTime: time.Now(),
	}

	world.Creatures = append(world.Creatures, creature)
}

// Simple helper functions for demonstration
func stringToBlockName(blockType blocks.BlockType) string {
	names := map[blocks.BlockType]string{
		blocks.AIR:         "Air",
		blocks.DIRT:        "Dirt",
		blocks.GRASS:       "Grass",
		blocks.STONE:       "Stone",
		blocks.WATER:       "Water",
		blocks.COAL_ORE:    "Coal Ore",
		blocks.IRON_ORE:    "Iron Ore",
		blocks.DIAMOND_ORE: "Diamond Ore",
	}
	if name, exists := names[blockType]; exists {
		return name
	}
	return "Unknown"
}

func stringToCreatureName(creatureType creatures.CreatureType) string {
	names := map[creatures.CreatureType]string{
		creatures.SLIME:  "Slime",
		creatures.SPIDER: "Spider",
		creatures.ZOMBIE: "Zombie",
	}
	if name, exists := names[creatureType]; exists {
		return name
	}
	return "Unknown"
}

func stringToOrganismName(organismType organisms.OrganismType) string {
	names := map[organisms.OrganismType]string{
		organisms.TREE:          "Tree",
		organisms.BUSH:          "Bush",
		organisms.FLOWER:        "Flower",
		organisms.MUSHROOM:      "Mushroom",
		organisms.VENUS_FLYTRAP: "Venus Flytrap",
	}
	if name, exists := names[organismType]; exists {
		return name
	}
	return "Unknown"
}

func stringToAudioName(audioType audio.AudioType) string {
	names := map[audio.AudioType]string{
		audio.AudioTypeMusic:   "Music",
		audio.AudioTypeSFX:     "Sound Effects",
		audio.AudioTypeAmbient: "Ambient",
	}
	if name, exists := names[audioType]; exists {
		return name
	}
	return "Unknown"
}

func getDefaultBlockColor(blockType blocks.BlockType) string {
	colors := map[blocks.BlockType]string{
		blocks.AIR:   "#FFFFFF",
		blocks.DIRT:  "#8B4513",
		blocks.GRASS: "#00FF00",
		blocks.STONE: "#808080",
		blocks.WATER: "#0000FF",
	}
	if color, exists := colors[blockType]; exists {
		return color
	}
	return "#FFFFFF"
}

func isBlockTransparent(blockType blocks.BlockType) bool {
	return blockType == blocks.AIR || blockType == blocks.WATER
}

func isBlockSolid(blockType blocks.BlockType) bool {
	return blockType != blocks.AIR
}

func getDefaultBlockHardness(blockType blocks.BlockType) float64 {
	hardness := map[blocks.BlockType]float64{
		blocks.AIR:     0.0,
		blocks.DIRT:    0.5,
		blocks.GRASS:   0.6,
		blocks.STONE:   1.5,
		blocks.WATER:   0.0,
		blocks.BEDROCK: 50.0,
		blocks.GLASS:   0.3,
	}
	if h, exists := hardness[blockType]; exists {
		return h
	}
	return 1.0
}

func getDefaultBlockLightLevel(blockType blocks.BlockType) int {
	light := map[blocks.BlockType]int{
		blocks.AIR:   0,
		blocks.TORCH: 14,
		blocks.WATER: 0,
		blocks.GLASS: 0,
	}
	if l, exists := light[blockType]; exists {
		return l
	}
	return 0
}

func isBlockFlammable(blockType blocks.BlockType) bool {
	flammable := map[blocks.BlockType]bool{
		blocks.LOG:    true,
		blocks.LEAVES: true,
		blocks.GRASS:  true,
		blocks.FLOWER: true,
		blocks.WOOL:   true,
		blocks.PLANK:  true,
	}
	if f, exists := flammable[blockType]; exists {
		return f
	}
	return false
}

func isBlockBreakable(blockType blocks.BlockType) bool {
	unbreakable := map[blocks.BlockType]bool{
		blocks.BEDROCK: true,
		blocks.WATER:   true,
		blocks.AIR:     false,
	}
	if ub, exists := unbreakable[blockType]; exists {
		return !ub
	}
	return true
}

func getDefaultCreatureHealth(creatureType creatures.CreatureType) float64 {
	health := map[creatures.CreatureType]float64{
		creatures.SLIME:  10,
		creatures.SPIDER: 15,
		creatures.ZOMBIE: 20,
	}
	if h, exists := health[creatureType]; exists {
		return h
	}
	return 10
}

func getDefaultCreatureDamage(creatureType creatures.CreatureType) float64 {
	damage := map[creatures.CreatureType]float64{
		creatures.SLIME:  2,
		creatures.SPIDER: 4,
		creatures.ZOMBIE: 6,
	}
	if d, exists := damage[creatureType]; exists {
		return d
	}
	return 2
}

func getDefaultCreatureSpeed(creatureType creatures.CreatureType) float64 {
	speed := map[creatures.CreatureType]float64{
		creatures.SLIME:  1.0,
		creatures.SPIDER: 1.5,
		creatures.ZOMBIE: 0.8,
	}
	if s, exists := speed[creatureType]; exists {
		return s
	}
	return 1.0
}

func getDefaultCreatureColor(creatureType creatures.CreatureType) string {
	colors := map[creatures.CreatureType]string{
		creatures.SLIME:  "#00FF00",
		creatures.SPIDER: "#8B4513",
		creatures.ZOMBIE: "#808080",
	}
	if color, exists := colors[creatureType]; exists {
		return color
	}
	return "#FFFFFF"
}

func isCreatureHostile(creatureType creatures.CreatureType) bool {
	return creatureType == creatures.ZOMBIE || creatureType == creatures.SPIDER
}

func getDefaultOrganismHeight(organismType organisms.OrganismType) float64 {
	height := map[organisms.OrganismType]float64{
		organisms.TREE:          8.0,
		organisms.BUSH:          2.0,
		organisms.FLOWER:        1.0,
		organisms.MUSHROOM:      1.5,
		organisms.VENUS_FLYTRAP: 1.0,
	}
	if h, exists := height[organismType]; exists {
		return h
	}
	return 2.0
}

func getDefaultOrganismColor(organismType organisms.OrganismType) string {
	colors := map[organisms.OrganismType]string{
		organisms.TREE:          "#8B4513",
		organisms.BUSH:          "#228B22",
		organisms.FLOWER:        "#FF69B4",
		organisms.MUSHROOM:      "#FFD700",
		organisms.VENUS_FLYTRAP: "#FF1493",
	}
	if color, exists := colors[organismType]; exists {
		return color
	}
	return "#FFFFFF"
}

func getDefaultAudioVolume(audioType audio.AudioType) float64 {
	volume := map[audio.AudioType]float64{
		audio.AudioTypeMusic:   0.7,
		audio.AudioTypeSFX:     1.0,
		audio.AudioTypeAmbient: 0.5,
	}
	if v, exists := volume[audioType]; exists {
		return v
	}
	return 0.8
}
