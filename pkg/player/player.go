package player

import (
	"math"
	"tesselbox/pkg/creatures"
	"tesselbox/pkg/items"
	"tesselbox/pkg/status"
	"tesselbox/pkg/world"
	"time"
)

const (
	// Physics constants
	Gravity      = 0.5
	PlayerSpeed  = 300.0 // Speed in pixels per second (framerate independent)
	JumpForce    = -8.0  // Jump force in pixels per second
	Friction     = 0.85
	TerminalVelX = 300.0
	TerminalVelY = 1200.0 // Increased for faster falling
	MiningRange  = 150.0
	PlayerWidth  = 40.0
	PlayerHeight = 40.0 // Square player
)

// Player represents a player in the game
type Player struct {
	X, Y          float64
	VX, VY        float64
	Width, Height float64

	// Movement state
	MovingLeft  bool
	MovingRight bool
	Jumping     bool
	OnGround    bool

	// Mining state
	Mining          bool
	MiningTarget    *world.Hexagon
	MiningProgress  float64
	MiningStartTime time.Time

	// Inventory (reference to inventory)
	SelectedSlot int

	// Health and stats
	Health    float64
	MaxHealth float64

	// Combat state
	Attacking      bool
	AttackProgress float64
	AttackCooldown float64
	LastAttackTime time.Time

	// Combat stats
	AttackDamage float64
	AttackRange  float64
	AttackSpeed  float64 // Attacks per second

	// Equipment
	Weapon     *items.Item
	Helmet     *items.Item
	Chestplate *items.Item
	Leggings   *items.Item
	Boots      *items.Item

	// Status effects
	StatusManager *status.StatusManager

	// Leveling system
	Level      int
	Experience int
	ExpToNext  int

	// Time tracking for delta time
	LastUpdateTime time.Time
}

// NewPlayer creates a new player at the specified position
func NewPlayer(x, y float64) *Player {
	return &Player{
		X:              x,
		Y:              y,
		VX:             0,
		VY:             0,
		Width:          PlayerWidth,
		Height:         PlayerHeight,
		Health:         100.0,
		MaxHealth:      100.0,
		SelectedSlot:   0,
		LastUpdateTime: time.Now(),

		// Combat stats
		AttackDamage: 10.0,
		AttackRange:  80.0,
		AttackSpeed:  1.0, // 1 attack per second

		// Status effects
		StatusManager: status.NewStatusManager(),

		// Leveling system
		Level:      1,
		Experience: 0,
		ExpToNext:  100, // Experience needed for level 2
	}
}

// Update updates the player's physics with delta time
// This fixes the lag issue by using actual elapsed time instead of hardcoded FPS
func (p *Player) Update(deltaTime float64) {
	// Clamp delta time to prevent physics explosions on frame drops
	if deltaTime > 0.1 {
		deltaTime = 0.1
	}
	if deltaTime < 0.001 {
		deltaTime = 0.001
	}

	// Apply horizontal movement with acceleration
	if p.MovingLeft {
		p.VX -= PlayerSpeed * deltaTime * 10 // Quick acceleration
	} else if p.MovingRight {
		p.VX += PlayerSpeed * deltaTime * 10
	} else {
		// Apply friction for smooth stopping
		p.VX *= Friction
	}

	// Clamp horizontal velocity
	if p.VX > TerminalVelX {
		p.VX = TerminalVelX
	} else if p.VX < -TerminalVelX {
		p.VX = -TerminalVelX
	}

	// Stop very small movements to prevent jitter
	if !p.MovingLeft && !p.MovingRight && p.VX > -0.1 && p.VX < 0.1 {
		p.VX = 0
	}

	// Apply gravity (framerate independent)
	p.VY += Gravity * deltaTime * 60 * 60 // Proper scaling for consistent gravity

	// Clamp vertical velocity (increased for faster falling)
	if p.VY > TerminalVelY {
		p.VY = TerminalVelY
	}

	// Jump with delta time
	if p.Jumping && p.OnGround {
		p.VY = JumpForce * 60 // Scale jump force to match gravity scaling
		p.OnGround = false
	}

	// Reset jumping flag
	p.Jumping = false

	// Update mining progress
	if p.Mining && p.MiningTarget != nil {
		if p.MiningStartTime.IsZero() {
			p.MiningStartTime = time.Now()
		}
		elapsed := time.Since(p.MiningStartTime).Seconds()
		p.MiningProgress = elapsed * 100.0 // 0-100% over time
		if p.MiningProgress >= 100.0 {
			p.MiningProgress = 0
			p.Mining = false
			p.MiningStartTime = time.Time{}
		}
	}

	// Update attack cooldown
	if p.AttackCooldown > 0 {
		p.AttackCooldown -= deltaTime
		if p.AttackCooldown < 0 {
			p.AttackCooldown = 0
		}
	}

	// Update attack animation progress
	if p.Attacking {
		p.AttackProgress += deltaTime * p.AttackSpeed * 2 // Animation speed
		if p.AttackProgress >= 1.0 {
			p.Attacking = false
			p.AttackProgress = 0
		}
	}
}

