// Package equipment implements the sophisticated clothing and armor system
// with material types, layering, and special equipment like wings.
package equipment

import (
	"fmt"
	"image/color"
)

// Ensure color package is properly imported

// EquipmentSlot represents different equipment slots
type EquipmentSlot int

const (
	SlotHelmet EquipmentSlot = iota
	SlotChestplate
	SlotLeggings
	SlotBoots
	SlotWings
	SlotGloves
	SlotCloak
	SlotAmulet
	SlotRing1
	SlotRing2
	SlotCount // Total number of slots
)

// ArmorMaterial represents different armor materials
type ArmorMaterial int

const (
	MaterialCloth ArmorMaterial = iota
	MaterialLeather
	MaterialChain
	MaterialIron
	MaterialGold
	MaterialDiamond
	MaterialNetherite
	MaterialDragon
)

// ArmorType represents the type of armor piece
type ArmorType int

const (
	ArmorLight ArmorType = iota
	ArmorMedium
	ArmorHeavy
)

// EquipmentItem represents a piece of equipment
type EquipmentItem struct {
	ID          string
	Name        string
	Description string
	Slot        EquipmentSlot
	Material    ArmorMaterial
	ArmorType   ArmorType
	IconColor   color.RGBA

	// Defense stats
	BaseDefense   float64 // Base damage reduction
	Durability    int
	MaxDurability int

	// Special properties
	FireResistant        bool
	ProjectileProof      bool
	BlastResistant       bool
	MagicResistant       bool
	GrantsFlight         bool // For wings
	GrantsWaterBreathing bool
	GrantsNightVision    bool

	// Movement modifiers
	MovementSpeedMod float64 // Multiplier (1.0 = normal)
	JumpHeightMod    float64 // Multiplier
	FallDamageMod    float64 // Damage multiplier (0 = no fall damage)

	// Set bonus
	SetName        string // Name of the armor set
	SetPieces      int    // Number of pieces needed for bonus
	SetBonusActive bool   // Whether set bonus is currently active
}

// EquipmentSet manages all equipped items
type EquipmentSet struct {
	Slots            [SlotCount]*EquipmentItem
	ActiveSetBonuses map[string]bool
}

// NewEquipmentSet creates a new empty equipment set
func NewEquipmentSet() *EquipmentSet {
	return &EquipmentSet{
		Slots:            [SlotCount]*EquipmentItem{},
		ActiveSetBonuses: make(map[string]bool),
	}
}

// EquipItem equips an item to a slot
func (es *EquipmentSet) EquipItem(item *EquipmentItem, slot EquipmentSlot) *EquipmentItem {
	if slot < 0 || slot >= SlotCount {
		return nil
	}

	oldItem := es.Slots[slot]
	es.Slots[slot] = item

	// Recalculate set bonuses
	es.recalculateSetBonuses()

	return oldItem // Return previously equipped item
}

// UnequipItem removes an item from a slot
func (es *EquipmentSet) UnequipItem(slot EquipmentSlot) *EquipmentItem {
	if slot < 0 || slot >= SlotCount {
		return nil
	}

	item := es.Slots[slot]
	es.Slots[slot] = nil

	// Recalculate set bonuses
	es.recalculateSetBonuses()

	return item
}

// GetItem returns the item in a slot
func (es *EquipmentSet) GetItem(slot EquipmentSlot) *EquipmentItem {
	if slot < 0 || slot >= SlotCount {
		return nil
	}
	return es.Slots[slot]
}

// GetTotalDefense calculates total defense from all equipped armor
func (es *EquipmentSet) GetTotalDefense() float64 {
	total := 0.0
	for _, item := range es.Slots {
		if item != nil && item.Slot != SlotWings {
			total += item.BaseDefense
		}
	}
	return total
}

// GetDamageReduction calculates damage reduction (0-1) based on total defense
func (es *EquipmentSet) GetDamageReduction() float64 {
	defense := es.GetTotalDefense()
	// Diminishing returns formula
	reduction := defense / (defense + 100.0)
	if reduction > 0.8 {
		reduction = 0.8 // Cap at 80% reduction
	}
	return reduction
}

