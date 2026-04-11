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
func (zs *ZombieSpawner) Update(deltaTime float64, player *player.Player, ambientLight float64) {
	// Check if it's night time (light < 0.3)
	isNight := ambientLight < 0.3

	// Spawn new zombies at night
	if isNight && len(zs.Zombies) < zs.MaxZombies {
		if time.Since(zs.LastSpawnTime) > zs.SpawnCooldown {
			// Check if player is in a dark area
			if zs.canSpawnAt(player.X, player.Y, ambientLight) {
				zombie := zs.spawnZombie(player.X, player.Y)
				if zombie != nil {
					zs.Zombies = append(zs.Zombies, zombie)
					zs.LastSpawnTime = time.Now()
				}
			}
		}
	}

	// Update existing zombies
	activeZombies := []*Zombie{}
	for _, zombie := range zs.Zombies {
		zombie.Update(deltaTime, player, ambientLight, zs.OnPlayerDamage)

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
	distance := zs.SpawnRadius*0.5 + rand.Float64()*zs.SpawnRadius*0.5

	sx := px + math.Cos(angle)*distance
	sy := py + math.Sin(angle)*distance

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

// NewZombie creates a new zombie
func NewZombie(id int, ztype ZombieType, x, y float64) *Zombie {
	zombie := &Zombie{
		ID:              fmt.Sprintf("zombie_%d", id),
		Type:            ztype,
		X:               x,
		Y:               y,
		Width:           40,
		Height:          60,
		SpawnTime:       time.Now(),
		IsAlive:         true,
		LightDamageRate: 10.0, // 10 damage per second in light
		State:           ZombieIdle,
	}

	// Set stats based on type
	switch ztype {
	case ZombieNormal:
		zombie.MaxHealth = 20
		zombie.Health = 20
		zombie.Damage = 5
		zombie.Speed = 80
		zombie.AttackRange = 50
		zombie.AttackCooldown = 1 * time.Second
	case ZombieFast:
		zombie.MaxHealth = 15
		zombie.Health = 15
		zombie.Damage = 3
		zombie.Speed = 140
		zombie.AttackRange = 50
		zombie.AttackCooldown = 800 * time.Millisecond
	case ZombieStrong:
		zombie.MaxHealth = 30
		zombie.Health = 30
		zombie.Damage = 10
		zombie.Speed = 60
		zombie.AttackRange = 60
		zombie.AttackCooldown = 1200 * time.Millisecond
	case ZombieTank:
		zombie.MaxHealth = 50
		zombie.Health = 50
		zombie.Damage = 6
		zombie.Speed = 40
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

	// AI behavior
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
			// Move towards player
			z.moveTowards(player.X, player.Y, deltaTime)
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
		if dist < 500 {
			z.moveTowards(player.X, player.Y, deltaTime)
		}
	}

	// Apply velocity
	z.X += z.VX * deltaTime
	z.Y += z.VY * deltaTime

	// Friction
	z.VX *= 0.9
	z.VY *= 0.9
}

// moveTowards moves zombie towards a target position
func (z *Zombie) moveTowards(tx, ty float64, deltaTime float64) {
	dx := tx - z.X
	dy := ty - z.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist > 0 {
		// Normalize and apply speed
		z.VX = (dx / dist) * z.Speed
		z.VY = (dy / dist) * z.Speed
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
