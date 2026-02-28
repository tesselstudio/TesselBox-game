package organisms

import (
	"math"
	"os"
	"tesselbox/pkg/hexagon"
	"time"

	"gopkg.in/yaml.v3"
)

// OrganismType represents the type of organism
type OrganismType int

const (
	TREE OrganismType = iota
	BUSH
	FLOWER
	MUSHROOM
	VENUS_FLYTRAP
)

var OrganismDefinitions = make(map[string]*OrganismJSON)

var OrganismTypeMap = map[string]OrganismType{
	"tree":          TREE,
	"bush":          BUSH,
	"flower":        FLOWER,
	"mushroom":      MUSHROOM,
	"venus_flytrap": VENUS_FLYTRAP,
}

// LoadOrganisms loads organism definitions from YAML
func LoadOrganisms() {
	data, err := os.ReadFile("config/organisms.yaml")
	if err != nil {
		return
	}
	var orgs map[string]*OrganismJSON
	err = yaml.Unmarshal(data, &orgs)
	if err != nil {
		panic(err)
	}
	OrganismDefinitions = orgs
}

func init() {
	LoadOrganisms()
}

// OrganismJSON represents the YAML structure for organism definitions
type OrganismJSON struct {
	ID         string                 `yaml:"id"`
	Name       string                 `yaml:"name"`
	Type       string                 `yaml:"type"`
	Appearance map[string]interface{} `yaml:"appearance"`
	Properties map[string]interface{} `yaml:"properties"`
	Behavior   map[string]interface{} `yaml:"behavior"`
	Function   map[string]interface{} `yaml:"function"`
	Drops      []string               `yaml:"drops"`
}

// Organism represents a living organism in the world
type Organism struct {
	ID             string
	Type           OrganismType
	TypeString     string
	X, Y           float64
	Hex            hexagon.Hexagon
	Health         float64
	MaxHealth      float64
	IsHostile      bool
	Damage         float64
	AttackRange    float64
	LastAttackTime int64 // Unix timestamp
}

// Tree represents a tree organism
type Tree struct {
	Organism
	LogCount  int
	LeafCount int
	HasFruit  bool
}

// Bush represents a bush organism
type Bush struct {
	Organism
	LeafType   string // "oak", "birch", etc.
	BerryCount int
}

// Flower represents a flower organism
type Flower struct {
	Organism
	FlowerType string // "red", "yellow", "blue", etc.
}

// Mushroom represents a hostile mushroom organism
type Mushroom struct {
	Organism
	SporeDamage    float64
	PoisonDuration int64
}

// VenusFlytrap represents a hostile carnivorous plant
type VenusFlytrap struct {
	Organism
	BiteDamage float64
	LureRange  float64
}

// CreateOrganism creates a new organism based on type
func CreateOrganism(typeStr string, x, y float64, hex hexagon.Hexagon) *Organism {
	props := OrganismDefinitions[typeStr]
	switch typeStr {
	case "tree":
		return createTree(x, y, hex, props)
	case "bush":
		return createBush(x, y, hex, props)
	case "flower":
		return createFlower(x, y, hex, props)
	case "mushroom":
		return createMushroom(x, y, hex, props)
	case "venus_flytrap":
		return createVenusFlytrap(x, y, hex, props)
	default:
		return nil
	}
}

// createTree creates a new tree
func createTree(x, y float64, hex hexagon.Hexagon, props *OrganismJSON) *Organism {
	maxHealth := props.Properties["maxHealth"].(float64)
	isHostile := props.Properties["isHostile"].(bool)
	damage := props.Properties["damage"].(float64)
	attackRange := props.Properties["attackRange"].(float64)
	logCount := int(props.Behavior["logCount"].(float64))
	leafCount := int(props.Behavior["leafCount"].(float64))
	hasFruit := props.Behavior["hasFruit"].(bool)
	tree := &Tree{
		Organism: Organism{
			ID:          generateID(),
			Type:        TREE,
			TypeString:  props.ID,
			X:           x,
			Y:           y,
			Hex:         hex,
			Health:      maxHealth,
			MaxHealth:   maxHealth,
			IsHostile:   isHostile,
			Damage:      damage,
			AttackRange: attackRange,
		},
		LogCount:  logCount,
		LeafCount: leafCount,
		HasFruit:  hasFruit,
	}
	return &tree.Organism
}

