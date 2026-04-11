// Package health implements a locational health system where different body parts
// have their own health pools with visual representation.
package health

import (
	"image/color"
	"time"
)

// BodyPart represents different parts of the body
type BodyPart int

const (
	PartHead BodyPart = iota
	PartTorso
	PartLeftArm
	PartRightArm
	PartLeftLeg
	PartRightLeg
	PartCount // Total number of body parts
)

// DamageType represents different types of damage
type DamageType int

const (
	DamagePhysical DamageType = iota
	DamageFire
	DamageMagic
	DamagePoison
	DamageFall
	DamageExplosion
)

// BodyPartHealth tracks health for a specific body part
type BodyPartHealth struct {
	Name           string
	Health         float64
	MaxHealth      float64
	ArmorCoverage  float64 // 0-1, how well protected this part is
	IsVital        bool    // If this part reaches 0, player dies
	LastDamageTime time.Time
	DamageEffects  []DamageEffect
}

// DamageEffect represents active damage effects on a body part
type DamageEffect struct {
	Type      DamageType
	Amount    float64
	Duration  time.Duration
	StartTime time.Time
}

// LocationalHealthSystem manages health for all body parts
type LocationalHealthSystem struct {
	Parts           [PartCount]*BodyPartHealth
	OverallHealth   float64 // Calculated from parts
	MaxOverallHealth float64
	RegenEnabled    bool
	RegenRate       float64
	LastDamageTime  time.Time
	RegenDelay      time.Duration
}

// NewLocationalHealthSystem creates a new health system
func NewLocationalHealthSystem() *LocationalHealthSystem {
	system := &LocationalHealthSystem{
		MaxOverallHealth: 100,
		RegenEnabled:     true,
		RegenRate:        0.25, // Health per second
		RegenDelay:       10 * time.Second,
		LastDamageTime:   time.Now(),
	}

	// Initialize body parts
	system.Parts[PartHead] = &BodyPartHealth{
		Name:          "Head",
		Health:        20,
		MaxHealth:     20,
		ArmorCoverage: 0.15,
		IsVital:       true,
	}
	system.Parts[PartTorso] = &BodyPartHealth{
		Name:          "Torso",
		Health:        40,
		MaxHealth:     40,
		ArmorCoverage: 0.35,
		IsVital:       true,
	}
	system.Parts[PartLeftArm] = &BodyPartHealth{
		Name:          "Left Arm",
		Health:        15,
		MaxHealth:     15,
		ArmorCoverage: 0.10,
		IsVital:       false,
	}
	system.Parts[PartRightArm] = &BodyPartHealth{
		Name:          "Right Arm",
		Health:        15,
		MaxHealth:     15,
		ArmorCoverage: 0.10,
		IsVital:       false,
	}
	system.Parts[PartLeftLeg] = &BodyPartHealth{
		Name:          "Left Leg",
		Health:        25,
		MaxHealth:     25,
		ArmorCoverage: 0.15,
		IsVital:       false,
	}
	system.Parts[PartRightLeg] = &BodyPartHealth{
		Name:          "Right Leg",
		Health:        25,
		MaxHealth:     25,
		ArmorCoverage: 0.15,
		IsVital:       false,
	}

	system.calculateOverallHealth()
	return system
}

// DamageBodyPart applies damage to a specific body part
func (lhs *LocationalHealthSystem) DamageBodyPart(part BodyPart, amount float64, damageType DamageType) float64 {
	if part < 0 || part >= PartCount {
		return 0
	}

	bodyPart := lhs.Parts[part]
	actualDamage := amount

	// Apply armor coverage reduction
	if bodyPart.ArmorCoverage > 0 {
		reduction := bodyPart.ArmorCoverage * 0.8 // Max 80% reduction from armor
		actualDamage = amount * (1 - reduction)
	}

	// Apply damage
	bodyPart.Health -= actualDamage
	if bodyPart.Health < 0 {
		bodyPart.Health = 0
	}

	bodyPart.LastDamageTime = time.Now()
	lhs.LastDamageTime = time.Now()

	// Add damage effect
	effect := DamageEffect{
		Type:      damageType,
		Amount:    actualDamage,
		Duration:  5 * time.Second,
		StartTime: time.Now(),
	}
	bodyPart.DamageEffects = append(bodyPart.DamageEffects, effect)

	lhs.calculateOverallHealth()
	return actualDamage
}

// HealBodyPart heals a specific body part
func (lhs *LocationalHealthSystem) HealBodyPart(part BodyPart, amount float64) {
	if part < 0 || part >= PartCount {
		return
	}

	bodyPart := lhs.Parts[part]
	bodyPart.Health += amount
	if bodyPart.Health > bodyPart.MaxHealth {
		bodyPart.Health = bodyPart.MaxHealth
	}

	lhs.calculateOverallHealth()
}

