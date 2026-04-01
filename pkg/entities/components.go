package entities

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// ============================================================================
// Core Component Types
// ============================================================================

// RenderComponent handles visual properties
type RenderComponent struct {
	Type        string       `yaml:"type"`
	Color       color.RGBA   `yaml:"color"`
	TopColor    color.RGBA   `yaml:"topColor,omitempty"`
	SideColor   color.RGBA   `yaml:"sideColor,omitempty"`
	Colors      []color.RGBA `yaml:"colors,omitempty"`
	Pattern     string       `yaml:"pattern"`
	Texture     string       `yaml:"texture,omitempty"`
	Visible     bool         `yaml:"visible"`
	LightLevel  int          `yaml:"lightLevel"`
	Transparent bool         `yaml:"transparent"`
	Scale       float64      `yaml:"scale"`
	Animated    bool         `yaml:"animated"`
}

func (c *RenderComponent) GetType() string { return "render" }
func (c *RenderComponent) Clone() Component {
	clone := *c
	if len(c.Colors) > 0 {
		clone.Colors = make([]color.RGBA, len(c.Colors))
		copy(clone.Colors, c.Colors)
	}
	return &clone
}
func (c *RenderComponent) Merge(other Component) {
	if rc, ok := other.(*RenderComponent); ok {
		if rc.Color != (color.RGBA{0, 0, 0, 0}) {
			c.Color = rc.Color
		}
		if rc.TopColor != (color.RGBA{0, 0, 0, 0}) {
			c.TopColor = rc.TopColor
		}
		if rc.SideColor != (color.RGBA{0, 0, 0, 0}) {
			c.SideColor = rc.SideColor
		}
		if len(rc.Colors) > 0 {
			c.Colors = rc.Colors
		}
		if rc.Pattern != "" {
			c.Pattern = rc.Pattern
		}
		if rc.Texture != "" {
			c.Texture = rc.Texture
		}
	}
}
func (c *RenderComponent) Validate() error {
	if c.LightLevel < 0 || c.LightLevel > 15 {
		return fmt.Errorf("light level must be between 0 and 15")
	}
	if c.Scale <= 0 {
		return fmt.Errorf("scale must be positive")
	}
	return nil
}

// PhysicsComponent handles physical properties
type PhysicsComponent struct {
	Type       string  `yaml:"type"`
	Hardness   float64 `yaml:"hardness"`
	Density    float64 `yaml:"density"`
	Solid      bool    `yaml:"solid"`
	Gravity    bool    `yaml:"gravity"`
	Viscosity  float64 `yaml:"viscosity"`
	Friction   float64 `yaml:"friction"`
	Bounciness float64 `yaml:"bounciness"`
	Collision  bool    `yaml:"collision"`
	Mass       float64 `yaml:"mass"`
}

func (c *PhysicsComponent) GetType() string { return "physics" }
func (c *PhysicsComponent) Clone() Component {
	clone := *c
	return &clone
}
func (c *PhysicsComponent) Merge(other Component) {
	if pc, ok := other.(*PhysicsComponent); ok {
		if pc.Hardness != 0 {
			c.Hardness = pc.Hardness
		}
		if pc.Density != 0 {
			c.Density = pc.Density
		}
		c.Solid = pc.Solid
		c.Gravity = pc.Gravity
		if pc.Viscosity != 0 {
			c.Viscosity = pc.Viscosity
		}
		if pc.Friction != 0 {
			c.Friction = pc.Friction
		}
		if pc.Bounciness != 0 {
			c.Bounciness = pc.Bounciness
		}
		c.Collision = pc.Collision
		if pc.Mass != 0 {
			c.Mass = pc.Mass
		}
	}
}
func (c *PhysicsComponent) Validate() error {
	if c.Hardness < 0 {
		return fmt.Errorf("hardness cannot be negative")
	}
	if c.Density < 0 {
		return fmt.Errorf("density cannot be negative")
	}
	if c.Viscosity < 0 {
		return fmt.Errorf("viscosity cannot be negative")
	}
	if c.Friction < 0 || c.Friction > 1 {
		return fmt.Errorf("friction must be between 0 and 1")
	}
	if c.Bounciness < 0 || c.Bounciness > 1 {
		return fmt.Errorf("bounciness must be between 0 and 1")
	}
	if c.Mass < 0 {
		return fmt.Errorf("mass cannot be negative")
	}
	return nil
}

