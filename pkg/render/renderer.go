// renderer.go
package render

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"time"

	"tesselbox/pkg/blocks"
	"tesselbox/pkg/menu"
	"tesselbox/pkg/player"
	"tesselbox/pkg/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Colors
var (
	White    = color.RGBA{255, 255, 255, 255}
	Black    = color.RGBA{0, 0, 0, 255}
	DarkGray = color.RGBA{64, 64, 64, 255}
	Gray     = color.RGBA{128, 128, 128, 128}
	Green    = color.RGBA{100, 200, 100, 255}
	Red      = color.RGBA{200, 100, 100, 255}
	Blue     = color.RGBA{50, 150, 255, 255}
)

// Screen dimensions and defaults
const (
	ScreenWidth           = 1280
	ScreenHeight          = 720
	DefaultRenderDistance = 4
)

// Game represents the game state
type Game struct {
	ScreenWidth     int
	ScreenHeight    int
	World           *world.World
	Player          *player.Player
	CameraX         float64
	CameraY         float64
	Mining          bool
	MiningProgress  float64
	MiningStartTime time.Time
	renderDistance  int
	Particles       []*Particle
	inMenu          bool
	mainMenu        *menu.Menu
	lastUpdateTime  time.Time
}

// Particle represents a simple particle effect
type Particle struct {
	X, Y    float64
	VX, VY  float64
	Life    int
	MaxLife int
	Color   color.Color
}

// NewGame creates a new game instance (no auth anymore)
func NewGame() *Game {
	game := &Game{
		ScreenWidth:    ScreenWidth,
		ScreenHeight:   ScreenHeight,
		World:          world.NewWorld("default"),
		CameraX:        0,
		CameraY:        0,
		Mining:         false,
		MiningProgress: 0,
		renderDistance: DefaultRenderDistance,
		Particles:      []*Particle{},
		inMenu:         true,
		mainMenu:       menu.NewMenu(),
		lastUpdateTime: time.Now(),
	}

	return game
}

// Layout returns the game's logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.ScreenWidth, g.ScreenHeight
}

// Update updates the game state
func (g *Game) Update() error {
	// Calculate delta time
	now := time.Now()
	deltaTime := now.Sub(g.lastUpdateTime).Seconds()
	g.lastUpdateTime = now

	// Handle menu interaction
	if g.inMenu {
		action := g.mainMenu.Update()
		if action != menu.ActionNone {
			g.handleMenuAction(action)
		}
		return nil
	}

	// Exit to menu with Escape
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.inMenu = true
		return nil
	}

	// Player movement input handling
	if g.Player != nil {
		g.Player.MovingLeft = ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
		g.Player.MovingRight = ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight)
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			g.Player.Jump()
		}

		g.Player.Update(deltaTime)

		// Update Camera to follow player
		g.CameraX = g.Player.X - float64(g.ScreenWidth)/2
		g.CameraY = g.Player.Y - float64(g.ScreenHeight)/2
	}

	// Mining logic
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		worldX := float64(mx) + g.CameraX
		worldY := float64(my) + g.CameraY

		hex := g.World.GetHexagonAt(worldX, worldY)
		if hex != nil && g.Player != nil {
			dx := hex.X - g.Player.X
			dy := hex.Y - g.Player.Y
			distanceSq := dx*dx + dy*dy
			if distanceSq <= player.MiningRange*player.MiningRange {
				if !g.Mining {
					g.Mining = true
					g.MiningProgress = 0
					g.MiningStartTime = time.Now()
				}
				g.MiningProgress += deltaTime
				if g.MiningProgress >= 0.5 {
					g.World.RemoveHexagonAt(worldX, worldY)
					g.createExplosion(hex.X, hex.Y, Gray)
					g.Mining = false
					g.MiningProgress = 0
				}
			} else {
				g.Mining = false
				g.MiningProgress = 0
			}
		}
	} else {
		g.Mining = false
		g.MiningProgress = 0
	}

	// Update Particles
	for i := len(g.Particles) - 1; i >= 0; i-- {
		p := g.Particles[i]
		p.X += p.VX * deltaTime
		p.Y += p.VY * deltaTime
		p.Life--
		if p.Life <= 0 {
			g.Particles = append(g.Particles[:i], g.Particles[i+1:]...)
		}
	}

	return nil
}

