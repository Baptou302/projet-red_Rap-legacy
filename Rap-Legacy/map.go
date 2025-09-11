package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Map struct{}

func NewMap() *Map {
	return &Map{}
}

func (m *Map) Draw(screen *ebiten.Image) {
	// Simple fond gris pour la ville
	ebitenutil.DrawRect(screen, 0, 0, 640, 480, color.RGBA{200, 200, 200, 255})
}
