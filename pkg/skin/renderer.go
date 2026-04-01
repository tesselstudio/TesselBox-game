package skin

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// PlayerRenderer renders the player with custom skin
type PlayerRenderer struct {
	skinData *SkinData
	scale    float64
	flipX    bool
}

// NewPlayerRenderer creates a new player renderer
func NewPlayerRenderer(skinData *SkinData) *PlayerRenderer {
	return &PlayerRenderer{
		skinData: skinData,
		scale:    1.0,
		flipX:    false,
	}
}

// SetScale sets the rendering scale
func (pr *PlayerRenderer) SetScale(scale float64) {
	pr.scale = scale
}

// SetFlipX sets horizontal flip
func (pr *PlayerRenderer) SetFlipX(flip bool) {
	pr.flipX = flip
}

// SetSkinData sets the skin data
func (pr *PlayerRenderer) SetSkinData(skinData *SkinData) {
	pr.skinData = skinData
}

// Draw renders the player at the specified position
func (pr *PlayerRenderer) Draw(screen *ebiten.Image, x, y float64) {
	if pr.skinData == nil {
		return
	}

	// Draw player body parts based on skin data
	pr.drawBody(screen, x, y)
	pr.drawHead(screen, x, y)
	pr.drawArms(screen, x, y)
	pr.drawLegs(screen, x, y)
}

// drawBody draws the player's body
func (pr *PlayerRenderer) drawBody(screen *ebiten.Image, x, y float64) {
	bodyWidth := 24.0 * pr.scale
	bodyHeight := 32.0 * pr.scale
	bodyX := x - bodyWidth/2
	bodyY := y - bodyHeight/2

	// Extract body region from skin (16:32 to 48:64 in skin coordinates)
	bodyRegion := pr.extractSkinRegion(16, 32, 32, 48)
	
	// Draw body
	pr.drawSkinRegion(screen, bodyRegion, bodyX, bodyY, bodyWidth, bodyHeight)
}

// drawHead draws the player's head
func (pr *PlayerRenderer) drawHead(screen *ebiten.Image, x, y float64) {
	headSize := 16.0 * pr.scale
	headX := x - headSize/2
	headY := y - bodyHeight - headSize/2 - 2*pr.scale // 2 pixels above body

	// Extract head region from skin (8:8 to 24:24 in skin coordinates)
	headRegion := pr.extractSkinRegion(8, 8, 24, 24)
	
	// Draw head
	pr.drawSkinRegion(screen, headRegion, headX, headY, headSize, headSize)
}

// drawArms draws the player's arms
func (pr *PlayerRenderer) drawArms(screen *ebiten.Image, x, y float64) {
	armWidth := 8.0 * pr.scale
	armHeight := 32.0 * pr.scale
	
	// Left arm
	leftArmX := x - bodyWidth/2 - armWidth - 2*pr.scale
	leftArmY := y - armHeight/2
	leftArmRegion := pr.extractSkinRegion(44, 32, 48, 64)
	pr.drawSkinRegion(screen, leftArmRegion, leftArmX, leftArmY, armWidth, armHeight)
	
	// Right arm
	rightArmX := x + bodyWidth/2 + 2*pr.scale
	rightArmY := y - armHeight/2
	rightArmRegion := pr.extractSkinRegion(52, 32, 56, 64)
	pr.drawSkinRegion(screen, rightArmRegion, rightArmX, rightArmY, armWidth, armHeight)
}

// drawLegs draws the player's legs
func (pr *PlayerRenderer) drawLegs(screen *ebiten.Image, x, y float64) {
	legWidth := 8.0 * pr.scale
	legHeight := 32.0 * pr.scale
	
	// Left leg
	leftLegX := x - legWidth - 2*pr.scale
	leftLegY := y + bodyHeight/2 - legHeight/2
	leftLegRegion := pr.extractSkinRegion(20, 52, 28, 64)
	pr.drawSkinRegion(screen, leftLegRegion, leftLegX, leftLegY, legWidth, legHeight)
	
	// Right leg
	rightLegX := x + 2*pr.scale
	rightLegY := y + bodyHeight/2 - legHeight/2
	rightLegRegion := pr.extractSkinRegion(36, 52, 44, 64)
	pr.drawSkinRegion(screen, rightLegRegion, rightLegX, rightLegY, legWidth, legHeight)
}

// extractSkinRegion extracts a region from the skin data
func (pr *PlayerRenderer) extractSkinRegion(x1, y1, x2, y2 int) [][]color.RGBA {
	width := x2 - x1
	height := y2 - y1
	
	region := make([][]color.RGBA, height)
	for y := 0; y < height; y++ {
		region[y] = make([]color.RGBA, width)
		for x := 0; x < width; x++ {
			skinX := x1 + x
			skinY := y1 + y
			
			if skinX >= 0 && skinX < pr.skinData.Width && 
			   skinY >= 0 && skinY < pr.skinData.Height {
				region[y][x] = pr.skinData.Pixels[skinY][skinX]
			} else {
				region[y][x] = color.RGBA{0, 0, 0, 0} // Transparent
			}
		}
	}
	
	return region
}

// drawSkinRegion draws a skin region at the specified position
func (pr *PlayerRenderer) drawSkinRegion(screen *ebiten.Image, region [][]color.RGBA, x, y, width, height float64) {
	regionHeight := len(region)
	if regionHeight == 0 {
		return
	}
	regionWidth := len(region[0])
	
	pixelWidth := width / float64(regionWidth)
	pixelHeight := height / float64(regionHeight)
	
	for py := 0; py < regionHeight; py++ {
		for px := 0; px < regionWidth; px++ {
			pixelColor := region[py][px]
			if pixelColor.A == 0 {
				continue // Skip transparent pixels
			}
			
			pixelX := x + float64(px)*pixelWidth
			pixelY := y + float64(py)*pixelHeight
			
			if pr.flipX {
				pixelX = x + width - float64(px)*pixelWidth - pixelWidth
			}
			
			ebitenutil.DrawRect(screen, pixelX, pixelY, pixelWidth, pixelHeight, pixelColor)
		}
	}
}

// Constants for body parts
const (
	bodyWidth  = 24.0
	bodyHeight = 32.0
	headSize   = 16.0
	armWidth   = 8.0
	armHeight  = 32.0
	legWidth   = 8.0
	legHeight  = 32.0
)

// DrawPlayerIcon draws a small player icon for UI
func DrawPlayerIcon(screen *ebiten.Image, skinData *SkinData, x, y, size float64) {
	if skinData == nil {
		return
	}
	
	renderer := NewPlayerRenderer(skinData)
	renderer.SetScale(size / 64.0) // Scale based on 64px base size
	renderer.Draw(screen, x+size/2, y+size/2)
}

// DrawPlayerPreview draws a player preview with animation
func (pr *PlayerRenderer) DrawPlayerPreview(screen *ebiten.Image, x, y float64, time float64) {
	// Add simple animation (bobbing up and down)
	bobAmount := math.Sin(time * 2) * 2 * pr.scale
	y += bobAmount
	
	pr.Draw(screen, x, y)
}
