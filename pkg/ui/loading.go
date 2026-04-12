package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// LoadingScreen shows generation progress
type LoadingScreen struct {
	IsVisible     bool
	Title         string
	Message       string
	Progress      float64 // 0.0 to 1.0
	Stage         int
	TotalStages   int
	StageNames    []string
	ScreenWidth   int
	ScreenHeight  int
}

// NewLoadingScreen creates a new loading screen
func NewLoadingScreen(width, height int) *LoadingScreen {
	return &LoadingScreen{
		IsVisible:    false,
		Title:        "Loading...",
		Message:      "Please wait",
		Progress:     0.0,
		Stage:        0,
		TotalStages:  4,
		StageNames:   []string{"Initializing", "Generating terrain", "Placing features", "Spawning entities"},
		ScreenWidth:  width,
		ScreenHeight: height,
	}
}

// Show displays the loading screen
func (l *LoadingScreen) Show(title string) {
	l.IsVisible = true
	l.Title = title
	l.Progress = 0.0
	l.Stage = 0
	l.Message = l.StageNames[0]
}

// Hide hides the loading screen
func (l *LoadingScreen) Hide() {
	l.IsVisible = false
}

// SetProgress updates the progress (0.0 to 1.0)
func (l *LoadingScreen) SetProgress(progress float64) {
	l.Progress = progress
	// Update stage based on progress
	stageIndex := int(progress * float64(l.TotalStages))
	if stageIndex >= l.TotalStages {
		stageIndex = l.TotalStages - 1
	}
	if stageIndex != l.Stage {
		l.Stage = stageIndex
		l.Message = l.StageNames[l.Stage]
	}
}

// SetMessage sets a custom message
func (l *LoadingScreen) SetMessage(msg string) {
	l.Message = msg
}

// Draw renders the loading screen
func (l *LoadingScreen) Draw(screen *ebiten.Image) {
	if !l.IsVisible {
		return
	}

	// Dark overlay
	vector.DrawFilledRect(screen, 0, 0, float32(l.ScreenWidth), float32(l.ScreenHeight), color.RGBA{0, 0, 0, 200}, false)

	// Center panel
	panelWidth := 400
	panelHeight := 200
	panelX := (l.ScreenWidth - panelWidth) / 2
	panelY := (l.ScreenHeight - panelHeight) / 2

	// Panel background
	vector.DrawFilledRect(screen, float32(panelX), float32(panelY), float32(panelWidth), float32(panelHeight), color.RGBA{40, 40, 40, 255}, false)

	// Panel border
	vector.StrokeRect(screen, float32(panelX), float32(panelY), float32(panelWidth), float32(panelHeight), 2, color.RGBA{100, 100, 100, 255}, false)

	// Title
	titleY := panelY + 40
	ebitenutil.DebugPrintAt(screen, l.Title, panelX+20, titleY)

	// Progress bar background
	barX := panelX + 50
	barY := panelY + 100
	barWidth := panelWidth - 100
	barHeight := 20
	vector.DrawFilledRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{60, 60, 60, 255}, false)

	// Progress bar fill
	fillWidth := int(float64(barWidth) * l.Progress)
	if fillWidth > 0 {
		vector.DrawFilledRect(screen, float32(barX), float32(barY), float32(fillWidth), float32(barHeight), color.RGBA{0, 200, 100, 255}, false)
	}

	// Progress bar border
	vector.StrokeRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), 1, color.RGBA{150, 150, 150, 255}, false)

	// Progress text
	progressText := fmt.Sprintf("%.0f%%", l.Progress*100)
	textX := barX + (barWidth-len(progressText)*7)/2
	ebitenutil.DebugPrintAt(screen, progressText, textX, barY+1)

	// Stage message
	msgY := barY + 40
	ebitenutil.DebugPrintAt(screen, l.Message, panelX+(panelWidth-len(l.Message)*7)/2, msgY)
}

// Update handles any animation updates
func (l *LoadingScreen) Update() {
	// Could add pulsing effects here if needed
}

// SetScreenSize updates the screen dimensions
func (l *LoadingScreen) SetScreenSize(width, height int) {
	l.ScreenWidth = width
	l.ScreenHeight = height
}