// handleMenuAction processes button clicks from the menu (cleaned, no auth)
func (g *Game) handleMenuAction(action menu.MenuAction) {
	switch action {
	case menu.ActionStartGame:
		g.inMenu = false
		g.spawnPlayer()
		g.World.GetChunksInRange(g.Player.X, g.Player.Y)

	case menu.ActionOpenSettings:
		g.mainMenu.SetSettingsMenu()

	case menu.ActionBack:
		g.mainMenu.SetMainMenu()

	case menu.ActionExit:
		os.Exit(0)

	// Removed login/logout actions
	default:
		fmt.Printf("Unhandled menu action: %v\n", action)
	}
}

func (g *Game) spawnPlayer() {
	g.Player = player.NewPlayer(100, 100)
}

func (g *Game) createExplosion(x, y float64, c color.Color) {
	for i := 0; i < 10; i++ {
		g.Particles = append(g.Particles, &Particle{
			X:       x,
			Y:       y,
			VX:      (rand.Float64() - 0.5) * 400,
			VY:      (rand.Float64() - 0.5) * 400,
			Life:    60,
			MaxLife: 60,
			Color:   c,
		})
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.inMenu {
		g.mainMenu.Draw(screen)
		return
	}

	screen.Fill(Blue)

	if g.Player != nil {
		g.World.GetChunksInRange(g.Player.X, g.Player.Y)
	}

	visibleHexagons := g.World.GetVisibleBlocks(g.CameraX+float64(g.ScreenWidth)/2, g.CameraY+float64(g.ScreenHeight)/2)
	for _, hex := range visibleHexagons {
		drawX := hex.X - g.CameraX
		drawY := hex.Y - g.CameraY
		ebitenutil.DrawRect(screen, drawX-20, drawY-20, 40, 40, hex.Color)
	}

	if g.Player != nil {
		pX := g.Player.X - g.CameraX
		pY := g.Player.Y - g.CameraY
		ebitenutil.DrawRect(screen, pX-20, pY-40, 40, 80, Red)
	}

	for _, p := range g.Particles {
		drawX := p.X - g.CameraX
		drawY := p.Y - g.CameraY
		ebitenutil.DrawRect(screen, drawX-1, drawY-1, 2, 2, p.Color)
	}

	g.drawUI(screen)
}

func (g *Game) drawUI(screen *ebiten.Image) {
	hotbarWidth := 400.0
	hotbarHeight := 40.0
	hotbarX := (float64(g.ScreenWidth) - hotbarWidth) / 2
	hotbarY := float64(g.ScreenHeight) - hotbarHeight - 10
	slotWidth := hotbarWidth / 9

	for i := 0; i < 9; i++ {
		slotX := hotbarX + float64(i)*slotWidth
		if g.Player != nil && i == g.Player.SelectedSlot {
			ebitenutil.DrawRect(screen, slotX, hotbarY, slotWidth-2, hotbarHeight, White)
		} else {
			ebitenutil.DrawRect(screen, slotX, hotbarY, slotWidth-2, hotbarHeight, Gray)
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// getBlockKeyFromType converts a BlockType to its string key
func getBlockKeyFromType(blockType blocks.BlockType) string {
	switch blockType {
	case blocks.AIR:
		return "air"
	case blocks.DIRT:
		return "dirt"
	case blocks.GRASS:
		return "grass"
	case blocks.STONE:
		return "stone"
	case blocks.SAND:
		return "sand"
	case blocks.WATER:
		return "water"
	case blocks.LOG:
		return "log"
	case blocks.LEAVES:
		return "leaves"
	case blocks.COAL_ORE:
		return "coal_ore"
	case blocks.IRON_ORE:
		return "iron_ore"
	case blocks.GOLD_ORE:
		return "gold_ore"
	case blocks.DIAMOND_ORE:
		return "diamond_ore"
	case blocks.BEDROCK:
		return "bedrock"
	case blocks.GLASS:
		return "glass"
	case blocks.BRICK:
		return "brick"
	case blocks.PLANK:
		return "plank"
	case blocks.CACTUS:
		return "cactus"
	default:
		return "dirt"
	}
}
