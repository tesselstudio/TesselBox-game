package gui

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"tesselbox/pkg/config"
)

// MenuScene represents the main menu screen using pixel art GUI
type MenuScene struct {
	// Screen state
	currentScreen string
	shouldExit    bool

	// World data
	worlds        []string
	worldSeeds    map[string]int64
	selectedWorld int

	// Create world form
	newWorldName string
	newWorldSeed string
	cursorField  int // 0 = name, 1 = seed

	// Settings
	soundEnabled    bool
	graphicsQuality string

	// Skin editor
	skinColors     []string
	selectedColor  int
	skinCanvas     [16][16]color.Color
	skinBrushColor color.Color
	skinTool       string
	customSkins    []CustomSkin
	selectedSkin   int
	skinPreview    *ebiten.Image

	// Plugins
	plugins        []PluginInfo
	selectedPlugin int

	// UI Components
	fontManager *FontManager
	buttons     []*PixelButton
	panels      []*PixelPanel
	toggles     []*PixelToggle

	// Selection cursor for keyboard navigation
	cursor int

	// Selected world data to pass to game
	SelectedWorldData WorldSelection
}

// PluginInfo holds plugin information for the plugin manager
type PluginInfo struct {
	Name        string
	Version     string
	Enabled     bool
	Description string
}

// WorldSelection holds the world data when starting the game
type WorldSelection struct {
	WorldName   string
	Seed        int64
	ShouldStart bool
}

// CustomSkin represents a saved custom skin with name and pixel data
type CustomSkin struct {
	Name   string
	Pixels [16][16]color.Color
}

// NewMenuScene creates a new pixel art menu scene
func NewMenuScene() *MenuScene {
	m := &MenuScene{
		currentScreen:   "main",
		soundEnabled:    true,
		graphicsQuality: "medium",
		skinColors:      []string{"Default", "Red", "Blue", "Green", "Purple"},
		selectedColor:   0,
		cursor:          0,
		fontManager:     NewFontManager(),
	}

	// Initialize with default worlds
	m.worlds = []string{"World 1", "World 2"}
	m.worldSeeds = map[string]int64{
		"World 1": 12345,
		"World 2": 67890,
	}

	// Initialize plugins
	m.plugins = []PluginInfo{
		{Name: "Minimap", Version: "1.0", Enabled: true, Description: "Shows a minimap"},
		{Name: "Auto-Save", Version: "2.1", Enabled: true, Description: "Auto-saves every 5 minutes"},
		{Name: "Debug Tools", Version: "0.5", Enabled: false, Description: "Developer debugging tools"},
	}

	// Initialize skin editor
	m.skinBrushColor = ColorPrimary
	m.skinTool = "brush"
	m.customSkins = []CustomSkin{
		{Name: "Default", Pixels: m.createDefaultSkinPixels()},
		{Name: "Steve", Pixels: m.createSteveSkinPixels()},
	}
	m.selectedSkin = 0
	m.skinCanvas = m.customSkins[0].Pixels
	m.skinPreview = ebiten.NewImage(64, 64)

	m.setupMainMenu()

	return m
}

// setupMainMenu creates the main menu buttons
func (m *MenuScene) setupMainMenu() {
	m.buttons = []*PixelButton{}
	m.panels = []*PixelPanel{}
	m.toggles = []*PixelToggle{}

	switch m.currentScreen {
	case "main":
		m.setupMainScreen()
	case "worlds":
		m.setupWorldsScreen()
	case "create_world":
		m.setupCreateWorldScreen()
	case "settings":
		m.setupSettingsScreen()
	case "skin_editor":
		m.setupSkinEditorScreen()
	case "plugins":
		m.setupPluginsScreen()
	}
}

