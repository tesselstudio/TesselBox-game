package gui

import (
	"image/color"
)

// Pixel Art Color Palette - Retro Gaming Aesthetic
var (
	// Background colors
	ColorBackground    = color.RGBA{0x1a, 0x1c, 0x2c, 0xff} // Dark blue-black
	ColorBackgroundDark = color.RGBA{0x10, 0x12, 0x20, 0xff} // Darker shade
	
	// Primary colors
	ColorPrimary       = color.RGBA{0x57, 0xff, 0x57, 0xff} // Bright green (selection)
	ColorPrimaryDark   = color.RGBA{0x3d, 0xb8, 0x3d, 0xff} // Darker green
	ColorSecondary     = color.RGBA{0xff, 0x8c, 0x42, 0xff} // Orange (accents)
	ColorSecondaryDark = color.RGBA{0xd4, 0x6a, 0x28, 0xff} // Darker orange
	
	// Status colors
	ColorDanger  = color.RGBA{0xff, 0x44, 0x44, 0xff} // Red (delete/warning)
	ColorWarning = color.RGBA{0xff, 0xaa, 0x00, 0xff} // Amber
	ColorInfo    = color.RGBA{0x42, 0xa5, 0xff, 0xff} // Blue
	ColorSuccess = color.RGBA{0x42, 0xff, 0x8c, 0xff} // Mint green
	
	// Text colors
	ColorText       = color.RGBA{0xf4, 0xf4, 0xf4, 0xff} // Off-white
	ColorTextDim    = color.RGBA{0x9c, 0xa3, 0xaf, 0xff} // Gray (inactive)
	ColorTextDark   = color.RGBA{0x6b, 0x72, 0x80, 0xff} // Dark gray
	
	// UI element colors
	ColorPanel      = color.RGBA{0x25, 0x28, 0x3d, 0xff} // Panel background
	ColorPanelLight = color.RGBA{0x35, 0x38, 0x4d, 0xff} // Lighter panel
	ColorBorder     = color.RGBA{0x4a, 0x4d, 0x5e, 0xff} // Border color
	ColorBorderDark = color.RGBA{0x1a, 0x1c, 0x2c, 0xff} // Dark border (for 3D effect)
	ColorHighlight  = color.RGBA{0x6a, 0x6d, 0x7e, 0xff} // Highlight (3D effect top)
	ColorShadow     = color.RGBA{0x0a, 0x0c, 0x18, 0xff} // Shadow (3D effect bottom)
)

// UI Dimensions - Pixel-perfect sizing
type UIDimensions struct {
	ButtonWidth      int
	ButtonHeight     int
	ButtonPadding    int
	PanelPadding     int
	BorderWidth      int
	CornerRadius     int
	TextPadding      int
	LineHeight       int
	ScrollBarWidth   int
	IconSize         int
}

// DefaultDims provides standard pixel art UI dimensions
var DefaultDims = UIDimensions{
	ButtonWidth:    200,
	ButtonHeight:   40,
	ButtonPadding:  10,
	PanelPadding:   20,
	BorderWidth:    2,
	CornerRadius:   4,
	TextPadding:    8,
	LineHeight:     24,
	ScrollBarWidth: 12,
	IconSize:       16,
}

// FontSizes - Bitmap font sizes
type FontSizes struct {
	Small  int // 8px
	Normal int // 16px
	Large  int // 24px
	Header int // 32px
}

// DefaultFontSizes provides standard font sizes
var DefaultFontSizes = FontSizes{
	Small:  8,
	Normal: 16,
	Large:  24,
	Header: 32,
}

// Animation timings
type AnimationTimings struct {
	ButtonPress   float64 // Duration of button press animation
	HoverFade     float64 // Duration of hover highlight fade
	MenuTransition float64 // Duration between menu screens
	CursorBlink   float64 // Cursor blink rate
}

// DefaultAnimationTimings provides standard animation speeds
var DefaultAnimationTimings = AnimationTimings{
	ButtonPress:    0.1,
	HoverFade:      0.15,
	MenuTransition: 0.3,
	CursorBlink:    0.5,
}

// Screen identifiers
const (
	ScreenMain        = "main"
	ScreenWorlds      = "worlds"
	ScreenCreateWorld = "create_world"
	ScreenDeleteWorld = "delete_world"
	ScreenSettings    = "settings"
	ScreenSkinEditor  = "skin_editor"
	ScreenPlugins     = "plugins"
)

// GetShadowColor returns shadow color with alpha based on depth
func GetShadowColor(depth int) color.RGBA {
	alpha := uint8(255)
	if depth > 0 {
		alpha = uint8(max(50, 255-depth*30))
	}
	return color.RGBA{0x0a, 0x0c, 0x18, alpha}
}

// max helper function
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