// UpdateWithCollision updates player position with collision detection
// This should be called after Update() with the nearby hexagons
func (p *Player) UpdateWithCollision(deltaTime float64, checkCollision func(float64, float64, float64, float64) bool) {
	// Clamp delta time
	if deltaTime > 0.1 {
		deltaTime = 0.1
	}
	if deltaTime < 0.001 {
		deltaTime = 0.001
	}

	// Update position with delta time
	p.X += p.VX * deltaTime
	p.Y += p.VY * deltaTime

	// Get player bounds
	minX, minY, maxX, maxY := p.GetBounds()

	// Check vertical collision (ground detection) - increased height for better ground detection
	bottomLeftCollision := checkCollision(minX, maxY, maxX, maxY+10)
	bottomRightCollision := checkCollision(minX+p.Width/2, maxY, maxX, maxY+10)

	if bottomLeftCollision || bottomRightCollision {
		// We hit the ground - snap to it
		p.VY = 0
		p.OnGround = true

		// Find the exact ground position by checking from current position upward
		groundY := p.Y
		for checkY := p.Y + p.Height; checkY > p.Y-10; checkY-- {
			if !checkCollision(minX, checkY, maxX, checkY+1) {
				groundY = checkY - p.Height
				break
			}
		}
		p.Y = groundY
	} else {
		// No ground below - player is falling
		p.OnGround = false
	}

	// Check horizontal collision (walls)
	if p.VX < 0 { // Moving left
		leftCollision := checkCollision(minX-1, minY+5, minX, maxY-5)
		if leftCollision {
			p.X = minX + 1
			p.VX = 0
		}
	} else if p.VX > 0 { // Moving right
		rightCollision := checkCollision(maxX, minY+5, maxX+1, maxY-5)
		if rightCollision {
			p.X = maxX - p.Width - 1
			p.VX = 0
		}
	}

	// Check ceiling collision (head bump)
	if p.VY < 0 { // Moving upward
		ceilingLeftCollision := checkCollision(minX, minY-1, minX+p.Width/2, minY)
		ceilingRightCollision := checkCollision(minX+p.Width/2, minY-1, maxX, minY)
		if ceilingLeftCollision || ceilingRightCollision {
			p.VY = 0
			p.Y = minY + 1
		}
	}
}

// GetCenter returns the center position of the player
func (p *Player) GetCenter() (float64, float64) {
	return p.X + p.Width/2.0, p.Y + p.Height/2.0
}

// GetPosition returns the top-left position of the player
func (p *Player) GetPosition() (float64, float64) {
	return p.X, p.Y
}

// SetPosition sets the player's position
func (p *Player) SetPosition(x, y float64) {
	p.X = x
	p.Y = y
}

// Move moves the player by the specified offset
func (p *Player) Move(dx, dy float64) {
	p.X += dx
	p.Y += dy
}

// GetVelocity returns the player's velocity
func (p *Player) GetVelocity() (float64, float64) {
	return p.VX, p.VY
}

// SetVelocity sets the player's velocity
func (p *Player) SetVelocity(vx, vy float64) {
	p.VX = vx
	p.VY = vy
}

// Jump makes the player jump if on ground
func (p *Player) Jump() {
	if p.OnGround {
		p.Jumping = true
	}
}

// IsOnGround returns true if the player is on the ground
func (p *Player) IsOnGround() bool {
	return p.OnGround
}

// SetOnGround sets the player's on-ground state
func (p *Player) SetOnGround(onGround bool) {
	p.OnGround = onGround
}

// GetBounds returns the player's bounding box
func (p *Player) GetBounds() (float64, float64, float64, float64) {
	return p.X, p.Y, p.X + p.Width, p.Y + p.Height
}

