package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Player struct {
	X, Y     float64
	Ego      int
	Flow     int
	Charisma int
}

func NewPlayer(x, y float64) *Player {
	return &Player{X: x, Y: y, Ego: 100, Flow: 10, Charisma: 5}
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.X += 2
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, p.X, p.Y, 32, 32, color.RGBA{0, 255, 0, 255})
}
