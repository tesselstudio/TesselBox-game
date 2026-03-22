// status.go
package status

import (
	"time"
)

// StatusEffectType represents different types of status effects
type StatusEffectType int

const (
	POISON StatusEffectType = iota
	BLEEDING
	STRENGTH_BUFF
	SPEED_BUFF
	DEFENSE_BUFF
)

// StatusEffect represents a status effect applied to an entity
type StatusEffect struct {
	Type      StatusEffectType
	Duration  time.Duration // How long the effect lasts
	Strength  float64       // Strength of the effect
	StartTime time.Time     // When the effect was applied
}

// IsExpired returns true if the status effect has expired
func (se *StatusEffect) IsExpired() bool {
	return time.Since(se.StartTime) >= se.Duration
}

// GetRemainingTime returns the remaining duration of the effect
func (se *StatusEffect) GetRemainingTime() time.Duration {
	elapsed := time.Since(se.StartTime)
	if elapsed >= se.Duration {
		return 0
	}
	return se.Duration - elapsed
}

// StatusManager manages status effects for an entity
type StatusManager struct {
	Effects []StatusEffect
}

// NewStatusManager creates a new status manager
func NewStatusManager() *StatusManager {
	return &StatusManager{
		Effects: make([]StatusEffect, 0),
	}
}

// ApplyEffect adds or refreshes a status effect
func (sm *StatusManager) ApplyEffect(effectType StatusEffectType, duration time.Duration, strength float64) {
	// Remove existing effects of the same type
	sm.RemoveEffect(effectType)

	// Add new effect
	effect := StatusEffect{
		Type:      effectType,
		Duration:  duration,
		Strength:  strength,
		StartTime: time.Now(),
	}

	sm.Effects = append(sm.Effects, effect)
}

// RemoveEffect removes all effects of the specified type
func (sm *StatusManager) RemoveEffect(effectType StatusEffectType) {
	var newEffects []StatusEffect
	for _, effect := range sm.Effects {
		if effect.Type != effectType {
			newEffects = append(newEffects, effect)
		}
	}
	sm.Effects = newEffects
}

// Update removes expired effects
func (sm *StatusManager) Update() {
	var activeEffects []StatusEffect
	for _, effect := range sm.Effects {
		if !effect.IsExpired() {
			activeEffects = append(activeEffects, effect)
		}
	}
	sm.Effects = activeEffects
}

// HasEffect returns true if the entity has the specified effect
func (sm *StatusManager) HasEffect(effectType StatusEffectType) bool {
	for _, effect := range sm.Effects {
		if effect.Type == effectType {
			return true
		}
	}
	return false
}

// GetEffectStrength returns the strength of the specified effect (0 if not present)
func (sm *StatusManager) GetEffectStrength(effectType StatusEffectType) float64 {
	for _, effect := range sm.Effects {
		if effect.Type == effectType {
			return effect.Strength
		}
	}
	return 0
}

// ApplyPeriodicEffects applies damage/healing from periodic effects like poison/bleeding
func (sm *StatusManager) ApplyPeriodicEffects(health *float64, maxHealth float64) {
	for _, effect := range sm.Effects {
		switch effect.Type {
		case POISON:
			// Poison deals damage over time
			damage := effect.Strength
			*health -= damage
			if *health < 0 {
				*health = 0
			}
		case BLEEDING:
			// Bleeding deals damage over time
			damage := effect.Strength
			*health -= damage
			if *health < 0 {
				*health = 0
			}
		}
	}
}

// GetDamageMultiplier returns damage multiplier from buffs
func (sm *StatusManager) GetDamageMultiplier() float64 {
	multiplier := 1.0
	for _, effect := range sm.Effects {
		switch effect.Type {
		case STRENGTH_BUFF:
			multiplier += effect.Strength
		}
	}
	return multiplier
}

// GetSpeedMultiplier returns speed multiplier from buffs
func (sm *StatusManager) GetSpeedMultiplier() float64 {
	multiplier := 1.0
	for _, effect := range sm.Effects {
		switch effect.Type {
		case SPEED_BUFF:
			multiplier += effect.Strength
		}
	}
	return multiplier
}

// GetDefenseBonus returns defense bonus from buffs
func (sm *StatusManager) GetDefenseBonus() float64 {
	bonus := 0.0
	for _, effect := range sm.Effects {
		switch effect.Type {
		case DEFENSE_BUFF:
			bonus += effect.Strength
		}
	}
	return bonus
}
