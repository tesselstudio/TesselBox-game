package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"tesselbox/pkg/auth"
	"tesselbox/pkg/blocks"
	"tesselbox/pkg/crafting"
	"tesselbox/pkg/hexagon"
	"tesselbox/pkg/items"
	"tesselbox/pkg/menu"
	"tesselbox/pkg/player"
	"tesselbox/pkg/world"
)

// stringToBlockType converts a string block type name to blocks.BlockType
func stringToBlockType(blockTypeStr string) blocks.BlockType {
	blockMap := map[string]blocks.BlockType{
		"air":         blocks.AIR,
		"dirt":        blocks.DIRT,
		"grass":       blocks.GRASS,
		"stone":       blocks.STONE,
		"sand":        blocks.SAND,
		"water":       blocks.WATER,
		"log":         blocks.LOG,
		"leaves":      blocks.LEAVES,
		"coal_ore":    blocks.COAL_ORE,
		"iron_ore":    blocks.IRON_ORE,
		"gold_ore":    blocks.GOLD_ORE,
		"diamond_ore": blocks.DIAMOND_ORE,
		"bedrock":     blocks.BEDROCK,
		"glass":       blocks.GLASS,
		"brick":       blocks.BRICK,
		"plank":       blocks.PLANK,
		"cactus":      blocks.CACTUS,
	}
	if bt, ok := blockMap[blockTypeStr]; ok {
		return bt
	}
	return blocks.AIR
}

// minFloat32 returns the minimum of two float32 values
func minFloat32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
	FPS          = 60
)

// Game represents the game state
type Game struct {
	// Game state
	world     *world.World
	player    *player.Player
	inventory *items.Inventory

	// Systems
	craftingSystem *crafting.CraftingSystem
	craftingUI     *crafting.CraftingUI
	menu           *menu.Menu
	authService    *auth.AuthService

	// Camera
	cameraX, cameraY float64

	// Mouse
	mouseX, mouseY int

	// Game state flags
	inMenu     bool
	inGame     bool
	inCrafting bool

	// Timing
	lastTime time.Time

	// For solid color drawing
	whiteImage *ebiten.Image
}

// NewGame creates a new game
func NewGame() *Game {
	// Create a 1x1 white image for solid color drawing
	whiteImage := ebiten.NewImage(1, 1)
	whiteImage.Fill(color.RGBA{255, 255, 255, 255})

	g := &Game{
		world:     world.NewWorld(), // Random seed
		player:    player.NewPlayer(0, 0),
		inventory: items.NewInventory(32),
		cameraX:   0,
		cameraY:   0,
		lastTime:  time.Now(),
		whiteImage: whiteImage,
	}

	// Initialize crafting system
	g.craftingSystem = crafting.NewCraftingSystem()
	if err := g.craftingSystem.LoadRecipes("config/crafting_recipes.json"); err != nil {
		log.Printf("Warning: Failed to load crafting recipes: %v", err)
	}
	g.craftingUI = crafting.NewCraftingUI(g.craftingSystem, g.inventory)

	// Initialize menu
	g.menu = menu.NewMenu()

	// Initialize OAuth (optional - will work without credentials)
	authConfig, err := auth.LoadConfig("config/oauth_config.json")
	if err != nil {
		log.Printf("Could not load OAuth config: %v. Running without authentication.", err)
	} else {
		g.authService, err = auth.NewAuthService(authConfig)
		if err != nil {
			log.Printf("Could not create auth service: %v. Running without authentication.", err)
			g.authService = nil // Ensure auth service is nil if it fails to create
		}
	}

	// Set default hotbar items
	defaultItems := items.DefaultHotbarItems()
	for i, item := range defaultItems {
		if i < len(g.inventory.Slots) {
			g.inventory.Slots[i] = item
		}
	}

	// Start in menu
	g.inMenu = true
	g.inGame = false

	return g
}

