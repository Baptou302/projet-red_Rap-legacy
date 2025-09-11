package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Enemy struct {
    X, Y float64
    Name string
    Ego int
}

func NewEnemy(x, y float64, name string) *Enemy {
    return &Enemy{X: x, Y: y, Name: name, Ego: 50}
}

func (e *Enemy) Draw(screen *ebiten.Image) {
    ebitenutil.DrawRect(screen, e.X, e.Y, 32, 32, color.RGBA{255, 0, 0, 255})
}
