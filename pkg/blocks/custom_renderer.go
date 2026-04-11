package blocks

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// CustomBlockRenderer handles custom block appearances
type CustomBlockRenderer struct {
	randSource *rand.Rand
}

// NewCustomBlockRenderer creates a new custom renderer
func NewCustomBlockRenderer() *CustomBlockRenderer {
	return &CustomBlockRenderer{
		randSource: rand.New(rand.NewSource(42)), // Deterministic seed
	}
}

// GenerateCustomTexture creates a custom texture for any block
func (cbr *CustomBlockRenderer) GenerateCustomTexture(blockType string, x, y int, customParams map[string]interface{}) *ebiten.Image {
	const textureSize = 64
	img := ebiten.NewImage(textureSize, textureSize)

	switch blockType {
	case "stone":
		cbr.generateStoneTexture(img, x, y, customParams)
	case "grass":
		cbr.generateGrassTexture(img, x, y, customParams)
	case "water":
		cbr.generateWaterTexture(img, x, y, customParams)
	case "custom_crystal":
		cbr.generateCrystalTexture(img, x, y, customParams)
	default:
		cbr.generateSolidTexture(img, color.RGBA{128, 128, 128, 255})
	}

	return img
}

// generateStoneTexture creates a stone texture with custom patterns
func (cbr *CustomBlockRenderer) generateStoneTexture(img *ebiten.Image, x, y int, params map[string]interface{}) {
	baseColor := color.RGBA{128, 128, 128, 255}

	// Get custom parameters
	roughness, _ := params["roughness"].(float64)
	if roughness == 0 {
		roughness = 0.3
	}

	// Generate stone texture with noise
	for px := 0; px < 64; px++ {
		for py := 0; py < 64; py++ {
			// Create noise pattern
			noise := (cbr.randSource.Float64() - 0.5) * roughness
			noise += math.Sin(float64(px+x)*0.1) * 0.1
			noise += math.Cos(float64(py+y)*0.1) * 0.1

			// Apply noise to base color
			r := float64(baseColor.R) + noise*50
			g := float64(baseColor.G) + noise*50
			b := float64(baseColor.B) + noise*50

			// Clamp values
			r = math.Max(0, math.Min(255, r))
			g = math.Max(0, math.Min(255, g))
			b = math.Max(0, math.Min(255, b))

			img.Set(px, py, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: 255,
			})
		}
	}
}

// generateGrassTexture creates grass with blade patterns
func (cbr *CustomBlockRenderer) generateGrassTexture(img *ebiten.Image, x, y int, params map[string]interface{}) {
	baseColor := color.RGBA{124, 169, 84, 255}
	darkColor := color.RGBA{94, 139, 54, 255}

	// Fill base color
	img.Fill(baseColor)

	// Add grass blades
	bladeCount := 20 + cbr.randSource.Intn(10)
	for i := 0; i < bladeCount; i++ {
		startX := cbr.randSource.Intn(64)
		startY := cbr.randSource.Intn(64)
		length := 3 + cbr.randSource.Intn(5)
		angle := cbr.randSource.Float64() * math.Pi * 2

		// Draw grass blade
		for j := 0; j < length; j++ {
			px := startX + int(math.Cos(angle)*float64(j))
			py := startY + int(math.Sin(angle)*float64(j))

			if px >= 0 && px < 64 && py >= 0 && py < 64 {
				img.Set(px, py, darkColor)
			}
		}
	}
}

// generateWaterTexture creates animated water effect
func (cbr *CustomBlockRenderer) generateWaterTexture(img *ebiten.Image, x, y int, params map[string]interface{}) {
	baseColor := color.RGBA{64, 164, 223, 180}

	// Create water with wave patterns
	for px := 0; px < 64; px++ {
		for py := 0; py < 64; py++ {
			// Wave pattern
			wave := math.Sin(float64(px+x)*0.2) * math.Cos(float64(py+y)*0.2)
			wave *= 0.3

			// Apply wave to color
			r := float64(baseColor.R) + wave*30
			g := float64(baseColor.G) + wave*30
			b := float64(baseColor.B) + wave*40

			// Clamp values
			r = math.Max(0, math.Min(255, r))
			g = math.Max(0, math.Min(255, g))
			b = math.Max(0, math.Min(255, b))

			img.Set(px, py, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: 180,
			})
		}
	}
}

// generateCrystalTexture creates a crystalline structure
func (cbr *CustomBlockRenderer) generateCrystalTexture(img *ebiten.Image, x, y int, params map[string]interface{}) {
	// Get crystal color from params or use default
	crystalColor, ok := params["color"].(color.RGBA)
	if !ok {
		crystalColor = color.RGBA{185, 242, 255, 255}
	}

	// Create crystalline pattern
	centerX, centerY := 32, 32

	for px := 0; px < 64; px++ {
		for py := 0; py < 64; py++ {
			dx := float64(px - centerX)
			dy := float64(py - centerY)
			distance := math.Sqrt(dx*dx + dy*dy)

			// Create crystalline facets
			angle := math.Atan2(dy, dx)
			facet := math.Sin(angle*6) * math.Cos(distance*0.3)

			// Apply facet shading
			brightness := 0.5 + facet*0.5
			r := float64(crystalColor.R) * brightness
			g := float64(crystalColor.G) * brightness
			b := float64(crystalColor.B) * brightness

			// Clamp values
			r = math.Max(0, math.Min(255, r))
			g = math.Max(0, math.Min(255, g))
			b = math.Max(0, math.Min(255, b))

			img.Set(px, py, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: 255,
			})
		}
	}
}

// generateSolidTexture creates a simple solid color texture
func (cbr *CustomBlockRenderer) generateSolidTexture(img *ebiten.Image, col color.RGBA) {
	img.Fill(col)
}

// Global custom renderer instance
var GlobalCustomRenderer = NewCustomBlockRenderer()
