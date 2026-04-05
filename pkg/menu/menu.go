package menu

import (
	"fmt"
	"image/color"
	"math"
	"tesselbox/pkg/blocks"
	"tesselbox/pkg/hexagon"
	"time"

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
	MenuBlockLibrary
	MenuWorldSelect
	MenuCreateWorld
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
	ActionOpenBlockLibrary
	ActionSelectBlock
	ActionOpenPluginManager
	ActionOpenSkinEditor
	ActionToggleSound
	ActionToggleFullscreen
	ActionChangeResolution
	ActionSelectWorld
	ActionCreateNewWorld
	ActionDeleteWorld
	ActionBackToWorldSelect
)

// WorldInfo represents information about a saved world
type WorldInfo struct {
	Name      string
	LastSaved string
	CreatedAt string
	Seed      int64
	GameMode  string
	Exists    bool
}

// MenuItem represents a menu option
type MenuItem struct {
	Text        string
	Action      MenuAction
	Position    int
	Hovered     bool
	Enabled     bool
	Tooltip     string
	Description string
}

// Menu represents the main menu system
type Menu struct {
	CurrentMenu   MenuType
	Items         []MenuItem
	SelectedItem  int
	SelectedBlock string // For block library menu
	CreativeMode  bool   // Whether creative mode is enabled

	// World selection data
	Worlds        []WorldInfo
	SelectedWorld int

	// Create world data
	NewWorldName string
	NewWorldSeed int64
	NewWorldMode string
	SelectedMode int

	// Scrolling for block library
	ScrollOffset    int
	VisibleItems    int
	MaxVisibleItems int

	// Visual properties
	Title           string
	BackgroundColor color.RGBA
	AccentColor     color.RGBA
	HoverColor      color.RGBA
	TextColor       color.RGBA
	DisabledColor   color.RGBA

	// Animation
	animTimer     float64
	rotationAngle float64

	// Tooltip system
	TooltipVisible     bool
	TooltipText        string
	TooltipX, TooltipY int
	TooltipTimer       float64
	TooltipDelay       float64
	TooltipAlpha       float64

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
		CurrentMenu:     MenuMain,                       // Start with main menu screen
		BackgroundColor: color.RGBA{15, 20, 35, 255},    // Darker blue background
		AccentColor:     color.RGBA{120, 180, 255, 255}, // Brighter blue accent
		HoverColor:      color.RGBA{180, 220, 255, 255}, // Light blue hover
		TextColor:       color.RGBA{255, 255, 255, 255}, // White text
		DisabledColor:   color.RGBA{128, 128, 128, 255}, // Gray disabled
		animTimer:       0.0,
		rotationAngle:   0.0,
		fadeAlpha:       1.0,
		transitioning:   false,
		whiteImage:      whiteImage,
		MaxVisibleItems: 8, // Show max 8 items at once
		ScrollOffset:    0,
		// Tooltip system initialization
		TooltipVisible: false,
		TooltipText:    "",
		TooltipTimer:   0.0,
		TooltipDelay:   0.5, // 0.5 second delay
		TooltipAlpha:   0.0,
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
		{Text: "START GAME", Action: ActionStartGame, Position: 0, Enabled: true,
			Tooltip: "Start a new game session", Description: "Begin your adventure in the hexagon world"},
	}

	if m.CreativeMode {
		m.Items = append(m.Items, MenuItem{Text: "BLOCK LIBRARY", Action: ActionOpenBlockLibrary, Position: len(m.Items), Enabled: true,
			Tooltip: "Browse and select blocks", Description: "Access all available block types for creative mode"})
		m.Items = append(m.Items, MenuItem{Text: "PLUGIN MANAGER", Action: ActionOpenPluginManager, Position: len(m.Items), Enabled: true,
			Tooltip: "Manage plugins and mods", Description: "Browse, install, and manage game plugins"})
		m.Items = append(m.Items, MenuItem{Text: "SKIN EDITOR", Action: ActionOpenSkinEditor, Position: len(m.Items), Enabled: true,
			Tooltip: "Customize player skin", Description: "Create and edit custom player skins"})
	}

	m.Items = append(m.Items, []MenuItem{
		{Text: "SETTINGS", Action: ActionOpenSettings, Position: len(m.Items), Enabled: true,
			Tooltip: "Configure game settings", Description: "Adjust audio, video and gameplay options"},
		{Text: "EXIT", Action: ActionExit, Position: len(m.Items), Enabled: true,
			Tooltip: "Exit the game", Description: "Close TesselBox and return to desktop"},
	}...)
	m.SelectedItem = 0
}