// InventoryComponent handles inventory and stacking
type InventoryComponent struct {
	Type              string                 `yaml:"type"`
	StackSize         int                    `yaml:"stackSize"`
	MaxDurability     int                    `yaml:"maxDurability"`
	CurrentDurability int                    `yaml:"currentDurability"`
	Container         bool                   `yaml:"container"`
	Slots             int                    `yaml:"slots"`
	Contents          map[string]int         `yaml:"contents,omitempty"`
	Weight            float64                `yaml:"weight"`
	Categories        []string               `yaml:"categories"`
	Properties        map[string]interface{} `yaml:"properties,omitempty"`
}

func (c *InventoryComponent) GetType() string { return "inventory" }
func (c *InventoryComponent) Clone() Component {
	clone := *c
	if len(c.Contents) > 0 {
		clone.Contents = make(map[string]int)
		for k, v := range c.Contents {
			clone.Contents[k] = v
		}
	}
	if len(c.Categories) > 0 {
		clone.Categories = make([]string, len(c.Categories))
		copy(clone.Categories, c.Categories)
	}
	if len(c.Properties) > 0 {
		clone.Properties = make(map[string]interface{})
		for k, v := range c.Properties {
			clone.Properties[k] = v
		}
	}
	return &clone
}
func (c *InventoryComponent) Merge(other Component) {
	if ic, ok := other.(*InventoryComponent); ok {
		if ic.StackSize != 0 {
			c.StackSize = ic.StackSize
		}
		if ic.MaxDurability != 0 {
			c.MaxDurability = ic.MaxDurability
		}
		if ic.CurrentDurability != 0 {
			c.CurrentDurability = ic.CurrentDurability
		}
		c.Container = ic.Container
		if ic.Slots != 0 {
			c.Slots = ic.Slots
		}
		if len(ic.Contents) > 0 {
			if c.Contents == nil {
				c.Contents = make(map[string]int)
			}
			for k, v := range ic.Contents {
				c.Contents[k] = v
			}
		}
		if ic.Weight != 0 {
			c.Weight = ic.Weight
		}
		if len(ic.Categories) > 0 {
			c.Categories = ic.Categories
		}
	}
}
func (c *InventoryComponent) Validate() error {
	if c.StackSize <= 0 {
		return fmt.Errorf("stack size must be positive")
	}
	if c.MaxDurability < 0 {
		return fmt.Errorf("max durability cannot be negative")
	}
	if c.CurrentDurability < 0 {
		return fmt.Errorf("current durability cannot be negative")
	}
	if c.Slots < 0 {
		return fmt.Errorf("slots cannot be negative")
	}
	if c.Weight < 0 {
		return fmt.Errorf("weight cannot be negative")
	}
	return nil
}

// BehaviorComponent handles AI and interactions
type BehaviorComponent struct {
	Type         string                 `yaml:"type"`
	AIType       string                 `yaml:"aiType"`
	Passive      bool                   `yaml:"passive"`
	Hostile      bool                   `yaml:"hostile"`
	Neutral      bool                   `yaml:"neutral"`
	BehaviorTree string                 `yaml:"behaviorTree,omitempty"`
	States       []string               `yaml:"states"`
	CurrentState string                 `yaml:"currentState"`
	SightRange   float64                `yaml:"sightRange"`
	HearingRange float64                `yaml:"hearingRange"`
	Speed        float64                `yaml:"speed"`
	JumpHeight   float64                `yaml:"jumpHeight"`
	Abilities    []string               `yaml:"abilities"`
	Properties   map[string]interface{} `yaml:"properties,omitempty"`
}

