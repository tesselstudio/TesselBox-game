// Package ui implements damage indicators and floating combat text
package ui

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// DamageTier represents the severity of damage
type DamageTier int

const (
	TierGreen  DamageTier = iota // Low damage
	TierYellow                   // Moderate damage
	TierRed                      // Severe damage
	TierPurple                   // Fatal/Instant death
)

// DamageIndicator represents a floating damage number
type DamageIndicator struct {
	X, Y       float64
	Amount     float64
	Tier       DamageTier
	IsCritical bool
	IsHeal     bool
	StartTime  time.Time
	Duration   time.Duration
	VelocityY  float64
}

// DamageIndicatorManager manages all active damage indicators
type DamageIndicatorManager struct {
	Indicators   []*DamageIndicator
	ScreenWidth  int
	ScreenHeight int
}

// NewDamageIndicatorManager creates a new damage indicator manager
func NewDamageIndicatorManager(screenWidth, screenHeight int) *DamageIndicatorManager {
	return &DamageIndicatorManager{
		Indicators:   make([]*DamageIndicator, 0),
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
	}
}

// SpawnDamageIndicator creates a new damage indicator
func (dim *DamageIndicatorManager) SpawnDamageIndicator(x, y float64, amount float64, tier DamageTier, isCritical bool) {
	indicator := &DamageIndicator{
		X:          x,
		Y:          y,
		Amount:     amount,
		Tier:       tier,
		IsCritical: isCritical,
		IsHeal:     false,
		StartTime:  time.Now(),
		Duration:   2 * time.Second,
		VelocityY:  -30, // Float upward
	}
	dim.Indicators = append(dim.Indicators, indicator)
}

// SpawnHealIndicator creates a healing indicator
func (dim *DamageIndicatorManager) SpawnHealIndicator(x, y float64, amount float64) {
	indicator := &DamageIndicator{
		X:         x,
		Y:         y,
		Amount:    amount,
		Tier:      TierGreen,
		IsHeal:    true,
		StartTime: time.Now(),
		Duration:  2 * time.Second,
		VelocityY: -30,
	}
	dim.Indicators = append(dim.Indicators, indicator)
}

// Update updates all damage indicators
func (dim *DamageIndicatorManager) Update(deltaTime float64) {
	now := time.Now()
	activeIndicators := make([]*DamageIndicator, 0)

	for _, indicator := range dim.Indicators {
		// Move upward
		indicator.Y += indicator.VelocityY * deltaTime

		// Check if expired
		if now.Sub(indicator.StartTime) < indicator.Duration {
			activeIndicators = append(activeIndicators, indicator)
		}
	}

	dim.Indicators = activeIndicators
}

// Draw renders all damage indicators
func (dim *DamageIndicatorManager) Draw(screen *ebiten.Image, cameraX, cameraY float64) {
	for _, indicator := range dim.Indicators {
		screenX := indicator.X - cameraX
		screenY := indicator.Y - cameraY

		// Skip if off screen
		if screenX < -50 || screenX > float64(dim.ScreenWidth)+50 ||
			screenY < -50 || screenY > float64(dim.ScreenHeight)+50 {
			continue
		}

		// Determine color based on tier
		var textColor color.RGBA
		if indicator.IsHeal {
			textColor = color.RGBA{0, 255, 0, 255} // Green for healing
		} else {
			switch indicator.Tier {
			case TierGreen:
				textColor = color.RGBA{0, 255, 0, 255} // Green
			case TierYellow:
				textColor = color.RGBA{255, 255, 0, 255} // Yellow
			case TierRed:
				textColor = color.RGBA{255, 0, 0, 255} // Red
			case TierPurple:
				textColor = color.RGBA{255, 0, 255, 255} // Purple
			}
		}

		// Format text
		var text string
		if indicator.IsHeal {
			text = fmt.Sprintf("+%.0f", indicator.Amount)
		} else {
			text = fmt.Sprintf("-%.0f", indicator.Amount)
		}

		// Draw text background for visibility
		ebitenutil.DrawRect(screen, screenX-2, screenY-2, float64(len(text)*8+4), 14, color.RGBA{0, 0, 0, 150})

		// Draw text (use ebitenutil for simplicity, in production use proper font rendering)
		// Note: ebitenutil.DebugPrintAt doesn't support color, so we use white
		// In a real implementation, you'd use a proper font rendering system
		ebitenutil.DebugPrintAt(screen, text, int(screenX), int(screenY))
		_ = textColor // Would be used with proper font rendering system

		// Draw "CRITICAL!" text for critical hits
		if indicator.IsCritical && !indicator.IsHeal {
			ebitenutil.DebugPrintAt(screen, "CRITICAL!", int(screenX-10), int(screenY-15))
		}
	}
}