// Update updates the game state
func (g *Game) Update() error {
	// Calculate delta time for framerate-independent movement
	currentTime := time.Now()
	deltaTime := currentTime.Sub(g.lastTime).Seconds()
	g.lastTime = currentTime

	// Handle menu state
	if g.inMenu {
		action := g.menu.Update()
		g.handleMenuAction(action)
		return nil
	}

	// Handle crafting UI
	if g.inCrafting {
		if err := g.craftingUI.Update(); err != nil {
			return err
		}

		// Handle escape to close crafting
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.inCrafting = false
		}
		return nil
	}

	// Handle game state
	if g.inGame {
		g.handleGameInput()

		// Update player with delta time (framerate-independent)
		g.player.Update(deltaTime)

		// Apply collision-aware position update
		nearbyHexagons := g.world.GetNearbyHexagons(g.player.X, g.player.Y, 300)
		g.player.UpdateWithCollision(deltaTime, func(minX, minY, maxX, maxY float64) bool {
			for _, hex := range nearbyHexagons {
				def := blocks.BlockDefinitions[getBlockKeyFromType(hex.BlockType)]
				if def == nil || !def.Solid {
					continue
				}

				hexMinX := hex.X - hex.Size
				hexMinY := hex.Y - hex.Size
				hexMaxX := hex.X + hex.Size
				hexMaxY := hex.Y + hex.Size

				if !(maxX < hexMinX || minX > hexMaxX || maxY < hexMinY || minY > hexMaxY) {
					return true
				}
			}
			return false
		})

		// Update camera to follow player
		centerX, centerY := g.player.GetCenter()
		g.cameraX = centerX - ScreenWidth/2
		g.cameraY = centerY - ScreenHeight/2

		// Generate chunks around player
		g.world.GetChunksInRange(centerX, centerY)
	}

	return nil
}

// handleMenuAction handles menu actions
func (g *Game) handleMenuAction(action menu.MenuAction) {
	switch action {
	case menu.ActionLoginGoogle:
		if g.authService == nil {
			log.Println("Auth service not available")
			return
		}
		// Start OAuth server and open browser for Google login
		if err := g.authService.StartServer(":8080"); err != nil {
			log.Printf("Warning: Failed to start OAuth server: %v", err)
		}
		authURL := g.authService.GetAuthURL(auth.ProviderGoogle)
		if authURL != "" {
			if err := openURL(authURL); err != nil {
				log.Printf("Warning: Failed to open browser: %v", err)
			}
		}

	case menu.ActionLoginGitHub:
		if g.authService == nil {
			log.Println("Auth service not available")
			return
		}
		// Start OAuth server and open browser for GitHub login
		if err := g.authService.StartServer(":8080"); err != nil {
			log.Printf("Warning: Failed to start OAuth server: %v", err)
		}
		authURL := g.authService.GetAuthURL(auth.ProviderGitHub)
		if authURL != "" {
			if err := openURL(authURL); err != nil {
				log.Printf("Warning: Failed to open browser: %v", err)
			}
		}

	case menu.ActionStartGame:
		g.inMenu = false
		g.inGame = true
		g.player.SetPosition(0, 0)
		g.player.SetVelocity(0, 0)

	case menu.ActionOpenSettings:
		g.menu.SetSettingsMenu()

	case menu.ActionOpenLogin:
		g.menu.SetLoginMenu()
		if g.authService != nil {
			// Start OAuth server when login menu is opened
			if err := g.authService.StartServer(":8080"); err != nil {
				log.Printf("Warning: Failed to start OAuth server: %v", err)
			}
			g.authService.Logout() // Clear any previous session
		}

	case menu.ActionLogout:
		if g.authService != nil {
			g.authService.Logout()
		}
		g.menu.SetMainMenu() // Return to main menu
		g.inGame = false
		g.inMenu = true

	case menu.ActionBack:
		if g.menu.CurrentMenu == menu.MenuSettings {
			g.menu.SetMainMenu()
		} else {
			g.menu.SetMainMenu()
		}

	case menu.ActionExit:
		// Stop the OAuth server before exiting
		if g.authService != nil {
			g.authService.StopServer()
		}
		os.Exit(0)
		return
	}
}

