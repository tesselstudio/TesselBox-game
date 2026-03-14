package menu

import (
	"image/color"
	"math"
	"tesselbox/pkg/hexagon"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// MenuType represents the type of menu screen
type MenuType int

const (
	MenuLogin MenuType = iota
	MenuMain
	MenuSettings
)

// MenuAction represents an action to be taken by the menu
type MenuAction int

const (
	ActionNone MenuAction = iota
	ActionLoginGoogle
	ActionLoginGitHub
	ActionStartGame
	ActionOpenSettings
	ActionOpenLogin
	ActionExit
	ActionBack
	ActionLogout
)

// MenuItem represents a menu option
type MenuItem struct {
	Text     string
	Action   MenuAction
	Position int
	Hovered  bool
	Enabled  bool
}

// Menu represents the main menu system
type Menu struct {
	CurrentMenu  MenuType
	Items        []MenuItem
	SelectedItem int

	// Visual properties
	Title           string
	BackgroundColor color.RGBA
	AccentColor     color.RGBA
	HoverColor      color.RGBA

	// Animation
	animTimer     float64
	rotationAngle float64

	// Transitions
	fadeAlpha     float64
	transitioning bool

	// For solid color drawing
	whiteImage *ebiten.Image
}

// NewMenu creates a new menu system
func NewMenu() *Menu {
	// Create a 1x1 white image for solid color drawing
	whiteImage := ebiten.NewImage(1, 1)
	whiteImage.Fill(color.RGBA{255, 255, 255, 255})

	menu := &Menu{
		CurrentMenu:     MenuMain, // Start with main menu screen
		BackgroundColor: color.RGBA{20, 25, 40, 255},
		AccentColor:     color.RGBA{100, 150, 200, 255},
		HoverColor:      color.RGBA{150, 200, 255, 255},
		animTimer:       0.0,
		rotationAngle:   0.0,
		fadeAlpha:       1.0,
		transitioning:   false,
		whiteImage:      whiteImage,
	}

	menu.SetMainMenu()
	return menu
}

// SetLoginMenu sets up the login menu
func (m *Menu) SetLoginMenu() {
	m.CurrentMenu = MenuLogin
	m.Title = "TESSELBOX"
	m.Items = []MenuItem{
		{Text: "LOGIN WITH GOOGLE", Action: ActionLoginGoogle, Position: 0, Enabled: true},
		{Text: "LOGIN WITH GITHUB", Action: ActionLoginGitHub, Position: 1, Enabled: true},
	}
	m.SelectedItem = 0
}

// SetMainMenu sets up the main menu items
func (m *Menu) SetMainMenu() {
	m.CurrentMenu = MenuMain
	m.Title = "TESSELBOX"
	m.Items = []MenuItem{
		{Text: "START GAME", Action: ActionStartGame, Position: 0, Enabled: true},
		{Text: "SETTINGS", Action: ActionOpenSettings, Position: 1, Enabled: true},
		{Text: "LOGOUT", Action: ActionLogout, Position: 2, Enabled: true},
		{Text: "EXIT", Action: ActionExit, Position: 3, Enabled: true},
	}
	m.SelectedItem = 0
}

// SetSettingsMenu sets up the settings menu
func (m *Menu) SetSettingsMenu() {
	m.CurrentMenu = MenuSettings
	m.Title = "SETTINGS"
	m.Items = []MenuItem{
		{Text: "Graphics: High", Action: ActionNone, Position: 0, Enabled: true},
		{Text: "Sound: On", Action: ActionNone, Position: 1, Enabled: true},
		{Text: "Controls: WASD", Action: ActionNone, Position: 2, Enabled: true},
		{Text: "BACK", Action: ActionBack, Position: 3, Enabled: true},
	}
	m.SelectedItem = 0
}

// Update handles menu input and updates
func (m *Menu) Update() MenuAction {
	// Update animations
	m.animTimer += 0.016 // Approx 60 FPS
	m.rotationAngle += 0.005

	// Handle fade transition
	if m.transitioning {
		m.fadeAlpha -= 0.05
		if m.fadeAlpha <= 0.0 {
			m.fadeAlpha = 0.0
			m.transitioning = false
		}
	}

	// Handle keyboard input
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		m.SelectedItem--
		if m.SelectedItem < 0 {
			m.SelectedItem = len(m.Items) - 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		m.SelectedItem++
		if m.SelectedItem >= len(m.Items) {
			m.SelectedItem = 0
		}
	}

	// Update hovered states
	for i := range m.Items {
		m.Items[i].Hovered = (i == m.SelectedItem)
	}

	// Handle selection
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return m.Items[m.SelectedItem].Action
	}

	// Handle mouse input
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		return m.handleMouseClick(mx, my)
	}

	// Update hover from mouse position
	mx, my := ebiten.CursorPosition()
	m.updateHoverFromMouse(mx, my)

	return ActionNone
}

