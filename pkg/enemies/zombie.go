// Package enemies implements hostile entities including zombies.
package enemies

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"tesselbox/pkg/gametime"
	"tesselbox/pkg/player"
)

// ZombieType represents different zombie variants
type ZombieType int

const (
	ZombieNormal ZombieType = iota
	ZombieFast
	ZombieStrong
	ZombieTank
)

// Zombie represents a hostile undead enemy
type Zombie struct {
	ID     string
	Type   ZombieType
	X, Y   float64
	VX, VY float64
	Width  float64
	Height float64

	// Combat stats
	Health         float64
	MaxHealth      float64
	Damage         float64
	Speed          float64
	AttackRange    float64
	AttackCooldown time.Duration
	LastAttack     time.Time

	// AI state
	Target    *player.Player
	State     ZombieState
	SpawnTime time.Time
	IsBurning bool
	IsAlive   bool

	// Light sensitivity
	LightDamageRate float64 // Damage per second when in light

	// Physics state (same as player)
	OnGround bool // For gravity/collision
}

// ZombieState represents AI states
type ZombieState int

const (
	ZombieIdle ZombieState = iota
	ZombieChasing
	ZombieAttacking
	ZombieBurning
	ZombieDying
)

// DamageCallback is called when a zombie deals damage to the player
type DamageCallback func(damage float64, zombieX, zombieY float64)

// ZombieSpawner manages zombie spawning
type ZombieSpawner struct {
	Zombies        []*Zombie
	MaxZombies     int
	SpawnRadius    float64
	DespawnRadius  float64
	LastSpawnTime  time.Time
	SpawnCooldown  time.Duration
	DayNightCycle  *gametime.DayNightCycle
	NextID         int
	OnPlayerDamage DamageCallback // Callback for when player takes damage
}

// NewZombieSpawner creates a new zombie spawner
func NewZombieSpawner(dayNight *gametime.DayNightCycle) *ZombieSpawner {
	return &ZombieSpawner{
		Zombies:       make([]*Zombie, 0),
		MaxZombies:    15,
		SpawnRadius:   800,
		DespawnRadius: 1500,
		SpawnCooldown: 3 * time.Second,
		DayNightCycle: dayNight,
		NextID:        1,
	}
}

// Update updates all zombies and handles spawning/despawning
func (zs *ZombieSpawner) Update(deltaTime float64, player *player.Player, ambientLight float64,
	checkCollision func(float64, float64, float64, float64) bool, worldSpawnFunc func(float64, float64) (float64, float64)) {

	// Check if it's night time (light < 0.3)
	isNight := ambientLight < 0.3

	// Spawn new zombies at night
	if isNight && len(zs.Zombies) < zs.MaxZombies {
		if time.Since(zs.LastSpawnTime) > zs.SpawnCooldown {
			// Spawn everywhere - find valid spawn positions like player
			zombie := zs.spawnZombieEverywhere(player.X, player.Y, worldSpawnFunc)
			if zombie != nil {
				zs.Zombies = append(zs.Zombies, zombie)
				zs.LastSpawnTime = time.Now()
			}
		}
	}

	// Update existing zombies
	activeZombies := []*Zombie{}
	for _, zombie := range zs.Zombies {
		zombie.Update(deltaTime, player, ambientLight, zs.OnPlayerDamage)

		// Apply collision detection if collision function provided
		if checkCollision != nil {
			zombie.UpdateWithCollision(deltaTime, checkCollision)
		}

		// Despawn if too far or dead
		distance := distance(zombie.X, zombie.Y, player.X, player.Y)
		if zombie.IsAlive && distance < zs.DespawnRadius {
			activeZombies = append(activeZombies, zombie)
		}
	}
	zs.Zombies = activeZombies
}

// canSpawnAt checks if a location is suitable for zombie spawning
func (zs *ZombieSpawner) canSpawnAt(px, py float64, ambientLight float64) bool {
	// Only spawn in dark areas (night time)
	return ambientLight < 0.3
}

// spawnZombie creates a new zombie near the player
func (zs *ZombieSpawner) spawnZombie(px, py float64) *Zombie {
	// Random position around player (but not too close)
	angle := rand.Float64() * 2 * math.Pi
	dist := zs.SpawnRadius*0.5 + rand.Float64()*zs.SpawnRadius*0.5

	sx := px + math.Cos(angle)*dist
	sy := py + math.Sin(angle)*dist

	// Determine zombie type based on random chance
	var ztype ZombieType
	r := rand.Float64()
	switch {
	case r < 0.6:
		ztype = ZombieNormal
	case r < 0.8:
		ztype = ZombieFast
	case r < 0.95:
		ztype = ZombieStrong
	default:
		ztype = ZombieTank
	}

	return NewZombie(zs.NextID, ztype, sx, sy)
}