// CanFly returns true if player has equipped wings that grant flight
func (es *EquipmentSet) CanFly() bool {
	wings := es.Slots[SlotWings]
	return wings != nil && wings.GrantsFlight
}

// GetMovementModifiers returns combined movement modifiers
func (es *EquipmentSet) GetMovementModifiers() (speed, jump, fallDamage float64) {
	speed = 1.0
	jump = 1.0
	fallDamage = 1.0

	for _, item := range es.Slots {
		if item != nil {
			speed *= item.MovementSpeedMod
			jump *= item.JumpHeightMod
			fallDamage *= item.FallDamageMod
		}
	}

	return speed, jump, fallDamage
}

// recalculateSetBonuses recalculates which set bonuses are active
func (es *EquipmentSet) recalculateSetBonuses() {
	// Clear existing bonuses
	es.ActiveSetBonuses = make(map[string]bool)

	// Count pieces per set
	setCounts := make(map[string]int)
	for _, item := range es.Slots {
		if item != nil && item.SetName != "" {
			setCounts[item.SetName]++
		}
	}

	// Check for active set bonuses
	for setName, count := range setCounts {
		// Find how many pieces are needed for bonus
		neededPieces := 4 // Default to full set
		for _, item := range es.Slots {
			if item != nil && item.SetName == setName && item.SetPieces > 0 {
				neededPieces = item.SetPieces
				break
			}
		}

		if count >= neededPieces {
			es.ActiveSetBonuses[setName] = true
		}
	}
}

// HasSetBonus returns true if a specific set bonus is active
func (es *EquipmentSet) HasSetBonus(setName string) bool {
	return es.ActiveSetBonuses[setName]
}

// TakeDurabilityDamage reduces durability on all armor when taking damage
func (es *EquipmentSet) TakeDurabilityDamage(amount int) {
	for _, item := range es.Slots {
		if item != nil && item.Durability > 0 {
			item.Durability -= amount
			if item.Durability < 0 {
				item.Durability = 0
			}
		}
	}
}

// GetEquippedCount returns number of equipped items
func (es *EquipmentSet) GetEquippedCount() int {
	count := 0
	for _, item := range es.Slots {
		if item != nil {
			count++
		}
	}
	return count
}

// MaterialProperties returns display properties for materials
func MaterialProperties(material ArmorMaterial) (name string, col color.RGBA, tier int) {
	switch material {
	case MaterialCloth:
		return "Cloth", color.RGBA{200, 180, 160, 255}, 1
	case MaterialLeather:
		return "Leather", color.RGBA{139, 90, 43, 255}, 2
	case MaterialChain:
		return "Chainmail", color.RGBA{169, 169, 169, 255}, 3
	case MaterialIron:
		return "Iron", color.RGBA{192, 192, 192, 255}, 4
	case MaterialGold:
		return "Gold", color.RGBA{255, 215, 0, 255}, 4
	case MaterialDiamond:
		return "Diamond", color.RGBA{0, 255, 255, 255}, 5
	case MaterialNetherite:
		return "Netherite", color.RGBA{50, 50, 50, 255}, 6
	case MaterialDragon:
		return "Dragon", color.RGBA{200, 50, 50, 255}, 7
	default:
		return "Unknown", color.RGBA{128, 128, 128, 255}, 0
	}
}

// CreateArmor creates a new armor piece with material-based stats
func CreateArmor(name string, slot EquipmentSlot, material ArmorMaterial, armorType ArmorType) *EquipmentItem {
	item := &EquipmentItem{
		ID:               fmt.Sprintf("%s_%d_%d", name, slot, material),
		Name:             name,
		Slot:             slot,
		Material:         material,
		ArmorType:        armorType,
		MaxDurability:    getMaterialDurability(material),
		Durability:       getMaterialDurability(material),
		MovementSpeedMod: 1.0,
		JumpHeightMod:    1.0,
		FallDamageMod:    1.0,
	}

	// Calculate base defense based on material, slot, and type
	item.BaseDefense = calculateDefense(material, slot, armorType)
	_, item.IconColor, _ = MaterialProperties(material)

	// Apply material special properties
	applyMaterialProperties(item, material)

	return item
}

