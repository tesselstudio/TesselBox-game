//go:build !tcp
// +build !tcp

package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"tesselbox/pkg/render"
)

func main() {
	game := render.NewGame()

	ebiten.SetWindowSize(render.ScreenWidth, render.ScreenHeight)
	ebiten.SetWindowTitle("Tesselbox - Go/Ebiten Version")

	if err := ebiten.RunGame(game); err != nil {
		log.Println(err)
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
	}
}