// GetHealth returns the player's current health
func (p *Player) GetHealth() float64 {
	return p.Health
}

// SetHealth sets the player's health
func (p *Player) SetHealth(health float64) {
	p.Health = health
	if p.Health > p.MaxHealth {
		p.Health = p.MaxHealth
	}
	if p.Health < 0 {
		p.Health = 0
	}
}

// GetMaxHealth returns the player's maximum health
func (p *Player) GetMaxHealth() float64 {
	return p.MaxHealth
}

// GetDefense returns the player's total defense from armor
func (p *Player) GetDefense() float64 {
	defense := 0.0

	armorPieces := []*items.Item{p.Helmet, p.Chestplate, p.Leggings, p.Boots}
	for _, armor := range armorPieces {
		if armor != nil {
			props := items.GetItemProperties(armor.Type)
			if props != nil && props.IsArmor {
				defense += props.ArmorDefense
			}
		}
	}

	return defense
}

// EquipItem equips an item to the appropriate slot
func (p *Player) EquipItem(item *items.Item) bool {
	if item == nil {
		return false
	}

	props := items.GetItemProperties(item.Type)
	if props == nil {
		return false
	}

	if props.IsWeapon {
		// Unequip current weapon first
		if p.Weapon != nil {
			p.UnequipWeapon()
		}
		p.Weapon = item
		return true
	}

	if props.IsArmor {
		switch props.ArmorType {
		case "helmet":
			if p.Helmet != nil {
				p.UnequipArmor("helmet")
			}
			p.Helmet = item
			return true
		case "chestplate":
			if p.Chestplate != nil {
				p.UnequipArmor("chestplate")
			}
			p.Chestplate = item
			return true
		case "leggings":
			if p.Leggings != nil {
				p.UnequipArmor("leggings")
			}
			p.Leggings = item
			return true
		case "boots":
			if p.Boots != nil {
				p.UnequipArmor("boots")
			}
			p.Boots = item
			return true
		}
	}

	return false
}

// UnequipWeapon removes the equipped weapon
func (p *Player) UnequipWeapon() *items.Item {
	weapon := p.Weapon
	p.Weapon = nil
	return weapon
}

// UnequipArmor removes equipped armor of the specified type
func (p *Player) UnequipArmor(armorType string) *items.Item {
	switch armorType {
	case "helmet":
		armor := p.Helmet
		p.Helmet = nil
		return armor
	case "chestplate":
		armor := p.Chestplate
		p.Chestplate = nil
		return armor
	case "leggings":
		armor := p.Leggings
		p.Leggings = nil
		return armor
	case "boots":
		armor := p.Boots
		p.Boots = nil
		return armor
	}
	return nil
}

// TakeDamage reduces the player's health, considering armor defense
func (p *Player) TakeDamage(amount float64, fromX, fromY, knockbackForce float64) {
	defense := p.GetDefense()
	actualDamage := amount - defense
	if actualDamage < 1 {
		actualDamage = 1 // Minimum damage
	}

	p.Health -= actualDamage
	if p.Health < 0 {
		p.Health = 0
	}

	// Apply knockback when taking damage
	if knockbackForce > 0 {
		p.ApplyKnockback(fromX, fromY, knockbackForce)
	}
}

// ApplyKnockback applies knockback force to the player
func (p *Player) ApplyKnockback(fromX, fromY float64, force float64) {
	// Calculate direction away from the knockback source
	dx := p.X + p.Width/2 - fromX
	dy := p.Y + p.Height/2 - fromY

	// Normalize direction
	distance := math.Sqrt(dx*dx + dy*dy)
	if distance > 0 {
		dx /= distance
		dy /= distance
	} else {
		// Default knockback direction if at same position
		dx = 1
		dy = -0.5
	}

	// Apply knockback velocity
	p.VX += dx * force
	p.VY += dy * force * 0.5 // Less vertical knockback

	// Cap knockback velocity
	maxKnockback := 200.0
	if p.VX > maxKnockback {
		p.VX = maxKnockback
	} else if p.VX < -maxKnockback {
		p.VX = -maxKnockback
	}
	if p.VY > maxKnockback {
		p.VY = maxKnockback
	} else if p.VY < -maxKnockback {
		p.VY = -maxKnockback
	}
}

// Heal increases the player's health
func (p *Player) Heal(amount float64) {
	p.Health += amount
	if p.Health > p.MaxHealth {
		p.Health = p.MaxHealth
	}
}