// setupMainScreen creates the main menu buttons
func (m *MenuScene) setupMainScreen() {
	menuItems := []string{"Play", "Create New World", "Skin Editor", "Plugin Manager", "Settings", "Exit"}
	centerX := 400.0
	startY := 220.0

	for i, item := range menuItems {
		btn := NewPixelButton(centerX-100, startY+float64(i)*70, 200, 45, item, func(idx int) func() {
			return func() {
				m.handleMainMenuSelection(idx)
			}
		}(i))

		if i == m.cursor {
			btn.State = ButtonStateHover
		}

		m.buttons = append(m.buttons, btn)
	}

	// Add title panel (positioned above buttons with more space)
	titlePanel := NewPixelPanel(centerX-180, 50, 360, 120, "")
	m.panels = append(m.panels, titlePanel)
}

// setupWorldsScreen creates the world selection screen
func (m *MenuScene) setupWorldsScreen() {
	centerX := 400.0
	startY := 150.0

	// World list
	for i, world := range m.worlds {
		seed := m.worldSeeds[world]
		text := world + " (seed: " + strconv.FormatInt(seed, 10) + ")"

		btn := NewPixelButton(centerX-150, startY+float64(i)*65, 300, 50, text, func(idx int) func() {
			return func() {
				m.selectedWorld = idx
				m.SelectedWorldData = WorldSelection{
					WorldName:   m.worlds[idx],
					Seed:        m.worldSeeds[m.worlds[idx]],
					ShouldStart: true,
				}
				m.shouldExit = true
			}
		}(i))

		if i == m.cursor {
			btn.State = ButtonStateHover
		}

		m.buttons = append(m.buttons, btn)
	}

	// Delete button (with more spacing)
	deleteBtn := NewPixelButton(centerX-100, startY+float64(len(m.worlds))*65+30, 200, 45, "Delete World", func() {
		if m.selectedWorld >= 0 && m.selectedWorld < len(m.worlds) {
			m.deleteWorld(m.selectedWorld)
		}
	})
	m.buttons = append(m.buttons, deleteBtn)

	// Back button
	backBtn := NewPixelButton(centerX-100, startY+float64(len(m.worlds))*65+90, 200, 45, "Back", func() {
		m.currentScreen = "main"
		m.cursor = 0
		m.setupMainMenu()
	})
	m.buttons = append(m.buttons, backBtn)
}

// setupCreateWorldScreen creates the world creation form
func (m *MenuScene) setupCreateWorldScreen() {
	centerX := 400.0
	startY := 180.0

	// Name field panel (larger to fit both fields)
	m.panels = append(m.panels, NewPixelPanel(centerX-220, startY-30, 440, 160, ""))

	// Create button
	createBtn := NewPixelButton(centerX-100, startY+150, 200, 45, "Create", func() {
		if m.newWorldName != "" {
			m.createWorld()
		}
	})
	m.buttons = append(m.buttons, createBtn)

	// Back button
	backBtn := NewPixelButton(centerX-100, startY+220, 200, 45, "Back", func() {
		m.currentScreen = "worlds"
		m.cursor = 0
		m.setupMainMenu()
	})
	m.buttons = append(m.buttons, backBtn)
}

// setupSettingsScreen creates the settings screen
func (m *MenuScene) setupSettingsScreen() {
	centerX := 400.0
	startY := 180.0

	// Sound toggle
	soundText := "Sound: OFF"
	if m.soundEnabled {
		soundText = "Sound: ON"
	}
	soundBtn := NewPixelButton(centerX-100, startY, 200, 45, soundText, func() {
		m.soundEnabled = !m.soundEnabled
		m.setupMainMenu()
	})
	m.buttons = append(m.buttons, soundBtn)

	// Graphics quality button
	graphicsBtn := NewPixelButton(centerX-100, startY+70, 200, 45, "Graphics: "+m.graphicsQuality, func() {
		switch m.graphicsQuality {
		case "low":
			m.graphicsQuality = "medium"
		case "medium":
			m.graphicsQuality = "high"
		case "high":
			m.graphicsQuality = "low"
		}
		m.setupMainMenu()
	})
	m.buttons = append(m.buttons, graphicsBtn)

	// Back button
	backBtn := NewPixelButton(centerX-100, startY+160, 200, 45, "Back", func() {
		m.currentScreen = "main"
		m.cursor = 0
		m.setupMainMenu()
	})
	m.buttons = append(m.buttons, backBtn)
}