// spawnZombieEverywhere spawns a zombie at a valid position using the world's spawn function
// This allows zombies to spawn everywhere with proper terrain, same as player spawning
func (zs *ZombieSpawner) spawnZombieEverywhere(playerX, playerY float64, worldSpawnFunc func(float64, float64) (float64, float64)) *Zombie {
	// Try multiple spawn positions
	for attempts := 0; attempts < 10; attempts++ {
		// Random position around player
		angle := rand.Float64() * 2 * math.Pi
		dist := zs.SpawnRadius*0.3 + rand.Float64()*zs.SpawnRadius*0.7

		tryX := playerX + math.Cos(angle)*dist
		tryY := playerY + math.Sin(angle)*dist

		// Use world spawn function to find valid ground position
		spawnX, spawnY := worldSpawnFunc(tryX, tryY)

		// Make sure we found a valid position above ground
		if spawnY < 10000 { // Valid spawn found (not the fallback max value)
			// Place zombie above ground like player
			zombieY := spawnY - 200

			// Determine zombie type based on random chance
			var ztype ZombieType
			r := rand.Float64()
			switch {
			case r < 0.6:
				ztype = ZombieNormal
			case r < 0.8:
				ztype = ZombieFast
			case r < 0.95:
				ztype = ZombieStrong
			default:
				ztype = ZombieTank
			}

			zombie := NewZombie(zs.NextID, ztype, spawnX, zombieY)
			zs.NextID++
			return zombie
		}
	}

	return nil // Could not find valid spawn
}

// NewZombie creates a new zombie - decaying version of player
func NewZombie(id int, ztype ZombieType, x, y float64) *Zombie {
	// Zombies are same size as player - 50x50 cubes
	zombie := &Zombie{
		ID:              fmt.Sprintf("zombie_%d", id),
		Type:            ztype,
		X:               x,
		Y:               y,
		Width:           50, // Same as player
		Height:          50, // Same as player (square cube)
		VX:              0,
		VY:              0,
		SpawnTime:       time.Now(),
		IsAlive:         true,
		LightDamageRate: 10.0,
		State:           ZombieIdle,
		OnGround:        false,
	}

	// Set stats based on type - slower than player, same size
	switch ztype {
	case ZombieNormal:
		zombie.MaxHealth = 20
		zombie.Health = 20
		zombie.Damage = 5
		zombie.Speed = 60 // Slower than player (300)
		zombie.AttackRange = 50
		zombie.AttackCooldown = 1 * time.Second
	case ZombieFast:
		zombie.MaxHealth = 15
		zombie.Health = 15
		zombie.Damage = 3
		zombie.Speed = 120 // Still slower than player
		zombie.AttackRange = 50
		zombie.AttackCooldown = 800 * time.Millisecond
	case ZombieStrong:
		zombie.MaxHealth = 30
		zombie.Health = 30
		zombie.Damage = 10
		zombie.Speed = 40
		zombie.AttackRange = 60
		zombie.AttackCooldown = 1200 * time.Millisecond
	case ZombieTank:
		zombie.MaxHealth = 50
		zombie.Health = 50
		zombie.Damage = 6
		zombie.Speed = 30
		zombie.AttackRange = 55
		zombie.AttackCooldown = 1500 * time.Millisecond
	}

	return zombie
}

// Update updates the zombie AI and state
func (z *Zombie) Update(deltaTime float64, player *player.Player, ambientLight float64, damageCallback DamageCallback) {
	if !z.IsAlive {
		return
	}

	// Check if zombie is in light - zombies die in light
	if ambientLight > 0.4 {
		// Take damage from light
		damage := z.LightDamageRate * deltaTime
		z.Health -= damage
		z.IsBurning = true
		z.State = ZombieBurning

		if z.Health <= 0 {
			z.Die()
			return
		}
	} else {
		z.IsBurning = false
		if z.State == ZombieBurning {
			z.State = ZombieIdle
		}
	}

	// Calculate distance to player
	dist := distance(z.X, z.Y, player.X, player.Y)

	// Apply gravity (same as player) - zombies fall
	const Gravity = 2.0
	const Friction = 0.85
	const TerminalVelY = 1200.0

	z.VY += Gravity * deltaTime * 60.0
	if z.VY > TerminalVelY {
		z.VY = TerminalVelY
	}

	// AI behavior - hex grid aligned movement (no diagonal/angled movement)
	switch z.State {
	case ZombieIdle:
		if dist < 500 { // Detection range
			z.State = ZombieChasing
			z.Target = player
		}

	case ZombieChasing:
		if dist > 600 {
			z.State = ZombieIdle
			z.Target = nil
		} else if dist < z.AttackRange {
			z.State = ZombieAttacking
		} else {
			// Move towards player on hex grid - only left/right (no angled movement)
			z.moveTowardsPlayerHex(player.X, deltaTime)
		}

	case ZombieAttacking:
		if dist > z.AttackRange*1.2 {
			z.State = ZombieChasing
		} else {
			// Attack player
			if time.Since(z.LastAttack) > z.AttackCooldown {
				z.attack(player, damageCallback)
			}
		}

	case ZombieBurning:
		// Still try to attack player while burning
		if dist < z.AttackRange && time.Since(z.LastAttack) > z.AttackCooldown {
			z.attack(player, damageCallback)
		}
		// Panic movement - but still hex-aligned
		z.VX += (rand.Float64() - 0.5) * z.Speed * 2 * deltaTime
	}

	// Apply velocity
	z.X += z.VX * deltaTime
	z.Y += z.VY * deltaTime

	// Apply friction
	z.VX *= Friction

	// Stop very small movements
	if z.VX > -0.1 && z.VX < 0.1 {
		z.VX = 0
	}
}