// IsAlive returns true if the player is alive
func (p *Player) IsAlive() bool {
	return p.Health > 0
}

// StartMining starts mining at the target hexagon
func (p *Player) StartMining(target *world.Hexagon) {
	p.Mining = true
	p.MiningTarget = target
	p.MiningProgress = 0
	p.MiningStartTime = time.Now()
}

// StopMining stops mining
func (p *Player) StopMining() {
	p.Mining = false
	p.MiningTarget = nil
	p.MiningProgress = 0
	p.MiningStartTime = time.Time{}
}

// GetMiningProgress returns the current mining progress (0-100)
func (p *Player) GetMiningProgress() float64 {
	return p.MiningProgress
}

// IsMining returns true if the player is currently mining
func (p *Player) IsMining() bool {
	return p.Mining
}

// GetMiningTarget returns the hexagon being mined
func (p *Player) GetMiningTarget() *world.Hexagon {
	return p.MiningTarget
}

// GetMiningRange returns the player's mining range
func (p *Player) GetMiningRange() float64 {
	return MiningRange
}

// DistanceTo returns the distance from the player to a point
func (p *Player) DistanceTo(x, y float64) float64 {
	centerX, centerY := p.GetCenter()
	dx := centerX - x
	dy := centerY - y
	return dx*dx + dy*dy // Return squared distance for efficiency
}

// CanReach returns true if the player can reach a point
func (p *Player) CanReach(x, y float64) bool {
	return p.DistanceTo(x, y) <= MiningRange*MiningRange
}

// SetSelectedSlot sets the currently selected inventory slot
func (p *Player) SetSelectedSlot(slot int) {
	p.SelectedSlot = slot
}

// GetSelectedSlot returns the currently selected inventory slot
func (p *Player) GetSelectedSlot() int {
	return p.SelectedSlot
}

// Attack performs a melee attack if cooldown allows
func (p *Player) Attack() bool {
	if p.AttackCooldown > 0 {
		return false // Still on cooldown
	}

	p.Attacking = true
	p.AttackProgress = 0

	// Use weapon speed if equipped, otherwise base speed
	attackSpeed := p.AttackSpeed
	if p.Weapon != nil {
		props := items.GetItemProperties(p.Weapon.Type)
		if props != nil && props.IsWeapon {
			attackSpeed = props.WeaponSpeed
		}
	}

	p.AttackCooldown = 1.0 / attackSpeed // Cooldown based on attack speed
	p.LastAttackTime = time.Now()

	return true
}

// IsAttackReady returns true if the player can attack
func (p *Player) IsAttackReady() bool {
	return p.AttackCooldown <= 0
}

// GetAttackSpeed returns the player's attack speed (attacks per second)
func (p *Player) GetAttackSpeed() float64 {
	speed := p.AttackSpeed

	// Use weapon speed if equipped
	if p.Weapon != nil {
		props := items.GetItemProperties(p.Weapon.Type)
		if props != nil && props.IsWeapon {
			speed = props.WeaponSpeed
		}
	}

	return speed
}

// GetAttackDamage returns the player's attack damage
func (p *Player) GetAttackDamage() float64 {
	damage := p.AttackDamage // Base damage

	// Add weapon damage if equipped
	if p.Weapon != nil {
		props := items.GetItemProperties(p.Weapon.Type)
		if props != nil && props.IsWeapon {
			damage = props.WeaponDamage
		}
	}

	// Apply status effect multipliers
	if p.StatusManager != nil {
		damage *= p.StatusManager.GetDamageMultiplier()
	}

	return damage
}

// DealDamage applies damage to a target
func (p *Player) DealDamage(target interface{}, damage float64) {
	switch t := target.(type) {
	case *creatures.Creature:
		// Damage a creature
		t.TakeDamage(damage)
		// Apply knockback from player to creature
		knockbackForce := 100.0
		playerCenterX, playerCenterY := p.GetCenter()
		t.ApplyKnockback(playerCenterX, playerCenterY, knockbackForce)
	default:
		// Placeholder for other damageable objects
	}
}

// GetAttackProgress returns the current attack animation progress (0-1)
func (p *Player) GetAttackProgress() float64 {
	return p.AttackProgress
}

// IsAttacking returns true if the player is currently attacking
func (p *Player) IsAttacking() bool {
	return p.Attacking
}

