package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/projet-red_rap-legacy/game" // <-- ici le package game
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("Rap Legacy")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
