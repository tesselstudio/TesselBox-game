// Package combat implements weapon swing mechanics
package combat

import (
	"math"
	"math/rand"
	"time"

	"tesselbox/pkg/enemies"
)

// WeaponSwing represents an active weapon swing attack
type WeaponSwing struct {
	Active     bool
	StartX     float64
	StartY     float64
	EndX       float64
	EndY       float64
	StartTime  time.Time
	Duration   time.Duration
	Damage     float64
	SwingAngle float64 // Angle of the swing arc
	SwingWidth float64 // Width of the swing arc in degrees
}

// WeaponSystem manages weapon attacks
type WeaponSystem struct {
	CurrentSwing  *WeaponSwing
	Cooldown      time.Duration
	LastSwingTime time.Time
	IsAttacking   bool
}

// NewWeaponSystem creates a new weapon system
func NewWeaponSystem() *WeaponSystem {
	return &WeaponSystem{
		Cooldown: 500 * time.Millisecond, // 0.5s between attacks
	}
}

// CanAttack checks if the weapon is ready to swing
func (ws *WeaponSystem) CanAttack() bool {
	return time.Since(ws.LastSwingTime) >= ws.Cooldown && !ws.IsAttacking
}

// StartSwing begins a new weapon swing
func (ws *WeaponSystem) StartSwing(startX, startY, targetX, targetY, damage float64) {
	if !ws.CanAttack() {
		return
	}

	// Calculate swing angle
	dx := targetX - startX
	dy := targetY - startY
	angle := math.Atan2(dy, dx) * 180 / math.Pi // Convert to degrees

	ws.CurrentSwing = &WeaponSwing{
		Active:     true,
		StartX:     startX,
		StartY:     startY,
		EndX:       targetX,
		EndY:       targetY,
		StartTime:  time.Now(),
		Duration:   200 * time.Millisecond, // Swing takes 0.2s
		Damage:     damage,
		SwingAngle: angle,
		SwingWidth: 90, // 90 degree arc
	}

	ws.IsAttacking = true
	ws.LastSwingTime = time.Now()
}

// Update updates the weapon swing state
func (ws *WeaponSystem) Update(deltaTime float64) {
	if ws.CurrentSwing != nil && ws.CurrentSwing.Active {
		if time.Since(ws.CurrentSwing.StartTime) > ws.CurrentSwing.Duration {
			ws.CurrentSwing.Active = false
			ws.IsAttacking = false
		}
	}
}

// GetSwingProgress returns the current swing progress (0-1)
func (ws *WeaponSystem) GetSwingProgress() float64 {
	if ws.CurrentSwing == nil || !ws.CurrentSwing.Active {
		return 0
	}

	elapsed := time.Since(ws.CurrentSwing.StartTime).Seconds()
	total := ws.CurrentSwing.Duration.Seconds()
	progress := elapsed / total

	if progress > 1 {
		return 1
	}
	return progress
}

// CheckHit checks if a target is hit by the current swing
func (ws *WeaponSystem) CheckHit(targetX, targetY, targetWidth, targetHeight float64) bool {
	if ws.CurrentSwing == nil || !ws.CurrentSwing.Active {
		return false
	}

	// Only check hits in the first half of the swing
	progress := ws.GetSwingProgress()
	if progress > 0.5 {
		return false // Already passed the hit window
	}

	// Calculate distance from swing start to target
	targetCenterX := targetX + targetWidth/2
	targetCenterY := targetY + targetHeight/2

	dx := targetCenterX - ws.CurrentSwing.StartX
	dy := targetCenterY - ws.CurrentSwing.StartY
	distance := math.Sqrt(dx*dx + dy*dy)

	// Max swing range
	maxRange := 150.0 // pixels
	if distance > maxRange {
		return false
	}

	// Calculate angle to target
	targetAngle := math.Atan2(dy, dx) * 180 / math.Pi

	// Check if target is within swing arc
	angleDiff := math.Abs(targetAngle - ws.CurrentSwing.SwingAngle)
	if angleDiff > 180 {
		angleDiff = 360 - angleDiff
	}

	// Half the swing width on each side
	if angleDiff <= ws.CurrentSwing.SwingWidth/2 {
		return true
	}

	return false
}

