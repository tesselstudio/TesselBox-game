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

const (
	// Zombie uses same physics as player
	ZombieSpeed  = 300.0 // Same as player
	ZombieJump   = -8.0  // Same as player jump force
	Friction     = 0.85
	TerminalVelX = 300.0
	TerminalVelY = 1200.0

	// Same dimensions as player
	ZombieWidth  = 50.0
	ZombieHeight = 50.0

	// Attack with weapon-like behavior (like player)
	AttackRange = 60.0
)

// ZombieType represents different zombie variants
type ZombieType int

const (
	ZombieNormal ZombieType = iota
	ZombieFast
	ZombieStrong
	ZombieTank
)

// Zombie represents a hostile undead enemy - looks and moves like player but green
type Zombie struct {
	ID     string
	Type   ZombieType
	X, Y   float64
	VX, VY float64
	Width  float64
	Height float64

	// Combat stats - player-like
	Health         float64
	MaxHealth      float64
	Damage         float64
	Speed          float64
	AttackRange    float64
	LastAttack     time.Time
	AttackCooldown time.Duration

	// AI state
	Target    *player.Player
	State     ZombieState
	SpawnTime time.Time
	IsBurning bool
	IsAlive   bool

	// Movement state like player
	OnGround    bool
	Jumping     bool
	MovingLeft  bool
	MovingRight bool

	// Light sensitivity
	LightDamageRate float64 // Damage per second when in light

	// Animation
	IsAttacking bool
	AttackTime  time.Time
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

// NewZombie creates a new zombie - green version of player
func NewZombie(id int, ztype ZombieType, x, y float64) *Zombie {
	zombie := &Zombie{
		ID:              fmt.Sprintf("zombie_%d", id),
		Type:            ztype,
		X:               x,
		Y:               y,
		Width:           ZombieWidth,  // Same as player
		Height:          ZombieHeight, // Same as player
		VX:              0,
		VY:              0,
		SpawnTime:       time.Now(),
		IsAlive:         true,
		LightDamageRate: 10.0,
		State:           ZombieIdle,
		OnGround:        false,
		MovingLeft:      false,
		MovingRight:     false,
		Jumping:         false,
		IsAttacking:     false,
	}

	// Set stats based on type - player-like speeds
	switch ztype {
	case ZombieNormal:
		zombie.MaxHealth = 20
		zombie.Health = 20
		zombie.Damage = 5
		zombie.Speed = ZombieSpeed // Same as player
		zombie.AttackRange = AttackRange
		zombie.AttackCooldown = 800 * time.Millisecond
	case ZombieFast:
		zombie.MaxHealth = 15
		zombie.Health = 15
		zombie.Damage = 3
		zombie.Speed = ZombieSpeed * 1.2 // Slightly faster
		zombie.AttackRange = AttackRange
		zombie.AttackCooldown = 500 * time.Millisecond
	case ZombieStrong:
		zombie.MaxHealth = 30
		zombie.Health = 30
		zombie.Damage = 10
		zombie.Speed = ZombieSpeed * 0.8 // Slower but stronger
		zombie.AttackRange = AttackRange * 1.2
		zombie.AttackCooldown = 1000 * time.Millisecond
	case ZombieTank:
		zombie.MaxHealth = 50
		zombie.Health = 50
		zombie.Damage = 6
		zombie.Speed = ZombieSpeed * 0.6 // Slow but tanky
		zombie.AttackRange = AttackRange
		zombie.AttackCooldown = 1500 * time.Millisecond
	}

	return zombie
}

// Update updates the zombie AI with player-like physics
func (z *Zombie) Update(deltaTime float64, player *player.Player, ambientLight float64, damageCallback DamageCallback) {
	if !z.IsAlive {
		return
	}

	// Clamp delta time to prevent physics issues
	if deltaTime > 0.1 {
		deltaTime = 0.1
	}
	if deltaTime < 0.001 {
		deltaTime = 0.001
	}

	// Check if zombie is in light - zombies die in light
	if ambientLight > 0.4 {
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

	// AI behavior - like player controlling with keyboard
	switch z.State {
	case ZombieIdle:
		z.MovingLeft = false
		z.MovingRight = false
		if dist < 500 { // Detection range
			z.State = ZombieChasing
			z.Target = player
		}

	case ZombieChasing:
		z.MovingLeft = false
		z.MovingRight = false
		if dist > 600 {
			z.State = ZombieIdle
			z.Target = nil
		} else if dist < z.AttackRange {
			z.State = ZombieAttacking
		} else {
			// Move towards player - player-like control
			z.moveTowardsPlayer(player, dist, deltaTime)
		}

	case ZombieAttacking:
		if dist > z.AttackRange*1.2 {
			z.State = ZombieChasing
		} else {
			// Attack like player with weapon
			if time.Since(z.LastAttack) > z.AttackCooldown {
				z.attack(player, damageCallback)
			}
		}
		z.MovingLeft = false
		z.MovingRight = false

	case ZombieBurning:
		// Panic but still move towards player
		z.moveTowardsPlayer(player, dist, deltaTime)
	}

	// Apply player-like physics
	z.applyPlayerPhysics(deltaTime)

	// Apply velocity
	z.X += z.VX * deltaTime
	z.Y += z.VY * deltaTime
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

// moveTowardsPlayer moves towards player with player-like controls
func (z *Zombie) moveTowardsPlayer(player *player.Player, dist float64, deltaTime float64) {
	dx := player.X - z.X

	// Player-like movement - left/right
	if math.Abs(dx) > 5 {
		if dx > 0 {
			z.MovingRight = true
			z.MovingLeft = false
			z.VX += z.Speed * deltaTime * 10
		} else {
			z.MovingLeft = true
			z.MovingRight = false
			z.VX -= z.Speed * deltaTime * 10
		}

		// Clamp velocity
		if z.VX > TerminalVelX {
			z.VX = TerminalVelX
		} else if z.VX < -TerminalVelX {
			z.VX = -TerminalVelX
		}
	} else {
		z.MovingLeft = false
		z.MovingRight = false
	}

	// Jump if player is above and can jump
	dy := player.Y - z.Y
	if dy < -30 && z.OnGround && !z.Jumping {
		z.VY = ZombieJump
		z.Jumping = true
		z.OnGround = false
	}
}

// applyPlayerPhysics applies player-like physics to zombie
func (z *Zombie) applyPlayerPhysics(deltaTime float64) {
	// Apply acceleration from movement state
	if z.MovingLeft {
		z.VX -= z.Speed * deltaTime * 10
	} else if z.MovingRight {
		z.VX += z.Speed * deltaTime * 10
	} else {
		// Apply friction for smooth stopping
		z.VX *= Friction
	}

	// Apply gravity
	z.VY += 2.0 * deltaTime * 60.0
	if z.VY > TerminalVelY {
		z.VY = TerminalVelY
	}

	// Stop very small movements to prevent jitter
	if !z.MovingLeft && !z.MovingRight && z.VX > -0.1 && z.VX < 0.1 {
		z.VX = 0
	}
}

// attack performs a player-like weapon attack on the player
func (z *Zombie) attack(player *player.Player, callback DamageCallback) {
	z.LastAttack = time.Now()
	z.IsAttacking = true
	z.AttackTime = time.Now()

	// Deal damage to player (like player hitting with weapon)
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
