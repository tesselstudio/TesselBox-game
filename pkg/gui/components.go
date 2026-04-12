package gui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ButtonState represents the current state of a button
type ButtonState int

const (
	ButtonStateNormal ButtonState = iota
	ButtonStateHover
	ButtonStatePressed
	ButtonStateDisabled
)

// PixelButton represents a 3D pixel art button
type PixelButton struct {
	X, Y          float64
	Width, Height float64
	Text          string
	State         ButtonState
	OnClick       func()

	// Visual properties
	BgColor     color.Color
	TextColor   color.Color
	BorderColor color.Color
}

// NewPixelButton creates a new pixel art button
func NewPixelButton(x, y, width, height float64, text string, onClick func()) *PixelButton {
	return &PixelButton{
		X:           x,
		Y:           y,
		Width:       width,
		Height:      height,
		Text:        text,
		State:       ButtonStateNormal,
		OnClick:     onClick,
		BgColor:     ColorPanel,
		TextColor:   ColorText,
		BorderColor: ColorBorder,
	}
}

// Update checks for hover/click state
func (b *PixelButton) Update() {
	if b.State == ButtonStateDisabled {
		return
	}

	mx, my := ebiten.CursorPosition()

	// Check if mouse is over button
	if float64(mx) >= b.X && float64(mx) <= b.X+b.Width &&
		float64(my) >= b.Y && float64(my) <= b.Y+b.Height {

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b.State = ButtonStatePressed
		} else {
			if b.State == ButtonStatePressed {
				// Button was just released - trigger click
				if b.OnClick != nil {
					b.OnClick()
				}
			}
			b.State = ButtonStateHover
		}
	} else {
		b.State = ButtonStateNormal
	}
}

// Draw renders the button with 3D pixel art style
func (b *PixelButton) Draw(screen *ebiten.Image, fm *FontManager) {
	// Determine colors based on state
	bgColor := b.BgColor
	textColor := b.TextColor
	topHighlight := ColorHighlight
	bottomShadow := ColorShadow
	offset := 0

	switch b.State {
	case ButtonStateHover:
		bgColor = ColorPanelLight
		topHighlight = ColorTextDim
	case ButtonStatePressed:
		bgColor = ColorPanelLight
		topHighlight = ColorShadow
		bottomShadow = ColorHighlight
		offset = 2 // Button moves down when pressed
	case ButtonStateDisabled:
		bgColor = ColorBackgroundDark
		textColor = ColorTextDark
	}

	x := b.X + float64(offset)
	y := b.Y + float64(offset)
	w := b.Width
	h := b.Height

	// Draw shadow (bottom/right) - 3D effect
	vector.DrawFilledRect(screen, float32(x+2), float32(y+2), float32(w), float32(h), bottomShadow, false)

	// Draw top highlight (top/left) - 3D effect
	vector.DrawFilledRect(screen, float32(x-1), float32(y-1), float32(w+1), float32(h+1), topHighlight, false)

	// Draw main button body
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), bgColor, false)

	// Draw border
	vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 1, b.BorderColor, false)

	// Draw text centered
	if fm != nil && b.Text != "" {
		textX := x + w/2
		textY := y + h/2 - 5 // Slight offset for vertical centering
		fm.DrawTextCentered(screen, b.Text, textX, textY, textColor, "normal")
	}
}

// IsHovered returns true if the mouse is over the button
func (b *PixelButton) IsHovered() bool {
	mx, my := ebiten.CursorPosition()
	return float64(mx) >= b.X && float64(mx) <= b.X+b.Width &&
		float64(my) >= b.Y && float64(my) <= b.Y+b.Height
}

// SetPosition updates the button position
func (b *PixelButton) SetPosition(x, y float64) {
	b.X = x
	b.Y = y
}

// PixelPanel represents a bordered container with shadow
type PixelPanel struct {
	X, Y          float64
	Width, Height float64
	Title         string
	BgColor       color.Color
	BorderColor   color.Color
}

// NewPixelPanel creates a new pixel art panel
func NewPixelPanel(x, y, width, height float64, title string) *PixelPanel {
	return &PixelPanel{
		X:           x,
		Y:           y,
		Width:       width,
		Height:      height,
		Title:       title,
		BgColor:     ColorPanel,
		BorderColor: ColorBorder,
	}
}