// handleGameInput handles game input
func (g *Game) handleGameInput() {
	// Handle keyboard input
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.MovingLeft = true
		g.player.MovingRight = false
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.MovingRight = true
		g.player.MovingLeft = false
	} else {
		g.player.MovingLeft = false
		g.player.MovingRight = false
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Jump()
	}

	// Open crafting menu
	if inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.inCrafting = true
		g.craftingUI.Toggle()
	}

	// Return to menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.inGame = false
		g.inMenu = true
		g.menu.SetMainMenu()
	}

	// Handle mouse input
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	// Mining with mouse click
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.handleMining()
	}

	// Block placement with right click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.handleBlockPlacement()
	}

	// Hotbar selection
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.inventory.SelectSlot(0)
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.inventory.SelectSlot(1)
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.inventory.SelectSlot(2)
	} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.inventory.SelectSlot(3)
	} else if inpututil.IsKeyJustPressed(ebiten.Key5) {
		g.inventory.SelectSlot(4)
	} else if inpututil.IsKeyJustPressed(ebiten.Key6) {
		g.inventory.SelectSlot(5)
	} else if inpututil.IsKeyJustPressed(ebiten.Key7) {
		g.inventory.SelectSlot(6)
	} else if inpututil.IsKeyJustPressed(ebiten.Key8) {
		g.inventory.SelectSlot(7)
	} else if inpututil.IsKeyJustPressed(ebiten.Key9) {
		g.inventory.SelectSlot(8)
	}

	// Scroll wheel for hotbar
	_, scrollY := ebiten.Wheel()
	if scrollY > 0 {
		g.inventory.PrevSlot()
	} else if scrollY < 0 {
		g.inventory.NextSlot()
	}
}

// handleMining handles block mining
func (g *Game) handleMining() {
	// Convert mouse position to world coordinates
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY

	// Find the hexagon at mouse position directly using world coordinates
	targetHex := g.world.GetHexagonAt(mouseWorldX, mouseWorldY)
	if targetHex == nil {
		return
	}

	// Check if player can reach the block
	if !g.player.CanReach(targetHex.X, targetHex.Y) {
		return
	}

	// Damage the block
	targetHex.TakeDamage(5.0)

	// Check if block is destroyed
	if targetHex.Health <= 0 {
		// Get the exact world position before removing
		x, y := targetHex.X, targetHex.Y
		g.world.RemoveHexagonAt(x, y)
		// Use item durability
		g.inventory.UseItem()
	}
}

// handleBlockPlacement handles block placement
func (g *Game) handleBlockPlacement() {
	// Get selected item
	selectedItem := g.inventory.GetSelectedItem()
	if selectedItem == nil || selectedItem.Type == items.NONE {
		return
	}

	// Get item properties
	props := items.GetItemProperties(selectedItem.Type)
	if props == nil || !props.IsPlaceable {
		return
	}

	// Convert mouse position to world coordinates
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY

	// Use the correct world position for placement
	centerX, centerY, _, _ := world.PixelToHexCenter(mouseWorldX, mouseWorldY)

	// Place block at the calculated position
	blockType := stringToBlockType(props.BlockType)
	g.world.AddHexagonAt(centerX, centerY, blockType)

	// Remove item from inventory
	g.inventory.RemoveItem(1)
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	if g.inMenu {
		g.menu.Draw(screen)
		return
	}

	if g.inCrafting {
		// Draw game in background
		g.drawGameScene(screen)
		// Draw crafting UI overlay
		g.craftingUI.Draw(screen)
		return
	}

	if g.inGame {
		g.drawGameScene(screen)
		g.drawUI(screen)
		g.drawDebugInfo(screen)
	}
}

// drawGameScene draws the game world and player
func (g *Game) drawGameScene(screen *ebiten.Image) {
	// Draw background
	screen.Fill(g.colorToRGB(135, 206, 235)) // Sky blue

	// Get visible blocks
	px, py := g.player.GetCenter()
	visibleBlocks := g.world.GetVisibleBlocks(px, py)

	// Draw blocks
	for _, block := range visibleBlocks {
		g.drawBlock(screen, block)
	}

	// Draw player
	g.drawPlayer(screen)
}