// handleMouseClick handles mouse clicks on menu items
func (m *Menu) handleMouseClick(mx, my int) MenuAction {
	screenWidth := 1280
	screenHeight := 720

	// Calculate menu item positions - match drawMenuItems exactly
	startY := screenHeight/2 - 100
	itemHeight := 80 // Match drawMenuItems
	itemWidth := 450 // Match drawMenuItems
	centerX := screenWidth / 2
	startX := centerX - itemWidth/2

	for i, item := range m.Items {
		itemY := startY + i*itemHeight

		// Check if click is on this item - use same dimensions as drawing
		if mx >= startX && mx <= startX+itemWidth &&
			my >= itemY && my <= itemY+itemHeight-10 {
			m.SelectedItem = i
			if item.Enabled {
				return item.Action
			}
		}
	}

	return ActionNone
}

// updateHoverFromMouse updates the selected item based on mouse position
func (m *Menu) updateHoverFromMouse(mx, my int) {
	screenWidth := 1280
	screenHeight := 720

	// Calculate menu item positions - match drawMenuItems exactly
	startY := screenHeight/2 - 100
	itemHeight := 80 // Match drawMenuItems
	itemWidth := 450 // Match drawMenuItems
	centerX := screenWidth / 2
	startX := centerX - itemWidth/2

	for i, item := range m.Items {
		itemY := startY + i*itemHeight

		// Check if mouse is hovering over this item - use same dimensions as drawing
		if mx >= startX && mx <= startX+itemWidth &&
			my >= itemY && my <= itemY+itemHeight-10 {
			if item.Enabled {
				m.SelectedItem = i
			}
			break
		}
	}
}

// Draw renders the menu
func (m *Menu) Draw(screen *ebiten.Image) {
	// Draw background
	screen.Fill(m.BackgroundColor)

	// Draw animated hexagonal background pattern
	m.drawHexBackground(screen)

	// Draw menu title
	m.drawTitle(screen)

	// Draw menu items
	m.drawMenuItems(screen)

	// Draw version info
	ebitenutil.DebugPrintAt(screen, "v2.0 - Hexagon Sandbox", 10, 690)
}

// drawHexBackground draws animated hexagonal background pattern
func (m *Menu) drawHexBackground(screen *ebiten.Image) {
	screenWidth := 1280
	screenHeight := 720

	// Draw a grid of hexagons with varying opacity
	hexSize := 40.0
	for y := 0; y < screenHeight+100; y += int(hexSize * 1.7) {
		for x := 0; x < screenWidth+100; x += int(hexSize * 1.5) {
			// Offset every other row
			offsetX := float64(0)
			if (y/int(hexSize*1.7))%2 == 1 {
				offsetX = hexSize * 0.75
			}

			// Calculate opacity based on position and animation
			distance := math.Sqrt(math.Pow(float64(x)-640, 2) + math.Pow(float64(y)-360, 2))
			opacity := math.Max(0, 1.0-distance/800.0)
			opacity *= 0.3 + math.Sin(m.animTimer+distance/100.0)*0.1

			// Draw hexagon
			corners := hexagon.GetHexCorners(float64(x)+offsetX, float64(y), hexSize)
			vertices := make([]ebiten.Vertex, len(corners))
			for i, corner := range corners {
				vertices[i] = ebiten.Vertex{
					DstX:   float32(corner[0]),
					DstY:   float32(corner[1]),
					ColorR: float32(m.AccentColor.R) / 255.0 * float32(opacity),
					ColorG: float32(m.AccentColor.G) / 255.0 * float32(opacity),
					ColorB: float32(m.AccentColor.B) / 255.0 * float32(opacity),
					ColorA: float32(opacity * 0.5),
				}
			}

			indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5}
			screen.DrawTriangles(vertices, indices, m.whiteImage, nil)
		}
	}
}

