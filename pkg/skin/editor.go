package skin

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// SkinEditor represents the skin editing interface
type SkinEditor struct {
	// Skin data
	skinData     *SkinData
	originalSkin *SkinData
	previewSkin  *SkinData

	// UI state
	editorMode    EditorMode
	selectedTool  Tool
	selectedColor color.RGBA
	brushSize     int
	brushOpacity  float64
	brushFlow     int
	blendMode     BlendMode

	// Canvas state
	canvasX, canvasY int
	canvasSize       int
	pixelSize        int
	zoomLevel        float64

	// Preview state
	previewX, previewY int
	previewSize        int
	previewAnimating   bool
	previewRotation    float64

	// Color palette
	palette         []color.RGBA
	paletteX        int
	paletteY        int
	selectedPalette int

	// Tool panel
	toolPanelX int
	toolPanelY int

	// Professional UI layout
	toolbarHeight   int
	sidebarWidth    int
	statusBarHeight int

	// Qt-inspired UI components
	mainWindow    *QtWindow
	menuBar       *QtMenuBar
	toolBar       *QtToolBar
	dockWidgets   []*QtDockWidget
	statusBar     *QtStatusBar
	centralWidget *QtCentralWidget

	// History for undo/redo
	history      []*SkinData
	historyIndex int
	maxHistory   int

	// Layers system
	layers       []*Layer
	activeLayer  int
	layerOpacity float64
	layerBlend   BlendMode
	layerVisible []bool

	// Drawing state
	isDrawing  bool
	lastPixelX int
	lastPixelY int

	// Animation
	animTimer   float64
	cursorBlink bool

	// File operations
	skinDirectory string
	currentSkin   string

	// Visual properties
	backgroundColor color.RGBA
	gridColor       color.RGBA
	selectionColor  color.RGBA

	// For solid color drawing
	whiteImage *ebiten.Image
}

// EditorMode represents different editing modes
type EditorMode int

const (
	ModeEdit EditorMode = iota
	ModePreview
	ModeColorPicker
)

// BlendMode represents different blending modes
type BlendMode int

const (
	BlendNormal BlendMode = iota
	BlendMultiply
	BlendScreen
	BlendOverlay
)

// Tool represents different editing tools
type Tool int

const (
	ToolPencil Tool = iota
	ToolEraser
	ToolFill
	ToolEyedropper
	ToolLine
	ToolRectangle
	ToolCircle
	ToolAirbrush
	ToolSpray
	ToolGradient
	ToolSelect
	ToolMove
	ToolSymmetry
)

// SkinData represents the player's skin data
type SkinData struct {
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	Width     int               `json:"width"`
	Height    int               `json:"height"`
	Pixels    [][]color.RGBA    `json:"pixels"`
	Metadata  map[string]string `json:"metadata"`
}

// Layer represents a drawing layer
type Layer struct {
	Name      string         `json:"name"`
	Visible   bool           `json:"visible"`
	Opacity   float64        `json:"opacity"`
	BlendMode BlendMode      `json:"blendMode"`
	Pixels    [][]color.RGBA `json:"pixels"`
	Locked    bool           `json:"locked"`
}

// Qt-inspired UI components for Ebiten

// QtWindow represents the main application window
type QtWindow struct {
	Title     string
	Width     int
	Height    int
	Widgets   []QtWidget
	MenuBar   *QtMenuBar
	ToolBar   *QtToolBar
	StatusBar *QtStatusBar
}

// QtMenuBar represents a menu bar
type QtMenuBar struct {
	Menus []QtMenu
}

// QtMenu represents a menu
type QtMenu struct {
	Title   string
	Actions []QtAction
}

// QtAction represents a menu action
type QtAction struct {
	Text     string
	Shortcut string
	Callback func()
}

// QtToolBar represents a toolbar
type QtToolBar struct {
	Actions             []QtAction
	X, Y, Width, Height int
	Movable             bool
}

// QtDockWidget represents a dockable panel
type QtDockWidget struct {
	Title               string
	Content             QtWidget
	X, Y, Width, Height int
	Floating            bool
	Visible             bool
}

// QtStatusBar represents a status bar
type QtStatusBar struct {
	Messages            []string
	X, Y, Width, Height int
}

// QtCentralWidget represents the main content area
type QtCentralWidget struct {
	Content             QtWidget
	X, Y, Width, Height int
}

// QtWidget represents a UI widget
type QtWidget interface {
	Draw(screen *ebiten.Image)
	Update() error
	HandleInput(mx, my int) bool
	GetArea() (int, int, int, int)
}

// QtButton represents a button widget
type QtButton struct {
	Text                string
	X, Y, Width, Height int
	Pressed             bool
	Hovered             bool
	Callback            func()
	Style               QtButtonStyle
}

// QtButtonStyle defines button appearance
type QtButtonStyle struct {
	BackgroundColor color.RGBA
	HoverColor      color.RGBA
	PressedColor    color.RGBA
	TextColor       color.RGBA
	BorderColor     color.RGBA
}

// QtPanel represents a panel widget
type QtPanel struct {
	X, Y, Width, Height int
	Children            []QtWidget
	Style               QtPanelStyle
}

// QtPanelStyle defines panel appearance
type QtPanelStyle struct {
	BackgroundColor color.RGBA
	BorderColor     color.RGBA
	BorderWidth     int
}

// QtLabel represents a text label
type QtLabel struct {
	Text                string
	X, Y, Width, Height int
	Style               QtLabelStyle
}

// QtLabelStyle defines label appearance
type QtLabelStyle struct {
	TextColor color.RGBA
	FontSize  int
}

// SkinConfig represents saved skin configuration
type SkinConfig struct {
	CurrentSkin string    `json:"currentSkin"`
	Skins       []string  `json:"skins"`
	LastUsed    time.Time `json:"lastUsed"`
}

const (
	SkinWidth   = 64
	SkinHeight  = 64
	CanvasSize  = 512
	PreviewSize = 128
	PaletteSize = 16
	MaxHistory  = 50
)

// NewSkinEditor creates a new skin editor
func NewSkinEditor() *SkinEditor {
	log.Printf("Creating new skin editor...")

	// Create a 1x1 white image for solid color drawing
	whiteImage := ebiten.NewImage(1, 1)
	whiteImage.Fill(color.RGBA{255, 255, 255, 255})

	editor := &SkinEditor{
		editorMode:       ModeEdit,
		selectedTool:     ToolPencil,
		selectedColor:    color.RGBA{255, 255, 255, 255},
		brushSize:        1,
		brushOpacity:     1.0,
		brushFlow:        100,
		blendMode:        BlendNormal,
		canvasX:          50,
		canvasY:          80, // Leave space for toolbar
		canvasSize:       CanvasSize,
		pixelSize:        CanvasSize / SkinWidth,
		zoomLevel:        1.0,
		previewX:         600,
		previewY:         100,
		previewSize:      PreviewSize,
		previewAnimating: true,
		paletteX:         50,
		paletteY:         600,
		toolPanelX:       600,
		toolPanelY:       400,
		toolbarHeight:    60,
		sidebarWidth:     200,
		statusBarHeight:  30,
		history:          make([]*SkinData, 0),
		historyIndex:     -1,
		maxHistory:       MaxHistory,
		layers:           make([]*Layer, 0),
		activeLayer:      0,
		layerOpacity:     1.0,
		layerBlend:       BlendNormal,
		layerVisible:     make([]bool, 0),
		skinDirectory:    "skins",
		backgroundColor:  color.RGBA{30, 30, 40, 255},
		gridColor:        color.RGBA{60, 60, 80, 255},
		selectionColor:   color.RGBA{100, 150, 255, 255},
		whiteImage:       whiteImage,
	}

	log.Printf("Initializing palette...")
	// Initialize default palette
	editor.initializePalette()

	log.Printf("Initializing layers...")
	// Initialize layers
	editor.initializeLayers()

	log.Printf("Initializing Qt UI...")
	// Initialize Qt-inspired UI components
	editor.initializeQtUI()

	log.Printf("Creating default skin...")
	// Create default skin
	editor.createDefaultSkin()

	log.Printf("Loading skin configuration...")
	// Load saved skins
	if err := editor.loadSkinConfig(); err != nil {
		log.Printf("Warning: Failed to load skin config: %v", err)
	}

	log.Printf("Skin editor created successfully")
	return editor
}