func (c *BehaviorComponent) GetType() string { return "behavior" }
func (c *BehaviorComponent) Clone() Component {
	clone := *c
	if len(c.States) > 0 {
		clone.States = make([]string, len(c.States))
		copy(clone.States, c.States)
	}
	if len(c.Abilities) > 0 {
		clone.Abilities = make([]string, len(c.Abilities))
		copy(clone.Abilities, c.Abilities)
	}
	if len(c.Properties) > 0 {
		clone.Properties = make(map[string]interface{})
		for k, v := range c.Properties {
			clone.Properties[k] = v
		}
	}
	return &clone
}
func (c *BehaviorComponent) Merge(other Component) {
	if bc, ok := other.(*BehaviorComponent); ok {
		if bc.AIType != "" {
			c.AIType = bc.AIType
		}
		c.Passive = bc.Passive
		c.Hostile = bc.Hostile
		c.Neutral = bc.Neutral
		if bc.BehaviorTree != "" {
			c.BehaviorTree = bc.BehaviorTree
		}
		if len(bc.States) > 0 {
			c.States = bc.States
		}
		if bc.CurrentState != "" {
			c.CurrentState = bc.CurrentState
		}
		if bc.SightRange != 0 {
			c.SightRange = bc.SightRange
		}
		if bc.HearingRange != 0 {
			c.HearingRange = bc.HearingRange
		}
		if bc.Speed != 0 {
			c.Speed = bc.Speed
		}
		if bc.JumpHeight != 0 {
			c.JumpHeight = bc.JumpHeight
		}
		if len(bc.Abilities) > 0 {
			c.Abilities = bc.Abilities
		}
	}
}
func (c *BehaviorComponent) Validate() error {
	if c.SightRange < 0 {
		return fmt.Errorf("sight range cannot be negative")
	}
	if c.HearingRange < 0 {
		return fmt.Errorf("hearing range cannot be negative")
	}
	if c.Speed < 0 {
		return fmt.Errorf("speed cannot be negative")
	}
	if c.JumpHeight < 0 {
		return fmt.Errorf("jump height cannot be negative")
	}
	return nil
}

// CraftingComponent handles crafting and recipes
type CraftingComponent struct {
	Type         string                 `yaml:"type"`
	Craftable    bool                   `yaml:"craftable"`
	Recipe       map[string]int         `yaml:"recipe,omitempty"`
	Results      map[string]int         `yaml:"results,omitempty"`
	CraftingTime time.Duration          `yaml:"craftingTime"`
	RequiredTool string                 `yaml:"requiredTool,omitempty"`
	Category     string                 `yaml:"category"`
	Tier         int                    `yaml:"tier"`
	Properties   map[string]interface{} `yaml:"properties,omitempty"`
}

func (c *CraftingComponent) GetType() string { return "crafting" }
func (c *CraftingComponent) Clone() Component {
	clone := *c
	if len(c.Recipe) > 0 {
		clone.Recipe = make(map[string]int)
		for k, v := range c.Recipe {
			clone.Recipe[k] = v
		}
	}
	if len(c.Results) > 0 {
		clone.Results = make(map[string]int)
		for k, v := range c.Results {
			clone.Results[k] = v
		}
	}
	if len(c.Properties) > 0 {
		clone.Properties = make(map[string]interface{})
		for k, v := range c.Properties {
			clone.Properties[k] = v
		}
	}
	return &clone
}
func (c *CraftingComponent) Merge(other Component) {
	if cc, ok := other.(*CraftingComponent); ok {
		c.Craftable = cc.Craftable
		if len(cc.Recipe) > 0 {
			if c.Recipe == nil {
				c.Recipe = make(map[string]int)
			}
			for k, v := range cc.Recipe {
				c.Recipe[k] = v
			}
		}
		if len(cc.Results) > 0 {
			if c.Results == nil {
				c.Results = make(map[string]int)
			}
			for k, v := range cc.Results {
				c.Results[k] = v
			}
		}
		if cc.CraftingTime != 0 {
			c.CraftingTime = cc.CraftingTime
		}
		if cc.RequiredTool != "" {
			c.RequiredTool = cc.RequiredTool
		}
		if cc.Category != "" {
			c.Category = cc.Category
		}
		if cc.Tier != 0 {
			c.Tier = cc.Tier
		}
	}
}
func (c *CraftingComponent) Validate() error {
	if c.CraftingTime < 0 {
		return fmt.Errorf("crafting time cannot be negative")
	}
	if c.Tier < 0 {
		return fmt.Errorf("tier cannot be negative")
	}
	return nil
}

