package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Constants for mining UI
const (
	MiningBarWidth   = 40.0
	MiningBarHeight  = 6.0
	MiningBarYOffset = 50.0
)

var (
	MiningBarBackgroundColor = color.RGBA{50, 50, 50, 200}
	MiningBarFillColor       = color.RGBA{255, 200, 0, 255}
)

// renderMiningBar draws the mining progress bar
func (g *Game) renderMiningBar(screen *ebiten.Image, blockX, blockY float64, progress float64) {
	if progress <= 0 || progress >= 1 {
		return
	}

	// Convert world position to screen position
	screenX := blockX - g.cameraX - MiningBarWidth/2
	screenY := blockY - g.cameraY - MiningBarYOffset

	// Draw background
	bgRect := ebiten.NewImage(int(MiningBarWidth), int(MiningBarHeight))
	bgRect.Fill(MiningBarBackgroundColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenX, screenY)
	screen.DrawImage(bgRect, op)

	// Draw fill
	fillWidth := int(MiningBarWidth * progress)
	if fillWidth > 0 {
		fillRect := ebiten.NewImage(fillWidth, int(MiningBarHeight))
		fillRect.Fill(MiningBarFillColor)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(screenX, screenY)
		screen.DrawImage(fillRect, op)
	}
}
