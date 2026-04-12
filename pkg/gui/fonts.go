package gui

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// FontManager handles loading and rendering of bitmap fonts
type FontManager struct {
	// Fallback to basic font if custom font fails
	normalFace font.Face
	smallFace  font.Face
	largeFace  font.Face
	headerFace font.Face
}

// NewFontManager creates a font manager with embedded pixel fonts
func NewFontManager() *FontManager {
	// Use basic font as fallback - in production, embed actual pixel fonts
	return &FontManager{
		normalFace: basicfont.Face7x13,
		smallFace:  basicfont.Face7x13,
		largeFace:  basicfont.Face7x13,
		headerFace: basicfont.Face7x13,
	}
}

// DrawText renders text at the specified position with color
func (fm *FontManager) DrawText(screen *ebiten.Image, str string, x, y float64, clr color.Color, size string) {
	var face font.Face
	switch size {
	case "small":
		face = fm.smallFace
	case "normal":
		face = fm.normalFace
	case "large":
		face = fm.largeFace
	case "header":
		face = fm.headerFace
	default:
		face = fm.normalFace
	}

	// For header size, draw text larger by repeating it with offset for pixel-art effect
	if size == "header" {
		// Draw shadow first
		text.Draw(screen, str, face, int(x)+2, int(y)+2, ColorShadow)
		// Draw main text 3x for bold effect
		text.Draw(screen, str, face, int(x), int(y), clr)
		text.Draw(screen, str, face, int(x)+1, int(y), clr)
		text.Draw(screen, str, face, int(x), int(y)+1, clr)
	} else {
		text.Draw(screen, str, face, int(x), int(y), clr)
	}
}

// DrawTextCentered renders text centered at the specified position
func (fm *FontManager) DrawTextCentered(screen *ebiten.Image, str string, x, y float64, clr color.Color, size string) {
	var face font.Face
	switch size {
	case "small":
		face = fm.smallFace
	case "normal":
		face = fm.normalFace
	case "large":
		face = fm.largeFace
	case "header":
		face = fm.headerFace
	default:
		face = fm.normalFace
	}

	// Calculate text width to center
	width := font.MeasureString(face, str).Round()
	xPos := int(x) - width/2
	yPos := int(y)

	// For header size, draw text larger for visibility
	if size == "header" {
		// Draw shadow first
		text.Draw(screen, str, face, xPos+2, yPos+2, ColorShadow)
		// Draw main text multiple times for bold effect
		text.Draw(screen, str, face, xPos, yPos, clr)
		text.Draw(screen, str, face, xPos+1, yPos, clr)
		text.Draw(screen, str, face, xPos, yPos+1, clr)
	} else {
		text.Draw(screen, str, face, xPos, yPos, clr)
	}
}

// DrawTextRight renders text right-aligned at the specified position
func (fm *FontManager) DrawTextRight(screen *ebiten.Image, str string, x, y float64, clr color.Color, size string) {
	var face font.Face
	switch size {
	case "small":
		face = fm.smallFace
	case "normal":
		face = fm.normalFace
	case "large":
		face = fm.largeFace
	case "header":
		face = fm.headerFace
	default:
		face = fm.normalFace
	}

	// Calculate text width for right alignment
	width := font.MeasureString(face, str).Round()

	text.Draw(screen, str, face, int(x)-width, int(y), clr)
}

// GetTextSize returns the dimensions of text
func (fm *FontManager) GetTextSize(str string, size string) (width, height float64) {
	var face font.Face
	switch size {
	case "small":
		face = fm.smallFace
	case "normal":
		face = fm.normalFace
	case "large":
		face = fm.largeFace
	case "header":
		face = fm.headerFace
	default:
		face = fm.normalFace
	}

	bounds := font.MeasureString(face, str)
	return float64(bounds.Round()), float64(face.Metrics().Height.Round())
}

// DrawPixelText renders text with pixel-perfect positioning
// Uses integer coordinates to ensure crisp rendering
func (fm *FontManager) DrawPixelText(screen *ebiten.Image, str string, x, y int, clr color.Color, size string) {
	fm.DrawText(screen, str, float64(x), float64(y), clr, size)
}

// DrawPixelTextCentered renders centered text with pixel-perfect positioning
func (fm *FontManager) DrawPixelTextCentered(screen *ebiten.Image, str string, x, y int, clr color.Color, size string) {
	fm.DrawTextCentered(screen, str, float64(x), float64(y), clr, size)
}

// DrawShadowedText renders text with a drop shadow for better visibility
func (fm *FontManager) DrawShadowedText(screen *ebiten.Image, str string, x, y float64, textClr, shadowClr color.Color, size string) {
	// Draw shadow offset by 2 pixels
	fm.DrawText(screen, str, x+2, y+2, shadowClr, size)
	// Draw main text
	fm.DrawText(screen, str, x, y, textClr, size)
}

// DrawShadowedTextCentered renders centered text with drop shadow
func (fm *FontManager) DrawShadowedTextCentered(screen *ebiten.Image, str string, x, y float64, textClr, shadowClr color.Color, size string) {
	var face font.Face
	switch size {
	case "small":
		face = fm.smallFace
	case "normal":
		face = fm.normalFace
	case "large":
		face = fm.largeFace
	case "header":
		face = fm.headerFace
	default:
		face = fm.normalFace
	}

	width := font.MeasureString(face, str).Round()
	xInt := int(x) - width/2
	yInt := int(y)

	// Draw shadow
	text.Draw(screen, str, face, xInt+2, yInt+2, shadowClr)

	// Draw main text
	text.Draw(screen, str, face, xInt, yInt, textClr)
}

// LogFontStatus logs the current font status (for debugging)
func (fm *FontManager) LogFontStatus() {
	log.Println("FontManager initialized with basic fallback font")
	log.Println("  - Small size:", DefaultFontSizes.Small)
	log.Println("  - Normal size:", DefaultFontSizes.Normal)
	log.Println("  - Large size:", DefaultFontSizes.Large)
	log.Println("  - Header size:", DefaultFontSizes.Header)
}