// setupSkinEditorScreen creates the skin editor drawing board screen
func (m *MenuScene) setupSkinEditorScreen() {
	// Main drawing canvas panel (left side)
	canvasPanel := NewPixelPanel(50, 80, 320, 320, "")
	m.panels = append(m.panels, canvasPanel)

	// Tools and colors panel (right side)
	toolsPanel := NewPixelPanel(400, 80, 350, 400, "")
	m.panels = append(m.panels, toolsPanel)

	// Tool selection buttons
	toolY := 100.0
	brushBtn := NewPixelButton(420, toolY, 100, 35, "Brush", func() {
		m.skinTool = "brush"
		m.setupMainMenu()
	})
	if m.skinTool == "brush" {
		brushBtn.BgColor = ColorHighlight
	}
	m.buttons = append(m.buttons, brushBtn)

	eraserBtn := NewPixelButton(530, toolY, 100, 35, "Eraser", func() {
		m.skinTool = "eraser"
		m.setupMainMenu()
	})
	if m.skinTool == "eraser" {
		eraserBtn.BgColor = ColorHighlight
	}
	m.buttons = append(m.buttons, eraserBtn)

	fillBtn := NewPixelButton(640, toolY, 100, 35, "Fill", func() {
		m.skinTool = "fill"
		m.setupMainMenu()
	})
	if m.skinTool == "fill" {
		fillBtn.BgColor = ColorHighlight
	}
	m.buttons = append(m.buttons, fillBtn)

	// Color palette
	colors := []color.Color{
		ColorPrimary,                   // Green
		color.RGBA{255, 0, 0, 255},     // Red
		color.RGBA{0, 0, 255, 255},     // Blue
		color.RGBA{255, 255, 0, 255},   // Yellow
		color.RGBA{255, 165, 0, 255},   // Orange
		color.RGBA{128, 0, 128, 255},   // Purple
		color.RGBA{0, 0, 0, 255},       // Black
		color.RGBA{255, 255, 255, 255}, // White
	}
	colorY := 150.0
	for i, c := range colors {
		col := c // capture for closure
		x := 420 + float64(i%4)*80
		y := colorY + float64(i/4)*45
		btn := NewPixelButton(x, y, 70, 35, "", func() {
			m.skinBrushColor = col
			m.setupMainMenu()
		})
		btn.BgColor = c
		m.buttons = append(m.buttons, btn)
	}

	// Skin slot buttons (custom skins)
	skinY := 260.0
	for i := 0; i < 4; i++ {
		idx := i
		var label string
		if idx < len(m.customSkins) {
			label = m.customSkins[idx].Name
		} else {
			label = fmt.Sprintf("Slot %d", idx+1)
		}
		skinBtn := NewPixelButton(420, skinY+float64(i)*45, 150, 35, label, func() {
			if idx < len(m.customSkins) {
				// Load existing skin
				m.selectedSkin = idx
				m.skinCanvas = m.customSkins[idx].Pixels
			}
		})
		if i == m.selectedSkin {
			skinBtn.BgColor = ColorHighlight
		}
		m.buttons = append(m.buttons, skinBtn)
	}

	// Save skin button
	saveBtn := NewPixelButton(580, skinY, 150, 35, "Save Skin", func() {
		m.saveCurrentSkin()
	})
	m.buttons = append(m.buttons, saveBtn)

	// New skin button
	newBtn := NewPixelButton(580, skinY+45, 150, 35, "New Skin", func() {
		m.createNewSkin()
	})
	m.buttons = append(m.buttons, newBtn)

	// Clear canvas button
	clearBtn := NewPixelButton(580, skinY+90, 150, 35, "Clear", func() {
		m.clearSkinCanvas()
	})
	m.buttons = append(m.buttons, clearBtn)

	// Back button
	backBtn := NewPixelButton(400, 500, 200, 45, "Back", func() {
		m.currentScreen = "main"
		m.cursor = 0
		m.setupMainMenu()
	})
	m.buttons = append(m.buttons, backBtn)
}