// Draw renders the panel
func (p *PixelPanel) Draw(screen *ebiten.Image, fm *FontManager) {
	x, y, w, h := float32(p.X), float32(p.Y), float32(p.Width), float32(p.Height)

	// Draw shadow
	vector.DrawFilledRect(screen, x+2, y+2, w, h, ColorShadow, false)

	// Draw main panel
	vector.DrawFilledRect(screen, x, y, w, h, p.BgColor, false)

	// Draw border
	vector.StrokeRect(screen, x, y, w, h, 1, p.BorderColor, false)

	// Draw title if provided
	if fm != nil && p.Title != "" {
		titleY := y + 15
		fm.DrawTextCentered(screen, p.Title, float64(x+w/2), float64(titleY), ColorText, "normal")

		// Draw title underline
		vector.StrokeLine(screen, x+10, titleY+10, x+w-10, titleY+10, 1, ColorBorder, false)
	}
}

// ContainsPoint checks if a point is inside the panel
func (p *PixelPanel) ContainsPoint(x, y float64) bool {
	return x >= p.X && x <= p.X+p.Width && y >= p.Y && y <= p.Y+p.Height
}

// PixelList represents a scrollable list of items
type PixelList struct {
	X, Y          float64
	Width, Height float64
	Items         []string
	SelectedIndex int
	HoverIndex    int
	ScrollOffset  int
	ItemHeight    float64
	OnSelect      func(index int)
}

// NewPixelList creates a new scrollable list
func NewPixelList(x, y, width, height float64, items []string, onSelect func(int)) *PixelList {
	return &PixelList{
		X:             x,
		Y:             y,
		Width:         width,
		Height:        height,
		Items:         items,
		SelectedIndex: -1,
		HoverIndex:    -1,
		ScrollOffset:  0,
		ItemHeight:    30,
		OnSelect:      onSelect,
	}
}

// Update handles input for the list
func (l *PixelList) Update() {
	mx, my := ebiten.CursorPosition()

	// Check if mouse is over list
	if float64(mx) >= l.X && float64(mx) <= l.X+l.Width &&
		float64(my) >= l.Y && float64(my) <= l.Y+l.Height {

		// Calculate which item is hovered
		relativeY := float64(my) - l.Y
		itemIndex := int(relativeY/l.ItemHeight) + l.ScrollOffset

		if itemIndex >= 0 && itemIndex < len(l.Items) {
			l.HoverIndex = itemIndex

			// Handle click
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				l.SelectedIndex = itemIndex
				if l.OnSelect != nil {
					l.OnSelect(itemIndex)
				}
			}
		} else {
			l.HoverIndex = -1
		}
	} else {
		l.HoverIndex = -1
	}

	// Handle scrolling
	_, scrollY := ebiten.Wheel()
	if scrollY > 0 {
		l.ScrollOffset = max(0, l.ScrollOffset-1)
	} else if scrollY < 0 {
		maxOffset := max(0, len(l.Items)-int(l.Height/l.ItemHeight))
		l.ScrollOffset = min(maxOffset, l.ScrollOffset+1)
	}
}