// drawUI draws the user interface
func (g *Game) drawUI(screen *ebiten.Image) {
	// Draw hotbar
	hotbarWidth := 400
	hotbarHeight := 50
	hotbarX := (ScreenWidth - hotbarWidth) / 2
	hotbarY := ScreenHeight - hotbarHeight - 20

	slotWidth := hotbarWidth / 8

	for i := 0; i < 8; i++ {
		slotX := hotbarX + i*slotWidth
		slotY := hotbarY

		// Draw slot background
		bgColor := g.colorToRGB(100, 100, 100)
		if i == g.inventory.Selected {
			bgColor = g.colorToRGB(150, 150, 150)
		}

		ebitenutil.DrawRect(screen, float64(slotX), float64(slotY), float64(slotWidth-2), float64(hotbarHeight), bgColor)

		// Draw slot border
		ebitenutil.DrawRect(screen, float64(slotX), float64(slotY), float64(slotWidth-2), 2, g.colorToRGB(0, 0, 0))
		ebitenutil.DrawRect(screen, float64(slotX), float64(slotY+hotbarHeight-2), float64(slotWidth-2), 2, g.colorToRGB(0, 0, 0))
		ebitenutil.DrawRect(screen, float64(slotX), float64(slotY), 2, float64(hotbarHeight), g.colorToRGB(0, 0, 0))
		ebitenutil.DrawRect(screen, float64(slotX+slotWidth-4), float64(slotY), 2, float64(hotbarHeight), g.colorToRGB(0, 0, 0))

		// Draw item if present
		if i < len(g.inventory.Slots) {
			item := g.inventory.Slots[i]
			if item.Type != items.NONE {
				itemColor := items.ItemColorByID(item.Type)
				ebitenutil.DrawRect(screen, float64(slotX+10), float64(slotY+10), float64(slotWidth-20), float64(hotbarHeight-20),
					g.colorToRGB(int(itemColor.R), int(itemColor.G), int(itemColor.B)))

				// Draw quantity indicator
				if item.Quantity > 1 {
					ebitenutil.DrawRect(screen, float64(slotX+slotWidth-20), float64(slotY+hotbarHeight-20), 8, 8, g.colorToRGB(255, 255, 255))
				}
			}
		}
	}

	// Draw health bar
	healthBarWidth := 200
	healthBarHeight := 20
	healthBarX := 20
	healthBarY := 20

	healthRatio := g.player.Health / g.player.MaxHealth

	// Background
	ebitenutil.DrawRect(screen, float64(healthBarX), float64(healthBarY), float64(healthBarWidth), float64(healthBarHeight), g.colorToRGB(50, 50, 50))

	// Health
	ebitenutil.DrawRect(screen, float64(healthBarX), float64(healthBarY), float64(int(float64(healthBarWidth)*healthRatio)), float64(healthBarHeight), g.colorToRGB(200, 50, 50))

	// Border
	ebitenutil.DrawRect(screen, float64(healthBarX), float64(healthBarY), float64(healthBarWidth), 2, g.colorToRGB(0, 0, 0))
	ebitenutil.DrawRect(screen, float64(healthBarX), float64(healthBarY+healthBarHeight-2), float64(healthBarWidth), 2, g.colorToRGB(0, 0, 0))
	ebitenutil.DrawRect(screen, float64(healthBarX), float64(healthBarY), 2, float64(healthBarHeight), g.colorToRGB(0, 0, 0))
	ebitenutil.DrawRect(screen, float64(healthBarX+healthBarWidth-2), float64(healthBarY), 2, float64(healthBarHeight), g.colorToRGB(0, 0, 0))
}

// drawDebugInfo draws debug information
func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	px, py := g.player.GetCenter()
	vx, vy := g.player.GetVelocity()

	info := fmt.Sprintf("Pos: (%.1f, %.1f)\nVel: (%.1f, %.1f)\nFPS: %.1f\nOnGround: %v\nDelta: %.4f",
		px, py, vx, vy, ebiten.ActualFPS(), g.player.IsOnGround(), time.Since(g.lastTime).Seconds())

	ebitenutil.DebugPrint(screen, info)
}

