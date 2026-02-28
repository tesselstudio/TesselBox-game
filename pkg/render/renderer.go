// renderer.go
package render

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"
	"time"

	"tesselbox/pkg/blocks"
	"tesselbox/pkg/creatures"
	"tesselbox/pkg/gametime"
	"tesselbox/pkg/hexagon"
	"tesselbox/pkg/items"
	"tesselbox/pkg/menu"
	"tesselbox/pkg/organisms"
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
	Particles       *ParticlePool
	inMenu          bool
	mainMenu        *menu.Menu
	lastUpdateTime  time.Time
	whiteImage      *ebiten.Image // Reusable image for drawing
	dayNightCycle   *gametime.DayNightCycle
}

// Particle represents a simple particle effect
type Particle struct {
	X, Y    float64
	VX, VY  float64
	Life    int
	MaxLife int
	Color   color.Color
	Active  bool // Whether this particle is currently in use
}

// ParticlePool manages a pool of reusable particles
type ParticlePool struct {
	particles    []*Particle
	maxParticles int
	nextIndex    int
}

// NewParticlePool creates a new particle pool with pre-allocated particles
func NewParticlePool(maxParticles int) *ParticlePool {
	pool := &ParticlePool{
		particles:    make([]*Particle, maxParticles),
		maxParticles: maxParticles,
		nextIndex:    0,
	}

	// Pre-allocate all particles
	for i := 0; i < maxParticles; i++ {
		pool.particles[i] = &Particle{Active: false}
	}

	return pool
}

// GetInactiveParticle returns an inactive particle from the pool, or nil if none available
func (pp *ParticlePool) GetInactiveParticle() *Particle {
	// Start from nextIndex and wrap around to find an inactive particle
	startIndex := pp.nextIndex
	for i := 0; i < pp.maxParticles; i++ {
		idx := (startIndex + i) % pp.maxParticles
		if !pp.particles[idx].Active {
			pp.nextIndex = (idx + 1) % pp.maxParticles
			return pp.particles[idx]
		}
	}
	return nil // No inactive particles available
}

// Update updates all active particles in the pool
func (pp *ParticlePool) Update(deltaTime float64) {
	for _, p := range pp.particles {
		if p.Active {
			p.X += p.VX * deltaTime
			p.Y += p.VY * deltaTime
			p.Life--
			if p.Life <= 0 {
				p.Active = false
			}
		}
	}
}

// GetActiveParticles returns a slice of all active particles
func (pp *ParticlePool) GetActiveParticles() []*Particle {
	active := make([]*Particle, 0, pp.maxParticles)
	for _, p := range pp.particles {
		if p.Active {
			active = append(active, p)
		}
	}
	return active
}

// GetActiveCount returns the number of currently active particles
func (pp *ParticlePool) GetActiveCount() int {
	count := 0
	for _, p := range pp.particles {
		if p.Active {
			count++
		}
	}
	return count
}

