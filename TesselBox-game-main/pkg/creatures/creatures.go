// creatures.go
package creatures

import (
	"fmt"
	"math"
	"math/rand"
	"tesselbox/pkg/hexagon"
	"tesselbox/pkg/items"
	"time"
)

// CreatureType represents different enemy types
type CreatureType int

const (
	SLIME CreatureType = iota
	SPIDER
	ZOMBIE
)

// AIState represents the current behavior state of a creature
type AIState int

const (
	IDLE AIState = iota
	WANDER
	CHASE
	FLEE
)

// Creature represents an enemy creature in the world
type Creature struct {
	ID         string
	Type       CreatureType
	X, Y       float64
	Hex        hexagon.Hexagon
	Health     float64
	MaxHealth  float64
	Damage     float64
	Speed      float64
	AIState    AIState
	TargetX    float64
	TargetY    float64
	LastUpdate time.Time
	SpawnTime  time.Time
	IsHostile  bool
}

// NewCreature creates a new creature of the specified type at the given position
func NewCreature(creatureType CreatureType, x, y float64, hex hexagon.Hexagon) *Creature {
	creature := &Creature{
		ID:         generateCreatureID(),
		Type:       creatureType,
		X:          x,
		Y:          y,
		Hex:        hex,
		AIState:    IDLE,
		LastUpdate: time.Now(),
		SpawnTime:  time.Now(),
		IsHostile:  true,
	}

	// Set type-specific properties
	switch creatureType {
	case SLIME:
		creature.MaxHealth = 20
		creature.Health = 20
		creature.Damage = 8
		creature.Speed = 0.5
	case SPIDER:
		creature.MaxHealth = 15
		creature.Health = 15
		creature.Damage = 6
		creature.Speed = 1.0
	case ZOMBIE:
		creature.MaxHealth = 30
		creature.Health = 30
		creature.Damage = 12
		creature.Speed = 0.3
	}

	return creature
}

// Update updates the creature with basic pathfinding around obstacles
func (c *Creature) Update(playerX, playerY float64, deltaTime float64, isBlocked func(x, y float64) bool) {
	// Simple AI logic
	distanceToPlayer := math.Sqrt((c.X-playerX)*(c.X-playerX) + (c.Y-playerY)*(c.Y-playerY))

	switch c.AIState {
	case IDLE:
		// Chance to start wandering
		if rand.Float64() < 0.01 { // 1% chance per update
			c.AIState = WANDER
			c.TargetX = c.X + (rand.Float64()-0.5)*200
			c.TargetY = c.Y + (rand.Float64()-0.5)*200
		}
		// If player is close, start chasing
		if distanceToPlayer < 100 && c.IsHostile {
			c.AIState = CHASE
		}
	case WANDER:
		// Move towards target with obstacle avoidance
		c.moveTowardsWithAvoidance(c.TargetX, c.TargetY, deltaTime, isBlocked)

		// Check if reached target
		dx := c.TargetX - c.X
		dy := c.TargetY - c.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < 20 {
			c.AIState = IDLE
		}

		// If player gets close while wandering, chase
		if distanceToPlayer < 100 && c.IsHostile {
			c.AIState = CHASE
		}
	case CHASE:
		// Move towards player with obstacle avoidance
		c.moveTowardsWithAvoidance(playerX, playerY, deltaTime, isBlocked)

		// If player gets too far, go back to idle
		if distanceToPlayer > 150 {
			c.AIState = IDLE
		}
	case FLEE:
		// Move away from player with obstacle avoidance
		awayX := c.X - (playerX - c.X)
		awayY := c.Y - (playerY - c.Y)
		c.moveTowardsWithAvoidance(awayX, awayY, deltaTime, isBlocked)

		// If far enough from player, go back to idle
		if distanceToPlayer > 200 {
			c.AIState = IDLE
		}
	}
}

// TakeDamage reduces creature health
func (c *Creature) TakeDamage(amount float64) {
	c.Health -= amount
	if c.Health < 0 {
		c.Health = 0
	}
}

// IsAlive returns true if the creature has health remaining
func (c *Creature) IsAlive() bool {
	return c.Health > 0
}

// AttackPlayer attempts to attack the player if in range
func (c *Creature) AttackPlayer(playerX, playerY float64, dealDamage func(damage float64, fromX, fromY float64, knockback float64)) {
	distance := math.Sqrt((c.X-playerX)*(c.X-playerX) + (c.Y-playerY)*(c.Y-playerY))

	// Check if in attack range
	attackRange := 50.0 // Creature attack range
	if distance <= attackRange && c.IsHostile {
		// Attack the player
		dealDamage(c.Damage, c.X, c.Y, 50.0) // Damage, from position, knockback
	}
}

// ApplyKnockback applies knockback force to a creature
func (c *Creature) ApplyKnockback(fromX, fromY float64, force float64) {
	// Calculate direction away from the knockback source
	dx := c.X - fromX
	dy := c.Y - fromY

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

	// Apply knockback velocity (creatures are simpler, just move them)
	c.X += dx * force * 0.1
	c.Y += dy * force * 0.1
}

// moveTowardsWithAvoidance moves towards a target while avoiding obstacles
func (c *Creature) moveTowardsWithAvoidance(targetX, targetY float64, deltaTime float64, isBlocked func(x, y float64) bool) {
	dx := targetX - c.X
	dy := targetY - c.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist < 5 {
		return // Already close enough
	}

	// Normalize direction
	dx /= dist
	dy /= dist

	// Check if direct path is blocked
	checkDistance := 20.0 // Distance to check ahead
	testX := c.X + dx*checkDistance
	testY := c.Y + dy*checkDistance

	if isBlocked(testX, testY) {
		// Direct path blocked, try alternatives
		// Try left perpendicular
		leftDx := -dy
		leftDy := dx
		testLeftX := c.X + leftDx*checkDistance
		testLeftY := c.Y + leftDy*checkDistance

		if !isBlocked(testLeftX, testLeftY) {
			// Left is clear, use left direction
			dx = leftDx
			dy = leftDy
		} else {
			// Left also blocked, try right perpendicular
			rightDx := dy
			rightDy := -dx
			testRightX := c.X + rightDx*checkDistance
			testRightY := c.Y + rightDy*checkDistance

			if !isBlocked(testRightX, testRightY) {
				// Right is clear, use right direction
				dx = rightDx
				dy = rightDy
			} else {
				// Both sides blocked, try moving away from obstacles
				// This is a simple fallback - just stop or move randomly
				return
			}
		}
	}

	// Move in the chosen direction
	c.X += dx * c.Speed * deltaTime
	c.Y += dy * c.Speed * deltaTime
}

// GetLootDrops returns the items dropped by this creature when killed
func (c *Creature) GetLootDrops() []items.Item {
	switch c.Type {
	case SLIME:
		return []items.Item{
			{Type: items.GEL, Quantity: 1, Durability: -1},
		}
	case SPIDER:
		return []items.Item{
			{Type: items.STRING, Quantity: 1, Durability: -1},
		}
	case ZOMBIE:
		return []items.Item{
			{Type: items.ROTTEN_FLESH, Quantity: 1, Durability: -1},
		}
	default:
		return []items.Item{}
	}
}

// generateCreatureID creates a unique ID for a creature
func generateCreatureID() string {
	return fmt.Sprintf("creature_%d", rand.Int63())
}