// drawBlock draws a hexagonal block
func (g *Game) drawBlock(screen *ebiten.Image, block *world.Hexagon) {
	if block.BlockType == blocks.AIR {
		return
	}

	props := blocks.BlockDefinitions[getBlockKey(block.BlockType)]
	if props == nil {
		return
	}

	// Calculate screen position
	screenX := block.X - g.cameraX
	screenY := block.Y - g.cameraY

	// Check if block is on screen
	if screenX < -100 || screenX > ScreenWidth+100 ||
		screenY < -100 || screenY > ScreenHeight+100 {
		return
	}

	// Get hexagon corners
	corners := hexagon.GetHexCorners(screenX, screenY, world.HexSize)

	// Create polygon vertices
	vertices := make([]ebiten.Vertex, len(corners))
	for i, corner := range corners {
		vertices[i] = ebiten.Vertex{
			DstX:   float32(corner[0]),
			DstY:   float32(corner[1]),
			ColorR: float32(props.Color.R) / 255.0,
			ColorG: float32(props.Color.G) / 255.0,
			ColorB: float32(props.Color.B) / 255.0,
			ColorA: float32(props.Color.A) / 255.0,
		}
	}

	// Draw filled hexagon
	indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5}

	// Check if mouse is hovering over this block
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY

	q, r := hexagon.PixelToHex(mouseWorldX, mouseWorldY, world.HexSize)
	hoverHex := hexagon.HexRound(q, r)

	// Get block's hex coordinates
	blockQ, blockR := hexagon.PixelToHex(block.X, block.Y, world.HexSize)
	blockHex := hexagon.HexRound(blockQ, blockR)

	if hoverHex.Q == blockHex.Q && hoverHex.R == blockHex.R {
		// Highlight hovered block
		for i := range vertices {
			vertices[i].ColorR = minFloat32(1.0, vertices[i].ColorR+0.2)
			vertices[i].ColorG = minFloat32(1.0, vertices[i].ColorG+0.2)
			vertices[i].ColorB = minFloat32(1.0, vertices[i].ColorB+0.2)
		}
	}

	// Apply damage darkening
	if block.Health < block.MaxHealth {
		damageRatio := block.Health / block.MaxHealth
		for i := range vertices {
			vertices[i].ColorR *= float32(math.Min(1.0, float64(damageRatio)))
			vertices[i].ColorG *= float32(math.Min(1.0, float64(damageRatio)))
			vertices[i].ColorB *= float32(math.Min(1.0, float64(damageRatio)))
		}
	}

	screen.DrawTriangles(vertices, indices, g.whiteImage, nil)
}

// drawPlayer draws the player
func (g *Game) drawPlayer(screen *ebiten.Image) {
	screenX := g.player.X - g.cameraX
	screenY := g.player.Y - g.cameraY

	// Draw player body
	ebitenutil.DrawRect(screen, screenX, screenY, float64(g.player.Width), float64(g.player.Height), g.colorToRGB(255, 100, 100))

	// Draw player head (simple representation)
	ebitenutil.DrawRect(screen, screenX+10, screenY+5, 20, 20, g.colorToRGB(255, 200, 150))
}

// Layout defines the game's layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// colorToRGB converts RGB values to a color
func (g *Game) colorToRGB(rVal, gVal, bVal int) color.RGBA {
	return color.RGBA{
		R: uint8(rVal),
		G: uint8(gVal),
		B: uint8(bVal),
		A: 255,
	}
}

func min(a, b float64) float64 {
	if a < b {
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

// getBlockKey is an alias for getBlockKeyFromType
func getBlockKey(blockType blocks.BlockType) string {
	return getBlockKeyFromType(blockType)
}

// openURL opens a URL in the default browser
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// Main function
func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Tesselbox v2.0 - Hexagon Sandbox")
	ebiten.SetTPS(FPS)

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
