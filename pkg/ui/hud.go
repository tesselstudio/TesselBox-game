// Package ui implements the Heads-Up Display for survival mode
package ui

import (
	"fmt"
	"image/color"

	"tesselbox/pkg/equipment"
	"tesselbox/pkg/gametime"
	"tesselbox/pkg/health"
	"tesselbox/pkg/survival"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// HUD represents the survival mode heads-up display
type HUD struct {
	ScreenWidth  int
	ScreenHeight int

	// System references
	SurvivalManager *survival.SurvivalManager
	EquipmentSet    *equipment.EquipmentSet
	HealthSystem    *health.LocationalHealthSystem
	DayNightCycle   *gametime.DayNightCycle

	// Bar positions
	HealthBarX, HealthBarY   float64
	HungerBarX, HungerBarY   float64
	ThirstBarX, ThirstBarY   float64
	StaminaBarX, StaminaBarY float64

	// Bar dimensions
	BarWidth   float64
	BarHeight  float64
	BarSpacing float64

	// Icon sizes
	IconSize float64
}

// NewHUD creates a new HUD instance
func NewHUD(screenWidth, screenHeight int, sm *survival.SurvivalManager, es *equipment.EquipmentSet, hs *health.LocationalHealthSystem, dnc *gametime.DayNightCycle) *HUD {
	hud := &HUD{
		ScreenWidth:     screenWidth,
		ScreenHeight:    screenHeight,
		SurvivalManager: sm,
		EquipmentSet:    es,
		HealthSystem:    hs,
		DayNightCycle:   dnc,

		BarWidth:   120,
		BarHeight:  12,
		BarSpacing: 4,
		IconSize:   16,
	}

	hud.calculateLayout()
	return hud
}

// calculateLayout positions the HUD elements
func (h *HUD) calculateLayout() {
	// Position bars at bottom-left corner
	startX := 10.0
	startY := float64(h.ScreenHeight) - 80.0

	h.HealthBarX = startX
	h.HealthBarY = startY
	h.HungerBarX = startX
	h.HungerBarY = startY + h.BarHeight + h.BarSpacing
	h.ThirstBarX = startX
	h.ThirstBarY = startY + 2*(h.BarHeight+h.BarSpacing)
	h.StaminaBarX = startX
	h.StaminaBarY = startY + 3*(h.BarHeight+h.BarSpacing)
}

// Draw renders the HUD
func (h *HUD) Draw(screen *ebiten.Image) {
	// Draw health bar (always shown if health system exists)
	h.drawHealthBar(screen)

	// Draw survival bars (only if survival manager exists)
	if h.SurvivalManager != nil {
		h.drawHungerBar(screen)
		h.drawThirstBar(screen)
		h.drawStaminaBar(screen)
	}

	// Draw day/night indicator
	h.drawDayNightIndicator(screen)

	// Draw defense indicator
	h.drawDefenseIndicator(screen)

	// Draw body part health diagram
	h.drawBodyPartHealth(screen)
}

// drawHealthBar draws the segmented health bar
func (h *HUD) drawHealthBar(screen *ebiten.Image) {
	// Background
	bgColor := color.RGBA{30, 30, 30, 200}
	ebitenutil.DrawRect(screen, h.HealthBarX, h.HealthBarY, h.BarWidth, h.BarHeight, bgColor)

	// Health fill
	if h.HealthSystem != nil {
		healthPct := h.HealthSystem.GetOverallHealthPercentage()
		fillWidth := h.BarWidth * healthPct

		// Color based on health level
		var fillColor color.RGBA
		if healthPct > 0.6 {
			fillColor = color.RGBA{220, 50, 50, 255} // Red
		} else if healthPct > 0.3 {
			fillColor = color.RGBA{220, 100, 50, 255} // Orange-red
		} else {
			fillColor = color.RGBA{150, 30, 30, 255} // Dark red (critical)
		}

		ebitenutil.DrawRect(screen, h.HealthBarX, h.HealthBarY, fillWidth, h.BarHeight, fillColor)

		// Draw segments for body parts
		segmentWidth := h.BarWidth / 6
		for i := 0; i < 6; i++ {
			x := h.HealthBarX + float64(i)*segmentWidth
			ebitenutil.DrawRect(screen, x, h.HealthBarY, 1, h.BarHeight, color.RGBA{0, 0, 0, 100})
		}
	}

	// Icon
	h.drawIcon(screen, h.HealthBarX-18, h.HealthBarY, color.RGBA{220, 50, 50, 255}, "HP")

	// Value text
	if h.HealthSystem != nil {
		text := fmt.Sprintf("%.0f", h.HealthSystem.OverallHealth)
		ebitenutil.DebugPrintAt(screen, text, int(h.HealthBarX+h.BarWidth+5), int(h.HealthBarY))
	}
}

// drawHungerBar draws the hunger bar
func (h *HUD) drawHungerBar(screen *ebiten.Image) {
	// Background
	bgColor := color.RGBA{30, 30, 30, 200}
	ebitenutil.DrawRect(screen, h.HungerBarX, h.HungerBarY, h.BarWidth, h.BarHeight, bgColor)

	// Hunger fill
	hungerPct := h.SurvivalManager.Hunger / h.SurvivalManager.MaxHunger
	fillWidth := h.BarWidth * hungerPct

	// Color based on hunger level (inverse - full is good)
	var fillColor color.RGBA
	if hungerPct > 0.6 {
		fillColor = color.RGBA{139, 90, 43, 255} // Brown (meat color)
	} else if hungerPct > 0.3 {
		fillColor = color.RGBA{160, 100, 50, 255} // Lighter brown
	} else {
		fillColor = color.RGBA{200, 50, 50, 255} // Red (starving)
	}

	ebitenutil.DrawRect(screen, h.HungerBarX, h.HungerBarY, fillWidth, h.BarHeight, fillColor)

	// Icon
	h.drawIcon(screen, h.HungerBarX-18, h.HungerBarY, color.RGBA{139, 90, 43, 255}, "HU")

	// Value text
	text := fmt.Sprintf("%.0f", h.SurvivalManager.Hunger)
	ebitenutil.DebugPrintAt(screen, text, int(h.HungerBarX+h.BarWidth+5), int(h.HungerBarY))

	// Starving indicator
	if h.SurvivalManager.IsStarving {
		ebitenutil.DebugPrintAt(screen, "STARVING!", int(h.HungerBarX), int(h.HungerBarY-12))
	}
}

// drawThirstBar draws the thirst bar
func (h *HUD) drawThirstBar(screen *ebiten.Image) {
	// Background
	bgColor := color.RGBA{30, 30, 30, 200}
	ebitenutil.DrawRect(screen, h.ThirstBarX, h.ThirstBarY, h.BarWidth, h.BarHeight, bgColor)

	// Thirst fill
	thirstPct := h.SurvivalManager.Thirst / h.SurvivalManager.MaxThirst
	fillWidth := h.BarWidth * thirstPct

	// Color based on thirst level (inverse - full is good)
	var fillColor color.RGBA
	if thirstPct > 0.6 {
		fillColor = color.RGBA{50, 100, 220, 255} // Blue
	} else if thirstPct > 0.3 {
		fillColor = color.RGBA{80, 120, 200, 255} // Lighter blue
	} else {
		fillColor = color.RGBA{150, 50, 50, 255} // Red (dehydrated)
	}

	ebitenutil.DrawRect(screen, h.ThirstBarX, h.ThirstBarY, fillWidth, h.BarHeight, fillColor)

	// Icon
	h.drawIcon(screen, h.ThirstBarX-18, h.ThirstBarY, color.RGBA{50, 100, 220, 255}, "TH")

	// Value text
	text := fmt.Sprintf("%.0f", h.SurvivalManager.Thirst)
	ebitenutil.DebugPrintAt(screen, text, int(h.ThirstBarX+h.BarWidth+5), int(h.ThirstBarY))

	// Dehydrated indicator
	if h.SurvivalManager.IsDehydrated {
		ebitenutil.DebugPrintAt(screen, "DEHYDRATED!", int(h.ThirstBarX), int(h.ThirstBarY-12))
	}
}

// drawStaminaBar draws the stamina bar
func (h *HUD) drawStaminaBar(screen *ebiten.Image) {
	// Background
	bgColor := color.RGBA{30, 30, 30, 200}
	ebitenutil.DrawRect(screen, h.StaminaBarX, h.StaminaBarY, h.BarWidth, h.BarHeight, bgColor)

	// Stamina fill
	staminaPct := h.SurvivalManager.Stamina / h.SurvivalManager.MaxStamina
	fillWidth := h.BarWidth * staminaPct

	// Color based on stamina level
	var fillColor color.RGBA
	if staminaPct > 0.6 {
		fillColor = color.RGBA{220, 220, 50, 255} // Yellow
	} else if staminaPct > 0.3 {
		fillColor = color.RGBA{200, 180, 50, 255} // Darker yellow
	} else {
		fillColor = color.RGBA{150, 100, 30, 255} // Brown/orange (exhausted)
	}

	ebitenutil.DrawRect(screen, h.StaminaBarX, h.StaminaBarY, fillWidth, h.BarHeight, fillColor)

	// Icon
	h.drawIcon(screen, h.StaminaBarX-18, h.StaminaBarY, color.RGBA{220, 220, 50, 255}, "ST")

	// Value text
	text := fmt.Sprintf("%.0f", h.SurvivalManager.Stamina)
	ebitenutil.DebugPrintAt(screen, text, int(h.StaminaBarX+h.BarWidth+5), int(h.StaminaBarY))
}

// drawDayNightIndicator draws the sun/moon indicator
func (h *HUD) drawDayNightIndicator(screen *ebiten.Image) {
	if h.DayNightCycle == nil {
		return
	}

	// Position at top-right
	x := float64(h.ScreenWidth) - 40
	y := 10.0
	size := 24.0

	// Determine if day or night
	isDay := h.DayNightCycle.AmbientLight > 0.5

	// Draw circle background
	bgColor := color.RGBA{30, 30, 30, 200}
	ebitenutil.DrawRect(screen, x, y, size, size, bgColor)

	// Draw sun or moon
	if isDay {
		// Sun (yellow circle)
		sunColor := color.RGBA{255, 255, 0, 255}
		ebitenutil.DrawRect(screen, x+4, y+4, size-8, size-8, sunColor)
	} else {
		// Moon (white crescent - simplified as white circle)
		moonColor := color.RGBA{200, 200, 220, 255}
		ebitenutil.DrawRect(screen, x+4, y+4, size-8, size-8, moonColor)
	}

	// Draw time text
	timeText := fmt.Sprintf("%.0f%%", h.DayNightCycle.AmbientLight*100)
	ebitenutil.DebugPrintAt(screen, timeText, int(x-10), int(y+size+2))
}

// drawDefenseIndicator shows total armor defense
func (h *HUD) drawDefenseIndicator(screen *ebiten.Image) {
	if h.EquipmentSet == nil {
		return
	}

	// Position at bottom-right
	x := float64(h.ScreenWidth) - 80
	y := float64(h.ScreenHeight) - 30

	// Calculate total defense
	defense := h.EquipmentSet.GetTotalDefense()
	speedMod, _, _ := h.EquipmentSet.GetMovementModifiers()

	// Draw defense value
	defenseText := fmt.Sprintf("DEF: %.1f", defense)
	ebitenutil.DebugPrintAt(screen, defenseText, int(x), int(y))

	// Draw movement modifier
	moveText := fmt.Sprintf("SPD: %.0f%%", speedMod*100)
	ebitenutil.DebugPrintAt(screen, moveText, int(x), int(y+12))
}

// drawIcon draws a simple icon for a bar
func (h *HUD) drawIcon(screen *ebiten.Image, x, y float64, col color.RGBA, label string) {
	size := h.IconSize
	// Background
	ebitenutil.DrawRect(screen, x, y, size, size, color.RGBA{50, 50, 50, 200})
	// Icon color
	ebitenutil.DrawRect(screen, x+2, y+2, size-4, size-4, col)
}

// Update recalculates layout if screen size changed
func (h *HUD) Update(screenWidth, screenHeight int) {
	if h.ScreenWidth != screenWidth || h.ScreenHeight != screenHeight {
		h.ScreenWidth = screenWidth
		h.ScreenHeight = screenHeight
		h.calculateLayout()
	}
}