// CritTier represents the critical hit tier
type CritTier int

const (
	CritTierGreen  CritTier = iota // Base damage (1x)
	CritTierYellow                 // Moderate (1.5x)
	CritTierRed                    // Severe (2x)
	CritTierPurple                 // Fatal (instant death)
)

// AttackResult contains the result of an attack
type AttackResult struct {
	Hit        bool
	Damage     float64
	IsCritical bool
	Tier       CritTier
	HitX       float64
	HitY       float64
}

// calculateCrit determines critical hit tier based on random chance and headshots
func calculateCrit(baseDamage, zombieHealth float64, isHeadshot bool) (float64, CritTier, bool) {
	// Headshots are always purple (instant death) for normal zombies
	if isHeadshot {
		return zombieHealth, CritTierPurple, true
	}

	// Random critical chance
	// 5% purple, 15% red, 30% yellow, 50% green
	roll := rand.Float64()

	switch {
	case roll < 0.05: // 5% - Purple (Fatal)
		return zombieHealth, CritTierPurple, true
	case roll < 0.20: // 15% - Red (Severe, 2x)
		return baseDamage * 2, CritTierRed, true
	case roll < 0.50: // 30% - Yellow (Moderate, 1.5x)
		return baseDamage * 1.5, CritTierYellow, true
	default: // 50% - Green (Low, base)
		return baseDamage, CritTierGreen, false
	}
}

// PerformAttack executes an attack and returns results for all hit targets
func (ws *WeaponSystem) PerformAttack(playerX, playerY, targetX, targetY, damage float64, zombies []*enemies.Zombie) []AttackResult {
	results := make([]AttackResult, 0)

	// Start the swing
	ws.StartSwing(playerX, playerY, targetX, targetY, damage)
	if !ws.IsAttacking {
		return results
	}

	// Check hits against all zombies
	for _, zombie := range zombies {
		if !zombie.IsAlive {
			continue
		}

		if ws.CheckHit(zombie.X, zombie.Y, zombie.Width, zombie.Height) {
			// Check if headshot (aiming at upper body/head area)
			targetCenterY := zombie.Y + zombie.Height/2
			isHeadshot := (targetCenterY - zombie.Y) < zombie.Height/3

			// Calculate critical hit
			finalDamage, tier, isCrit := calculateCrit(damage, zombie.Health, isHeadshot)

			result := AttackResult{
				Hit:        true,
				Damage:     finalDamage,
				IsCritical: isCrit,
				Tier:       tier,
				HitX:       zombie.X + zombie.Width/2,
				HitY:       zombie.Y + zombie.Height/2,
			}
			results = append(results, result)
		}
	}

	return results
}

// IsSwinging returns true if a swing is currently active
func (ws *WeaponSystem) IsSwinging() bool {
	return ws.IsAttacking && ws.CurrentSwing != nil && ws.CurrentSwing.Active
}

// GetSwingArc returns the current swing arc for visualization
func (ws *WeaponSystem) GetSwingArc() (startAngle, endAngle float64, radius float64) {
	if ws.CurrentSwing == nil {
		return 0, 0, 0
	}

	halfWidth := ws.CurrentSwing.SwingWidth / 2
	startAngle = ws.CurrentSwing.SwingAngle - halfWidth
	endAngle = ws.CurrentSwing.SwingAngle + halfWidth
	radius = 150 // Same as max range

	return startAngle, endAngle, radius
}

// InterruptSwing cancels the current swing
func (ws *WeaponSystem) InterruptSwing() {
	if ws.CurrentSwing != nil {
		ws.CurrentSwing.Active = false
	}
	ws.IsAttacking = false
}