// ToolComponent handles tool properties and effectiveness
type ToolComponent struct {
	Type          string                 `yaml:"type"`
	ToolType      string                 `yaml:"toolType"`
	Power         float64                `yaml:"power"`
	Efficiency    float64                `yaml:"efficiency"`
	Durability    int                    `yaml:"durability"`
	MaxDurability int                    `yaml:"maxDurability"`
	Effective     []string               `yaml:"effective"`
	Enchantments  []string               `yaml:"enchantments,omitempty"`
	Properties    map[string]interface{} `yaml:"properties,omitempty"`
}

func (c *ToolComponent) GetType() string { return "tool" }
func (c *ToolComponent) Clone() Component {
	clone := *c
	if len(c.Effective) > 0 {
		clone.Effective = make([]string, len(c.Effective))
		copy(clone.Effective, c.Effective)
	}
	if len(c.Enchantments) > 0 {
		clone.Enchantments = make([]string, len(c.Enchantments))
		copy(clone.Enchantments, c.Enchantments)
	}
	if len(c.Properties) > 0 {
		clone.Properties = make(map[string]interface{})
		for k, v := range c.Properties {
			clone.Properties[k] = v
		}
	}
	return &clone
}
func (c *ToolComponent) Merge(other Component) {
	if tc, ok := other.(*ToolComponent); ok {
		if tc.ToolType != "" {
			c.ToolType = tc.ToolType
		}
		if tc.Power != 0 {
			c.Power = tc.Power
		}
		if tc.Efficiency != 0 {
			c.Efficiency = tc.Efficiency
		}
		if tc.Durability != 0 {
			c.Durability = tc.Durability
		}
		if tc.MaxDurability != 0 {
			c.MaxDurability = tc.MaxDurability
		}
		if len(tc.Effective) > 0 {
			c.Effective = tc.Effective
		}
		if len(tc.Enchantments) > 0 {
			c.Enchantments = tc.Enchantments
		}
	}
}
func (c *ToolComponent) Validate() error {
	if c.Power < 0 {
		return fmt.Errorf("tool power cannot be negative")
	}
	if c.Efficiency < 0 {
		return fmt.Errorf("tool efficiency cannot be negative")
	}
	if c.Durability < 0 {
		return fmt.Errorf("tool durability cannot be negative")
	}
	if c.MaxDurability < 0 {
		return fmt.Errorf("tool max durability cannot be negative")
	}
	return nil
}

// CombatComponent handles combat properties
type CombatComponent struct {
	Type         string                 `yaml:"type"`
	WeaponType   string                 `yaml:"weaponType"`
	Damage       float64                `yaml:"damage"`
	Range        float64                `yaml:"range"`
	Speed        float64                `yaml:"speed"`
	Armor        float64                `yaml:"armor"`
	ArmorType    string                 `yaml:"armorType"`
	Health       float64                `yaml:"health"`
	MaxHealth    float64                `yaml:"maxHealth"`
	Regeneration float64                `yaml:"regeneration"`
	Effects      []string               `yaml:"effects,omitempty"`
	Properties   map[string]interface{} `yaml:"properties,omitempty"`
}