// initializePalette sets up the default color palette
func (se *SkinEditor) initializePalette() {
	se.palette = []color.RGBA{
		// Basic colors
		{0, 0, 0, 255},       // Black
		{255, 255, 255, 255}, // White
		{128, 128, 128, 255}, // Gray
		{192, 192, 192, 255}, // Light gray
		{64, 64, 64, 255},    // Dark gray

		// Primary colors
		{255, 0, 0, 255},   // Red
		{255, 128, 0, 255}, // Orange
		{255, 255, 0, 255}, // Yellow
		{0, 255, 0, 255},   // Green
		{0, 255, 255, 255}, // Cyan
		{0, 0, 255, 255},   // Blue
		{128, 0, 255, 255}, // Purple
		{255, 0, 255, 255}, // Magenta

		// Extended colors
		{165, 42, 42, 255},   // Brown
		{255, 192, 203, 255}, // Pink
		{0, 128, 0, 255},     // Dark green
		{255, 165, 0, 255},   // Gold
		{0, 128, 128, 255},   // Teal
		{128, 0, 128, 255},   // Maroon

		// Skin tones
		{255, 220, 177, 255}, // Light skin
		{255, 198, 144, 255}, // Medium skin
		{255, 160, 122, 255}, // Dark skin
		{222, 184, 135, 255}, // Tan
		{210, 180, 140, 255}, // Beige

		// Custom color slots (empty for user customization)
		{255, 255, 255, 255}, // Custom 1
		{255, 255, 255, 255}, // Custom 2
		{255, 255, 255, 255}, // Custom 3
		{255, 255, 255, 255}, // Custom 4
	}
}

// initializeLayers sets up the default layers
func (se *SkinEditor) initializeLayers() {
	// Create base layer
	baseLayer := &Layer{
		Name:      "Base",
		Visible:   true,
		Opacity:   1.0,
		BlendMode: BlendNormal,
		Pixels:    make([][]color.RGBA, SkinHeight),
		Locked:    false,
	}

	// Initialize base layer pixels
	for y := 0; y < SkinHeight; y++ {
		baseLayer.Pixels[y] = make([]color.RGBA, SkinWidth)
		for x := 0; x < SkinWidth; x++ {
			baseLayer.Pixels[y][x] = color.RGBA{0, 0, 0, 0} // Transparent
		}
	}

	se.layers = []*Layer{baseLayer}
	se.activeLayer = 0
	se.layerVisible = []bool{true}

	log.Printf("Layers initialized with %d layers", len(se.layers))
}

// initializeQtUI sets up the Qt-inspired UI components
func (se *SkinEditor) initializeQtUI() {
	se.mainWindow = se.createMainWindow()
	se.menuBar = se.createMenuBar()
	se.toolBar = se.createToolBar()
	se.dockWidgets = se.createDockWidgets()
	se.statusBar = se.createStatusBar()
	se.centralWidget = se.createCentralWidget()

	log.Printf("Qt UI components initialized")
}

// createMainWindow creates the main window
func (se *SkinEditor) createMainWindow() *QtWindow {
	return &QtWindow{
		Title:   "TesselBox Skin Editor",
		Width:   1280,
		Height:  720,
		Widgets: make([]QtWidget, 0),
	}
}

// createMenuBar creates the menu bar
func (se *SkinEditor) createMenuBar() *QtMenuBar {
	return &QtMenuBar{
		Menus: []QtMenu{},
	}
}

// createToolBar creates the toolbar
func (se *SkinEditor) createToolBar() *QtToolBar {
	return &QtToolBar{
		Actions: []QtAction{},
		X:       0, Y: 60, Width: 1280, Height: 40,
		Movable: true,
	}
}

// createDockWidgets creates dockable widgets
func (se *SkinEditor) createDockWidgets() []*QtDockWidget {
	return make([]*QtDockWidget, 0)
}

// createStatusBar creates the status bar
func (se *SkinEditor) createStatusBar() *QtStatusBar {
	return &QtStatusBar{
		Messages: []string{"Ready"},
		X:        0, Y: 690, Width: 1280, Height: 30,
	}
}

// createCentralWidget creates the central widget
func (se *SkinEditor) createCentralWidget() *QtCentralWidget {
	return &QtCentralWidget{
		X: 200, Y: 100, Width: 880, Height: 590,
	}
}

// addLayer adds a new layer
func (se *SkinEditor) addLayer(name string) {
	newLayer := &Layer{
		Name:      name,
		Visible:   true,
		Opacity:   1.0,
		BlendMode: BlendNormal,
		Pixels:    make([][]color.RGBA, SkinHeight),
		Locked:    false,
	}

	// Initialize layer pixels
	for y := 0; y < SkinHeight; y++ {
		newLayer.Pixels[y] = make([]color.RGBA, SkinWidth)
		for x := 0; x < SkinWidth; x++ {
			newLayer.Pixels[y][x] = color.RGBA{0, 0, 0, 0} // Transparent
		}
	}

	se.layers = append(se.layers, newLayer)
	se.layerVisible = append(se.layerVisible, true)
	se.activeLayer = len(se.layers) - 1

	log.Printf("Added layer: %s (total: %d)", name, len(se.layers))
}

// deleteLayer removes the current layer
func (se *SkinEditor) deleteLayer() {
	if len(se.layers) <= 1 {
		return // Cannot delete the last layer
	}

	if se.activeLayer >= 0 && se.activeLayer < len(se.layers) {
		// Remove layer
		se.layers = append(se.layers[:se.activeLayer], se.layers[se.activeLayer+1:]...)
		se.layerVisible = append(se.layerVisible[:se.activeLayer], se.layerVisible[se.activeLayer+1:]...)

		// Adjust active layer
		if se.activeLayer >= len(se.layers) {
			se.activeLayer = len(se.layers) - 1
		}

		log.Printf("Deleted layer, total: %d", len(se.layers))
	}
}

// setActiveLayer sets the active layer
func (se *SkinEditor) setActiveLayer(index int) {
	if index >= 0 && index < len(se.layers) {
		se.activeLayer = index
		log.Printf("Active layer set to: %s", se.layers[index].Name)
	}
}

// blendColors blends two colors with opacity
func (se *SkinEditor) blendColors(base, overlay color.RGBA, opacity float64) color.RGBA {
	r := float64(base.R)*(1.0-opacity) + float64(overlay.R)*opacity
	g := float64(base.G)*(1.0-opacity) + float64(overlay.G)*opacity
	b := float64(base.B)*(1.0-opacity) + float64(overlay.B)*opacity
	a := float64(base.A)*(1.0-opacity) + float64(overlay.A)*opacity

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}

// mergeLayers combines all visible layers into the skin data
func (se *SkinEditor) mergeLayers() {
	// Clear skin data
	for y := 0; y < SkinHeight; y++ {
		for x := 0; x < SkinWidth; x++ {
			se.skinData.Pixels[y][x] = color.RGBA{0, 0, 0, 0} // Transparent
		}
	}

	// Merge layers from bottom to top
	for i, layer := range se.layers {
		if !se.layerVisible[i] || !layer.Visible {
			continue
		}

		for y := 0; y < SkinHeight; y++ {
			for x := 0; x < SkinWidth; x++ {
				layerPixel := layer.Pixels[y][x]
				if layerPixel.A == 0 {
					continue // Skip transparent pixels
				}

				// Apply layer opacity and blend mode
				blended := se.blendColors(
					se.skinData.Pixels[y][x],
					layerPixel,
					layer.Opacity,
				)
				se.skinData.Pixels[y][x] = blended
			}
		}
	}
}

