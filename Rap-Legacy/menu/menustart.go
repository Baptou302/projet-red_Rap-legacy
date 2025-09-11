package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	menuOption int
	startGame  bool
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if g.menuOption > 0 {
			g.menuOption--
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if g.menuOption < 2 {
			g.menuOption++
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if g.menuOption == 0 {
			g.startGame = true
		} else if g.menuOption == 1 {
			log.Println("Options sélectionnées")
		} else if g.menuOption == 2 {
			log.Println("Quitter")
			return ebiten.Termination
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	if g.startGame {
		ebitenutil.DebugPrint(screen, "Le jeu commence !")
		return
	}

	options := []string{"Jouer", "Options", "Quitter"}
	for i, option := range options {
		text := option
		if i == g.menuOption {
			text = "> " + option
		}
		ebitenutil.DebugPrintAt(screen, text, 50, 50+i*30)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Menu de démarrage du jeu")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
