package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Map représente la map du jeu
type Map struct {
	bgImage *ebiten.Image
}

// NewMap charge l'image de fond
func NewMap() *Map {
	// Utilise la fonction LoadImage globale
	return &Map{
		bgImage: LoadImage("assets/image3.png"),
	}
}

// Draw affiche l'image de fond
func (m *Map) Draw(screen *ebiten.Image) {
	if m.bgImage != nil {
		opts := &ebiten.DrawImageOptions{}
		screen.DrawImage(m.bgImage, opts)
	}
}