// setupPluginsScreen creates the plugin manager screen
func (m *MenuScene) setupPluginsScreen() {
	centerX := 400.0
	startY := 120.0

	// Plugin toggles
	for i, plugin := range m.plugins {
		status := "[OFF]"
		if plugin.Enabled {
			status = "[ON]"
		}
		text := plugin.Name + " " + status + " v" + plugin.Version

		btn := NewPixelButton(centerX-180, startY+float64(i)*70, 360, 55, text, func(idx int) func() {
			return func() {
				m.plugins[idx].Enabled = !m.plugins[idx].Enabled
				m.setupMainMenu()
			}
		}(i))

		m.buttons = append(m.buttons, btn)
	}

	// Back button
	backBtn := NewPixelButton(centerX-100, startY+float64(len(m.plugins))*70+30, 200, 45, "Back", func() {
		m.currentScreen = "main"
		m.cursor = 0
		m.setupMainMenu()
	})
	m.buttons = append(m.buttons, backBtn)
}

// handleMainMenuSelection handles main menu item selection
func (m *MenuScene) handleMainMenuSelection(index int) {
	switch index {
	case 0: // Play
		m.currentScreen = "worlds"
		m.cursor = 0
		m.setupMainMenu()
	case 1: // Create New World
		m.currentScreen = "create_world"
		m.newWorldName = ""
		m.newWorldSeed = ""
		m.cursorField = 0
		m.setupMainMenu()
	case 2: // Skin Editor
		m.currentScreen = "skin_editor"
		m.setupMainMenu()
	case 3: // Plugin Manager
		m.currentScreen = "plugins"
		m.setupMainMenu()
	case 4: // Settings
		m.currentScreen = "settings"
		m.setupMainMenu()
	case 5: // Exit
		m.shouldExit = true
	}
}

// createWorld creates a new world with the given name and seed
func (m *MenuScene) createWorld() {
	var seed int64
	if m.newWorldSeed == "" {
		seed = time.Now().UnixNano()
	} else {
		parsed, err := strconv.ParseInt(m.newWorldSeed, 10, 64)
		if err != nil {
			seed = time.Now().UnixNano()
		} else {
			seed = parsed
		}
	}

	m.worlds = append(m.worlds, m.newWorldName)
	m.worldSeeds[m.newWorldName] = seed

	m.currentScreen = "worlds"
	m.setupMainMenu()
}

// deleteWorld removes a world from the list
func (m *MenuScene) deleteWorld(index int) {
	if index < 0 || index >= len(m.worlds) {
		return
	}

	worldToDelete := m.worlds[index]

	// Remove from worlds slice
	m.worlds = append(m.worlds[:index], m.worlds[index+1:]...)

	// Remove from seeds map
	delete(m.worldSeeds, worldToDelete)

	// Reset cursor
	if m.cursor >= len(m.worlds) {
		m.cursor = 0
	}

	m.setupMainMenu()
}

// Layout implements ebiten.Game interface
func (m *MenuScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}

// Update handles input and updates the menu state
func (m *MenuScene) Update() error {
	// Handle quit
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if m.currentScreen == "main" {
			m.shouldExit = true
		} else {
			m.currentScreen = "main"
			m.cursor = 0
			m.setupMainMenu()
		}
	}

	// Handle keyboard navigation
	if m.currentScreen != "create_world" {
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
			if m.cursor > 0 {
				m.cursor--
				m.setupMainMenu()
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			if m.cursor < len(m.buttons)-1 {
				m.cursor++
				m.setupMainMenu()
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			if m.cursor >= 0 && m.cursor < len(m.buttons) {
				m.buttons[m.cursor].OnClick()
			}
		}
	} else {
		// Create world screen text input
		if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
			m.cursorField = 1 - m.cursorField
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			m.createWorld()
		}

		// Handle text input
		if m.cursorField == 0 {
			m.newWorldName = m.handleTextInput(m.newWorldName, 32)
		} else {
			m.newWorldSeed = m.handleTextInput(m.newWorldSeed, 20)
		}
	}

	// Update all buttons and check for mouse hover
	mx, my := ebiten.CursorPosition()
	for i, btn := range m.buttons {
		// Check if mouse is hovering over this button
		if float64(mx) >= btn.X && float64(mx) <= btn.X+btn.Width &&
			float64(my) >= btn.Y && float64(my) <= btn.Y+btn.Height {
			m.cursor = i
		}
		btn.Update()
	}

	if m.shouldExit {
		return ebiten.Termination
	}

	return nil
}