// SetBlockLibraryMenu sets up the block library menu
func (m *Menu) SetBlockLibraryMenu() {
	m.CurrentMenu = MenuBlockLibrary
	m.Title = "BLOCK LIBRARY"
	m.Items = []MenuItem{}
	m.SelectedBlock = ""
	m.ScrollOffset = 0

	// Add all available blocks as menu items
	for blockName := range blocks.BlockDefinitions {
		m.Items = append(m.Items, MenuItem{
			Text:     blockName,
			Action:   ActionNone, // Will handle selection differently
			Position: len(m.Items),
			Enabled:  true,
		})
	}

	if len(m.Items) > 0 {
		m.SelectedItem = 0
		// Calculate visible items
		m.VisibleItems = len(m.Items)
		if m.VisibleItems > m.MaxVisibleItems {
			m.VisibleItems = m.MaxVisibleItems
		}
	}
}

// SetSettingsMenu sets up the settings menu
func (m *Menu) SetSettingsMenu() {
	m.CurrentMenu = MenuSettings
	m.Title = "SETTINGS"
	m.Items = []MenuItem{
		{Text: "TOGGLE SOUND", Action: ActionToggleSound, Position: 0, Enabled: true},
		{Text: "TOGGLE FULLSCREEN", Action: ActionToggleFullscreen, Position: 1, Enabled: true},
		{Text: "CHANGE RESOLUTION", Action: ActionChangeResolution, Position: 2, Enabled: true},
		{Text: "BACK", Action: ActionBack, Position: 3, Enabled: true},
	}
	m.SelectedItem = 0
}