// createDefaultSkin creates a default player skin
func (se *SkinEditor) createDefaultSkin() {
	log.Printf("Creating default skin...")

	skin := &SkinData{
		Name:      "Default",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Width:     SkinWidth,
		Height:    SkinHeight,
		Pixels:    make([][]color.RGBA, SkinHeight),
		Metadata:  make(map[string]string),
	}

	// Initialize pixels with transparent background
	for y := 0; y < SkinHeight; y++ {
		skin.Pixels[y] = make([]color.RGBA, SkinWidth)
		for x := 0; x < SkinWidth; x++ {
			skin.Pixels[y][x] = color.RGBA{0, 0, 0, 0} // Transparent
		}
	}

	// Draw a simple default skin (humanoid figure)
	se.drawDefaultHumanoid(skin)

	se.skinData = skin
	se.originalSkin = se.copySkin(skin)
	se.previewSkin = se.copySkin(skin)
	se.addToHistory(skin)

	log.Printf("Default skin created successfully")
}

// drawDefaultHumanoid draws a simple humanoid figure on the skin
func (se *SkinEditor) drawDefaultHumanoid(skin *SkinData) {
	// Head (light skin tone)
	headColor := color.RGBA{255, 220, 177, 255}
	for y := 8; y < 16; y++ {
		for x := 24; x < 40; x++ {
			if x >= 0 && x < SkinWidth && y >= 0 && y < SkinHeight {
				skin.Pixels[y][x] = headColor
			}
		}
	}

	// Eyes
	eyeColor := color.RGBA{0, 0, 0, 255}
	skin.Pixels[10][28] = eyeColor
	skin.Pixels[10][36] = eyeColor

	// Body (blue shirt)
	bodyColor := color.RGBA{0, 100, 200, 255}
	for y := 16; y < 32; y++ {
		for x := 20; x < 44; x++ {
			if x >= 0 && x < SkinWidth && y >= 0 && y < SkinHeight {
				skin.Pixels[y][x] = bodyColor
			}
		}
	}

	// Arms
	armColor := color.RGBA{255, 220, 177, 255}
	// Left arm
	for y := 16; y < 28; y++ {
		for x := 16; x < 20; x++ {
			if x >= 0 && x < SkinWidth && y >= 0 && y < SkinHeight {
				skin.Pixels[y][x] = armColor
			}
		}
	}
	// Right arm
	for y := 16; y < 28; y++ {
		for x := 44; x < 48; x++ {
			if x >= 0 && x < SkinWidth && y >= 0 && y < SkinHeight {
				skin.Pixels[y][x] = armColor
			}
		}
	}

	// Legs (brown pants)
	legColor := color.RGBA{139, 69, 19, 255}
	// Left leg
	for y := 32; y < 48; y++ {
		for x := 24; x < 32; x++ {
			if x >= 0 && x < SkinWidth && y >= 0 && y < SkinHeight {
				skin.Pixels[y][x] = legColor
			}
		}
	}
	// Right leg
	for y := 32; y < 48; y++ {
		for x := 32; x < 40; x++ {
			if x >= 0 && x < SkinWidth && y >= 0 && y < SkinHeight {
				skin.Pixels[y][x] = legColor
			}
		}
	}
}

// Update handles skin editor updates
func (se *SkinEditor) Update() error {
	// Update animations
	se.animTimer += 0.016
	if se.animTimer > 0.5 {
		se.animTimer = 0
		se.cursorBlink = !se.cursorBlink
	}

	// Update preview rotation
	if se.previewAnimating {
		se.previewRotation += 0.02
	}

	// Handle input based on mode
	switch se.editorMode {
	case ModeEdit:
		se.handleEditInput()
	case ModePreview:
		se.handlePreviewInput()
	case ModeColorPicker:
		se.handleColorPickerInput()
	}

	return nil
}

// isOverUI checks if mouse is over UI elements
func (se *SkinEditor) isOverUI(mx, my int) bool {
	// Check if over tool panel
	if mx >= se.toolPanelX && mx <= se.toolPanelX+150 &&
		my >= se.toolPanelY && my <= se.toolPanelY+200 {
		return true
	}

	// Check if over palette
	if mx >= se.paletteX && mx <= se.paletteX+PaletteSize*20 &&
		my >= se.paletteY && my <= se.paletteY+20 {
		return true
	}

	// Check if over preview area
	if mx >= se.previewX && mx <= se.previewX+se.previewSize &&
		my >= se.previewY && my <= se.previewY+se.previewSize {
		return true
	}

	return false
}

// handleEditInput handles input in edit mode
func (se *SkinEditor) handleEditInput() {
	// Mouse input
	mx, my := ebiten.CursorPosition()

	// Handle tool selection first (only on click)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		se.handleToolSelection(mx, my)
		se.handlePaletteSelection(mx, my)
	}

	// Handle drawing (only when mouse is held down)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Only draw if not clicking on UI elements
		if !se.isOverUI(mx, my) {
			se.handleDrawing(mx, my)
		}
	} else {
		se.isDrawing = false
		se.lastPixelX = -1
		se.lastPixelY = -1
	}

	// Keyboard shortcuts
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		se.saveSkin()
	}

	// Tool shortcuts
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		se.selectedTool = ToolPencil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		se.selectedTool = ToolEraser
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		se.selectedTool = ToolFill
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		se.selectedTool = ToolEyedropper
	}

	// Brush size
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		se.brushSize = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		se.brushSize = 2
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		se.brushSize = 3
	}

	// Undo/Redo
	if inpututil.IsKeyJustPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		se.undo()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyY) {
		se.redo()
	}

	// Zoom
	if inpututil.IsKeyJustPressed(ebiten.KeyControl) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		se.zoomLevel = 1.0
	}

	// Layer management shortcuts
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		layerName := fmt.Sprintf("Layer %d", len(se.layers)+1)
		se.addLayer(layerName)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		se.deleteLayer()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		// Toggle layer visibility
		if se.activeLayer >= 0 && se.activeLayer < len(se.layerVisible) {
			se.layerVisible[se.activeLayer] = !se.layerVisible[se.activeLayer]
			se.layers[se.activeLayer].Visible = se.layerVisible[se.activeLayer]
			se.mergeLayers()
			se.previewSkin = se.copySkin(se.skinData)
		}
	}
	// Layer selection with number keys
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		se.setActiveLayer(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key5) && len(se.layers) > 1 {
		se.setActiveLayer(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key6) && len(se.layers) > 2 {
		se.setActiveLayer(2)
	}

	// Mouse wheel for zoom
	_, scrollY := ebiten.Wheel()
	if scrollY != 0 {
		se.zoomLevel *= math.Pow(1.1, -scrollY)
		se.zoomLevel = math.Max(0.5, math.Min(3.0, se.zoomLevel))
	}
}