// Update updates the health system (regeneration, effects)
func (lhs *LocationalHealthSystem) Update(deltaTime float64) {
	if !lhs.RegenEnabled {
		return
	}

	// Check if enough time has passed since last damage
	if time.Since(lhs.LastDamageTime) < lhs.RegenDelay {
		return
	}

	// Regenerate all body parts
	for _, part := range lhs.Parts {
		if part != nil && part.Health < part.MaxHealth {
			part.Health += lhs.RegenRate * deltaTime
			if part.Health > part.MaxHealth {
				part.Health = part.MaxHealth
			}
		}
	}

	// Clean up expired damage effects
	for _, part := range lhs.Parts {
		if part != nil {
			activeEffects := []DamageEffect{}
			for _, effect := range part.DamageEffects {
				if time.Since(effect.StartTime) < effect.Duration {
					activeEffects = append(activeEffects, effect)
				}
			}
			part.DamageEffects = activeEffects
		}
	}

	lhs.calculateOverallHealth()
}

// calculateOverallHealth recalculates overall health from body parts
func (lhs *LocationalHealthSystem) calculateOverallHealth() {
	totalHealth := 0.0
	totalMax := 0.0

	for _, part := range lhs.Parts {
		if part != nil {
			totalHealth += part.Health
			totalMax += part.MaxHealth
		}
	}

	lhs.OverallHealth = totalHealth
	lhs.MaxOverallHealth = totalMax
}

// IsAlive returns true if player is alive
func (lhs *LocationalHealthSystem) IsAlive() bool {
	// Check if any vital parts are at 0
	for _, part := range lhs.Parts {
		if part != nil && part.IsVital && part.Health <= 0 {
			return false
		}
	}
	return lhs.OverallHealth > 0
}

// GetPartHealthPercentage returns health percentage for a body part
func (lhs *LocationalHealthSystem) GetPartHealthPercentage(part BodyPart) float64 {
	if part < 0 || part >= PartCount || lhs.Parts[part] == nil {
		return 0
	}
	p := lhs.Parts[part]
	if p.MaxHealth <= 0 {
		return 0
	}
	return p.Health / p.MaxHealth
}

// GetOverallHealthPercentage returns overall health percentage
func (lhs *LocationalHealthSystem) GetOverallHealthPercentage() float64 {
	if lhs.MaxOverallHealth <= 0 {
		return 0
	}
	return lhs.OverallHealth / lhs.MaxOverallHealth
}

// GetStatusEffects returns all active status effects
func (lhs *LocationalHealthSystem) GetStatusEffects() []DamageEffect {
	allEffects := []DamageEffect{}
	for _, part := range lhs.Parts {
		if part != nil {
			allEffects = append(allEffects, part.DamageEffects...)
		}
	}
	return allEffects
}

// SetArmorCoverage updates armor coverage for a body part
func (lhs *LocationalHealthSystem) SetArmorCoverage(part BodyPart, coverage float64) {
	if part < 0 || part >= PartCount || lhs.Parts[part] == nil {
		return
	}
	if coverage < 0 {
		coverage = 0
	}
	if coverage > 1 {
		coverage = 1
	}
	lhs.Parts[part].ArmorCoverage = coverage
}

// BodyPartPosition returns the visual position for a body part in the UI
func BodyPartPosition(part BodyPart) (x, y float64, width, height float64) {
	// Returns positions for a body diagram (140x200 rectangle)
	switch part {
	case PartHead:
		return 50, 10, 40, 40 // Head at top center
	case PartTorso:
		return 40, 55, 60, 70 // Torso below head
	case PartLeftArm:
		return 10, 55, 25, 60 // Left arm
	case PartRightArm:
		return 105, 55, 25, 60 // Right arm
	case PartLeftLeg:
		return 45, 130, 25, 60 // Left leg
	case PartRightLeg:
		return 70, 130, 25, 60 // Right leg
	default:
		return 0, 0, 0, 0
	}
}

// GetBodyPartColor returns the color for a body part based on health
func GetBodyPartColor(healthPercentage float64) color.RGBA {
	if healthPercentage > 0.7 {
		return color.RGBA{100, 255, 100, 255} // Green - healthy
	} else if healthPercentage > 0.4 {
		return color.RGBA{255, 255, 100, 255} // Yellow - injured
	} else if healthPercentage > 0.2 {
		return color.RGBA{255, 150, 50, 255} // Orange - badly injured
	} else {
		return color.RGBA{255, 50, 50, 255} // Red - critical
	}
}

// GetPartName returns the display name for a body part
func GetPartName(part BodyPart) string {
	names := []string{"Head", "Torso", "Left Arm", "Right Arm", "Left Leg", "Right Leg"}
	if part < 0 || part >= PartCount {
		return "Unknown"
	}
	return names[part]
}

// IsPartInjured returns true if body part is significantly damaged
func (lhs *LocationalHealthSystem) IsPartInjured(part BodyPart) bool {
	return lhs.GetPartHealthPercentage(part) < 0.5
}

// GetMostDamagedPart returns the most damaged body part
func (lhs *LocationalHealthSystem) GetMostDamagedPart() (BodyPart, float64) {
	mostDamaged := PartHead
	lowestHealth := 1.0

	for i, part := range lhs.Parts {
		if part != nil {
			health := lhs.GetPartHealthPercentage(BodyPart(i))
			if health < lowestHealth {
				lowestHealth = health
				mostDamaged = BodyPart(i)
			}
		}
	}

	return mostDamaged, lowestHealth
}
