// Package ui implements the death screen interface
package ui

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// DeathScreen represents the death screen UI
type DeathScreen struct {
	Active      bool
	DeathTime   time.Time
	CauseOfDeath string
	ScreenWidth  int
	ScreenHeight int

	// Button states
	RespawnHovered   bool
	MainMenuHovered  bool

	// Callbacks
	OnRespawn   func()
	OnMainMenu  func()
}

// NewDeathScreen creates a new death screen
func NewDeathScreen(screenWidth, screenHeight int) *DeathScreen {
	return &DeathScreen{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		Active:       false,
	}
}

// Trigger triggers the death screen
func (ds *DeathScreen) Trigger(cause string) {
	ds.Active = true
	ds.DeathTime = time.Now()
	ds.CauseOfDeath = cause
}

// Update handles input for the death screen
func (ds *DeathScreen) Update() error {
	if !ds.Active {
		return nil
	}

	// Get mouse position
	mx, my := ebiten.CursorPosition()

	// Check button hovers
	buttonWidth := 200.0
	buttonHeight := 50.0
	centerX := float64(ds.ScreenWidth) / 2
	respawnY := float64(ds.ScreenHeight)/2 + 50
	menuY := float64(ds.ScreenHeight)/2 + 120

	ds.RespawnHovered = mx >= int(centerX-buttonWidth/2) && mx <= int(centerX+buttonWidth/2) &&
		my >= int(respawnY) && my <= int(respawnY+buttonHeight)
	ds.MainMenuHovered = mx >= int(centerX-buttonWidth/2) && mx <= int(centerX+buttonWidth/2) &&
		my >= int(menuY) && my <= int(menuY+buttonHeight)

	// Handle clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if ds.RespawnHovered && ds.OnRespawn != nil {
			ds.OnRespawn()
			ds.Active = false
		}
		if ds.MainMenuHovered && ds.OnMainMenu != nil {
			ds.OnMainMenu()
			ds.Active = false
		}
	}

	// Also support keyboard shortcuts
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		if ds.OnRespawn != nil {
			ds.OnRespawn()
			ds.Active = false
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if ds.OnMainMenu != nil {
			ds.OnMainMenu()
			ds.Active = false
		}
	}

	return nil
}

// Draw renders the death screen
func (ds *DeathScreen) Draw(screen *ebiten.Image) {
	if !ds.Active {
		return
	}

	// Draw semi-transparent dark overlay
	overlayColor := color.RGBA{20, 0, 0, 230}
	ebitenutil.DrawRect(screen, 0, 0, float64(ds.ScreenWidth), float64(ds.ScreenHeight), overlayColor)

	// Draw "YOU DIED" text (large, centered)
	titleY := float64(ds.ScreenHeight)/2 - 100
	ebitenutil.DebugPrintAt(screen, "YOU DIED", ds.ScreenWidth/2-40, int(titleY))

	// Draw cause of death
	if ds.CauseOfDeath != "" {
		causeText := "Cause: " + ds.CauseOfDeath
		ebitenutil.DebugPrintAt(screen, causeText, ds.ScreenWidth/2-60, int(titleY+40))
	}

	// Draw buttons
	buttonWidth := 200.0
	buttonHeight := 50.0
	centerX := float64(ds.ScreenWidth) / 2

	// Respawn button
	respawnY := float64(ds.ScreenHeight)/2 + 50
	respawnColor := color.RGBA{50, 150, 50, 255}
	if ds.RespawnHovered {
		respawnColor = color.RGBA{70, 200, 70, 255}
	}
	ds.drawButton(screen, centerX-buttonWidth/2, respawnY, buttonWidth, buttonHeight, respawnColor, "RESPAWN (R)")

	// Main Menu button
	menuY := float64(ds.ScreenHeight)/2 + 120
	menuColor := color.RGBA{100, 100, 100, 255}
	if ds.MainMenuHovered {
		menuColor = color.RGBA{150, 150, 150, 255}
	}
	ds.drawButton(screen, centerX-buttonWidth/2, menuY, buttonWidth, buttonHeight, menuColor, "MAIN MENU (ESC)")

	// Draw instructions at bottom
	ebitenutil.DebugPrintAt(screen, "Click a button or press R to respawn, ESC for menu", ds.ScreenWidth/2-150, ds.ScreenHeight-50)
}

// drawButton draws a button with the given properties
func (ds *DeathScreen) drawButton(screen *ebiten.Image, x, y, width, height float64, bgColor color.RGBA, text string) {
	// Button background
	ebitenutil.DrawRect(screen, x, y, width, height, bgColor)

	// Button border
	borderColor := color.RGBA{200, 200, 200, 255}
	ebitenutil.DrawRect(screen, x, y, width, 2, borderColor)
	ebitenutil.DrawRect(screen, x, y+height-2, width, 2, borderColor)
	ebitenutil.DrawRect(screen, x, y, 2, height, borderColor)
	ebitenutil.DrawRect(screen, x+width-2, y, 2, height, borderColor)

	// Button text (centered)
	textX := int(x + width/2 - float64(len(text)*4))
	textY := int(y + height/2 - 6)
	ebitenutil.DebugPrintAt(screen, text, textX, textY)
}

// IsActive returns whether the death screen is currently showing
func (ds *DeathScreen) IsActive() bool {
	return ds.Active
}

// Hide hides the death screen
func (ds *DeathScreen) Hide() {
	ds.Active = false
}

// SetScreenSize updates the screen size
func (ds *DeathScreen) SetScreenSize(width, height int) {
	ds.ScreenWidth = width
	ds.ScreenHeight = height
}