// createBush creates a new bush
func createBush(x, y float64, hex hexagon.Hexagon, props *OrganismJSON) *Organism {
	maxHealth := props.Properties["maxHealth"].(float64)
	isHostile := props.Properties["isHostile"].(bool)
	damage := props.Properties["damage"].(float64)
	attackRange := props.Properties["attackRange"].(float64)
	leafType := props.Behavior["leafType"].(string)
	berryCount := int(props.Behavior["berryCount"].(float64))
	bush := &Bush{
		Organism: Organism{
			ID:          generateID(),
			Type:        BUSH,
			TypeString:  props.ID,
			X:           x,
			Y:           y,
			Hex:         hex,
			Health:      maxHealth,
			MaxHealth:   maxHealth,
			IsHostile:   isHostile,
			Damage:      damage,
			AttackRange: attackRange,
		},
		LeafType:   leafType,
		BerryCount: berryCount,
	}
	return &bush.Organism
}

// createFlower creates a new flower
func createFlower(x, y float64, hex hexagon.Hexagon, props *OrganismJSON) *Organism {
	maxHealth := props.Properties["maxHealth"].(float64)
	isHostile := props.Properties["isHostile"].(bool)
	damage := props.Properties["damage"].(float64)
	attackRange := props.Properties["attackRange"].(float64)
	flowerType := props.Behavior["flowerType"].(string)
	flower := &Flower{
		Organism: Organism{
			ID:          generateID(),
			Type:        FLOWER,
			TypeString:  props.ID,
			X:           x,
			Y:           y,
			Hex:         hex,
			Health:      maxHealth,
			MaxHealth:   maxHealth,
			IsHostile:   isHostile,
			Damage:      damage,
			AttackRange: attackRange,
		},
		FlowerType: flowerType,
	}
	return &flower.Organism
}

// createMushroom creates a new hostile mushroom
func createMushroom(x, y float64, hex hexagon.Hexagon, props *OrganismJSON) *Organism {
	maxHealth := props.Properties["maxHealth"].(float64)
	isHostile := props.Properties["isHostile"].(bool)
	damage := props.Properties["damage"].(float64)
	attackRange := props.Properties["attackRange"].(float64)
	sporeDamage := props.Properties["sporeDamage"].(float64)
	poisonDuration := int64(props.Properties["poisonDuration"].(float64))
	mushroom := &Mushroom{
		Organism: Organism{
			ID:          generateID(),
			Type:        MUSHROOM,
			TypeString:  props.ID,
			X:           x,
			Y:           y,
			Hex:         hex,
			Health:      maxHealth,
			MaxHealth:   maxHealth,
			IsHostile:   isHostile,
			Damage:      damage,
			AttackRange: attackRange,
		},
		SporeDamage:    sporeDamage,
		PoisonDuration: poisonDuration,
	}
	return &mushroom.Organism
}

// createVenusFlytrap creates a new hostile venus flytrap
func createVenusFlytrap(x, y float64, hex hexagon.Hexagon, props *OrganismJSON) *Organism {
	maxHealth := props.Properties["maxHealth"].(float64)
	isHostile := props.Properties["isHostile"].(bool)
	damage := props.Properties["damage"].(float64)
	attackRange := props.Properties["attackRange"].(float64)
	biteDamage := props.Properties["biteDamage"].(float64)
	lureRange := props.Properties["lureRange"].(float64)
	venusFlytrap := &VenusFlytrap{
		Organism: Organism{
			ID:          generateID(),
			Type:        VENUS_FLYTRAP,
			TypeString:  props.ID,
			X:           x,
			Y:           y,
			Hex:         hex,
			Health:      maxHealth,
			MaxHealth:   maxHealth,
			IsHostile:   isHostile,
			Damage:      damage,
			AttackRange: attackRange,
		},
		BiteDamage: biteDamage,
		LureRange:  lureRange,
	}
	return &venusFlytrap.Organism
}

