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
	skinData      *SkinData
	originalSkin  *SkinData
	previewSkin  *SkinData
	
	// UI state
	editorMode    EditorMode
	selectedTool  Tool
	selectedColor color.RGBA
	brushSize     int
	
	// Canvas state
	canvasX, canvasY int
	canvasSize      int
	pixelSize       int
	zoomLevel       float64
	
	// Preview state
	previewX, previewY int
	previewSize       int
	previewAnimating  bool
	previewRotation   float64
	
	// Color palette
	palette       []color.RGBA
	paletteX      int
	paletteY      int
	selectedPalette int
	
	// Tool panel
	toolPanelX    int
	toolPanelY    int
	
	// History for undo/redo
	history       []*SkinData
	historyIndex  int
	maxHistory    int
	
	// Drawing state
	isDrawing     bool
	lastPixelX    int
	lastPixelY    int
	
	// Animation
	animTimer     float64
	cursorBlink   bool
	
	// File operations
	skinDirectory string
	currentSkin   string
	
	// Visual properties
	backgroundColor color.RGBA
	gridColor      color.RGBA
	selectionColor color.RGBA
	
	// For solid color drawing
	whiteImage     *ebiten.Image
}

// EditorMode represents different editing modes
type EditorMode int

const (
	ModeEdit EditorMode = iota
	ModePreview
	ModeColorPicker
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

// SkinConfig represents saved skin configuration
type SkinConfig struct {
	CurrentSkin string    `json:"currentSkin"`
	Skins       []string  `json:"skins"`
	LastUsed    time.Time `json:"lastUsed"`
}

const (
	SkinWidth  = 64
	SkinHeight = 64
	CanvasSize = 512
	PreviewSize = 128
	PaletteSize = 16
	MaxHistory = 50
)

// NewSkinEditor creates a new skin editor
func NewSkinEditor() *SkinEditor {
	log.Printf("Creating new skin editor...")
	
	// Create a 1x1 white image for solid color drawing
	whiteImage := ebiten.NewImage(1, 1)
	whiteImage.Fill(color.RGBA{255, 255, 255, 255})

	editor := &SkinEditor{
		editorMode:      ModeEdit,
		selectedTool:    ToolPencil,
		selectedColor:   color.RGBA{255, 255, 255, 255},
		brushSize:       1,
		canvasX:         50,
		canvasY:         50,
		canvasSize:      CanvasSize,
		pixelSize:       CanvasSize / SkinWidth,
		zoomLevel:       1.0,
		previewX:        600,
		previewY:        100,
		previewSize:     PreviewSize,
		previewAnimating: true,
		paletteX:        50,
		paletteY:        600,
		toolPanelX:      600,
		toolPanelY:      300,
		history:         make([]*SkinData, 0),
		historyIndex:    -1,
		maxHistory:      MaxHistory,
		skinDirectory:   "skins",
		backgroundColor:  color.RGBA{30, 30, 40, 255},
		gridColor:        color.RGBA{60, 60, 80, 255},
		selectionColor:  color.RGBA{100, 150, 255, 255},
		whiteImage:       whiteImage,
	}

	log.Printf("Initializing palette...")
	// Initialize default palette
	editor.initializePalette()
	
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
		{0, 0, 0, 255},         // Black
		{255, 255, 255, 255},   // White
		{128, 128, 128, 255},   // Gray
		{192, 192, 192, 255},   // Light gray
		{64, 64, 64, 255},      // Dark gray
		{255, 0, 0, 255},       // Red
		{255, 128, 0, 255},     // Orange
		{255, 255, 0, 255},     // Yellow
		{0, 255, 0, 255},       // Green
		{0, 255, 255, 255},     // Cyan
		{0, 0, 255, 255},       // Blue
		{128, 0, 255, 255},     // Purple
		{255, 0, 255, 255},     // Magenta
		{165, 42, 42, 255},     // Brown
		{255, 192, 203, 255},   // Pink
		{0, 128, 0, 255},       // Dark green
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
		// log.Printf("Drew pixel with color: RGB(%d,%d,%d)", 
		//	se.selectedColor.R, se.selectedColor.G, se.selectedColor.B)
	case ToolEraser:
		se.erasePixel(skinX, skinY)
		// log.Printf("Erased pixel at: %d,%d", skinX, skinY)
	case ToolEyedropper:
		se.pickColor(skinX, skinY)
	case ToolFill:
		se.fillArea(skinX, skinY)
	}
	
	se.isDrawing = true
	se.lastPixelX = skinX
	se.lastPixelY = skinY
}

// drawPixel draws a pixel at the specified position
func (se *SkinEditor) drawPixel(x, y int) {
	// Draw with brush size
	for dy := -se.brushSize / 2; dy <= se.brushSize/2; dy++ {
		for dx := -se.brushSize / 2; dx <= se.brushSize/2; dx++ {
			px, py := x+dx, y+dy
			if px >= 0 && px < SkinWidth && py >= 0 && py < SkinHeight {
				// Check if within circular brush
				if dx*dx+dy*dy <= (se.brushSize/2)*(se.brushSize/2) {
					se.skinData.Pixels[py][px] = se.selectedColor
				}
			}
		}
	}
	
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
	// Draw canvas
	se.drawCanvas(screen)
	
	// Draw preview
	se.drawPreview(screen)
	
	// Draw color palette
	se.drawPalette(screen)
	
	// Draw tool panel
	se.drawToolPanel(screen)
	
	// Draw UI elements
	se.drawUI(screen)
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
		}
		
		toolInfo := fmt.Sprintf("Tool: %s | Brush: %d | RGB(%d,%d,%d)", 
			toolName, se.brushSize, 
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
	
	log.Printf("Loaded skin: %s", skinName)
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
