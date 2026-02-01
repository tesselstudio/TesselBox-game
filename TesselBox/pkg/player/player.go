package player

import (
	"tesselbox/pkg/world"
)

const (
	Gravity       = 0.5
	PlayerSpeed   = 5.0
	JumpForce     = -12.0
	Friction      = 0.8
	MiningRange   = 150.0
	PlayerWidth   = 40.0
	PlayerHeight  = 80.0
)

// Player represents a player in the game
type Player struct {
	X, Y          float64
	VX, VY        float64
	Width, Height float64
	
	// Movement state
	MovingLeft   bool
	MovingRight  bool
	Jumping      bool
	OnGround     bool
	
	// Mining state
	Mining       bool
	MiningTarget *world.Hexagon
	MiningProgress float64
	
	// Inventory (reference to inventory)
	SelectedSlot int
	
	// Health and stats
	Health       float64
	MaxHealth    float64
}

// NewPlayer creates a new player at the specified position
func NewPlayer(x, y float64) *Player {
	return &Player{
		X:          x,
		Y:          y,
		VX:         0,
		VY:         0,
		Width:      PlayerWidth,
		Height:     PlayerHeight,
		Health:     100.0,
		MaxHealth:  100.0,
		SelectedSlot: 0,
	}
}

// Update updates the player's physics
func (p *Player) Update(deltaTime float64) {
	// Apply horizontal movement
	if p.MovingLeft {
		p.VX = -PlayerSpeed
	} else if p.MovingRight {
		p.VX = PlayerSpeed
	} else {
		p.VX *= Friction // Apply friction
	}
	
	// Apply gravity
	p.VY += Gravity
	
	// Jump
	if p.Jumping && p.OnGround {
		p.VY = JumpForce
		p.OnGround = false
	}
	
	// Apply terminal velocity
	if p.VY > 20.0 {
		p.VY = 20.0
	}
	if p.VX > PlayerSpeed {
		p.VX = PlayerSpeed
	} else if p.VX < -PlayerSpeed {
		p.VX = -PlayerSpeed
	}
	
	// Reset jumping flag
	p.Jumping = false
	
	// Update mining progress
	if p.Mining && p.MiningTarget != nil {
		p.MiningProgress += deltaTime * 10.0 // Mining speed
		if p.MiningProgress >= 100.0 {
			p.MiningProgress = 0
			p.Mining = false
		}
	} else {
		p.MiningProgress = 0
		p.Mining = false
	}
}

// UpdateWithCollision updates player position with collision detection
// This should be called after Update() with the nearby hexagons
func (p *Player) UpdateWithCollision(deltaTime float64, checkCollision func(float64, float64, float64, float64) bool) {
	// Update position
	p.X += p.VX * deltaTime
	p.Y += p.VY * deltaTime
	
	// Get player bounds
	minX, minY, maxX, maxY := p.GetBounds()
	
	// Check vertical collision (ground detection)
	// Check bottom edge
	bottomLeftCollision := checkCollision(minX, maxY, maxX, maxY+1)
	bottomRightCollision := checkCollision(minX+p.Width/2, maxY, maxX, maxY+1)
	
	if bottomLeftCollision || bottomRightCollision {
		// We hit the ground - snap to it
		p.VY = 0
		p.OnGround = true
		
		// Find the exact ground position by checking from current position upward
		groundY := p.Y
		for checkY := p.Y + p.Height; checkY > p.Y - 10; checkY-- {
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

// TakeDamage reduces the player's health
func (p *Player) TakeDamage(amount float64) {
	p.Health -= amount
	if p.Health < 0 {
		p.Health = 0
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
}

// StopMining stops mining
func (p *Player) StopMining() {
	p.Mining = false
	p.MiningTarget = nil
	p.MiningProgress = 0
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