// UpdateCombat updates combat-related state (called from main Update)
func (p *Player) UpdateCombat(deltaTime float64) {
	// Update combat cooldown
	if p.AttackCooldown > 0 {
		p.AttackCooldown -= deltaTime
		if p.AttackCooldown < 0 {
			p.AttackCooldown = 0
		}
	}

	// Update attack animation progress
	if p.Attacking {
		attackSpeed := p.GetAttackSpeed()
		p.AttackProgress += deltaTime * attackSpeed
		if p.AttackProgress >= 1.0 {
			p.Attacking = false
			p.AttackProgress = 0
		}
	}

	// Update status effects
	if p.StatusManager != nil {
		p.StatusManager.Update()
		p.StatusManager.ApplyPeriodicEffects(&p.Health, p.MaxHealth)
	}
}

// AreaAttack performs an area attack damaging multiple creatures
func (p *Player) AreaAttack(radius float64) bool {
	if p.AttackCooldown > 0 {
		return false // Still on cooldown
	}

	// This would need access to the world to get nearby creatures
	// For now, just set up the attack - the caller will handle the damage

	p.Attacking = true
	p.AttackProgress = 0
	p.AttackCooldown = 2.0 // Longer cooldown for area attack

	return true
}

// ApplyStatusEffect applies a status effect to a target
func (p *Player) ApplyStatusEffect(target interface{}, effectType status.StatusEffectType, duration time.Duration, strength float64) {
	// This would need the target's status manager
	// For creatures, we'd need to add StatusManager to creatures too
}

// SelfBuff applies a buff to the player
func (p *Player) SelfBuff(buffType status.StatusEffectType, duration time.Duration, strength float64) bool {
	if p.StatusManager == nil {
		return false
	}

	p.StatusManager.ApplyEffect(buffType, duration, strength)
	return true
}

// WeaponSpecial performs a weapon-specific special ability
func (p *Player) WeaponSpecial() bool {
	if p.Weapon == nil || p.AttackCooldown > 0 {
		return false
	}

	props := items.GetItemProperties(p.Weapon.Type)
	if props == nil || !props.IsWeapon {
		return false
	}

	// Different abilities based on weapon type
	switch props.WeaponType {
	case "melee":
		// Power attack - extra damage
		return p.PowerAttack()
	case "ranged":
		// Multi-shot
		return p.MultiShot()
	case "magic":
		// Magic blast
		return p.MagicBlast()
	}

	return false
}

// PowerAttack performs a powerful melee attack
func (p *Player) PowerAttack() bool {
	p.Attacking = true
	p.AttackProgress = 0
	p.AttackCooldown = 3.0 // Long cooldown
	return true
}

// MultiShot performs multiple ranged attacks
func (p *Player) MultiShot() bool {
	p.Attacking = true
	p.AttackProgress = 0
	p.AttackCooldown = 2.5 // Medium cooldown
	return true
}

// MagicBlast performs a magical area attack
func (p *Player) MagicBlast() bool {
	p.Attacking = true
	p.AttackProgress = 0
	p.AttackCooldown = 4.0 // Very long cooldown
	return true
}

// GainExperience adds experience points and handles leveling
func (p *Player) GainExperience(amount int) {
	p.Experience += amount

	// Check for level up
	for p.Experience >= p.ExpToNext {
		p.LevelUp()
	}
}

// LevelUp increases the player's level and stats
func (p *Player) LevelUp() {
	p.Level++
	p.Experience -= p.ExpToNext

	// Increase experience needed for next level (exponential growth)
	p.ExpToNext = int(float64(p.ExpToNext) * 1.5)

	// Increase base stats
	p.MaxHealth += 10.0
	p.Health = p.MaxHealth // Full heal on level up
	p.AttackDamage += 2.0
	p.AttackSpeed += 0.1

	// Cap attack speed at reasonable maximum
	if p.AttackSpeed > 3.0 {
		p.AttackSpeed = 3.0
	}
}

// GetLevel returns the player's current level
func (p *Player) GetLevel() int {
	return p.Level
}

// GetExperience returns the player's current experience
func (p *Player) GetExperience() int {
	return p.Experience
}

// GetExpToNext returns experience needed for next level
func (p *Player) GetExpToNext() int {
	return p.ExpToNext
}

// GetAttackRange returns the player's attack range
func (p *Player) GetAttackRange() float64 {
	return p.AttackRange
}