// handleTextInput processes text input for form fields
func (m *MenuScene) handleTextInput(current string, maxLen int) string {
	// This is a simplified version - in production, use ebiten's input handling
	return current
}

// Draw renders the menu scene
func (m *MenuScene) Draw(screen *ebiten.Image) {
	// Clear screen with background color
	screen.Fill(ColorBackground)

	// Draw panels first (background layer)
	for _, panel := range m.panels {
		panel.Draw(screen, m.fontManager)
	}

	// Draw title for main menu (on top of title panel)
	if m.currentScreen == "main" {
		// Draw pixel-art style title using large blocks
		m.drawPixelTitle(screen, "TESSELBOX", 400, 85, ColorPrimary)
		m.fontManager.DrawTextCentered(screen, "Sandbox Game", 400, 135, ColorTextDim, "normal")
	}

	// Draw buttons (top layer)
	for i, btn := range m.buttons {
		// Highlight cursor selection (but don't override pressed state)
		if i == m.cursor && btn.State != ButtonStatePressed {
			btn.State = ButtonStateHover
		}
		btn.Draw(screen, m.fontManager)
	}

	// Draw create world form
	if m.currentScreen == "create_world" {
		m.drawCreateWorldForm(screen)
	}

	// Draw help text
	m.drawHelpText(screen)
}

// drawCreateWorldForm draws the world creation form
func (m *MenuScene) drawCreateWorldForm(screen *ebiten.Image) {
	centerX := 400.0
	startY := 200.0

	// Name field
	labelColor := ColorText
	if m.cursorField == 0 {
		labelColor = ColorPrimary
	}
	m.fontManager.DrawText(screen, "World Name:", centerX-150, startY, labelColor, "normal")

	// Input box
	vector.DrawFilledRect(screen, float32(centerX-150), float32(startY+15), 300, 30, ColorPanel, false)
	vector.StrokeRect(screen, float32(centerX-150), float32(startY+15), 300, 30, 1, ColorBorder, false)

	// Input text
	inputText := m.newWorldName
	if m.cursorField == 0 && time.Now().UnixNano()%1000000000 < 500000000 {
		inputText += "_"
	}
	m.fontManager.DrawText(screen, inputText, centerX-140, startY+22, ColorText, "normal")

	// Seed field
	seedColor := ColorText
	if m.cursorField == 1 {
		seedColor = ColorPrimary
	}
	m.fontManager.DrawText(screen, "Seed (optional):", centerX-150, startY+70, seedColor, "normal")

	// Seed input box
	vector.DrawFilledRect(screen, float32(centerX-150), float32(startY+85), 300, 30, ColorPanel, false)
	vector.StrokeRect(screen, float32(centerX-150), float32(startY+85), 300, 30, 1, ColorBorder, false)

	// Seed input text
	seedText := m.newWorldSeed
	if seedText == "" {
		seedText = "(random)"
	}
	if m.cursorField == 1 && time.Now().UnixNano()%1000000000 < 500000000 {
		seedText += "_"
	}
	m.fontManager.DrawText(screen, seedText, centerX-140, startY+92, ColorText, "normal")
}

// drawHelpText draws keyboard help at the bottom
func (m *MenuScene) drawHelpText(screen *ebiten.Image) {
	helpText := ""
	switch m.currentScreen {
	case "main":
		helpText = "↑/↓: Navigate  Enter: Select  ESC: Quit"
	case "create_world":
		helpText = "Tab: Switch fields  Enter: Create  ESC: Cancel"
	default:
		helpText = "↑/↓: Navigate  Enter: Select  ESC: Back"
	}

	m.fontManager.DrawTextCentered(screen, helpText, 400, 550, ColorTextDim, "small")
}

// ShouldExit returns true if the menu should close
func (m *MenuScene) ShouldExit() bool {
	return m.shouldExit
}