// NewGame creates a new game instance (no auth anymore)
func NewGame() *Game {
	// Create a white image for drawing hexagons
	whiteImage := ebiten.NewImage(1, 1)
	whiteImage.Fill(color.RGBA{255, 255, 255, 255})

	game := &Game{
		ScreenWidth:    ScreenWidth,
		ScreenHeight:   ScreenHeight,
		World:          world.NewWorld("default"),
		CameraX:        0,
		CameraY:        0,
		Mining:         false,
		MiningProgress: 0,
		renderDistance: DefaultRenderDistance,
		Particles:      NewParticlePool(1000), // Pre-allocate 1000 particles
		inMenu:         true,
		mainMenu:       menu.NewMenu(),
		lastUpdateTime: time.Now(),
		whiteImage:     whiteImage,
		dayNightCycle:  gametime.NewDayNightCycle(600.0), // 10 minutes day/night cycle
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

		// Update weapon direction based on mouse position
		mx, my := ebiten.CursorPosition()
		mouseWorldX := float64(mx) + g.CameraX
		mouseWorldY := float64(my) + g.CameraY
		dx := mouseWorldX - g.Player.X
		dy := mouseWorldY - g.Player.Y
		distance := math.Sqrt(dx*dx + dy*dy)
		if distance > 0 {
			g.Player.WeaponInstance.DirectionX = dx / distance
			g.Player.WeaponInstance.DirectionY = dy / distance
			g.Player.WeaponInstance.Angle = math.Atan2(dy, dx)
		}

		// Automatic weapon attack when mouse is over an organism
		if g.Player.WeaponInstance != nil && !g.Player.WeaponInstance.Swinging {
			organism := g.World.GetOrganismAt(mouseWorldX, mouseWorldY, 20.0) // 20 pixel tolerance
			if organism != nil && organism.IsAlive() {
				// Check if organism is within weapon range from player
				orgDx := organism.X - g.Player.X
				orgDy := organism.Y - g.Player.Y
				orgDistance := math.Sqrt(orgDx*orgDx + orgDy*orgDy)
				if orgDistance <= g.Player.WeaponInstance.Range {
					g.Player.WeaponInstance.Swinging = true
					g.Player.WeaponInstance.SwingProgress = 0
				}
			}
		}

		// Update weapon swing animation
		if g.Player.WeaponInstance.Swinging {
			g.Player.WeaponInstance.SwingProgress += deltaTime * g.Player.WeaponInstance.SwingSpeed
			if g.Player.WeaponInstance.SwingProgress >= 1.0 {
				g.Player.WeaponInstance.Swinging = false
				g.Player.WeaponInstance.SwingProgress = 0
			}

			// Check for weapon-organism collisions during swing
			nearbyOrganisms := g.World.GetNearbyOrganisms(g.Player.X, g.Player.Y, g.Player.WeaponInstance.Range+50)
			for _, org := range nearbyOrganisms {
				if org.IsAlive() {
					// Simple distance check from player to organism (could be improved with swing arc)
					dx := org.X - g.Player.X
					dy := org.Y - g.Player.Y
					distance := math.Sqrt(dx*dx + dy*dy)
					if distance <= g.Player.WeaponInstance.Range {
						// Damage the organism
						org.TakeDamage(g.Player.WeaponInstance.Damage)
						g.createExplosion(org.X, org.Y, Red)
						break // Only hit one organism per swing
					}
				}
			}
		}

		// Mining logic with right mouse button
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
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

		// Apply collision detection with hexagonal tiles
		nearbyHexagons := g.World.GetNearbyHexagons(g.Player.X, g.Player.Y, 300)
		g.Player.UpdateWithCollision(deltaTime, func(minX, minY, maxX, maxY float64) bool {
			for _, hex := range nearbyHexagons {
				def := blocks.BlockDefinitions[getBlockKeyFromType(hex.BlockType)]
				if def == nil || !def.Solid {
					continue
				}

				// Check collision with hexagon using point-in-hexagon test
				if g.checkRectHexCollision(minX, minY, maxX, maxY, hex) {
					return true
				}
			}
			return false
		})

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
	g.Particles.Update(deltaTime)

	// Update creatures if we have a world and player
	if g.World != nil && g.Player != nil {
		playerX, playerY := g.Player.GetCenter()

		// Spawn creatures periodically
		g.World.SpawnCreatures(g.dayNightCycle, playerX, playerY)

		// Update all creatures with AI
		g.World.UpdateCreatures(playerX, playerY, deltaTime)

		// Handle creature attacks on player
		for _, creature := range g.World.Creatures {
			if creature.IsAlive() {
				creature.AttackPlayer(playerX, playerY, func(damage, fromX, fromY, knockback float64) {
					if g.Player != nil {
						g.Player.TakeDamage(damage, fromX, fromY, knockback)
					}
				})
			}
		}

		// Handle player attacks on creatures
		if g.Player.IsAttacking() {
			attackRange := g.Player.GetAttackRange()
			nearbyCreatures := g.World.GetCreaturesInArea(playerX, playerY, attackRange)
			for _, creature := range nearbyCreatures {
				if creature.IsAlive() {
					// Check if creature is in attack direction (simplified)
					dx := creature.X - playerX
					dy := creature.Y - playerY
					distance := math.Sqrt(dx*dx + dy*dy)
					if distance <= attackRange {
						damage := g.Player.GetAttackDamage()
						g.Player.DealDamage(creature, damage)

						// Create hit effect
						g.createExplosion(creature.X, creature.Y, Red)
						break // Only hit one creature per attack
					}
				}
			}
		}

		// Remove dead creatures and handle loot
		for _, creature := range g.World.Creatures {
			if !creature.IsAlive() {
				// Drop loot
				loot := creature.GetLootDrops()
				for _, item := range loot {
					// Add to player's inventory (simplified - just add to inventory)
					if g.Player != nil {
						// For now, just print what was dropped
						// In a full implementation, this would add to inventory or drop on ground
						fmt.Printf("Creature dropped: %s x%d\n", items.ItemNameByID(item.Type), item.Quantity)
					}
				}

				// Create death particles
				g.createExplosion(creature.X, creature.Y, color.RGBA{128, 128, 128, 255})
			}
		}
		g.World.RemoveDeadCreatures()
	}

	// Handle dead organisms and remove them
	if g.World != nil {
		var aliveOrganisms []*organisms.Organism
		for _, org := range g.World.Organisms {
			if !org.IsAlive() {
				// Drop items
				drops := organisms.GetDrops(org)
				for _, itemName := range drops {
					// For now, just print what was dropped (similar to creatures)
					fmt.Printf("Organism dropped: %s\n", itemName)
				}

				// Create death particles
				g.createExplosion(org.X, org.Y, color.RGBA{128, 128, 128, 255})
			} else {
				aliveOrganisms = append(aliveOrganisms, org)
			}
		}
		g.World.Organisms = aliveOrganisms
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
		particle := g.Particles.GetInactiveParticle()
		if particle == nil {
			break // No available particles in pool
		}

		particle.X = x
		particle.Y = y
		particle.VX = (rand.Float64() - 0.5) * 400
		particle.VY = (rand.Float64() - 0.5) * 400
		particle.Life = 60
		particle.MaxLife = 60
		particle.Color = c
		particle.Active = true
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
		if hex.BlockType == blocks.AIR {
			continue
		}
		g.drawHexagon(screen, hex)
	}

	if g.Player != nil {
		// Draw square player
		pX := g.Player.X - g.CameraX
		pY := g.Player.Y - g.CameraY
		ebitenutil.DrawRect(screen, pX, pY, g.Player.Width, g.Player.Height, Red)

		// Draw weapon
		if g.Player.WeaponInstance != nil {
			g.drawWeapon(screen, g.Player.WeaponInstance, g.Player.X, g.Player.Y, g.CameraX, g.CameraY)
		}
	}

	// Draw creatures
	if g.World != nil {
		for _, creature := range g.World.Creatures {
			cX := creature.X - g.CameraX
			cY := creature.Y - g.CameraY

			// Skip if off-screen
			if cX < -50 || cX > float64(g.ScreenWidth)+50 ||
				cY < -50 || cY > float64(g.ScreenHeight)+50 {
				continue
			}

			// Choose color based on creature type
			var creatureColor color.RGBA
			switch creature.Type {
			case creatures.SLIME:
				creatureColor = color.RGBA{0, 255, 0, 255} // Green
			case creatures.SPIDER:
				creatureColor = color.RGBA{0, 0, 0, 255} // Black
			case creatures.ZOMBIE:
				creatureColor = color.RGBA{0, 100, 0, 255} // Dark green
			default:
				creatureColor = color.RGBA{128, 128, 128, 255} // Gray
			}

			// Draw creature as a circle (approximated with rectangle for simplicity)
			size := 20.0
			ebitenutil.DrawRect(screen, cX-size/2, cY-size/2, size, size, creatureColor)

			// Draw health bar above creature if damaged
			if creature.Health < creature.MaxHealth {
				barWidth := size
				barHeight := 4.0
				barX := cX - barWidth/2
				barY := cY - size/2 - 8

				// Background (red)
				ebitenutil.DrawRect(screen, barX, barY, barWidth, barHeight, color.RGBA{255, 0, 0, 255})
				// Health (green)
				healthWidth := barWidth * (creature.Health / creature.MaxHealth)
				ebitenutil.DrawRect(screen, barX, barY, healthWidth, barHeight, color.RGBA{0, 255, 0, 255})
			}
		}
	}

	// Draw organisms
	if g.World != nil {
		for _, org := range g.World.Organisms {
			oX := org.X - g.CameraX
			oY := org.Y - g.CameraY

			// Skip if off-screen
			if oX < -50 || oX > float64(g.ScreenWidth)+50 ||
				oY < -50 || oY > float64(g.ScreenHeight)+50 {
				continue
			}

			// Choose color based on organism type
			var orgColor color.RGBA
			switch org.Type {
			case organisms.TREE:
				orgColor = color.RGBA{139, 69, 19, 255} // Brown for tree trunk
			case organisms.BUSH:
				orgColor = color.RGBA{34, 139, 34, 255} // Forest green for bush
			case organisms.FLOWER:
				orgColor = color.RGBA{255, 0, 255, 255} // Magenta for flower
			default:
				orgColor = color.RGBA{128, 128, 128, 255} // Gray
			}

			// Draw organism as a small rectangle
			size := 10.0
			ebitenutil.DrawRect(screen, oX-size/2, oY-size/2, size, size, orgColor)

			// Draw health bar above organism if damaged
			if org.Health < org.MaxHealth {
				barWidth := size
				barHeight := 3.0
				barX := oX - barWidth/2
				barY := oY - size/2 - 6

				// Background (red)
				ebitenutil.DrawRect(screen, barX, barY, barWidth, barHeight, color.RGBA{255, 0, 0, 255})
				// Health (green)
				healthWidth := barWidth * (org.Health / org.MaxHealth)
				ebitenutil.DrawRect(screen, barX, barY, healthWidth, barHeight, color.RGBA{0, 255, 0, 255})
			}
		}
	}

	for _, p := range g.Particles.GetActiveParticles() {
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

// drawHexagon draws a hexagonal tile
func (g *Game) drawHexagon(screen *ebiten.Image, hex *world.Hexagon) {
	// Calculate screen position
	screenX := hex.X - g.CameraX
	screenY := hex.Y - g.CameraY

	// Check if hexagon is on screen
	if screenX < -100 || screenX > float64(g.ScreenWidth)+100 ||
		screenY < -100 || screenY > float64(g.ScreenHeight)+100 {
		return
	}

	// Get hexagon corners
	corners := hexagon.GetHexCorners(screenX, screenY, world.HexSize)

	// Get block color
	def := blocks.BlockDefinitions[getBlockKeyFromType(hex.BlockType)]
	if def == nil {
		return
	}

	// Create polygon vertices
	vertices := make([]ebiten.Vertex, len(corners))
	for i, corner := range corners {
		r, g, b, a := def.Color.RGBA()
		vertices[i] = ebiten.Vertex{
			DstX:   float32(corner[0]),
			DstY:   float32(corner[1]),
			ColorR: float32(r) / 65535.0,
			ColorG: float32(g) / 65535.0,
			ColorB: float32(b) / 65535.0,
			ColorA: float32(a) / 65535.0,
		}
	}

	// Apply damage darkening
	if hex.Health < hex.MaxHealth {
		damageRatio := hex.Health / hex.MaxHealth
		for i := range vertices {
			vertices[i].ColorR *= float32(damageRatio)
			vertices[i].ColorG *= float32(damageRatio)
			vertices[i].ColorB *= float32(damageRatio)
		}
	}

	// Draw filled hexagon using triangles
	indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5}
	screen.DrawTriangles(vertices, indices, g.whiteImage, nil)

	// Hexagon outline can be added later if needed
	// For now, the filled hexagon is sufficient
}

// checkRectHexCollision checks if a rectangle collides with a hexagon
func (g *Game) checkRectHexCollision(minX, minY, maxX, maxY float64, hex *world.Hexagon) bool {
	// Quick bounding box check first
	hexRadius := hex.Size
	hexMinX := hex.X - hexRadius
	hexMaxX := hex.X + hexRadius
	hexMinY := hex.Y - hexRadius
	hexMaxY := hex.Y + hexRadius

	if maxX < hexMinX || minX > hexMaxX || maxY < hexMinY || minY > hexMaxY {
		return false
	}

	// Check if any corner of the rectangle is inside the hexagon
	corners := [][2]float64{
		{minX, minY}, {maxX, minY}, {maxX, maxY}, {minX, maxY},
	}

	for _, corner := range corners {
		if g.pointInHexagon(corner[0], corner[1], hex) {
			return true
		}
	}

	// Check if any hexagon corner is inside the rectangle
	for _, hexCorner := range hex.Corners {
		if hexCorner[0] >= minX && hexCorner[0] <= maxX &&
			hexCorner[1] >= minY && hexCorner[1] <= maxY {
			return true
		}
	}

	// Check if rectangle edges intersect with hexagon edges
	rectEdges := [][4]float64{
		{minX, minY, maxX, minY}, // top
		{maxX, minY, maxX, maxY}, // right
		{maxX, maxY, minX, maxY}, // bottom
		{minX, maxY, minX, minY}, // left
	}

	for _, edge := range rectEdges {
		for i := 0; i < len(hex.Corners); i++ {
			next := (i + 1) % len(hex.Corners)
			if g.lineIntersect(edge[0], edge[1], edge[2], edge[3],
				hex.Corners[i][0], hex.Corners[i][1],
				hex.Corners[next][0], hex.Corners[next][1]) {
				return true
			}
		}
	}

	return false
}

// pointInHexagon checks if a point is inside a hexagon
func (g *Game) pointInHexagon(x, y float64, hex *world.Hexagon) bool {
	dx := x - hex.X
	dy := y - hex.Y
	distanceSq := dx*dx + dy*dy
	radiusSq := hex.Size * hex.Size

	// Quick circle check first
	if distanceSq > radiusSq {
		return false
	}

	// More precise hexagon check using ray casting
	inside := false
	j := len(hex.Corners) - 1
	for i := 0; i < len(hex.Corners); i++ {
		xi, yi := hex.Corners[i][0], hex.Corners[i][1]
		xj, yj := hex.Corners[j][0], hex.Corners[j][1]

		if ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}

	return inside
}

// lineIntersect checks if two line segments intersect
func (g *Game) lineIntersect(x1, y1, x2, y2, x3, y3, x4, y4 float64) bool {
	denom := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
	if denom == 0 {
		return false
	}

	t := ((x1-x3)*(y3-y4) - (y1-y3)*(x3-x4)) / denom
	u := -((x1-x2)*(y1-y3) - (y1-y2)*(x1-x3)) / denom

	return t >= 0 && t <= 1 && u >= 0 && u <= 1
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

// drawWeapon renders the player's weapon with visual representation
func (g *Game) drawWeapon(screen *ebiten.Image, weapon *player.Weapon, playerX, playerY, cameraX, cameraY float64) {
	// Calculate weapon base position (player center)
	baseX := playerX + g.Player.Width/2
	baseY := playerY + g.Player.Height/2

	// Calculate weapon end point
	endX := baseX + weapon.DirectionX*weapon.Length
	endY := baseY + weapon.DirectionY*weapon.Length

	// Apply swing animation
	if weapon.Swinging {
		swingAngle := weapon.Angle + (weapon.SwingProgress-0.5)*math.Pi*0.5 // Swing arc
		swingDirX := math.Cos(swingAngle)
		swingDirY := math.Sin(swingAngle)
		endX = baseX + swingDirX*weapon.Length
		endY = baseY + swingDirY*weapon.Length
	}

	// Convert to screen coordinates
	screenBaseX := baseX - cameraX
	screenBaseY := baseY - cameraY
	screenEndX := endX - cameraX
	screenEndY := endY - cameraY

	// Draw weapon as a thick line with a blade end
	bladeColor := color.RGBA{255, 255, 255, 255} // White for blade tip

	// Draw weapon shaft (thicker line)
	shaftWidth := 4.0
	perpX := -weapon.DirectionY * shaftWidth / 2
	perpY := weapon.DirectionX * shaftWidth / 2

	// Create a rectangle for the weapon shaft
	shaftPoints := []ebiten.Vertex{
		{DstX: float32(screenBaseX + perpX), DstY: float32(screenBaseY + perpY), ColorR: 0.75, ColorG: 0.75, ColorB: 0.75, ColorA: 1.0},
		{DstX: float32(screenBaseX - perpX), DstY: float32(screenBaseY - perpY), ColorR: 0.75, ColorG: 0.75, ColorB: 0.75, ColorA: 1.0},
		{DstX: float32(screenEndX - perpX), DstY: float32(screenEndY - perpY), ColorR: 0.75, ColorG: 0.75, ColorB: 0.75, ColorA: 1.0},
		{DstX: float32(screenEndX + perpX), DstY: float32(screenEndY + perpY), ColorR: 0.75, ColorG: 0.75, ColorB: 0.75, ColorA: 1.0},
	}

	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(shaftPoints, indices, g.whiteImage, nil)

	// Draw blade tip (brighter, larger)
	bladeSize := 8.0
	if weapon.Swinging {
		bladeSize = bladeSize * (1.0 + weapon.SwingProgress*0.5) // Grow during swing
	}

	ebitenutil.DrawRect(screen, screenEndX-bladeSize/2, screenEndY-bladeSize/2, bladeSize, bladeSize, bladeColor)

	// Add glow effect when swinging
	if weapon.Swinging {
		glowAlpha := 0.3 * (1.0 - weapon.SwingProgress) // Fade out during swing
		glowColor := color.RGBA{255, 255, 200, uint8(glowAlpha * 255)}
		ebitenutil.DrawRect(screen, screenEndX-bladeSize, screenEndY-bladeSize, bladeSize*2, bladeSize*2, glowColor)
	}
}
