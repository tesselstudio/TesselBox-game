package debug

import (
	"fmt"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// PerformanceProfiler tracks and displays performance metrics
type PerformanceProfiler struct {
	// Metrics
	fps           float64
	frameTime     time.Duration
	drawCalls     int
	memoryAlloc   uint64
	memorySys     uint64
	goroutines    int

	// Tracking
	lastFrameTime time.Time
	frameCount    int
	lastFPSUpdate time.Time

	// Settings
	Enabled       bool
	ShowDetailed  bool
}

// NewPerformanceProfiler creates a new profiler
func NewPerformanceProfiler() *PerformanceProfiler {
	return &PerformanceProfiler{
		Enabled:      false,
		ShowDetailed: false,
		lastFPSUpdate: time.Now(),
	}
}

// Update updates performance metrics
func (pp *PerformanceProfiler) Update() {
	if !pp.Enabled {
		return
	}

	now := time.Now()

	// Calculate frame time
	if !pp.lastFrameTime.IsZero() {
		pp.frameTime = now.Sub(pp.lastFrameTime)
	}
	pp.lastFrameTime = now

	// Update FPS every 500ms
	pp.frameCount++
	if now.Sub(pp.lastFPSUpdate) >= 500*time.Millisecond {
		pp.fps = float64(pp.frameCount) * 2.0 // Multiply by 2 because we measure over 500ms
		pp.frameCount = 0
		pp.lastFPSUpdate = now
	}

	// Get memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	pp.memoryAlloc = m.Alloc / 1024 / 1024 // MB
	pp.memorySys = m.Sys / 1024 / 1024     // MB

	// Get goroutine count
	pp.goroutines = runtime.NumGoroutine()
}

// RecordDrawCall should be called for each draw operation
func (pp *PerformanceProfiler) RecordDrawCall() {
	if pp.Enabled {
		pp.drawCalls++
	}
}

// ResetDrawCalls resets the draw call counter (call at start of frame)
func (pp *PerformanceProfiler) ResetDrawCalls() {
	pp.drawCalls = 0
}

// Draw renders the debug overlay
func (pp *PerformanceProfiler) Draw(screen *ebiten.Image) {
	if !pp.Enabled {
		return
	}

	// Background
	msg := ""

	// Basic info (always shown when enabled)
	msg += fmt.Sprintf("FPS: %.1f | Frame: %.2fms\n", pp.fps, float64(pp.frameTime.Microseconds())/1000.0)

	if pp.ShowDetailed {
		msg += fmt.Sprintf("Memory: %dMB / %dMB\n", pp.memoryAlloc, pp.memorySys)
		msg += fmt.Sprintf("Goroutines: %d\n", pp.goroutines)
		msg += fmt.Sprintf("Draw Calls: %d\n", pp.drawCalls)
		msg += fmt.Sprintf("TPS: %.1f\n", ebiten.ActualTPS())
	}

	// Draw with shadow for readability
	x, y := 10.0, 10.0
	ebitenutil.DebugPrintAt(screen, msg, int(x)+1, int(y)+1) // Shadow
	ebitenutil.DebugPrintAt(screen, msg, int(x), int(y))   // Text
}

// ToggleEnabled toggles the profiler on/off
func (pp *PerformanceProfiler) ToggleEnabled() {
	pp.Enabled = !pp.Enabled
}

// ToggleDetailed toggles detailed information
func (pp *PerformanceProfiler) ToggleDetailed() {
	pp.ShowDetailed = !pp.ShowDetailed
}

// SetEnabled sets the enabled state
func (pp *PerformanceProfiler) SetEnabled(enabled bool) {
	pp.Enabled = enabled
}

// GetMetrics returns current metrics as a string
func (pp *PerformanceProfiler) GetMetrics() string {
	return fmt.Sprintf("FPS: %.1f | Frame: %.2fms | Mem: %dMB | Goroutines: %d",
		pp.fps,
		float64(pp.frameTime.Microseconds())/1000.0,
		pp.memoryAlloc,
		pp.goroutines)
}

// IsSlowFrame returns true if the last frame took longer than target
func (pp *PerformanceProfiler) IsSlowFrame(targetFPS float64) bool {
	targetFrameTime := time.Second / time.Duration(targetFPS)
	return pp.frameTime > targetFrameTime
}
