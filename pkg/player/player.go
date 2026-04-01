package player

import (
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
	MiningRange  = 2000.0 // Increased range for better block placement
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
	MovingUp    bool
	MovingDown  bool
	Jumping     bool
	OnGround    bool
	IsFlying    bool

	// Mining state
	Mining          bool
	MiningTarget    *world.Hexagon
	MiningProgress  float64
	MiningStartTime time.Time

	// Inventory (reference to inventory)
	SelectedSlot int

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
		SelectedSlot:   0,
		LastUpdateTime: time.Now(),
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

	// Apply gravity (framerate independent) - skip if flying
	if !p.IsFlying {
		p.VY += Gravity * deltaTime * 60 * 60 // Proper scaling for consistent gravity
	}

	// Handle vertical movement when flying
	if p.IsFlying {
		if p.MovingUp {
			p.VY -= PlayerSpeed * deltaTime * 10
		} else if p.MovingDown {
			p.VY += PlayerSpeed * deltaTime * 10
		} else {
			// Apply friction for smooth stopping in air
			p.VY *= Friction
		}
		// Clamp vertical velocity when flying
		if p.VY > TerminalVelX {
			p.VY = TerminalVelX
		} else if p.VY < -TerminalVelX {
			p.VY = -TerminalVelX
		}
		// Stop very small vertical movements to prevent jitter
		if !p.MovingUp && !p.MovingDown && p.VY > -0.1 && p.VY < 0.1 {
			p.VY = 0
		}
	}

	// Clamp vertical velocity (increased for faster falling)
	if p.VY > TerminalVelY {
		p.VY = TerminalVelY
	} else if p.VY < -TerminalVelY {
		p.VY = -TerminalVelY
	}

	// Jump with delta time
	if p.Jumping && p.OnGround {
		p.VY = JumpForce * 60 // Scale jump force to match gravity scaling
		p.OnGround = false
	}

	// Reset jumping flag
	p.Jumping = false
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
	squaredDistance := p.DistanceTo(x, y)
	maxSquaredDistance := MiningRange * MiningRange
	return squaredDistance <= maxSquaredDistance
}

// SetSelectedSlot sets the currently selected inventory slot
func (p *Player) SetSelectedSlot(slot int) {
	p.SelectedSlot = slot
}

// GetSelectedSlot returns the currently selected inventory slot
func (p *Player) GetSelectedSlot() int {
	return p.SelectedSlot
}

// SetFlying sets the flying state
func (p *Player) SetFlying(flying bool) {
	p.IsFlying = flying
	if flying {
		p.OnGround = false // Can't be on ground while flying
	}
}

// IsFlying returns true if the player is flying
func (p *Player) GetIsFlying() bool {
	return p.IsFlying
}

// SetMovingUp sets the upward movement state
func (p *Player) SetMovingUp(up bool) {
	p.MovingUp = up
}

// SetMovingDown sets the downward movement state
func (p *Player) SetMovingDown(down bool) {
	p.MovingDown = down
}