// GetOrganismBlocks returns the blocks that make up an organism
func GetOrganismBlocks(org *Organism) []hexagon.Hexagon {
	blocks := []hexagon.Hexagon{}

	switch org.TypeString {
	case "tree":
		props := OrganismDefinitions[org.TypeString]
		logCount := int(props.Behavior["logCount"].(float64))
		// Tree trunk (vertical)
		for i := 0; i < logCount; i++ {
			hex := hexagon.AxialToHex(org.Hex.Q, org.Hex.R-i)
			blocks = append(blocks, hex)
		}
		// Leaves (around the top)
		topHex := hexagon.AxialToHex(org.Hex.Q, org.Hex.R-logCount)
		neighbors := hexagon.HexNeighbors(topHex)
		blocks = append(blocks, neighbors...)
		// More leaves layer
		topHex2 := hexagon.AxialToHex(org.Hex.Q, org.Hex.R-logCount+1)
		neighbors2 := hexagon.HexNeighbors(topHex2)
		blocks = append(blocks, neighbors2...)

	case "bush":
		// Bush is a single block
		blocks = append(blocks, org.Hex)

	case "flower":
		// Flower is a single block
		blocks = append(blocks, org.Hex)
	}

	return blocks
}

// TakeDamage damages an organism
func (org *Organism) TakeDamage(amount float64) bool {
	org.Health -= amount
	if org.Health <= 0 {
		org.Health = 0
		return true // Organism destroyed
	}
	return false
}

// Heal heals an organism
func (org *Organism) Heal(amount float64) {
	org.Health += amount
	if org.Health > org.MaxHealth {
		org.Health = org.MaxHealth
	}
}

// IsAlive returns true if the organism is still alive
func (org *Organism) IsAlive() bool {
	return org.Health > 0
}

// GetPosition returns the organism's position
func (org *Organism) GetPosition() (float64, float64) {
	return org.X, org.Y
}

// GetHex returns the organism's hex position
func (org *Organism) GetHex() hexagon.Hexagon {
	return org.Hex
}

// generateID generates a unique ID for an organism
func generateID() string {
	// Simple ID generation - in production use UUID
	return "org_" + randomString(8)
}

// randomString generates a random string of the given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)] // Simplified for deterministic output
	}
	return string(result)
}

// GetDrops returns the items dropped when an organism is destroyed
func GetDrops(org *Organism) []string {
	if props, ok := OrganismDefinitions[org.TypeString]; ok {
		return props.Drops
	}
	return []string{}
}

// CanAttack checks if an organism can attack based on cooldown
func (org *Organism) CanAttack() bool {
	if !org.IsHostile {
		return false
	}

	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	attackCooldown := int64(2000) // 2 seconds cooldown

	return currentTime-org.LastAttackTime >= attackCooldown
}

// AttackPlayer attempts to attack the player if in range
func (org *Organism) AttackPlayer(playerX, playerY float64, dealDamage func(damage float64, fromX, fromY float64, knockback float64)) bool {
	if !org.IsHostile || !org.CanAttack() {
		return false
	}

	// Calculate distance to player
	dx := playerX - org.X
	dy := playerY - org.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	// Check if in attack range
	if distance <= org.AttackRange {
		// Attack the player
		dealDamage(org.Damage, org.X, org.Y, 30.0) // Damage, from position, knockback
		org.LastAttackTime = time.Now().UnixNano() / int64(time.Millisecond)
		return true
	}

	return false
}

// UpdateCombat updates organism combat behavior
func (org *Organism) UpdateCombat(playerX, playerY float64, deltaTime float64, dealDamage func(damage float64, fromX, fromY float64, knockback float64)) {
	if !org.IsHostile {
		return
	}

	// Try to attack player if in range
	org.AttackPlayer(playerX, playerY, dealDamage)
}