// handleDrawing handles pixel drawing
func (se *SkinEditor) handleDrawing(mx, my int) {
	// Check if mouse is over canvas
	if mx < se.canvasX || mx > se.canvasX+se.canvasSize ||
		my < se.canvasY || my > se.canvasY+se.canvasSize {
		return
	}

	// Convert to skin coordinates
	skinX := int((float64(mx-se.canvasX) / float64(se.canvasSize)) * float64(SkinWidth))
	skinY := int((float64(my-se.canvasY) / float64(se.canvasSize)) * float64(SkinHeight))

	// Clamp to skin bounds
	skinX = max(0, min(SkinWidth-1, skinX))
	skinY = max(0, min(SkinHeight-1, skinY))

	// Log drawing attempt (only for debugging - remove in production)
	if se.lastPixelX != skinX || se.lastPixelY != skinY {
		// log.Printf("Drawing at skin coords: %d,%d from mouse: %d,%d", skinX, skinY, mx, my)
	}

	// Save to history only when starting to draw at a new position
	if !se.isDrawing || (skinX != se.lastPixelX || skinY != se.lastPixelY) {
		se.saveToHistory()
	}

	switch se.selectedTool {
	case ToolPencil:
		se.drawPixel(skinX, skinY)
	case ToolEraser:
		se.erasePixel(skinX, skinY)
	case ToolEyedropper:
		se.pickColor(skinX, skinY)
	case ToolFill:
		se.fillArea(skinX, skinY)
	case ToolAirbrush:
		se.drawAirbrush(skinX, skinY)
	case ToolSpray:
		se.drawSpray(skinX, skinY)
	case ToolGradient:
		se.drawGradient(skinX, skinY)
	case ToolSelect:
		se.handleSelection(skinX, skinY)
	case ToolMove:
		se.handleMove(skinX, skinY)
	case ToolSymmetry:
		se.drawSymmetry(skinX, skinY)
	}

	se.isDrawing = true
	se.lastPixelX = skinX
	se.lastPixelY = skinY
}

// drawPixel draws a pixel at the specified position
func (se *SkinEditor) drawPixel(x, y int) {
	// Check if we have an active layer
	if se.activeLayer < 0 || se.activeLayer >= len(se.layers) {
		return
	}

	layer := se.layers[se.activeLayer]
	if layer.Locked {
		return
	}

	// Draw with brush size on active layer
	for dy := -se.brushSize / 2; dy <= se.brushSize/2; dy++ {
		for dx := -se.brushSize / 2; dx <= se.brushSize/2; dx++ {
			px, py := x+dx, y+dy
			if px >= 0 && px < SkinWidth && py >= 0 && py < SkinHeight {
				// Check if within circular brush
				if dx*dx+dy*dy <= (se.brushSize/2)*(se.brushSize/2) {
					layer.Pixels[py][px] = se.selectedColor
				}
			}
		}
	}

	// Merge layers and update preview
	se.mergeLayers()
	se.previewSkin = se.copySkin(se.skinData)
}

// erasePixel erases a pixel at the specified position
func (se *SkinEditor) erasePixel(x, y int) {
	// Erase with brush size
	for dy := -se.brushSize / 2; dy <= se.brushSize/2; dy++ {
		for dx := -se.brushSize / 2; dx <= se.brushSize/2; dx++ {
			px, py := x+dx, y+dy
			if px >= 0 && px < SkinWidth && py >= 0 && py < SkinHeight {
				// Check if within circular brush
				if dx*dx+dy*dy <= (se.brushSize/2)*(se.brushSize/2) {
					se.skinData.Pixels[py][px] = color.RGBA{0, 0, 0, 0} // Transparent
				}
			}
		}
	}

	se.previewSkin = se.copySkin(se.skinData)
}

// pickColor picks color from the specified position
func (se *SkinEditor) pickColor(x, y int) {
	if x >= 0 && x < SkinWidth && y >= 0 && y < SkinHeight {
		se.selectedColor = se.skinData.Pixels[y][x]
		se.selectedTool = ToolPencil // Switch back to pencil
	}
}

// fillArea fills an area with the selected color
func (se *SkinEditor) fillArea(startX, startY int) {
	if startX < 0 || startX >= SkinWidth || startY < 0 || startY >= SkinHeight {
		return
	}

	se.saveToHistory()

	targetColor := se.skinData.Pixels[startY][startX]
	if targetColor.R == se.selectedColor.R &&
		targetColor.G == se.selectedColor.G &&
		targetColor.B == se.selectedColor.B &&
		targetColor.A == se.selectedColor.A {
		return // Already the same color
	}

	// Flood fill algorithm with safety limits
	stack := [][2]int{{startX, startY}}
	visited := make(map[[2]int]bool)
	maxIterations := SkinWidth * SkinHeight // Prevent infinite loops
	iterations := 0

	for len(stack) > 0 && iterations < maxIterations {
		iterations++

		last := stack[len(stack)-1]
		x, y := last[0], last[1]
		stack = stack[:len(stack)-1]

		key := [2]int{x, y}
		if visited[key] {
			continue
		}
		visited[key] = true

		if x < 0 || x >= SkinWidth || y < 0 || y >= SkinHeight {
			continue
		}

		currentColor := se.skinData.Pixels[y][x]
		if currentColor.R != targetColor.R ||
			currentColor.G != targetColor.G ||
			currentColor.B != targetColor.B ||
			currentColor.A != targetColor.A {
			continue
		}

		se.skinData.Pixels[y][x] = se.selectedColor

		// Add neighbors
		stack = append(stack, [2]int{x + 1, y})
		stack = append(stack, [2]int{x - 1, y})
		stack = append(stack, [2]int{x, y + 1})
		stack = append(stack, [2]int{x, y - 1})
	}

	if iterations >= maxIterations {
		log.Printf("Fill algorithm stopped after %d iterations to prevent infinite loop", maxIterations)
	}

	se.previewSkin = se.copySkin(se.skinData)
}

// handleToolSelection handles tool selection from UI
func (se *SkinEditor) handleToolSelection(mx, my int) {
	// Check if clicking on tool panel
	if mx >= se.toolPanelX && mx <= se.toolPanelX+150 &&
		my >= se.toolPanelY && my <= se.toolPanelY+200 {

		toolIndex := (my - se.toolPanelY) / 30
		tools := []Tool{ToolPencil, ToolEraser, ToolFill, ToolEyedropper, ToolLine, ToolRectangle, ToolCircle}

		if toolIndex >= 0 && toolIndex < len(tools) {
			se.selectedTool = tools[toolIndex]
		}
	}
}

// handlePaletteSelection handles color palette selection
func (se *SkinEditor) handlePaletteSelection(mx, my int) {
	// Check if clicking on palette
	if mx >= se.paletteX && mx <= se.paletteX+PaletteSize*20 &&
		my >= se.paletteY && my <= se.paletteY+PaletteSize*20 {

		paletteX := (mx - se.paletteX) / 20
		paletteY := (my - se.paletteY) / 20

		if paletteX >= 0 && paletteX < PaletteSize && paletteY >= 0 && paletteY < 1 {
			index := paletteX
			if index < len(se.palette) {
				se.selectedColor = se.palette[index]
				se.selectedPalette = index
			}
		}
	}
}

// handlePreviewInput handles input in preview mode
func (se *SkinEditor) handlePreviewInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		se.editorMode = ModeEdit
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		se.previewAnimating = !se.previewAnimating
	}
}

// handleColorPickerInput handles input in color picker mode
func (se *SkinEditor) handleColorPickerInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		se.editorMode = ModeEdit
	}
}

// Draw renders the skin editor
func (se *SkinEditor) Draw(screen *ebiten.Image) {
	// Safety check
	if se == nil || se.skinData == nil {
		ebitenutil.DebugPrintAt(screen, "SKIN EDITOR ERROR - Not initialized", 10, 10)
		return
	}

	// Draw background
	screen.Fill(se.backgroundColor)

	// Draw based on current mode
	switch se.editorMode {
	case ModeEdit:
		se.drawEditMode(screen)
	case ModePreview:
		se.drawPreviewMode(screen)
	case ModeColorPicker:
		se.drawColorPickerMode(screen)
	}
}

// drawEditMode renders the editing interface
func (se *SkinEditor) drawEditMode(screen *ebiten.Image) {
	// Draw professional toolbar
	se.drawToolbar(screen)

	// Draw sidebar
	se.drawSidebar(screen)

	// Draw canvas
	se.drawCanvas(screen)

	// Draw preview
	se.drawPreview(screen)

	// Draw status bar
	se.drawStatusBar(screen)
}