// drawTitle draws the menu title
func (m *Menu) drawTitle(screen *ebiten.Image) {
	screenWidth := 1280

	// Draw title background hexagon
	titleX := float64(screenWidth) / 2
	titleY := 150.0
	hexSize := 80.0

	// Rotate hexagon
	corners := hexagon.GetHexCorners(titleX, titleY, hexSize)
	rotatedCorners := make([][2]float64, len(corners))
	for i, corner := range corners {
		rotatedCorners[i] = rotatePoint(corner[0], corner[1], titleX, titleY, m.rotationAngle)
	}

	// Draw hexagon background
	vertices := make([]ebiten.Vertex, len(rotatedCorners))
	for i, corner := range rotatedCorners {
		vertices[i] = ebiten.Vertex{
			DstX:   float32(corner[0]),
			DstY:   float32(corner[1]),
			ColorR: float32(m.AccentColor.R) / 255.0,
			ColorG: float32(m.AccentColor.G) / 255.0,
			ColorB: float32(m.AccentColor.B) / 255.0,
			ColorA: 0.8,
		}
	}

	indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5}
	screen.DrawTriangles(vertices, indices, m.whiteImage, nil)

	// Draw title text
	// Using debug print for simplicity - in production, use proper font rendering
	ebitenutil.DebugPrintAt(screen, m.Title, int(titleX)-40, int(titleY)-10)
}

// drawMenuItems draws the menu options with larger text
func (m *Menu) drawMenuItems(screen *ebiten.Image) {
	screenWidth := 1280
	screenHeight := 720

	startY := screenHeight/2 - 100
	itemHeight := 80 // Increased height for larger text
	itemWidth := 450 // Increased width
	centerX := screenWidth / 2
	startX := centerX - itemWidth/2

	for i, item := range m.Items {
		itemY := startY + i*itemHeight

		// Determine color based on state
		bgColor := color.RGBA{50, 60, 80, 200}
		borderColor := m.AccentColor

		if item.Hovered {
			bgColor = color.RGBA{80, 100, 130, 230}
			borderColor = m.HoverColor
		}

		if !item.Enabled {
			bgColor = color.RGBA{40, 40, 50, 150}
			borderColor = color.RGBA{80, 80, 80, 255}
		}

		// Draw hexagon-shaped button background
		m.drawHexButton(screen, float64(startX), float64(itemY), float64(itemWidth), float64(itemHeight-10), bgColor, borderColor)

		// Draw item text with larger size - draw multiple times for thicker text
		// Center text in button
		textX := startX + (itemWidth-len(item.Text)*8)/2
		textY := itemY + 30 // Centered vertically in button

		// Draw text multiple times with slight offsets for thicker appearance
		for dx := 0; dx < 2; dx++ {
			for dy := 0; dy < 2; dy++ {
				ebitenutil.DebugPrintAt(screen, item.Text, textX+dx, textY+dy)
			}
		}
	}
}

// drawHexButton draws a hexagon-shaped button
func (m *Menu) drawHexButton(screen *ebiten.Image, x, y, width, height float64, bgColor, borderColor color.RGBA) {
	// Draw rounded rectangle with hexagon-like corners
	ebitenutil.DrawRect(screen, x+10, y, width-20, height, bgColor)

	// Draw border
	ebitenutil.DrawRect(screen, x+10, y, width-20, 3, borderColor)
	ebitenutil.DrawRect(screen, x+10, y+height-3, width-20, 3, borderColor)

	// Draw decorative hexagon points
	leftPoint := x
	rightPoint := x + width
	centerY := y + height/2

	// Left point
	m.drawHexagonPoint(screen, leftPoint, centerY, borderColor)
	// Right point
	m.drawHexagonPoint(screen, rightPoint, centerY, borderColor)
}

// drawHexagonPoint draws a hexagon point decoration
func (m *Menu) drawHexagonPoint(screen *ebiten.Image, x, y float64, color color.RGBA) {
	pointSize := 10.0
	ebitenutil.DrawRect(screen, x-pointSize/2, y-pointSize/2, pointSize, pointSize, color)
}

// rotatePoint rotates a point around a center
func rotatePoint(x, y, cx, cy, angle float64) [2]float64 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	dx := x - cx
	dy := y - cy

	rotatedX := cx + dx*cos - dy*sin
	rotatedY := cy + dx*sin + dy*cos

	return [2]float64{rotatedX, rotatedY}
}