// Update handles menu input and updates
func (m *Menu) Update() MenuAction {
	// Update animations
	m.animTimer += 0.016 // Approx 60 FPS
	m.rotationAngle += 0.005

	// Update tooltips
	m.updateTooltip(0.016)

	// Handle fade transition
	if m.transitioning {
		m.fadeAlpha -= 0.05
		if m.fadeAlpha <= 0.0 {
			m.fadeAlpha = 0.0
			m.transitioning = false
		}
	}

	// Handle ESC key for returning to main menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if m.CurrentMenu == MenuBlockLibrary || m.CurrentMenu == MenuSettings {
			return ActionBack
		}
	}

	// Handle keyboard input
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		m.SelectedItem--
		if m.SelectedItem < 0 {
			m.SelectedItem = len(m.Items) - 1
		}
		// Handle scrolling for block library
		if m.CurrentMenu == MenuBlockLibrary && len(m.Items) > m.MaxVisibleItems {
			if m.SelectedItem < m.ScrollOffset {
				m.ScrollOffset = m.SelectedItem
			} else if m.SelectedItem >= m.ScrollOffset+m.MaxVisibleItems {
				m.ScrollOffset = m.SelectedItem - m.MaxVisibleItems + 1
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		m.SelectedItem++
		if m.SelectedItem >= len(m.Items) {
			m.SelectedItem = 0
			m.ScrollOffset = 0 // Reset scroll when wrapping to top
		}
		// Handle scrolling for block library
		if m.CurrentMenu == MenuBlockLibrary && len(m.Items) > m.MaxVisibleItems {
			if m.SelectedItem < m.ScrollOffset {
				m.ScrollOffset = m.SelectedItem
			} else if m.SelectedItem >= m.ScrollOffset+m.MaxVisibleItems {
				m.ScrollOffset = m.SelectedItem - m.MaxVisibleItems + 1
			}
		}
	}

	// Handle mouse wheel scrolling for block library
	if m.CurrentMenu == MenuBlockLibrary && len(m.Items) > m.MaxVisibleItems {
		_, scrollY := ebiten.Wheel()
		if scrollY > 0 {
			// Scroll up
			m.ScrollOffset--
			if m.ScrollOffset < 0 {
				m.ScrollOffset = 0
			}
			// Adjust selected item if needed
			if m.SelectedItem < m.ScrollOffset {
				m.SelectedItem = m.ScrollOffset
			}
		} else if scrollY < 0 {
			// Scroll down
			maxScroll := len(m.Items) - m.MaxVisibleItems
			if m.ScrollOffset < maxScroll {
				m.ScrollOffset++
			}
			// Adjust selected item if needed
			if m.SelectedItem >= m.ScrollOffset+m.MaxVisibleItems {
				m.SelectedItem = m.ScrollOffset + m.MaxVisibleItems - 1
			}
		}
	}

	// Update hovered states
	for i := range m.Items {
		m.Items[i].Hovered = (i == m.SelectedItem)
	}

	// Handle selection
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if m.CurrentMenu == MenuBlockLibrary {
			// Select the block
			m.SelectedBlock = m.Items[m.SelectedItem].Text
			return ActionSelectBlock
		} else {
			return m.Items[m.SelectedItem].Action
		}
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

	// Determine which items are visible based on scroll offset
	visibleStart := m.ScrollOffset
	visibleEnd := m.ScrollOffset + m.VisibleItems
	if visibleEnd > len(m.Items) {
		visibleEnd = len(m.Items)
	}

	for i := visibleStart; i < visibleEnd; i++ {
		item := m.Items[i]
		visibleIndex := i - m.ScrollOffset
		itemY := startY + visibleIndex*itemHeight

		// Check if click is on this item - use same dimensions as drawing
		if mx >= startX && mx <= startX+itemWidth &&
			my >= itemY && my <= itemY+itemHeight-10 {
			m.SelectedItem = i
			if item.Enabled {
				if m.CurrentMenu == MenuBlockLibrary {
					m.SelectedBlock = item.Text
					return ActionSelectBlock
				} else {
					return item.Action
				}
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

	// Determine which items are visible based on scroll offset
	visibleStart := m.ScrollOffset
	visibleEnd := m.ScrollOffset + m.VisibleItems
	if visibleEnd > len(m.Items) {
		visibleEnd = len(m.Items)
	}

	for i := visibleStart; i < visibleEnd; i++ {
		item := m.Items[i]
		visibleIndex := i - m.ScrollOffset
		itemY := startY + visibleIndex*itemHeight

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

	// Draw tooltips
	m.drawTooltip(screen)

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

	// Determine which items are visible based on scroll offset
	visibleStart := m.ScrollOffset
	visibleEnd := m.ScrollOffset + m.VisibleItems
	if visibleEnd > len(m.Items) {
		visibleEnd = len(m.Items)
	}

	// Draw scroll indicators if needed
	if m.CurrentMenu == MenuBlockLibrary && len(m.Items) > m.MaxVisibleItems {
		// Draw up arrow if not at top
		if m.ScrollOffset > 0 {
			ebitenutil.DebugPrintAt(screen, "▲", centerX-5, startY-30)
		}
		// Draw down arrow if not at bottom
		if m.ScrollOffset < len(m.Items)-m.MaxVisibleItems {
			lastVisibleY := startY + (m.VisibleItems-1)*itemHeight + itemHeight
			ebitenutil.DebugPrintAt(screen, "▼", centerX-5, lastVisibleY+10)
		}
		// Draw scroll indicator text
		scrollText := fmt.Sprintf("%d/%d", m.ScrollOffset+1, len(m.Items))
		ebitenutil.DebugPrintAt(screen, scrollText, centerX-20, startY-50)
	}

	for i := visibleStart; i < visibleEnd; i++ {
		item := m.Items[i]
		visibleIndex := i - m.ScrollOffset
		itemY := startY + visibleIndex*itemHeight

		// Determine color based on state
		bgColor := color.RGBA{40, 60, 90, 230} // Darker, more opaque background
		borderColor := m.AccentColor

		if item.Hovered {
			bgColor = color.RGBA{80, 120, 180, 250} // Brighter hover
			borderColor = m.HoverColor
		}

		if !item.Enabled {
			bgColor = color.RGBA{30, 30, 40, 180} // Darker disabled
			borderColor = m.DisabledColor
		}

		// Draw hexagon-shaped button background
		m.drawHexButton(screen, float64(startX), float64(itemY), float64(itemWidth), float64(itemHeight-10), bgColor, borderColor)

		// Draw item text with better visibility
		// Center text in button
		textX := startX + (itemWidth-len(item.Text)*8)/2
		textY := itemY + 30 // Centered vertically in button

		// Draw text multiple times with slight offsets for thicker appearance
		for dx := 0; dx < 3; dx++ {
			for dy := 0; dy < 3; dy++ {
				ebitenutil.DebugPrintAt(screen, item.Text, textX+dx, textY+dy)
			}
		}

		// Also draw a larger version for better visibility
		ebitenutil.DebugPrintAt(screen, item.Text, textX, textY)
	}
}

// drawHexButton draws a hexagon-shaped button
func (m *Menu) drawHexButton(screen *ebiten.Image, x, y, width, height float64, bgColor, borderColor color.RGBA) {
	// Draw rounded rectangle with hexagon-like corners - make it more visible
	ebitenutil.DrawRect(screen, x+5, y, width-10, height, bgColor)

	// Draw border - make it thicker
	ebitenutil.DrawRect(screen, x+5, y, width-10, 5, borderColor)
	ebitenutil.DrawRect(screen, x+5, y+height-5, width-10, 5, borderColor)

	// Draw decorative hexagon points
	leftPoint := x
	rightPoint := x + width
	centerY := y + height/2

	// Left point - make it larger
	m.drawHexagonPoint(screen, leftPoint, centerY, borderColor)
	// Right point - make it larger
	m.drawHexagonPoint(screen, rightPoint, centerY, borderColor)
}

// drawHexagonPoint draws a hexagon point decoration
func (m *Menu) drawHexagonPoint(screen *ebiten.Image, x, y float64, color color.RGBA) {
	pointSize := 15.0 // Make it larger
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

// drawTooltip renders the tooltip if visible
func (m *Menu) drawTooltip(screen *ebiten.Image) {
	if !m.TooltipVisible || m.TooltipText == "" {
		return
	}

	// Calculate tooltip dimensions
	padding := 10
	maxWidth := 300
	lines := m.wrapText(m.TooltipText, maxWidth)

	lineHeight := 20
	tooltipWidth := maxWidth + padding*2
	tooltipHeight := len(lines)*lineHeight + padding*2

	// Ensure tooltip stays on screen
	tooltipX := m.TooltipX
	tooltipY := m.TooltipY

	if tooltipX+tooltipWidth > 1280 {
		tooltipX = 1280 - tooltipWidth - 10
	}
	if tooltipY+tooltipHeight > 720 {
		tooltipY = 720 - tooltipHeight - 10
	}
	if tooltipX < 0 {
		tooltipX = 10
	}
	if tooltipY < 0 {
		tooltipY = 10
	}

	// Draw tooltip background with alpha
	bgColor := color.RGBA{40, 40, 50, uint8(m.TooltipAlpha * 200)}
	borderColor := color.RGBA{120, 180, 255, uint8(m.TooltipAlpha * 255)}

	m.drawHexButton(screen, float64(tooltipX), float64(tooltipY), float64(tooltipWidth), float64(tooltipHeight), bgColor, borderColor)

	// Draw tooltip text
	for i, line := range lines {
		ebitenutil.DebugPrintAt(screen, line, tooltipX+padding, tooltipY+padding+i*lineHeight)
	}
}

// updateTooltip updates the tooltip state based on hover
func (m *Menu) updateTooltip(deltaTime float64) {
	mx, my := ebiten.CursorPosition()

	// Check if hovering over any menu item
	hoveringItem := false
	for i, item := range m.Items {
		if !item.Enabled || item.Tooltip == "" {
			continue
		}

		// Calculate item position (same as drawMenuItems)
		screenWidth := 1280
		screenHeight := 720
		startY := screenHeight/2 - 100
		itemHeight := 80
		itemWidth := 450
		centerX := screenWidth / 2
		startX := centerX - itemWidth/2

		// Check if this item is visible
		if i < m.ScrollOffset || i >= m.ScrollOffset+m.VisibleItems {
			continue
		}

		visibleIndex := i - m.ScrollOffset
		itemY := startY + visibleIndex*itemHeight

		// Check if mouse is hovering over this item
		if mx >= startX && mx <= startX+itemWidth &&
			my >= itemY && my <= itemY+itemHeight-10 {

			if !m.TooltipVisible || m.TooltipText != item.Tooltip {
				m.TooltipTimer = 0.0
				m.TooltipText = item.Tooltip
				m.TooltipX = mx + 15
				m.TooltipY = my - 30
			}

			m.TooltipTimer += deltaTime
			if m.TooltipTimer >= m.TooltipDelay {
				m.TooltipVisible = true
				// Fade in effect
				if m.TooltipAlpha < 1.0 {
					m.TooltipAlpha = min(1.0, m.TooltipAlpha+deltaTime*3)
				}
			}

			hoveringItem = true
			break
		}
	}

	// Hide tooltip if not hovering
	if !hoveringItem {
		if m.TooltipVisible {
			// Fade out effect
			m.TooltipAlpha -= deltaTime * 3
			if m.TooltipAlpha <= 0 {
				m.TooltipVisible = false
				m.TooltipAlpha = 0
				m.TooltipText = ""
			}
		} else {
			m.TooltipTimer = 0
		}
	}
}

// wrapText wraps text to fit within maxWidth
func (m *Menu) wrapText(text string, maxWidth int) []string {
	words := []string{}
	currentWord := ""

	for _, char := range text {
		if char == ' ' {
			if currentWord != "" {
				words = append(words, currentWord)
				currentWord = ""
			}
		} else {
			currentWord += string(char)
		}
	}
	if currentWord != "" {
		words = append(words, currentWord)
	}

	lines := []string{}
	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		// Rough estimate of text width (6 pixels per character)
		if len(testLine)*6 > maxWidth {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = word
			} else {
				lines = append(lines, word)
			}
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// SetWorldSelectMenu sets up world selection menu
func (m *Menu) SetWorldSelectMenu() {
	m.CurrentMenu = MenuWorldSelect
	m.Title = "SELECT WORLD"

	// Load world list
	worlds, err := m.loadWorldList()
	if err != nil {
		worlds = []WorldInfo{}
	}

	m.Worlds = worlds
	m.SelectedWorld = 0

	// Create menu items
	m.Items = []MenuItem{}
	for i, world := range worlds {
		if world.Exists {
			m.Items = append(m.Items, MenuItem{
				Text:     world.Name,
				Action:   ActionSelectWorld,
				Position: i,
				Enabled:  true,
				Tooltip: fmt.Sprintf("Last saved: %s\nCreated: %s\nSeed: %d\nMode: %s",
					world.LastSaved, world.CreatedAt, world.Seed, world.GameMode),
			})
		} else {
			m.Items = append(m.Items, MenuItem{
				Text:     world.Name,
				Action:   ActionSelectWorld,
				Position: i,
				Enabled:  false,
				Tooltip:  "New world - will be created when selected",
			})
		}
	}

	// Add navigation items
	m.Items = append(m.Items, MenuItem{
		Text:     "CREATE NEW WORLD",
		Action:   ActionCreateNewWorld,
		Position: len(worlds),
		Enabled:  true,
		Tooltip:  "Create a brand new world",
	})

	m.Items = append(m.Items, MenuItem{
		Text:     "BACK",
		Action:   ActionBack,
		Position: len(worlds) + 1,
		Enabled:  true,
		Tooltip:  "Return to main menu",
	})

	m.SelectedItem = 0
}

// UpdateCreateWorldMenu handles input for world creation
func (m *Menu) UpdateCreateWorldMenu() MenuAction {
	// Handle keyboard input
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		m.SelectedItem--
		if m.SelectedItem < 0 {
			m.SelectedItem = 4
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		m.SelectedItem++
		if m.SelectedItem > 4 {
			m.SelectedItem = 0
		}
	}

	// Update menu items
	m.updateCreateWorldItems()

	// Handle selection
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return ActionCreateNewWorld
	}

	// Handle back
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ActionBackToWorldSelect
	}

	return ActionNone
}

// editWorldName handles text input for world name
func (m *Menu) editWorldName() {
	presetNames := []string{"New World", "My World", "Survival Island", "Creative Paradise", "Hexagon World"}
	for i, name := range presetNames {
		if name == m.NewWorldName {
			nextIndex := (i + 1) % len(presetNames)
			m.NewWorldName = presetNames[nextIndex]
			break
		}
	}
}

// editSeed handles seed input and randomization
func (m *Menu) editSeed() {
	if m.NewWorldSeed < 0 {
		m.NewWorldSeed = time.Now().Unix() % 1000000
	} else {
		knownSeeds := []int64{12345, 67890, 42, 1337, 9001, 31415926}
		for i, seed := range knownSeeds {
			if seed == m.NewWorldSeed {
				nextIndex := (i + 1) % len(knownSeeds)
				m.NewWorldSeed = knownSeeds[nextIndex]
				break
			}
		}
	}
}

// updateCreateWorldItems updates the menu items display
func (m *Menu) updateCreateWorldItems() {
	m.Items = []MenuItem{
		{Text: fmt.Sprintf("Name: %s", m.NewWorldName), Action: ActionNone, Position: 0, Enabled: true, Tooltip: "Current: " + m.NewWorldName},
		{Text: fmt.Sprintf("Seed: %d", m.NewWorldSeed), Action: ActionNone, Position: 1, Enabled: true, Tooltip: "Current: " + fmt.Sprintf("%d", m.NewWorldSeed)},
		{Text: fmt.Sprintf("Mode: %s", m.NewWorldMode), Action: ActionNone, Position: 2, Enabled: true, Tooltip: "Current: " + m.NewWorldMode},
		{Text: "CREATE WORLD", Action: ActionCreateNewWorld, Position: 3, Enabled: true, Tooltip: "Create and enter world"},
		{Text: "BACK", Action: ActionBackToWorldSelect, Position: 4, Enabled: true, Tooltip: "Return to world selection"},
	}
	m.SelectedItem = 3
}

// loadWorldList loads the list of saved worlds
func (m *Menu) loadWorldList() ([]WorldInfo, error) {
	// For now, return basic info - in a full implementation,
	// this would read from world metadata files
	worlds := []WorldInfo{
		{
			Name:      "New World",
			LastSaved: "Never",
			CreatedAt: time.Now().Format("2006-01-02"),
			Seed:      time.Now().Unix() % 1000000,
			GameMode:  "Creative",
			Exists:    false,
		},
	}

	return worlds, nil
}