// drawCanvas renders the main editing canvas
func (se *SkinEditor) drawCanvas(screen *ebiten.Image) {
	// Draw canvas background
	ebitenutil.DrawRect(screen, float64(se.canvasX), float64(se.canvasY),
		float64(se.canvasSize), float64(se.canvasSize), color.RGBA{20, 20, 30, 255})

	// Draw pixels
	for y := 0; y < SkinHeight; y++ {
		for x := 0; x < SkinWidth; x++ {
			pixelColor := se.skinData.Pixels[y][x]

			// Skip transparent pixels
			if pixelColor.A == 0 {
				continue
			}

			// Calculate pixel position with zoom
			pixelSize := int(float64(se.pixelSize) * se.zoomLevel)
			px := se.canvasX + x*pixelSize
			py := se.canvasY + y*pixelSize

			// Draw pixel
			ebitenutil.DrawRect(screen, float64(px), float64(py),
				float64(pixelSize), float64(pixelSize), pixelColor)
		}
	}

	// Draw grid
	if se.zoomLevel > 1.0 {
		gridSize := int(float64(se.pixelSize) * se.zoomLevel)
		for i := 0; i <= SkinWidth; i++ {
			x := se.canvasX + i*gridSize
			ebitenutil.DrawLine(screen, float64(x), float64(se.canvasY),
				float64(x), float64(se.canvasY+se.canvasSize), se.gridColor)
		}
		for i := 0; i <= SkinHeight; i++ {
			y := se.canvasY + i*gridSize
			ebitenutil.DrawLine(screen, float64(se.canvasX), float64(y),
				float64(se.canvasX+se.canvasSize), float64(y), se.gridColor)
		}
	}

	// Draw canvas border
	ebitenutil.DrawRect(screen, float64(se.canvasX), float64(se.canvasY),
		float64(se.canvasSize), float64(se.canvasSize), color.RGBA{100, 100, 120, 255})
}

// drawPreview renders the 3D preview
func (se *SkinEditor) drawPreview(screen *ebiten.Image) {
	// Draw preview background
	ebitenutil.DrawRect(screen, float64(se.previewX), float64(se.previewY),
		float64(se.previewSize), float64(se.previewSize), color.RGBA{40, 40, 50, 255})

	// Draw preview title
	ebitenutil.DebugPrintAt(screen, "PREVIEW", se.previewX, se.previewY-20)

	// Draw simple 2D preview (could be enhanced to 3D)
	previewScale := float64(se.previewSize) / float64(SkinWidth)
	for y := 0; y < SkinHeight; y++ {
		for x := 0; x < SkinWidth; x++ {
			pixelColor := se.previewSkin.Pixels[y][x]
			if pixelColor.A == 0 {
				continue
			}

			px := se.previewX + int(float64(x)*previewScale)
			py := se.previewY + int(float64(y)*previewScale)

			ebitenutil.DrawRect(screen, float64(px), float64(py),
				previewScale, previewScale, pixelColor)
		}
	}

	// Draw preview border
	ebitenutil.DrawRect(screen, float64(se.previewX), float64(se.previewY),
		float64(se.previewSize), float64(se.previewSize), color.RGBA{80, 80, 100, 255})
}

// drawPalette renders the color palette
func (se *SkinEditor) drawPalette(screen *ebiten.Image) {
	// Draw palette background
	ebitenutil.DrawRect(screen, float64(se.paletteX), float64(se.paletteY),
		float64(PaletteSize*20), float64(20), color.RGBA{30, 30, 40, 255})

	// Draw palette title
	ebitenutil.DebugPrintAt(screen, "COLOR PALETTE", se.paletteX, se.paletteY-20)

	// Draw color swatches
	for i, paletteColor := range se.palette {
		x := se.paletteX + i*20
		y := se.paletteY

		// Draw color
		ebitenutil.DrawRect(screen, float64(x), float64(y), 20, 20, paletteColor)

		// Highlight selected color
		if i == se.selectedPalette {
			ebitenutil.DrawRect(screen, float64(x), float64(y), 20, 20, se.selectionColor)
		}
	}

	// Draw current color indicator
	currentColorX := se.paletteX + PaletteSize*20 + 20
	ebitenutil.DrawRect(screen, float64(currentColorX), float64(se.paletteY), 40, 40, se.selectedColor)
	ebitenutil.DebugPrintAt(screen, "Current", currentColorX, se.paletteY-20)
}

// drawToolPanel renders the tool selection panel
func (se *SkinEditor) drawToolPanel(screen *ebiten.Image) {
	// Draw tool panel background
	ebitenutil.DrawRect(screen, float64(se.toolPanelX), float64(se.toolPanelY), 150, 200, color.RGBA{30, 30, 40, 255})

	// Draw panel title
	ebitenutil.DebugPrintAt(screen, "TOOLS", se.toolPanelX, se.toolPanelY-20)

	// Draw tools
	tools := []struct {
		name string
		tool Tool
		key  string
	}{
		{"Pencil", ToolPencil, "B"},
		{"Eraser", ToolEraser, "E"},
		{"Fill", ToolFill, "F"},
		{"Eyedropper", ToolEyedropper, "I"},
		{"Line", ToolLine, "L"},
		{"Rectangle", ToolRectangle, "R"},
		{"Circle", ToolCircle, "C"},
		{"Airbrush", ToolAirbrush, "A"},
		{"Spray", ToolSpray, "S"},
		{"Gradient", ToolGradient, "G"},
		{"Select", ToolSelect, "X"},
		{"Move", ToolMove, "M"},
		{"Symmetry", ToolSymmetry, "Y"},
	}

	for i, toolInfo := range tools {
		y := se.toolPanelY + i*30

		// Highlight selected tool
		if se.selectedTool == toolInfo.tool {
			ebitenutil.DrawRect(screen, float64(se.toolPanelX), float64(y), 150, 25, se.selectionColor)
		}

		// Draw tool name
		toolText := fmt.Sprintf("%s [%s]", toolInfo.name, toolInfo.key)
		ebitenutil.DebugPrintAt(screen, toolText, se.toolPanelX+5, y+5)
	}

	// Draw brush size indicator
	brushY := se.toolPanelY + len(tools)*30 + 20
	brushText := fmt.Sprintf("Brush Size: %d (1-3)", se.brushSize)
	ebitenutil.DebugPrintAt(screen, brushText, se.toolPanelX+5, brushY)
}