// ShouldStartGame returns true if the game should start
func (m *MenuScene) ShouldStartGame() bool {
	return m.SelectedWorldData.ShouldStart
}

// GetWorldSelection returns the selected world data
func (m *MenuScene) GetWorldSelection() WorldSelection {
	return m.SelectedWorldData
}

// GetTesselboxDir returns the storage directory
func getTesselboxDir() string {
	return config.GetTesselboxDir()
}

// createDefaultSkinPixels creates a default green character skin
func (m *MenuScene) createDefaultSkinPixels() [16][16]color.Color {
	var pixels [16][16]color.Color
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			pixels[y][x] = color.RGBA{34, 139, 34, 255} // Forest green
		}
	}
	return pixels
}

// createSteveSkinPixels creates a Steve-like skin with face and body
func (m *MenuScene) createSteveSkinPixels() [16][16]color.Color {
	var pixels [16][16]color.Color
	// Skin color (face/hands)
	skinColor := color.RGBA{255, 224, 189, 255}
	// Shirt color (blue)
	shirtColor := color.RGBA{65, 105, 225, 255}
	// Pants color (dark blue)
	pantsColor := color.RGBA{25, 25, 112, 255}

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			switch {
			case y < 4: // Head/face
				pixels[y][x] = skinColor
			case y < 10: // Shirt/body
				pixels[y][x] = shirtColor
			default: // Pants/legs
				pixels[y][x] = pantsColor
			}
		}
	}
	return pixels
}

// saveCurrentSkin saves the current canvas to custom skins
func (m *MenuScene) saveCurrentSkin() {
	if m.selectedSkin < len(m.customSkins) {
		m.customSkins[m.selectedSkin].Pixels = m.skinCanvas
	}
}

// createNewSkin creates a new empty custom skin
func (m *MenuScene) createNewSkin() {
	newSkin := CustomSkin{
		Name:   fmt.Sprintf("Skin %d", len(m.customSkins)+1),
		Pixels: m.createDefaultSkinPixels(),
	}
	m.customSkins = append(m.customSkins, newSkin)
	m.selectedSkin = len(m.customSkins) - 1
	m.skinCanvas = newSkin.Pixels
	m.setupMainMenu()
}

// clearSkinCanvas clears the current skin canvas
func (m *MenuScene) clearSkinCanvas() {
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			m.skinCanvas[y][x] = color.RGBA{255, 255, 255, 255} // White
		}
	}
}

// drawPixelTitle draws a large pixel-art style title using blocks
func (m *MenuScene) drawPixelTitle(screen *ebiten.Image, text string, x, y float64, clr color.Color) {
	// Simple 5x5 pixel font for each letter
	letterWidth := 24
	letterHeight := 32
	spacing := 4

	startX := x - float64(len(text)*letterWidth)/2 + float64(len(text)*spacing)/2

	for i, ch := range text {
		lx := startX + float64(i*(letterWidth+spacing))
		switch ch {
		case 'T':
			m.drawLetterT(screen, lx, y, letterWidth, letterHeight, clr)
		case 'E':
			m.drawLetterE(screen, lx, y, letterWidth, letterHeight, clr)
		case 'S':
			m.drawLetterS(screen, lx, y, letterWidth, letterHeight, clr)
		case 'L':
			m.drawLetterL(screen, lx, y, letterWidth, letterHeight, clr)
		case 'B':
			m.drawLetterB(screen, lx, y, letterWidth, letterHeight, clr)
		case 'O':
			m.drawLetterO(screen, lx, y, letterWidth, letterHeight, clr)
		case 'X':
			m.drawLetterX(screen, lx, y, letterWidth, letterHeight, clr)
		}
	}
}

// Pixel letter drawing functions
func (m *MenuScene) drawLetterT(screen *ebiten.Image, x, y float64, w, h int, clr color.Color) {
	// Top bar
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h/5), clr, false)
	// Vertical bar
	vector.DrawFilledRect(screen, float32(x)+float32(w)/2-3, float32(y), 6, float32(h), clr, false)
}

