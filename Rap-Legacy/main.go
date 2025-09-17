package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/projet-red_rap-legacy/game"
)

func main() {
	SetGameIcon("assets/icon.png")

	g := game.NewGame()
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("Rap Legacy")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