// drawUI renders UI elements
func (se *SkinEditor) drawUI(screen *ebiten.Image) {
	// Draw title
	title := fmt.Sprintf("SKIN EDITOR - %s", se.skinData.Name)
	ebitenutil.DebugPrintAt(screen, title, 10, 10)

	// Draw instructions
	instructions := []string{
		"B: Pencil  E: Eraser  F: Fill  I: Eyedropper",
		"1-3: Brush Size  Ctrl+Z: Undo  Ctrl+Y: Redo",
		"Mouse Wheel: Zoom  ESC: Save & Exit",
	}

	for i, instruction := range instructions {
		ebitenutil.DebugPrintAt(screen, instruction, 10, 680-i*20)
	}

	// Draw cursor position
	mx, my := ebiten.CursorPosition()
	if mx >= se.canvasX && mx <= se.canvasX+se.canvasSize &&
		my >= se.canvasY && my <= se.canvasY+se.canvasSize {

		skinX := int((float64(mx-se.canvasX) / float64(se.canvasSize)) * float64(SkinWidth))
		skinY := int((float64(my-se.canvasY) / float64(se.canvasSize)) * float64(SkinHeight))

		posText := fmt.Sprintf("Position: %d, %d | Mouse: %d, %d | Drawing: %v", skinX, skinY, mx, my, se.isDrawing)
		ebitenutil.DebugPrintAt(screen, posText, 10, 30)

		// Draw tool info
		toolName := ""
		switch se.selectedTool {
		case ToolPencil:
			toolName = "Pencil"
		case ToolEraser:
			toolName = "Eraser"
		case ToolFill:
			toolName = "Fill"
		case ToolEyedropper:
			toolName = "Eyedropper"
		case ToolLine:
			toolName = "Line"
		case ToolRectangle:
			toolName = "Rectangle"
		case ToolCircle:
			toolName = "Circle"
		case ToolAirbrush:
			toolName = "Airbrush"
		case ToolSpray:
			toolName = "Spray"
		case ToolGradient:
			toolName = "Gradient"
		case ToolSelect:
			toolName = "Select"
		case ToolMove:
			toolName = "Move"
		case ToolSymmetry:
			toolName = "Symmetry"
		}

		toolInfo := fmt.Sprintf("Tool: %s | Brush: %d | Opacity: %.1f | RGB(%d,%d,%d)",
			toolName, se.brushSize, se.brushOpacity,
			se.selectedColor.R, se.selectedColor.G, se.selectedColor.B)
		ebitenutil.DebugPrintAt(screen, toolInfo, 10, 50)
	}
}

// drawPreviewMode renders the preview mode
func (se *SkinEditor) drawPreviewMode(screen *ebiten.Image) {
	// Full screen preview
	previewScale := math.Min(float64(ScreenWidth)/float64(SkinWidth),
		float64(ScreenHeight)/float64(SkinHeight)) * 0.8

	previewX := (ScreenWidth - int(float64(SkinWidth)*previewScale)) / 2
	previewY := (ScreenHeight - int(float64(SkinHeight)*previewScale)) / 2

	for y := 0; y < SkinHeight; y++ {
		for x := 0; x < SkinWidth; x++ {
			pixelColor := se.previewSkin.Pixels[y][x]
			if pixelColor.A == 0 {
				continue
			}

			px := previewX + int(float64(x)*previewScale)
			py := previewY + int(float64(y)*previewScale)

			ebitenutil.DrawRect(screen, float64(px), float64(py),
				previewScale, previewScale, pixelColor)
		}
	}

	// Instructions
	ebitenutil.DebugPrintAt(screen, "PREVIEW MODE - SPACE: Toggle Animation  ESC: Back", 10, 10)
}

// drawColorPickerMode renders the color picker mode
func (se *SkinEditor) drawColorPickerMode(screen *ebiten.Image) {
	// Advanced color picker interface
	ebitenutil.DebugPrintAt(screen, "COLOR PICKER MODE - ESC: Back", 10, 10)

	// Draw current color large
	ebitenutil.DrawRect(screen, 100, 100, 200, 200, se.selectedColor)

	// Draw RGB values
	rText := fmt.Sprintf("R: %d", se.selectedColor.R)
	gText := fmt.Sprintf("G: %d", se.selectedColor.G)
	bText := fmt.Sprintf("B: %d", se.selectedColor.B)
	aText := fmt.Sprintf("A: %d", se.selectedColor.A)

	ebitenutil.DebugPrintAt(screen, rText, 320, 100)
	ebitenutil.DebugPrintAt(screen, gText, 320, 120)
	ebitenutil.DebugPrintAt(screen, bText, 320, 140)
	ebitenutil.DebugPrintAt(screen, aText, 320, 160)
}

// saveToHistory saves current state to history
func (se *SkinEditor) saveToHistory() {
	// Remove any states after current index
	if se.historyIndex < len(se.history)-1 {
		se.history = se.history[:se.historyIndex+1]
	}

	// Add current state
	se.history = append(se.history, se.copySkin(se.skinData))
	se.historyIndex++

	// Limit history size and prevent memory leaks
	if len(se.history) > se.maxHistory {
		// Remove oldest entries
		se.history = se.history[1:]
		se.historyIndex--
	}

	// Prevent excessive memory usage
	if len(se.history) > se.maxHistory*2 { // Emergency cleanup
		se.history = se.history[len(se.history)-se.maxHistory:]
		se.historyIndex = len(se.history) - 1
		log.Printf("History emergency cleanup - reduced to %d entries", len(se.history))
	}
}

// undo restores previous state
func (se *SkinEditor) undo() {
	if se.historyIndex > 0 {
		se.historyIndex--
		se.skinData = se.copySkin(se.history[se.historyIndex])
		se.previewSkin = se.copySkin(se.skinData)
	}
}

// redo restores next state
func (se *SkinEditor) redo() {
	if se.historyIndex < len(se.history)-1 {
		se.historyIndex++
		se.skinData = se.copySkin(se.history[se.historyIndex])
		se.previewSkin = se.copySkin(se.skinData)
	}
}

// copySkin creates a deep copy of skin data
func (se *SkinEditor) copySkin(skin *SkinData) *SkinData {
	newSkin := &SkinData{
		Name:      skin.Name,
		CreatedAt: skin.CreatedAt,
		UpdatedAt: time.Now(),
		Width:     skin.Width,
		Height:    skin.Height,
		Pixels:    make([][]color.RGBA, skin.Height),
		Metadata:  make(map[string]string),
	}

	// Copy metadata
	for k, v := range skin.Metadata {
		newSkin.Metadata[k] = v
	}

	// Copy pixels
	for y := 0; y < skin.Height; y++ {
		newSkin.Pixels[y] = make([]color.RGBA, skin.Width)
		for x := 0; x < skin.Width; x++ {
			newSkin.Pixels[y][x] = skin.Pixels[y][x]
		}
	}

	return newSkin
}

// SaveSkin saves the current skin to file (public method)
func (se *SkinEditor) SaveSkin() error {
	return se.saveSkin()
}

// saveSkin saves the current skin to file
func (se *SkinEditor) saveSkin() error {
	// Create skin directory if it doesn't exist
	if err := os.MkdirAll(se.skinDirectory, 0755); err != nil {
		log.Printf("Failed to create skins directory: %v", err)
		return err
	}

	// Update skin data
	se.skinData.UpdatedAt = time.Now()

	// Save skin file
	skinFile := filepath.Join(se.skinDirectory, se.skinData.Name+".json")
	data, err := json.MarshalIndent(se.skinData, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal skin data: %v", err)
		return err
	}

	if err := os.WriteFile(skinFile, data, 0644); err != nil {
		log.Printf("Failed to save skin file: %v", err)
		return err
	}

	// Update skin config
	if err := se.updateSkinConfig(); err != nil {
		log.Printf("Failed to update skin config: %v", err)
		return err
	}

	log.Printf("Skin saved: %s", se.skinData.Name)
	return nil
}

// loadSkinConfig loads the skin configuration
func (se *SkinEditor) loadSkinConfig() error {
	// Ensure skins directory exists
	if err := os.MkdirAll(se.skinDirectory, 0755); err != nil {
		log.Printf("Failed to create skins directory: %v", err)
		return err
	}

	configFile := filepath.Join(se.skinDirectory, "config.json")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create default config
		config := &SkinConfig{
			CurrentSkin: "Default",
			Skins:       []string{"Default"},
			LastUsed:    time.Now(),
		}

		data, _ := json.MarshalIndent(config, "", "  ")
		if err := os.WriteFile(configFile, data, 0644); err != nil {
			log.Printf("Failed to create skin config: %v", err)
			return err
		}
		return nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("Failed to read skin config: %v", err)
		return err
	}

	var config SkinConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Printf("Failed to parse skin config: %v", err)
		return err
	}

	// Load current skin if it exists
	if config.CurrentSkin != "" {
		se.loadSkin(config.CurrentSkin)
	}

	return nil
}