func (c *CombatComponent) GetType() string { return "combat" }
func (c *CombatComponent) Clone() Component {
	clone := *c
	if len(c.Effects) > 0 {
		clone.Effects = make([]string, len(c.Effects))
		copy(clone.Effects, c.Effects)
	}
	if len(c.Properties) > 0 {
		clone.Properties = make(map[string]interface{})
		for k, v := range c.Properties {
			clone.Properties[k] = v
		}
	}
	return &clone
}
func (c *CombatComponent) Merge(other Component) {
	if cc, ok := other.(*CombatComponent); ok {
		if cc.WeaponType != "" {
			c.WeaponType = cc.WeaponType
		}
		if cc.Damage != 0 {
			c.Damage = cc.Damage
		}
		if cc.Range != 0 {
			c.Range = cc.Range
		}
		if cc.Speed != 0 {
			c.Speed = cc.Speed
		}
		if cc.Armor != 0 {
			c.Armor = cc.Armor
		}
		if cc.ArmorType != "" {
			c.ArmorType = cc.ArmorType
		}
		if cc.Health != 0 {
			c.Health = cc.Health
		}
		if cc.MaxHealth != 0 {
			c.MaxHealth = cc.MaxHealth
		}
		if cc.Regeneration != 0 {
			c.Regeneration = cc.Regeneration
		}
		if len(cc.Effects) > 0 {
			c.Effects = cc.Effects
		}
	}
}
func (c *CombatComponent) Validate() error {
	if c.Damage < 0 {
		return fmt.Errorf("damage cannot be negative")
	}
	if c.Range < 0 {
		return fmt.Errorf("range cannot be negative")
	}
	if c.Speed < 0 {
		return fmt.Errorf("speed cannot be negative")
	}
	if c.Armor < 0 {
		return fmt.Errorf("armor cannot be negative")
	}
	if c.Health < 0 {
		return fmt.Errorf("health cannot be negative")
	}
	if c.MaxHealth < 0 {
		return fmt.Errorf("max health cannot be negative")
	}
	if c.Regeneration < 0 {
		return fmt.Errorf("regeneration cannot be negative")
	}
	return nil
}

// ============================================================================
// Component Registration
// ============================================================================

func init() {
	// Register all component types for dynamic loading
	RegisterComponent("render", &RenderComponent{})
	RegisterComponent("physics", &PhysicsComponent{})
	RegisterComponent("inventory", &InventoryComponent{})
	RegisterComponent("behavior", &BehaviorComponent{})
	RegisterComponent("crafting", &CraftingComponent{})
	RegisterComponent("tool", &ToolComponent{})
	RegisterComponent("combat", &CombatComponent{})
}

// ============================================================================
// Utility Functions
// ============================================================================

// ColorFromSlice converts a slice of uint8 to color.RGBA
func ColorFromSlice(slice []uint8) color.RGBA {
	if len(slice) >= 4 {
		return color.RGBA{slice[0], slice[1], slice[2], slice[3]}
	} else if len(slice) >= 3 {
		return color.RGBA{slice[0], slice[1], slice[2], 255}
	}
	return color.RGBA{255, 255, 255, 255}
}

// GenerateProceduralTexture creates a procedural texture from color palette
func GenerateProceduralTexture(colors []color.RGBA, seed int64) *ebiten.Image {
	if len(colors) == 0 {
		return nil
	}

	const size = 64
	img := ebiten.NewImage(size, size)

	// Use seed for deterministic generation
	rng := NewSimpleRNG(seed)

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			idx := rng.Intn(len(colors))
			img.Set(x, y, colors[idx])
		}
	}

	return img
}

// SimpleRNG is a simple random number generator for deterministic textures
type SimpleRNG struct {
	seed int64
}

func NewSimpleRNG(seed int64) *SimpleRNG {
	return &SimpleRNG{seed: seed}
}

func (r *SimpleRNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	r.seed = (r.seed*1103515245 + 12345) & 0x7fffffff
	return int(r.seed % int64(n))
}

// Clamp clamps a value between min and max
func Clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}