// CreateWings creates wings that grant flight
func CreateWings(name string, material ArmorMaterial) *EquipmentItem {
	wings := &EquipmentItem{
		ID:               fmt.Sprintf("wings_%s_%d", name, material),
		Name:             name + " Wings",
		Description:      "Wings that allow flight",
		Slot:             SlotWings,
		Material:         material,
		ArmorType:        ArmorLight,
		GrantsFlight:     true,
		MaxDurability:    getMaterialDurability(material) / 2, // Wings are fragile
		Durability:       getMaterialDurability(material) / 2,
		MovementSpeedMod: 1.1, // Slight speed boost
		FallDamageMod:    0.5, // Reduced fall damage
	}

	// Material-specific wing properties
	switch material {
	case MaterialCloth:
		wings.BaseDefense = 1
		wings.IconColor = color.RGBA{255, 255, 255, 255}
	case MaterialLeather:
		wings.BaseDefense = 2
		wings.IconColor = color.RGBA{139, 90, 43, 255}
	case MaterialDragon:
		wings.BaseDefense = 8
		wings.FireResistant = true
		wings.IconColor = color.RGBA{200, 50, 50, 255}
	default:
		wings.BaseDefense = 2
	}

	return wings
}

// Helper functions
func getMaterialDurability(material ArmorMaterial) int {
	switch material {
	case MaterialCloth:
		return 50
	case MaterialLeather:
		return 80
	case MaterialChain:
		return 120
	case MaterialIron:
		return 165
	case MaterialGold:
		return 112
	case MaterialDiamond:
		return 363
	case MaterialNetherite:
		return 500
	case MaterialDragon:
		return 800
	default:
		return 50
	}
}

func calculateDefense(material ArmorMaterial, slot EquipmentSlot, armorType ArmorType) float64 {
	// Base defense by material
	materialDefense := map[ArmorMaterial]float64{
		MaterialCloth:     2,
		MaterialLeather:   4,
		MaterialChain:     6,
		MaterialIron:      8,
		MaterialGold:      6,
		MaterialDiamond:   12,
		MaterialNetherite: 16,
		MaterialDragon:    20,
	}

	// Slot multiplier (chestplate gives more defense than boots)
	slotMultiplier := map[EquipmentSlot]float64{
		SlotHelmet:     0.8,
		SlotChestplate: 1.2,
		SlotLeggings:   1.0,
		SlotBoots:      0.6,
		SlotGloves:     0.3,
		SlotCloak:      0.4,
	}

	// Type multiplier
	typeMultiplier := map[ArmorType]float64{
		ArmorLight:  0.8,
		ArmorMedium: 1.0,
		ArmorHeavy:  1.3,
	}

	base := materialDefense[material]
	mult := slotMultiplier[slot] * typeMultiplier[armorType]

	return base * mult
}

func applyMaterialProperties(item *EquipmentItem, material ArmorMaterial) {
	switch material {
	case MaterialLeather:
		item.MovementSpeedMod = 1.05 // Leather is light
	case MaterialChain:
		// Chain has no movement penalty
	case MaterialIron:
		item.MovementSpeedMod = 0.95
	case MaterialGold:
		item.MovementSpeedMod = 1.0
		item.MagicResistant = true // Gold has magic resistance
	case MaterialDiamond:
		item.MovementSpeedMod = 0.90
		item.ProjectileProof = true
	case MaterialNetherite:
		item.MovementSpeedMod = 0.85
		item.FireResistant = true
		item.BlastResistant = true
		item.ProjectileProof = true
	case MaterialDragon:
		item.MovementSpeedMod = 1.0 // Dragon scales are light but strong
		item.FireResistant = true
		item.ProjectileProof = true
		item.BlastResistant = true
		item.GrantsNightVision = true
	}

	if item.ArmorType == ArmorHeavy {
		item.MovementSpeedMod *= 0.9 // Heavy armor slows more
	}
}