// Draw renders the list
func (l *PixelList) Draw(screen *ebiten.Image, fm *FontManager) {
	// Draw panel background
	vector.DrawFilledRect(screen, float32(l.X), float32(l.Y), float32(l.Width), float32(l.Height), ColorPanel, false)
	vector.StrokeRect(screen, float32(l.X), float32(l.Y), float32(l.Width), float32(l.Height), 1, ColorBorder, false)

	// Calculate visible range
	visibleCount := int(l.Height / l.ItemHeight)
	startIdx := l.ScrollOffset
	endIdx := min(len(l.Items), startIdx+visibleCount)

	// Draw items
	for i := startIdx; i < endIdx; i++ {
		itemY := l.Y + float64(i-startIdx)*l.ItemHeight

		// Determine item background color
		bgColor := ColorPanel
		if i == l.SelectedIndex {
			bgColor = ColorPrimary
		} else if i == l.HoverIndex {
			bgColor = ColorPanelLight
		}

		// Draw item background
		vector.DrawFilledRect(screen, float32(l.X), float32(itemY), float32(l.Width), float32(l.ItemHeight-1), bgColor, false)

		// Draw item text
		if fm != nil {
			textColor := ColorText
			if i == l.SelectedIndex {
				textColor = ColorBackground
			}
			textX := l.X + 10
			textY := itemY + l.ItemHeight/2 - 5
			fm.DrawText(screen, l.Items[i], textX, textY, textColor, "normal")
		}
	}

	// Draw scrollbar if needed
	if len(l.Items) > visibleCount {
		scrollbarX := l.X + l.Width - float64(DefaultDims.ScrollBarWidth)
		scrollbarHeight := l.Height * (float64(visibleCount) / float64(len(l.Items)))
		scrollbarY := l.Y + (l.Height-scrollbarHeight)*(float64(l.ScrollOffset)/float64(len(l.Items)-visibleCount))

		vector.DrawFilledRect(screen, float32(scrollbarX), float32(l.Y), float32(DefaultDims.ScrollBarWidth), float32(l.Height), ColorBackgroundDark, false)
		vector.DrawFilledRect(screen, float32(scrollbarX+2), float32(scrollbarY), float32(DefaultDims.ScrollBarWidth-4), float32(scrollbarHeight), ColorBorder, false)
	}
}

// PixelToggle represents a toggle switch
type PixelToggle struct {
	X, Y     float64
	Width    float64
	Label    string
	Checked  bool
	OnToggle func(bool)
}

// NewPixelToggle creates a new toggle switch
func NewPixelToggle(x, y float64, label string, checked bool, onToggle func(bool)) *PixelToggle {
	return &PixelToggle{
		X:        x,
		Y:        y,
		Width:    200,
		Label:    label,
		Checked:  checked,
		OnToggle: onToggle,
	}
}

// Update handles input
func (t *PixelToggle) Update() {
	mx, my := ebiten.CursorPosition()

	// Click area includes label + toggle box
	if float64(mx) >= t.X && float64(mx) <= t.X+t.Width &&
		float64(my) >= t.Y && float64(my) <= t.Y+25 {

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			t.Checked = !t.Checked
			if t.OnToggle != nil {
				t.OnToggle(t.Checked)
			}
		}
	}
}

// Draw renders the toggle
func (t *PixelToggle) Draw(screen *ebiten.Image, fm *FontManager) {
	boxSize := float32(16)
	x := float32(t.X)
	y := float32(t.Y)

	// Draw label
	if fm != nil {
		fm.DrawText(screen, t.Label, t.X+25, t.Y+3, ColorText, "normal")
	}

	// Draw toggle box
	if t.Checked {
		vector.DrawFilledRect(screen, x, y, boxSize, boxSize, ColorPrimary, false)
	} else {
		vector.DrawFilledRect(screen, x, y, boxSize, boxSize, ColorPanel, false)
	}

	vector.StrokeRect(screen, x, y, boxSize, boxSize, 1, ColorBorder, false)

	// Draw checkmark if checked
	if t.Checked {
		checkColor := ColorBackground
		// Simple checkmark: two lines
		vector.StrokeLine(screen, x+3, y+8, x+7, y+12, 2, checkColor, false)
		vector.StrokeLine(screen, x+7, y+12, x+13, y+4, 2, checkColor, false)
	}
}

// minFloat64 returns the minimum of two float64 values
func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// maxFloat64 returns the maximum of two float64 values
func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// DrawDebugGrid draws a pixel grid for debugging alignment
func DrawDebugGrid(screen *ebiten.Image, spacing int) {
	w, h := screen.Size()
	gridColor := color.RGBA{255, 0, 0, 50}

	for x := 0; x < w; x += spacing {
		vector.StrokeLine(screen, float32(x), 0, float32(x), float32(h), 1, gridColor, false)
	}
	for y := 0; y < h; y += spacing {
		vector.StrokeLine(screen, 0, float32(y), float32(w), float32(y), 1, gridColor, false)
	}
}

// DrawPixelText is a convenience function that doesn't require FontManager
func DrawPixelText(screen *ebiten.Image, text string, x, y int, clr color.Color) {
	ebitenutil.DebugPrintAt(screen, text, x, y)
}