// UpdateWithCollision updates zombie position with collision detection
// Similar to player's UpdateWithCollision
func (z *Zombie) UpdateWithCollision(deltaTime float64, checkCollision func(float64, float64, float64, float64) bool) {
	// Get zombie bounds
	minX, minY, maxX, maxY := z.GetBounds()

	// Check vertical collision (ground detection) - check from zombie's feet downward
	feetY := maxY // Zombie's feet position
	groundCheckDistance := 5.0

	bottomLeftCollision := checkCollision(minX, feetY, minX+z.Width/2, feetY+groundCheckDistance)
	bottomRightCollision := checkCollision(minX+z.Width/2, feetY, maxX, feetY+groundCheckDistance)
	bottomCenterCollision := checkCollision(minX+z.Width/2, feetY, minX+z.Width/2+1, feetY+groundCheckDistance)

	if bottomLeftCollision || bottomRightCollision || bottomCenterCollision {
		// We hit the ground - stop falling and snap to ground
		if z.VY > 0 { // Only if moving downward
			z.VY = 0
			z.OnGround = true

			// Find exact ground position
			for checkY := feetY; checkY <= feetY+groundCheckDistance; checkY += 1.0 {
				if checkCollision(minX, checkY, maxX, checkY+1) {
					z.Y = checkY - z.Height
					break
				}
			}
		}
	} else {
		// No ground below - zombie is falling
		z.OnGround = false
	}

	// Check horizontal collision (walls)
	if z.VX < 0 { // Moving left
		leftCollision := checkCollision(minX-1, minY+5, minX, maxY-5)
		if leftCollision {
			z.X = minX + 1
			z.VX = 0
		}
	} else if z.VX > 0 { // Moving right
		rightCollision := checkCollision(maxX, minY+5, maxX+1, maxY-5)
		if rightCollision {
			z.X = maxX - z.Width - 1
			z.VX = 0
		}
	}

	// Check ceiling collision (head bump)
	if z.VY < 0 { // Moving upward
		ceilingLeftCollision := checkCollision(minX, minY-1, minX+z.Width/2, minY)
		ceilingRightCollision := checkCollision(minX+z.Width/2, minY-1, maxX, minY)
		if ceilingLeftCollision || ceilingRightCollision {
			z.VY = 0
			z.Y = minY + 1
		}
	}
}

// moveTowardsPlayerHex moves zombie towards player on hex grid (only horizontal)
// Zombies move left/right like players, falling with gravity, no diagonal/angled movement
func (z *Zombie) moveTowardsPlayerHex(playerX float64, deltaTime float64) {
	dx := playerX - z.X

	// Only move horizontally on hex grid - no angled movement
	if math.Abs(dx) > 5 { // Small threshold to prevent jitter
		// Accelerate towards player
		if dx > 0 {
			z.VX += z.Speed * deltaTime * 10 // Move right
		} else {
			z.VX -= z.Speed * deltaTime * 10 // Move left
		}

		// Clamp velocity
		const TerminalVelX = 300.0
		if z.VX > TerminalVelX {
			z.VX = TerminalVelX
		} else if z.VX < -TerminalVelX {
			z.VX = -TerminalVelX
		}
	}
}

// attack performs an attack on the player
func (z *Zombie) attack(player *player.Player, callback DamageCallback) {
	z.LastAttack = time.Now()
	player.TakeDamage(z.Damage)
	if callback != nil {
		callback(z.Damage, z.X, z.Y)
	}
}

// TakeDamage applies damage to the zombie
func (z *Zombie) TakeDamage(amount float64) {
	z.Health -= amount
	if z.Health <= 0 {
		z.Die()
	}
}

// Die kills the zombie
func (z *Zombie) Die() {
	z.IsAlive = false
	z.State = ZombieDying
	z.Health = 0
}

// GetBounds returns the zombie's bounding box
func (z *Zombie) GetBounds() (minX, minY, maxX, maxY float64) {
	return z.X, z.Y, z.X + z.Width, z.Y + z.Height
}

// GetCenter returns the zombie's center position
func (z *Zombie) GetCenter() (float64, float64) {
	return z.X + z.Width/2, z.Y + z.Height/2
}

// IsNightTime returns true if it's night time based on ambient light
func IsNightTime(ambientLight float64) bool {
	return ambientLight < 0.3
}

// Helper function
func distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}