// Clear removes all indicators
func (dim *DamageIndicatorManager) Clear() {
	dim.Indicators = make([]*DamageIndicator, 0)
}

// ScreenFlash represents a screen flash effect when player takes damage
type ScreenFlash struct {
	Active    bool
	StartTime time.Time
	Duration  time.Duration
	Color     color.RGBA
	Intensity float64
}

// NewScreenFlash creates a new screen flash
func NewScreenFlash() *ScreenFlash {
	return &ScreenFlash{
		Duration: 500 * time.Millisecond,
		Color:    color.RGBA{255, 0, 0, 100},
	}
}

// Trigger triggers a screen flash
func (sf *ScreenFlash) Trigger(col color.RGBA, duration time.Duration) {
	sf.Active = true
	sf.StartTime = time.Now()
	sf.Color = col
	sf.Duration = duration
}

// Update updates the screen flash
func (sf *ScreenFlash) Update() {
	if sf.Active && time.Since(sf.StartTime) > sf.Duration {
		sf.Active = false
	}
}

// Draw renders the screen flash overlay
func (sf *ScreenFlash) Draw(screen *ebiten.Image, screenWidth, screenHeight int) {
	if !sf.Active {
		return
	}

	// Calculate fade based on time
	elapsed := time.Since(sf.StartTime).Seconds()
	total := sf.Duration.Seconds()
	progress := elapsed / total
	if progress > 1 {
		progress = 1
	}

	// Fade out
	alpha := uint8(float64(sf.Color.A) * (1 - progress))
	flashColor := color.RGBA{sf.Color.R, sf.Color.G, sf.Color.B, alpha}

	// Draw overlay
	ebitenutil.DrawRect(screen, 0, 0, float64(screenWidth), float64(screenHeight), flashColor)
}

// HitDirection represents the direction a hit came from
type HitDirection int

const (
	DirFront HitDirection = iota
	DirBack
	DirLeft
	DirRight
)

// DirectionalHitIndicator shows an arrow indicating hit direction
type DirectionalHitIndicator struct {
	Direction HitDirection
	StartTime time.Time
	Duration  time.Duration
}

// DirectionalHitManager manages directional hit indicators
type DirectionalHitManager struct {
	Indicators []*DirectionalHitIndicator
}

// NewDirectionalHitManager creates a new manager
func NewDirectionalHitManager() *DirectionalHitManager {
	return &DirectionalHitManager{
		Indicators: make([]*DirectionalHitIndicator, 0),
	}
}

// TriggerHit triggers a directional hit indicator
func (dhm *DirectionalHitManager) TriggerHit(direction HitDirection) {
	indicator := &DirectionalHitIndicator{
		Direction: direction,
		StartTime: time.Now(),
		Duration:  1 * time.Second,
	}
	dhm.Indicators = append(dhm.Indicators, indicator)
}

// Update updates all indicators
func (dhm *DirectionalHitManager) Update() {
	now := time.Now()
	active := make([]*DirectionalHitIndicator, 0)

	for _, indicator := range dhm.Indicators {
		if now.Sub(indicator.StartTime) < indicator.Duration {
			active = append(active, indicator)
		}
	}

	dhm.Indicators = active
}

// Draw renders directional indicators around the screen edges
func (dhm *DirectionalHitManager) Draw(screen *ebiten.Image, screenWidth, screenHeight int) {
	centerX := float64(screenWidth) / 2
	centerY := float64(screenHeight) / 2

	for _, indicator := range dhm.Indicators {
		// Calculate position based on direction
		var x, y float64
		switch indicator.Direction {
		case DirFront:
			x, y = centerX, float64(screenHeight)-50
		case DirBack:
			x, y = centerX, 50
		case DirLeft:
			x, y = 50, centerY
		case DirRight:
			x, y = float64(screenWidth)-50, centerY
		}

		// Draw arrow (simplified as colored rectangle)
		arrowColor := color.RGBA{255, 0, 0, 200}
		ebitenutil.DrawRect(screen, x-10, y-10, 20, 20, arrowColor)

		// Direction label
		var label string
		switch indicator.Direction {
		case DirFront:
			label = "^"
		case DirBack:
			label = "v"
		case DirLeft:
			label = "<"
		case DirRight:
			label = ">"
		}
		ebitenutil.DebugPrintAt(screen, label, int(x-3), int(y-5))
	}
}
