// Package survival implements survival mode mechanics including game mode management,
// health systems, and survival-specific gameplay rules.
package survival

import (
	"time"

	"tesselbox/pkg/items"
	"tesselbox/pkg/player"
)

// GameMode represents the current game mode
type GameMode int

const (
	ModeCreative GameMode = iota
	ModeSurvival
	ModeHardcore
)

// SurvivalManager manages survival mode mechanics
type SurvivalManager struct {
	Mode              GameMode
	Player            *player.Player
	Inventory         *items.Inventory

	// Survival mechanics
	Hunger            float64
	MaxHunger         float64
	Thirst            float64
	MaxThirst         float64
	Stamina           float64
	MaxStamina        float64

	// Health regeneration
	HealthRegenRate   float64 // Health per second when conditions met
	LastDamageTime    time.Time
	RegenDelay        time.Duration // Time after damage before regen starts

	// Difficulty settings
	HungerDecayRate   float64 // Hunger points lost per second
	ThirstDecayRate   float64 // Thirst points lost per second

	// Status effects
	IsStarving        bool
	IsDehydrated      bool
	CanRegenerate     bool
}

// NewSurvivalManager creates a new survival manager
func NewSurvivalManager(mode GameMode, p *player.Player, inv *items.Inventory) *SurvivalManager {
	sm := &SurvivalManager{
		Mode:            mode,
		Player:          p,
		Inventory:       inv,
		MaxHunger:       100.0,
		MaxThirst:       100.0,
		MaxStamina:      100.0,
		Hunger:          100.0,
		Thirst:          100.0,
		Stamina:         100.0,
		HealthRegenRate: 0.5,    // 0.5 health per second
		RegenDelay:      10 * time.Second,
		HungerDecayRate: 0.02,  // Slow hunger decay
		ThirstDecayRate: 0.03,  // Slightly faster thirst decay
		LastDamageTime:  time.Now(),
		CanRegenerate:   true,
	}

	// Apply difficulty settings
	sm.applyDifficultySettings()

	return sm
}

// applyDifficultySettings adjusts rates based on game mode
func (sm *SurvivalManager) applyDifficultySettings() {
	switch sm.Mode {
	case ModeSurvival:
		sm.HungerDecayRate = 0.02
		sm.ThirstDecayRate = 0.03
		sm.HealthRegenRate = 0.5
	case ModeHardcore:
		sm.HungerDecayRate = 0.04
		sm.ThirstDecayRate = 0.05
		sm.HealthRegenRate = 0.25
		sm.RegenDelay = 15 * time.Second
	case ModeCreative:
		sm.HungerDecayRate = 0
		sm.ThirstDecayRate = 0
		sm.HealthRegenRate = 100.0 // Instant healing
		sm.CanRegenerate = true
	}
}

// Update updates survival mechanics
func (sm *SurvivalManager) Update(deltaTime float64) {
	if sm.Mode == ModeCreative {
		// In creative mode, keep everything full
		sm.Hunger = sm.MaxHunger
		sm.Thirst = sm.MaxThirst
		sm.Stamina = sm.MaxStamina
		if sm.Player.Health < sm.Player.MaxHealth {
			sm.Player.Health = sm.Player.MaxHealth
		}
		return
	}

	// Decay hunger and thirst
	sm.Hunger -= sm.HungerDecayRate * deltaTime
	sm.Thirst -= sm.ThirstDecayRate * deltaTime

	// Clamp values
	if sm.Hunger < 0 {
		sm.Hunger = 0
		sm.IsStarving = true
	} else if sm.Hunger > sm.MaxHunger {
		sm.Hunger = sm.MaxHunger
		sm.IsStarving = false
	}

	if sm.Thirst < 0 {
		sm.Thirst = 0
		sm.IsDehydrated = true
	} else if sm.Thirst > sm.MaxThirst {
		sm.Thirst = sm.MaxThirst
		sm.IsDehydrated = false
	}

	// Regenerate stamina
	if sm.Stamina < sm.MaxStamina {
		sm.Stamina += 5.0 * deltaTime // Fast stamina regen
		if sm.Stamina > sm.MaxStamina {
			sm.Stamina = sm.MaxStamina
		}
	}

	// Health regeneration
	sm.updateHealthRegeneration(deltaTime)

	// Starvation/dehydration damage
	if sm.IsStarving {
		sm.Player.TakeDamage(0.5 * deltaTime)
	}
	if sm.IsDehydrated {
		sm.Player.TakeDamage(1.0 * deltaTime)
	}
}

// updateHealthRegeneration handles regenerative health
func (sm *SurvivalManager) updateHealthRegeneration(deltaTime float64) {
	// Check if enough time has passed since last damage
	timeSinceDamage := time.Since(sm.LastDamageTime)
	if timeSinceDamage < sm.RegenDelay {
		return
	}

	// Only regen if hunger and thirst are above 50%
	if sm.Hunger < sm.MaxHunger*0.5 || sm.Thirst < sm.MaxThirst*0.5 {
		return
	}

	// Regenerate health
	if sm.Player.Health < sm.Player.MaxHealth {
		sm.Player.Heal(sm.HealthRegenRate * deltaTime)
	}
}

// OnPlayerDamaged should be called when player takes damage
func (sm *SurvivalManager) OnPlayerDamaged() {
	sm.LastDamageTime = time.Now()
}

// UseStamina consumes stamina for actions
func (sm *SurvivalManager) UseStamina(amount float64) bool {
	if sm.Stamina >= amount {
		sm.Stamina -= amount
		return true
	}
	return false
}

// EatFood restores hunger
func (sm *SurvivalManager) EatFood(amount float64) {
	sm.Hunger += amount
	if sm.Hunger > sm.MaxHunger {
		sm.Hunger = sm.MaxHunger
	}
}

// Drink restores thirst
func (sm *SurvivalManager) Drink(amount float64) {
	sm.Thirst += amount
	if sm.Thirst > sm.MaxThirst {
		sm.Thirst = sm.MaxThirst
	}
}

// GetSurvivalStats returns current survival stats
func (sm *SurvivalManager) GetSurvivalStats() map[string]float64 {
	return map[string]float64{
		"hunger":      sm.Hunger,
		"max_hunger":  sm.MaxHunger,
		"thirst":      sm.Thirst,
		"max_thirst":  sm.MaxThirst,
		"stamina":     sm.Stamina,
		"max_stamina": sm.MaxStamina,
		"health":      sm.Player.Health,
		"max_health":  sm.Player.MaxHealth,
	}
}

// SetGameMode changes the game mode
func (sm *SurvivalManager) SetGameMode(mode GameMode) {
	sm.Mode = mode
	sm.applyDifficultySettings()
}

// IsSurvival returns true if in survival or hardcore mode
func (sm *SurvivalManager) IsSurvival() bool {
	return sm.Mode == ModeSurvival || sm.Mode == ModeHardcore
}

// IsHardcore returns true if in hardcore mode
func (sm *SurvivalManager) IsHardcore() bool {
	return sm.Mode == ModeHardcore
}