// updateSkinConfig updates the skin configuration file
func (se *SkinEditor) updateSkinConfig() error {
	configFile := filepath.Join(se.skinDirectory, "config.json")

	// Read existing config
	var config SkinConfig
	if data, err := os.ReadFile(configFile); err == nil {
		json.Unmarshal(data, &config)
	}

	// Update config
	config.CurrentSkin = se.skinData.Name
	config.LastUsed = time.Now()

	// Add to skins list if not present
	found := false
	for _, skinName := range config.Skins {
		if skinName == se.skinData.Name {
			found = true
			break
		}
	}
	if !found {
		config.Skins = append(config.Skins, se.skinData.Name)
	}

	// Save config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}

// loadSkin loads a specific skin from file
func (se *SkinEditor) loadSkin(skinName string) error {
	skinFile := filepath.Join(se.skinDirectory, skinName+".json")

	data, err := os.ReadFile(skinFile)
	if err != nil {
		log.Printf("Failed to load skin %s: %v", skinName, err)
		// If skin doesn't exist, create default skin
		if skinName == "Default" {
			se.createDefaultSkin()
			return nil
		}
		return err
	}

	var skin SkinData
	if err := json.Unmarshal(data, &skin); err != nil {
		log.Printf("Failed to parse skin %s: %v", skinName, err)
		return err
	}

	se.skinData = &skin
	se.previewSkin = se.copySkin(&skin)
	se.addToHistory(&skin)

	log.Printf("Skin loaded: %s", skinName)
	return nil
}

// addToHistory adds a skin state to history
func (se *SkinEditor) addToHistory(skin *SkinData) {
	se.history = append(se.history, se.copySkin(skin))
	se.historyIndex = len(se.history) - 1
}

// GetSkinData returns the current skin data
func (se *SkinEditor) GetSkinData() *SkinData {
	return se.skinData
}

// SetSkinData sets the current skin data
func (se *SkinEditor) SetSkinData(skin *SkinData) {
	se.skinData = se.copySkin(skin)
	se.previewSkin = se.copySkin(skin)
	se.addToHistory(skin)
}

// drawToolbar renders the professional toolbar
func (se *SkinEditor) drawToolbar(screen *ebiten.Image) {
	// Draw toolbar background
	toolbarColor := color.RGBA{45, 45, 60, 255}
	ebitenutil.DrawRect(screen, 0, 0, 1280, float64(se.toolbarHeight), toolbarColor)

	// Draw toolbar separator
	ebitenutil.DrawRect(screen, 0, float64(se.toolbarHeight-1), 1280, 1, color.RGBA{80, 80, 100, 255})

	// Draw file operations section
	fileX := 10
	ebitenutil.DebugPrintAt(screen, "File", fileX, 20)
	ebitenutil.DebugPrintAt(screen, "[N] New [S] Save [O] Open", fileX, 35)

	// Draw edit operations section
	editX := 150
	ebitenutil.DebugPrintAt(screen, "Edit", editX, 20)
	ebitenutil.DebugPrintAt(screen, "[Ctrl+Z] Undo [Ctrl+Y] Redo", editX, 35)

	// Draw view operations section
	viewX := 350
	ebitenutil.DebugPrintAt(screen, "View", viewX, 20)
	zoomText := fmt.Sprintf("Zoom: %.0f%%", se.zoomLevel*100)
	ebitenutil.DebugPrintAt(screen, zoomText, viewX, 35)

	// Draw current tool info
	toolX := 500
	toolName := ""
	switch se.selectedTool {
	case ToolPencil:
		toolName = "Pencil"
	case ToolEraser:
		toolName = "Eraser"
	case ToolFill:
		toolName = "Fill"
	case ToolEyedropper:
		toolName = "Eyedropper"
	case ToolAirbrush:
		toolName = "Airbrush"
	case ToolSpray:
		toolName = "Spray"
	case ToolGradient:
		toolName = "Gradient"
	case ToolSymmetry:
		toolName = "Symmetry"
	default:
		toolName = "Unknown"
	}
	ebitenutil.DebugPrintAt(screen, "Tool: "+toolName, toolX, 20)
	brushText := fmt.Sprintf("Size: %d  Opacity: %.0f%%", se.brushSize, se.brushOpacity*100)
	ebitenutil.DebugPrintAt(screen, brushText, toolX, 35)

	// Draw skin info
	infoX := 750
	skinName := se.skinData.Name
	if len(skinName) > 20 {
		skinName = skinName[:17] + "..."
	}
	ebitenutil.DebugPrintAt(screen, "Skin: "+skinName, infoX, 20)
	layerText := fmt.Sprintf("Layer %d/%d", se.activeLayer+1, len(se.layers))
	ebitenutil.DebugPrintAt(screen, layerText, infoX, 35)
}

// drawSidebar renders the professional sidebar with tools and layers
func (se *SkinEditor) drawSidebar(screen *ebiten.Image) {
	sidebarX := 1280 - se.sidebarWidth

	// Draw sidebar background
	sidebarColor := color.RGBA{40, 40, 55, 255}
	ebitenutil.DrawRect(screen, float64(sidebarX), float64(se.toolbarHeight),
		float64(se.sidebarWidth), 720-float64(se.toolbarHeight+se.statusBarHeight), sidebarColor)

	// Draw sidebar separator
	ebitenutil.DrawRect(screen, float64(sidebarX-1), float64(se.toolbarHeight), 1, 720-float64(se.toolbarHeight+se.statusBarHeight), color.RGBA{80, 80, 100, 255})

	// Draw tools section
	toolsY := se.toolbarHeight + 20
	ebitenutil.DebugPrintAt(screen, "TOOLS", sidebarX+10, toolsY)

	// Compact tool grid
	tools := []struct {
		name string
		tool Tool
		key  string
	}{
		{"Pencil", ToolPencil, "B"}, {"Eraser", ToolEraser, "E"},
		{"Fill", ToolFill, "F"}, {"Eye", ToolEyedropper, "I"},
		{"Air", ToolAirbrush, "A"}, {"Spray", ToolSpray, "S"},
		{"Grad", ToolGradient, "G"}, {"Sym", ToolSymmetry, "Y"},
	}

	for i, toolInfo := range tools {
		toolX := sidebarX + 10 + (i%4)*45
		toolY := toolsY + 25 + (i/4)*25

		// Highlight selected tool
		if se.selectedTool == toolInfo.tool {
			ebitenutil.DrawRect(screen, float64(toolX-2), float64(toolY-2), 40, 20, se.selectionColor)
		}

		// Draw tool name and key
		toolText := fmt.Sprintf("%s[%s]", toolInfo.name, toolInfo.key)
		ebitenutil.DebugPrintAt(screen, toolText, toolX, toolY)
	}

	// Draw color section
	colorY := toolsY + 80
	ebitenutil.DebugPrintAt(screen, "COLOR", sidebarX+10, colorY)

	// Draw current color
	ebitenutil.DrawRect(screen, float64(sidebarX+10), float64(colorY+20), 40, 40, se.selectedColor)

	// Draw RGB values
	rText := fmt.Sprintf("R:%d G:%d B:%d", se.selectedColor.R, se.selectedColor.G, se.selectedColor.B)
	ebitenutil.DebugPrintAt(screen, rText, sidebarX+60, colorY+25)

	// Draw mini palette (8 colors)
	paletteY := colorY + 70
	for i := 0; i < 8 && i < len(se.palette); i++ {
		paletteX := sidebarX + 10 + (i%4)*25
		paletteYPos := paletteY + (i/4)*25

		ebitenutil.DrawRect(screen, float64(paletteX), float64(paletteYPos), 20, 20, se.palette[i])

		if i == se.selectedPalette {
			ebitenutil.DrawRect(screen, float64(paletteX), float64(paletteYPos), 20, 20, se.selectionColor)
		}
	}

	// Draw layers section
	layersY := paletteY + 60
	ebitenutil.DebugPrintAt(screen, "LAYERS", sidebarX+10, layersY)

	// Draw layer list (compact)
	for i, layer := range se.layers {
		if i >= 4 { // Show max 4 layers
			break
		}
		layerYPos := layersY + 25 + i*20

		// Highlight active layer
		if i == se.activeLayer {
			ebitenutil.DrawRect(screen, float64(sidebarX+5), float64(layerYPos-2), float64(se.sidebarWidth-10), 18, se.selectionColor)
		}

		// Layer visibility
		visColor := color.RGBA{60, 60, 60, 255}
		if se.layerVisible[i] && layer.Visible {
			visColor = color.RGBA{0, 200, 0, 255}
		}
		ebitenutil.DrawRect(screen, float64(sidebarX+8), float64(layerYPos+2), 8, 8, visColor)

		// Layer name (truncated)
		layerName := layer.Name
		if len(layerName) > 12 {
			layerName = layerName[:9] + "..."
		}
		if layer.Locked {
			layerName += "🔒"
		}
		ebitenutil.DebugPrintAt(screen, layerName, sidebarX+20, layerYPos)
	}

	// Draw layer controls
	controlsY := layersY + 110
	ebitenutil.DebugPrintAt(screen, "[+]Add [-]Del [V]Vis", sidebarX+10, controlsY)
}

// drawStatusBar renders the professional status bar
func (se *SkinEditor) drawStatusBar(screen *ebiten.Image) {
	statusY := 720 - se.statusBarHeight

	// Draw status bar background
	statusColor := color.RGBA{35, 35, 50, 255}
	ebitenutil.DrawRect(screen, 0, float64(statusY), 1280, float64(se.statusBarHeight), statusColor)

	// Draw status bar separator
	ebitenutil.DrawRect(screen, 0, float64(statusY-1), 1280, 1, color.RGBA{80, 80, 100, 255})

	// Draw cursor position
	mx, my := ebiten.CursorPosition()
	if mx >= se.canvasX && mx <= se.canvasX+se.canvasSize &&
		my >= se.canvasY && my <= se.canvasY+se.canvasSize {

		skinX := int((float64(mx-se.canvasX) / float64(se.canvasSize)) * float64(SkinWidth))
		skinY := int((float64(my-se.canvasY) / float64(se.canvasSize)) * float64(SkinHeight))

		posText := fmt.Sprintf("Position: %d, %d", skinX, skinY)
		ebitenutil.DebugPrintAt(screen, posText, 10, statusY+8)
	}

	// Draw drawing status
	statusText := "Ready"
	if se.isDrawing {
		statusText = "Drawing..."
	}
	ebitenutil.DebugPrintAt(screen, statusText, 150, statusY+8)

	// Draw help text
	ebitenutil.DebugPrintAt(screen, "ESC: Save & Exit | Mouse Wheel: Zoom | 1-3: Brush Size", 400, statusY+8)
}

// Missing drawing tool implementations

// drawAirbrush draws with soft edges
func (se *SkinEditor) drawAirbrush(x, y int) {
	// Check if we have an active layer
	if se.activeLayer < 0 || se.activeLayer >= len(se.layers) {
		return
	}

	layer := se.layers[se.activeLayer]
	if layer.Locked {
		return
	}

	// Apply opacity to brush
	opacity := se.brushOpacity
	for dy := -se.brushSize; dy <= se.brushSize; dy++ {
		for dx := -se.brushSize; dx <= se.brushSize; dx++ {
			distance := math.Sqrt(float64(dx*dx + dy*dy))
			if distance <= float64(se.brushSize) {
				px, py := x+dx, y+dy
				if px >= 0 && px < SkinWidth && py >= 0 && py < SkinHeight {
					// Apply opacity blending
					existing := layer.Pixels[py][px]
					blended := se.blendColors(existing, se.selectedColor, opacity)
					layer.Pixels[py][px] = blended
				}
			}
		}
	}
}

// drawSpray creates scattered pixels
func (se *SkinEditor) drawSpray(x, y int) {
	// Check if we have an active layer
	if se.activeLayer < 0 || se.activeLayer >= len(se.layers) {
		return
	}

	layer := se.layers[se.activeLayer]
	if layer.Locked {
		return
	}

	// Deterministic spray pattern (no random)
	sprayCount := se.brushFlow / 10
	for i := 0; i < sprayCount; i++ {
		// Create deterministic spray pattern
		angle := float64(i) / float64(sprayCount) * 2 * math.Pi
		sprayX := x + int(math.Cos(angle)*float64(se.brushSize))
		sprayY := y + int(math.Sin(angle)*float64(se.brushSize))

		if sprayX >= 0 && sprayX < SkinWidth && sprayY >= 0 && sprayY < SkinHeight {
			// Fixed opacity for spray effect
			opacity := se.brushOpacity * 0.6
			existing := layer.Pixels[sprayY][sprayX]
			blended := se.blendColors(existing, se.selectedColor, opacity)
			layer.Pixels[sprayY][sprayX] = blended
		}
	}
}

// drawGradient creates color gradient
func (se *SkinEditor) drawGradient(x, y int) {
	// Check if we have an active layer
	if se.activeLayer < 0 || se.activeLayer >= len(se.layers) {
		return
	}

	layer := se.layers[se.activeLayer]
	if layer.Locked {
		return
	}

	// Simple linear gradient from center
	for dy := -se.brushSize; dy <= se.brushSize; dy++ {
		for dx := -se.brushSize; dx <= se.brushSize; dx++ {
			distance := math.Sqrt(float64(dx*dx + dy*dy))
			if distance <= float64(se.brushSize) {
				px, py := x+dx, y+dy
				if px >= 0 && px < SkinWidth && py >= 0 && py < SkinHeight {
					// Gradient based on distance
					ratio := distance / float64(se.brushSize)
					gradient := se.interpolateColor(se.selectedColor, color.RGBA{0, 0, 0, 0}, ratio)
					layer.Pixels[py][px] = gradient
				}
			}
		}
	}
}

// handleSelection manages selection area
func (se *SkinEditor) handleSelection(x, y int) {
	// Selection tool not implemented
}

// handleMove moves selected area
func (se *SkinEditor) handleMove(x, y int) {
	// Move tool not implemented
}

// drawSymmetry draws with mirror symmetry
func (se *SkinEditor) drawSymmetry(x, y int) {
	// Check if we have an active layer
	if se.activeLayer < 0 || se.activeLayer >= len(se.layers) {
		return
	}

	layer := se.layers[se.activeLayer]
	if layer.Locked {
		return
	}

	// Draw on both sides of center
	centerX := SkinWidth / 2
	mirrorX := centerX - (x - centerX)

	// Draw original pixel
	if x >= 0 && x < SkinWidth && y >= 0 && y < SkinHeight {
		layer.Pixels[y][x] = se.selectedColor
	}

	// Draw mirrored pixel
	if mirrorX >= 0 && mirrorX < SkinWidth && y >= 0 && y < SkinHeight {
		layer.Pixels[y][mirrorX] = se.selectedColor
	}
}

// interpolateColor creates smooth color transition
func (se *SkinEditor) interpolateColor(start, end color.RGBA, ratio float64) color.RGBA {
	r := float64(start.R)*(1.0-ratio) + float64(end.R)*ratio
	g := float64(start.G)*(1.0-ratio) + float64(end.G)*ratio
	b := float64(start.B)*(1.0-ratio) + float64(end.B)*ratio
	a := float64(start.A)*(1.0-ratio) + float64(end.A)*ratio

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}

// Utility functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Constants for screen dimensions (should match game constants)
const (
	ScreenWidth  = 1280
	ScreenHeight = 720
)