func (m *MenuScene) drawLetterE(screen *ebiten.Image, x, y float64, w, h int, clr color.Color) {
	// Left vertical
	vector.DrawFilledRect(screen, float32(x), float32(y), 6, float32(h), clr, false)
	// Top bar
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h/5), clr, false)
	// Middle bar
	vector.DrawFilledRect(screen, float32(x), float32(y+float64(h)*2/5), float32(w), float32(h/5), clr, false)
	// Bottom bar
	vector.DrawFilledRect(screen, float32(x), float32(y+float64(h)*4/5), float32(w), float32(h/5), clr, false)
}

func (m *MenuScene) drawLetterS(screen *ebiten.Image, x, y float64, w, h int, clr color.Color) {
	// Top bar
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w-3), float32(h/5), clr, false)
	// Top left vertical
	vector.DrawFilledRect(screen, float32(x), float32(y), 6, float32(h/2), clr, false)
	// Middle bar
	vector.DrawFilledRect(screen, float32(x), float32(y+float64(h)*2/5), float32(w), float32(h/5), clr, false)
	// Bottom right vertical
	vector.DrawFilledRect(screen, float32(x)+float32(w)-6, float32(y)+float32(h)*2/5, 6, float32(h)/2, clr, false)
	// Bottom bar
	vector.DrawFilledRect(screen, float32(x+3), float32(y+float64(h)*4/5), float32(w-3), float32(h/5), clr, false)
}

func (m *MenuScene) drawLetterL(screen *ebiten.Image, x, y float64, w, h int, clr color.Color) {
	// Left vertical
	vector.DrawFilledRect(screen, float32(x), float32(y), 6, float32(h), clr, false)
	// Bottom bar
	vector.DrawFilledRect(screen, float32(x), float32(y+float64(h)*4/5), float32(w), float32(h/5), clr, false)
}

func (m *MenuScene) drawLetterB(screen *ebiten.Image, x, y float64, w, h int, clr color.Color) {
	// Left vertical
	vector.DrawFilledRect(screen, float32(x), float32(y), 6, float32(h), clr, false)
	// Top bar
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w-3), float32(h/5), clr, false)
	// Middle bar
	vector.DrawFilledRect(screen, float32(x), float32(y+float64(h)*2/5), float32(w-6), float32(h/5), clr, false)
	// Bottom bar
	vector.DrawFilledRect(screen, float32(x), float32(y+float64(h)*4/5), float32(w-3), float32(h/5), clr, false)
	// Right top vertical
	vector.DrawFilledRect(screen, float32(x)+float32(w)-6, float32(y), 6, float32(h)/2-2, clr, false)
	// Right bottom vertical
	vector.DrawFilledRect(screen, float32(x)+float32(w)-6, float32(y)+float32(h)*2/5+2, 6, float32(h)/2-2, clr, false)
}

func (m *MenuScene) drawLetterO(screen *ebiten.Image, x, y float64, w, h int, clr color.Color) {
	// Left vertical
	vector.DrawFilledRect(screen, float32(x), float32(y+3), 6, float32(h-6), clr, false)
	// Right vertical
	vector.DrawFilledRect(screen, float32(x)+float32(w)-6, float32(y)+3, 6, float32(h)-6, clr, false)
	// Top bar
	vector.DrawFilledRect(screen, float32(x+3), float32(y), float32(w-6), 3, clr, false)
	// Bottom bar
	vector.DrawFilledRect(screen, float32(x+3), float32(y+float64(h)-3), float32(w-6), 3, clr, false)
}

func (m *MenuScene) drawLetterX(screen *ebiten.Image, x, y float64, w, h int, clr color.Color) {
	// Draw X using diagonal strokes - two crossing diagonal bars
	thickness := float32(6)
	// Diagonal from top-left to bottom-right (use line segments)
	vector.StrokeLine(screen, float32(x), float32(y), float32(x)+float32(w), float32(y)+float32(h), thickness, clr, false)
	// Diagonal from top-right to bottom-left
	vector.StrokeLine(screen, float32(x)+float32(w), float32(y), float32(x), float32(y)+float32(h), thickness, clr, false)
}
